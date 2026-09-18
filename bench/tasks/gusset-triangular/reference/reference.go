// Package reference is what "make the spokes triangular" meant: a solid
// corner gusset (a corbel) under each spoke, filling the corner between the
// spoke's underside and the hub wall, running the spoke's full length so
// every layer of the spoke lands on something.
//
// The spokes' own rectangular cross-section is UNCHANGED. The gusset's
// sloped face is 45 degrees, from the spoke's outer tip (x = hubR+spokeOut,
// z = spokeZ0) down and in to the hub wall (x = hubR, z = 0), which is both
// self-supporting when printed hub-axis-vertical and lands the corbel's root
// on the build plate.
//
// Hidden spec: answers.md. Incident: mcweed I-12, U290 "I meant like a solid
// corner gusset when I said triangular".
package reference

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
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

// spokeLen runs from the hub AXIS so the box welds to the curved hub wall
// with no edge gap (mcweed I-12, U304-U308).
const spokeLen = hubR + spokeOut

// gussetRun is the horizontal reach of the 45 degree face, from the hub
// surface out to the spoke tip. spokeZ0 == gussetRun makes it exactly 45.
const gussetRun = spokeOut

func Build() *solid.Solid {
	hub := solid.Cylinder(hubH, hubR, 0).TranslateZ(hubH / 2)

	spoke := solid.Box(v3.XYZ(spokeLen, spokeW, spokeT), 0).
		Translate(v3.XYZ(spokeLen/2, 0, spokeZ0+spokeT/2))

	// Corbel profile in (x, z): flat top against the spoke underside, a
	// vertical root buried in the hub, and the 45 degree sloped face from
	// (hubR, 0) up to the spoke tip (hubR+gussetRun, spokeZ0).
	gusset := shape.Polygon([]v2.Vec{
		v2.XY(0, 0),
		v2.XY(hubR, 0),
		v2.XY(hubR+gussetRun, spokeZ0),
		v2.XY(0, spokeZ0),
	}).Extrude(spokeW).RotateX(90)

	arm := spoke.Union(gusset)
	return hub.Union(arm.RotateUnionZ(spokes, solid.RotateZMatrix(360.0/spokes)))
}
