package main

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/layout"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// overall build controls
const shrink = 1.0 / 0.98 // 2% Al shrinkage
const core_print = false  // add the core print to the wheel
const pie_print = false   // create a 1/n pie segment (n = number of webs)

// draft angles (radians)
var draft_angle = 4.0 * math.Pi / 180       // standard overall draft
var core_draft_angle = 10.0 * math.Pi / 180 // draft angle for the core print

// nominal size values (mm)
const wheel_diameter = units.MillimetresPerInch * 8.0
const hub_diameter = 40.0
const hub_height = 53.0
const shaft_diameter = 21.0
const shaft_length = 45.0
const wall_height = 35.0
const wall_thickness = 4.0
const plate_thickness = 7.0
const web_width = 2.0
const web_height = 25.0
const core_height = 15.0
const number_of_webs = 6

// derived values
const wheel_radius = wheel_diameter / 2
const hub_radius = hub_diameter / 2
const shaft_radius = shaft_diameter / 2
const web_length = wheel_radius - wall_thickness - hub_radius

func wheel_profile() *shape.Shape {
	draft0 := (hub_height - plate_thickness) * math.Tan(draft_angle)
	draft1 := (wall_height - plate_thickness) * math.Tan(draft_angle)
	draft2 := wall_height * math.Tan(draft_angle)
	draft3 := core_height * math.Tan(core_draft_angle)

	s := shape.NewPoly()

	if core_print {
		s.Add(0, 0)
		s.Add(wheel_radius+draft2, 0)
		s.Add(wheel_radius, wall_height).Smooth(1.0, 5)
		s.Add(wheel_radius-wall_thickness, wall_height).Smooth(1.0, 5)
		s.Add(wheel_radius-wall_thickness-draft1, plate_thickness).Smooth(2.0, 5)
		s.Add(hub_radius+draft0, plate_thickness).Smooth(2.0, 5)
		s.Add(hub_radius, hub_height).Smooth(2.0, 5)
		s.Add(shaft_radius, hub_height)
		s.Add(shaft_radius-draft3, hub_height+core_height)
		s.Add(0, hub_height+core_height)
	} else {
		s.Add(0, 0)
		s.Add(wheel_radius+draft2, 0)
		s.Add(wheel_radius, wall_height).Smooth(1.0, 5)
		s.Add(wheel_radius-wall_thickness, wall_height).Smooth(1.0, 5)
		s.Add(wheel_radius-wall_thickness-draft1, plate_thickness).Smooth(2.0, 5)
		s.Add(hub_radius+draft0, plate_thickness).Smooth(2.0, 5)
		s.Add(hub_radius, hub_height).Smooth(2.0, 5)
		s.Add(shaft_radius, hub_height)
		s.Add(shaft_radius, hub_height-shaft_length)
		s.Add(0, hub_height-shaft_length)
	}

	return s.Build()
}

func web_profile() *shape.Shape {
	draft := web_height * math.Tan(draft_angle)
	x0 := web_width + draft
	x1 := web_width

	s := shape.NewPoly()
	s.Add(-x0, 0)
	s.Add(-x1, web_height).Smooth(1.0, 3)
	s.Add(x1, web_height).Smooth(1.0, 3)
	s.Add(x0, 0)

	return s.Build()
}

func wheel_pattern() *solid.Solid {
	// build reinforcing webs
	web := web_profile().Extrude(web_length)
	// equivalent of: m = Translate(0, plate_thickness, hub_radius+web_length/2)
	//                m = RotateX(90) * m  -> rotate by 90 around X axis
	web = web.Translate(v3.YZ(plate_thickness, hub_radius+web_length/2)).RotateX(90)

	var wheel *solid.Solid
	if pie_print {
		web = web.RotateZ(120)
		wheel = wheel_profile().RevolveAngle(60)
	} else {
		web = web.RotateZ(90)
		web = web.RotateCopyZ(number_of_webs)
		wheel = wheel_profile().Revolve()
	}

	// union with blend
	return solid.SmoothUnion(solid.PolyMin(wall_thickness), wheel, web)
}

func core_profile() *shape.Shape {
	draft := core_height * math.Tan(core_draft_angle)

	s := shape.NewPoly()
	s.Add(0, 0)
	s.Add(shaft_radius-draft, 0)
	s.Add(shaft_radius, core_height)
	s.Add(shaft_radius, core_height+shaft_length).Smooth(2.0, 3)
	s.Add(0, core_height+shaft_length)

	return s.Build()
}

func core_box() *solid.Solid {
	w := 4.2 * shaft_radius
	d := 1.2 * shaft_radius
	h := (core_height + shaft_length) * 1.1
	box := solid.Box(v3.XYZ(h, w, d), 0)

	// holes in the box
	dy := w * 0.37
	dx := h * 0.4
	hole_radius := ((3.0 / 16.0) * units.MillimetresPerInch) / 2.0
	holes := solid.Cylinder(d, hole_radius, 0).Multi(layout.RectCorners(2*dx, 2*dy)...)

	box = box.Cut(holes)

	// build the core
	core := core_profile().Revolve().
		RotateY(-90).
		Translate(v3.XZ(h/2, d/2))

	return box.Cut(core)
}

func main() {
	s0 := wheel_pattern().ScaleUniform(shrink)
	s0.STL("wheel.stl", 2.0)
	// DXF slice (non-STL output)
	slice := shape.SliceOf(s0, v3.XYZ(0, 0, 15), v3.XYZ(0, 0, 1))
	slice.ToDXF("wheel.dxf", 200)

	s1 := core_box().ScaleUniform(shrink)
	s1.STL("core_box.stl", 2.0)
}
