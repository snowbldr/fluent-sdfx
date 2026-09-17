// Package asfailed: the boss moved, the hole did not.
//
// The single most likely real slip on this edit. The start code holds one
// constant, bossX, that places BOTH bosses and BOTH holes; anyone editing it
// has to split it, and the easy half-edit is to give the boss cylinder a new
// x and leave the hole cylinder on the old one. The result is a blind boss
// with no hole through it, sitting 4 mm from a 3 mm hole that now goes
// through bare plate. It looks right in a top view and is invisible in a
// bounding-box check.
//
// Incident: mcweed I-2 (U209-U225). Five revisions and two wasted prints on
// one small boss, because each pass changed something other than the thing
// the user named. mined4: "features the model invented unasked were the ones
// that failed" -- and the flip side, every untouched feature the model
// disturbed on the way past. The prompt here says "everything else stays
// exactly as it is", and the left boss's hole was part of the left boss.
package asfailed

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

	leftBossX  = -14.0 // the boss did move
	leftHoleX  = -10.0 // ... and its hole was left behind at the old x
	rightBossX = 10.0

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
		Cut(pocket, hole.TranslateX(leftHoleX), hole.TranslateX(rightBossX))
}
