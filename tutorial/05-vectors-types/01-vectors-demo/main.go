// Vectors demo: build a small assembly using v2, v3, and p2 constructors.
//
// Three pegs around a base disc, positioned in polar form via p2.RT() and
// converted to cartesian via v2.FromP2 + .ToV3(z). Touches every vector
// constructor flavour we use elsewhere in the docs.
package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
	"github.com/snowbldr/fluent-sdfx/vec/p2"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	// 2D base profile via v2.XY: a rounded square 30mm on a side.
	base := shape.Rect(v2.XY(30, 30), 4)

	// Three peg positions in polar form (degrees → radians via units.DtoR).
	pegs := make([]*solid.Solid, 3)
	for i := range pegs {
		theta := units.DtoR(90 + float64(i*120))
		pos := v2.FromP2(p2.RT(11, theta))
		pegs[i] = solid.Cylinder(8, 1.5, 0.4).Translate(pos.ToV3(2))
	}

	// Lift the whole part up by 2 in Z via v3.Z (single-axis constructor).
	part := base.Extrude(2).Translate(v3.Z(0)).Union(pegs...)
	part.STL("out.stl", 4.0)
}
