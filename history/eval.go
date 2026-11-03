package history

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/snowbldr/fluent-sdfx/anchor"
	"github.com/snowbldr/fluent-sdfx/sculpt"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/sdfx/sdf"
)

// Options are a runtime's evaluation settings. Two runtimes with different
// options are two caches: a preview runtime and a final one can serve the
// same history side by side.
type Options struct {
	// BakeEvery samples a solid onto a grid after this many closed-form
	// warps (sculpt.Pull) in a chain, so a long run of pulls stays fast to
	// evaluate and mesh. 0 never bakes. A bake is a cache, never a step:
	// the surface moves by a fraction of a cell, the field keeps its
	// values (it is sampled, not redistanced), and editing anything before
	// it throws it away.
	BakeEvery int
	// BakeCells is the bake's resolution, cells along the longest side.
	// 0 means 128.
	BakeCells int
}

// Runtime evaluates histories and caches every step's result by a hash of
// what the step is and everything it depends on, so evaluating again after
// an edit recomputes only the edited step and what follows from it. Safe
// for concurrent use; evaluations take turns.
type Runtime struct {
	opts  Options
	mu    sync.Mutex
	cache map[string]*entry
	gen   int
}

// keepGens is how many evaluations an unused cache entry survives: enough
// that undoing an edit, or flipping a step off and on, finds the old
// results still there.
const keepGens = 16

type entry struct {
	val     reflect.Value
	status  Status
	msg     string
	depth   int
	baked   bool
	millis  float64
	lastGen int
	// consumes names the steps whose value this one used whole, so they
	// are not bodies of their own; inputConsumed is its implicit input.
	consumes      []string
	inputConsumed bool
}

// NewRuntime returns a runtime with empty caches.
func NewRuntime(o Options) *Runtime {
	if o.BakeCells == 0 {
		o.BakeCells = 128
	}
	return &Runtime{opts: o, cache: map[string]*entry{}}
}

// Status is a step's outcome.
type Status string

const (
	OK      Status = "ok"      // evaluated
	Failed  Status = "error"   // evaluated and failed; Error says why and usually how to fix it
	Blocked Status = "blocked" // not evaluated: something it needs failed, is off with nothing to pass on, or does not exist
	Off     Status = "off"     // skipped; its input passes through to the steps that use it
)

// warps are the closed-form ops that leave the field a scaled bound and
// slow every evaluation after them; BakeEvery counts these.
var warps = map[string]bool{"sculpt.Pull": true}

// stepInfo is what a step depends on, found from its definition alone.
type stepInfo struct {
	input  int      // the step whose value is this step's input, or -1
	deps   []int    // every earlier step it references, input included
	params []string // params it references
}

type blockedErr struct{ msg string }

func (e blockedErr) Error() string { return e.msg }

// Eval evaluates the history. It never fails as a whole: every step gets
// a status, and the steps that could be evaluated were.
func (rt *Runtime) Eval(h *History) *Result {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	rt.gen++

	sc := &scope{rt: rt, h: h, names: map[string]int{}}
	for i, s := range h.Steps {
		if s.Name != "" {
			sc.names[s.Name] = i
		}
	}
	sc.infos = make([]stepInfo, len(h.Steps))
	for i := range h.Steps {
		sc.infos[i] = sc.analyze(i)
	}
	sc.slots = make([]slot, len(h.Steps))
	sc.hashes = make([]string, len(h.Steps))
	res := &Result{Steps: make([]StepResult, len(h.Steps))}

	for i, s := range h.Steps {
		sc.i = i
		sr := StepResult{Index: i, Name: s.Name, Op: s.Op, By: s.By, From: s.From}
		if s.Anchor != nil {
			sr.Op = "anchor"
		}
		info := sc.infos[i]
		sc.hashes[i] = sc.hash(i)
		if s.Off {
			sr.Status = Off
			if in := info.input; in >= 0 && sc.slots[in].val.IsValid() {
				sc.slots[i] = slot{val: sc.slots[in].val, status: Off, depth: sc.slots[in].depth}
				sc.hashes[i] = sc.hashes[in]
			} else {
				sc.slots[i] = slot{status: Off}
			}
			sr.Kind = kindOf(sc.slots[i].val)
			res.Steps[i] = sr
			continue
		}
		e, hit := rt.cache[sc.hashes[i]]
		if !hit {
			start := time.Now()
			e = sc.run(i)
			e.millis = float64(time.Since(start).Microseconds()) / 1000
			rt.cache[sc.hashes[i]] = e
		}
		e.lastGen = rt.gen
		sc.slots[i] = slot{val: e.val, status: e.status, depth: e.depth, e: e}
		sr.Status, sr.Error, sr.Baked, sr.Millis, sr.Cached = e.status, e.msg, e.baked, e.millis, hit
		sr.Kind = kindOf(e.val)
		if e.status == OK {
			describe(&sr, e.val)
		}
		res.Steps[i] = sr
	}
	for k, e := range rt.cache {
		if e.lastGen < rt.gen-keepGens {
			delete(rt.cache, k)
		}
	}
	res.bodies(sc)
	return res
}

