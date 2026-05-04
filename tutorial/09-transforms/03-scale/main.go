// Transforms: Scale takes a per-axis vector; ScaleUniform a single factor.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	// Stretch a sphere along Z and squish along Y.
	solid.Sphere(8).Scale(v3.XYZ(1, 0.5, 1.5)).STL("out.stl", 4.0)
}
