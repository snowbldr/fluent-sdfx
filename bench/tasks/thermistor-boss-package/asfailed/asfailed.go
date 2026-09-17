// Package asfailed: the model assumed a RADIAL two-lead package and put the
// two lead holes side by side, 1.3 mm apart horizontally, at mid-height
// (mcweed U30: "Right now the wires are right next to each other").
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func Build() *solid.Solid {
	core := solid.Cylinder(30, 10, 0)
	boss := solid.Cylinder(8, 4, 0).Translate(v3.X(10))
	hole := solid.Cylinder(8, 0.4, 0).RotateY(90).Translate(v3.X(10))
	holes := hole.Translate(v3.Y(0.65)).Union(hole.Translate(v3.Y(-0.65)))
	return core.Union(boss).Cut(holes)
}