type slot struct {
	val    reflect.Value // invalid when the step produced nothing usable
	status Status
	depth  int    // warps since the last grid, for BakeEvery
	e      *entry // the evaluation, for what it consumed
}

type scope struct {
	rt     *Runtime
	h      *History
	names  map[string]int
	infos  []stepInfo
	slots  []slot
	hashes []string
	i      int // the step being evaluated

	input    reflect.Value // the current step's input, for inline anchors
	consumed map[string]bool
}

// prev is the nearest earlier step that is not an anchor, or -1.
func (sc *scope) prev(i int) int {
	for j := i - 1; j >= 0; j-- {
		if sc.h.Steps[j].Anchor == nil {
			return j
		}
	}
	return -1
}

// implicitInput says whether a step with no `on` takes the previous step
// as its input: always for a method, and for a function whose first
// parameter is a form (a solid, grid or shape) that args does not name.
func implicitInput(s Step) bool {
	if s.Anchor != nil || s.On != "" {
		return false
	}
	if !strings.Contains(s.Op, ".") {
		return true
	}
	c, ok := lookupFunc(s.Op)
	if !ok {
		return false
	}
	t := c.fn.Type()
	if t.NumIn() == 0 || !formType(t.In(0)) || (t.IsVariadic() && t.NumIn() == 1) {
		return false
	}
	switch a := s.Args.(type) {
	case nil:
		return true
	case map[string]any:
		_, named := lookupArg(a, c.params[0])
		return !named
	}
	return false
}

var (
	solidType = reflect.TypeOf((*solid.Solid)(nil))
	gridType  = reflect.TypeOf((*sculpt.Grid)(nil))
	shapeType = reflect.TypeOf((*shape.Shape)(nil))
	pointType = reflect.TypeOf(anchor.Point{})
	vec3Type  = reflect.TypeOf(v3.Vec{})
	vec2Type  = reflect.TypeOf(v2.Vec{})
	anchoredT = reflect.TypeOf(solid.AnchoredSolid{})
	placeType = reflect.TypeOf(solid.Placement{})
	sdf3Type  = reflect.TypeOf((*sdf.SDF3)(nil)).Elem()
	sdf2Type  = reflect.TypeOf((*sdf.SDF2)(nil)).Elem()
	floatType = reflect.TypeOf(0.0)
	vec3sType = reflect.TypeOf([]v3.Vec{})
)

func formType(t reflect.Type) bool {
	switch t {
	case solidType, gridType, shapeType, sdf3Type, sdf2Type:
		return true
	}
	return false
}

func (sc *scope) analyze(i int) stepInfo {
	s := sc.h.Steps[i]
	info := stepInfo{input: -1}
	deps, params := map[int]bool{}, map[string]bool{}
	var walk func(v any)
	ref := func(r string) {
		if strings.HasPrefix(r, "$") {
			params[r[1:]] = true
			return
		}
		base, _, _ := strings.Cut(r, ".")
		if j, ok := sc.names[base]; ok && j < i {
			deps[j] = true
		}
	}
	walkAnchor := func(a *AnchorDef) {
		if a == nil {
			return
		}
		if a.On != "" {
			ref(a.On)
		}
		walk(a.Nearest)
		walk(a.Offset)
		if a.Ray != nil {
			walk(a.Ray.From)
			walk(a.Ray.Dir)
			walk(a.Ray.Toward)
		}
	}
	walk = func(v any) {
		switch t := v.(type) {
		case string:
			ref(t)
		case map[string]any:
			if a, ok := t["anchor"]; ok && len(t) == 1 {
				if def, err := decodeAnchor(a); err == nil {
					walkAnchor(def)
				}
				return
			}
			for _, x := range t {
				walk(x)
			}
		case []any:
			for _, x := range t {
				walk(x)
			}
		}
	}
	switch {
	case s.Anchor != nil:
		walkAnchor(s.Anchor)
		if s.Anchor.On == "" {
			info.input = sc.prev(i)
		} else if j, ok := sc.names[s.Anchor.On]; ok && j < i {
			info.input = j
		}
	case s.On != "":
		ref(s.On)
		if j, ok := sc.names[strings.SplitN(s.On, ".", 2)[0]]; ok && j < i {
			info.input = j
		}
	case implicitInput(s):
		info.input = sc.prev(i)
	}
	if info.input >= 0 {
		deps[info.input] = true
	}
	walk(s.Args)
	for _, c := range s.Then {
		walk(c.Args)
	}
	for j := range deps {
		info.deps = append(info.deps, j)
	}
	for p := range params {
		info.params = append(info.params, p)
	}
	sort.Ints(info.deps)
	sort.Strings(info.params)
	return info
}

