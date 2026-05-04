package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const airIntakeRadius = 0.5 * 5.125 * units.MillimetresPerInch
const airIntakeWall = (3.0 / 16.0) * units.MillimetresPerInch
const airIntakeHeight = 1.375 * units.MillimetresPerInch
const airIntakeHole = 0.5 * (5.0 / 16.0) * units.MillimetresPerInch

func airIntakeCover() *solid.Solid {
	const h0 = 2.0 * (airIntakeHeight + airIntakeWall)
	const r0 = airIntakeRadius + airIntakeWall
	body := solid.Cylinder(h0, r0, 0.75*airIntakeWall)

	const h1 = 2.0 * airIntakeHeight
	const r1 = airIntakeRadius
	cavity := solid.Cylinder(h1, r1, 0)

	const h2 = h0
	const r2 = airIntakeHole
	hole := solid.Cylinder(h2, r2, 0)

	return body.Cut(cavity.Union(hole)).CutPlane(v3.XYZ(0, 0, 0), v3.XYZ(0, 0, 1))
}

const dX = 0.5 * 5.625 * units.MillimetresPerInch
const dY0 = 0.5 * 4.25 * units.MillimetresPerInch
const dY1 = 0.5 * 5.16 * units.MillimetresPerInch
const holeRadius = 0.5 * (5.0 / 16.0) * units.MillimetresPerInch
const holeClearance = 1.05

const plateX = (2.0 * dX) + 20.0
const plateY = (2.0 * dY1) + 20.0
const plateZ = 4.0

func blockOffPlate() *solid.Solid {
	plate := shape.Rect(v2.XY(plateX, plateY), 1.0*plateZ)
	hole := shape.Circle(holeClearance * holeRadius)

	posn := []v2.Vec{
		v2.XY(dX, dY0), v2.XY(-dX, -dY0), v2.XY(dX, -dY0), v2.XY(-dX, dY0),
		v2.XY(dX, dY1), v2.XY(-dX, -dY1), v2.XY(dX, -dY1), v2.XY(-dX, dY1),
	}
	holes := hole.Multi(posn...)

	return plate.Cut(holes).Extrude(plateZ)
}

func main() {
	blockOffPlate().STL("plate.stl", 3.0)
	airIntakeCover().STL("air.stl", 3.0)
}
