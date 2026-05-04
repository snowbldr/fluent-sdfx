// 3D solids: a rounded box via solid.Box(size, round).
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	solid.Box(v3.XYZ(20, 16, 10), 1.5).STL("out.stl", 4.0)
}
