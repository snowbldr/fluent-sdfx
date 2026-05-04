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

const baseThickness = 3.0

func boardStandoffs() *solid.Solid {
	pillarHeight := 14.0
	zOfs := 0.5 * (pillarHeight + baseThickness)
	x := 82.0
	y := 54.0
	x0 := -34.0
	y0 := -0.5 * y
	return obj.Standoff3D(obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 4.5,
		HoleDepth:      11.0,
		HoleDiameter:   2.6,
		NumberWebs:     2,
		WebHeight:      10,
		WebDiameter:    12,
		WebWidth:       3.5,
	}).Multi(v3.XYZ(x0, y0, zOfs), v3.XYZ(x0+x, y0, zOfs), v3.XYZ(x0, y0+y, zOfs), v3.XYZ(x0+x, y0+y, zOfs))
}

func bezelStandoffs() *solid.Solid {
	pillarHeight := 22.0
	zOfs := 0.5 * (pillarHeight + baseThickness)
	x := 140.0
	y := 55.0
	x0 := -0.5 * x
	y0 := -0.5 * y
	return obj.Standoff3D(obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      11.0,
		HoleDiameter:   2.4,
	}).Multi(v3.XYZ(x0, y0, zOfs), v3.XYZ(x0+x, y0, zOfs), v3.XYZ(x0, y0+y, zOfs), v3.XYZ(x0+x, y0+y, zOfs))
}

func speakerHoles(d float64, ofs v2.Vec) *shape.Shape {
	holeRadius := 1.7
	s0 := shape.Circle(holeRadius)
	s1 := obj.BoltCircle2D(holeRadius, d*0.3, 6)
	return s0.Union(s1).Translate(ofs)
}

func speakerHolder(d float64, ofs v2.Vec) *solid.Solid {
	thickness := 3.0
	zOfs := 0.5 * (thickness + baseThickness)
	return obj.Washer3D(obj.WasherParms{
		Thickness:   thickness,
		InnerRadius: 0.5 * d,
		OuterRadius: 0.5 * (d + 4.0),
		Remove:      0.3,
	}).
		RotateZ(180).
		Translate(v3.XYZ(ofs.X, ofs.Y, zOfs))
}

func bezel() *solid.Solid {
	speakerOfs := v2.XY(60, 14)
	speakerDiameter := 20.3

	// bezel
	b0 := shape.Rect(v2.XY(150, 65), 2)

	// lcd cutout
	l0 := shape.Rect(v2.XY(60, 46), 2)

	// camera cutout
	c0 := shape.Circle(7.25).Translate(v2.X(42))

	// led hole cutout
	c1 := shape.Circle(2).Translate(v2.XY(44, -20))

	// speaker holes cutout
	c2 := speakerHoles(speakerDiameter, speakerOfs)

	// extrude the bezel
	s0 := b0.Cut(l0.Union(c0, c1, c2)).Extrude(baseThickness)

	// add the board standoffs
	s0 = s0.Union(boardStandoffs())

	// add the bezel standoffs (with foot rounding)
	s1 := solid.SmoothUnion(solid.PolyMin(3.0), s0, bezelStandoffs())

	// speaker holder
	s3 := speakerHolder(speakerDiameter, speakerOfs)

	return s1.Union(s3)
}

func main() {
	bezel().ScaleUniform(shrink).STL("bezel.stl", 3.3)
}
