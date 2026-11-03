package history

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"sync"
)

// callable is one function or method a step can call, with what the
// source knows about it.
type callable struct {
	key    string        // solid.Sphere, solid.Solid.Translate, or a registered name
	fn     reflect.Value // for a method, bound to its receiver
	params []string      // parameter names in order; nil when unknown
	doc    string
	method bool
	reg    *registered // set for a registered function
}

type registered struct {
	name   string
	fn     reflect.Value
	params []string
	doc    string
}

var registry = struct {
	sync.RWMutex
	m map[string]*registered
}{m: map[string]*registered{}}

// Register makes a Go function an op, called from a history as
// `op: name` like any library function. Reflection reads its signature, so
// argument conversion and the palette treat it as a built-in. params names
// the parameters in order, for args given by name; leave them out for a
// function whose only parameter is a struct, and the args map fills the
// struct's fields by name. doc is the one-line description an agent or a UI
// shows for it.
//
// The name must contain a dot (my.Bracket) and must not already be an op.
// Register panics on a bad registration: it runs at program start, where a
// panic is the clearest report.
//
//	history.Register("my.Bracket", Bracket, "A bracket with a slot for the hinge pin.", "width", "height", "slot")
func Register(name string, fn any, doc string, params ...string) {
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		panic(fmt.Sprintf("history.Register(%q): %T is not a function", name, fn))
	}
	if !strings.Contains(name, ".") {
		panic(fmt.Sprintf("history.Register(%q): the name needs a dot, like my.Bracket, so it cannot be mistaken for a method", name))
	}
	if _, builtin := genFuncs[name]; builtin {
		panic(fmt.Sprintf("history.Register(%q): that is a library function", name))
	}
	t := v.Type()
	if t.NumOut() < 1 || t.NumOut() > 2 || (t.NumOut() == 2 && t.Out(1) != errorType) {
		panic(fmt.Sprintf("history.Register(%q): a step returns one value, or a value and an error; %s does not", name, t))
	}
	if len(params) == 0 {
		for i := 0; i < t.NumIn(); i++ {
			params = append(params, fmt.Sprintf("arg%d", i))
		}
	}
	if len(params) != t.NumIn() {
		panic(fmt.Sprintf("history.Register(%q): %d parameter names for a function of %d parameters", name, len(params), t.NumIn()))
	}
	registry.Lock()
	defer registry.Unlock()
	if _, dup := registry.m[name]; dup {
		panic(fmt.Sprintf("history.Register(%q): registered twice", name))
	}
	registry.m[name] = &registered{name: name, fn: v, params: params, doc: doc}
}

// importOf returns the import path and identifier of a registered
// function, for the Go export. A closure or method value has none.
func (r *registered) importOf() (pkgPath, ident string, ok bool) {
	full := runtime.FuncForPC(r.fn.Pointer()).Name() // github.com/me/parts.Bracket
	slash := strings.LastIndex(full, "/")
	dot := strings.Index(full[slash+1:], ".")
	if dot < 0 {
		return "", "", false
	}
	pkgPath, ident = full[:slash+1+dot], full[slash+1+dot+1:]
	if strings.ContainsAny(ident, ".()") || pkgPath == "main" {
		return "", "", false
	}
	return pkgPath, ident, true
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// lookupFunc finds a package-level function by op name.
func lookupFunc(op string) (*callable, bool) {
	registry.RLock()
	r, ok := registry.m[op]
	registry.RUnlock()
	if ok {
		return &callable{key: op, fn: r.fn, params: r.params, doc: r.doc, reg: r}, true
	}
	fn, ok := genFuncs[op]
	if !ok {
		return nil, false
	}
	s := genSigs[op]
	return &callable{key: op, fn: reflect.ValueOf(fn), params: s.params, doc: s.doc}, true
}

// lookupMethod finds a method on a value.
func lookupMethod(recv reflect.Value, name string) (*callable, error) {
	m := recv.MethodByName(name)
	if !m.IsValid() && recv.Kind() != reflect.Pointer && recv.CanAddr() {
		m = recv.Addr().MethodByName(name)
	}
	if !m.IsValid() && recv.Type() == gridType {
		// A grid is a solid too: Shell or Translate on a grid means on its Solid().
		if s := asSolid(recv); s != nil {
			return lookupMethod(reflect.ValueOf(s), name)
		}
	}
	if !m.IsValid() {
		return nil, fmt.Errorf("a %s has no method %s%s", typeName(recv.Type()), name, suggest(name, methodNames(recv.Type())))
	}
	key := typeKey(recv.Type()) + "." + name
	s := genSigs[key]
	return &callable{key: key, fn: m, params: s.params, doc: s.doc, method: true}, nil
}

// typeKey is pkg.Type for a type or a pointer to one: solid.Solid.
func typeKey(t reflect.Type) string {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return path.Base(t.PkgPath()) + "." + t.Name()
}

// typeName is how a type reads in a message: solid.Solid, v3.Vec, float64.
func typeName(t reflect.Type) string {
	return strings.TrimPrefix(t.String(), "*")
}

func methodNames(t reflect.Type) []string {
	var out []string
	pt := t
	if t.Kind() != reflect.Pointer {
		pt = reflect.PointerTo(t)
	}
	for i := 0; i < pt.NumMethod(); i++ {
		out = append(out, pt.Method(i).Name)
	}
	return out
}

// suggest returns ", did you mean X?" for the closest name, or "".
func suggest(name string, among []string) string {
	best, bestD := "", 1<<30
	for _, c := range among {
		if d := editDistance(strings.ToLower(name), strings.ToLower(c)); d < bestD {
			best, bestD = c, d
		}
	}
	if best != "" && bestD <= 1+len(name)/3 {
		return fmt.Sprintf("; did you mean %s?", best)
	}
	return ""
}

func editDistance(a, b string) int {
	prev := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur := make([]int, len(b)+1)
		cur[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			cur[j] = min(prev[j]+1, cur[j-1]+1, prev[j-1]+cost)
		}
		prev = cur
	}
	return prev[len(b)]
}

// funcNames is every function op name, library and registered, sorted.
func funcNames() []string {
	var out []string
	for k := range genFuncs {
		out = append(out, k)
	}
	registry.RLock()
	for k := range registry.m {
		out = append(out, k)
	}
	registry.RUnlock()
	sort.Strings(out)
	return out
}
