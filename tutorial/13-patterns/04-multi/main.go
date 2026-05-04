// Patterns: Multi places copies of a solid at an explicit list of
// positions. Variadic — drop the positions in directly, or pass a slice
// with `Multi(positions...)`.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	solid.Sphere(1.5).Multi(
		v3.XYZ(0, 0, 0),
		v3.XYZ(8, 0, 0),
		v3.XYZ(4, 7, 0),
		v3.XYZ(-4, 7, 0),
		v3.XYZ(-8, 0, 0),
		v3.XYZ(-4, -7, 0),
		v3.XYZ(4, -7, 0),
	).STL("out.stl", 8.0)
}
