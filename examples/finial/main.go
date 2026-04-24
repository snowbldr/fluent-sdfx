package main

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func square1(l float64) *shape.Shape {
	return shape.Polygon(shape.Nagon(4, l*math.Sqrt(0.5)))
}

func square2(l float64) *shape.Shape {
	h := l * 0.5
	r := l * 0.1
	n := 5

	p := shape.NewPoly()
	p.Add(h, -h).Smooth(r, n)
	p.Add(h, h).Smooth(r, n)
	p.Add(-h, h).Smooth(r, n)
	p.Add(-h, -h).Smooth(r, n)
	p.Close()
	return p.Build()
}

func finialFrom(base2d *shape.Shape) *solid.Solid {
	baseHeight := 20.0
	columnRadius := 15.0
	columnHeight := 60.0
	ballRadius := 45.0
	columnOfs := (columnHeight + baseHeight) / 2
	ballOfs := (baseHeight / 2) + columnHeight + ballRadius*0.8
	round := ballRadius / 5

	s1 := shape.Circle(columnRadius)
	column3d := base2d.LoftTo(s1, columnHeight, 0).Translate(v3.XYZ(0, 0, columnOfs))
	ball3d := solid.Sphere(ballRadius).Translate(v3.XYZ(0, 0, ballOfs))
	base3d := base2d.Extrude(baseHeight)

	bc3d := solid.SmoothUnion(solid.PolyMin(round), column3d, ball3d)
	return bc3d.Union(base3d)
}

func finial1() *solid.Solid {
	return finialFrom(square1(100.0))
}

func finial2() *solid.Solid {
	return finialFrom(square2(100.0))
}

func main() {
	finial1().STL("f1.stl", 3.0)
	finial2().STL("f2.stl", 3.0)
}
