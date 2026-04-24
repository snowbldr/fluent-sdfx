package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%

func holder() *solid.Solid {
	// dimensions taken from loadcell
	const xLoadcell = 34.0
	const yLoadCell = 34.0
	const zLoadCell = 3.0
	const rLoadCell = 8.0
	const innerMargin = 4.0

	// dimensions to outside body
	const outerMargin = 4.0
	const bodyHeight = 2.0 * 8.0

	// body
	bodySize := v2.XY(xLoadcell+2.0*outerMargin, yLoadCell+2.0*outerMargin)
	bodyRadius := rLoadCell + outerMargin
	body2d := shape.Rect(bodySize, bodyRadius)
	body3d := body2d.ExtrudeRounded(bodyHeight, 2.0)

	// tabs
	tabX := 15.0
	tabSize := v2.XY(bodySize.X+2.0*tabX, 0.5*bodySize.Y)
	tabHeight := bodyHeight * 0.75
	tab2d := shape.Rect(tabSize, bodyRadius*0.25)
	tab3d := tab2d.ExtrudeRounded(tabHeight, 2.0)

	// screw holes
	screw0 := obj.CounterSunkHole3D(tabHeight, 2.0)
	screwOfs := 0.5*(bodySize.X+tabX) + 1.0
	screwL := screw0.Translate(v3.X(-screwOfs))
	screwR := screw0.Translate(v3.X(screwOfs))
	screw3d := screwL.Union(screwR)

	// inner hole
	holeSize := v2.XY(xLoadcell-2.0*innerMargin, yLoadCell-2.0*innerMargin)
	holeRadius := rLoadCell - innerMargin
	hole3d := shape.Rect(holeSize, holeRadius).Extrude(bodyHeight)

	// recess
	recessSize := v2.XY(xLoadcell, yLoadCell)
	recess3d := shape.Rect(recessSize, rLoadCell).Extrude(zLoadCell)
	zOfs := 0.5 * (bodyHeight - zLoadCell)
	recess3d = recess3d.Translate(v3.Z(zOfs))

	// wire recess
	wireSize := v3.XYZ(2.0, 2.0, 3.0*outerMargin)
	wire3d := solid.Box(wireSize, 0).RotateX(90)
	zOfs = 0.5 * (bodyHeight - wireSize.X)
	yOfs := 0.5 * (yLoadCell + outerMargin)
	wire3d = wire3d.Translate(v3.YZ(yOfs, zOfs))

	// union body + tabs with fillet
	h := solid.SmoothUnion(solid.PolyMin(2.0), body3d, tab3d)
	// remove the holes
	h = h.Cut(hole3d, recess3d, screw3d, wire3d)
	// cut along xy plane
	h = h.CutPlane(v3.Zero, v3.Z(1))

	return h
}

func main() {
	holder().ScaleUniform(shrink).STL("holder.stl", 3.0)
}