// hash is the cache key of step i: its definition without the fields that
// do not change the result (name, provenance, off), the hashes of every
// step it depends on, and the values of the params it uses.
func (sc *scope) hash(i int) string {
	s := sc.h.Steps[i]
	def, _ := json.Marshal(struct {
		Op     string
		On     string
		Args   any
		Then   []Call
		Anchor *AnchorDef
		Mirror string
	}{s.Op, s.On, s.Args, s.Then, s.Anchor, s.Mirror})
	hw := sha256.New()
	hw.Write(def)
	info := sc.infos[i]
	for _, d := range info.deps {
		fmt.Fprintf(hw, "|%s", sc.hashes[d])
	}
	if info.input >= 0 {
		fmt.Fprintf(hw, "|in:%s", sc.hashes[info.input])
	}
	for _, p := range info.params {
		v, _ := json.Marshal(sc.h.Params[p])
		fmt.Fprintf(hw, "|$%s=%s", p, v)
	}
	return hex.EncodeToString(hw.Sum(nil))
}

// run evaluates step i. Panics inside the library become the step's error.
func (sc *scope) run(i int) (e *entry) {
	s := sc.h.Steps[i]
	e = &entry{status: OK}
	sc.consumed = map[string]bool{}
	defer func() {
		if r := recover(); r != nil {
			e.val, e.status, e.msg = reflect.Value{}, Failed, fmt.Sprint(r)
		}
		if e.status == OK {
			for n := range sc.consumed {
				e.consumes = append(e.consumes, n)
			}
			sort.Strings(e.consumes)
		}
	}()
	fail := func(err error) *entry {
		e.val = reflect.Value{}
		var b blockedErr
		if errors.As(err, &b) {
			e.status, e.msg = Blocked, b.msg
		} else {
			e.status, e.msg = Failed, err.Error()
		}
		return e
	}

	info := sc.infos[i]
	var in reflect.Value
	if s.Anchor == nil && (s.On != "" || info.input >= 0) {
		var err error
		if s.On != "" {
			in, err = sc.ref(s.On)
			if !strings.Contains(s.On, ".") {
				sc.consumed[s.On] = true
			}
		} else {
			in, err = sc.slotValue(info.input)
			e.inputConsumed = true
		}
		if err != nil {
			return fail(fmt.Errorf("on: %w", err))
		}
	} else if s.Anchor == nil && implicitInput(s) {
		return fail(blockedErr{fmt.Sprintf("%s needs an input and there is no step before it", s.Op)})
	}
	sc.input = in

	var val reflect.Value
	var err error
	if s.Anchor != nil {
		var surf reflect.Value
		if s.Anchor.On == "" {
			if info.input < 0 {
				return fail(blockedErr{"an anchor needs a surface: give on, or put a step before it"})
			}
			surf, err = sc.slotValue(info.input)
			if err != nil {
				return fail(err)
			}
		}
		var p anchor.Point
		p, err = sc.anchorPoint(s.Anchor, surf)
		if err != nil {
			return fail(err)
		}
		val = reflect.ValueOf(p)
	} else {
		val, err = sc.apply(s, in)
		if err != nil {
			return fail(err)
		}
	}

	// Warps since the last grid, through every solid this step built on.
	for _, d := range info.deps {
		if sc.slots[d].depth > e.depth && isForm(sc.slots[d].val) {
			e.depth = sc.slots[d].depth
		}
	}
	if warps[s.Op] {
		e.depth++
		if s.Mirror != "" {
			e.depth++
		}
	}
	if o := sc.rt.opts; o.BakeEvery > 0 && e.depth >= o.BakeEvery && val.Type() == solidType {
		val = reflect.ValueOf(sculpt.Sample(val.Interface().(*solid.Solid), o.BakeCells).Solid())
		e.depth, e.baked = 0, true
	}
	if val.Type() == gridType {
		e.depth = 0
	}
	e.val = val
	return e
}

