// Quickstart step 2: drill one hole down the middle with Cut.
package main

import "github.com/snowbldr/fluent-sdfx/solid"

func main() {
	solid.Cylinder(20, 10, 0).
		Cut(solid.Cylinder(25, 2, 0)).
		STL("out.stl", 3.0)
}
