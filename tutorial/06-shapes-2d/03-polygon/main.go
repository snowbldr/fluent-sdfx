// 2D shapes: an arbitrary polygon from a list of vertices.
package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

func main() {
	shape.Polygon([]v2.Vec{
		v2.XY(-12, -10),
		v2.XY(12, -10),
		v2.XY(15, 0),
		v2.XY(8, 12),
		v2.XY(-8, 12),
		v2.XY(-15, 0),
	}).Extrude(1).STL("out.stl", 5)
}
