// Smooth blends: SmoothArray repeats a solid in a grid where each
// neighbour is smooth-unioned to its peers, producing soft webs between
// what would otherwise be discrete copies.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	solid.Sphere(4).
		SmoothArray(4, 4, 1, v3.XYZ(7, 7, 0), solid.RoundMin(2.5)).
		STL("out.stl", 5.0)
}
