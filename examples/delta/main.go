package main

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%

const upperArmWidth = 30.0
const upperArmRadius0 = 15.0
const upperArmRadius1 = 5.0
const upperArmRadius2 = 3.9 * 0.5
const upperArmLength = 100.0

func upperArm() *solid.Solid {
	const upperArmThickness = 5.0 * 2.0
	const gussetThickness = 0.5

	// body
	body := solid.Extrude(shape.FlatFlankCam(upperArmLength, upperArmRadius0, upperArmRadius1), upperArmThickness)

	// end cylinder
	c0 := solid.Cylinder(upperArmWidth*2.0, upperArmRadius1, 0).Translate(v3.Y(upperArmLength))

	// end cylinder hole
	c1 := solid.Cylinder(upperArmWidth*2.0, upperArmRadius2, 0).Translate(v3.Y(upperArmLength))

	// gusset
	const dx = upperArmWidth * 2.0 * 0.4
	const dy = upperArmLength * 0.6
	g := shape.NewPoly()
	g.Add(-dx, dy)
	g.Add(dx, dy)
	g.Add(0, 0)
	gusset := solid.Extrude(g.Build(), upperArmThickness*gussetThickness).
		RotateY(90).
		Translate(v3.Y(upperArmLength - dy))

	// servo mounting
	horn := solid.Extrude(obj.ServoHorn(obj.ServoHornParms{
		CenterRadius: 3,
		NumHoles:     4,
		CircleRadius: 14 * 0.5,
		HoleRadius:   1.9,
	}), upperArmThickness)

	const hornRadius = 10
	const hornThickness = 2.3
	hornBody := solid.Cylinder(hornThickness, hornRadius, 0).
		Translate(v3.Z((upperArmThickness - hornThickness) * 0.5))

	// body + cylinder, then gusset with PolyMin fillet
	s := solid.SmoothUnion(solid.PolyMin(upperArmThickness*gussetThickness), body.Union(c0), gusset)
	// remove the holes
	s = s.Cut(c1.Union(horn, hornBody))
	// cut in half
	s = s.CutPlane(v3.Zero, v3.Z(1))

	return s
}

const servoMountUprightLength = 66.0
const servoMountBaseLength = 35.0
const servoMountThickness = 3.5
const servoMountWidth = 35.0
const servoMountHoleRadius = 2.4

func servoMountHoles(h float64) *solid.Solid {
	hole := solid.Cylinder(h, servoMountHoleRadius, 0).
		Translate(v3.X((servoMountBaseLength + servoMountThickness) * 0.5))
	dx := (servoMountBaseLength * 0.5) - servoMountThickness - 4.0
	dy := (servoMountWidth * 0.5) - servoMountThickness - 6.0
	return hole.Multi([]v3.Vec{v3.XYZ(dx, dy, 0), v3.XYZ(-dx, dy, 0), v3.XYZ(dx, -dy, 0), v3.XYZ(-dx, -dy, 0)})
}

func servoMount() *solid.Solid {
	const servoOffset = servoMountUprightLength - 20.0

	m := shape.NewPoly()
	m.Add(0, 0)
	m.Add(servoMountBaseLength, 0)
	m.Add(servoMountBaseLength, servoMountThickness)
	m.Add(servoMountThickness, servoMountUprightLength)
	m.Add(0, servoMountUprightLength)
	mount := solid.Extrude(m.Build(), servoMountWidth)

	// cavity
	c := shape.NewPoly()
	c.Add(servoMountThickness, servoMountThickness)
	c.Add(servoMountBaseLength, servoMountThickness)
	c.Add(servoMountThickness, servoMountUprightLength)
	cavity := solid.Extrude(c.Build(), servoMountWidth-2*servoMountThickness)

	mountSolid := mount.Cut(cavity).RotateX(90)

	// base holes
	holes := servoMountHoles(servoMountThickness).
		Translate(v3.Z(servoMountThickness * 0.5))
	mountSolid = mountSolid.Cut(holes)

	// servo
	k := obj.ServoLookup("annimos_ds3218")
	servo := solid.Extrude(obj.Servo2D(*k, 2.1), servoMountThickness).
		RotateY(90).
		Translate(v3.XZ(servoMountThickness*0.5, servoOffset))

	return mountSolid.Cut(servo)
}

const baseSide = 150
const baseThickness = 7
const basePillarHeight = 20
const baseHoleRadius = 7

