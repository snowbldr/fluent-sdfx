package main

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// overall build controls
const casting = false // add allowances, remove machined features

// scaling
const desired_scale = 1.25
const al_shrink = 1.0 / 0.99   // ~1%
const pla_shrink = 1.0 / 0.998 //~0.2%
const abs_shrink = 1.0 / 0.995 //~0.5%

const shrink = desired_scale * al_shrink * pla_shrink

const general_round = 0.1

// exhaust bosses

const eb_side_radius = 5.0 / 32.0
const eb_main_radius = 5.0 / 16.0
const eb_hole_radius = 3.0 / 16.0
const eb_c2c_distance = 13.0 / 16.0
const eb_distance = eb_c2c_distance / 2.0

const eb_x_offset = 0.5*(head_length+eb_height) - eb_height0
const eb_y_offset = (head_width / 2.0) - eb_distance - eb_side_radius
const eb_z_offset = 1.0 / 16.0

const eb_height0 = 1.0 / 16.0
const eb_height1 = 1.0 / 8.0
const eb_height = eb_height0 + eb_height1

func exhaust_boss(mode string, x_ofs float64) *solid.Solid {
	var s0 *shape.Shape

	switch mode {
	case "body":
		s0 = shape.Flange1(eb_distance, eb_main_radius, eb_side_radius)
	case "hole":
		s0 = shape.Circle(eb_hole_radius)
	default:
		panic("bad mode")
	}

	s1 := solid.Extrude(s0, eb_height)
	// original: rotateZ(90) then rotateY(90), then translate
	return s1.RotateZ(90).RotateY(90).Translate(v3.XYZ(x_ofs, eb_y_offset, eb_z_offset))
}

func exhaust_bosses(mode string) *solid.Solid {
	return exhaust_boss(mode, eb_x_offset).Union(exhaust_boss(mode, -eb_x_offset))
}

// spark plug bosses

const sp2sp_distance = 1.0 + (5.0 / 8.0)

var sp_theta = units.DtoR(30)

const sp_boss_r1 = 21.0 / 64.0
const sp_boss_r2 = 15.0 / 32.0
const sp_boss_h1 = 0.79
const sp_boss_h2 = 0.94
const sp_boss_h3 = 2

const sp_hole_d = 21.0 / 64.0
const sp_hole_r = sp_hole_d / 2.0
const sp_hole_h = 1.0

const sp_cb_h1 = 1.0
const sp_cb_h2 = 2.0
const sp_cb_r = 5.0 / 16.0

var sp_hyp = sp_hole_h + sp_cb_r*math.Tan(sp_theta)
var sp_y_ofs = sp_hyp*math.Cos(sp_theta) - head_width/2
var sp_z_ofs = -sp_hyp * math.Sin(sp_theta)

func sparkplug(mode string, x_ofs float64) *solid.Solid {
	var profile *shape.Shape
	switch mode {
	case "boss":
		p := shape.NewPoly()
		p.Add(0, 0)
		p.Add(sp_boss_r1, 0)
		p.Add(sp_boss_r1, sp_boss_h1).Smooth(sp_boss_r1*0.3, 3)
		p.Add(sp_boss_r2, sp_boss_h2).Smooth(sp_boss_r2*0.3, 3)
		p.Add(sp_boss_r2, sp_boss_h3)
		p.Add(0, sp_boss_h3)
		profile = p.Build()
	case "hole":
		vlist := []v2.Vec{v2.XY(0, 0), v2.XY(sp_hole_r, 0), v2.XY(sp_hole_r, sp_hole_h), v2.XY(0, sp_hole_h)}
		profile = shape.Polygon(vlist)
	case "counterbore":
		p := shape.NewPoly()
		p.Add(0, sp_cb_h1)
		p.Add(sp_cb_r, sp_cb_h1).Smooth(sp_cb_r/6.0, 3)
		p.Add(sp_cb_r, sp_cb_h2)
		p.Add(0, sp_cb_h2)
		profile = p.Build()
	default:
		panic("bad mode")
	}
	s := solid.Revolve(profile)
	// original: rotateX(Pi/2 - sp_theta), then translate
	// use radians -> degrees
	angleDeg := (math.Pi/2 - sp_theta) * 180 / math.Pi
	return s.RotateX(angleDeg).Translate(v3.XYZ(x_ofs, sp_y_ofs, sp_z_ofs))
}

