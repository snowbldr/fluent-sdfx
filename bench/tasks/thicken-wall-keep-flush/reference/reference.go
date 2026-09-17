// Package reference: all three asks done together, which means the boss
// MOVES.
//
// The wall goes 4 -> 6 thick with the outside pinned at r = 20, so the bore
// goes 16 -> 14: the wall grows inward and every radius measured from the
// bore moves with it.
//
//   - "keep the boss flush against the bore" -> its outer face is now at
//     r = 14, not 16.
//   - "keep it 6 deep measured from the bore" -> its inner face is now at
//     r = 14 - 6 = 8, not 10. The boss spans r 8..14 and still stands 6 mm
//     proud of the bore wall, 10 wide and 12 tall at the same angle and
//     height.
//   - the 5 mm radial hole still runs from outside the wall clean through
//     the wall and the boss, so its inner end follows the boss's new inner
//     face out to r = 8.
//
// Incident: mined1 #6/#18 (mcweed U15/U39) and the U6 18-item punch list -
// the user sends several asks in one message and they are coupled; the
// model has to notice that changing the datum (the bore) moves everything
// dimensioned from it.
package reference

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	outerR = 20.0 // unchanged: the wall grows inward
	wallT  = 6.0  // was 4
	boreR  = outerR - wallT
	height = 24.0

	bossOuterR = 14.0 // flush against the new bore
	bossInnerR = 8.0  // still 6 deep measured from the bore
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
