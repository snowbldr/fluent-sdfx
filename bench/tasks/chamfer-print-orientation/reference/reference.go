// Package reference is the intended result once the print orientation is
// known: the part stands on its -Y face, so the tab's unsupported underside
// is its -Y face. The chamfer is a 45 degree wedge filling the corner
// between the tab's -Y face and the block's +X face, running the tab's full
// 3 mm height in Z. Nothing else changes.
//
// Incident: mcweed U18 -- the model carried a standing-print assumption
// through five revisions and hung 45 degree gables off the wrong faces;
// "the two halves of the body will print lying down" ended it in one turn.
package reference

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
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
	tabY0 = -tabW / 2  // the tab's -Y face, y = -4: "down" in print orientation
)

func Build() *solid.Solid {
	block := solid.Box(v3.XYZ(blockW, blockD, blockH), 0).TranslateZ(blockH / 2)

	// Cantilever: x 10..16, y -4..4, z 5..8.
	tab := solid.Box(v3.XYZ(tabL, tabW, tabH), 0).
		Translate(v3.XYZ(faceX+tabL/2, 0, tabZ0+tabH/2))

	// 45 degree support wedge in the XY plane, under the tab's -Y face:
	// full depth (tabL = 6) against the block wall, tapering to nothing at
	// the tab's tip. Its 6 mm leg lands exactly on the block's -Y face,
	// i.e. on the build plate. Extruded through the tab's 3 mm Z height.
	wedge := shape.Polygon([]v2.Vec{
		v2.XY(faceX, tabY0),      // inside corner, at the block wall
		v2.XY(faceX+tabL, tabY0), // tab tip, wedge thickness 0
		v2.XY(faceX, tabY0-tabL), // 45 degrees down the block wall
	}).Extrude(tabH).TranslateZ(tabZ0 + tabH/2)

	return block.Union(tab, wedge)
}
