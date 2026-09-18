// Package reference: all THREE items on the punch list done.
//
//  1. 2 mm chamfer on all four top edges. Each chamfer is a 45 degree plane
//     that meets the top face 2 mm in from the side and the side face 2 mm
//     down from the top, e.g. for the +X edge the plane x + z = 15 + 10 - 2.
//     Bottom edges stay square.
//  2. A 3 mm diameter hole straight through, vertical, on the block centre.
//  3. A 10 (X) x 4 (Z) slot, 1 mm deep, in the FRONT face - the low-Y face
//     at y = -10 - centred on that face, long axis along X.
//
// Incident: mcweed U25 ("seems you also missed the feedback about the ribs
// in the socket") and the U6 18-item punch list. Long multi-item messages
// lose items; the last item in the list is the one that goes missing.
package reference

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	blockX = 30.0 // x -15..15
	blockY = 20.0 // y -10..10
	blockZ = 10.0 // z 0..10

	chamfer = 2.0 // 45 degree, measured along each face

	holeD = 3.0
	pad   = 1.0 // through-cutter overshoot at each end

	slotLen = 10.0 // along X
	slotHt  = 4.0  // along Z
	slotDep = 1.0  // into the front face, along +Y
)

// chamferTop returns s with a 45 degree chamfer cut off one top edge. n is
// the outward horizontal direction of that edge (+X, -X, +Y or -Y).
func chamferTop(s *solid.Solid, n v3.Vec, faceHalf float64) *solid.Solid {
	// Keep the side of the plane whose normal points back into the block:
	// (-n - Z). The plane passes through the point chamfer mm below the top
	// on that side face.
	pt := n.MulScalar(faceHalf).Add(v3.Z(blockZ - chamfer))
	return s.CutPlane(pt, n.Neg().Sub(v3.Z(1)).Normalize())
}

func Build() *solid.Solid {
	block := solid.Box(v3.XYZ(blockX, blockY, blockZ), 0).TranslateZ(blockZ / 2)

	// 1. chamfer all four top edges
	block = chamferTop(block, v3.X(1), blockX/2)
	block = chamferTop(block, v3.X(-1), blockX/2)
	block = chamferTop(block, v3.Y(1), blockY/2)
	block = chamferTop(block, v3.Y(-1), blockY/2)

	// 2. vertical through hole on the centre
	hole := solid.Cylinder(blockZ+2*pad, holeD/2, 0).TranslateZ(blockZ / 2)

	// 3. slot in the front (low-Y) face, 1 deep, centred, long axis X
	slot := solid.Box(v3.XYZ(slotLen, slotDep+pad, slotHt), 0).
		Translate(v3.XYZ(0, -blockY/2-pad/2+slotDep/2, blockZ/2))

	return block.Cut(hole, slot)
}
