package validate_test

import (
	"fmt"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Where is the surface, exactly? March a ray and bisect the crossing — this
// measures the real geometry, unlike a raw field reading after a boolean.
func ExampleSurfaceAlong() {
	cube := solid.Box(v3.XYZ(10, 10, 10), 0)
	d, ok := validate.SurfaceAlong(cube, v3.Vec{}, v3.X(1), 20, 0)
	fmt.Printf("half width = %.3f mm (found=%v)", d, ok)
	// Output: half width = 5.000 mm (found=true)
}

// The outer and bore radii of a turned part, read back out of the solid at a
// chosen angle and height.
func ExampleSurfaceR() {
	tube := solid.Cylinder(20, 5, 0).Cut(solid.Cylinder(30, 3, 0))
	bore, _ := validate.SurfaceR(tube, 0, 0, 0, 20)
	outer, _ := validate.SurfaceR(tube, 0, 0, 4, 20)
	fmt.Printf("bore r = %.3f, outer r = %.3f", bore, outer)
	// Output: bore r = 3.000, outer r = 5.000
}

// The top surface at a plan-view point: march down from above the part.
func ExampleSurfaceZ() {
	part := solid.Box(v3.XYZ(20, 20, 10), 0).ZeroZ().
		Union(solid.Cylinder(6, 2, 0).TranslateZ(13))
	onBoss, _ := validate.SurfaceZ(part, 0, 0, 50, -1)
	onPlate, _ := validate.SurfaceZ(part, 8, 8, 50, -1)
	fmt.Printf("boss top = %.3f, plate top = %.3f", onBoss, onPlate)
	// Output: boss top = 16.000, plate top = 10.000
}

// Opening measures a void through a point in air — here the gap between two
// bodies; WallThickness measures the material run through a point inside.
func ExampleOpening() {
	pair := solid.Box(v3.XYZ(10, 10, 10), 0).TranslateX(-6.5).
		Union(solid.Box(v3.XYZ(10, 10, 10), 0).TranslateX(6.5))
	gap, ok := validate.Opening(pair, v3.Vec{}, v3.X(1), 20)
	fmt.Printf("gap = %.3f mm (found=%v)", gap, ok)
	// Output: gap = 3.000 mm (found=true)
}

// The wall of a tube, measured radially through a point in the wall.
func ExampleWallThickness() {
	tube := solid.Cylinder(20, 5, 0).Cut(solid.Cylinder(30, 3, 0))
	w, ok := validate.WallThickness(tube, v3.XYZ(4, 0, 0), v3.X(1), 20)
	fmt.Printf("wall = %.3f mm (found=%v)", w, ok)
	// Output: wall = 2.000 mm (found=true)
}

// A screening sweep for accidental thin skins over a whole region.
func ExampleMinWall() {
	tube := solid.Cylinder(20, 5, 0).Cut(solid.Cylinder(30, 3, 0))
	r := validate.MinWall(tube, v3.Box{Min: v3.XYZ(-5, -5, 0), Max: v3.XYZ(5, 5, 0)}, 1, 20)
	fmt.Printf("thinnest run = %.3f mm over %d samples", r.Min, r.Samples)
	// Output: thinnest run = 2.000 mm over 40 samples
}

// Measure, then assert, in the standard "name = got (want w ± tol)" form.
func ExampleRequireNear() {
	t := &testing.T{} // in real code this is the test's own *testing.T
	tube := solid.Cylinder(20, 5, 0).Cut(solid.Cylinder(30, 3, 0))
	bore, _ := validate.SurfaceR(tube, 0, 0, 0, 20)
	validate.RequireNear(t, "bore radius", bore, 3.0, 0.01)
	fmt.Println("failed:", t.Failed())
	// Output: failed: false
}
