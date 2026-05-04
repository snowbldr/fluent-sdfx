// Patterns: LineOf places copies of a solid along a line from p0 to p1.
//
// The pattern string is a single character per slot — 'x' places a copy,
// '.' skips it. Length of the string is the number of slots.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	solid.Sphere(2).
		LineOf(v3.XYZ(-15, 0, 0), v3.XYZ(15, 0, 0), "xxx.xxx").
		STL("out.stl", 6.0)
}
