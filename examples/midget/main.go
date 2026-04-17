package main

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const shrink = 1.0 / 0.999 // PLA ~0.1%

//-----------------------------------------------------------------------------
// Cylinder Pattern and Core Box

const cylinderBaseOffset = 3.0 / 16.0
const cylinderBaseThickness = 0.25
const cylinderWaistLength = 0.75
const cylinderBodyLength = 1.75
const cylinderCoreLength = 4.0 + (7.0 / 16.0)

const cylinderInnerRadius = 1.0 * 0.5
const cylinderWaistRadius = 1.5 * 0.5
const cylinderBodyRadius = 2.0 * 0.5

func cylinderBase() *solid.Solid {
	const draft = 3.0

	const x = cylinderBodyRadius * 2.0
	const y = cylinderBaseThickness * 2.0
	const z = cylinderBodyRadius

	const round = 0.125

	k := obj.TruncRectPyramidParms{
		Size:        v3.XYZ(x, y, z),
		BaseAngle:   units.DtoR(90 - draft),
		BaseRadius:  round,
		RoundRadius: round * 1.5,
	}
	base0 := obj.TruncRectPyramid3D(k)
	base1 := base0.MirrorXY()
	base := base0.Union(base1)
	base = base.CutPlane(v3.XYZ(0, 0, 0), v3.XYZ(0, 1, 0))
	return base.RotateX(90)
}

func cylinderPattern(core, split bool) *solid.Solid {
	draft := math.Tan(units.DtoR(3.0))
	const smooth0 = 0.125
	const smooth1 = smooth0 * 0.5
	const smoothN = 5

	const l0 = cylinderBaseOffset + cylinderBaseThickness + cylinderWaistLength
	const l1 = cylinderBodyLength
	const l2 = cylinderCoreLength

	const r0 = cylinderInnerRadius
	const r1 = cylinderWaistRadius
	const r2 = cylinderBodyRadius

	// cylinder body
	p := shape.NewPoly()
	p.Add(0, 0)
	p.Add(r1, draft*r1).Rel().Smooth(smooth1, smoothN)
	p.Add(0, l0).Rel().Smooth(smooth0, smoothN)
	p.Add(r2-r1, draft*(r2-r1)).Rel().Smooth(smooth0, smoothN)
	p.Add(0, l1).Rel().Smooth(smooth1, smoothN)
	p.Add(-r2, draft*r2).Rel()
	body := solid.Revolve(p.Build())

	// cylinder base
	base := cylinderBase().Translate(v3.XYZ(0, 0, cylinderBaseOffset))

	body = body.Union(base)

	// core print
	p = shape.NewPoly()
	p.Add(0, -0.75)
	p.Add(r0, draft*r0).Rel().Smooth(smooth1, smoothN)
	p.Add(0, l2).Rel().Smooth(smooth1, smoothN)
	p.Add(-r0, draft*r0).Rel()
	corePrint := solid.Revolve(p.Build())

	var cylinder *solid.Solid
	if core {
		cylinder = body.Union(corePrint)
	} else {
		cylinder = body.Cut(corePrint)
	}

	if split {
		cylinder = cylinder.CutPlane(v3.XYZ(0, 0, 0), v3.XYZ(0, 1, 0))
	}

	return cylinder
}

//-----------------------------------------------------------------------------
// Crankcase Pattern and Core Box

const crankcaseOuterRadius = 1.0 + (5.0 / 16.0)
const crankcaseInnerRadius = 1.0 + (1.0 / 8.0)
const crankcaseOuterHeight = 7.0 / 8.0
const boltLugRadius = 0.5 * (7.0 / 16.0)

func mountLugs() *solid.Solid {
	const draft = 3.0
	const thickness = 0.25

	k := obj.TruncRectPyramidParms{
		Size:        v3.XYZ(4.75, thickness, crankcaseOuterHeight),
		BaseAngle:   units.DtoR(90 - draft),
		BaseRadius:  crankcaseOuterHeight * 0.1,
		RoundRadius: crankcaseOuterHeight * 0.1,
	}

	return obj.TruncRectPyramid3D(k).Translate(v3.XYZ(0, thickness*0.5, 0))
}

func cylinderMount() *solid.Solid {
	const draft = 3.0

	k := obj.TruncRectPyramidParms{
		Size:        v3.XYZ(2.0, 5.0/16.0, 1+(3.0/16.0)),
		BaseAngle:   units.DtoR(90 - draft),
		BaseRadius:  crankcaseOuterHeight * 0.1,
		RoundRadius: crankcaseOuterHeight * 0.1,
	}

	return obj.TruncRectPyramid3D(k).Translate(v3.XYZ(0, crankcaseInnerRadius, 0))
}

func boltLugs() *solid.Solid {
	const draft = 3.0

	k := obj.TruncRectPyramidParms{
		Size:        v3.XYZ(0, 0, crankcaseOuterHeight),
		BaseAngle:   units.DtoR(90 - draft),
		BaseRadius:  boltLugRadius,
		RoundRadius: crankcaseOuterHeight * 0.1,
	}
	lug := obj.TruncRectPyramid3D(k)

	r := crankcaseOuterRadius
	d := r * math.Cos(units.DtoR(45))
	dy0 := 0.75
	dx0 := -math.Sqrt(r*r - dy0*dy0)
	positions := []v3.Vec{v3.XYZ(dx0, dy0, 0), v3.XYZ(1.0, 13.0/16.0, 0), v3.XYZ(-d, -d, 0), v3.XYZ(d, -d, 0)}

	return lug.Multi(positions)
}

func basePattern() *solid.Solid {
	const draft = 3.0

	k := obj.TruncRectPyramidParms{
		Size:        v3.XYZ(0, 0, crankcaseOuterHeight),
		BaseAngle:   units.DtoR(90 - draft),
		BaseRadius:  crankcaseOuterRadius,
		RoundRadius: crankcaseOuterHeight * 0.1,
	}
	body := obj.TruncRectPyramid3D(k)

	bl := boltLugs()
	ml := mountLugs()
	s := solid.SmoothUnion(solid.PolyMin(0.1), body, bl, ml)

	// cleanup top artifacts
	s = s.CutPlane(v3.XYZ(0, 0, crankcaseOuterHeight), v3.XYZ(0, 0, -1))

	// add cylinder mount
	cm := cylinderMount()
	s = solid.SmoothUnion(solid.PolyMin(0.1), s, cm)

	// cleanup bottom
	s = s.CutPlane(v3.XYZ(0, 0, 0), v3.XYZ(0, 0, 1))

	return s
}

func ccFrontPattern() *solid.Solid {
	return basePattern()
}

//-----------------------------------------------------------------------------

func main() {
	const scale = shrink * units.MillimetresPerInch

	cp := cylinderPattern(true, true)
	cp.ScaleUniform(scale).ToSTL("cylinder_pattern.stl", 330)

	ccfp := ccFrontPattern()
	ccfp.ScaleUniform(scale).ToSTL("crankcase_front.stl", 300)
}
