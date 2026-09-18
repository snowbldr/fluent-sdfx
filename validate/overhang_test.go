package validate

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// ohFakeTB captures what a Require* helper reports without failing the real
// test, so the failing path can be asserted like any other result.
type ohFakeTB struct {
	testing.TB
	logs   []string
	errors []string
}

func (f *ohFakeTB) Helper()         {}
func (f *ohFakeTB) Log(args ...any) { f.logs = append(f.logs, fmt.Sprint(args...)) }
func (f *ohFakeTB) Logf(format string, args ...any) {
	f.logs = append(f.logs, fmt.Sprintf(format, args...))
}
func (f *ohFakeTB) Errorf(format string, args ...any) {
	f.errors = append(f.errors, fmt.Sprintf(format, args...))
}

// ohBox is a box sitting on the plate, centred on the Z axis in X and Y.
func ohBox(sx, sy, sz, z0 float64) *solid.Solid {
	return solid.Box(v3.XYZ(sx, sy, sz), 0).Translate(v3.XYZ(0, 0, z0+sz/2))
}

// ohT is the canonical failing part: a 6 mm stem from the plate to z=10
// carrying a 16 mm cross-bar from z=10 to z=14, so 5 mm of bar hangs
// unsupported on each side.
func ohT() *solid.Solid {
	return ohBox(6, 6, 10, 0).Union(ohBox(16, 6, 4, 10))
}

// ohGusset is a 45° triangular gusset under one side of ohT's cross-bar:
// the block x∈[3,8], z∈[5,10] cut by the plane through (3,0,5) with normal
// (-1,0,1), which keeps the material above the 45° diagonal.
func ohGusset(sign float64) *solid.Solid {
	blk := solid.Box(v3.XYZ(5, 6, 5), 0).Translate(v3.XYZ(sign*5.5, 0, 7.5))
	return blk.CutPlane(v3.XYZ(sign*3, 0, 5), v3.XYZ(-sign, 0, 1))
}

func TestOverhangsBoxHasNone(t *testing.T) {
	// A plain box: the only downward face is its bottom, which is on the
	// plate and exempt. The worst downward angle is still 90°.
	r := Overhangs(ohBox(20, 20, 10, 0), 6, OverhangOptions{})
	if r.Faces != 0 || r.Area != 0 {
		t.Errorf("box: %d faces / %.3f mm² unsupported, want none\n%s", r.Faces, r.Area, r)
	}
	if math.Abs(r.WorstDeg-90) > 1 {
		t.Errorf("box worst downward angle = %.2f°, want 90", r.WorstDeg)
	}
	if r.WorstAt.Z > 0.5 {
		t.Errorf("box worst downward face at z=%.3f, want the bottom face (z≈0)", r.WorstAt.Z)
	}
}

func TestOverhangsTCrossBarFails(t *testing.T) {
	// The bar underside is a 90° ceiling from |x|=3 (the stem) out to
	// |x|=8. The lateral probe forgives 0.8 mm of it on each side, so
	// roughly 2*(8-3.8)*6 ≈ 50 mm² should be reported.
	r := Overhangs(ohT(), 8, OverhangOptions{})
	t.Log(r)
	if r.Faces == 0 {
		t.Fatalf("T cross-bar underside reported as supported: %s", r)
	}
	if r.Area < 35 || r.Area > 65 {
		t.Errorf("failing area = %.2f mm², want ≈50 (2*(8-3.8)*6)", r.Area)
	}
	if math.Abs(r.WorstDeg-90) > 1 {
		t.Errorf("worst downward angle = %.2f°, want 90 (flat ceiling)", r.WorstDeg)
	}
	worst, at := ohWorstFailing(r.Failing)
	if math.Abs(worst-90) > 1 || math.Abs(at.Z-10) > 0.5 {
		t.Errorf("worst failing face = %.2f° at z=%.3f, want 90° at z≈10", worst, at.Z)
	}
	if !strings.Contains(r.String(), "unsupported faces over threshold") {
		t.Errorf("String() = %q", r.String())
	}
}

