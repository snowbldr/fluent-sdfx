// Package asfailed: what a competent model builds when it does not ask.
//
// Incident: mcweed mined1 -- "Silent assumptions persisted across revisions
// because nothing forced them to be stated (print orientation x5 revs,
// thermistor package, standing-vs-lying, 'user visible')". Print orientation
// was finally stated at U18, after five geometry revisions; the 0.2 mm nozzle
// at U83, after four. Every one of those was a fact the user held and the
// model guessed at instead of asking for.
//
// This is a plausible, well-formed, perfectly printable L bracket. It is
// wrong in all three of the ways the prompt left open, and each one is a
// thing the user knew and would have said if asked:
//
//  1. FLAT mounting face. The tube is 40 mm OD, so a flat face rocks on it
//     and the two bolts cannot pull the bracket up tight. The reference has a
//     20 mm radius concave saddle scooped 3 mm into the face (and into the
//     root of the leg, which also has to clear the tube).
//  2. Bolt spacing guessed at 20 mm instead of the 24 mm the user's tube
//     clamp is drilled for, so neither hole lines up.
//  3. Slot turned 90 degrees: the model put the long axis front to back,
//     along the leg's projection, and shortened it to 18 mm so it would fit
//     in the 20 mm leg. A 20 mm wide strap will not pass through a 3 mm wide
//     opening.
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plateW = 30.0
	plateH = 40.0
	plateT = 5.0

	legW = 30.0
	legL = 20.0
	legT = 5.0

	boltR   = 1.7
	boltSep = 20.0 // GUESSED: the user's clamp is drilled 24 apart
	boltZ   = 22.5

	slotL = 3.0  // GUESSED orientation: 3 across, 18 front to back
	slotW = 18.0 // shortened from 20.5 so it fits the 20 mm leg
	slotY = 10.0
)

func Build() *solid.Solid {
	plate := solid.Box(v3.XYZ(plateW, plateT, plateH), 0).
		Translate(v3.YZ(plateT/2, plateH/2))
	leg := solid.Box(v3.XYZ(legW, legL, legT), 0).
		Translate(v3.YZ(legL/2, legT/2))

	// No saddle cut at all: the mounting face is left flat.

	bolt := solid.Cylinder(plateT+4, boltR, 0).RotateX(90).Translate(v3.Y(plateT / 2))
	boltUpper := bolt.Translate(v3.Z(boltZ + boltSep/2))
	boltLower := bolt.Translate(v3.Z(boltZ - boltSep/2))

	slot := solid.Box(v3.XYZ(slotL, slotW, legT+4), 0).
		Translate(v3.YZ(slotY, legT/2))

	return plate.Union(leg).Cut(boltUpper, boltLower, slot)
}
