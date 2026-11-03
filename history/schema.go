package history

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// History is a part as data: named parameters and the steps that build
// it, in order. It is the file on disk, the thing a UI or an agent edits,
// and the only state there is. Deleting a step and evaluating again is
// undo.
type History struct {
	// Header is the comment block at the top of the file, kept through a
	// load and save. Every other comment is lost on save.
	Header string `yaml:"-" json:"header,omitempty"`
	// Params are named numbers or vectors, used in args as $name.
	Params map[string]any `yaml:"params,omitempty" json:"params,omitempty"`
	Steps  []Step         `yaml:"steps" json:"steps"`
}

// Step is one operation, or one named anchor.
//
// An operation names a library function (`op: solid.Sphere`), a method on
// its input (`op: Translate`), or a registered function. Its input is the
// step named by `on`; for a method, or a function whose first parameter is
// a solid and is not given in args, an omitted `on` means the previous
// step. Args are by parameter name (a map) or by position (a list).
//
// An anchor step (`anchor:`) has no op. Its value is a point on a surface
// with the normal there, which later steps use by name.
type Step struct {
	Name   string     `yaml:"name,omitempty" json:"name,omitempty"`
	Op     string     `yaml:"op,omitempty" json:"op,omitempty"`
	On     string     `yaml:"on,omitempty" json:"on,omitempty"`
	Args   any        `yaml:"args,omitempty" json:"args,omitempty"`
	Then   []Call     `yaml:"then,omitempty" json:"then,omitempty"`
	Anchor *AnchorDef `yaml:"anchor,omitempty" json:"anchor,omitempty"`
	// Mirror runs the step a second time on its own result with every
	// vector argument mirrored across the plane normal to this axis: x,
	// y or z. A pull at the left cheek, mirrored in x, also pulls the right.
	Mirror string `yaml:"mirror,omitempty" json:"mirror,omitempty"`
	// Off keeps the step in the file but skips it. A step that used it gets
	// its input instead, so turning a step off is as if it never happened.
	Off bool `yaml:"off,omitempty" json:"off,omitempty"`
	// By is who made the step: user (the default) or agent.
	By string `yaml:"by,omitempty" json:"by,omitempty"`
	// From is an opaque reference to what produced the step, such as the
	// chat or prompt (chat/123). Stored and returned exactly; never read.
	From string `yaml:"from,omitempty" json:"from,omitempty"`
}

// Call is one method in a step's `then` chain.
type Call struct {
	Op   string `yaml:"op" json:"op"`
	Args any    `yaml:"args,omitempty" json:"args,omitempty"`
}

// AnchorDef finds a point on a surface. Exactly one of Ray or Nearest.
type AnchorDef struct {
	// On is the step whose surface the anchor is found on. In an anchor
	// step it defaults to the previous step; inline in a step's args, to
	// that step's input.
	On      string  `yaml:"on,omitempty" json:"on,omitempty"`
	Ray     *RayDef `yaml:"ray,omitempty" json:"ray,omitempty"`
	Nearest any     `yaml:"nearest,omitempty" json:"nearest,omitempty"`
	// Offset moves the found point along the normal: positive out of the
	// material, negative into it.
	Offset any `yaml:"offset,omitempty" json:"offset,omitempty"`
}

// RayDef is a ray from a point, along a direction or toward a point.
type RayDef struct {
	From   any `yaml:"from" json:"from"`
	Dir    any `yaml:"dir,omitempty" json:"dir,omitempty"`
	Toward any `yaml:"toward,omitempty" json:"toward,omitempty"`
}

var nameRE = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)

// Parse reads a history from YAML (or JSON, which is YAML). Unknown keys
// are an error, so a misspelt `radious` fails here and not silently.
func Parse(b []byte) (*History, error) {
	h := &History{Header: header(b)}
	dec := yaml.NewDecoder(bytes.NewReader(b))
	dec.KnownFields(true)
	if err := dec.Decode(h); err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("history: %w", err)
	}
	if err := h.Validate(); err != nil {
		return nil, err
	}
	return h, nil
}

// Load reads and parses a history file.
func Load(path string) (*History, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	h, err := Parse(b)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return h, nil
}

