// Booleans: Intersect keeps only the volume common to both solids.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	solid.Box(v3.XYZ(20, 20, 20), 1).
		Intersect(solid.Sphere(13)).
		STL("out.stl", 4.0)
}
