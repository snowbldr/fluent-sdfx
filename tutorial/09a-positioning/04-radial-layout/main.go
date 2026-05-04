// Positioning: layout.Polar(radius, n) returns n positions evenly
// spaced on a circle in the XY plane. Spread them into the variadic
// .Multi(positions...) to stamp out a ring of copies. Equivalent to
// hand-writing a sin/cos loop, but a one-liner.
package main

import (
	"github.com/snowbldr/fluent-sdfx/layout"
	"github.com/snowbldr/fluent-sdfx/solid"
)

func main() {
	pillar := solid.Cylinder(20, 2, 0).BottomAt(0)
	pillars := pillar.Multi(layout.Polar(15, 8)...)
	pillars.STL("out.stl", 5.0)
}