func sparkplugs(mode string) *solid.Solid {
	x_ofs := 0.5 * sp2sp_distance
	return sparkplug(mode, x_ofs).Union(sparkplug(mode, -x_ofs))
}

// valve bosses

const valve_diameter = 1.0 / 4.0
const valve_radius = valve_diameter / 2.0
const valve_y_offset = 1.0 / 8.0
const valve_wall = 5.0 / 32.0
const v2v_distance = 1.0 / 2.0

var valve_draft = units.DtoR(5)

func valve(d float64, mode string) *solid.Solid {
	var s *solid.Solid
	h := head_height - cylinder_height

	switch mode {
	case "boss":
		delta := h * math.Tan(valve_draft)
		r1 := valve_radius + valve_wall
		r0 := r1 + delta
		s = solid.Cone(h, r0, r1, 0)
	case "hole":
		s = solid.Cylinder(h, valve_radius, 0)
	default:
		panic("bad mode")
	}

	z_ofs := cylinder_height / 2
	return s.Translate(v3.XYZ(d, valve_y_offset, z_ofs))
}

func valve_set(d float64, mode string) *solid.Solid {
	delta := v2v_distance / 2
	// blend the pair with PolyMin
	s := solid.SmoothUnion(solid.PolyMin(general_round), valve(-delta, mode), valve(delta, mode))
	return s.Translate(v3.X(d))
}

func valve_sets(mode string) *solid.Solid {
	delta := c2c_distance / 2
	return valve_set(-delta, mode).Union(valve_set(delta, mode))
}

// cylinder domes (or full base)

const cylinder_height = 3.0 / 16.0
const cylinder_diameter = 1.0 + (1.0 / 8.0)
const cylinder_wall = 1.0 / 4.0
const cylinder_radius = cylinder_diameter / 2.0

const dome_radius = cylinder_wall + cylinder_radius
const dome_height = cylinder_wall + cylinder_height

var dome_draft = units.DtoR(5)

const c2c_distance = 1.0 + (3.0 / 8.0)

func cylinder_head_part(d float64, mode string) *solid.Solid {
	var s *solid.Solid

	switch mode {
	case "dome":
		z_ofs := (head_height - dome_height) / 2
		extra_z := general_round * 2
		s = solid.Cylinder(dome_height+extra_z, dome_radius, general_round)
		s = s.Translate(v3.XZ(d, -z_ofs-extra_z))
	case "chamber":
		z_ofs := (head_height - cylinder_height) / 2
		s = solid.Cylinder(cylinder_height, cylinder_radius, 0)
		s = s.Translate(v3.XZ(d, -z_ofs))
	default:
		panic("bad mode")
	}
	return s
}

func cylinder_heads(mode string) *solid.Solid {
	x_ofs := c2c_distance / 2
	a := cylinder_head_part(-x_ofs, mode)
	b := cylinder_head_part(x_ofs, mode)
	if mode == "dome" {
		return solid.SmoothUnion(solid.PolyMin(general_round), a, b)
	}
	return a.Union(b)
}

// cylinder studs: location, bosses and holes

const stud_hole_radius = 1.0 / 16.0
const stud_boss_radius = 3.0 / 16.0
const stud_hole_dy = 11.0 / 16.0
const stud_hole_dx0 = 7.0 / 16.0
const stud_hole_dx1 = 1.066

var stud_locations = []v2.Vec{
	v2.XY(stud_hole_dx0+stud_hole_dx1, 0),
	v2.XY(stud_hole_dx0+stud_hole_dx1, stud_hole_dy),
	v2.XY(stud_hole_dx0+stud_hole_dx1, -stud_hole_dy),
	v2.XY(stud_hole_dx0, stud_hole_dy),
	v2.XY(stud_hole_dx0, -stud_hole_dy),
	v2.XY(-stud_hole_dx0-stud_hole_dx1, 0),
	v2.XY(-stud_hole_dx0-stud_hole_dx1, stud_hole_dy),
	v2.XY(-stud_hole_dx0-stud_hole_dx1, -stud_hole_dy),
	v2.XY(-stud_hole_dx0, stud_hole_dy),
	v2.XY(-stud_hole_dx0, -stud_hole_dy),
}

