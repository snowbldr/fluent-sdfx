// Package asfailed: what the model built first -- a 4.5 wide, 4 deep lane
// running most of the length of the block instead of a small groove where
// the wire actually is. It removes nearly six times the material the wire
// needs and thins the block under the lane to half its thickness.
//
// Incident: mcweed U35 -- "we don't need that huge box of material right?
// Just a small box where the ferrule is" (4.5 x 14 lane where a 2.4 mm
// groove sufficed).
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	blockW = 30.0
	blockD = 20.0
	blockH = 8.0

	holeR = 1.5

	laneW = 4.5  // OVERSIZED: nearly twice the ferrule
	laneD = 4.0  // OVERSIZED: half the block thickness
	laneX = 12.0 // OVERSIZED: the lane runs back past the hole, not out from it
	over  = 1.0
)

func Build() *solid.Solid {
	block := solid.Box(v3.XYZ(blockW, blockD, blockH), 0).TranslateZ(blockH / 2)
	hole := solid.Cylinder(blockH+2, holeR, 0).TranslateZ(blockH / 2)

	laneLen := laneX + blockW/2 + over // from x = -12 out past the +X edge
	lane := solid.Box(v3.XYZ(laneLen, laneW, laneD+over), 0).
		Translate(v3.XYZ(laneLen/2-laneX, 0, blockH-laneD+(laneD+over)/2))

	return block.Cut(hole, lane)
}
