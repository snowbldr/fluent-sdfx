// 2D → 3D: linear extrusion of a 2D profile.
package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

func main() {
	shape.Rect(v2.XY(20, 16), 2).Extrude(8).STL("out.stl", 4.0)
}
