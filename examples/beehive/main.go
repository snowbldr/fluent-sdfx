package main

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%

const wheelRadius = 6.5 * 0.5 * units.MillimetresPerInch
const wheelThickness = 0.25 * units.MillimetresPerInch

// wheelRetainer returns a retaining clip for the entrance wheel.
func wheelRetainer() *solid.Solid {
	size := v3.XYZ(
		1.75*units.MillimetresPerInch,
		1.5*units.MillimetresPerInch,
		1.5*wheelThickness)

	const round = 0.25 * units.MillimetresPerInch
	const holeRadius = 7 * 0.5
	const clearance = 1

	s2d := shape.Rect(v2.XY(size.X, size.Y), round)

	hole := shape.Circle(holeRadius).Translate(v2.XY(0, 0.25*size.Y))
	s2d = s2d.Cut(hole)

	s3d := solid.Extrude(s2d, size.Z).Translate(v3.XYZ(0, wheelRadius, 0))

	t := wheelThickness * 0.9
	ofs := 0.5 * (t - size.Z)
	wheel := solid.Cylinder(t, wheelRadius+clearance, 0).Translate(v3.XYZ(0, 0, ofs))

	return s3d.Cut(wheel)
}

// entrance0 returns an open entrance
func entrance0(size v3.Vec) *solid.Solid {
	r := size.Y * 0.5
	s0 := shape.Line(size.X-(2*r), r)
	return solid.Extrude(s0, size.Z)
}

// entrance1 returns a vent entrance
func entrance1(size v3.Vec) *solid.Solid {
	const rows = 3
	const cols = 16
	const holeRadius = 3.2 * 0.5

	hole := shape.Circle(holeRadius)

	size.X -= 2 * holeRadius
	size.Y -= 2 * holeRadius
	dx := size.X / (cols - 1)
	dy := size.Y / (rows - 1)
	xOfs := -size.X / 2
	yOfs := size.Y / 2

	positions := []v2.Vec{}
	x := xOfs
	for i := 0; i < cols; i++ {
		y := yOfs
		for j := 0; j < rows; j++ {
			positions = append(positions, v2.XY(x, y))
			y -= dy
		}
		x += dx
	}
	s := hole.Multi(positions)

	return solid.Extrude(s, size.Z)
}

// entranceWheel returns a rotating entrance for a swarm trap.
func entranceWheel() *solid.Solid {
	plate := solid.Cylinder(wheelThickness, wheelRadius, 0)

	hole := solid.Cylinder(wheelThickness, 2.5, 0)

	entranceSize := v3.XYZ(
		4*units.MillimetresPerInch,
		0.5*units.MillimetresPerInch,
		wheelThickness)

	const k = 1.6
	ofs := k * entranceSize.X * 0.5 * math.Tan(units.DtoR(30))

	// open entrance
	e0 := entrance0(entranceSize).Translate(v3.XYZ(0, ofs, 0))

	// vent entrance
	e1 := entrance1(entranceSize).Translate(v3.XYZ(0, ofs, 0)).RotateZ(120)

	return plate.Cut(e0, e1, hole)
}

func holePattern(n int) string {
	s := make([]byte, n)
	for i := range s {
		s[i] = byte('x')
	}
	return string(s)
}

func entranceReducer() *solid.Solid {
	const zSize = 0.25 * units.MillimetresPerInch
	const xSize = 6.0 * units.MillimetresPerInch
	const ySize = 1.9 * units.MillimetresPerInch

	k := obj.PanelParms{
		Size:         v2.XY(xSize, ySize),
		CornerRadius: 5.0,
	}
	s := obj.Panel2D(k)

	const holeRadius = (3.0 / 16.0) * units.MillimetresPerInch
	hole := shape.Line(2*holeRadius, holeRadius).Rotate(90)

	const entranceSize = 4.0 * units.MillimetresPerInch
	const n = 6
	const gap = (entranceSize - (n * holeRadius)) / (n + 1)
	const yOfs = -ySize * 0.5
	const xOfs = (n - 1) * (holeRadius + gap) * 0.5
	p0 := v2.XY(-xOfs, yOfs)
	p1 := v2.XY(xOfs+holeRadius+gap, yOfs)
	hole = hole.LineOf(p0, p1, holePattern(n))

	return solid.Extrude(s.Cut(hole), zSize)
}

func angleHole() *solid.Solid {
	const l = 1.25 * units.MillimetresPerInch
	const t = 0.125 * units.MillimetresPerInch
	const r = 0.125 * units.MillimetresPerInch

	k := obj.AngleParams{
		X:          obj.AngleLeg{Length: l, Thickness: t},
		Y:          obj.AngleLeg{Length: l, Thickness: t},
		RootRadius: r,
		Length:     12 * units.MillimetresPerInch,
	}

	return obj.Angle3D(k).Translate(v3.XYZ(-0.5*l, -0.5*l, 0))
}

func antCap() *solid.Solid {
	// angle hole
	angle3d := angleHole()

	// outer cap
	capHeight := 1.75 * units.MillimetresPerInch
	capRadius1 := 0.5 * 2.5 * units.MillimetresPerInch
	capRadius0 := 0.5 * 4.0 * units.MillimetresPerInch
	hat0 := solid.Cone(capHeight, capRadius0, capRadius1, 0)

	// inner cap
	const capWall = 0.25 * units.MillimetresPerInch
	capHeight -= capWall
	capRadius1 -= capWall
	capRadius0 -= capWall
	hat1 := solid.Cone(capHeight, capRadius0, capRadius1, 0).Translate(v3.XYZ(0, 0, -0.5*capWall))

	return hat0.Cut(angle3d, hat1)
}

func main() {
	entranceReducer().ScaleUniform(shrink).ToSTL("reducer.stl", 300)

	entranceWheel().ScaleUniform(shrink).ToSTL("wheel.stl", 300)

	wheelRetainer().ScaleUniform(shrink).ToSTL("retainer.stl", 300)

	antCap().ScaleUniform(shrink).ToSTL("antcap.stl", 300)
}
