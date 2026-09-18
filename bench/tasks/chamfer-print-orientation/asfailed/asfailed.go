// Package asfailed: the model never asked how the part prints, assumed the
// default "sitting on -Z" orientation, and hung the 45 degree chamfer under
// the tab's -Z face. In the real print orientation (standing on -Y) that
// wedge supports nothing and the tab's actual underside still needs
// supports.
//
// Incident: mcweed U18 -- "you have some funny bits on the housing, angled
// roofs where it's not appropriate. the two halves of the body will print
// lying down"; the model's own words: "an assumption I carried through five
// revisions without asking".
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	blockW = 20.0
	blockD = 20.0
	blockH = 10.0

	tabL  = 6.0
	tabW  = 8.0
	tabH  = 3.0
	tabZ0 = 5.0

	faceX = blockW / 2
)

func Build() *solid.Solid {
	block := solid.Box(v3.XYZ(blockW, blockD, blockH), 0).TranslateZ(blockH / 2)
	tab := solid.Box(v3.XYZ(tabL, tabW, tabH), 0).
		Translate(v3.XYZ(faceX+tabL/2, 0, tabZ0+tabH/2))

	// WRONG FACE: a 45 degree wedge in the XZ plane under the tab's -Z
	// face, running from the tab tip back to the wall and clipped where it
	// meets the block's bottom at z=0. Extruded across the tab's width.
	wedge := shape.Polygon([]v2.Vec{
		v2.XY(faceX, tabZ0),      // inside corner at the block wall
		v2.XY(faceX+tabL, tabZ0), // tab tip, wedge thickness 0
		v2.XY(faceX+1, 0),        // 45 degrees down, clipped at the plate
		v2.XY(faceX, 0),
	}).Extrude(tabW).RotateX(90)

	return block.Union(tab, wedge)
}
