// Package start is the part before the change: a 20x20x10 block sitting on
// z=0 with a cantilevered tab sticking out of its +X face at mid height.
//
// Nothing in the geometry says which way the part goes on the build plate.
// It actually prints standing on its -Y face, so the tab's unsupported
// "bottom" is its -Y side, not its -Z side (see answers.md).
package start

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	blockW = 20.0 // X
	blockD = 20.0 // Y
	blockH = 10.0 // Z, block spans z 0..10

	tabL  = 6.0 // X, how far the tab cantilevers off the +X face
	tabW  = 8.0 // Y
	tabH  = 3.0 // Z
	tabZ0 = 5.0 // tab spans z 5..8

	faceX = blockW / 2 // the +X face of the block, x = 10
)

func Build() *solid.Solid {
	block := solid.Box(v3.XYZ(blockW, blockD, blockH), 0).TranslateZ(blockH / 2)

	// Cantilever: x 10..16, y -4..4, z 5..8.
	tab := solid.Box(v3.XYZ(tabL, tabW, tabH), 0).
		Translate(v3.XYZ(faceX+tabL/2, 0, tabZ0+tabH/2))

	return block.Union(tab)
}
