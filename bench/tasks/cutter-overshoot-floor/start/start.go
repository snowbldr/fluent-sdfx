// Package start is the part before the change: a plain disc, radius 15 and
// 4 mm thick, sitting on the build plate (z = 0 .. 4).
package start

import "github.com/snowbldr/fluent-sdfx/solid"

const (
	discR = 15.0
	discH = 4.0 // z = 0 .. 4
)

func Build() *solid.Solid {
	return solid.Cylinder(discH, discR, 0).ZeroZ()
}
