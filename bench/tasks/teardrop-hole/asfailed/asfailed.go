// Package asfailed: what a model that ignores the printability clause builds
// (or one that does not know the teardrop form): a plain round hole.
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func Build() *solid.Solid {
	block := solid.Box(v3.XYZ(20, 20, 10), 0)
	hole := solid.Cylinder(22, 3, 0).RotateY(90)
	return block.Cut(hole)
}
