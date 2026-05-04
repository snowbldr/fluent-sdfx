// Quickstart step 1: a plain cylinder.
package main

import "github.com/snowbldr/fluent-sdfx/solid"

func main() {
	solid.Cylinder(20, 10, 0).STL("out.stl", 3.0)
}
