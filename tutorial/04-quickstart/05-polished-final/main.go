// Quickstart step 5: print-shrinkage compensation and mesh decimation for a smaller STL.
package main

import (
	"github.com/snowbldr/fluent-sdfx/layout"
	"github.com/snowbldr/fluent-sdfx/solid"
)

const shrink = 1.0 / 0.999 // PLA shrinks ~0.1% on cooling.

func main() {
	solid.Cylinder(20, 10, 1).
		Cut(solid.Cylinder(25, 2, 0).Multi(layout.Polar(5, 4)...)).
		ScaleUniform(shrink).
		// 0.5 keeps half the triangles after meshoptimizer decimation.
		STL("out.stl", 3.0, 0.5)
}