// apply runs an op step: the op, its mirror, then its `then` chain.
func (sc *scope) apply(s Step, in reflect.Value) (reflect.Value, error) {
	var c *callable
	method := !strings.Contains(s.Op, ".")
	if method {
		if !in.IsValid() {
			return reflect.Value{}, fmt.Errorf("%s is a method and needs an input: give on", s.Op)
		}
		var err error
		if c, err = lookupMethod(in, s.Op); err != nil {
			return reflect.Value{}, err
		}
	} else {
		var ok bool
		if c, ok = lookupFunc(s.Op); !ok {
			return reflect.Value{}, fmt.Errorf("no op %s%s", s.Op, suggest(s.Op, funcNames()))
		}
	}
	var first *reflect.Value
	if !method && in.IsValid() {
		first = &in
	}
	out, args, err := sc.call(c, first, s.Args, "args")
	if err != nil {
		return reflect.Value{}, err
	}

	if s.Mirror != "" {
		if !in.IsValid() {
			return reflect.Value{}, fmt.Errorf("mirror applies a step twice to its input, and %s has none", s.Op)
		}
		axis := strings.Index("xyz", s.Mirror)
		mirrored := make([]reflect.Value, len(args))
		for k, a := range args {
			mirrored[k] = mirrorVal(a, axis)
		}
		if method {
			c2, err := lookupMethod(out, s.Op)
			if err != nil {
				return reflect.Value{}, err
			}
			out, err = invoke(c2, mirrored)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("mirror: %w", err)
			}
		} else {
			first, err := assign(out, c.fn.Type().In(0))
			if err != nil {
				return reflect.Value{}, fmt.Errorf("mirror: %w", err)
			}
			if out, err = invoke(c, append([]reflect.Value{first}, mirrored...)); err != nil {
				return reflect.Value{}, fmt.Errorf("mirror: %w", err)
			}
		}
	}

	for k, t := range s.Then {
		mc, err := lookupMethod(out, t.Op)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("then[%d]: %w", k, err)
		}
		if out, _, err = sc.call(mc, nil, t.Args, fmt.Sprintf("then[%d].args", k)); err != nil {
			return reflect.Value{}, fmt.Errorf("then[%d] %s: %w", k, t.Op, err)
		}
	}
	return out, nil
}

// call converts args to c's parameters and calls it. first, when given,
// fills the first parameter. It returns the result and the converted
// arguments other than first, which a mirror reuses.
func (sc *scope) call(c *callable, first *reflect.Value, raw any, where string) (reflect.Value, []reflect.Value, error) {
	t := c.fn.Type()
	n := t.NumIn()
	start := 0
	var in []reflect.Value
	if first != nil {
		if n == 0 {
			return reflect.Value{}, nil, fmt.Errorf("%s takes no input", c.key)
		}
		v, err := assign(*first, t.In(0))
		if err != nil {
			return reflect.Value{}, nil, fmt.Errorf("on: %w", err)
		}
		in = append(in, v)
		start = 1
	}
	var args []reflect.Value
	add := func(v reflect.Value) {
		args = append(args, v)
		if !sc.input.IsValid() && isForm(v) {
			sc.input = v // the first solid argument is the surface an inline anchor is on
		}
	}
	nonVar := n
	if t.IsVariadic() {
		nonVar = n - 1
	}

	switch a := raw.(type) {
	case nil, map[string]any:
		m, _ := a.(map[string]any)
		if c.params == nil && len(m) > 0 {
			return reflect.Value{}, nil, fmt.Errorf("%s takes args by position (a list); its parameter names are not known", c.key)
		}
		// One struct parameter and args that do not name it: args are the
		// struct's fields, as in obj.Bolt or a registered func(p Params).
		if n-start == 1 && !t.IsVariadic() && isParamStruct(t.In(start)) {
			if _, named := lookupArg(m, c.params[start]); !named {
				v, err := sc.convert(m, t.In(start), where)
				if err != nil {
					return reflect.Value{}, nil, err
				}
				add(v)
				break
			}
		}
		// A key that is no parameter is the likelier mistake than a missing
		// one (a typo makes both), so it is reported first.
		used := map[string]bool{}
		for k := start; k < n; k++ {
			if key, ok := lookupArg(m, c.params[k]); ok {
				used[key] = true
			}
		}
		for key := range m {
			if !used[key] {
				return reflect.Value{}, nil, fmt.Errorf("%s has no parameter %s%s%s", c.key, key, suggest(key, c.params[start:]), paramList(c, t, start))
			}
		}
		for k := start; k < n; k++ {
			name := c.params[k]
			key, ok := lookupArg(m, name)
			if k >= nonVar { // variadic
				if !ok {
					continue
				}
				list, isList := m[key].([]any)
				if !isList {
					list = []any{m[key]}
				}
				for j, x := range list {
					v, err := sc.convert(x, t.In(k).Elem(), fmt.Sprintf("%s.%s[%d]", where, name, j))
					if err != nil {
						return reflect.Value{}, nil, err
					}
					add(v)
				}
				continue
			}
			if !ok {
				return reflect.Value{}, nil, fmt.Errorf("%s needs %s (%s)%s", c.key, name, typeName(t.In(k)), paramList(c, t, start))
			}
			v, err := sc.convert(m[key], t.In(k), where+"."+name)
			if err != nil {
				return reflect.Value{}, nil, err
			}
			add(v)
		}
	case []any:
		if len(a) < nonVar-start || (!t.IsVariadic() && len(a) > n-start) {
			return reflect.Value{}, nil, fmt.Errorf("%s takes %d args after its input, got %d%s", c.key, n-start, len(a), paramList(c, t, start))
		}
		for j, x := range a {
			k := start + j
			pt := reflect.Type(nil)
			if k >= nonVar {
				pt = t.In(n - 1).Elem()
			} else {
				pt = t.In(k)
			}
			v, err := sc.convert(x, pt, fmt.Sprintf("%s[%d]", where, j))
			if err != nil {
				return reflect.Value{}, nil, err
			}
			add(v)
		}
	}
	out, err := invoke(c, append(in, args...))
	return out, args, err
}

