// Package reference: a 45 degree chamfer UNDER each tab catch. Each is a
// wedge filling the concave corner between the catch's underside (z = 6)
// and the block's +X face (x = 10): 3 mm tall where it meets the face,
// tapering to nothing at the catch's tip, so the whole underside is carried
// at 45 degrees and prints without supports on the block's -Z face. The
// catches themselves are untouched -- same 3 x 4 x 2 envelope, sharp top
// edge, full height at the tip.
//
// Incident: mcweed U69 -- "I can't print this without supports. the tab
// catches need to have a 45 degree chamfer on the bottom"; and mined4 I-5
// U523, "cuff under shoulder has no chamfer -> unprintable".
package reference

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

	tabL  = 3.0
	tabW  = 4.0
	tabH  = 2.0
	tabZ0 = 6.0
	tabY  = 6.0

	faceX = blockW / 2
)

func Build() *solid.Solid {
	block := solid.Box(v3.XYZ(blockW, blockD, blockH), 0).TranslateZ(blockH / 2)

	tab := solid.Box(v3.XYZ(tabL, tabW, tabH), 0).
		Translate(v3.XYZ(faceX+tabL/2, 0, tabZ0+tabH/2))

	// 45 degree support wedge, profile in the XZ plane: along the catch
	// underside from the face out to the tip, then straight back down the
	// face by the catch's own 3 mm reach. Extruded across the catch width.
	wedge := shape.Polygon([]v2.Vec{
		v2.XY(faceX, tabZ0),      // inside corner, at the block face
		v2.XY(faceX+tabL, tabZ0), // catch tip, wedge thickness 0
		v2.XY(faceX, tabZ0-tabL), // 45 degrees down the block face
	}).Extrude(tabW).RotateX(90)

	return block.Union(tab.Multi(v3.Y(tabY), v3.Y(-tabY)), wedge.Multi(v3.Y(tabY), v3.Y(-tabY)))
}
