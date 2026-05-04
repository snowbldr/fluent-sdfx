// Modifiers: Elongate stretches a solid along the given axes by inserting
// a flat-walled "extrusion" between the two halves of the SDF. Great for
// turning a cylinder into a rounded slot, or a sphere into a pill.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	solid.Sphere(6).Elongate(v3.XYZ(10, 0, 0)).STL("out.stl", 6.0)
}
