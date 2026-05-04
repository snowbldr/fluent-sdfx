// 2D → 3D: full revolution of a 2D profile around the Y axis.
//
// The profile lives in the XY plane; +X is the radius from the axis.
package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

func main() {
	// A simple cup-like profile: rectangle offset from the axis.
	shape.Rect(v2.XY(8, 12), 1.0).
		Translate(v2.X(8)).
		Revolve().
		STL("out.stl", 5.0)
}
