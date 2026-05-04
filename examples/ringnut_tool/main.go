package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%

const innerDiameter = 132.0
const ringWidth = 19.0
const outerDiameter = innerDiameter + (2.0 * ringWidth)
const ringHeight = 16.0
const topGap = 90.0
const screwDiameter = 25.4 * (3.0 / 16.0)
const screwX = (topGap * 0.5) + (screwDiameter * 1.5)
const screwY = innerDiameter * 0.22

const numTabs = 36
const tabDepth = 3.5
const tabWidth = 3.5

const sideThickness = 2.5 * tabDepth
const topThickness = 2.0 * tabDepth

func outerBody() *solid.Solid {
	h := (ringHeight + topThickness) * 2.0
	r := (outerDiameter * 0.5) + sideThickness
	round := topThickness * 0.5
	return solid.Cylinder(h, r, round)
}

func innerCavity() *solid.Solid {
	h := ringHeight * 2.0
	r := outerDiameter * 0.5
	round := ringHeight * 0.1
	s0 := solid.Cylinder(h, r, round)

	h = (ringHeight + topThickness) * 2.0
	r = innerDiameter * 0.5
	s1 := solid.Cylinder(h, r, 0).
		CutPlane(v3.XYZ(topGap*0.5, 0, 0), v3.XYZ(-1, 0, 0)).
		CutPlane(v3.XYZ(-topGap*0.5, 0, 0), v3.XYZ(1, 0, 0))

	return s0.Union(s1)
}

func tab() *solid.Solid {
	size := v3.XYZ(
		tabWidth,
		ringWidth+tabDepth,
		(ringHeight+tabDepth)*2.0)

	yofs := (size.Y + innerDiameter) * 0.5
	return solid.Box(size, 0).Translate(v3.XYZ(0, yofs, 0))
}

func tabs() *solid.Solid {
	t := tab()
	thetaDeg := 360.0 / numTabs
	return t.RotateUnionZ(numTabs, solid.RotateZMatrix(thetaDeg))
}

func screwHole() *solid.Solid {
	l := ringHeight + topThickness
	r := screwDiameter * 0.5

	s := obj.CounterSunkHole3D(l, r)

	zofs := (l * 0.5) + ringHeight
	return s.Translate(v3.XYZ(0, 0, -zofs))
}

func screwHoles() *solid.Solid {
	s := screwHole()
	return s.Multi(v3.XYZ(screwX, screwY, 0),
		v3.XYZ(-screwX, screwY, 0),
		v3.XYZ(screwX, -screwY, 0),
		v3.XYZ(-screwX, -screwY, 0),
		v3.XYZ(screwX, 0, 0),
		v3.XYZ(-screwX, 0, 0))
}

func tool() *solid.Solid {
	body := outerBody()
	cavity := innerCavity()
	tabs := tabs()
	screws := screwHoles()

	s := body.Cut(cavity.Union(tabs, screws))
	return s.CutPlane(v3.XYZ(0, 0, 0), v3.XYZ(0, 0, -1))
}

func main() {
	tool().ScaleUniform(shrink).STL("tool.stl", 3.0)
}
