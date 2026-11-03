package history

import (
	"fmt"
	"reflect"

	"github.com/snowbldr/fluent-sdfx/anchor"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Result is one evaluation: every step's outcome and the bodies the
// history makes.
type Result struct {
	Steps  []StepResult `json:"steps"`
	Bodies []Body       `json:"bodies"`
}

// StepResult is one step's outcome, in a form an agent can act on: the
// status, and for a failure the message, which says what went wrong and
// usually how to fix it.
type StepResult struct {
	Index  int    `json:"index"`
	Name   string `json:"name,omitempty"`
	Op     string `json:"op"`
	Status Status `json:"status"`
	Error  string `json:"error,omitempty"`
	// Kind is what the step made: solid, grid, shape, anchor, placement or
	// value.
	Kind   string        `json:"kind,omitempty"`
	Anchor *anchor.Point `json:"anchor,omitempty"` // for an anchor step
	Value  any           `json:"value,omitempty"`  // for a number, vector or other plain value
	By     string        `json:"by,omitempty"`
	From   string        `json:"from,omitempty"`
	Cached bool          `json:"cached,omitempty"` // reused from an earlier evaluation
	Baked  bool          `json:"baked,omitempty"`  // sampled onto a grid to keep later steps fast
	Millis float64       `json:"millis"`           // time to evaluate, when not cached

	value reflect.Value
}

// Body is a solid the history makes: a step whose solid no later step
// went on to use. A plain chain makes one body, its last step; a part
// built from pieces that are never joined makes several.
type Body struct {
	Step  int          `json:"step"`
	Name  string       `json:"name"` // the step's name, or #index
	Solid *solid.Solid `json:"-"`
}

// Err returns the first failed or blocked step as an error, or nil.
func (r *Result) Err() error {
	for _, s := range r.Steps {
		if s.Status == Failed || s.Status == Blocked {
			return fmt.Errorf("step %s: %s: %s", s.label(), s.Status, s.Error)
		}
	}
	return nil
}

func (s StepResult) label() string {
	if s.Name != "" {
		return s.Name
	}
	return fmt.Sprintf("#%d", s.Index)
}

// Value returns the value of the step with this name, or nil.
func (r *Result) Value(name string) any {
	for _, s := range r.Steps {
		if s.Name == name && s.value.IsValid() {
			return s.value.Interface()
		}
	}
	return nil
}

// SolidOf returns the named step's value as a solid (a grid is
// converted), or nil.
func (r *Result) SolidOf(name string) *solid.Solid {
	for _, s := range r.Steps {
		if s.Name == name {
			return asSolid(s.value)
		}
	}
	return nil
}

// Solid returns every body as one solid, or nil when there is none.
func (r *Result) Solid() *solid.Solid {
	switch len(r.Bodies) {
	case 0:
		return nil
	case 1:
		return r.Bodies[0].Solid
	}
	ss := make([]*solid.Solid, len(r.Bodies))
	for i, b := range r.Bodies {
		ss[i] = b.Solid
	}
	return solid.UnionAll(ss...)
}

func asSolid(v reflect.Value) *solid.Solid {
	if !isForm(v) {
		return nil
	}
	s, err := assign(v, solidType)
	if err != nil {
		return nil
	}
	return s.Interface().(*solid.Solid)
}

func (r *Result) bodies(sc *scope) {
	// An off step is transparent: a step that uses it uses its input, and
	// it is never a body itself.
	through := func(j int) int {
		for j >= 0 && sc.h.Steps[j].Off {
			j = sc.infos[j].input
		}
		return j
	}
	used := make([]bool, len(r.Steps))
	for i := range r.Steps {
		e := sc.slots[i].e
		if sc.h.Steps[i].Off || e == nil || !e.val.IsValid() {
			continue // a failed step consumed nothing: what it would have used still shows
		}
		var js []int
		if e.inputConsumed {
			js = append(js, sc.infos[i].input)
		}
		for _, n := range e.consumes {
			if j, ok := sc.names[n]; ok {
				js = append(js, j)
			}
		}
		for _, j := range js {
			if j = through(j); j >= 0 {
				used[j] = true
			}
		}
	}
	for i := range r.Steps {
		r.Steps[i].value = sc.slots[i].val
		if used[i] || sc.h.Steps[i].Off {
			continue
		}
		if s := asSolid(sc.slots[i].val); s != nil {
			r.Bodies = append(r.Bodies, Body{Step: i, Name: r.Steps[i].label(), Solid: s})
		}
	}
}

// Hit is where a pick ray meets a body.
type Hit struct {
	Body  Body         `json:"body"`
	Point anchor.Point `json:"point"`
	Dist  float64      `json:"dist"` // along the ray
}

// Pick casts a ray at the bodies and returns the nearest hit. This is the
// SDF's own surface, not a mesh's, so the point is exact to the field; a
// UI picks with it and stores the ray as an anchor, not the hit.
func (r *Result) Pick(from, dir v3.Vec) (Hit, bool) {
	var best Hit
	found := false
	for _, b := range r.Bodies {
		p, ok := tryRay(b.Solid, from, dir)
		if !ok || p.Pos.Sub(from).Dot(dir) < 0 {
			continue
		}
		d := p.Pos.Sub(from).Length()
		if !found || d < best.Dist {
			best, found = Hit{Body: b, Point: p, Dist: d}, true
		}
	}
	return best, found
}

func tryRay(s *solid.Solid, from, dir v3.Vec) (p anchor.Point, ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	return anchor.Ray(s, from, dir), true
}

// AnchorInfo is a point a step can refer to by name.
type AnchorInfo struct {
	Ref    string  `json:"ref"` // as written in a step: NoseTip, head.Top
	Step   int     `json:"step"`
	Pos    v3.Vec  `json:"pos"`
	Normal *v3.Vec `json:"normal,omitempty"` // surface anchors only
}

// bodyAnchors are the bounding-box anchors listed for every named solid.
var bodyAnchors = []string{
	"Top", "Bottom", "Left", "Right", "Front", "Back",
	"TopFront", "TopBack", "TopLeft", "TopRight", "BottomFront", "BottomBack", "BottomLeft", "BottomRight",
	"FrontLeft", "FrontRight", "BackLeft", "BackRight",
	"TopFrontLeft", "TopFrontRight", "TopBackLeft", "TopBackRight",
	"BottomFrontLeft", "BottomFrontRight", "BottomBackLeft", "BottomBackRight",
}

// Anchors lists every anchor a new step could use: the anchor steps, with
// their normals, and the bounding-box anchors of every named solid.
func (r *Result) Anchors() []AnchorInfo {
	var out []AnchorInfo
	for _, s := range r.Steps {
		if s.Name == "" || !s.value.IsValid() {
			continue
		}
		if s.Anchor != nil {
			n := s.Anchor.Normal
			out = append(out, AnchorInfo{Ref: s.Name, Step: s.Index, Pos: s.Anchor.Pos, Normal: &n})
			continue
		}
		sol := asSolid(s.value)
		if sol == nil {
			continue
		}
		for _, a := range bodyAnchors {
			v := reflect.ValueOf(sol).MethodByName(a).Call(nil)[0].Interface().(solid.AnchoredSolid)
			out = append(out, AnchorInfo{Ref: s.Name + "." + a, Step: s.Index, Pos: v.Point})
		}
	}
	return out
}
