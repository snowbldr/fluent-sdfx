// 2D shapes: a rounded rectangle.
package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

func main() {
	shape.Rect(v2.XY(30, 20), 3).Extrude(1).STL("out.stl", 5)
}
