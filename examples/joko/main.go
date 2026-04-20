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

// small end
const radiusOuterSmall = 1.0
const radiusInnerSmall = 0.55
const smallThickness = 1.0

// big end
const radiusOuterBig = 1.89
const radiusInnerBig = 2.90 * 0.5

const armWidth0 = 0.4
const armWidth1 = 0.5

const smallLength = 3.0
const overallLength = 9.75
const overallHeight = 4.0

var theta0 = units.DtoR(65.0) * 0.5
var theta1 = 0.5*units.Pi - theta0

const filletRadius0 = 0.25
const filletRadius1 = 0.50

const shaftRadius = 0.55
const keyRadius = 0.77
const keyWidth = 0.35

// derived
const centerToCenter = overallLength - radiusOuterBig - radiusOuterSmall

func planView() *shape.Shape {
	sOuter := shape.FlatFlankCam(centerToCenter, radiusOuterBig, radiusOuterSmall)
	sInner := sOuter.Offset(-armWidth0)
	s0 := sOuter.Cut(sInner)

	s1 := shape.Circle(radiusOuterSmall).Translate(v2.XY(0, centerToCenter))

	s2 := obj.Washer2D(obj.WasherParms{
		InnerRadius: radiusInnerBig,
		OuterRadius: radiusOuterBig,
	})

	s3 := shape.SmoothUnion(solid.PolyMin(0.3), s0, s1, s2)
	return sOuter.Intersect(s3)
}

const smoothSteps = 5

func sideView() *shape.Shape {
	dx0 := smallThickness * 0.5
	dy1 := smallLength
	dx2 := (overallHeight - smallThickness) * 0.5
	dy2 := dx2 * math.Tan(theta1)
	dy3 := overallLength - smallLength - dy2
	dx4 := -armWidth1
	dy5 := -dy3 + (armWidth1 / math.Cos(theta1)) - armWidth1*math.Tan(theta1)
	dx6 := armWidth1 - overallHeight*0.5
	dy6 := dx6 / math.Tan(theta0)

	p := shape.NewPoly()
	p.Add(dx0, 0)
	p.Add(0, dy1).Rel().Smooth(filletRadius1, smoothSteps)
	p.Add(dx2, dy2).Rel().Smooth(filletRadius1, smoothSteps)
	p.Add(0, dy3).Rel()
	p.Add(dx4, 0).Rel()
	p.Add(0, dy5).Rel().Smooth(filletRadius1, smoothSteps)
	p.Add(dx6, dy6).Rel().Smooth(filletRadius0, smoothSteps)
	// mirror
	p.Add(dx6, -dy6).Rel().Smooth(filletRadius1, smoothSteps)
	p.Add(0, -dy5).Rel()
	p.Add(dx4, 0).Rel()
	p.Add(0, -dy3).Rel().Smooth(filletRadius1, smoothSteps)
	p.Add(dx2, -dy2).Rel().Smooth(filletRadius1, smoothSteps)
	p.Add(0, -dy1).Rel()
	return p.Build()
}

func shaft() *solid.Solid {
	return obj.Keyway3D(obj.KeywayParameters{
		ShaftRadius: shaftRadius,
		KeyRadius:   keyRadius,
		KeyWidth:    keyWidth,
		ShaftLength: overallHeight,
	}).
		RotateY(-90).
		RotateX(-30).
		Translate(v3.XYZ(0, radiusOuterSmall, 0))
}

func part() *solid.Solid {
	side3d := solid.Extrude(sideView(), radiusOuterBig*2.0)

	plan3d := solid.Extrude(planView(), overallHeight).
		RotateZ(180).
		Translate(v3.XYZ(0, centerToCenter+radiusOuterSmall, 0)).
		RotateY(90)

	return plan3d.Intersect(side3d).Cut(shaft())
}

var _ = radiusInnerSmall // keep constant referenced

func main() {
	part().STL("part.stl", 3.0)
}
