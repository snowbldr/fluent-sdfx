// Output: cellsPerMM=8 — final-quality render. Smooth curves, much
// larger STL and longer build time.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	solid.Box(v3.XYZ(20, 20, 20), 4).
		Cut(solid.Sphere(11)).
		STL("out.stl", 8.0)
}
