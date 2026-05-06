// 2D → 3D: extrusion that twists the profile around Z over the height.
//
// twistDeg is in degrees, like every angle in the public API.
package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
)

func main() {
	shape.Star(12, 5, 5).TwistExtrude(20, 90).STL("out.stl", 5.0)
}
