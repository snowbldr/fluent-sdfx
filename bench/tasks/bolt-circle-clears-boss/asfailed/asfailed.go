// Package asfailed is the incident construction: the clearance was measured
// to the hole's CENTRE instead of to the hole's edge, so the bolt circle
// came out at bossOuterR + gap = 16.0 rather than 18.2.
//
// Why it is wrong: every hole now spans r = 13.8 .. 18.2, so its inboard
// edge crosses r = 14 and eats into the boss's footprint - the flange under
// the boss wall is broken through, the 2 mm clearance the user asked for is
// -0.2 mm, and there is no land left for the boss to stand on. The user's
// rule named a face and an edge ("from the hole's edge to the boss's outer
// face"); reading it as a centre distance loses exactly one hole radius.
//
// Incident: mcweed U66 "on the bowl mold board, the holes are a bit too big
// and are overlapping, also two sections don't have holes" - the hole
// pattern was placed by its centres without checking the hole's edge against
// the feature it had to clear.
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/layout"
	"github.com/snowbldr/fluent-sdfx/solid"
)

const (
	flangeR = 30.0
	flangeH = 6.0
	boreR   = 8.0

	bossInnerR = 10.0
	bossOuterR = 14.0
	bossH      = 10.0

	holeCount = 6
	holeR     = 2.2
	bossGap   = 2.0

	// WRONG: hole CENTRE held 2 mm off the boss instead of the hole's edge.
	boltCircleR = bossOuterR + bossGap // 16.0, should be 18.2
)

func Build() *solid.Solid {
	flange := solid.Cylinder(flangeH, flangeR, 0).ZeroZ()

	boss := solid.Cylinder(bossH, bossOuterR, 0).
		Cut(solid.Cylinder(bossH+2, bossInnerR, 0)).
		BottomAt(flangeH)

	bore := solid.Cylinder(flangeH+2, boreR, 0).BottomAt(-1)

	holes := solid.Cylinder(flangeH+2, holeR, 0).BottomAt(-1).
		Multi(layout.Polar(boltCircleR, holeCount)...)

	return flange.Union(boss).Cut(bore, holes)
}
