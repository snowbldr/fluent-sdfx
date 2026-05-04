// 2D → 3D: extrusion that twists the profile around Z over the height.
//
// twist is in radians. units.DtoR converts from degrees if you'd rather.
package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/units"
)

func main() {
	shape.Star(12, 5, 5).TwistExtrude(20, units.DtoR(90)).STL("out.stl", 5.0)
}
