// Package start: one of the housing tubes before the screw boss is added.
// A plain vertical tube, outer r 10, inner r 8, 30 tall, centred on the
// origin.
package start

import "github.com/snowbldr/fluent-sdfx/solid"

const (
	outerR = 10.0
	innerR = 8.0
	tubeH  = 30.0
)

func Build() *solid.Solid {
	return solid.Cylinder(tubeH, outerR, 0).Cut(solid.Cylinder(tubeH+2, innerR, 0))
}
