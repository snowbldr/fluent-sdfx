// Package asfailed: the model did items 1 and 2 of the three-item punch list
// and silently DROPPED item 3, the slot in the front face. The chamfer and
// the hole are exactly right, so everything the model reported on checks out
// - the failure is an omission, not an error.
//
// Incident: mcweed U25, "seems you also missed the feedback about the ribs in
// the socket"; U6's 18-item punch list, where single-clause items went
// missing. Long multi-item corrections lose items, usually the last one.
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	blockX  = 30.0
	blockY  = 20.0
	blockZ  = 10.0
	chamfer = 2.0
	holeD   = 3.0
	pad     = 1.0
)

func chamferTop(s *solid.Solid, n v3.Vec, faceHalf float64) *solid.Solid {
	pt := n.MulScalar(faceHalf).Add(v3.Z(blockZ - chamfer))
	return s.CutPlane(pt, n.Neg().Sub(v3.Z(1)).Normalize())
}

func Build() *solid.Solid {
	block := solid.Box(v3.XYZ(blockX, blockY, blockZ), 0).TranslateZ(blockZ / 2)

	block = chamferTop(block, v3.X(1), blockX/2)
	block = chamferTop(block, v3.X(-1), blockX/2)
	block = chamferTop(block, v3.Y(1), blockY/2)
	block = chamferTop(block, v3.Y(-1), blockY/2)

	hole := solid.Cylinder(blockZ+2*pad, holeD/2, 0).TranslateZ(blockZ / 2)

	// NOTE: the 10 x 4 x 1 deep slot in the front face is missing.
	return block.Cut(hole)
}
