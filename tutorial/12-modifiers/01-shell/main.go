// Modifiers: Shell hollows out a solid, leaving a wall of the given thickness.
//
// Shell is closed on all sides; cut the top off to actually turn it into
// a cup. `BottomAt(7)` drops the cut tool so its bottom face sits at z=7
// — 3mm below the shelled cylinder's top — without thinking about the
// tool's own height.
package main

import "github.com/snowbldr/fluent-sdfx/solid"

func main() {
	solid.Cylinder(20, 12, 1).
		Shell(1.5).
		Cut(solid.Cylinder(4, 13, 0).BottomAt(7)).
		STL("out.stl", 6.0)
}
