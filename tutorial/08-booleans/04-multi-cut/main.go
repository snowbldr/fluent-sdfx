// Booleans: Cut is variadic — pass any number of tools and they're all
// removed from the body in a single operation. `layout.Polar(7, 4)`
// returns 4 positions on a 7mm-radius circle to spread into Multi for
// the outer ring of holes; the central hole is its own argument.
package main

import (
	"github.com/snowbldr/fluent-sdfx/layout"
	"github.com/snowbldr/fluent-sdfx/solid"
)

func main() {
	hole := solid.Cylinder(25, 1.5, 0)
	solid.Cylinder(20, 12, 1).
		Cut(
			hole, // central
			hole.Multi(layout.Polar(7, 4)...),
		).STL("out.stl", 4.0)
}
