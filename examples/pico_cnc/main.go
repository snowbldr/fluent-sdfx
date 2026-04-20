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

const baseHoleDiameter = 3.5

//-----------------------------------------------------------------------------
// keypad panel

func keypadPanel() *solid.Solid {
	const panelThickness = 5.5
	const panelX = 75
	const panelYa = 25
	const panelYb = 45
	const panelY = 2 * (panelYa + panelYb)

	// key hole
	const holeRadius = (22.0 + 1.5) * 0.5

	hole0 := solid.Cylinder(panelThickness, holeRadius, 0)
	hole1 := hole0.Translate(v3.Y(panelYb))
	hole2 := hole0.Translate(v3.Y(-panelYb))

	return obj.Panel3D(obj.PanelParms{
		Size:         v2.XY(panelX, panelY),
		CornerRadius: 4,
		HoleDiameter: baseHoleDiameter,
		HoleMargin:   [4]float64{7, 7, 7, 7},
		HolePattern:  [4]string{"x", "xx", "x", "xx"},
		Thickness:    panelThickness,
	}).Cut(hole0.Union(hole1, hole2))
}

//-----------------------------------------------------------------------------
// pcb mount base for rs232/ttl serial converter

func serialConverter() *solid.Solid {
	pcb := v3.XYZ(21.5, 40.0, 1.5).Add(v3.YZ(0.8, 0.4))

	wallThickness := 5.0
	innerBox := v3.XYZ(pcb.X, pcb.Y-3.0, 15)
	outerBox := innerBox.Add(v3.XYZ(wallThickness, 2.0*wallThickness, wallThickness))

	outer := solid.Box(outerBox, 0.5*wallThickness)
	inner := solid.Box(innerBox, 0)

	// body
	s := outer.Cut(inner)
	s = s.CutPlane(v3.X(0.5*innerBox.X), v3.X(-1))
	s = s.CutPlane(v3.Z(0.5*innerBox.Z), v3.Z(-1))

	// base mounting hole
	hole0 := solid.Cylinder(10*wallThickness, baseHoleDiameter*0.5, 0).
		Translate(v3.Y(0.35 * innerBox.Y))
	hole1 := hole0.MirrorXZ()
	holes := hole0.Union(hole1)

	// pcb
	board := solid.Box(pcb, 0)

	s = s.Cut(holes)
	s = s.Cut(board)

	return s
}

//-----------------------------------------------------------------------------
// pcb mount base for pico cnc

const baseThickness = 3.0
const pcbX = 92.0
const pcbY = 94.5
const pcbHoleMargin = 3.5

func picoCncStandoffs() *solid.Solid {
	const pillarHeight = 15.0
	const zOfs = 0.5 * (pillarHeight + baseThickness)
	const dx = pcbX - (2.0 * pcbHoleMargin)
	const dy = pcbY - (2.0 * pcbHoleMargin)

	return obj.Standoff3D(obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4, // #4 screw
	}).Multi(v3.VecSet{v3.XYZ(0, 0, zOfs), v3.XYZ(dx, 0, zOfs), v3.XYZ(0, dy, zOfs), v3.XYZ(dx, dy, zOfs)})
}

func picoCnc() *solid.Solid {
	const holeMargin = 3.0
	const baseX = pcbX + (2.0 * holeMargin)
	const baseY = pcbY + (2.0 * holeMargin)
	const cutoutMargin = 12.0
	const cutoutX = baseX - (2.0 * cutoutMargin)
	const cutoutY = baseY - (2.0 * cutoutMargin)

	// base
	s0 := obj.Panel2D(obj.PanelParms{
		Size:         v2.XY(baseX, baseY),
		CornerRadius: 5.0,
		HoleDiameter: baseHoleDiameter,
		HoleMargin:   [4]float64{6.0, 6.0, 6.0, 6.0},
		HolePattern:  [4]string{".x...x", ".x...x", ".x...x", ".x...x"},
	})

	// cutouts
	c0 := shape.Rect(v2.XY(cutoutX, cutoutY), 3.0)

	// extrude the base
	s2 := solid.Extrude(s0.Cut(c0), baseThickness)

	const xOfs = (0.5 * baseX) - holeMargin - pcbHoleMargin
	const yOfs = (0.5 * baseY) - holeMargin - pcbHoleMargin
	s2 = s2.Translate(v3.XY(xOfs, yOfs))

	// add the standoffs
	s3 := picoCncStandoffs()

	return solid.SmoothUnion(solid.PolyMin(3.0), s2, s3)
}

//-----------------------------------------------------------------------------
// pen holder

func penHolder() *solid.Solid {
	const holderHeight = 20.0
	const holderWidth = 25.0
	const shaftRadius = 8.0 * 0.5
	const bossDiameter = 6.0

	// spring
	k := obj.SpringParms{
		Width:         holderWidth,                       // width of spring
		Height:        holderHeight,                      // height of spring (3d only)
		WallThickness: 1,                                 // thickness of wall
		Diameter:      5,                                 // diameter of spring turn
		NumSections:   3,                                 // number of spring sections
		Boss:          [2]float64{2.0 * bossDiameter, 8}, // boss sizes
	}
	spring := obj.Spring3D(k)

	// shaft hole
	shaft := solid.Cylinder(obj.SpringLength(k), shaftRadius, 0).
		RotateY(90)

	// shaft screw boss
	bossT := obj.ThreadedCylinder(obj.ThreadedCylinderParms{
		Height:    0.5 * holderHeight,
		Diameter:  bossDiameter,
		Thread:    "unc_8_32",
		Tolerance: 0,
	}).Translate(v3.Z(30))

	return spring.Union(bossT).Cut(shaft)
}

//-----------------------------------------------------------------------------

func main() {
	picoCnc().ScaleUniform(shrink).STL("pico_cnc.stl", 3.0)

	serialConverter().ScaleUniform(shrink).STL("serial.stl", 3.0)

	keypadPanel().ScaleUniform(shrink).STL("keypad_panel.stl", 3.0)

	penHolder().ScaleUniform(shrink).STL("pen_holder.stl", 3.0)
}