var servoY = -baseSide * math.Tan(30.0*math.Pi/180) * 0.5
var servoX = 25.0 - upperArmWidth*0.5

func deltaBase() *solid.Solid {
	// servo holes
	holes := servoMountHoles(baseThickness).
		Translate(v3.XYZ(servoX, servoY, -baseThickness*0.5)).
		RotateCopyZ(3)

	// base
	base := solid.Cylinder(baseThickness, baseSide*0.5*1.05, 0).
		Translate(v3.Z(-baseThickness * 0.5))

	// pillars
	pillars := obj.Standoff3D(obj.StandoffParms{
		PillarHeight:   basePillarHeight,
		PillarDiameter: 15,
		HoleDepth:      15,
		HoleDiameter:   3,
	}).
		RotateX(180).
		Translate(v3.YZ(-baseSide*0.4, -(0.5*basePillarHeight + baseThickness))).
		RotateCopyZ(3)

	// hole for servo wires
	baseHole := solid.Cylinder(baseThickness, baseHoleRadius, 0).
		Translate(v3.Z(-baseThickness * 0.5))
	holes = holes.Union(baseHole)

	// base + pillars blended, minus holes
	blended := solid.SmoothUnion(solid.PolyMin(baseThickness), base, pillars)
	return blended.Cut(holes)
}

const rodRadius = (6.0 + 0.5) * 0.5
const holderRadius = (11.7 + 0.2) * 0.5
const holderHeight = 2.6

func rodEnd() *solid.Solid {
	const endRadius = holderRadius * 1.5
	const endHeight = rodRadius * 2.0 * 1.5
	const round = endHeight * 0.1
	end := solid.Cylinder(endHeight, endRadius, round)

	const endX = 3 * endHeight
	box := solid.Box(v3.XYZ(endX, endHeight, endHeight), round).
		Translate(v3.X(0.5 * endX))

	const rodHole = (endX - endRadius) * 0.9
	const ofsX = endX - 0.5*rodHole
	rod := solid.Cylinder(rodHole, rodRadius, 0).
		RotateY(90).
		Translate(v3.X(ofsX))

	holder := solid.Cylinder(holderHeight, holderRadius, 0).
		Translate(v3.Z((endHeight - holderHeight) * 0.5))

	// end + box with fillets
	s := solid.SmoothUnion(solid.PolyMin(endRadius*0.1), end, box)

	// bump removal
	s = s.CutPlane(v3.Z(-endHeight*0.5), v3.Z(1))
	s = s.CutPlane(v3.Z(endHeight*0.5), v3.Z(-1))

	// remove the cavities
	return s.Cut(holder.Union(rod))
}

const platformSide = 50
const platformThickness = 10.0

func platform() *solid.Solid {
	pHalf := platformSide * 0.5
	pShort := pHalf / math.Sqrt(3)
	pLong := 2 * pShort

	c0 := v3.XYZ(0, -pLong, 0)
	c1 := v3.XYZ(pHalf, pShort, 0)
	c2 := v3.XYZ(-pHalf, pShort, 0)

	pp := shape.NewPoly()
	pp.Add(c0.X, c0.Y)
	pp.Add(c1.X, c1.Y)
	pp.Add(c2.X, c2.Y)
	platform := solid.Extrude(pp.Build(), platformThickness)

	// connection arms
	arm0 := obj.Pipe3D(platformThickness*0.5, upperArmRadius2, upperArmWidth).RotateY(90).Translate(c0)
	arm1 := arm0.RotateZ(120)
	arm2 := arm0.RotateZ(-120)

	s := solid.SmoothUnion(solid.PolyMin(platformThickness*0.7), platform, arm0, arm1, arm2)

	// bump removal
	s = s.CutPlane(v3.Z(-platformThickness*0.5), v3.Z(1))
	s = s.CutPlane(v3.Z(platformThickness*0.5), v3.Z(-1))

	return s
}

func main() {
	upperArm().ScaleUniform(shrink).STL("arm.stl", 5.0)
	servoMount().ScaleUniform(shrink).STL("servomount.stl", 2.5)
	deltaBase().ScaleUniform(shrink).STL("base.stl", 3.0)
	rodEnd().ScaleUniform(shrink).STL("rodend.stl", 1.0)
	platform().ScaleUniform(shrink).STL("platform.stl", 3.0)
}
