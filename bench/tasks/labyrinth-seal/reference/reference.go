// Package reference: a tongue-and-groove labyrinth seal along the seam.
//
// Tongue on plate A (the lower half): 2 wide, 3 tall, the full 30 run of the
// seam, on the centre line. Groove in plate B (the upper half): the SAME
// tongue grown by the clearance — 0.15 per side and 0.15 at the bottom, so
// 2 + 2*0.15 = 2.3 wide and 3 + 0.15 = 3.15 deep. Without that clearance the
// two halves are an interference fit and never close on the seam, which is
// the whole point of the seal.
//
// Incident: mcweed I-10 (U473-U476) "the ends actually have a bit of a
// gap... wondered if we could add a nice chonky labyrinth seal there" —
// implemented as "tongue 2x3 with chamfer, groove +0.15 clearance".
package reference

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plateLen = 30.0
	plateW   = 20.0
	plateT   = 5.0
	layoutX  = 40.0

	tongueW = 2.0 // across the seam (Y)
	tongueH = 3.0 // above the seam plane (Z)

	clearance = 0.15 // per side, and at the bottom of the groove

	grooveW = tongueW + 2*clearance // 2.3
	grooveD = tongueH + clearance   // 3.15

	// bury/overrun keep boolean faces off the plate faces they meet, so no
	// tool face terminates flush inside material (mcweed U299, I-3).
	bury    = 0.5
	overrun = 1.0
)

func Build() *solid.Solid {
	plateA := solid.Box(v3.XYZ(plateLen, plateW, plateT), 0).Translate(v3.Z(-plateT / 2))
	plateB := solid.Box(v3.XYZ(plateLen, plateW, plateT), 0).Translate(v3.XZ(layoutX, plateT/2))

	// Tongue: stands z 0..tongueH, rooted bury deep inside plate A.
	tongue := solid.Box(v3.XYZ(plateLen, tongueW, tongueH+bury), 0).
		Translate(v3.Z((tongueH - bury) / 2))

	// Groove: cut z 0..grooveD up into plate B, the cutter running overrun
	// past both ends of the plate and overrun below the mating face.
	groove := solid.Box(v3.XYZ(plateLen+2*overrun, grooveW, grooveD+overrun), 0).
		Translate(v3.XZ(layoutX, (grooveD-overrun)/2))

	return plateA.Union(plateB, tongue).Cut(groove)
}
