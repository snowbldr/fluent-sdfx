package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const tachoRadius = 0.5 * 3.5 * units.MillimetresPerInch
const bracketHeight = 15.0
const bracketWidth = 10.0
const tabWidth = 20.0
const tabLength = 20.0
const slotWidth = 4.0
const screwRadius = 1.1 * 0.5 * (5.0 / 32.0) * units.MillimetresPerInch

func tachoBracket() *solid.Solid {
	// outer bracket
	const outerRadius = tachoRadius + bracketWidth
	body := shape.Circle(outerRadius)

	// inner hole
	hole := shape.Circle(tachoRadius)

	// side tabs
	tabs := shape.Rect(v2.XY(2.0*(outerRadius+tabLength), tabWidth), 0.07*(tabWidth+tabLength))

	// slot
	l := bracketWidth + tabLength
	slot := shape.Rect(v2.XY(l, slotWidth), 0).Translate(v2.XY(0.5*l+tachoRadius, 0))

	// panel hole
	const xOfs = tachoRadius + bracketWidth + 0.5*tabLength
	panelHole := shape.Circle(screwRadius).Translate(v2.XY(-xOfs, 0))

	// outer body with smooth union
	s3 := shape.SmoothUnion(solid.PolyMin(bracketWidth), body, tabs)

	// remove the holes
	s4 := s3.Cut(hole.Union(slot, panelHole))
	bracket := solid.Extrude(s4, bracketHeight)

	// clamp hole
	clampHole := solid.Cylinder(1.1*tabWidth, screwRadius, 0).
		RotateX(90).
		Translate(v3.XYZ(xOfs, 0, 0))

	return bracket.Cut(clampHole)
}

func main() {
	tachoBracket().ToSTL("tacho.stl", 300)
}
