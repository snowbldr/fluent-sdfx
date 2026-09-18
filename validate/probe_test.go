package validate

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// probeTB captures what a Require* helper reports without failing the real
// test, so the negative cases can be asserted.
type probeTB struct {
	testing.TB
	logs []string
	errs []string
}

func (r *probeTB) Helper()                    {}
func (r *probeTB) Log(args ...any)            { r.logs = append(r.logs, fmt.Sprint(args...)) }
func (r *probeTB) Logf(f string, args ...any) { r.logs = append(r.logs, fmt.Sprintf(f, args...)) }
func (r *probeTB) Error(args ...any)          { r.errs = append(r.errs, fmt.Sprint(args...)) }
func (r *probeTB) Errorf(format string, args ...any) {
	r.errs = append(r.errs, fmt.Sprintf(format, args...))
}
func (r *probeTB) has(sub string, in []string) bool {
	for _, s := range in {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// A 10 mm cube centred on the origin: material for |x|,|y|,|z| < 5.
func probeCube() *solid.Solid { return solid.Box(v3.XYZ(10, 10, 10), 0) }

// The sign convention this whole package rests on: negative inside.
func TestProbeSignConvention(t *testing.T) {
	cube := probeCube()
	if got := At(cube, v3.Vec{}); got >= 0 {
		t.Fatalf("At at the centre of a box = %+.3f, want negative (inside material)", got)
	}
	if got := At(cube, v3.XYZ(20, 0, 0)); got <= 0 {
		t.Fatalf("At well outside a box = %+.3f, want positive (air)", got)
	}
}

func TestProbes(t *testing.T) {
	cube := probeCube()
	tests := []struct {
		probe      Probe
		wantPass   bool
		wantInside bool // expected sign of the reading, negative = inside
	}{
		{In("centre", v3.Vec{}), true, true},
		{In("just inside +X face", v3.XYZ(4.9, 0, 0)), true, true},
		{Out("just outside +X face", v3.XYZ(5.1, 0, 0)), true, false},
		{Out("beyond the corner", v3.XYZ(6, 6, 6)), true, false},
		{In("wrongly expected inside", v3.XYZ(5.1, 0, 0)), false, false},
		{Out("wrongly expected outside", v3.Vec{}), false, true},
		{Probe{Name: "deep inside, 1 mm margin", P: v3.Vec{}, Inside: true, Margin: 1}, true, true},
		{Probe{Name: "barely inside, 1 mm margin", P: v3.XYZ(4.9, 0, 0), Inside: true, Margin: 1}, false, true},
		{Probe{Name: "barely outside, 1 mm margin", P: v3.XYZ(5.1, 0, 0), Inside: false, Margin: 1}, false, false},
		{Probe{Name: "far outside, 1 mm margin", P: v3.XYZ(10, 0, 0), Inside: false, Margin: 1}, true, false},
	}
	probes := make([]Probe, len(tests))
	for i, tc := range tests {
		probes[i] = tc.probe
	}
	got := Probes(cube, probes)
	if len(got) != len(tests) {
		t.Fatalf("Probes returned %d results, want %d", len(got), len(tests))
	}
	for i, tc := range tests {
		r := got[i]
		if r.Name != tc.probe.Name {
			t.Errorf("result %d name = %q, want %q (order must be preserved)", i, r.Name, tc.probe.Name)
		}
		if r.Passed != tc.wantPass {
			t.Errorf("%s: Passed = %v, want %v (%s)", tc.probe.Name, r.Passed, tc.wantPass, r)
		}
		if (r.SDF < 0) != tc.wantInside {
			t.Errorf("%s: sdf = %+.3f, want %s", tc.probe.Name, r.SDF, probeState(tc.wantInside))
		}
	}
}

func TestProbeResultString(t *testing.T) {
	tests := []struct {
		r    ProbeResult
		want string
	}{
		{
			ProbeResult{Probe: Probe{Name: "bore wall", P: v3.XYZ(1, 2, 3), Inside: true}, SDF: -0.312, Passed: true},
			"[ok] bore wall (1.000, 2.000, 3.000) want=inside got=inside sdf=-0.312",
		},
		{
			ProbeResult{Probe: Probe{Name: "bore wall", P: v3.XYZ(1, 2, 3), Inside: true}, SDF: 0.045, Passed: false},
			"[FAIL] bore wall (1.000, 2.000, 3.000) want=inside got=outside sdf=+0.045",
		},
		{
			ProbeResult{Probe: Probe{Name: "clearance", P: v3.Vec{}, Inside: false, Margin: 0.2}, SDF: 0.5, Passed: true},
			"[ok] clearance (0.000, 0.000, 0.000) want=outside got=outside sdf=+0.500 margin=0.200",
		},
		{
			ProbeResult{Probe: Probe{Name: "lid:centre", P: v3.XYZ(-1, 0, 0), Inside: true}, SDF: math.NaN()},
			"[FAIL] lid:centre (-1.000, 0.000, 0.000) want=inside got=? sdf=NaN (no such solid)",
		},
	}
	for _, tc := range tests {
		if got := tc.r.String(); got != tc.want {
			t.Errorf("String() =\n  %q\nwant\n  %q", got, tc.want)
		}
	}
}

func TestRequireProbes(t *testing.T) {
	cube := probeCube()
	t.Run("all pass", func(t *testing.T) {
		rec := &probeTB{}
		RequireProbes(rec, cube, []Probe{In("centre", v3.Vec{}), Out("outside", v3.XYZ(9, 0, 0))})
		if len(rec.errs) != 0 {
			t.Errorf("unexpected failures: %v", rec.errs)
		}
		if !rec.has("2/2 probes passed", rec.logs) {
			t.Errorf("missing summary line in %v", rec.logs)
		}
	})
	t.Run("one fails", func(t *testing.T) {
		rec := &probeTB{}
		RequireProbes(rec, cube, []Probe{In("centre", v3.Vec{}), In("wrong", v3.XYZ(9, 0, 0))})
		if len(rec.errs) != 1 {
			t.Errorf("errors = %v, want exactly one", rec.errs)
		}
		if !rec.has("1/2 probes passed", rec.logs) {
			t.Errorf("missing summary line in %v", rec.logs)
		}
		if !rec.has("[FAIL] wrong", rec.errs) {
			t.Errorf("failure line not reported: %v", rec.errs)
		}
	})
}

func TestProbesIn(t *testing.T) {
	solids := map[string]*solid.Solid{
		"cube":     probeCube(),
		"cylinder": solid.Cylinder(20, 5, 0),
	}
	tests := []struct {
		probe    NamedProbe
		wantName string
		wantPass bool
	}{
		{NamedProbe{Name: "centre", Solid: "cube", P: v3.Vec{}, Inside: true}, "cube:centre", true},
		{NamedProbe{Name: "corner", Solid: "cube", P: v3.XYZ(4.9, 4.9, 4.9), Inside: false}, "cube:corner", false},
		{NamedProbe{Name: "on axis", Solid: "cylinder", P: v3.XYZ(0, 0, 9), Inside: true}, "cylinder:on axis", true},
		{NamedProbe{Name: "above end", Solid: "cylinder", P: v3.XYZ(0, 0, 11), Inside: false}, "cylinder:above end", true},
		{NamedProbe{Name: "centre", Solid: "missing", P: v3.Vec{}, Inside: true}, "missing:centre", false},
	}
	probes := make([]NamedProbe, len(tests))
	for i, tc := range tests {
		probes[i] = tc.probe
	}
	got := ProbesIn(solids, probes)
	for i, tc := range tests {
		if got[i].Name != tc.wantName {
			t.Errorf("result %d name = %q, want %q", i, got[i].Name, tc.wantName)
		}
		if got[i].Passed != tc.wantPass {
			t.Errorf("%s: Passed = %v, want %v (%s)", tc.wantName, got[i].Passed, tc.wantPass, got[i])
		}
	}
	if !math.IsNaN(got[len(got)-1].SDF) {
		t.Errorf("probe on a missing solid: sdf = %v, want NaN", got[len(got)-1].SDF)
	}
}

func TestRequireProbesIn(t *testing.T) {
	solids := map[string]*solid.Solid{"cube": probeCube()}
	rec := &probeTB{}
	RequireProbesIn(rec, solids, []NamedProbe{
		{Name: "centre", Solid: "cube", P: v3.Vec{}, Inside: true},
		{Name: "centre", Solid: "lid", P: v3.Vec{}, Inside: true},
	})
	if len(rec.errs) != 1 {
		t.Errorf("errors = %v, want exactly one (the missing solid)", rec.errs)
	}
	if !rec.has("1/2 probes passed", rec.logs) {
		t.Errorf("missing summary line in %v", rec.logs)
	}
	if !rec.has("no such solid", rec.errs) {
		t.Errorf("missing-solid failure not explained: %v", rec.errs)
	}
}

func TestInOut(t *testing.T) {
	p := In("a", v3.X(1))
	if !p.Inside || p.Name != "a" || p.P != v3.X(1) || p.Margin != 0 {
		t.Errorf("In = %+v, want Inside true at (1,0,0)", p)
	}
	if q := Out("b", v3.Y(2)); q.Inside || q.Name != "b" || q.P != v3.Y(2) {
		t.Errorf("Out = %+v, want Inside false at (0,2,0)", q)
	}
}
