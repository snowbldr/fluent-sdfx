package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%

const frontPanelThickness = 3.0
const frontPanelLength = 170.0
const frontPanelHeight = 50.0
const frontPanelYOffset = 15.0

const baseWidth = 50.0
const baseLength = 170.0
const baseThickness = 3.0

const baseFootWidth = 10.0
const baseFootCornerRadius = 3.0

const pcbWidth = 50.0
const pcbLength = 160.0

const pillarHeight = 16.8

// multiple standoffs
func standoffs() *solid.Solid {
	zOfs := 0.5 * (pillarHeight + baseThickness)

	// from the board mechanicals
	positions := []v3.Vec{
		v3.XYZ(3.5, 10.0, zOfs),   // H1
		v3.XYZ(3.5, 40.0, zOfs),   // H2
		v3.XYZ(54.0, 40.0, zOfs),  // H3
		v3.XYZ(156.5, 10.0, zOfs), // H4
		v3.XYZ(156.5, 40.0, zOfs), // H6
		v3.XYZ(44.0, 10.0, zOfs),  // H7
		v3.XYZ(116.0, 10.0, zOfs), // H8
	}

	return obj.Standoff3D(obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4,
	}).Multi(positions)
}

// base returns the base mount.
func base() *solid.Solid {
	base2d := obj.Panel2D(obj.PanelParms{
		Size:         v2.XY(baseLength, baseWidth),
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{7.0, 20.0, 7.0, 20.0},
		HolePattern:  [4]string{"xx", "x", "xx", "x"},
	})

	// cutout
	l := baseLength - (2.0 * baseFootWidth)
	w := 18.0
	cutoutYOfs := 0.5 * (baseWidth - pcbWidth)
	cutout := shape.Rect(v2.XY(l, w), baseFootCornerRadius).Translate(v2.Y(cutoutYOfs))

	xOfs := 0.5 * pcbLength
	yOfs := pcbWidth - (0.5 * baseWidth)
	s2 := base2d.Cut(cutout).Extrude(baseThickness).Translate(v3.XY(xOfs, yOfs))

	// standoffs + body blended with fillet
	return solid.SmoothUnion(solid.PolyMin(3.0), s2, standoffs())
}

// front panel cutouts

// button positions
const pbX = 53.0

var pb0 = v2.XY(pbX, 0.8)
var pb1 = v2.XY(pbX+5.334, 0.8)

// panelCutouts returns the 2D front panel cutouts
func panelCutouts() *shape.Shape {
	sMidi := shape.Circle(0.5 * 17.0)
	sJack0 := shape.Circle(0.5 * 11.5)
	sJack1 := shape.Circle(0.5 * 5.5)
	sLed := shape.Rect(v2.XY(1.6, 1.6), 0)

	sButton := obj.FingerButton2D(obj.FingerButtonParms{
		Width:  4.0,
		Gap:    0.6,
		Length: 20.0,
	}).Rotate(-90)

	jackX := 123.0
	midiX := 18.8
	ledX := 62.9

	type panelHole struct {
		center v2.Vec
		hole   *shape.Shape
	}

	holes := []panelHole{
		{v2.XY(midiX, 10.2), sMidi},                         // MIDI DIN Jack
		{v2.XY(midiX+20.32, 10.2), sMidi},                   // MIDI DIN Jack
		{v2.XY(jackX, 8.14), sJack0},                        // 1/4" Stereo Jack
		{v2.XY(jackX+19.5, 8.14), sJack0},                   // 1/4" Stereo Jack
		{v2.XY(107.6, 2.3), sJack1},                         // 3.5 mm Headphone Jack
		{v2.XY(ledX, 0.5), sLed},                            // LED
		{v2.XY(ledX+3.635, 0.5), sLed},                      // LED
		{pb0, sButton},                                      // Push Button
		{pb1, sButton},                                      // Push Button
		{v2.XY(84.1, 1.0), shape.Rect(v2.XY(16.0, 7.5), 0)}, // micro SD card
		{v2.XY(96.7, 1.0), shape.Rect(v2.XY(11.0, 7.5), 0)}, // micro USB connector
		{v2.XY(73.1, 7.1), shape.Rect(v2.XY(7.5, 15.0), 0)}, // fullsize USB connector
	}

	shapes := make([]*shape.Shape, len(holes))
	for i, h := range holes {
		shapes[i] = h.hole.Translate(h.center)
	}

	return shapes[0].Union(shapes[1:]...)
}

// frontPanel returns the front panel mount.
func frontPanel() *solid.Solid {
	xOfs := 0.5 * pcbLength
	yOfs := (0.5 * frontPanelHeight) - frontPanelYOffset
	panel2d := obj.Panel2D(obj.PanelParms{
		Size:         v2.XY(frontPanelLength, frontPanelHeight),
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"xx", "x", "xx", "x"},
	}).Translate(v2.XY(xOfs, yOfs))

	// Add buttons to the finger button
	bHeight := 4.0
	b := solid.Cylinder(bHeight, 1.4, 0)

	return panel2d.Cut(panelCutouts()).Extrude(frontPanelThickness).Union(
		b.Translate(pb0.ToV3(-0.5*bHeight)),
		b.Translate(pb1.ToV3(-0.5*bHeight)),
	)
}

func main() {
	// front panel
	s0 := frontPanel()
	sx := s0.RotateY(180)
	sx.ScaleUniform(shrink).STL("panel.stl", 4.0)

	// base
	s1 := base()
	s1.ScaleUniform(shrink).STL("base.stl", 4.0)

	// both together
	s0moved := s0.Translate(v3.Y(80))
	s3 := s0moved.Union(s1)
	s3.ScaleUniform(shrink).STL("panel_and_base.stl", 4.0)
}
