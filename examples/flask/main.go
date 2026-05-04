package main

import (
	"fmt"
	"math"

	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const shrink = 1.0 / 0.999 // PLA ~0.1%

const wallThickness = 5.0
const padThickness = 5.0
const padWidth = 60.0
const padDraft = 30.0
const cornerThickness = 7.0
const cornerLength = 30.0
const keyDepth = 4.0
const keyDraft = 60.0
const keyRatio = 0.85
const sideDraft = 3.0
const lugBaseThickness = 3.0
const lugBaseDraft = 15.0
const lugHeight = 28.0
const lugThickness = 14.0
const lugDraft = 5.0
const lugOffset = 1.5
const holeRadius = 1.5

const lugBaseWidth = padWidth * 0.95

// alignmentHoles returns an SDF3 for the alignment holes between the flask and pin lugs pattern.
func alignmentHoles() *solid.Solid {
	w := lugBaseWidth
	h := (lugBaseThickness + padThickness + wallThickness + cornerLength) * 2.0
	xofs := w * 0.8 * 0.5
	cylinder := solid.Cylinder(h, holeRadius, 0)
	return cylinder.Multi(v3.XYZ(xofs, 0, 0), v3.XYZ(-xofs, 0, 0))
}

// pinLug returns a single pin lug.
func pinLug(w float64) *solid.Solid {
	return obj.TruncRectPyramid3D(obj.TruncRectPyramidParms{
		Size:        v3.XYZ(w, lugThickness, lugHeight),
		BaseAngle:   units.DtoR(90 - lugDraft),
		BaseRadius:  lugThickness * 0.5,
		RoundRadius: lugThickness * 0.1,
	})
}

// pinLugs returns the pin lugs pattern.
func pinLugs() *solid.Solid {
	w := lugBaseWidth
	base := obj.TruncRectPyramid3D(obj.TruncRectPyramidParms{
		Size:        v3.XYZ(w, w, lugBaseThickness),
		BaseAngle:   units.DtoR(90 - lugBaseDraft),
		BaseRadius:  lugThickness*0.5 + lugOffset,
		RoundRadius: lugBaseThickness * 0.25,
	})

	pinWidth := w - 2.0*lugOffset
	pin := pinLug(pinWidth)
	yofs := 0.5 * (pinWidth - lugThickness)
	pin0 := pin.Translate(v3.XYZ(0, yofs, lugBaseThickness))
	pin1 := pin.Translate(v3.XYZ(0, -yofs, lugBaseThickness))

	s := solid.SmoothUnion(solid.PolyMin(lugBaseThickness*0.75), base, pin0, pin1)

	holes := alignmentHoles()
	return s.Cut(holes)
}

// sandKey returns an internal sand key.
func sandKey(size v3.Vec) *solid.Solid {
	theta := units.DtoR(90 - keyDraft)
	return obj.TruncRectPyramid3D(obj.TruncRectPyramidParms{
		Size:        size,
		BaseAngle:   theta,
		BaseRadius:  keyDepth / math.Tan(theta),
		RoundRadius: size.X * 0.5,
	})
}

// oddSide returns odd sides at either end of the flask pattern.
func oddSide(height float64) *solid.Solid {
	theta45 := units.DtoR(45)

	d := cornerLength * math.Cos(theta45)
	sx := 2.0*d + cornerThickness
	sy := height*1.1 + 2.0*d
	sz := d

	base := obj.TruncRectPyramid3D(obj.TruncRectPyramidParms{
		Size:        v3.XYZ(sx, sy, sz),
		BaseAngle:   theta45,
		BaseRadius:  0.5 * sx,
		RoundRadius: 0,
	})

	// mounting/pull holes
	h := 3.0 * d
	yofs := (height*1.1 - cornerThickness) * 0.5
	hole := solid.Cylinder(h, holeRadius, 0)
	holes := hole.Multi(v3.XYZ(0, yofs, 0), v3.XYZ(0, -yofs, 0))

	// hook into internal sand key
	sx2 := 0.8 * sx
	sy2 := height * keyRatio * 0.99
	sz2 := keyDepth
	key := sandKey(v3.XYZ(sx2, sy2, sz2)).Translate(v3.XYZ(0.5*sx2, 0, 0))

	return base.Union(key).Cut(holes)
}

// sideDraftProfile returns the 2d profile for the side draft of the flask pattern.
func sideDraftProfile(height float64) *shape.Shape {
	h0 := keyDepth + wallThickness + cornerLength
	w0 := height * 0.5
	w1 := w0 + w0
	w2 := w0 - h0*math.Tan(units.DtoR(sideDraft))

	p := shape.NewPoly()
	p.Add(w0, 0)
	p.Add(w1, 0)
	p.Add(w1, h0)
	p.Add(w2, h0)

	s0 := p.Build()
	s1 := s0.MirrorY()
	return s0.Union(s1)
}

// flaskSideProfile returns a half 2D extrusion profile for the flask.
func flaskSideProfile(width float64) *shape.Shape {
	theta45 := units.DtoR(45)
	theta135 := units.DtoR(135)
	theta225 := units.DtoR(225)

	w0 := width * 0.5
	w1 := padWidth * 0.5
	w2 := w1 + padThickness*math.Tan(units.DtoR(padDraft))

	h0 := keyDepth + wallThickness
	h1 := keyDepth + wallThickness + padThickness

	l0 := cornerLength + cornerThickness - (keyDepth+wallThickness)/math.Sin(theta45)

	r0 := cornerThickness * 0.25
	r1 := cornerThickness
	r2 := padThickness * 0.4

	p := shape.NewPoly()
	p.Add(0, 0)
	p.Add(w0, 0)
	p.Add(cornerLength, theta45).Polar().Rel().Smooth(r0, 4)
	p.Add(cornerThickness, theta135).Polar().Rel().Smooth(r0, 4)
	p.Add(l0, theta225).Polar().Rel().Smooth(r1, 4)
	p.Add(w2, h0).Smooth(r2, 3)
	p.Add(w1, h1).Smooth(r2, 3)
	p.Add(0, h1)

	return p.Build()
}

// pullHoles returns pull holes for the flask.
func pullHoles(width float64) *solid.Solid {
	h := (wallThickness + keyDepth) * 2.0
	xofs := width * 0.9 * 0.5
	hole := solid.Cylinder(h, holeRadius, 0)
	return hole.Multi(v3.XYZ(xofs, 0, 0), v3.XYZ(-xofs, 0, 0))
}

func flaskHalf(width, height float64) *solid.Solid {
	return flaskSideProfile(width).Extrude(height)
}

func flaskSide(width, height float64) *solid.Solid {
	side0 := flaskHalf(width, height)
	side1 := side0.MirrorYZ()
	flaskBody := side0.Union(side1)

	w := width + 2.0*cornerLength

	key := sandKey(v3.XYZ(w, height*keyRatio, keyDepth)).RotateX(-90)

	sd := sideDraftProfile(height)
	sideDraft3D := sd.Extrude(w).RotateY(90)

	aHoles := alignmentHoles().RotateX(90)
	pHoles := pullHoles(width).RotateX(90)

	return flaskBody.Cut(key, sideDraft3D, aHoles, pHoles)
}

func main() {
	widths := []float64{150, 200, 250, 300}
	height := 95.0
	for _, w := range widths {
		s := flaskSide(w, height).RotateX(-sideDraft).ScaleUniform(shrink)
		name := fmt.Sprintf("flask_%d.stl", int(w))
		s.STL(name, 3.0)
	}

	pinLugs().ScaleUniform(shrink).STL("pins.stl", 1.2)
	oddSide(height).ScaleUniform(shrink).STL("odd_side.stl", 3.0)
}
