package main

import (
	"fmt"
	"os"

	"github.com/snowbldr/fluent-sdfx/history"
	"gopkg.in/yaml.v3"
)

// legacy reads a history written by the first version of this tool (a
// base form and a list of pulls) and returns it as a history, so a file
// made before the runtime existed still opens.
func legacy(path string) (*history.History, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var old struct {
		Base  string `yaml:"base"`
		Steps []struct {
			Op     string    `yaml:"op"`
			At     []float64 `yaml:"at"`
			By     []float64 `yaml:"by"`
			Radius float64   `yaml:"radius"`
			Off    bool      `yaml:"off"`
		} `yaml:"steps"`
	}
	dec := yaml.NewDecoder(bytesReader(b))
	dec.KnownFields(true)
	if err := dec.Decode(&old); err != nil || old.Base == "" {
		return nil, fmt.Errorf("not a legacy history")
	}
	h, err := starter(old.Base)
	if err != nil {
		return nil, err
	}
	for _, st := range old.Steps {
		if st.Op != "pull" {
			return nil, fmt.Errorf("legacy history: unknown op %q", st.Op)
		}
		h.Steps = append(h.Steps, history.Step{
			Op:   "sculpt.Pull",
			Args: map[string]any{"at": floats(st.At), "by": floats(st.By), "radius": st.Radius},
			Off:  st.Off,
		})
	}
	return h, nil
}

func floats(xs []float64) []any {
	out := make([]any, len(xs))
	for i, x := range xs {
		out[i] = x
	}
	return out
}
