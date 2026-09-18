// Package asfailed: the model read "drill a hole" as a THROUGH hole and
// padded the cutter at both ends, the house idiom for through features. The
// 2 mm floor is gone and the hole carries on through the 3 mm base plate,
// so the "closed end" the user asked for is open at both ends.
//
// This is the mcweed "grooves punched through the 1.25 floor" class of bug
// (U23: "My bug - grooves punched through the 1.25 floor"): a cutter sized
// for clean booleans rather than for the feature's actual depth.
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plateT = 3.0
	bossR  = 5.0
	bossH  = 7.5
	holeR  = 1.0
)

func Build() *solid.Solid {
	plate := solid.Box(v3.XYZ(20, 20, plateT), 0).TranslateZ(-plateT / 2)
	boss := solid.Cylinder(bossH, bossR, 0).TranslateZ(bossH / 2)

	// WRONG: spans z = -4 .. 8.5, straight through boss AND plate.
	hole := solid.Cylinder(12.5, holeR, 0).TranslateZ(2.25)

	return plate.Union(boss).Cut(hole)
}
