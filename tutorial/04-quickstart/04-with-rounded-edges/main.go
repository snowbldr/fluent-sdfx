// Quickstart step 4: round the body's top and bottom edges by setting the round parameter.
package main

import (
	"github.com/snowbldr/fluent-sdfx/layout"
	"github.com/snowbldr/fluent-sdfx/solid"
)

func main() {
	solid.Cylinder(20, 10, 1).
		Cut(solid.Cylinder(25, 2, 0).Multi(layout.Polar(5, 4)...)).
		STL("out.stl", 3.0)
}
