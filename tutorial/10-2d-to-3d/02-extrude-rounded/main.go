// 2D → 3D: extrusion with rounded top and bottom edges.
package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

func main() {
	shape.Rect(v2.XY(20, 16), 2).ExtrudeRounded(8, 1.5).STL("out.stl", 5.0)
}