// Save writes the history to path atomically: a reader never sees half a
// file, and a watcher sees one change.
func (h *History) Save(path string) error {
	b, err := h.Marshal()
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// Marshal returns the history as YAML: the header, then params and steps,
// with lists of numbers written on one line.
func (h *History) Marshal() ([]byte, error) {
	if err := h.Validate(); err != nil {
		return nil, err
	}
	var n yaml.Node
	if err := n.Encode(h); err != nil {
		return nil, err
	}
	flowScalarLists(&n)
	flowSmallMaps(&n, true)
	var b bytes.Buffer
	if h.Header != "" {
		b.WriteString(h.Header)
		if !strings.HasSuffix(h.Header, "\n") {
			b.WriteByte('\n')
		}
	}
	enc := yaml.NewEncoder(&b)
	enc.SetIndent(2)
	if err := enc.Encode(&n); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func flowScalarLists(n *yaml.Node) {
	if n.Kind == yaml.SequenceNode && len(n.Content) > 0 {
		flat := true
		for _, c := range n.Content {
			if c.Kind != yaml.ScalarNode {
				flat = false
				break
			}
		}
		if flat {
			n.Style = yaml.FlowStyle
		}
	}
	for _, c := range n.Content {
		flowScalarLists(c)
	}
}

// flowSmallMaps writes a short map of plain values on one line, op first:
// args: {radius: 12}, {op: solid.RoundMin, args: {k: 3}}. A step itself
// and the top level stay one key per line.
func flowSmallMaps(n *yaml.Node, block bool) {
	if n.Kind == yaml.DocumentNode {
		for _, c := range n.Content {
			flowSmallMaps(c, true)
		}
		return
	}
	if n.Kind == yaml.MappingNode {
		// The steps list's items stay block; everything below may flow.
		for k := 0; k+1 < len(n.Content); k += 2 {
			key, val := n.Content[k], n.Content[k+1]
			if key.Value == "steps" && val.Kind == yaml.SequenceNode {
				for _, item := range val.Content {
					flowSmallMaps(item, true)
				}
				continue
			}
			flowSmallMaps(val, false)
		}
		for k := 0; k < len(n.Content); k += 2 {
			// yaml.v3 quotes keys YAML 1.1 read as booleans; ours (on, off)
			// are keys, and YAML 1.2 reads them as the strings they are.
			if key := n.Content[k]; key.Style == yaml.DoubleQuotedStyle && (key.Value == "on" || key.Value == "off") {
				key.Style = 0
			}
		}
		if !block {
			opFirst(n)
		}
		if !block && width(n) <= 64 {
			n.Style = yaml.FlowStyle
		}
		return
	}
	for _, c := range n.Content {
		flowSmallMaps(c, false)
	}
}

// opFirst orders a map with an op key as op, args, then the rest.
func opFirst(n *yaml.Node) {
	rank := func(k string) int {
		switch k {
		case "op":
			return 0
		case "args":
			return 1
		}
		return 2
	}
	type kv struct{ k, v *yaml.Node }
	var pairs []kv
	hasOp := false
	for k := 0; k+1 < len(n.Content); k += 2 {
		pairs = append(pairs, kv{n.Content[k], n.Content[k+1]})
		hasOp = hasOp || n.Content[k].Value == "op"
	}
	if !hasOp {
		return
	}
	sort.SliceStable(pairs, func(i, j int) bool { return rank(pairs[i].k.Value) < rank(pairs[j].k.Value) })
	n.Content = n.Content[:0]
	for _, p := range pairs {
		n.Content = append(n.Content, p.k, p.v)
	}
}

// width is about how wide a node is on one line, or a large number when
// it holds anything that is not itself on one line.
func width(n *yaml.Node) int {
	switch n.Kind {
	case yaml.ScalarNode:
		return len(n.Value)
	case yaml.SequenceNode, yaml.MappingNode:
		if n.Style != yaml.FlowStyle && n.Kind == yaml.SequenceNode {
			return 1 << 20
		}
		w := 2
		for _, c := range n.Content {
			if c.Kind == yaml.MappingNode && c.Style != yaml.FlowStyle {
				return 1 << 20
			}
			w += width(c) + 2
		}
		return w
	}
	return 1 << 20
}

// header is the leading run of comment lines.
func header(b []byte) string {
	var out strings.Builder
	sc := bufio.NewScanner(bytes.NewReader(b))
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasPrefix(strings.TrimSpace(line), "#") {
			break
		}
		out.WriteString(line)
		out.WriteByte('\n')
	}
	return out.String()
}

// Validate checks what can be checked without evaluating: names are
// identifiers and unique, every step is an op or an anchor, and the
// provenance and mirror fields hold known values.
func (h *History) Validate() error {
	for name := range h.Params {
		if !nameRE.MatchString(name) {
			return fmt.Errorf("history: param %q: a name is letters, digits and _, starting with a letter", name)
		}
	}
	seen := map[string]int{}
	for i, s := range h.Steps {
		at := fmt.Sprintf("history: step %d", i)
		if s.Name != "" {
			at += fmt.Sprintf(" (%s)", s.Name)
			if !nameRE.MatchString(s.Name) {
				return fmt.Errorf("%s: a name is letters, digits and _, starting with a letter", at)
			}
			if j, dup := seen[s.Name]; dup {
				return fmt.Errorf("%s: step %d already has this name", at, j)
			}
			seen[s.Name] = i
		}
		switch {
		case s.Op == "" && s.Anchor == nil:
			return fmt.Errorf("%s: needs an op or an anchor", at)
		case s.Op != "" && s.Anchor != nil:
			return fmt.Errorf("%s: has both an op and an anchor; make them two steps", at)
		}
		if s.Anchor != nil {
			if (s.Anchor.Ray == nil) == (s.Anchor.Nearest == nil) {
				return fmt.Errorf("%s: an anchor needs exactly one of ray or nearest", at)
			}
			if s.Mirror != "" || len(s.Then) > 0 || s.Args != nil || s.On != "" {
				return fmt.Errorf("%s: an anchor step takes only anchor (put on inside it)", at)
			}
		}
		switch s.Mirror {
		case "", "x", "y", "z":
		default:
			return fmt.Errorf("%s: mirror is x, y or z, not %q", at, s.Mirror)
		}
		switch s.By {
		case "", "user", "agent":
		default:
			return fmt.Errorf("%s: by is user or agent, not %q", at, s.By)
		}
		switch s.Args.(type) {
		case nil, map[string]any, []any:
		default:
			return fmt.Errorf("%s: args is a map by name or a list by position", at)
		}
	}
	return nil
}

// Index returns the index of the step with this name, or -1.
func (h *History) Index(name string) int {
	for i, s := range h.Steps {
		if s.Name == name {
			return i
		}
	}
	return -1
}

// Insert puts a step at index i; i == len(Steps) appends.
func (h *History) Insert(i int, s Step) error {
	if i < 0 || i > len(h.Steps) {
		return fmt.Errorf("history: insert at %d: there are %d steps", i, len(h.Steps))
	}
	next := append(append(append([]Step{}, h.Steps[:i]...), s), h.Steps[i:]...)
	return h.replace(next)
}

// Set replaces step i.
func (h *History) Set(i int, s Step) error {
	if err := h.check(i); err != nil {
		return err
	}
	next := append([]Step{}, h.Steps...)
	next[i] = s
	return h.replace(next)
}

// Remove deletes step i. Steps that used it by name are blocked until
// they are changed; steps that used it as their implicit previous step
// get the step before it.
func (h *History) Remove(i int) error {
	if err := h.check(i); err != nil {
		return err
	}
	return h.replace(append(append([]Step{}, h.Steps[:i]...), h.Steps[i+1:]...))
}

// Move moves step from to index to, shifting the steps between.
func (h *History) Move(from, to int) error {
	if err := h.check(from); err != nil {
		return err
	}
	if err := h.check(to); err != nil {
		return err
	}
	next := append([]Step{}, h.Steps...)
	s := next[from]
	next = append(next[:from], next[from+1:]...)
	next = append(next[:to], append([]Step{s}, next[to:]...)...)
	return h.replace(next)
}

// SetOff turns step i off or back on.
func (h *History) SetOff(i int, off bool) error {
	if err := h.check(i); err != nil {
		return err
	}
	h.Steps[i].Off = off
	return nil
}

func (h *History) check(i int) error {
	if i < 0 || i >= len(h.Steps) {
		return fmt.Errorf("history: no step %d: there are %d steps", i, len(h.Steps))
	}
	return nil
}

// replace swaps in a new step list if it validates, and leaves the
// history unchanged if it does not.
func (h *History) replace(steps []Step) error {
	old := h.Steps
	h.Steps = steps
	if err := h.Validate(); err != nil {
		h.Steps = old
		return err
	}
	return nil
}

// Clone returns a deep copy, so an edit can be previewed without touching
// the original.
func (h *History) Clone() *History {
	c := &History{Header: h.Header}
	if h.Params != nil {
		c.Params = deepCopy(h.Params).(map[string]any)
	}
	c.Steps = make([]Step, len(h.Steps))
	for i, s := range h.Steps {
		c.Steps[i] = s.Clone()
	}
	return c
}

// Clone returns a deep copy of the step.
func (s Step) Clone() Step {
	c := s
	c.Args = deepCopy(s.Args)
	c.Then = nil
	for _, t := range s.Then {
		c.Then = append(c.Then, Call{Op: t.Op, Args: deepCopy(t.Args)})
	}
	if s.Anchor != nil {
		a := *s.Anchor
		a.Nearest = deepCopy(a.Nearest)
		a.Offset = deepCopy(a.Offset)
		if a.Ray != nil {
			r := *a.Ray
			r.From, r.Dir, r.Toward = deepCopy(r.From), deepCopy(r.Dir), deepCopy(r.Toward)
			a.Ray = &r
		}
		c.Anchor = &a
	}
	return c
}

func deepCopy(v any) any {
	switch t := v.(type) {
	case map[string]any:
		c := make(map[string]any, len(t))
		for k, x := range t {
			c[k] = deepCopy(x)
		}
		return c
	case []any:
		c := make([]any, len(t))
		for i, x := range t {
			c[i] = deepCopy(x)
		}
		return c
	}
	return v
}
