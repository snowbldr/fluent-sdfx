package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const shrink = 1.0 / 0.999 // PLA ~0.1%

const pipeClearance = 1.01
const pipeDiameter = 48.45 * pipeClearance
const baseThickness = 3.0
const pipeWall = 3.0
const pipeLength = 30.0

var baseSize = v2.XY(77.0, 77.0)
var pipeOffset = v2.XY(0, 0)

const pipeRadius = 0.5 * pipeDiameter
const pipeFillet = 0.95 * pipeWall

func flange() *solid.Solid {
	base := obj.Panel2D(obj.PanelParms{
		Size:         baseSize,
		CornerRadius: 18.0,
		HoleDiameter: 3.5, // #6 screw
		HoleMargin:   [4]float64{12.0, 12.0, 12.0, 12.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	}).Extrude(2.0 * baseThickness)

	outerPipe := solid.Cylinder(2.0*pipeLength, pipeRadius+pipeWall, 0.0).
		Translate(pipeOffset.ToV3(0))
	innerPipe := solid.Cylinder(2.0*pipeLength, pipeRadius, 0.0).
		Translate(pipeOffset.ToV3(0))

	return solid.SmoothUnion(solid.PolyMin(pipeFillet), base, outerPipe).
		Cut(innerPipe).
		CutPlane(v3.XYZ(0, 0, 0), v3.XYZ(0, 0, 1))
}

func main() {
	flange().ScaleUniform(shrink).STL("flange.stl", 3.0)
}
