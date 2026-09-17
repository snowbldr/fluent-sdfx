// Package asfailed is what the model actually built: a plain rectangular
// box spanning centre to centre, its width set by the SMALLER tube's
// diameter so that it at least does not overhang that one.
//
// It is not an angled piece, so it cannot be tangent to both. Where it meets
// the big tube its flat face is 3 mm inboard of the tangent plane, leaving a
// visible step / notch on each side — "the current look is bad", "extend all
// the way between both tubes smoothly, this looks bad" (mcweed U46).
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	bigX   = -15.0
	bigR   = 8.0
	smallX = 15.0
	smallR = 5.0
	tubeH  = 20.0
)

func Build() *solid.Solid {
	big := solid.Cylinder(tubeH, bigR, 0).Translate(v3.X(bigX))
	small := solid.Cylinder(tubeH, smallR, 0).Translate(v3.X(smallX))

	// Straight box, centre to centre, width = the smaller tube's diameter.
	slab := solid.Box(v3.XYZ(smallX-bigX, 2*smallR, tubeH), 0)

	return big.Union(small, slab)
}
