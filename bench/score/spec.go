// Package score compares a candidate solid against a reference solid and
// reports staged metrics: compiled, built, bounding box, volumetric IoU,
// surface deviation, intent probes, and (optionally) mesh validity.
//
// It is the scoring half of the description-to-geometry benchmark in
// bench/. See bench/README.md for the task format.
package score

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Probe is a point in the reference frame with the expected material state.
// SDF < 0 is inside material.
type Probe struct {
	Name   string     `json:"name"`
	P      [3]float64 `json:"p"`
	Inside bool       `json:"inside"`
}

// Pass holds the thresholds that decide whether a candidate passes.
type Pass struct {
	IoUMin   float64 `json:"iou_min"`    // default 0.97
	MaxDevMM float64 `json:"max_dev_mm"` // default 0.3
	ExtraMax float64 `json:"extra_max"`  // unrequested volume as a fraction of reference volume; default 0.05
}

// Spec is the machine-readable half of a task (task.json).
type Spec struct {
	Name     string   `json:"name"`
	Category string   `json:"category"` // vocabulary, assumption, orientation, anchor, api-trap, printability, verification, over-engineering, multi-ask
	Tags     []string `json:"tags"`
	Mode     string   `json:"mode"`   // build | modify | clarify
	Source   string   `json:"source"` // incident pointer, e.g. "mcweed U209-U225 (gang-boss length)"
	Summary  string   `json:"summary"`

	// ExpectsQuestion marks a task whose ideal first response is a
	// clarifying question. The runner scores whether the model asked.
	ExpectsQuestion bool `json:"expects_question"`

	// Align controls how the candidate is registered to the reference
	// before comparison: "none" (shared frame; default for modify tasks)
	// or "center" (translate candidate so bbox centers coincide; default
	// for build tasks where the description does not fix an origin).
	Align string `json:"align"`

	// Grid is the number of sample cells along the longest axis of the
	// union bounding box. Default 96.
	Grid int `json:"grid"`

	// MeshCellsPerMM > 0 enables the mesh stage (watertight, volume,
	// connected components) at that resolution. Default 0 (skipped).
	MeshCellsPerMM float64 `json:"mesh_cells_per_mm"`

	Probes []Probe `json:"probes"`
	Pass   Pass    `json:"pass"`
}

// Defaults fills unset fields.
func (s *Spec) Defaults() {
	if s.Align == "" {
		if s.Mode == "build" {
			s.Align = "center"
		} else {
			s.Align = "none"
		}
	}
	if s.Grid == 0 {
		s.Grid = 96
	}
	if s.Pass.IoUMin == 0 {
		s.Pass.IoUMin = 0.97
	}
	if s.Pass.MaxDevMM == 0 {
		s.Pass.MaxDevMM = 0.3
	}
	if s.Pass.ExtraMax == 0 {
		s.Pass.ExtraMax = 0.05
	}
}

// LoadSpec reads task.json from a task directory (or a direct path).
func LoadSpec(path string) (Spec, error) {
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		path = filepath.Join(path, "task.json")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return Spec{}, err
	}
	var s Spec
	if err := json.Unmarshal(b, &s); err != nil {
		return Spec{}, fmt.Errorf("%s: %w", path, err)
	}
	if s.Name == "" {
		s.Name = filepath.Base(filepath.Dir(path))
	}
	s.Defaults()
	return s, nil
}
