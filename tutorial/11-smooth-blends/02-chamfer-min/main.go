// Smooth blends: ChamferMin produces a 45° chamfer at the union seam
// instead of a fillet. Useful when the design language is faceted rather
// than rounded.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	solid.Box(v3.XYZ(20, 20, 2), 1).
		UnderneathOf(solid.Cylinder(10, 4, 0).Bottom()).
		SmoothAdd(solid.ChamferMin(1.25)).
		STL("out.stl", 5.0)
}