func head_stud_holes() *solid.Solid {
	c := shape.Circle(stud_hole_radius).Multi(stud_locations)
	return solid.Extrude(c, head_height)
}

// head walls

const head_length = 4.30 / 1.25
const head_width = 2.33 / 1.25
const head_height = 7.0 / 8.0
const head_corner_round = (5.0 / 32.0) / 1.25
const head_wall_thickness = 0.154

func head_wall_outer_2d() *shape.Shape {
	return shape.Rect(v2.XY(head_length, head_width), head_corner_round)
}

func head_wall_inner_2d() *shape.Shape {
	l := head_length - (2 * head_wall_thickness)
	w := head_width - (2 * head_wall_thickness)
	s0 := shape.Rect(v2.XY(l, w), 0)
	s1 := shape.Circle(stud_boss_radius).Multi(stud_locations)
	return shape.SmoothCut(solid.PolyMax(general_round), s0, s1)
}

func head_envelope() *solid.Solid {
	s0 := shape.Rect(v2.XY(head_length+2*eb_height1, head_width), 0)
	return solid.Extrude(s0, head_height)
}

func head_wall() *solid.Solid {
	s := head_wall_outer_2d().Cut(head_wall_inner_2d())
	return solid.Extrude(s, head_height)
}

// manifolds

const manifold_radius = 4.5 / 16.0
const manifold_hole_radius = 1.0 / 8.0
const inlet_theta = 30.2564
const exhaust_theta = 270.0 + 13.9736
const exhaust_x_offset = (c2c_distance / 2) + (v2v_distance / 2)
const inlet_x_offset = (c2c_distance / 2) - (v2v_distance / 2)

func manifold_set(r float64) *solid.Solid {
	const h = 2

	// exhaust cylinder: build, translate by +Z h/2, rotateX(-90), rotateZ(exhaust_theta), translate
	s_ex := solid.Cylinder(h, r, 0).
		Translate(v3.Z(h / 2)).
		RotateX(-90).
		RotateZ(exhaust_theta).
		Translate(v3.XYZ(exhaust_x_offset, valve_y_offset, eb_z_offset))

	s_in := solid.Cylinder(h, r, 0).
		Translate(v3.Z(h / 2)).
		RotateX(-90).
		RotateZ(inlet_theta).
		Translate(v3.XYZ(inlet_x_offset, valve_y_offset, eb_z_offset))

	return s_ex.Union(s_in)
}

func manifolds(mode string) *solid.Solid {
	var r float64
	switch mode {
	case "body":
		r = manifold_radius
	case "hole":
		r = manifold_hole_radius
	default:
		panic("bad mode")
	}
	s0 := manifold_set(r)
	s1 := s0.MirrorYZ()
	if mode == "body" {
		return solid.SmoothUnion(solid.PolyMin(general_round), s0, s1)
	}
	return s0.Union(s1)
}

func additive() *solid.Solid {
	s := solid.SmoothUnion(
		solid.PolyMin(general_round),
		head_wall(),
		cylinder_heads("dome"),
		valve_sets("boss"),
		sparkplugs("boss"),
		manifolds("body"),
		exhaust_bosses("body"),
	)

	s = s.Cut(sparkplugs("counterbore"))

	// cleanup the blending artifacts on the outside
	s = s.Intersect(head_envelope())

	return s
}

func subtractive() *solid.Solid {
	if casting {
		return nil
	}
	return solid.UnionAll(
		cylinder_heads("chamber"),
		head_stud_holes(),
		valve_sets("hole"),
		sparkplugs("hole"),
		manifolds("hole"),
		exhaust_bosses("hole"),
	)
}

func main() {
	result := additive()
	if sub := subtractive(); sub != nil {
		result = result.Cut(sub)
	}
	result.ScaleUniform(shrink).STL("head.stl", 4.0)
}
