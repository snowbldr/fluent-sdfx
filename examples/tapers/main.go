package main

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
)

func taper1() *solid.Solid {
	pitch := 0.50
	radius := 2.0
	length := 5.0
	taper := 20.0 * math.Pi / 180

	isoThread := shape.ISOThread(radius, pitch, true)

	s0 := isoThread.Screw(length, taper, pitch, 7)
	s1 := isoThread.Screw(length, taper, pitch, -7)

	return s0.Union(s1)
}

func taper2() *solid.Solid {
	pitch := 0.50
	radius := 2.0
	length := 10.0
	taper := 3.0 * math.Pi / 180

	isoThread := shape.ISOThread(radius, pitch, true)
	return isoThread.Screw(length, taper, pitch, 1)
}

func main() {
	taper1().STL("taper1.stl", 3.0)
	taper2().STL("taper2.stl", 3.0)
}