// invoke calls c, turning a panic or a returned error into an error.
func invoke(c *callable, in []reflect.Value) (out reflect.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	res := c.fn.Call(in)
	if len(res) == 2 && !res[1].IsNil() {
		return reflect.Value{}, res[1].Interface().(error)
	}
	if len(res) == 0 {
		return reflect.Value{}, fmt.Errorf("%s returns nothing", c.key)
	}
	if k := res[0].Kind(); (k == reflect.Pointer || k == reflect.Interface) && res[0].IsNil() {
		return reflect.Value{}, fmt.Errorf("%s returned nothing", c.key)
	}
	return res[0], nil
}

func paramList(c *callable, t reflect.Type, start int) string {
	if c.params == nil {
		return ""
	}
	var ps []string
	for k := start; k < t.NumIn(); k++ {
		pt := t.In(k)
		if t.IsVariadic() && k == t.NumIn()-1 {
			ps = append(ps, c.params[k]+" ..."+typeName(pt.Elem()))
		} else {
			ps = append(ps, c.params[k]+" "+typeName(pt))
		}
	}
	return "; it takes " + strings.Join(ps, ", ")
}

// lookupArg finds a parameter in args by exact name, then ignoring case.
func lookupArg(m map[string]any, name string) (string, bool) {
	if _, ok := m[name]; ok {
		return name, true
	}
	for k := range m {
		if strings.EqualFold(k, name) {
			return k, true
		}
	}
	return "", false
}

func isParamStruct(t reflect.Type) bool {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t.Kind() == reflect.Struct && t != vec3Type && t != vec2Type && t != pointType && t != anchoredT && t != placeType
}

func isForm(v reflect.Value) bool {
	if !v.IsValid() {
		return false
	}
	switch v.Type() {
	case solidType, gridType:
		return true
	}
	return false
}

func mirrorVal(v reflect.Value, axis int) reflect.Value {
	switch v.Type() {
	case vec3Type:
		p := v.Interface().(v3.Vec)
		p.Set(axis, -p.Get(axis))
		return reflect.ValueOf(p)
	case vec3sType:
		ps := append([]v3.Vec{}, v.Interface().([]v3.Vec)...)
		for k := range ps {
			ps[k].Set(axis, -ps[k].Get(axis))
		}
		return reflect.ValueOf(ps)
	}
	return v
}

// slotValue is step j's value for a step that uses it.
func (sc *scope) slotValue(j int) (reflect.Value, error) {
	sl := sc.slots[j]
	label := sc.label(j)
	switch {
	case sl.val.IsValid():
		return sl.val, nil
	case sl.status == Off:
		return reflect.Value{}, blockedErr{fmt.Sprintf("needs %s, which is off and has no input to pass on", label)}
	case sl.status == Failed:
		return reflect.Value{}, blockedErr{fmt.Sprintf("needs %s, which failed", label)}
	default:
		return reflect.Value{}, blockedErr{fmt.Sprintf("needs %s, which is blocked", label)}
	}
}

