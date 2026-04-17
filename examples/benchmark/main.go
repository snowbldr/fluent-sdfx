package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	shape.Circle(5).Benchmark("circle SDF2")
	shape.FlatFlankCam(30, 20, 5).Benchmark("cam1 SDF2")
	shape.ThreeArcCam(30, 20, 5, 200).Benchmark("cam2 SDF2")
	shape.Polygon(shape.Nagon(6, 10.0)).Benchmark("poly6 SDF2")
	shape.Polygon(shape.Nagon(12, 10.0)).Benchmark("poly12 SDF2")
	solid.Box(v3.XYZ(10, 20, 30), 1).Benchmark("box SDF3")
}
