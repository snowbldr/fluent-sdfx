// 3D solids: a sphere via solid.Sphere(radius).
package main

import "github.com/snowbldr/fluent-sdfx/solid"

func main() {
	solid.Sphere(12).STL("out.stl", 4.0)
}