func (sc *scope) label(j int) string {
	if n := sc.h.Steps[j].Name; n != "" {
		return "step " + n
	}
	return fmt.Sprintf("step #%d", j)
}

// ref resolves a reference: a step name, then any members, each a method
// with no arguments or a field: head, head.Top, NoseTip.Normal.
func (sc *scope) ref(r string) (reflect.Value, error) {
	parts := strings.Split(r, ".")
	j, ok := sc.names[parts[0]]
	if !ok {
		return reflect.Value{}, blockedErr{fmt.Sprintf("needs step %s, which does not exist", parts[0])}
	}
	if j >= sc.i {
		return reflect.Value{}, fmt.Errorf("step %s comes at or after this one; a step can only use earlier steps", parts[0])
	}
	v, err := sc.slotValue(j)
	if err != nil {
		return v, err
	}
	for _, m := range parts[1:] {
		if v, err = member(v, m); err != nil {
			return v, fmt.Errorf("%s: %w", r, err)
		}
	}
	return v, nil
}

func member(v reflect.Value, name string) (reflect.Value, error) {
	if m := v.MethodByName(name); m.IsValid() && m.Type().NumIn() == 0 && m.Type().NumOut() == 1 {
		return invoke(&callable{key: name, fn: m}, nil)
	}
	s := v
	for s.Kind() == reflect.Pointer {
		s = s.Elem()
	}
	if s.Kind() == reflect.Struct {
		if f := s.FieldByName(name); f.IsValid() && f.CanInterface() {
			return f, nil
		}
	}
	var names []string
	for _, m := range methodNames(v.Type()) {
		if mm := v.MethodByName(m); mm.IsValid() && mm.Type().NumIn() == 0 && mm.Type().NumOut() == 1 {
			names = append(names, m)
		}
	}
	return reflect.Value{}, fmt.Errorf("a %s has no %s%s", typeName(v.Type()), name, suggest(name, names))
}

