package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func ferriteMount() *solid.Solid {
	const rodRadius = 10.4 * 0.5
	const baseSize = 20.0
	const rodHeight = 25.0
	const WallThickness = 3.0
	const holderDepth = 6.0
	const holderRadius = WallThickness + rodRadius
	const holderLength = holderDepth + WallThickness

	const baseX = baseSize + 2.0*holderRadius
	const baseY = holderRadius
	const baseZ = 20.0

	// support wall (triangle offset)
	wall := obj.IsocelesTriangle2D(baseSize, rodHeight).Offset(holderRadius).Extrude(WallThickness)
	// base
	base := solid.Box(v3.XYZ(baseX, baseY, baseZ), 0)
	// holder
	holder := solid.Cylinder(holderLength, holderRadius, 0)
	rodHole := solid.Cylinder(holderDepth, rodRadius, 0)

	return base.Translate(v3.YZ(-0.5*(baseY+rodHeight), 0.5*(baseZ-WallThickness))).Union(
		wall,
		holder.Cut(rodHole.Translate(v3.Z(0.5*(holderLength-holderDepth)))).
			Translate(v3.YZ(rodHeight*0.5, 0.5*(holderLength-WallThickness))),
	).CutPlane(v3.Y(-0.5*rodHeight-WallThickness), v3.Y(1))
}

const screwHoleRadius = 3.7 * 0.5
const shaftRadius = 8.0 * 0.5

func vcapMountHole(length float64) *solid.Solid {
	const screwOffset = 14.0 * 0.5
	sh := obj.ChamferedHole3D(length, screwHoleRadius, 0.5)
	h0 := sh.Translate(v3.X(screwOffset))
	h1 := sh.Translate(v3.X(-screwOffset))
	h2 := solid.Cylinder(length, shaftRadius+0.4, 0)
	return h0.Union(h1, h2)
}

func vcapShaftHole(length float64) *solid.Solid {
	const tipRadius = 6.7 * 0.5
	const tipFlatToFlat = 4.6
	const tipLength = 2.5

	tip := solid.Cylinder(tipLength, tipRadius, 0)
	xOfs := tipFlatToFlat * 0.5
	tip = tip.CutPlane(v3.X(xOfs), v3.X(-1))
	tip = tip.CutPlane(v3.X(-xOfs), v3.X(1))
	tip = tip.Translate(v3.Z(0.5 * (length - tipLength)))

	// countersink
	const csRadius = 8.0 * 0.5
	const csLength = 3.0
	cs := solid.Cylinder(csLength, csRadius, 0).
		Translate(v3.Z(-0.5 * (length - csLength)))

	// screw hole
	hole := solid.Cylinder(length, screwHoleRadius, 0)

	return hole.Union(tip, cs)
}

const mountThickness = 5.0

func vcapKnob() *solid.Solid {
	const knobRadius = 40.0 * 0.5
	const knobWidth = 15.0
	const shaftLength = mountThickness - 1.3
	const totalLength = knobWidth + shaftLength

	knob := solid.Cylinder(knobWidth, knobRadius, 2.0)
	knurl := obj.KnurledHead3D(knobRadius, knobWidth*0.67, 3.0)
	shaft := solid.Cylinder(totalLength, shaftRadius, 0)
	hole := vcapShaftHole(totalLength)

	return knob.Union(knurl, shaft.Translate(v3.Z(0.5*shaftLength))).
		Cut(hole.Translate(v3.Z(0.5 * shaftLength)))
}

func vcapMount() *solid.Solid {
	const length = 45.0

	mount := solid.Box(v3.XYZ(length, length, mountThickness), 0)
	holes := vcapMountHole(mountThickness)
	return mount.Cut(holes)
}

func main() {
	knob := vcapKnob()
	knob.STL("vc_knob.stl", 5.0)

	mount := vcapMount()
	mount.STL("vc_mount.stl", 5.0)

	fm := ferriteMount()
	fm.STL("fr_mount.stl", 5.0)
}
