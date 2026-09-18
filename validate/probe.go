package validate

import (
	"fmt"
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Probe is one named point in the solid's own coordinate frame together with
// the material state expected there: Inside true means At(s, P) must be
// negative (inside material), false means it must be positive (air).
type Probe struct {
	Name   string  // what the point means, e.g. "bore wall at mid height"
	P      v3.Vec  // the point, in the solid's frame, mm
	Inside bool    // true: expect At(s, P) < 0; false: expect At(s, P) > 0
	Margin float64 // optional: also require |At(s, P)| >= Margin mm (0 = any distance)
}

// ProbeResult is the outcome of one Probe: the signed distance found there
// (negative inside material) and whether it met the expectation, including
// the Margin when one was given.
type ProbeResult struct {
	Probe
	SDF    float64 // the signed distance at P, mm; NaN if the solid was missing
	Passed bool
}

// String returns the one-line report for this probe, e.g.
// "[ok] bore wall (3.000, 0.000, 5.000) want=inside got=inside sdf=-0.312",
// with a "margin=" suffix when the probe asked for one.
func (r ProbeResult) String() string {
	if math.IsNaN(r.SDF) {
		return fmt.Sprintf("[FAIL] %s (%.3f, %.3f, %.3f) want=%s got=? sdf=NaN (no such solid)",
			r.Name, r.P.X, r.P.Y, r.P.Z, probeState(r.Inside))
	}
	tag := "[FAIL]"
	if r.Passed {
		tag = "[ok]"
	}
	line := fmt.Sprintf("%s %s (%.3f, %.3f, %.3f) want=%s got=%s sdf=%+.3f",
		tag, r.Name, r.P.X, r.P.Y, r.P.Z, probeState(r.Inside), probeState(r.SDF < 0), r.SDF)
	if r.Margin > 0 {
		line += fmt.Sprintf(" margin=%.3f", r.Margin)
	}
	return line
}

// Probes evaluates every probe against s in parallel and returns one result
// per probe, in the same order. A probe passes when the sign of the field
// matches Probe.Inside (negative is inside material) and, if Margin is
// non-zero, |sdf| >= Margin.
func Probes(s *solid.Solid, probes []Probe) []ProbeResult {
	pts := make([]v3.Vec, len(probes))
	for i, p := range probes {
		pts[i] = p.P
	}
	f := evalAll(s, pts)
	out := make([]ProbeResult, len(probes))
	for i, p := range probes {
		out[i] = ProbeResult{Probe: p, SDF: f[i], Passed: probePassed(p, f[i])}
	}
	return out
}

// RequireProbes evaluates probes against s, logs one report line per probe
// plus an "N/M probes passed" summary, and calls t.Errorf for every probe
// whose material state (or margin) was wrong.
func RequireProbes(t testing.TB, s *solid.Solid, probes []Probe) {
	t.Helper()
	probeReport(t, Probes(s, probes))
}

// NamedProbe is a Probe that also names which solid of a set it applies to,
// for tables that check several parts of an assembly at once. Solid is the
// map key; Name is the point's own name.
type NamedProbe struct {
	Name   string  // what the point means
	Solid  string  // key into the solids map passed to ProbesIn
	P      v3.Vec  // the point, in that solid's frame, mm
	Inside bool    // true: expect At(solids[Solid], P) < 0
	Margin float64 // optional: also require |sdf| >= Margin mm
}

// ProbesIn evaluates each NamedProbe against solids[probe.Solid] and returns
// results whose Name is "solid:name". A probe naming a solid that is not in
// the map fails with SDF = NaN rather than panicking.
func ProbesIn(solids map[string]*solid.Solid, probes []NamedProbe) []ProbeResult {
	out := make([]ProbeResult, len(probes))
	for i, np := range probes {
		out[i] = ProbeResult{
			Probe: Probe{Name: np.Solid + ":" + np.Name, P: np.P, Inside: np.Inside, Margin: np.Margin},
			SDF:   math.NaN(),
		}
	}
	parallel(len(probes), func(i int) {
		s, ok := solids[probes[i].Solid]
		if !ok || s == nil {
			return
		}
		f := At(s, probes[i].P)
		out[i].SDF = f
		out[i].Passed = probePassed(out[i].Probe, f)
	})
	return out
}

// RequireProbesIn evaluates a multi-solid probe table, logs one report line
// per probe plus an "N/M probes passed" summary, and calls t.Errorf for every
// failure. Result names are "solid:name".
func RequireProbesIn(t testing.TB, solids map[string]*solid.Solid, probes []NamedProbe) {
	t.Helper()
	probeReport(t, ProbesIn(solids, probes))
}

// In returns a Probe expecting p to be inside material (At < 0).
func In(name string, p v3.Vec) Probe { return Probe{Name: name, P: p, Inside: true} }

// Out returns a Probe expecting p to be in air (At > 0).
func Out(name string, p v3.Vec) Probe { return Probe{Name: name, P: p, Inside: false} }

// probePassed reports whether a field reading meets a probe's expectation.
func probePassed(p Probe, sdf float64) bool {
	if math.IsNaN(sdf) {
		return false
	}
	if (sdf < 0) != p.Inside {
		return false
	}
	if p.Margin > 0 && math.Abs(sdf) < p.Margin {
		return false
	}
	return true
}

// probeState names a material state for the report line.
func probeState(inside bool) string {
	if inside {
		return "inside"
	}
	return "outside"
}

// probeReport logs every result, then the summary, then errors per failure.
func probeReport(t testing.TB, rs []ProbeResult) {
	t.Helper()
	passed := 0
	for _, r := range rs {
		t.Log(r.String())
		if r.Passed {
			passed++
		}
	}
	t.Logf("%d/%d probes passed", passed, len(rs))
	for _, r := range rs {
		if !r.Passed {
			t.Errorf("%s", r.String())
		}
	}
}
