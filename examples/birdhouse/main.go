package main

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const width = 120.0
const height = 85.0
const thickness = 2.0
const hookHeight = 10.0
const holeFactor = 0.9

func holeRadius() float64 {
	b := width * 0.5
	h := height
	a := math.Sqrt((b * b) + (h * h))
	return b * math.Sqrt((a-b)/(a+b))
}

func hook() *solid.Solid {
	return obj.Washer3D(obj.WasherParms{
		Thickness:   thickness,
		InnerRadius: hookHeight * 0.5,
		OuterRadius: hookHeight,
		Remove:      0.5,
	}).RotateY(90).Translate(v3.XYZ(0, 0, height+thickness))
}

func frame() *solid.Solid {
	pts := []v2.Vec{v2.XY(width/2, 0), v2.XY(0, height), v2.XY(-width/2, 0)}
	s := shape.Polygon(pts)
	sOuter := s.Offset(2 * thickness)
	sInner := s.Offset(thickness)
	f2d := sOuter.Cut(sInner)
	return f2d.Extrude(width * 1.1).RotateX(90)
}

func hole() *solid.Solid {
	r := holeRadius()
	return solid.Cylinder(2*width, r*holeFactor, 0).
		RotateX(90).
		Translate(v3.XYZ(0, 0, r))
}

func cross(s *solid.Solid) *solid.Solid {
	return s.Union(s.RotateZ(90))
}

func birdhouse() *solid.Solid {
	return cross(frame()).Cut(cross(hole())).Union(hook())
}

func main() {
	birdhouse().STL("birdhouse.stl", 3.0)
}
