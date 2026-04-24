package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%

const panelThickness = 2.5 // mm

func standoff(h float64) *solid.Solid {
	return obj.Standoff3D(obj.StandoffParms{
		PillarHeight:   h,
		PillarDiameter: 8,
		HoleDepth:      10,
		HoleDiameter:   2.4, // #4 screw
	})
}

// halfBreadBoardStandoffs returns the standoffs for an adafruit 1/2 breadboard.
func halfBreadBoardStandoffs(h float64) *solid.Solid {
	s := standoff(h)
	positions := []v3.Vec{v3.XYZ(0, -1450*units.Mil, 0), v3.XYZ(0, 1450*units.Mil, 0)}
	return s.Multi(positions)
}

// pot0 returns the panel hole/indent for a potentiometer
func pot0() *solid.Solid {
	return obj.PanelHole3D(obj.PanelHoleParms{
		Diameter:  9.4,
		Thickness: panelThickness,
		Indent:    v3.XYZ(2, 4, 2),
		Offset:    11.0,
	})
}

// pot1 returns the panel hole/indent for a potentiometer
func pot1() *solid.Solid {
	return obj.PanelHole3D(obj.PanelHoleParms{
		Diameter:  7.2,
		Thickness: panelThickness,
		Indent:    v3.XYZ(2, 2, 1.5),
		Offset:    7.0,
	})
}

// spdt returns the panel hole/indent for a spdt switch
func spdt() *solid.Solid {
	return obj.PanelHole3D(obj.PanelHoleParms{
		Diameter:  6.2,
		Thickness: panelThickness,
		Indent:    v3.XYZ(2, 2, 1.5),
		Offset:    5.4,
	})
}

// led returns the panel hole for an led bezel
func led() *solid.Solid {
	return obj.PanelHole3D(obj.PanelHoleParms{
		Diameter:  7.0,
		Thickness: panelThickness,
	})
}

// jack35 returns the panel hole/indent for a 3.5 mm audio jack
func jack35() *solid.Solid {
	return obj.PanelHole3D(obj.PanelHoleParms{
		Diameter:  6.4,
		Thickness: panelThickness,
		Indent:    v3.XYZ(2, 2, 1.5),
		Offset:    4.9,
	})
}

// powerBoardMount returns a pcb mount for a SynthRotek Noise Filtering Power Distribution Board.
func powerBoardMount() *solid.Solid {
	const baseThickness = 3
	const standoffHeight = 10
	const xSpace = 0.9 * units.MillimetresPerInch
	const ySpace = 1.1 * units.MillimetresPerInch

	s0 := standoff(standoffHeight)
	// 4x2 sections
	const zOfs = (baseThickness + standoffHeight) * 0.5
	positions := []v3.Vec{
		v3.XYZ(-1.5*xSpace, -0.5*ySpace, zOfs),
		v3.XYZ(-1.5*xSpace, 0.5*ySpace, zOfs),
		v3.XYZ(-0.5*xSpace, -0.5*ySpace, zOfs),
		v3.XYZ(-0.5*xSpace, 0.5*ySpace, zOfs),
		v3.XYZ(0.5*xSpace, -0.5*ySpace, zOfs),
		v3.XYZ(0.5*xSpace, 0.5*ySpace, zOfs),
		v3.XYZ(1.5*xSpace, -0.5*ySpace, zOfs),
		v3.XYZ(1.5*xSpace, 0.5*ySpace, zOfs),
	}
	s1 := s0.Multi(positions)

	// base
	const baseX = (4 - 0.1) * xSpace
	const baseY = 2.0 * ySpace
	s2 := obj.Panel2D(obj.PanelParms{
		Size:         v2.XY(baseX, baseY),
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	})

	// cutout
	c0 := shape.Rect(v2.XY(3*xSpace, 0.5*ySpace), 3.0)
	s3 := s2.Cut(c0).Extrude(baseThickness)

	return solid.SmoothUnion(solid.PolyMin(3.0), s3, s1)
}

var psuSize = v3.XYZ(98, 129, 38)

// rt65b returns a model of a Meanwell RT-65B PSU
func rt65b() *solid.Solid {
	body := solid.Box(psuSize, 0).Translate(psuSize.MulScalar(0.5))

	s0 := obj.CounterBoredHole3D(12, 3.8*0.5, 10.6*0.5, 3.5)

	// vertical screws
	vs0 := s0.RotateY(180).Translate(v3.XYZ(31, 4.5+73.5, 0))
	vs1 := vs0.Translate(v3.XYZ(33, 0, 0))

	// horizontal screws
	hs0 := s0.RotateY(90).Translate(v3.XYZ(psuSize.X, 32, 38-18.5))
	hs1 := hs0.Translate(v3.XYZ(0, 77, 9))
	hs2 := hs0.Translate(v3.XYZ(0, 77, -9))

	return body.Union(vs0, vs1, hs0, hs1, hs2)
}

