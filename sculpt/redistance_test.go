package sculpt

import (
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// The redistance is the primitive everything else leans on, so it is
// checked against the one field whose true distance is known everywhere: a
// sphere. Every corner of the redistanced grid, deep interior included,
// must match |p| − R to within a fraction of a cell.
func TestRedistanceMatchesExactSphere(t *testing.T) {
	const R = 10.0
	g := Sample(solid.Sphere(R), 64)
	h := g.CellSize()

	// Scale the magnitudes first so the test proves the values were
	// recomputed, not copied through. (A uniform scale keeps the
	// interpolated zero crossings where they are; a non-uniform one would
	// move the surface the samples describe.)
	scrambled := g.clone()
	for i, v := range scrambled.vals {
		scrambled.vals[i] = 3 * v
	}
	r := scrambled.Redistance()

	maxErr, sumErr, deepErr := 0.0, 0.0, 0.0
	for id, v := range r.vals {
		want := r.posOf(id).Length() - R
		e := math.Abs(float64(v) - want)
		sumErr += e
		if e > maxErr {
			maxErr = e
		}
		if want < -R/2 && e > deepErr {
			deepErr = e
		}
	}
	mean := sumErr / float64(len(r.vals))
	t.Logf("cell %.3f mm: max error %.4f (%.2f cells), mean %.4f, deep-interior max %.4f", h, maxErr, maxErr/h, mean, deepErr)
	if maxErr > 0.75*h {
		t.Errorf("redistanced sphere is off by %.4f mm, more than 0.75 cells (%.4f)", maxErr, 0.75*h)
	}
	if mean > 0.15*h {
		t.Errorf("mean error %.4f mm is more than 0.15 cells", mean)
	}
}

// The scramble above destroys the band's own values, so this is the same
// comparison on the field as sampled: the band interpolation must do no
// worse than the scramble.
func TestRedistanceKeepsExactField(t *testing.T) {
	g := Sample(solid.Sphere(10), 48)
	r := g.Redistance()
	maxErr := 0.0
	for i := range g.vals {
		if e := math.Abs(float64(g.vals[i] - r.vals[i])); e > maxErr {
			maxErr = e
		}
	}
	t.Logf("exact vs redistanced: max |Δ| %.4f mm (%.2f cells)", maxErr, maxErr/g.CellSize())
	if maxErr > 0.6*g.CellSize() {
		t.Errorf("redistance moved an exact field by %.4f mm", maxErr)
	}
}

// A twisted column is the inexact field from the findings: sdfx measures
// distance within the slice, not to the slanted surface, so Shell comes out
// thinner than stated along the wall's normal. After Redistance the wall is
// the stated thickness. Measured along the normal at the middle of a side.
func TestRedistanceFixesTwistedShellThickness(t *testing.T) {
	column := solid.TwistExtrude(shape.Rect(v2.XY(20, 20), 0), 40, 90)
	const wall = 2.0

	raw := column.Shell(wall)
	fixed := Sample(column, 160).Redistance().Solid().Shell(wall)

	// A point on the twisted surface a third of the way up, mid-side.
	from := v3.XYZ(0, 0, -6)
	d, ok := validate.SurfaceAlong(column, from, v3.X(1), 20, 0.05)
	if !ok {
		t.Fatal("no surface found along +X")
	}
	p := from.Add(v3.X(d))
	n := validate.Gradient(column, p)
	tFixed, ok1 := validate.WallThickness(fixed, p, n, 5)
	tRaw, ok2 := validate.WallThickness(raw, p, n, 5)
	if !ok1 || !ok2 {
		t.Fatalf("wall not measurable: fixed %v raw %v", ok1, ok2)
	}
	t.Logf("wall along normal at (%.2f, %.2f, %.2f): stated %.2f, raw twist %.3f, redistanced %.3f", p.X, p.Y, p.Z, wall, tRaw, tFixed)
	validate.RequireNear(t, "redistanced wall", tFixed, wall, 0.2)
	if math.Abs(tRaw-wall) < math.Abs(tFixed-wall) {
		t.Errorf("the raw twisted shell (%.3f) measured closer to %.2f than the redistanced one (%.3f); the test no longer shows the trap", tRaw, wall, tFixed)
	}
}

// Smoothing a sphere at any radius below its own leaves it alone: a ball
// rolls over a bigger ball everywhere. This pins the four-step morphology
// down to an identity on the one shape where it must be one.
func TestSmoothLeavesSphereAlone(t *testing.T) {
	g := Sample(solid.Sphere(10), 64)
	s := g.Smooth(2).Solid()
	for _, dir := range []v3.Vec{v3.X(1), v3.Y(-1), v3.Z(1), v3.XYZ(1, 1, 1)} {
		d, ok := validate.SurfaceAlong(s, v3.Zero, dir, 15, 0.1)
		if !ok {
			t.Fatalf("no surface along %v", dir)
		}
		validate.RequireNear(t, "sphere radius after smooth", d, 10, 0.6*g.CellSize())
	}
	validate.RequireNear(t, "centre depth", validate.At(s, v3.Zero), -10, g.CellSize())
}

// Grid accounting that the scenarios rely on.
func TestGridVolumeAndCells(t *testing.T) {
	g := Sample(solid.Box(v3.XYZ(10, 20, 30), 0), 60)
	c := g.Cells()
	if c[2] != 60+2*samplePad || c[0] != 20+2*samplePad {
		t.Errorf("cells = %v, want longest axis 60 plus padding", c)
	}
	validate.RequireNear(t, "box volume", g.Volume(), 6000, 6000*0.03)
}
