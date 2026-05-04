// Positioning: layout.RectCorners returns the 4 XY corners of a
// rectangle centered on the origin — perfect for standoffs around a
// panel. layout.Grid does the same for an arbitrary nx*ny grid.
package main

import (
	"github.com/snowbldr/fluent-sdfx/layout"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	panel := solid.Box(v3.XYZ(80, 50, 2), 0).BottomAt(0)
	standoff := solid.Cylinder(10, 3, 0).BottomAt(2)

	standoff.Multi(layout.RectCorners(80-12, 50-12)...).
		Union(panel).STL("out.stl", 5.0)
}
