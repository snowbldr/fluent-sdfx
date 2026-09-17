// Package start is the turn key before the printability fix: a vertical hub
// cylinder with three radial spoke boxes at 120 degrees, each a plain
// rectangular section cantilevered off the hub wall.
//
// Printed hub-axis-vertical, every spoke is a 10 mm unsupported cantilever
// with a flat underside at z = 10 — the defect the user is complaining about
// (mcweed I-12, U290).
package start

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	hubR = 6.0  // hub (centre pole) radius
	hubH = 20.0 // hub height, z 0..20

	spokeW   = 3.0  // tangential width, along Y for the 0 degree spoke
	spokeT   = 3.0  // vertical thickness (Z)
	spokeOut = 10.0 // reach beyond the hub surface
	spokeZ0  = 10.0 // spoke underside, z 10..13
	spokes   = 3
)

// spokeLen runs from the hub AXIS, not from the hub surface, so the box is
// buried in the hub and welds to the curved wall with no edge gap. A flat
// slab root face stopped at x = hubR only touches the cylinder on the mid
// plane and floats proud at the corners (mcweed I-12, U304-U308).
const spokeLen = hubR + spokeOut

func Build() *solid.Solid {
	hub := solid.Cylinder(hubH, hubR, 0).TranslateZ(hubH / 2)

	spoke := solid.Box(v3.XYZ(spokeLen, spokeW, spokeT), 0).
		Translate(v3.XYZ(spokeLen/2, 0, spokeZ0+spokeT/2))

	return hub.Union(spoke.RotateUnionZ(spokes, solid.RotateZMatrix(360.0/spokes)))
}
