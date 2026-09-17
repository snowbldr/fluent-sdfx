// Package reference: registration keys that are ASYMMETRIC, so the halves
// only go together one way.
//
// Two nubs on A's top face at A-relative (+8, 0) and (-8, 0); two pockets in
// B's top face at the MIRRORED positions, B-relative (-8, 0) and (+8, 0),
// so that when B is flipped over onto A the +8 nub meets the -8 pocket.
//
// The asymmetry: the +8 nub is a 3 diameter ROUND post, the -8 nub is a 3x3
// SQUARE post, and each pocket matches its own partner (3.3 round, 3.3x3.3
// square) with 0.3 of clearance. Rotate the assembly 180 degrees about Z and
// the square nub lands over the round pocket, which cannot accept it, so the
// wrong orientation simply does not close. The shape difference is also
// visible across the room, which is what the user asked for.
//
// Incident: mcweed U482 ("the keys are on the wrong side ... the key needs to
// be asymmetric, visually hard to tell which way is which"). The vocabulary
// table from that mining run: "'asymmetric' key -> must only assemble one
// way -> round vs square".
package reference

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plateX = 24.0
	plateY = 24.0
	plateT = 4.0
	plateA = -20.0
	plateB = 20.0

	keyOff  = 8.0 // key offset from its plate's centre, along X
	nubD    = 3.0 // nub diameter / square across flats
	nubH    = 2.0 // nub height above the top face
	pockD   = 3.3 // pocket diameter / square across flats: 0.3 clearance
	pockDep = 2.2 // pocket depth below the top face: 0.2 clearance
	pad     = 0.5 // pocket cutter overshoot above the top face
)

func plate(cx float64) *solid.Solid {
	return solid.Box(v3.XYZ(plateX, plateY, plateT), 0).Translate(v3.XYZ(cx, 0, plateT/2))
}

func Build() *solid.Solid {
	// Half A: a ROUND nub at A+8 and a SQUARE nub at A-8.
	roundNub := solid.Cylinder(nubH, nubD/2, 0).
		Translate(v3.XYZ(plateA+keyOff, 0, plateT+nubH/2))
	squareNub := solid.Box(v3.XYZ(nubD, nubD, nubH), 0).
		Translate(v3.XYZ(plateA-keyOff, 0, plateT+nubH/2))

	// Half B: pockets at the mirrored offsets, each matching ITS partner.
	// A's +8 round nub mates with B's -8 pocket, so that one is round.
	cutH := pockDep + pad
	roundPock := solid.Cylinder(cutH, pockD/2, 0).
		Translate(v3.XYZ(plateB-keyOff, 0, plateT-pockDep+cutH/2))
	squarePock := shape.Rect(v2.XY(pockD, pockD), 0).Extrude(cutH).
		Translate(v3.XYZ(plateB+keyOff, 0, plateT-pockDep+cutH/2))

	return plate(plateA).Union(plate(plateB), roundNub, squareNub).
		Cut(roundPock, squarePock)
}
