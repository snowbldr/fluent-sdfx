// Parametric helpers: a Panel3D with rounded corners and mounting holes.
//
// HoleMargin and HolePattern are 4-element arrays in [top, right, bottom,
// left] order. Each pattern string places a hole per character: 'x' = hole,
// '.' = skip.
package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

func main() {
	obj.Panel3D(obj.PanelParms{
		Size:         v2.XY(60, 40),
		CornerRadius: 4,
		HoleDiameter: 3,
		HoleMargin:   [4]float64{5, 5, 5, 5},
		HolePattern:  [4]string{"x.x", "x", "x.x", "x"},
		Thickness:    3,
	}).STL("out.stl", 6.0)
}
