package sculpt

import (
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Scenario: a hollow bauble with a hanging loop. The ornament was modelled
// as a shell to save material, and now it is going into a plaster mould,
// where the pattern has to be solid. Bucket-fill the void from its centre.
//
// Outer radius 15.75 (a 1.5 mm wall straddling r=15), loop a torus standing
// upright on top with its tube rooted in the wall.
const (
	baubleOuter = 15.75
	loopMajor   = 3.0
	loopMinor   = 0.8
	loopZ       = 18.3 // tube bottom at 14.5: inside the wall, not through it
)

func bauble() *solid.Solid {
	ball := solid.Sphere(15).Shell(1.5)
	loop := solid.Torus(loopMajor, loopMinor).RotateX(90).TranslateZ(loopZ)
	return ball.Union(loop)
}

func TestOrnamentFillVoid(t *testing.T) {
	g := Sample(bauble(), 96)
	h := g.CellSize()
	t.Logf("cell %.3f mm, %v cells", h, g.Cells())

	if !g.Enclosed(v3.Zero) {
		t.Fatal("the bauble's void reads as open; the fill test cannot proceed")
	}
	if n := len(g.Components()); n != 1 {
		t.Errorf("hollow bauble with loop: %d bodies, want 1 (loop must be rooted in the wall)", n)
	}

	filled := g.Fill(v3.Zero)
	fs := filled.Solid()

	// The centre is now deep inside: as far from the surface as the outer radius.
	validate.RequireNear(t, "depth at centre", validate.At(fs, v3.Zero), -baubleOuter, 2*h)

	// The outer surface has not moved, in several directions.
	for _, dir := range []v3.Vec{v3.X(1), v3.Y(-1), v3.Z(-1), v3.XYZ(1, -1, 0.5)} {
		d, ok := validate.SurfaceAlong(fs, v3.Zero, dir, 20, 0.1)
		if !ok {
			t.Fatalf("no surface along %v after fill", dir)
		}
		validate.RequireNear(t, "outer radius after fill", d, baubleOuter, 0.6*h)
	}

	// The loop is still there and its hole is still open.
	tubeTop := v3.Z(loopZ + loopMajor)
	if v := validate.At(fs, tubeTop); v > -0.5*loopMinor {
		t.Errorf("loop tube centre reads %.3f, want deep inside (≤ %.3f)", v, -0.5*loopMinor)
	}
	if v := validate.At(fs, v3.Z(loopZ)); v < 1.0 {
		t.Errorf("loop hole centre reads %.3f, want clearly in air (hole radius %.1f)", v, loopMajor-loopMinor)
	}

	// Volume: a solid ball plus the part of the loop outside the wall.
	ball := 4.0 / 3.0 * math.Pi * math.Pow(baubleOuter, 3)
	loop := 2 * math.Pi * math.Pi * loopMajor * loopMinor * loopMinor
	t.Logf("volume: hollow %.0f, filled %.0f, solid ball %.0f (+ loop ≤ %.0f)", g.Volume(), filled.Volume(), ball, loop)
	if v := filled.Volume(); v < ball*0.98 || v > (ball+loop)*1.02 {
		t.Errorf("filled volume %.0f outside [%.0f, %.0f]", v, ball*0.98, (ball+loop)*1.02)
	}
	if n := len(filled.Components()); n != 1 {
		t.Errorf("filled bauble: %d bodies, want 1", n)
	}

	// The air outside was not touched: compare the two fields where both are
	// true distances, a cell or more from the surface.
	pts := []v3.Vec{v3.X(18), v3.Y(-20), v3.Z(-19), v3.XYZ(12, 12, 0), v3.XYZ(0, 0, 23.5)}
	validate.RequireIdentical(t, g.Solid(), fs, pts, h)
}

// The same bauble with a vent hole in the bottom is not enclosed, and the
// fill refuses rather than flooding the room.
func TestOrnamentFillRefusesOpenVoid(t *testing.T) {
	vent := solid.Cylinder(10, 1.5, 0).TranslateZ(-15)
	g := Sample(bauble().Cut(vent), 96)
	if g.Enclosed(v3.Zero) {
		t.Fatal("vented bauble reads as enclosed")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Error("Fill on an open void did not panic")
		} else {
			t.Logf("Fill refused: %v", r)
		}
	}()
	g.Fill(v3.Zero)
}

// Seeding a fill from inside the wall is a mistake and says so.
func TestFillRefusesSeedInMaterial(t *testing.T) {
	g := Sample(bauble(), 48)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Fill seeded in material did not panic")
		}
	}()
	g.Fill(v3.X(15))
}