// convert turns a value from the file into a value of type t.
func (sc *scope) convert(raw any, t reflect.Type, where string) (reflect.Value, error) {
	wrap := func(err error) error {
		var b blockedErr
		if errors.As(err, &b) {
			return blockedErr{where + ": " + b.msg}
		}
		return fmt.Errorf("%s: %w", where, err)
	}
	switch r := raw.(type) {
	case string:
		if strings.HasPrefix(r, "$") {
			p, ok := sc.h.Params[r[1:]]
			if !ok {
				return reflect.Value{}, wrap(fmt.Errorf("no param %s", r[1:]))
			}
			return sc.convert(p, t, where)
		}
		if t.Kind() == reflect.String {
			return reflect.ValueOf(r).Convert(t), nil
		}
		v, err := sc.ref(r)
		if err != nil {
			return reflect.Value{}, wrap(err)
		}
		out, err := assign(v, t)
		if err != nil {
			return reflect.Value{}, wrap(err)
		}
		// The step uses this one's solid: whole, or through a body anchor
		// that a placement will union with (plate.Top in On).
		if base, member, _ := strings.Cut(r, "."); member == "" || t == anchoredT {
			sc.consumed[base] = true
		}
		return out, nil
	case map[string]any:
		if op, ok := r["op"].(string); ok && (t.Kind() != reflect.Struct || !hasField(t, "op")) {
			v, err := sc.nested(op, r)
			if err != nil {
				return reflect.Value{}, wrap(err)
			}
			out, err := assign(v, t)
			if err != nil {
				return reflect.Value{}, wrap(err)
			}
			return out, nil
		}
		if a, ok := r["anchor"]; ok && len(r) == 1 {
			def, err := decodeAnchor(a)
			if err != nil {
				return reflect.Value{}, wrap(err)
			}
			p, err := sc.anchorPoint(def, sc.input)
			if err != nil {
				return reflect.Value{}, wrap(err)
			}
			out, err := assign(reflect.ValueOf(p), t)
			if err != nil {
				return reflect.Value{}, wrap(err)
			}
			return out, nil
		}
		if t.Kind() == reflect.Pointer && t.Elem().Kind() == reflect.Struct {
			v, err := sc.convert(r, t.Elem(), where)
			if err != nil {
				return v, err
			}
			p := reflect.New(t.Elem())
			p.Elem().Set(v)
			return p, nil
		}
		if t.Kind() == reflect.Struct {
			out := reflect.New(t).Elem()
			for k, x := range r {
				f, ok := fieldFor(t, k)
				if !ok {
					return reflect.Value{}, wrap(fmt.Errorf("a %s has no field %s%s", typeName(t), k, suggest(k, fieldNames(t))))
				}
				v, err := sc.convert(x, f.Type, where+"."+k)
				if err != nil {
					return reflect.Value{}, err
				}
				out.FieldByIndex(f.Index).Set(v)
			}
			return out, nil
		}
		return reflect.Value{}, wrap(fmt.Errorf("want a %s, got a map", typeName(t)))
	case []any:
		switch t {
		case vec3Type, vec2Type:
			want := 3
			if t == vec2Type {
				want = 2
			}
			if len(r) != want {
				return reflect.Value{}, wrap(fmt.Errorf("a %s is %d numbers, got %d", typeName(t), want, len(r)))
			}
			out := reflect.New(t).Elem()
			for k, x := range r {
				f, err := sc.convert(x, floatType, fmt.Sprintf("%s[%d]", where, k))
				if err != nil {
					return reflect.Value{}, err
				}
				out.Field(k).SetFloat(f.Float())
			}
			return out, nil
		}
		switch t.Kind() {
		case reflect.Slice:
			out := reflect.MakeSlice(t, len(r), len(r))
			for k, x := range r {
				v, err := sc.convert(x, t.Elem(), fmt.Sprintf("%s[%d]", where, k))
				if err != nil {
					return reflect.Value{}, err
				}
				out.Index(k).Set(v)
			}
			return out, nil
		case reflect.Array:
			if len(r) != t.Len() {
				return reflect.Value{}, wrap(fmt.Errorf("want %d values, got %d", t.Len(), len(r)))
			}
			out := reflect.New(t).Elem()
			for k, x := range r {
				v, err := sc.convert(x, t.Elem(), fmt.Sprintf("%s[%d]", where, k))
				if err != nil {
					return reflect.Value{}, err
				}
				out.Index(k).Set(v)
			}
			return out, nil
		}
		return reflect.Value{}, wrap(fmt.Errorf("want a %s, got a list", typeName(t)))
	case bool:
		if t.Kind() == reflect.Bool {
			return reflect.ValueOf(r).Convert(t), nil
		}
		return reflect.Value{}, wrap(fmt.Errorf("want a %s, got %v", typeName(t), r))
	case nil:
		return reflect.Value{}, wrap(fmt.Errorf("no value; want a %s", typeName(t)))
	}
	f, ok := toFloat(raw)
	if !ok {
		return reflect.Value{}, wrap(fmt.Errorf("cannot use %v (%T)", raw, raw))
	}
	switch t.Kind() {
	case reflect.Float64, reflect.Float32:
		return reflect.ValueOf(f).Convert(t), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if f != float64(int64(f)) {
			return reflect.Value{}, wrap(fmt.Errorf("want a whole number, got %v", f))
		}
		return reflect.ValueOf(int64(f)).Convert(t), nil
	}
	return reflect.Value{}, wrap(fmt.Errorf("want a %s, got the number %v", typeName(t), f))
}

func toFloat(raw any) (float64, bool) {
	switch n := raw.(type) {
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint64:
		return float64(n), true
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	}
	return 0, false
}

func hasField(t reflect.Type, name string) bool {
	_, ok := fieldFor(t, name)
	return ok
}

func fieldFor(t reflect.Type, name string) (reflect.StructField, bool) {
	if f, ok := t.FieldByName(name); ok && f.IsExported() {
		return f, true
	}
	for i := 0; i < t.NumField(); i++ {
		if f := t.Field(i); f.IsExported() && strings.EqualFold(f.Name, name) {
			return f, true
		}
	}
	return reflect.StructField{}, false
}

func fieldNames(t reflect.Type) []string {
	var out []string
	for i := 0; i < t.NumField(); i++ {
		if f := t.Field(i); f.IsExported() {
			out = append(out, f.Name)
		}
	}
	return out
}

// nested evaluates an inline {op: ..., args: ...}: a function, never a
// method, since an inline value has no input.
func (sc *scope) nested(op string, m map[string]any) (reflect.Value, error) {
	for k := range m {
		if k != "op" && k != "args" {
			return reflect.Value{}, fmt.Errorf("an inline op has op and args, not %s", k)
		}
	}
	c, ok := lookupFunc(op)
	if !ok {
		return reflect.Value{}, fmt.Errorf("no op %s%s", op, suggest(op, funcNames()))
	}
	saved := sc.input
	defer func() { sc.input = saved }()
	out, _, err := sc.call(c, nil, m["args"], op)
	return out, err
}

