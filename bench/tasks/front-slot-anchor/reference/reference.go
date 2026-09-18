// Package reference is the intended result: a 4 mm wide, 2 mm deep slot
// milled across the FRONT (low-Y) face of the block, running the full X
// width and centred on the face height.
//
// The anchor that matters: FrontAt(v) pins the solid's LOW-Y face at v
// (and BackAt(v) pins its HIGH-Y face). The cutter is made twice as deep
// as the slot and its low-Y face is pinned one slot-depth OUTSIDE the
// block, so exactly slotDeep of it lands inside the material.
package reference

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	blockW = 30.0 // X
	blockD = 20.0 // Y: front face y = -10, back face y = +10
	blockH = 10.0 // Z

	slotWide = 4.0 // across the face, in Z, centred on the face (z = -2 .. +2)
	slotDeep = 2.0 // into the block from the front face (y = -10 .. -8)
	slotOver = 2.0 // cutter overshoot past each end of the X run
)

func Build() *solid.Solid {
	block := solid.Box(v3.XYZ(blockW, blockD, blockH), 0)

	// Cutter: longer than the block in X so the slot runs right out of both
	// ends; 2 x slotDeep thick in Y so half of it sits outside the block.
	// FrontAt pins the LOW-Y face at y = -12, so the cutter occupies
	// y = -12 .. -8 and removes y = -10 .. -8: a 2 mm deep slot in the
	// front face.
	cutter := solid.Box(v3.XYZ(blockW+2*slotOver, 2*slotDeep, slotWide), 0).
		FrontAt(-blockD/2 - slotDeep)

	return block.Cut(cutter)
}
