// Patterns: Orient places copies of a solid pointed in a list of
// directions. The 'base' argument is the receiver's natural orientation.
//
// Useful for things like radial fins, antennae, or mounting tabs that need
// to point along arbitrary vectors.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	// A single fin pointing along +X by default, oriented to eight
	// directions in the XY plane.
	solid.Box(v3.XYZ(8, 1, 4), 0.2).TranslateX(8).
		Orient(v3.X(1), []v3.Vec{
			v3.XYZ(1, 0, 0),
			v3.XYZ(0, 1, 0),
			v3.XYZ(-1, 0, 0),
			v3.XYZ(0, -1, 0),
			v3.XYZ(1, 1, 0).Normalize(),
			v3.XYZ(-1, 1, 0).Normalize(),
			v3.XYZ(-1, -1, 0).Normalize(),
			v3.XYZ(1, -1, 0).Normalize(),
		}).STL("out.stl", 6.0)
}
