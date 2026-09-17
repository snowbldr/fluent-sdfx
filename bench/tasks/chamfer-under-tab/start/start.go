// Package start is the part before the change: a 20x20x10 block on z=0 with
// two tab catches standing proud of its +X face. Each catch is 4 wide (Y),
// 3 long (X) and 2 tall (Z), at y = +-6, spanning z 6..8. Printed on the
// block's -Z face the catch undersides are flat horizontal overhangs.
package start

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	blockW = 20.0 // X
	blockD = 20.0 // Y
	blockH = 10.0 // Z, block spans z 0..10

	tabL  = 3.0 // X, how far each catch stands off the face
	tabW  = 4.0 // Y
	tabH  = 2.0 // Z
	tabZ0 = 6.0 // catches span z 6..8
	tabY  = 6.0 // catch centres, +-6

	faceX = blockW / 2 // the +X face of the block, x = 10
)

func Build() *solid.Solid {
	block := solid.Box(v3.XYZ(blockW, blockD, blockH), 0).TranslateZ(blockH / 2)

	tab := solid.Box(v3.XYZ(tabL, tabW, tabH), 0).
		Translate(v3.XYZ(faceX+tabL/2, 0, tabZ0+tabH/2))

	return block.Union(tab.Multi(v3.Y(tabY), v3.Y(-tabY)))
}
