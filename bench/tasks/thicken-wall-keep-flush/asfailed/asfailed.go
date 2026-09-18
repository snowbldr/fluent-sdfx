// Package asfailed: the model took the three asks as three independent
// edits. It thickened the wall correctly - outside still r = 20, bore now
// r = 14 - and then left the boss and the hole at the numbers they already
// had, because it read "keep the boss flush" and "keep it 6 deep" and
// "the hole stays 5" as "don't touch them".
//
// The result renders fine and every number the model reports is the number
// the user said: wall 6, boss 10 x 6 x 12, hole 5 through. It is still
// wrong. Both boss radii were dimensioned from the bore, and the bore
// moved:
//
//   - the boss's outer face is still at r = 16, so 2 mm of it is now buried
//     inside the thicker wall instead of flush with the bore;
//   - the boss only stands 4 mm proud of the new bore, not 6;
//   - the hole's inner end is 2 mm short of where the boss's inner face
//     should be.
//
// Incident: mined1 #6/#18 and the U6 punch list - coupled asks in one
// message, each done in isolation. The datum moved and nothing dimensioned
// from it followed.
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	outerR = 20.0
	wallT  = 6.0 // the one ask that was done
	boreR  = outerR - wallT
	height = 24.0

	bossOuterR = 16.0 // left where it was: now buried in the wall
	bossInnerR = 10.0 // left where it was: only 4 deep from the new bore
	bossW      = 10.0
	bossH      = 12.0
	bossZ      = 12.0

	holeD = 5.0
	holeZ = 12.0

	over = 2.0
)

func Build() *solid.Solid {
	housing := solid.Cylinder(height, outerR, 0).TranslateZ(height / 2)
	bore := solid.Cylinder(height+2*over, boreR, 0).TranslateZ(height / 2)

	bossThick := bossOuterR + over - bossInnerR
	boss := solid.Box(v3.XYZ(bossThick, bossW, bossH), 0).
		Translate(v3.XYZ(bossInnerR+bossThick/2, 0, bossZ))

	holeLen := (outerR + over) - (bossInnerR - over)
	hole := solid.Cylinder(holeLen, holeD/2, 0).
		RotateY(90).
		Translate(v3.XYZ(bossInnerR-over+holeLen/2, 0, holeZ))

	return housing.Cut(bore).Union(boss).Cut(hole)
}
