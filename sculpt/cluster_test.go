package sculpt

import (
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Scenario: a form built from blobs, as organic parts are: a body with a
// smaller sphere grown out of its side (leaving a sharp crease), a block
// set into its top (leaving sharp edges), and a bead that was meant to sit
// on the block but was placed too high and floats. Drop the bead, then
// smooth the whole form at one radius so creases fill and edges round.
const (
	clBodyR = 10.0
	clArmR  = 4.0
	clArmX  = 10.0 // arm centre: overlaps the body by 4
	clBlock = 10.0 // cube edge
	clBlkZ  = 10.0 // cube centre: top face at 15, lower half in the body
	clBeadR = 2.0
	clBeadZ = 20.0 // bottom at 18, block top at 15: floating
	clRound = 2.0
)

func cluster() *solid.Solid {
	body := solid.Sphere(clBodyR)
	arm := solid.Sphere(clArmR).TranslateX(clArmX)
	block := solid.Box(v3.XYZ(clBlock, clBlock, clBlock), 0).TranslateZ(clBlkZ)
	bead := solid.Sphere(clBeadR).TranslateZ(clBeadZ)
	return body.Union(arm, block, bead)
}

// The crease where the arm meets the body, in the XY plane: on both spheres.
func clusterCrease() v3.Vec {
	x := (clBodyR*clBodyR - clArmR*clArmR + clArmX*clArmX) / (2 * clArmX)
	return v3.XYZ(x, math.Sqrt(clBodyR*clBodyR-x*x), 0)
}

func TestClusterFindsAndDropsTheBead(t *testing.T) {
	g := Sample(cluster(), 128)
	h := g.CellSize()

	comps := g.Components()
	if len(comps) != 2 {
		t.Fatalf("%d bodies, want 2 (form and floating bead)", len(comps))
	}
	bead := comps[1]
	t.Logf("bodies: %d cells and %d cells; bead bounds %v, deepest point %v", comps[0].Cells, bead.Cells, bead.Bounds, bead.Deep)
	wantBounds := v3.Box{Min: v3.XYZ(-clBeadR, -clBeadR, clBeadZ-clBeadR), Max: v3.XYZ(clBeadR, clBeadR, clBeadZ+clBeadR)}
	if !bead.Bounds.Equals(wantBounds, 1.5*h) {
		t.Errorf("bead bounds %v, want %v within %.2f", bead.Bounds, wantBounds, 1.5*h)
	}
	if bead.Deep.Sub(v3.Z(clBeadZ)).Length() > h {
		t.Errorf("bead's deepest point %v is not its centre", bead.Deep)
	}

	kept := g.Keep(v3.Zero)
	if n := len(kept.Components()); n != 1 {
		t.Errorf("after Keep: %d bodies, want 1", n)
	}
	ks := kept.Solid()
	// Where the bead was is now air, as far from material as the block top.
	validate.RequireNear(t, "air where the bead was", validate.At(ks, v3.Z(clBeadZ)), clBeadZ-(clBlkZ+clBlock/2), 2*h)
	// The form itself did not move.
	d, _ := validate.SurfaceAlong(ks, v3.Zero, v3.X(-1), 15, 0.1)
	validate.RequireNear(t, "body radius after Keep", d, clBodyR, 0.6*h)
	validate.RequireNear(t, "block top after Keep", validate.At(ks, v3.Z(clBlkZ+clBlock/2)), 0, 0.6*h)

	// Keeping the bead instead keeps only the bead.
	onlyBead := g.Keep(v3.Z(clBeadZ))
	if n := len(onlyBead.Components()); n != 1 {
		t.Errorf("Keep(bead): %d bodies", n)
	}
	validate.RequireNear(t, "bead volume", onlyBead.Volume(), 4.0/3.0*math.Pi*math.Pow(clBeadR, 3), 0.25*4.0/3.0*math.Pi*math.Pow(clBeadR, 3))
	if validate.Inside(onlyBead.Solid(), v3.Zero) {
		t.Error("Keep(bead) left the body")
	}
}

func TestKeepRefusesSeedInAir(t *testing.T) {
	g := Sample(cluster(), 48)
	defer func() {
		if r := recover(); r == nil {
			t.Error("Keep seeded in air did not panic")
		}
	}()
	g.Keep(v3.Z(17)) // the gap between block and bead
}

// Smooth at radius 2: the crease fills to the depth a 2 mm ball leaves,
// the block's edges round to a 2 mm fillet, flat faces and the far side of
// the body stay put, and the form stays one body.
func TestClusterSmoothFillsCreasesAndRoundsEdges(t *testing.T) {
	g := Sample(cluster(), 128).Keep(v3.Zero)
	h := g.CellSize()
	sm := g.Smooth(clRound)
	s := sm.Solid()
	before := g.Solid()

	// Concave crease between arm and body. The closing pass is a ball of
	// radius r rolling in the crease: its centre sits r outside both
	// spheres, and the crease point ends up inside the ball by its distance
	// to that centre less the ball's radius.
	c := clusterCrease()
	bx := ((clBodyR+clRound)*(clBodyR+clRound) - (clArmR+clRound)*(clArmR+clRound) + clArmX*clArmX) / (2 * clArmX)
	ballC := v3.XYZ(bx, math.Sqrt((clBodyR+clRound)*(clBodyR+clRound)-bx*bx), 0)
	wantFill := c.Sub(ballC).Length() - clRound
	validate.RequireNear(t, "crease before smooth is on the surface", validate.At(before, c), 0, 0.6*h)
	t.Logf("crease at (%.2f, %.2f, 0): a %.0f mm ball leaves it %.3f deep; got %.3f", c.X, c.Y, clRound, wantFill, -validate.At(s, c))
	validate.RequireNear(t, "crease fill depth", -validate.At(s, c), wantFill, 0.1)

	// Convex edge of the block: the top edge's midpoint ends up outside the
	// fillet by r(√2 − 1).
	edge := v3.XYZ(clBlock/2, 0, clBlkZ+clBlock/2)
	wantOut := clRound * (math.Sqrt2 - 1)
	t.Logf("block edge: expected %.3f outside the fillet, got %.3f", wantOut, validate.At(s, edge))
	validate.RequireNear(t, "block edge rounded", validate.At(s, edge), wantOut, 0.1)

	// The flat top stays flat and the far side of the body stays put.
	validate.RequireNear(t, "block top face", validate.At(s, v3.Z(clBlkZ+clBlock/2)), 0, 0.6*h)
	d, _ := validate.SurfaceAlong(s, v3.Zero, v3.X(-1), 15, 0.1)
	validate.RequireNear(t, "body radius after smooth", d, clBodyR, 0.6*h)
	tip, _ := validate.SurfaceAlong(s, v3.X(clArmX), v3.X(1), 8, 0.1)
	validate.RequireNear(t, "arm tip after smooth", tip, clArmR, 0.6*h)

	if n := len(sm.Components()); n != 1 {
		t.Errorf("smoothed form is %d bodies", n)
	}
}

// A stem thinner than twice the radius does not survive Smooth, by design;
// the test records that so nobody is surprised by it later.
func TestSmoothRemovesThinFeatures(t *testing.T) {
	stem := solid.Cylinder(12, 0.8, 0).TranslateZ(clBodyR + 5)
	form := solid.Sphere(clBodyR).Union(stem)
	g := Sample(form, 96)
	sm := g.Smooth(clRound)
	if validate.Inside(sm.Solid(), v3.Z(clBodyR+8)) {
		t.Error("a 1.6 mm stem survived a 2 mm smooth")
	}
	if !validate.Inside(sm.Solid(), v3.Zero) {
		t.Error("the body did not survive")
	}
}