func TestOverhangsTWithGussetPasses(t *testing.T) {
	// The same T with a 45° gusset under each end of the bar: every
	// downward face is now either at threshold or has material 1 mm below
	// it, so nothing is reported.
	s := ohT().Union(ohGusset(1), ohGusset(-1))
	r := Overhangs(s, 8, OverhangOptions{})
	t.Log(r)
	if r.Area > 1 {
		t.Errorf("gusseted T: %.3f mm² unsupported over %d faces, want ≈0\n%s", r.Area, r.Faces, r)
	}
	if r.WorstDeg < 40 {
		t.Errorf("worst downward angle = %.2f°, want ≥45 (the gusset undersides are 45° and the bar ends are steeper)", r.WorstDeg)
	}
}

func TestOverhangsNearBedExemption(t *testing.T) {
	// A 0.3 mm tunnel through the bottom of a block: its ceiling is a flat
	// 90° face, exempt because it is within BedTol of the plate.
	block := ohBox(20, 10, 5, 0)
	tunnel := solid.Box(v3.XYZ(6, 12, 0.3), 0).Translate(v3.XYZ(0, 0, 0.15))
	s := block.Cut(tunnel)

	if r := Overhangs(s, 8, OverhangOptions{}); r.Faces != 0 {
		t.Errorf("0.3 mm tunnel ceiling reported: %s", r)
	}
	// With both the near-bed exemption and the plate-reaching probe turned
	// off (BedTol < 0), the same ceiling is a genuine unsupported face.
	r := Overhangs(s, 8, OverhangOptions{BedTol: -1})
	if r.Area < 20 {
		t.Errorf("with BedTol disabled the tunnel ceiling should be reported, got %.3f mm² over %d faces", r.Area, r.Faces)
	}
}

func TestOverhangsThresholdAndOptions(t *testing.T) {
	s := ohT()
	cases := []struct {
		name string
		o    OverhangOptions
		want bool // want unsupported area reported
	}{
		{"default 45°", OverhangOptions{}, true},
		{"89° threshold still catches the 90° ceiling", OverhangOptions{MaxDeg: 89}, true},
		{"per-height threshold exempts the bar", OverhangOptions{Threshold: func(z float64) float64 {
			if z > 9 {
				return 90
			}
			return 45
		}}, false},
		{"deep support probe finds the stem", OverhangOptions{SupportDepth: 12}, false},
		{"wide lateral probe finds the stem", OverhangOptions{Lateral: 6}, false},
		{"no lateral probe reports slightly more", OverhangOptions{Lateral: -1}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := Overhangs(s, 8, tc.o)
			if got := r.Area > 0.5; got != tc.want {
				t.Errorf("area = %.3f mm² over %d faces, want reported=%v", r.Area, r.Faces, tc.want)
			}
		})
	}

	// Disabling the lateral probe can only add failing area, never remove it.
	wide := Overhangs(s, 8, OverhangOptions{})
	narrow := Overhangs(s, 8, OverhangOptions{Lateral: -1})
	if narrow.Area < wide.Area {
		t.Errorf("straight-down-only area %.3f mm² < with-lateral area %.3f mm²", narrow.Area, wide.Area)
	}
}

func TestRequireOverhangArea(t *testing.T) {
	// Positive: a gusseted T inside a 1 mm² budget must not fail.
	good := ohT().Union(ohGusset(1), ohGusset(-1))
	RequireOverhangArea(t, good, 8, OverhangOptions{}, 1)

	// Negative: the bare T must fail the same budget and log the report.
	f := &ohFakeTB{}
	RequireOverhangArea(f, ohT(), 8, OverhangOptions{}, 1)
	if len(f.errors) != 1 {
		t.Fatalf("errors = %v, want exactly one", f.errors)
	}
	if !strings.Contains(f.errors[0], "exceeds budget") {
		t.Errorf("error = %q, want it to name the budget", f.errors[0])
	}
	if len(f.logs) != 1 || !strings.Contains(f.logs[0], "overhangs:") {
		t.Errorf("logs = %v, want the report line", f.logs)
	}
}

func TestOverhangReportStringNoFailures(t *testing.T) {
	r := OverhangReport{WorstDeg: 44.2, WorstAt: v3.XYZ(1, 2, 3)}
	if s := r.String(); !strings.Contains(s, "none unsupported over threshold") || !strings.Contains(s, "44.2°") {
		t.Errorf("String() = %q", s)
	}
}
