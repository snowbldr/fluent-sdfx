package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func arrow1() *solid.Solid {
	k := obj.ArrowParms{
		Axis:  [2]float64{50, 1},
		Head:  [2]float64{5, 2},
		Tail:  [2]float64{5, 2},
		Style: "cb",
	}
	return obj.Arrow3D(k)
}

func axes(min, max v3.Vec) *solid.Solid {
	return obj.Axes3D(min, max)
}

func main() {
	arrow1().ToSTL("arrow1.stl", 300)
	axes(v3.XYZ(-10, -10, -10), v3.XYZ(10, 20, 20)).ToSTL("axes1.stl", 300)
	axes(v3.XYZ(-10, -20, -30), v3.XYZ(0, 0, 0)).ToSTL("axes2.stl", 300)
	axes(v3.XYZ(0, 0, 0), v3.XYZ(500, 500, 1000)).ToSTL("axes3.stl", 300)
}
