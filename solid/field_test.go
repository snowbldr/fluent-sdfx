package solid

import (
	"math"
	"testing"

	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
)

// radiusAt finds the surface radius along +x at height z by bisection.
func radiusAt(s *Solid, z float64) float64 {
	lo, hi := 0.0, 100.0
	for i := 0; i < 80; i++ {
		mid := (lo + hi) / 2
		if s.SDF3.Evaluate(v3sdf.Vec{X: mid, Y: 0, Z: z}) < 0 {
			lo = mid
		} else {
			hi = mid
		}
	}
	return (lo + hi) / 2
}

// Below z0 a straightened cone must be a constant-radius prism, and that
// radius must equal the cone's own radius at z0.
func TestStraightenBelowIsPrismatic(t *testing.T) {
	// Cone: r=10 at the bottom (z=-10), r=2 at the top (z=+10).
	cone := Cone(20, 10, 2, 0)
	const z0 = 0.0
	want := radiusAt(cone, z0)

	st := cone.StraightenBelow(z0)
	for _, z := range []float64{-9, -5, -1, -0.001} {
		if got := radiusAt(st, z); math.Abs(got-want) > 1e-6 {
			t.Errorf("z=%.3f: radius %.9f, want the z0 radius %.9f", z, got, want)
		}
	}
}

// Above z0 the field must be untouched.
func TestStraightenBelowLeavesTopAlone(t *testing.T) {
	cone := Cone(20, 10, 2, 0)
	st := cone.StraightenBelow(0)
	for _, z := range []float64{0.5, 3, 7, 9.5} {
		a, b := radiusAt(cone, z), radiusAt(st, z)
		if math.Abs(a-b) > 1e-9 {
			t.Errorf("z=%.1f: straightened %.9f != original %.9f", z, b, a)
		}
	}
}

// The prism must stop at the solid's own bottom, not run to infinity.
func TestStraightenBelowFloorsAtOriginalBottom(t *testing.T) {
	cone := Cone(20, 10, 2, 0)
	bot := cone.Bounds().Min.Z
	st := cone.StraightenBelow(0)
	if d := st.SDF3.Evaluate(v3sdf.Vec{X: 0, Y: 0, Z: bot - 1}); d <= 0 {
		t.Errorf("a point 1mm below the original bottom reads %.9f, want > 0 (air)", d)
	}
	if d := st.SDF3.Evaluate(v3sdf.Vec{X: 0, Y: 0, Z: bot + 1}); d >= 0 {
		t.Errorf("a point 1mm above the original bottom reads %.9f, want < 0 (solid)", d)
	}
}

// A z0 at or below the solid's bottom is a no-op on the visible shape.
func TestStraightenBelowAtBottomIsNoOp(t *testing.T) {
	cone := Cone(20, 10, 2, 0)
	st := cone.StraightenBelow(cone.Bounds().Min.Z)
	for _, z := range []float64{-9, 0, 9} {
		a, b := radiusAt(cone, z), radiusAt(st, z)
		if math.Abs(a-b) > 1e-6 {
			t.Errorf("z=%.1f: %.9f != %.9f", z, b, a)
		}
	}
}

func TestWithBoundsEnlarges(t *testing.T) {
	s := Box(v3.XYZ(2, 2, 2), 0)
	big := NewBox3(v3.XYZ(0, 0, 0), v3.XYZ(10, 10, 10))
	got := s.WithBounds(big).Bounds()
	if math.Abs(got.Min.X+5) > 1e-9 || math.Abs(got.Max.X-5) > 1e-9 {
		t.Errorf("bounds = %v, want the supplied ±5 box", got)
	}
}

// The field must be completely unchanged — WithBounds reports, it does not
// reshape.
func TestWithBoundsLeavesFieldUntouched(t *testing.T) {
	s := Box(v3.XYZ(2, 2, 2), 0)
	w := s.WithBounds(NewBox3(v3.XYZ(0, 0, 0), v3.XYZ(10, 10, 10)))
	for _, p := range []v3sdf.Vec{
		{X: 0, Y: 0, Z: 0}, {X: 0.9, Y: 0, Z: 0}, {X: 1.1, Y: 0, Z: 0}, {X: 4, Y: 4, Z: 4},
	} {
		if a, b := s.SDF3.Evaluate(p), w.SDF3.Evaluate(p); a != b {
			t.Errorf("at %v: %.17g != %.17g", p, b, a)
		}
	}
}

// Asking for a box smaller than the geometry would silently clip the
// render, so the result is the union, never the shrunken box.
func TestWithBoundsNeverShrinksBelowTheSurface(t *testing.T) {
	s := Box(v3.XYZ(10, 10, 10), 0)
	got := s.WithBounds(NewBox3(v3.XYZ(0, 0, 0), v3.XYZ(2, 2, 2))).Bounds()
	if math.Abs(got.Min.X+5) > 1e-9 || math.Abs(got.Max.X-5) > 1e-9 {
		t.Errorf("bounds = %v, want the solid's own ±5 box preserved", got)
	}
}
