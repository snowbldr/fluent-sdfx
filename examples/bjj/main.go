package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// center hole
const ch_d = 0.755 * units.MillimetresPerInch
const ch_r = ch_d / 2.0

func bushing() *solid.Solid {
	// R6-2RS 3/8 x 7/8 x 9/32 bearing
	bearing_outer_od := (7.0 / 8.0) * units.MillimetresPerInch
	bearing_inner_id := (3.0 / 8.0) * units.MillimetresPerInch
	bearing_inner_od := 12.0
	bearing_thickness := (9.0 / 32.0) * units.MillimetresPerInch

	clearance := 0.0

	r0 := 2.3
	r1 := (bearing_outer_od + bearing_inner_od) / 4.0
	r2 := (bearing_inner_id / 2.0) - clearance

	h0 := 3.0
	h1 := h0 + bearing_thickness + 1.0

	p := shape.NewPoly()
	p.Add(r0, 0)
	p.Add(r1, 0)
	p.Add(r1, h0)
	p.Add(r2, h0)
	p.Add(r2, h1)
	p.Add(r0, h1)

	return p.Build().Revolve()
}

func plateHoles2D() *shape.Shape {
	d := 17.0
	h := shape.Circle(1.2)
	return h.Multi([]v2.Vec{v2.XY(d, d), v2.XY(-d, -d), v2.XY(-d, d), v2.XY(d, -d)})
}

const rod_r = (1.0 / 16.0) * units.MillimetresPerInch * 1.10

func lockingRod() *solid.Solid {
	l := 62.0
	s0 := shape.Circle(rod_r)
	s1 := shape.Rect(v2.XY(2*rod_r, rod_r), 0).Translate(v2.XY(0, -0.5*rod_r))
	return s0.Union(s1).Extrude(l)
}

func plate() *solid.Solid {
	r := (16.0 * gear_module / 2.0) * 0.83
	h := 5.0

	s0 := solid.Cylinder(h, r, 0)

	ph := plateHoles2D()
	s1 := ph.Extrude(h)

	s2 := solid.Cylinder(h, ch_r, 0)

	lr := lockingRod()
	s3 := lr.RotateX(-90).Translate(v3.XYZ(0, 0, h/2-rod_r))

	return s0.Cut(s1, s2, s3)
}

var gear_module = 80.0 / 16.0
var pressure_angle = units.DtoR(20)
var involute_facets = 10

func gears() *solid.Solid {
	g_height := 10.0

	g0 := obj.InvoluteGear(obj.InvoluteGearParms{
		NumberTeeth:   12,
		Module:        gear_module,
		PressureAngle: pressure_angle,
		Facets:        involute_facets,
	}).Extrude(g_height)

	g1 := obj.InvoluteGear(obj.InvoluteGearParms{
		NumberTeeth:   16,
		Module:        gear_module,
		PressureAngle: pressure_angle,
		Facets:        involute_facets,
	}).Extrude(g_height)

	s0 := g0.Translate(v3.XYZ(0, 0, g_height/2.0))
	s1 := g1.Translate(v3.XYZ(0, 0, -g_height/2.0))

	s2 := solid.Cylinder(2.0*g_height, ch_r, 0)

	ph := plateHoles2D()
	screw_depth := 10.0
	s3 := ph.Extrude(screw_depth).Translate(v3.XYZ(0, 0, screw_depth/2.0-g_height))

	return s0.Union(s1).Cut(s2, s3)
}

func main() {
	bushing().STL("bushing.stl", 1.0)
	gears().STL("gear.stl", 3.0)
	plate().STL("plate.stl", 3.0)
}
