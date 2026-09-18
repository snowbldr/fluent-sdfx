// Package asfailed: the pocket sized at the plug's nominal pose.
//
// Two mistakes, both of which survive a casual check:
//
//  1. The clearance is applied as a HORIZONTAL offset of 0.2 instead of a
//     horizontal offset of 0.2/cos(6 deg). On a 6 degree draft that is only
//     0.0011 mm of error -- invisible, and mentioned here only because it is
//     the error everybody looks for.
//
//  2. The pocket is made plugH + noseGap = 12.5 deep, i.e. the floor is put
//     0.5 mm under the nose with the plug sitting FLUSH in the pocket. But
//     nothing holds the plug flush: the walls are parallel to its faces, so
//     it keeps sinking until all four gaps close, which is
//     clearPerp/sin(6 deg) = 1.913 mm past flush. The nose lands on the
//     floor 1.4 mm before the taper grabs, so the part is located by the
//     floor instead of the wedge and the joint rocks. The pocket is 1.913
//     mm too shallow.
//
// Render it and it looks right: a drafted plug, a drafted pocket, clearance
// all round, air under the nose. Drop the plug in at the nominal pose and
// every static check passes. This is mcweed I-7 (U533-U541) again -- a fit
// verified at the pose it starts in, not the pose it reaches.
package asfailed

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plugBase  = 16.0
	draftDeg  = 6.0
	plugH     = 12.0
	clearPerp = 0.2
	noseGap   = 0.5

	blockW = 24.0
	blockH = 16.0

	plugX   = -20.0
	socketX = 20.0

	overCut = 2.0
)

var (
	draftRad = draftDeg * math.Pi / 180
	tanDraft = math.Tan(draftRad)

	plugBaseHalf = plugBase / 2
	plugTopHalf  = plugBaseHalf - plugH*tanDraft

	// WRONG: 0.2 sideways, not 0.2 perpendicular to the sloped face.
	clearHoriz = clearPerp

	// WRONG: no allowance for how far the wedge sinks past flush.
	pocketDepth  = plugH + noseGap // 12.5, should be 14.413
	pocketFloorZ = blockH - pocketDepth

	pocketMouthHalf = plugBaseHalf + clearHoriz
	pocketFloorHalf = pocketMouthHalf - pocketDepth*tanDraft
)

func Build() *solid.Solid {
	plug := solid.PyramidFrustum(plugBaseHalf, plugTopHalf, plugH).TranslateX(plugX)

	block := solid.Box(v3.XYZ(blockW, blockW, blockH), 0).BottomAt(0)

	pocket := solid.PyramidFrustum(
		pocketMouthHalf+overCut*tanDraft,
		pocketFloorHalf,
		pocketDepth+overCut,
	).MirrorXY().BottomAt(pocketFloorZ)

	socket := block.Cut(pocket).TranslateX(socketX)

	return plug.Union(socket)
}
