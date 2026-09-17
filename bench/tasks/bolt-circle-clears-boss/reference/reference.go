// Package reference is the intended result: six M4 clearance holes through
// the flange, evenly spaced on a bolt circle whose radius is DERIVED from
// the clearance rule rather than given.
//
// The rule is stated edge-to-face: the hole's EDGE must stand at least 2 mm
// off the boss's outer face. The hole's edge is one hole radius inboard of
// the hole's centre, so the bolt circle is
//
//	bossOuterR + gap + holeR = 14 + 2 + 2.2 = 18.2
//
// Two plausible readings are both wrong. Measuring to the hole's CENTRE
// gives 16.0 and drives the holes back under the boss's footprint - the hole
// then spans r = 13.8 .. 18.2 and undercuts the boss wall. Using M4's
// nominal 4 mm instead of the 4.4 mm clearance diameter gives 18.0, which
// looks right, passes IoU and passes max_dev, and leaves the hole 1.8 mm off
// the boss instead of 2.
//
// The rim check is the second half of the ask and is what makes this a
// verification task: the bolt circle is pushed out only as far as the boss
// forces it, and then 18.2 + 2.2 + 3 = 23.4 has to fit inside r = 30 before
// the radius is committed to. It does, with room to spare, so the derived
// radius stands. Six holes on an 18.2 circle are 18.2 apart centre to
// centre against a 4.4 diameter, so they cannot overlap either.
//
// Incident: mcweed U66 "on the bowl mold board, the holes are a bit too big
// and are overlapping, also two sections don't have holes" - hole patterns
// laid out without checking the hole against its neighbours or against the
// feature it had to clear.
package reference

import (
	"github.com/snowbldr/fluent-sdfx/layout"
	"github.com/snowbldr/fluent-sdfx/solid"
)

const (
	flangeR = 30.0
	flangeH = 6.0 // z = 0 .. 6
	boreR   = 8.0

	bossInnerR = 10.0
	bossOuterR = 14.0
	bossH      = 10.0 // z = 6 .. 16

	holeCount = 6
	holeR     = 2.2 // M4 clearance, 4.4 diameter
	bossGap   = 2.0 // hole EDGE to boss outer face
	rimLand   = 3.0 // flange left outside the hole edge at the rim

	// Derived, not given: the nearest the hole's edge may come to the boss
	// is bossOuterR + bossGap, and the centre is holeR further out again.
	boltCircleR = bossOuterR + bossGap + holeR // 18.2

	// The check that has to pass before that radius is committed to.
	rimNeeded = boltCircleR + holeR + rimLand // 23.4 <= flangeR
)

func Build() *solid.Solid {
	flange := solid.Cylinder(flangeH, flangeR, 0).ZeroZ()

	boss := solid.Cylinder(bossH, bossOuterR, 0).
		Cut(solid.Cylinder(bossH+2, bossInnerR, 0)).
		BottomAt(flangeH)

	bore := solid.Cylinder(flangeH+2, boreR, 0).BottomAt(-1)

	// One cutter, taller than the flange so both ends break out, copied
	// round the bolt circle. layout.Polar puts the first copy on +X.
	holes := solid.Cylinder(flangeH+2, holeR, 0).BottomAt(-1).
		Multi(layout.Polar(boltCircleR, holeCount)...)

	return flange.Union(boss).Cut(bore, holes)
}