// assign makes v fit type t, with the conversions a history leans on: an
// anchor where a point is wanted, a grid where a solid is wanted.
func assign(v reflect.Value, t reflect.Type) (reflect.Value, error) {
	if !v.IsValid() {
		return v, fmt.Errorf("no value")
	}
	vt := v.Type()
	if vt.AssignableTo(t) {
		return v, nil
	}
	switch {
	case t == vec3Type && vt == anchoredT:
		return v.FieldByName("Point"), nil
	case t == vec3Type && vt == pointType:
		return v.FieldByName("Pos"), nil
	case vt == gridType && (t == solidType || t == sdf3Type):
		s := v.Interface().(*sculpt.Grid).Solid()
		return assign(reflect.ValueOf(s), t)
	case vt == anchoredT && t == solidType:
		return v.FieldByName("Solid"), nil
	case t == floatType && vt.ConvertibleTo(t) && vt.Kind() >= reflect.Int && vt.Kind() <= reflect.Float64:
		return v.Convert(t), nil
	case vt.Kind() == reflect.Func && t.Kind() == reflect.Func && vt.ConvertibleTo(t):
		return v.Convert(t), nil
	}
	return reflect.Value{}, fmt.Errorf("this is a %s, and a %s is wanted", typeName(vt), typeName(t))
}

func decodeAnchor(raw any) (*AnchorDef, error) {
	b, err := json.Marshal(raw)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(strings.NewReader(string(b)))
	dec.DisallowUnknownFields()
	var d AnchorDef
	if err := dec.Decode(&d); err != nil {
		return nil, fmt.Errorf("anchor: %w", err)
	}
	if (d.Ray == nil) == (d.Nearest == nil) {
		return nil, fmt.Errorf("an anchor needs exactly one of ray or nearest")
	}
	return &d, nil
}

// anchorPoint resolves an anchor definition. surf is the surface used when
// the definition has no `on`.
func (sc *scope) anchorPoint(d *AnchorDef, surf reflect.Value) (p anchor.Point, err error) {
	if d.On != "" {
		v, err := sc.ref(d.On)
		if err != nil {
			return p, fmt.Errorf("anchor on: %w", err)
		}
		surf = v
	}
	if !surf.IsValid() {
		return p, fmt.Errorf("anchor: no surface to find it on; give on")
	}
	sv, err := assign(surf, solidType)
	if err != nil {
		return p, fmt.Errorf("anchor on: %w", err)
	}
	s := sv.Interface().(*solid.Solid)
	vec := func(raw any, where string) (v3.Vec, error) {
		v, err := sc.convert(raw, vec3Type, where)
		if err != nil {
			return v3.Vec{}, err
		}
		return v.Interface().(v3.Vec), nil
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	switch {
	case d.Ray != nil:
		from, err := vec(d.Ray.From, "anchor.ray.from")
		if err != nil {
			return p, err
		}
		switch {
		case d.Ray.Dir != nil && d.Ray.Toward == nil:
			dir, err := vec(d.Ray.Dir, "anchor.ray.dir")
			if err != nil {
				return p, err
			}
			p = anchor.Ray(s, from, dir)
		case d.Ray.Toward != nil && d.Ray.Dir == nil:
			to, err := vec(d.Ray.Toward, "anchor.ray.toward")
			if err != nil {
				return p, err
			}
			p = anchor.Toward(s, from, to)
		default:
			return p, fmt.Errorf("anchor.ray needs exactly one of dir or toward")
		}
	default:
		near, err := vec(d.Nearest, "anchor.nearest")
		if err != nil {
			return p, err
		}
		p = anchor.Nearest(s, near)
	}
	if d.Offset != nil {
		off, err := sc.convert(d.Offset, floatType, "anchor.offset")
		if err != nil {
			return p, err
		}
		p.Pos = p.Offset(off.Float())
	}
	return p, nil
}

func kindOf(v reflect.Value) string {
	if !v.IsValid() {
		return ""
	}
	switch v.Type() {
	case solidType:
		return "solid"
	case gridType:
		return "grid"
	case shapeType:
		return "shape"
	case pointType:
		return "anchor"
	case placeType:
		return "placement"
	}
	return "value"
}

// describe fills the parts of a step result an agent can read without the
// geometry: an anchor's point, a measured value.
func describe(sr *StepResult, v reflect.Value) {
	switch sr.Kind {
	case "anchor":
		p := v.Interface().(anchor.Point)
		sr.Anchor = &p
	case "value":
		switch v.Kind() {
		case reflect.Bool, reflect.Int, reflect.Int64, reflect.Float64, reflect.String:
			sr.Value = v.Interface()
		default:
			if v.Type() == vec3Type || v.Type() == vec3sType {
				sr.Value = v.Interface()
			}
		}
	}
	sr.value = v
}
