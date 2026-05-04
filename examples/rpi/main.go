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

const shrink = 1.0 / 0.999 // PLA ~0.1%

const displayAngle = 15.0 // degrees from vertical
var tanTheta = math.Tan(units.DtoR(displayAngle))
var invCosTheta = 1.0 / math.Cos(units.DtoR(displayAngle))

const filletRadius = 10.0

const baseHeight = 8.0
const baseWidth = 100.0
const baseLength = 160.0

const baseFootX = 30.0
const baseFootY = 15.0
const baseHoleRadius = 2.0

var baseHolePosn = v2.XY(0.7, 0.8)

const supportPosn = 0.25 // fraction of baseWidth
const supportHeight = 120.0
const supportThickness = 5.0
const supportLength = 20.0

const webSize = 7.0
const webLength = 5.0

// 4 x M3 mounting holes on display
const displayW = 126.2
const displayH = 65.65
const displayHoleRadius = 0.5 * 3.9
const displayPosn = 0.7 // fraction of supportHeight

// sideProfile returns the 2d web/support profile
func sideProfile(t float64) *shape.Shape {
	p := shape.NewPoly()
	p.Add(0, 0)
	p.Add(baseWidth, 0).Rel()
	p.Add(-baseHeight*tanTheta, baseHeight).Rel()
	p.Add(-baseWidth*supportPosn, 0).Rel().Smooth(filletRadius, 5)
	p.Add(-supportHeight*tanTheta, supportHeight).Rel()
	p.Add(-invCosTheta*(supportThickness+t), 0).Rel()
	p.Add(tanTheta*(supportHeight-t), t-supportHeight).Rel().Smooth(filletRadius, 7)
	p.Add(0, baseHeight+t)
	p.Add(0, 0)
	return p.Build()
}

func webs() *solid.Solid {
	s2d := sideProfile(webSize)
	l := webLength
	s := s2d.Extrude(l)
	ofs := 0.5 * (baseLength - l)
	s0 := s.Translate(v3.XYZ(0, 0, ofs))
	s1 := s.Translate(v3.XYZ(0, 0, -ofs))
	return s0.Union(s1)
}

func supports() *solid.Solid {
	s2d := sideProfile(0)
	l := supportLength + webLength
	s := s2d.Extrude(l)
	ofs := 0.5 * (baseLength - l)
	s0 := s.Translate(v3.XYZ(0, 0, ofs))
	s1 := s.Translate(v3.XYZ(0, 0, -ofs))
	return s0.Union(s1)
}

func baseProfile() *shape.Shape {
	p := shape.NewPoly()
	p.Add(0, 0)
	p.Add(baseWidth, 0).Rel()
	p.Add(-baseHeight*tanTheta, baseHeight).Rel()
	p.Add(0, baseHeight)
	p.Add(0, 0)
	return p.Build()
}

func base() *solid.Solid {
	return baseProfile().Extrude(baseLength)
}

func baseCutout() *solid.Solid {
	holeSize := v2.XY(baseLength-2*baseFootX, baseWidth-2*baseFootY)
	s2d := shape.Rect(holeSize, filletRadius)
	s := s2d.Extrude(baseHeight)
	s = s.RotateX(90).RotateY(90).Translate(v3.XYZ(0.5*baseWidth, 0.5*baseHeight, 0))
	return s
}

func baseHole() *solid.Solid {
	s := obj.CounterSunkHole3D(baseHeight, baseHoleRadius)
	return s.Translate(v3.XYZ(0, 0, 0.5*baseHeight))
}

func baseHoles() *solid.Solid {
	s := baseHole()

	dx := 0.5 * baseHolePosn.X * baseWidth
	dy := 0.5 * baseHolePosn.Y * baseLength

	holes := s.Multi(v3.XYZ(dx, dy, 0), v3.XYZ(-dx, dy, 0), v3.XYZ(dx, -dy, 0), v3.XYZ(-dx, -dy, 0))
	return holes.RotateX(-90).Translate(v3.XYZ(0.5*baseWidth, 0, 0))
}

func displayHoles() *solid.Solid {
	s := solid.Cylinder(2*supportThickness, displayHoleRadius, 0)

	dx := 0.5 * displayW
	dy := 0.5 * displayH

	holes := s.Multi(v3.XYZ(dx, dy, 0), v3.XYZ(-dx, dy, 0), v3.XYZ(dx, -dy, 0), v3.XYZ(-dx, -dy, 0))
	holes = holes.RotateY(90).RotateZ(15)

	yOfs := displayPosn * supportHeight
	xOfs := (1-supportPosn)*baseWidth - (baseHeight * tanTheta) - (yOfs * tanTheta)
	return holes.Translate(v3.XYZ(xOfs, yOfs, 0))
}

func DisplayStand() *solid.Solid {
	return base().Union(webs(), supports()).Cut(baseCutout(), baseHoles(), displayHoles())
}

func main() {
	DisplayStand().ScaleUniform(shrink).STL("display_stand.stl", 3.0)
}
