// Package reference: a drafted square plug and the socket it wedges into,
// side by side in one solid (plug at x = -20, socket at x = +20).
//
// The pocket walls are the plug's drafted faces offset outward by 0.2
// measured PERPENDICULAR to the face, which is a horizontal offset of
// 0.2/cos(6 deg) = 0.2011 -- a distinction worth only 0.001 mm here.
//
// The part of the description that actually costs material is the pocket
// DEPTH. The joint is a wedge fit: the walls of the pocket are parallel to
// the plug's faces, so when the plug is pushed past flush by d, every face
// of the plug moves outward by d*tan(draft) and all four gaps close at
// once. The plug therefore does not stop at the nominal flush pose -- it
// keeps sinking until
//
//	d * tan(draft) = clearPerp / cos(draft)   =>   d = clearPerp / sin(draft)
//
// which at 6 deg and 0.2 mm of clearance is 1.913 mm, nearly ten times the
// clearance itself. The floor of the pocket must clear the nose by 0.5 mm
// at THAT pose, so the pocket is 12 + 1.913 + 0.5 = 14.413 deep, not 12.5.
//
// Incident: mcweed I-7 (U533-U541), the pin shoulder that jammed the
// fingers: a fit checked at its nominal pose and never at the pose the
// parts actually reach. "Motion-interference is only caught by sweeping
// the actual motion; static fit checks miss it."
package reference

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plugBase  = 16.0 // plug is 16 square at its base, z = 0
	draftDeg  = 6.0  // draft per side; the plug narrows going up
	plugH     = 12.0 // plug height
	clearPerp = 0.2  // clearance measured perpendicular to a drafted face
	noseGap   = 0.5  // clear air under the plug's nose at the jam pose

	blockW = 24.0 // socket block, square in plan
	blockH = 16.0 // socket block height, z = 0..16

	plugX   = -20.0 // plug centre
	socketX = 20.0  // socket centre

	overCut = 2.0 // pocket tool run past the block's top face
)

var (
	draftRad = draftDeg * math.Pi / 180
	tanDraft = math.Tan(draftRad)
	sinDraft = math.Sin(draftRad)
	cosDraft = math.Cos(draftRad)

	plugBaseHalf = plugBase / 2
	plugTopHalf  = plugBaseHalf - plugH*tanDraft // 6.7387

	// Perpendicular clearance on a face drafted draftDeg off vertical is a
	// horizontal offset of clearPerp/cos(draft), not clearPerp.
	clearHoriz = clearPerp / cosDraft // 0.20110

	// How far past flush the plug sinks before the tapered walls grab it.
	sink = clearPerp / sinDraft // 1.91335

	pocketDepth  = plugH + sink + noseGap // 14.4133, measured from the mouth
	pocketFloorZ = blockH - pocketDepth   // 1.5867

	// Pocket half-widths: mouth is the plug's base plus the horizontal
	// clearance; the pocket narrows going down at the same draft angle.
	pocketMouthHalf = plugBaseHalf + clearHoriz              // 8.2011
	pocketFloorHalf = pocketMouthHalf - pocketDepth*tanDraft // 6.6862
)

func Build() *solid.Solid {
	// Plug: exact square frustum, base on z = 0, narrowing going up.
	plug := solid.PyramidFrustum(plugBaseHalf, plugTopHalf, plugH).TranslateX(plugX)

	block := solid.Box(v3.XYZ(blockW, blockW, blockH), 0).BottomAt(0)

	// Pocket: the same frustum widening upward, run overCut past the top
	// face so the cut is not coincident with it. Built base-down then
	// flipped, since PyramidFrustum's first half-width is its widest.
	pocket := solid.PyramidFrustum(
		pocketMouthHalf+overCut*tanDraft, // half-width overCut above the mouth
		pocketFloorHalf,
		pocketDepth+overCut,
	).MirrorXY().BottomAt(pocketFloorZ)

	socket := block.Cut(pocket).TranslateX(socketX)

	return plug.Union(socket)
}
