// Package asfailed: the model chamfered the wrong edge. It put a 45 degree
// cut on the TOP outer edge of each catch -- a feature that is already
// self-supporting -- and left the underside a flat unsupported overhang,
// which is the surface the user was pointing at. The part still cannot be
// printed without supports, and the catch has lost engagement height at its
// tip into the bargain.
//
// Incident: mcweed U69 -- "I can't print this without supports. the tab
// catches need to have a 45 degree chamfer on the bottom".
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

	// WRONG EDGE: a 45 degree cutter taking the top outer corner off the
	// catch. Profile in the XZ plane, extruded across the catch width plus
	// clearance.
	chamfer := shape.Polygon([]v2.Vec{
		v2.XY(faceX+1.5, tabZ0+tabH+0.5),
		v2.XY(faceX+4.0, tabZ0+tabH+0.5),
		v2.XY(faceX+4.0, tabZ0),
	}).Extrude(tabW + 2).RotateX(90)

	return block.Union(tab.Multi(v3.Y(tabY), v3.Y(-tabY))).
		Cut(chamfer.Multi(v3.Y(tabY), v3.Y(-tabY)))
}
