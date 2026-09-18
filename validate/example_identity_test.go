package validate_test

import (
	"fmt"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// The refactor guard: a rebuilt solid must have the same field as the old
// one at every sample point.
func ExampleIdentical() {
	old := solid.Box(v3.XYZ(10, 10, 10), 0).Cut(solid.Cylinder(20, 2, 0)).TranslateX(5)
	rebuilt := solid.Box(v3.XYZ(10, 10, 10), 0).TranslateX(2).
		Cut(solid.Cylinder(20, 2, 0).TranslateX(2)).TranslateX(3)
	pts := validate.GridPoints(v3.Box{Min: v3.XYZ(-2, -6, -6), Max: v3.XYZ(12, 6, 6)}, 2)
	fmt.Println(validate.Identical(old, rebuilt, pts, 1e-9))
	// Output: [ok] identical: max |Δfield| = 0.000000 at (-2.000, -6.000, -6.000) over 392 points
}

// Prove a change was local: ChangedRegion returns the box containing every
// surface point at which the two solids disagree.
func ExampleChangedRegion() {
	before := solid.Box(v3.XYZ(20, 20, 20), 0)
	after := before.Cut(solid.Cylinder(10, 2, 0).TranslateZ(10)) // a blind bore in the top face
	box, n := validate.ChangedRegion(before, after, 2, 0.05)
	fmt.Printf("changed=%v, z from %.0f to %.0f", n > 0, box.Min.Z, box.Max.Z)
	// Output: changed=true, z from 5 to 10
}

// The print-plate assertion: after ZeroZ the part's lowest point is z = 0.
func ExampleRequireStandsOn() {
	t := &testing.T{} // in real code this is the test's own *testing.T
	part := solid.Box(v3.XYZ(20, 20, 10), 0).ZeroZ()
	validate.RequireStandsOn(t, part, 0, 1e-6)
	fmt.Println("failed:", t.Failed())
	// Output: failed: false
}

// Bounds bookkeeping, for the "did that transform move what I thought"
// question.
func ExampleRequireBounds() {
	t := &testing.T{} // in real code this is the test's own *testing.T
	part := solid.Box(v3.XYZ(20, 20, 10), 0).ZeroZ()
	validate.RequireBounds(t, part, v3.Box{Min: v3.XYZ(-10, -10, 0), Max: v3.XYZ(10, 10, 10)}, 1e-6)
	fmt.Println("failed:", t.Failed())
	// Output: failed: false
}

// GridPoints fills a box with sample points for Identical and friends.
func ExampleGridPoints() {
	pts := validate.GridPoints(v3.Box{Min: v3.Vec{}, Max: v3.XYZ(2, 2, 2)}, 1)
	fmt.Printf("%d points, first %v, last %v", len(pts), pts[0], pts[len(pts)-1])
	// Output: 27 points, first {0 0 0}, last {2 2 2}
}
