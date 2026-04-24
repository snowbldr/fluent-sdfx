package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Part A
func cc18a() {
	p := shape.NewPoly()
	p.Add(0, 0)
	p.Add(175, units.DtoR(-15)).Polar().Rel()
	p.Add(130, 0).Rel()
	p.Add(0, -25).Rel()
	p.Add(80, 0).Rel()
	p.Add(0, 25).Rel()
	p.Add(75, 0).Rel()
	p.Add(0, -75).Rel()
	p.Add(115, units.DtoR(-105)).Polar().Rel()
	p.Add(-50, 0).Rel()
	p.Add(150, units.DtoR(-195)).Polar().Rel().Arc(-120, 15)
	p.Add(100, units.DtoR(-150)).Polar().Rel()
	p.Add(-60, 0).Rel()
	p.Add(-10, 0).Rel()
	p.Add(-30, 0).Rel()
	p.Add(0, 135).Rel()
	p.Add(-60, 0).Rel()
	p.Close()
	render.Poly(p.Raw(), "cc18a.dxf")
}

// Part B
func cc18b() *solid.Solid {
	p := shape.NewPoly()
	p.Add(0, 0)
	p.Add(6, 0)
	p.Add(6, 19).Smooth(0.5, 5)
	p.Add(8, 19)
	p.Add(8, 21)
	p.Add(6, 21)
	p.Add(6, 20)
	p.Add(0, 20)
	vpipe_3d := p.Build().Revolve()

	// bolt circle for the top flange
	top_holes_3d := obj.BoltCircle3D(
		2.0,
		0.5/2.0,
		14.50/2.0,
		6,
	).RotateZ(30).Translate(v3.XYZ(0, 0, 1.0+19.0))
	vpipe_3d = vpipe_3d.Cut(top_holes_3d)

	// horizontal pipe
	p = shape.NewPoly()
	p.Add(0, 0)
	p.Add(5, 0)
	p.Add(5, 12).Smooth(0.5, 5)
	p.Add(8, 12)
	p.Add(8, 14)
	p.Add(6, 14)
	p.Add(6, 14.35)
	p.Add(0, 14.35)
	hpipe_3d := p.Build().Revolve()

	side_holes_3d := obj.BoltCircle3D(
		2.0,
		1.0/2.0,
		14.0/2.0,
		4,
	).RotateZ(45).Translate(v3.XYZ(0, 0, 1.0+12.0))
	hpipe_3d = hpipe_3d.Cut(side_holes_3d)

	hpipe_3d = hpipe_3d.RotateY(90).Union(hpipe_3d.RotateY(-90))
	hpipe_3d = hpipe_3d.Translate(v3.XYZ(0, 0, 9))

	s := solid.SmoothUnion(solid.PolyMin(1.0), hpipe_3d, vpipe_3d)

	// vertical blind hole
	vertical_hole_3d := solid.Cylinder(19.0, 9.0/2.0, 0.0).
		Translate(v3.XYZ(0, 0, 19.0/2.0+1))

	// horizontal through hole
	horizontal_hole_3d := solid.Cylinder(28.70, 9.0/2.0, 0.0).
		RotateY(90).
		Translate(v3.XYZ(0, 0, 9))

	return s.Cut(vertical_hole_3d, horizontal_hole_3d)
}

// Part C
func cc18c() *solid.Solid {
	// build the tabs
	tab_3d := solid.Box(v3.XYZ(43, 12, 20), 0).
		Translate(v3.XYZ(43.0/2.0, 0, 0))

	tab_hole_3d := solid.Cylinder(12, 7.0/2.0, 0).
		RotateX(90).
		Translate(v3.XYZ(35, 0, 0))
	tab_3d = tab_3d.Cut(tab_hole_3d)
	tab_3d = tab_3d.RotateCopyZ(3)

	// central body
	body_3d := solid.Cylinder(20, 26.3, 0)
	body_3d = solid.SmoothUnion(solid.PolyMin(2.0), body_3d, tab_3d)
	body_3d = body_3d.CutPlane(v3.XYZ(0, 0, -10), v3.XYZ(0, 0, 1))
	body_3d = body_3d.CutPlane(v3.XYZ(0, 0, 10), v3.XYZ(0, 0, -1))

	// sleeve
	r_outer := 42.3 / 2.0
	pts := []v2.Vec{v2.XY(0, 0), v2.XY(r_outer, 0), v2.XY(r_outer, 29), v2.XY(r_outer-1.0, 30), v2.XY(0, 30)}
	sleeve_3d := shape.Polygon(pts).Revolve().
		Translate(v3.XYZ(0, 0, -10))
	body_3d = body_3d.Union(sleeve_3d)

	sleeve_hole_3d := solid.Cylinder(30, 36.5/2.0, 0).
		Translate(v3.XYZ(0, 0, 5))
	return body_3d.Cut(sleeve_hole_3d)
}
