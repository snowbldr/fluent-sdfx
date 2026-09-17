// Package asfailed is what the model actually built when the user said
// "make the spokes triangular": it read "triangular" as the spoke's CROSS
// SECTION and gave each spoke a gabled (house-shaped) roof, then added no
// gusset at all.
//
// This is wrong twice over. The spoke underside is still a flat 10 mm
// cantilever at z = spokeZ0 with nothing under it, so it still needs
// supports — the whole point of the request. And it destroys the spoke's
// square top corners, which the user never asked to change.
//
// Incident: mcweed I-12, U290 "The key is not printable as is. I meant like
// a solid corner gusset when I said triangular." Six user corrections.
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	hubR = 6.0
	hubH = 20.0

	spokeW   = 3.0
	spokeT   = 3.0
	spokeOut = 10.0
	spokeZ0  = 10.0
	spokes   = 3
)

const spokeLen = hubR + spokeOut

func Build() *solid.Solid {
	hub := solid.Cylinder(hubH, hubR, 0).TranslateZ(hubH / 2)

	// Gabled cross-section in (y, z): base spokeW wide on the spoke
	// underside, apex spokeT above it on the centre line. Extruded along the
	// spoke, then rolled so profile-x -> world Y and profile-y -> world Z.
	gable := shape.Polygon([]v2.Vec{
		v2.XY(-spokeW/2, 0),
		v2.XY(spokeW/2, 0),
		v2.XY(0, spokeT),
	}).Extrude(spokeLen).RotateX(90).RotateZ(90).
		Translate(v3.XYZ(spokeLen/2, 0, spokeZ0))

	return hub.Union(gable.RotateUnionZ(spokes, solid.RotateZMatrix(360.0/spokes)))
}
