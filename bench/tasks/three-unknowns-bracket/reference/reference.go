// Package reference: the L bracket once the three things the prompt left out
// have been asked about and answered (see answers.md).
//
// Frame: +X across the bracket's 30 mm width, +Y away from the tube, +Z up
// the tube. The tube stands vertically at -Y; the mounting face is the plane
// y = 0 and the saddle plate occupies y 0..5. The outstanding leg projects
// from the bottom edge, out to y = 20, and is 5 mm thick in Z.
//
// The three answers that a model has to ask for:
//  1. the tube is 40 mm OD, so the mounting face is a concave cylindrical
//     saddle of radius 20, not flat;
//  2. the two M3 clearance holes are 3.4 mm and 24 mm apart, stacked on the
//     vertical centre line of the saddle face;
//  3. the strap slot is 20.5 x 3 with its LONG AXIS HORIZONTAL, across the
//     bracket, through the outstanding leg.
package reference

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plateW = 30.0 // X, across the tube
	plateH = 40.0 // Z, up the tube
	plateT = 5.0  // Y, stock thickness before the saddle is scooped out

	legW = 30.0 // X, same width as the plate
	legL = 20.0 // Y, how far the leg stands off the mounting face
	legT = 5.0  // Z, leg thickness; the leg occupies z 0..5

	// The saddle: the tube's own outside surface, cut out of the mounting
	// face. Axis vertical, parallel to the tube, tubeAxisY behind the face,
	// so the scoop is saddleDrop deep on the centre line and runs out to
	// nothing at x = +/- sqrt(tubeR^2 - tubeAxisY^2) = +/- 10.54.
	tubeR      = 20.0                  // 40 mm OD tube
	saddleDrop = 3.0                   // scoop depth on the centre line
	tubeAxisY  = -(tubeR - saddleDrop) // -17: the tube's axis, behind the face

	// Two M3 clearance holes, stacked up the centre line of the saddle face,
	// centred on the clear part of the face above the leg (z 5..40).
	boltR   = 1.7                    // 3.4 mm clearance for M3
	boltSep = 24.0                   // centre to centre, up the tube
	boltZ   = legT + (plateH-legT)/2 // 22.5, centre of the pair
	boltLen = plateT + 4             // overshoot both faces so the cut is clean

	// The strap slot, through the outstanding leg. Long axis HORIZONTAL:
	// 20.5 across the bracket for the 20 mm strap's width, 3 for its
	// thickness.
	slotL = 20.5     // X, across the bracket
	slotW = 3.0      // Y, front to back
	slotY = legL / 2 // 10: centred along the leg's projection
)

func Build() *solid.Solid {
	// Saddle plate: y 0..5, z 0..40, centred on x.
	plate := solid.Box(v3.XYZ(plateW, plateT, plateH), 0).
		Translate(v3.YZ(plateT/2, plateH/2))

	// Outstanding leg: y 0..20, z 0..5, centred on x.
	leg := solid.Box(v3.XYZ(legW, legL, legT), 0).
		Translate(v3.YZ(legL/2, legT/2))

	// The tube itself, used as the cutting tool: it scoops the mounting face
	// AND the root of the leg, because the tube bulges 3 mm past the plane of
	// the mounting face and both legs have to clear it.
	tube := solid.Cylinder(plateH*2, tubeR, 0).
		Translate(v3.YZ(tubeAxisY, plateH/2))

	// Bolt holes: along Y, through the plate, on the centre line, 24 apart.
	bolt := solid.Cylinder(boltLen, boltR, 0).RotateX(90).Translate(v3.Y(plateT / 2))
	boltUpper := bolt.Translate(v3.Z(boltZ + boltSep/2))
	boltLower := bolt.Translate(v3.Z(boltZ - boltSep/2))

	// Strap slot: through the leg's 5 mm thickness, long axis across X.
	slot := solid.Box(v3.XYZ(slotL, slotW, legT+4), 0).
		Translate(v3.YZ(slotY, legT/2))

	return plate.Union(leg).Cut(tube, boltUpper, boltLower, slot)
}
