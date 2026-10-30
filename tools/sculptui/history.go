package main

import (
	"fmt"
	"os"

	"github.com/snowbldr/fluent-sdfx/sculpt"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"gopkg.in/yaml.v3"
)

// History is the whole model as data: a base form and the steps applied
// to it, in order. It is the file on disk, the thing the UI edits, and the
// only state the server keeps. Deleting a step and re-running is undo.
type History struct {
	Base  string `yaml:"base" json:"base"`   // figure, sphere or box
	Steps []Step `yaml:"steps" json:"steps"` // applied in order
}

// Step is one operation. Only pull exists in this test loop.
type Step struct {
	Op     string    `yaml:"op" json:"op"`
	At     []float64 `yaml:"at,flow" json:"at"`
	By     []float64 `yaml:"by,flow" json:"by"`
	Radius float64   `yaml:"radius" json:"radius"`
	Off    bool      `yaml:"off,omitempty" json:"off,omitempty"` // skipped when true: the step stays in the file
}

func newHistory() History { return History{Base: "figure"} }

func loadHistory(path string) (History, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return History{}, err
	}
	var h History
	if err := yaml.Unmarshal(b, &h); err != nil {
		return History{}, fmt.Errorf("%s: %w", path, err)
	}
	if h.Base == "" {
		h.Base = "figure"
	}
	return h, nil
}

const header = `# fluent-sdfx sculpt test loop.
# Written by the UI on every change; edit it by hand and the view reloads.
# base: figure | sphere | box. Each step runs in order; off: true skips one.
`

func (h History) save(path string) error {
	b, err := yaml.Marshal(h)
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append([]byte(header), b...), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func base(name string) (*solid.Solid, error) {
	switch name {
	case "figure":
		// The figure from sculpt's scenario tests: a body with a head
		// blended on top. -Y is the front.
		body := solid.Sphere(12)
		head := solid.Sphere(8).TranslateZ(14)
		return solid.SmoothUnion(solid.RoundMin(3), body, head), nil
	case "sphere":
		return solid.Sphere(12), nil
	case "box":
		return solid.Box(v3.XYZ(24, 24, 24), 3), nil
	}
	return nil, fmt.Errorf("unknown base %q (want figure, sphere or box)", name)
}

func vec(name string, xs []float64) (v3.Vec, error) {
	if len(xs) != 3 {
		return v3.Vec{}, fmt.Errorf("%s needs three numbers, got %d", name, len(xs))
	}
	return v3.XYZ(xs[0], xs[1], xs[2]), nil
}

// StepError is a step that failed, by index. The steps before it still
// built; it and the ones after it were not applied.
type StepError struct {
	Index int    `json:"index"`
	Msg   string `json:"msg"`
}

// build evaluates the history. It returns the solid from every step that
// succeeded and, if one failed, which one and why: the red-node behaviour,
// where everything up to the broken step still renders.
func (h History) build() (s *solid.Solid, stepErr *StepError, err error) {
	s, err = base(h.Base)
	if err != nil {
		return nil, nil, err
	}
	for i, st := range h.Steps {
		if st.Off {
			continue
		}
		next, e := applyStep(s, st)
		if e != nil {
			return s, &StepError{Index: i, Msg: e.Error()}, nil
		}
		s = next
	}
	return s, nil, nil
}

func applyStep(s *solid.Solid, st Step) (out *solid.Solid, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	switch st.Op {
	case "pull":
		at, err := vec("at", st.At)
		if err != nil {
			return nil, err
		}
		by, err := vec("by", st.By)
		if err != nil {
			return nil, err
		}
		return sculpt.Pull(s, at, by, st.Radius), nil
	}
	return nil, fmt.Errorf("unknown op %q", st.Op)
}
