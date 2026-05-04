// 2D → 3D: partial revolution — a wedge of a full revolve.
//
// angleDeg is in degrees, measured from +X around the Y axis.
package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

func main() {
	shape.Rect(v2.XY(8, 12), 1.0).
		Translate(v2.X(8)).
		RevolveAngle(270).
		STL("out.stl", 5.0)
}
