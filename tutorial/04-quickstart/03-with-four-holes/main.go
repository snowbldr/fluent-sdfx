// Quickstart step 3: drill four holes via Multi — places copies of the
// tool at the given positions, then Cut subtracts them all in one op.
//
// `layout.Polar(5, 4)` returns 4 positions evenly spaced on a circle of
// radius 5 in the XY plane — ready to spread straight into Multi.
package main

import (
	"github.com/snowbldr/fluent-sdfx/layout"
	"github.com/snowbldr/fluent-sdfx/solid"
)

func main() {
	solid.Cylinder(20, 10, 0).
		Cut(solid.Cylinder(25, 2, 0).Multi(layout.Polar(5, 4)...)).
		STL("out.stl", 3.0)
}
