package main

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func hex() *shape.Shape {
	return shape.Polygon(shape.Nagon(6, 20)).Offset(8)
}

func extrude1() *solid.Solid {
	h := hex()

	sLinear := h.Extrude(100)
	sFwd := h.TwistExtrude(100, units.Tau)
	sRev := h.TwistExtrude(100, -units.Tau)
	sCombo := sFwd.Union(sRev)

	d := 60.0
	return solid.UnionAll(
		sLinear.Translate(v3.Y(-1.5*d)),
		sFwd.Translate(v3.Y(-0.5*d)),
		sRev.Translate(v3.Y(0.5*d)),
		sCombo.Translate(v3.Y(1.5*d)),
	)
}

func extrude2() *solid.Solid {
	h := hex()

	s0 := h.ScaleExtrude(80, v2.XY(0.25, 0.5))
	s1 := h.ScaleTwistExtrude(80, math.Pi, v2.XY(0.25, 0.5))

	d := 30.0
	return s0.Translate(v3.Y(-d)).Union(s1.Translate(v3.Y(d)))
}

func main() {
	extrude1().STL("extrude1.stl", 3.0)
	extrude2().STL("extrude2.stl", 3.0)
}
