// Package start is the part before the change: a 30x20x8 block with a
// vertical through hole on the axis, where the ferrule sits. The wire has
// to leave the hole along the top face toward +X, but there is no channel
// for it yet.
package start

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	blockW = 30.0 // X
	blockD = 20.0 // Y
	blockH = 8.0  // Z, block spans z 0..8

	holeR = 1.5 // vertical through hole at the origin
)

func Build() *solid.Solid {
	block := solid.Box(v3.XYZ(blockW, blockD, blockH), 0).TranslateZ(blockH / 2)
	hole := solid.Cylinder(blockH+2, holeR, 0).TranslateZ(blockH / 2)
	return block.Cut(hole)
}
