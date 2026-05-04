package main

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// CAD Challenge #16 Part A
func cc16a() *solid.Solid {
	base_w := 4.5
	base_d := 2.0
	base_h := 0.62
	base_radius := 0.5

	slot_l := 0.5 * 2
	slot_r := 0.38 / 2.0

	base_2d := shape.Rect(v2.XY(base_w, base_d), base_radius)
	slot := shape.Line(slot_l, slot_r)
	slot0 := slot.Translate(v2.XY(base_w/2, 0))
	slot1 := slot.Translate(v2.XY(-base_w/2, 0))
	slots := slot0.Union(slot1)
	base_2d = base_2d.Cut(slots)
	base_3d := base_2d.Extrude(base_h)

	hole_h := 0.75
	block_radius := 1.0
	block_w := 0.62
	block_l := base_h + 2.0*hole_h
	y_ofs := (base_d - block_w) / 2

	hole_radius := 0.625 / 2.0
	cb_radius := 1.25 / 2.0
	cb_depth := 0.12

	block_2d := shape.Line(block_l, block_radius)
	block_2d = block_2d.CutLine(v2.XY(0, 0), v2.XY(0, 1))
	block_3d := block_2d.Extrude(block_w)

	cb_3d := obj.CounterBoredHole3D(block_w, hole_radius, cb_radius, cb_depth).Translate(v3.XYZ(block_l/2, 0, 0))
	block_3d = block_3d.Cut(cb_3d)

	// rotate: first RotateX(-90), then RotateY(-90), then translate
	block_3d = block_3d.RotateX(-90).RotateY(-90).Translate(v3.XYZ(0, y_ofs, 0))

	return base_3d.Union(block_3d)
}

// CAD Challenge #16 Part B
func cc16b() *solid.Solid {
	base_w := 120.0
	base_d := 80.0
	base_h := 24.0
	base_radius := 25.0

	base_2d := shape.Rect(v2.XY(base_w, 2*base_d), base_radius)
	base_2d = base_2d.CutLine(v2.XY(0, 0), v2.XY(-1, 0))
	base_2d = base_2d.Translate(v2.XY(0, -base_d/2))

	base_hole_r := 14.0 / 2.0
	base_hole_yofs := (base_d / 2.0) - 25.0
	base_hole_xofs := (base_w / 2.0) - 25.0
	holes := []v2.Vec{v2.XY(base_hole_xofs, base_hole_yofs), v2.XY(-base_hole_xofs, base_hole_yofs)}
	c := shape.Circle(base_hole_r)
	holes_2d := c.Multi(holes...)
	base_2d = base_2d.Cut(holes_2d)

	// slotted hole
	slot_l := 20.0
	slot_r := 9.0
	slot_2d := shape.Line(slot_l, slot_r).Rotate(90).Translate(v2.XY(0, slot_l/2))
	base_2d = base_2d.Cut(slot_2d)

	base_3d := base_2d.Extrude(base_h)

	// rails
	rail_w := 15.0
	rail_zofs := (base_h - rail_w) / 2.0
	rail_3d := solid.Box(v3.XYZ(rail_w, base_d, rail_w), 0)
	rail0_3d := rail_3d.Translate(v3.XYZ(base_hole_xofs, 0, -rail_zofs))
	rail1_3d := rail_3d.Translate(v3.XYZ(-base_hole_xofs, 0, -rail_zofs))
	base_3d = base_3d.Cut(rail0_3d, rail1_3d)

	// surface recess
	recess_w := 40.0
	recess_h := 2.0
	recess_zofs := (base_h / 2.0) - recess_h
	recess := []v2.Vec{v2.XY(0, 0), v2.XY(recess_w, 0), v2.XY(recess_w+recess_h, recess_h), v2.XY(0, recess_h)}
	recess_2d := shape.Polygon(recess)
	recess_3d := recess_2d.Extrude(base_w).
		RotateX(90).
		RotateZ(-90).
		Translate(v3.XYZ(0, recess_w, recess_zofs))
	base_3d = base_3d.Cut(recess_3d)

	// Tool Support
	support_h := 109.0 - base_h
	support_w := 24.0
	support_base_w := 14.0
	support_theta := math.Atan(support_h / (support_w - support_base_w))
	support_xofs := support_h / math.Tan(support_theta)

	facets := 5
	support := shape.NewPoly()
	support.Add(base_w/2, -1)
	support.Add(base_w/2, 0)
	support.Add(base_hole_xofs, 0).Smooth(5.0, facets)
	support.Add(base_hole_xofs+support_xofs, support_h).Smooth(25.0, 3*facets)
	support.Add(-base_hole_xofs-support_xofs, support_h).Smooth(25.0, 3*facets)
	support.Add(-base_hole_xofs, 0).Smooth(5.0, facets)
	support.Add(-base_w/2, 0)
	support.Add(-base_w/2, -1)
	support_2d := support.Build()
	support_3d := support_2d.Extrude(support_w)

	// chamfered hole
	hole_h := 84.0 - base_h
	hole_r := 35.0 / 2.0
	chamfer_d := 2.0
	hole_3d := obj.ChamferedHole3D(support_w, hole_r, chamfer_d).Translate(v3.XYZ(0, hole_h, 0))
	support_3d = support_3d.Cut(hole_3d)

	// cut the sloped face of the support
	support_3d = support_3d.CutPlane(
		v3.XYZ(0, support_h, -support_w/2),
		v3.XYZ(0, math.Cos(support_theta), math.Sin(support_theta)),
	)

	// position the support
	support_yofs := (base_d - support_w) / 2.0
	support_3d = support_3d.RotateX(90).Translate(v3.XYZ(0, -support_yofs, base_h/2))

	// Gussets
	gusset_l := 20.0
	gusset_w := 3.0
	gusset_xofs := 37.0 / 2.0
	gusset_h := 12.53

	gusset_yofs := base_d / 2.0
	gusset_yofs -= support_base_w
	gusset_yofs -= gusset_h / math.Tan(support_theta)
	gusset_yofs -= gusset_h

	gusset := shape.NewPoly()
	gusset.Add(gusset_l, 0)
	gusset.Add(0, 0).Smooth(20.0, facets)
	gusset.Add(-gusset_l, gusset_l)
	gusset.Add(-gusset_l, 0)
	gusset_2d := gusset.Build()

	gusset_3d := gusset_2d.Extrude(gusset_w).
		RotateX(90).
		RotateZ(90).
		Translate(v3.XYZ(0, -gusset_yofs, base_h/2))

	gusset0_3d := gusset_3d.Translate(v3.XYZ(gusset_xofs, 0, 0))
	gusset1_3d := gusset_3d.Translate(v3.XYZ(-gusset_xofs, 0, 0))
	gusset_3d = gusset0_3d.Union(gusset1_3d)

	return base_3d.Union(support_3d, gusset_3d)
}
