package layout_test

import (
	"github.com/snowbldr/fluent-sdfx/layout"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// 6 holes evenly spaced on a 20mm bolt circle, drilled through a disc.
func ExamplePolar() {
	disc := solid.Cylinder(2, 30, 0)
	hole := solid.Cylinder(4, 2, 0)
	plate := disc.Cut(hole.Multi(layout.Polar(20, 6)...))
	// plate.STL("flange.stl", 5.0)
	_ = plate
}

// 4x4 grid of pegs spaced 10mm apart.
func ExampleGrid() {
	peg := solid.Cylinder(5, 1, 0)
	pegs := peg.Multi(layout.Grid(10, 10, 4, 4)...)
	// pegs.STL("pegs.stl", 5.0)
	_ = pegs
}

// PCB standoffs at the 4 corners of an 80x50mm panel.
func ExampleRectCorners() {
	standoff := solid.Cylinder(8, 3, 0)
	mounts := standoff.Multi(layout.RectCorners(80, 50)...)
	// mounts.STL("mounts.stl", 5.0)
	_ = mounts
}

// A small sphere placed at every corner of a 40x40x40 box.
func ExampleBoxCorners() {
	marker := solid.Sphere(2)
	markers := marker.Multi(layout.BoxCorners(v3.XYZ(40, 40, 40))...)
	// markers.STL("markers.stl", 5.0)
	_ = markers
}
