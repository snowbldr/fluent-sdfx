// 3D solids: a cylinder via solid.Cylinder(height, radius, round).
package main

import "github.com/snowbldr/fluent-sdfx/solid"

func main() {
	solid.Cylinder(20, 10, 1).STL("out.stl", 3.0)
}
