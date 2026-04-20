package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func arrow1() *solid.Solid {
	return obj.Arrow3D(obj.ArrowParms{
		Axis:  [2]float64{50, 1},
		Head:  [2]float64{5, 2},
		Tail:  [2]float64{5, 2},
		Style: "cb",
	})
}

func axes(min, max v3.Vec) *solid.Solid {
	return obj.Axes3D(min, max)
}

func main() {
	arrow1().STL("arrow1.stl", 3.0)
	axes(v3.XYZ(-10, -10, -10), v3.XYZ(10, 20, 20)).STL("axes1.stl", 3.0)
	axes(v3.XYZ(-10, -20, -30), v3.XYZ(0, 0, 0)).STL("axes2.stl", 3.0)
	axes(v3.XYZ(0, 0, 0), v3.XYZ(500, 500, 1000)).STL("axes3.stl", 3.0)
}
