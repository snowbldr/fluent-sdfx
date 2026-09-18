// Package reference: the left boss has moved 4 mm further out, from x = -10
// to x = -14, and its through hole has moved with it. That is the whole
// change. Everything else -- the right boss and its hole, the 2 mm chamfer on
// all four top edges of the plate, the 6 x 20 x 2 pocket, the plate outline
// and thickness -- is bit-for-bit what it was.
//
// Incident: mcweed I-2 (U217-U225) and mined4: the model kept adding things
// nobody asked for -- a seal collar, a lead-in taper, a boss pad -- until
// "no taper ffs", and "features the model invented unasked were the ones that
// failed". The other half of the same failure is collateral damage: a
// one-number edit done by rebuilding the part, which quietly moves or
// re-rounds something that was already validated on a printed part.
package reference

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plateW  = 44.0
	plateD  = 30.0
	plateH  = 5.0
	chamfer = 2.0

	bossR = 5.0
	bossH = 8.0
	holeR = 1.5

	leftBossX  = -14.0 // CHANGED: was -10, moved 4 mm further out
	rightBossX = 10.0  // unchanged

	pocketW = 6.0
	pocketD = 20.0
	pocketZ = 2.0
)

func chamferedPlate() *solid.Solid {
	base := solid.Box(v3.XYZ(plateW, plateD, plateH-chamfer), 0).
		TranslateZ((plateH - chamfer) / 2)
	band := shape.Rect(v2.XY(plateW, plateD), 0).
		ScaleExtrude(chamfer, v2.XY((plateW-2*chamfer)/plateW, (plateD-2*chamfer)/plateD)).
		TranslateZ(plateH - chamfer/2)
	return base.Union(band)
}

func Build() *solid.Solid {
	plate := chamferedPlate()

	boss := solid.Cylinder(bossH, bossR, 0).TranslateZ(plateH + bossH/2)
	hole := solid.Cylinder(plateH+bossH+2, holeR, 0).TranslateZ((plateH + bossH) / 2)

	pocket := solid.Box(v3.XYZ(pocketW, pocketD, pocketZ*2), 0).TranslateZ(plateH)

	return plate.
		Union(boss.TranslateX(leftBossX), boss.TranslateX(rightBossX)).
		Cut(pocket, hole.TranslateX(leftBossX), hole.TranslateX(rightBossX))
}
