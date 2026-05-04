// 2D → 3D: loft transitions between two 2D profiles over a height.
//
// LoftTo blends from the receiver shape (bottom) to a target shape (top).
// The two profiles can have very different geometry — here a square base
// transitions to a circular top.
package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

func main() {
	shape.Rect(v2.XY(20, 20), 0).
		LoftTo(shape.Circle(8), 18, 1).
		STL("out.stl", 5.0)
}
