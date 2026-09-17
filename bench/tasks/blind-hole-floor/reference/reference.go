// Package reference: the boss with a BLIND hole - a hole with a floor.
//
// "5.5 deep from the top of the boss, leaving a 2 floor" fixes both ends:
// the hole mouth is the boss top (z = 7.5) and its flat bottom is at
// z = 7.5 - 5.5 = 2.0, leaving 2 mm of solid boss under it and the plate
// below completely untouched. The cutter is padded ABOVE the boss top only;
// its lower end is exactly the hole floor and must not be padded.
//
// Incident vocabulary: mcweed's recurring "blind hole / pocket with a floor"
// vs "through hole". Cutters padded at both ends (the house style for
// through features, see llms.txt "always make cutting tools longer than the
// body") silently punch a blind hole straight through.
package reference

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plateX = 20.0
	plateY = 20.0
	plateT = 3.0 // z -3..0

	bossR = 5.0
	bossH = 7.5 // z 0..7.5

	holeR     = 1.0               // 2 mm diameter
	holeDepth = 5.5               // measured down from the boss top face
	floorT    = 2.0               // bossH - holeDepth: solid boss left under the hole
	holeZBot  = bossH - holeDepth // 2.0
	pad       = 1.0               // cutter overshoot ABOVE the boss top only
)

func Build() *solid.Solid {
	plate := solid.Box(v3.XYZ(plateX, plateY, plateT), 0).TranslateZ(-plateT / 2)
	boss := solid.Cylinder(bossH, bossR, 0).TranslateZ(bossH / 2)

	// Cutter spans z = holeZBot .. bossH+pad: open at the top, closed at
	// z = 2.0 so a floorT-thick floor survives.
	cutH := (bossH + pad) - holeZBot
	hole := solid.Cylinder(cutH, holeR, 0).TranslateZ(holeZBot + cutH/2)

	return plate.Union(boss).Cut(hole)
}
