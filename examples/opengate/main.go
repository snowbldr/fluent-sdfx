package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%

const baseThickness = 3
const pillarHeight = 8

const pcbX = 116
const pcbY = 61

const baseX = pcbX + 30
const baseY = pcbY + 20

func standoffs() *solid.Solid {
	// standoffs with screw holes
	s := obj.Standoff3D(obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      pillarHeight,
		HoleDiameter:   2.4, // #4 screw
	})

	positions0 := []v3.Vec{v3.Zero, v3.X(pcbX), v3.XY(pcbX, pcbY), v3.Y(pcbY)}

	xOfs := -0.5 * pcbX
	yOfs := -0.5 * pcbY
	zOfs := 0.5 * (pillarHeight + baseThickness)

	return s.Multi(positions0).Translate(v3.XYZ(xOfs, yOfs, zOfs))
}

func mainBoard() *solid.Solid {
	// base
	baseShape := obj.Panel2D(obj.PanelParms{
		Size:         v2.XY(baseX, baseY),
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	})

	// cutout
	cutoutShape := obj.Panel2D(obj.PanelParms{
		Size:         v2.XY(baseX-40, baseY-40),
		CornerRadius: 5.0,
	})

	// extrude the base, add the standoffs with polymin fillet
	return solid.SmoothUnion(
		solid.PolyMin(3.0),
		solid.Extrude(baseShape.Cut(cutoutShape), baseThickness),
		standoffs(),
	)
}

func main() {
	mainBoard().ScaleUniform(shrink).STL("main_board.stl", 3.0)
}
