// Package start is the part before the change: a bare hub, a vertical
// cylinder of radius 8 standing on the build plate (z 0..10). The spoke
// has not been attached yet.
package start

import "github.com/snowbldr/fluent-sdfx/solid"

const (
	hubR = 8.0  // hub outer radius
	hubH = 10.0 // hub height, z 0..10
)

func Build() *solid.Solid {
	// solid.Cylinder is centred on the origin; lift it so z runs 0..10.
	return solid.Cylinder(hubH, hubR, 0).TranslateZ(hubH / 2)
}
