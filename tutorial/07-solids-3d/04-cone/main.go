// 3D solids: a truncated cone — bottom radius 10, top radius 4, height 18.
package main

import "github.com/snowbldr/fluent-sdfx/solid"

func main() {
	solid.Cone(18, 10, 4, 0.5).STL("out.stl", 4.0)
}
