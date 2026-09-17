// Package start is the collar before the cam groove is cut: a plain tube,
// outer r 10, inner r 6, 20 tall, centred on the origin.
package start

import "github.com/snowbldr/fluent-sdfx/solid"

const (
	outerR = 10.0
	innerR = 6.0
	tubeH  = 20.0
)

func Build() *solid.Solid {
	// The bore cutter is 2 mm taller than the tube so its end faces exit the
	// body instead of terminating inside it.
	return solid.Cylinder(tubeH, outerR, 0).Cut(solid.Cylinder(tubeH+2, innerR, 0))
}
