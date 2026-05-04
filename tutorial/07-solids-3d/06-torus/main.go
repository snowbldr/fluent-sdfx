// 3D solids: a torus with a major radius (centre-to-tube distance) and a
// minor radius (tube radius).
package main

import "github.com/snowbldr/fluent-sdfx/solid"

func main() {
	solid.Torus(12, 3).STL("out.stl", 4.0)
}
