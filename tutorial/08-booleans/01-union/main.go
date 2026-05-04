// Booleans: Union joins two solids into one. Overlapping volume merges.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	solid.Box(v3.XYZ(16, 16, 16), 1).
		Union(solid.Sphere(11).TranslateX(8)).
		STL("out.stl", 4.0)
}
