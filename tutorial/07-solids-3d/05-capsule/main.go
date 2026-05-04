// 3D solids: a capsule — a cylinder with hemispherical caps. The 'height'
// is the axial length between cap centres; total length = height + 2*radius.
package main

import "github.com/snowbldr/fluent-sdfx/solid"

func main() {
	solid.Capsule(20, 6).STL("out.stl", 4.0)
}
