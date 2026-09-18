// Package start: a solid twisted plug. A 10 x 6 rectangle with 1 mm rounded
// corners, twist-extruded 20 mm tall through 180 degrees, centred on the
// origin (z -10..10). Solid all the way through - nothing is hollowed yet.
package start

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

const (
	profX     = 10.0 // profile long axis
	profY     = 6.0  // profile short axis
	profRound = 1.0  // corner round
	plugH     = 20.0 // z -10..10
	twistDeg  = 180.0
)

func Profile() *shape.Shape {
	return shape.Rect(v2.XY(profX, profY), profRound)
}

func Build() *solid.Solid {
	return Profile().TwistExtrude(plugH, twistDeg)
}
