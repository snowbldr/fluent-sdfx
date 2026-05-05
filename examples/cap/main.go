package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
)

const wallThickness = 2.0
const innerDiameter = 75.5
const innerHeight = 15.0

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%

func tubeCap() *solid.Solid {
	h := innerHeight + wallThickness
	r := (innerDiameter * 0.5) + wallThickness
	outer := solid.Cylinder(h, r, 1.0)

	ih := innerHeight
	ir := innerDiameter * 0.5
	inner := solid.Cylinder(ih, ir, 1.0)

	return outer.Top().On(inner.Top()).Cut()
}

func main() {
	tubeCap().ScaleUniform(shrink).STL("cap.stl", 1.2)
}
