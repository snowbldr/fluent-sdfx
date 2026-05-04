// Patterns: Array repeats a solid on a regular grid.
//
// Array(numX, numY, numZ, step) — counts along each axis, with `step`
// being the distance between cells.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	solid.Sphere(2).Array(5, 4, 1, v3.XYZ(5, 5, 0)).STL("out.stl", 6.0)
}
