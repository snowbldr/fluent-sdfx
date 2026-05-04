// 2D → 3D: extrusion that scales the profile linearly over the height.
//
// scale is the (x, y) multiplier at the top of the extrusion. The bottom
// is the original size; the top is original × scale.
package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

func main() {
	shape.Rect(v2.XY(16, 16), 1).ScaleExtrude(20, v2.XY(0.4, 0.4)).STL("out.stl", 5.0)
}
