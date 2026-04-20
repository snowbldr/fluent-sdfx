package main

import (
	"fmt"
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

const scale = units.MillimetresPerInch

const draft = 2.0

const pinRadius = 0.5 * 4.4 * units.InchesPerMillimetre

type msParms struct {
	name          string
	size          float64
	width         float64
	wallThickness float64
	webThickness  float64
	holeRadius    float64
	holeOffset    float64
	allowance     float64
	pinRadius     float64
	nose          float64
}

// envelope for the outside machined/cast surfaces
func envelope(k *msParms, machined bool) *solid.Solid {
	c := k.nose
	l := k.size - c
	s0 := shape.Polygon([]v2.Vec{v2.XY(0, 0), v2.XY(l, 0), v2.XY(l, c), v2.XY(c, l), v2.XY(0, l)})

	if machined {
		return solid.Extrude(s0, k.width)
	}

	// cast
	s1 := solid.Extrude(s0, k.width+2.0*k.allowance)
	s2 := walls(k)
	return s1.Union(s2)
}

// wall returns an outside wall with casting draft of length l
func wall(k *msParms, l float64) *solid.Solid {
	ofs := (k.wallThickness - k.allowance) * 0.5
	return obj.TruncRectPyramid3D(obj.TruncRectPyramidParms{
		Size:        v3.XYZ(l, k.wallThickness+k.allowance, (k.width*0.5)+k.allowance),
		BaseAngle:   units.DtoR(90.0 - draft),
		BaseRadius:  (k.wallThickness + k.allowance) * 0.5,
		RoundRadius: k.wallThickness * 0.25,
	}).Translate(v3.XYZ(0, ofs, 0))
}

func walls(k *msParms) *solid.Solid {
	k0 := math.Sqrt(2.0)
	k1 := 1 + k0
	k2 := 2 + k0
	r := 0.5 * (k.wallThickness + k.allowance)

	// build the x-wall
	l0 := k.size + (k2 * k.allowance) - (k0 * r)
	w := wall(k, l0)
	ofs := 0.5*l0 - k.allowance
	w0 := w.Translate(v3.XYZ(ofs, 0, 0))

	// build the y-wall
	w1 := w0.MirrorXeqY()

	// build the 45-wall
	l1 := (k0 * k.size) + (2.0 * k1 * k.allowance) - (2.0 * k0 * r)
	w = wall(k, l1)
	ofs = 0.5 * k.size
	w2 := w.RotateZ(135).Translate(v3.XYZ(ofs, ofs, 0))

	// flipped walls
	w0f := w0.MirrorXY()
	w1f := w1.MirrorXY()
	w2f := w2.MirrorXY()

	return w0.Union(w1, w2, w0f, w1f, w2f)
}

// webHole returns the clamp hole within the web.
func webHole(k *msParms) *shape.Shape {
	r := k.holeRadius + 0.5*k.webThickness
	l := 2.0*k.holeOffset + k.webThickness
	s := shape.Line(l, r)
	return s.CutLine(v2.XY(0, 0), v2.XY(0, 1))
}

// web2d returns the 2d internal web.
func web2d(k *msParms) *shape.Shape {
	ofs := k.wallThickness * 0.9
	l := k.size - ofs*(2.0+math.Sqrt(2.0))

	s := shape.Polygon([]v2.Vec{v2.XY(ofs, ofs), v2.XY(ofs+l, ofs), v2.XY(ofs, ofs+l)})

	if k.holeRadius == 0 {
		return s
	}

	hole := webHole(k)

	k0 := k.size * 0.3
	k1 := k.size * 0.6
	k2 := k.size * 0.5

	h0 := hole.Translate(v2.XY(0, k0))
	h1 := hole.Translate(v2.XY(0, k1))

	holeR := hole.Rotate(90)
	h2 := holeR.Translate(v2.XY(k0, 0))
	h3 := holeR.Translate(v2.XY(k1, 0))

	holeR2 := hole.Rotate(135)
	h4 := holeR2.Translate(v2.XY(k2, k2))

	return s.Cut(h0.Union(h1, h2, h3, h4))
}

// web returns the internal web.
func web(k *msParms) *solid.Solid {
	s0 := web2d(k)
	return solid.ExtrudeRounded(s0, k.webThickness, 0.5*k.webThickness)
}

func corner90(k *msParms) *solid.Solid {
	r := 2.0 * k.wallThickness
	ofs := 0.8 * r
	return obj.TruncRectPyramid3D(obj.TruncRectPyramidParms{
		Size:        v3.XYZ(2.0*r, 2.0*r, (k.width*0.5)+k.allowance),
		BaseAngle:   units.DtoR(90.0 - 3.0*draft),
		BaseRadius:  r,
		RoundRadius: k.wallThickness * 0.25,
	}).Translate(v3.XYZ(ofs, ofs, 0))
}

func corner45(k *msParms) *solid.Solid {
	r := 2.3 * k.wallThickness
	dy := 0.7 * r
	dx := dy * (1.0 + math.Sqrt(2.0))
	return obj.TruncRectPyramid3D(obj.TruncRectPyramidParms{
		Size:        v3.XYZ(2.0*r, 2.0*r, (k.width*0.5)+k.allowance),
		BaseAngle:   units.DtoR(90.0 - 3.0*draft),
		BaseRadius:  r,
		RoundRadius: k.wallThickness * 0.25,
	}).Translate(v3.XYZ(k.size-dx, dy, 0))
}

func corners(k *msParms) *solid.Solid {
	c0 := corner90(k)

	c1 := corner45(k)
	c2 := c1.MirrorXeqY()

	c0f := c0.MirrorXY()
	c1f := c1.MirrorXY()
	c2f := c2.MirrorXY()

	return c0.Union(c0f, c1, c1f, c2, c2f)
}

func pin(k *msParms) *solid.Solid {
	return solid.Cylinder(k.width*0.8, k.pinRadius, 0)
}

// pins returns split-casting alignment pins
func pins(k *msParms) *solid.Solid {
	if k.pinRadius == 0 {
		return nil
	}

	// pin at 90 degree corner
	ofs := 1.5 * k.wallThickness
	p0 := pin(k).Translate(v3.XYZ(ofs, ofs, 0))

	// pins at 45 degree corners
	dy := 1.5 * k.wallThickness
	dx := dy * (1.0 + math.Sqrt(2.0))
	p1 := pin(k).Translate(v3.XYZ(k.size-dx, dy, 0))
	p2 := p1.MirrorXeqY()

	return p0.Union(p1, p2)
}

func mSquare(k *msParms, machined bool) {
	web := web(k)

	walls := walls(k)

	corners := corners(k)

	s := solid.SmoothUnion(solid.PolyMin(k.webThickness), web, walls, corners)

	// remove pin cavities
	if p := pins(k); p != nil {
		s = s.Cut(p)
	}

	// cleanup with the outside envelope
	env := envelope(k, machined)
	s = s.Intersect(env)

	s = s.ScaleUniform(shrink * scale)
	s.STL(fmt.Sprintf("%s.stl", k.name), 3.0)

	sUpper := s.CutPlane(v3.XYZ(0, 0, 0), v3.XYZ(0, 0, 1))
	sUpper.STL(fmt.Sprintf("%s_upper.stl", k.name), 3.0)

	sLower := s.CutPlane(v3.XYZ(0, 0, 0), v3.XYZ(0, 0, -1))
	sLower.STL(fmt.Sprintf("%s_lower.stl", k.name), 3.0)
}

func main() {
	k := &msParms{
		name:          "ms6",
		size:          6.0,
		width:         2.0,
		wallThickness: 0.25,
		webThickness:  0.25,
		holeRadius:    0.5,
		holeOffset:    0.75,
		allowance:     0.0625,
		pinRadius:     pinRadius,
		nose:          0.25,
	}

	mSquare(k, false)
}