// psuMount returns a mount for a Meanwell RT-65B PSU
func psuMount() *solid.Solid {
	// base
	const baseThickness = 6
	baseSize := v2.XY(135, 145)
	base := obj.Panel3D(obj.PanelParms{
		Size:         baseSize,
		CornerRadius: 5.0,
		HoleDiameter: 4.0,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
		Thickness:    baseThickness,
	})

	// cutout 0
	c0 := obj.Panel3D(obj.PanelParms{
		Size:         v2.XY(90, 55),
		CornerRadius: 4.0,
		Thickness:    baseThickness,
	}).Translate(v3.XYZ(-10, -27.5, 0))

	// cutout 1
	c1 := obj.Panel3D(obj.PanelParms{
		Size:         v2.XY(90, 30),
		CornerRadius: 4.0,
		Thickness:    baseThickness,
	}).Translate(v3.XYZ(-10, 40, 0))

	// upright mount
	uprightSize := v2.XY(psuSize.Z+baseThickness*0.5, baseSize.Y)
	upright := obj.Panel3D(obj.PanelParms{
		Size:         uprightSize,
		CornerRadius: 3.0,
		Thickness:    baseThickness,
	}).RotateY(90).Translate(v3.XYZ(psuSize.X+baseThickness, 0, uprightSize.X).MulScalar(0.5))

	psu := rt65b().Translate(v3.XYZ(-psuSize.X, -psuSize.Y, baseThickness).MulScalar(0.5))

	return base.Union(upright).Cut(psu, c0, c1)
}

// powerPanel returns a mounting panel for a ac-14-f16a power connector.
func powerPanel() *solid.Solid {
	const baseThickness = 4

	s := obj.Panel2D(obj.PanelParms{
		Size:         v2.XY(85, 95),
		CornerRadius: 5.0,
		HoleDiameter: 4.0,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	})

	// panel cutout
	c0 := shape.Rect(v2.XY(28, 48), 3)

	// mounting holes
	hole := shape.Circle(0.5 * 4.5)
	c1 := hole.Translate(v2.XY(20, 0))
	c2 := hole.Translate(v2.XY(-20, 0))

	cutouts := c0.Union(c1, c2)

	return s.Cut(cutouts).Extrude(baseThickness)
}

// powerPanelRouting returns a routing pattern for the power panel.
func powerPanelRouting() *solid.Solid {
	const baseThickness = 4

	s := obj.Panel2D(obj.PanelParms{
		Size:         v2.XY(85, 95),
		CornerRadius: 5.0,
		HoleDiameter: 4.0,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	})

	c := shape.Rect(v2.XY(55, 65), 3)

	return s.Cut(c).Extrude(baseThickness)
}

// arPanel returns the panel for an attack/release module.
func arPanel() *solid.Solid {
	// 3u x 12hp panel
	panel := obj.EuroRackPanel3D(obj.EuroRackParms{
		U:            3,
		HP:           12,
		CornerRadius: 3,
		HoleDiameter: 3.6,
		Thickness:    panelThickness,
		Ridge:        true,
	})

	// breadboard standoffs
	const standoffHeight = 25
	so := halfBreadBoardStandoffs(standoffHeight).
		Translate(v3.XYZ(0, 3, (panelThickness+standoffHeight)*0.5))
	s := solid.SmoothUnion(solid.PolyMin(2), panel, so)

	// push button
	pb := solid.Box(v3.XYZ(13.2, 10.8, panelThickness), 0)

	// cv input/output
	cv := jack35()
	cv0 := cv.Translate(v3.XYZ(-20, -45, 0))
	cv1 := cv.Translate(v3.XYZ(20, -45, 0))

	// LED
	ledS := led().Translate(v3.XYZ(0, -45, 0))

	// attack/release pots
	pot := pot0()
	pot0s := pot.Translate(v3.XYZ(-15, 25, 0))
	pot1s := pot.Translate(v3.XYZ(15, 25, 0))

	// spdt switch
	spdtS := spdt().Translate(v3.XYZ(0, -22, 0))

	return s.Cut(pb, cv0, cv1, ledS, pot0s, pot1s, spdtS)
}

// bbPanel returns a panel for mounting a half bread board.
func bbPanel() *solid.Solid {
	panel := obj.EuroRackPanel3D(obj.EuroRackParms{
		U:            3,
		HP:           12,
		CornerRadius: 3,
		HoleDiameter: 3.6,
		Thickness:    panelThickness,
		Ridge:        true,
	})

	const standoffHeight = 12
	so := halfBreadBoardStandoffs(standoffHeight).
		Translate(v3.XYZ(0, 3, (panelThickness+standoffHeight)*0.5))
	return solid.SmoothUnion(solid.PolyMin(2), panel, so)
}

func main() {
	arPanel().ScaleUniform(shrink).STL("ar_panel.stl", 3.0)
	powerBoardMount().ScaleUniform(shrink).STL("pwr_mount.stl", 3.0)
	powerPanel().ScaleUniform(shrink).STL("pwr_panel.stl", 3.0)
	powerPanelRouting().ScaleUniform(shrink).STL("pwr_panel_routing.stl", 3.0)
	psuMount().ScaleUniform(shrink).STL("psu_mount.stl", 3.0)
	bbPanel().ScaleUniform(shrink).STL("bb_panel.stl", 3.0)
}
