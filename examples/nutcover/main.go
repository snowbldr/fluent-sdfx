package main

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const nutFlat2Flat = 19.0        // measured flat 2 flat nut size
const recessHeight = 20.0        // recess within cover
const wallThickness = 2.0        // wall thickness
const counterBoreDiameter = 23.0 // diameter of washer counterbore
const counterBoreDepth = 2.0     // depth of washer counterbore
const nutFit = 1.01              // press fit on nut

func hexRadius(f2f float64) float64 {
	return f2f / (2.0 * math.Cos(30*math.Pi/180))
}

func cover() *solid.Solid {
	r := (hexRadius(nutFlat2Flat) * nutFit) + wallThickness
	h := recessHeight + wallThickness
	return solid.Cylinder(2*h, r, 0.1*r)
}

func recess() *solid.Solid {
	r := hexRadius(nutFlat2Flat) * nutFit
	h := recessHeight
	return obj.HexHead3D(r, 2*h, "")
}

func counterbore() *solid.Solid {
	r := counterBoreDiameter * 0.5
	h := counterBoreDepth
	return solid.Cylinder(2*h, r, 0)
}

func nutcover() *solid.Solid {
	rec := recess()
	result := cover().Cut(rec.Union(counterbore()))
	return result.CutPlane(v3.Zero, v3.Z(1))
}

func main() {
	nutcover().ToSTL("cover.stl", 150)
}
