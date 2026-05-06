package solid_test

import (
	"fmt"

	"github.com/snowbldr/fluent-sdfx/layout"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// A 20x10x5 mm box, sharp corners.
func ExampleBox() {
	part := solid.Box(v3.XYZ(20, 10, 5), 0)
	// part.STL("box.stl", 5.0)
	_ = part
}

// A 10mm-tall cylinder of radius 5, with a 1mm fillet at top and bottom.
func ExampleCylinder_chamfered() {
	part := solid.Cylinder(10, 5, 1)
	// part.STL("chamfered_cyl.stl", 5.0)
	_ = part
}

// Drill 4 evenly spaced holes through a disc on a 6mm bolt circle.
func ExampleSolid_Cut_polarHoles() {
	disc := solid.Cylinder(2, 10, 0)
	hole := solid.Cylinder(4, 1, 0)
	plate := disc.Cut(hole.Multi(layout.Polar(6, 4)...))
	// plate.STL("plate.stl", 10.0)
	_ = plate
}

// Stack a sphere on top of a cylinder with a 1mm gap.
func ExampleSolid_OnTopOf() {
	body := solid.Cylinder(10, 5, 0)
	cap := solid.Sphere(3)
	part := cap.OnTopOf(body.Top(), 1).Union()
	// part.STL("stack.stl", 5.0)
	_ = part
}

// Sit a part flat on the build plate (Z = 0).
func ExampleSolid_BottomAt() {
	part := solid.Box(v3.XYZ(10, 10, 10), 0).BottomAt(0)
	fmt.Printf("%.1f", part.Bottom().Point.Z)
	// Output: 0.0
}

// Smooth-blend a sphere onto a cylinder with a 2mm fillet at the seam.
func ExampleSolid_SmoothUnion() {
	body := solid.Cylinder(10, 5, 0)
	cap := solid.Sphere(4).TranslateZ(5)
	part := body.SmoothUnion(solid.PolyMin(2), cap)
	// part.STL("blend.stl", 5.0)
	_ = part
}
