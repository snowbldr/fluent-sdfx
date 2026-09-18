package solid

import (
	"math"
	"testing"

	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
)

// surfaceHalfWidth finds, by bisection along +x at height z, where the
// frustum's surface actually sits.
func surfaceHalfWidth(s *Solid, z float64) float64 {
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

// The whole reason PyramidFrustum exists: the half-width must be exactly
// linear in z. A Loft of two squares lerps distance fields instead and
// pulls the corners in, which is what this guards against.
func TestPyramidFrustumHalfWidthIsExactlyLinear(t *testing.T) {
	const halfBase, halfTop, h = 10.0, 4.0, 12.0
	f := PyramidFrustum(halfBase, halfTop, h)

	// Strictly inside: at exactly z=0 or z=h the probe sits on the end
	// face and the bisection has nothing to bracket.
	for _, frac := range []float64{0.001, 0.25, 0.5, 0.75, 0.999} {
		z := frac * h
		want := halfBase + (halfTop-halfBase)*frac
		got := surfaceHalfWidth(f, z)
		if math.Abs(got-want) > 1e-6 {
			t.Errorf("half-width at z=%.2f (%.0f%%) = %.9f, want %.9f", z, frac*100, got, want)
		}
	}
}

// Corners are where a lerped field loses the most, so check the diagonal
// too: at height z the corner must sit at halfWidth*sqrt(2) from the axis.
func TestPyramidFrustumCornersAreSquare(t *testing.T) {
	const halfBase, halfTop, h = 8.0, 3.0, 10.0
	f := PyramidFrustum(halfBase, halfTop, h)

	for _, frac := range []float64{0.1, 0.5, 0.9} {
		z := frac * h
		hw := halfBase + (halfTop-halfBase)*frac
		// Just inside the corner must be solid; just outside must be air.
		in := v3sdf.Vec{X: hw - 1e-3, Y: hw - 1e-3, Z: z}
		out := v3sdf.Vec{X: hw + 1e-3, Y: hw + 1e-3, Z: z}
		if d := f.SDF3.Evaluate(in); d > 0 {
			t.Errorf("corner just inside at z=%.2f reads %.9f, want <= 0", z, d)
		}
		if d := f.SDF3.Evaluate(out); d < 0 {
			t.Errorf("corner just outside at z=%.2f reads %.9f, want >= 0", z, d)
		}
	}
}

func TestPyramidFrustumBounds(t *testing.T) {
	f := PyramidFrustum(10, 4, 12)
	b := f.Bounds()
	if math.Abs(b.Min.Z) > 1e-9 {
		t.Errorf("base should sit at z=0, got %v", b.Min.Z)
	}
	if math.Abs(b.Max.Z-12) > 1e-9 {
		t.Errorf("top should sit at z=h, got %v", b.Max.Z)
	}
}

// A degenerate frustum with equal radii is a plain box.
func TestPyramidFrustumEqualRadiiIsABox(t *testing.T) {
	f := PyramidFrustum(5, 5, 10)
	box := Box(v3.XYZ(10, 10, 10), 0).BottomAt(0)
	for _, z := range []float64{0.5, 5, 9.5} {
		got, want := surfaceHalfWidth(f, z), surfaceHalfWidth(box, z)
		if math.Abs(got-want) > 1e-6 {
			t.Errorf("z=%.1f: frustum %.9f, box %.9f", z, got, want)
		}
	}
}

// A frustum whose top is wider than its base must actually widen. An
// earlier version sized its blank by halfBase alone, so PyramidFrustum(4, 10,
// 12) came out as an 8 mm prism with no error.
func TestPyramidFrustumWidens(t *testing.T) {
	f := PyramidFrustum(4, 10, 12)
	if b := f.Bounds(); b.Max.X < 10-1e-9 {
		t.Fatalf("widening frustum clipped to %v", b)
	}
	in := func(x, z float64) bool { return f.Evaluate(v3sdf.Vec{X: x, Y: 0, Z: z}) < 0 }
	if !in(9, 11.5) {
		t.Error("(9, 0, 11.5) should be inside the wide top")
	}
	if in(9, 0.5) {
		t.Error("(9, 0, 0.5) should be outside the narrow base")
	}
	// Half-width is linear in z: at z = 6 it is 7.
	if !in(6.9, 6) || in(7.1, 6) {
		t.Error("half-width at z=6 is not 7")
	}
}

func TestPyramidFrustumRejectsNonPositive(t *testing.T) {
	for _, c := range [][3]float64{{0, 5, 10}, {5, 0, 10}, {5, 5, 0}, {-1, 5, 10}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("PyramidFrustum%v did not panic", c)
				}
			}()
			PyramidFrustum(c[0], c[1], c[2])
		}()
	}
}
