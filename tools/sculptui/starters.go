package main

import (
	"fmt"

	"github.com/snowbldr/fluent-sdfx/history"
)

const header = `# fluent-sdfx sculpt test loop: a history the UI writes on every change.
# Edit it by hand and the view reloads. Steps run in order; off: true skips one.
`

// starter returns the steps a new history begins with.
func starter(form string) (*history.History, error) {
	h := &history.History{Header: header}
	switch form {
	case "figure":
		// The figure from sculpt's scenario tests: a body with a head
		// blended on top. -Y is the front.
		h.Steps = []history.Step{
			{Name: "body", Op: "solid.Sphere", Args: map[string]any{"radius": 12}},
			{Name: "head", Op: "solid.Sphere", Args: map[string]any{"radius": 8},
				Then: []history.Call{{Op: "TranslateZ", Args: []any{14}}}},
			{Name: "figure", Op: "solid.SmoothUnion", Args: map[string]any{
				"min":    map[string]any{"op": "solid.RoundMin", "args": map[string]any{"k": 3}},
				"solids": []any{"body", "head"},
			}},
		}
	case "sphere":
		h.Steps = []history.Step{{Name: "ball", Op: "solid.Sphere", Args: map[string]any{"radius": 12}}}
	case "box":
		h.Steps = []history.Step{{Name: "block", Op: "solid.Box", Args: map[string]any{"size": []any{24, 24, 24}, "round": 3}}}
	default:
		return nil, fmt.Errorf("unknown form %q (want figure, sphere or box)", form)
	}
	return h, nil
}
