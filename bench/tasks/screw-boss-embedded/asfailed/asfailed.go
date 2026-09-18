// Package asfailed is the "knuckle tower" the model actually built: the boss
// parked ON the tube's outer surface (axis at x = outerR) and unioned on raw,
// with no intersection against the parent cylinder.
//
// Result: a bump standing up to 2.2 mm proud of the tube at mid length and
// 4.5 mm proud at its ends, and no added material on the inside at all — the
// screw bore is drilled through fresh air outside the wall rather than
// through a thickened wall. "what are those bumps on the outsides of the
// cylinders?... those should be embedded in the tubes", "sticking out in a
// not nice way so they need to be inset further", "still popping out of the
// cylinder in a not nice way" (mcweed U19, U43, U46).
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	outerR = 10.0
	innerR = 8.0
	tubeH  = 30.0

	bossR   = 2.2
	bossX   = 10.0 // WRONG: sat on the outer surface instead of sunk into the wall
	bossZ   = 12.0
	bossLen = 12.0
	boreR   = 1.0
	boreLen = 16.0
)

func Build() *solid.Solid {
	tube := solid.Cylinder(tubeH, outerR, 0).Cut(solid.Cylinder(tubeH+2, innerR, 0))

	// No .Intersect(envelope): nothing trims the boss to the tube.
	boss := solid.Cylinder(bossLen, bossR, 0).RotateX(90).
		Translate(v3.XZ(bossX, bossZ))

	bore := solid.Cylinder(boreLen, boreR, 0).RotateX(90).
		Translate(v3.XZ(bossX, bossZ))

	return tube.Union(boss).Cut(bore)
}
