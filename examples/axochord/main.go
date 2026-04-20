package main

import (
	"math"
	"strings"

	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

var bR0 = 13.0 * 0.5          // major radius
var bR1 = 7.0 * 0.5           // minor radius
var bH0 = 6.0                 // cavity height for button body
var bH1 = 1.5                 // thru panel thickness
var bDeltaV = 22.0            // vertical inter-button distance
var bDeltaH = 20.0            // horizontal inter-button distance
var bTheta = units.DtoR(20.0) // button angle

const buttonsV = 3 // number of vertical buttons
const buttonsH = 3 //12 // number of horizontal buttons

func buttonCavity() *solid.Solid {
	p := shape.NewPoly()
	p.Add(0, -(bH0 + bH1))
	p.Add(bR0, 0).Rel()
	p.Add(0, bH0).Rel()
	p.Add(bR1-bR0, 0).Rel()
	p.Add(0, bH1).Rel()
	p.Add(bR0-bR1, 0).Rel()
	p.Add(0, bH0).Rel()
	p.Add(bR1-bR0, 0).Rel()
	p.Add(0, bH1).Rel()
	p.Add(-bR1, 0).Rel()
	return solid.Revolve(p.Build())
}

func buttons() *solid.Solid {
	// single key column
	d := buttonsV * bDeltaV
	p := v3.XYZ(-math.Sin(bTheta)*d, math.Cos(bTheta)*d, 0)
	bc := buttonCavity()
	col := bc.LineOf(v3.Zero, p, strings.Repeat("x", buttonsV))
	// multiple key columns
	d = buttonsH * bDeltaH
	p = v3.XYZ(d, 0, 0)
	matrix := col.LineOf(v3.Zero, p, strings.Repeat("x", buttonsH))
	// centered on the origin
	d = (buttonsV - 1) * bDeltaV
	dx := 0.5 * (((buttonsH - 1) * bDeltaH) - (d * math.Sin(bTheta)))
	dy := 0.5 * d * math.Cos(bTheta)
	return matrix.Translate(v3.XYZ(-dx, -dy, 0))
}

// https://geekhack.org/index.php?topic=47744.0
// https://cdn.sparkfun.com/datasheets/Components/Switches/MX%20Series.pdf

var cherryD0 = 0.551 * units.MillimetresPerInch
var cherryD1 = 0.614 * units.MillimetresPerInch
var cherryD2 = 0.1378 * units.MillimetresPerInch
var cherryD3 = 0.0386 * units.MillimetresPerInch

func cherryMX() *shape.Shape {
	cherryOfs := ((cherryD0 / 2.0) - cherryD3) - (cherryD2 / 2.0)

	r0 := shape.Rect(v2.XY(cherryD0, cherryD0), 0)
	r1 := shape.Rect(v2.XY(cherryD1, cherryD2), 0)

	r2 := r1.Translate(v2.XY(0, cherryOfs))
	r3 := r1.Translate(v2.XY(0, -cherryOfs))

	r4 := r2.Union(r3)
	r5 := r4.Rotate(90)

	return r0.Union(r4, r5)
}

func panel() *solid.Solid {
	v := (buttonsV - 1) * bDeltaV
	vx := float64(v) * math.Sin(bTheta)
	vy := float64(v) * math.Cos(bTheta)

	sx := (float64(buttonsH-1)*bDeltaH + vx) * 1.5
	sy := vy * 1.9

	return solid.Extrude(obj.Panel2D(obj.PanelParms{
		Size:         v2.XY(sx, sy),
		CornerRadius: 5.0,
		HoleDiameter: 3.0,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"xx", "x", "xx", "x"},
	}), 2.0*(bH0+bH1))
}

func main() {
	p := panel()
	b := buttons()
	s := p.Cut(b)
	upper := s.CutPlane(v3.Zero, v3.XYZ(0, 0, 1))
	lower := s.CutPlane(v3.Zero, v3.XYZ(0, 0, -1))

	upper.STL("upper.stl", 4.0)
	lower.STL("lower.stl", 4.0)

	cherryMX().ToDXF("plate.dxf", 400)
}
