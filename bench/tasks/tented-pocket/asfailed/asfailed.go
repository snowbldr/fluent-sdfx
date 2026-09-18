// Package asfailed: the model cut a plain rectangular pocket with a FLAT
// ceiling at z=12. In the stated print orientation (on the block's -Z face)
// that ceiling is a 10 x 12 mm horizontal overhang -- exactly the thing the
// prompt asked to avoid; it sags or needs supports inside a pocket you
// cannot reach to clean out.
//
// Incident: mined4 I-3 (U561 ff) -- unsupported ceilings survived until the
// user opened the slicer's overhang view; "printable in the export
// orientation" is a constraint the user never restates.
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	blockW = 24.0
	blockD = 24.0
	blockH = 20.0

	pocketW   = 10.0
	pocketH   = 8.0
	pocketDep = 12.0
	pocketZ0  = 4.0
	over      = 2.0
)

func Build() *solid.Solid {
	block := solid.Box(v3.XYZ(blockW, blockD, blockH), 0).TranslateZ(blockH / 2)

	depth := pocketDep + over
	pocket := solid.Box(v3.XYZ(pocketW, depth, pocketH), 0).
		Translate(v3.XYZ(0, -depth/2, pocketZ0+pocketH/2))

	return block.Cut(pocket)
}
