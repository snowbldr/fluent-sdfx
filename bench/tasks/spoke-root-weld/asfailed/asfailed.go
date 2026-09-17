// Package asfailed: what the model actually built (mcweed I-12, U304-U308).
//
// The spoke's root face is placed exactly ON the hub radius, x = 8. That is
// correct only on the midplane: at y = 0 the hub surface is at x = 8, but at
// the slab edges y = +/-1.5 the hub surface has fallen back to
// x = sqrt(8^2 - 1.5^2) = 7.8581, so a wedge up to 0.142 mm deep and 10 mm
// tall is left OPEN at each root corner. The model verified the weld with a
// midplane cross-section, which passes, and declared the join sound
// ("That gap shouldn't be possible per the code").
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	hubR = 8.0
	hubH = 10.0

	spokeT   = 3.0
	spokeOut = 10.0
	spokeH   = 10.0

	spokeRootX = hubR // WRONG: the slab starts at the hub RADIUS, not past the chord
	spokeTipX  = hubR + spokeOut
	spokeLen   = spokeTipX - spokeRootX
)

func Build() *solid.Solid {
	hub := solid.Cylinder(hubH, hubR, 0).TranslateZ(hubH / 2)

	spoke := solid.Box(v3.XYZ(spokeLen, spokeT, spokeH), 0).
		Translate(v3.XYZ((spokeRootX+spokeTipX)/2, 0, spokeH/2))

	return hub.Union(spoke)
}
