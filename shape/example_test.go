package shape_test

import (
	"github.com/snowbldr/fluent-sdfx/layout"
	"github.com/snowbldr/fluent-sdfx/shape"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

// A 40x20 mm rectangle with 2mm rounded corners.
func ExampleRect() {
	profile := shape.Rect(v2.XY(40, 20), 2)
	// profile.ToDXF("rect.dxf", 200)
	_ = profile
}

// A rectangular plate profile with a 4-up grid of mounting holes.
func ExampleShape_Cut_holes() {
	plate := shape.Rect(v2.XY(60, 40), 2)
	hole := shape.Circle(1.5)
	profile := plate.Cut(hole.Multi(layout.Grid2(20, 10, 2, 2)...))
	// profile.ToDXF("plate.dxf", 200)
	_ = profile
}

// Stack a small label rectangle on top of a base in 2D.
func ExampleShape_OnTopOf() {
	base := shape.Rect(v2.XY(40, 10), 1)
	label := shape.Rect(v2.XY(20, 6), 0)
	profile := label.OnTopOf(base.Top(), 1).Union()
	// profile.ToSVG("label.svg", 200)
	_ = profile
}
