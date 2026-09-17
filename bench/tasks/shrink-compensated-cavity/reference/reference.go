// Package reference is the intended mould block: the fired puck, grown for
// shrinkage, sunk mouth-down into the top face of a 34 x 34 x 14 block.
//
// The relationship the whole task turns on: the casting shrinks 4% LINEAR on
// firing, so the fired size is 0.96 x the cavity. To land on the nominal
// puck the cavity must be the puck DIVIDED by 0.96 — a factor of
// 1/0.96 = 1.0416666..., not 1.04. On the 10 mm radius that is 10.416667,
// not 10.400; the 0.0167 mm between those two is far under any sane grid, so
// only a probe on the exact SDF can see it.
//
// Three more consequences fall out of "the cavity IS the grown puck":
//
//   - The growth is uniform, so the 1 mm rim chamfer grows too (1.041667).
//     A chamfer applied at its nominal 1 mm after the scale is wrong.
//   - The puck is scaled about its own axis and THEN positioned. Scaling a
//     positioned solid scales its offset from the world origin as well and
//     slides the cavity out of the block face.
//   - The cavity is the puck's negative, so the puck's blind BORE is a POST
//     of block material standing inside the cavity, and the post's radius and
//     height are scaled like everything else. The puck lies mouth-down, so
//     the post stands on the cavity floor (mouth-up would leave it hanging in
//     mid air off the open mouth, and would make the mouth narrower than the
//     cavity below it — an undercut that cannot demould).
//
// Incident: mcweed U64 mold jig rev 1 ("the mold negatives now have a little
// bottom on them... we need some kind of lip there so that it creates a
// better seal") and U278, where a caliper read 88.0 mm against a CAD 87.16
// and the silicone-to-plaster oversize became a process rule; plus U225's
// "taper those cavities so green ceramic teeth release" demould draft.
package reference

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	blockX = 34.0
	blockY = 34.0
	blockH = 14.0 // z = 0 .. 14

	// The puck at finished, fired size.
	partR     = 10.0 // 20 across
	partH     = 8.0
	boreR     = 3.0 // 6 mm hole
	boreDepth = 5.0 // into the top face
	chamfer   = 1.0 // 45 deg, round the top rim

	// 4% linear shrink on firing: fired = cavity * shrink, so
	// cavity = part / shrink.
	shrink = 0.96
	grow   = 1.0 / shrink // 1.0416666..., NOT 1.04
)

// puck is the fired part at nominal size: a meridian profile revolved about
// Z. In the profile +X is the radius and Y becomes world Z, so the puck sits
// with its flat bottom on z = 0 and its bored, chamfered face up at z = 8.
func puck() *solid.Solid {
	return shape.Polygon([]v2.Vec{
		v2.XY(0, 0),                   // on the axis, bottom face
		v2.XY(partR, 0),               // out to the rim
		v2.XY(partR, partH-chamfer),   // straight outer wall
		v2.XY(partR-chamfer, partH),   // 45 deg chamfer on the top rim
		v2.XY(boreR, partH),           // top face, in to the bore
		v2.XY(boreR, partH-boreDepth), // bore wall, 5 deep
		v2.XY(0, partH-boreDepth),     // bore floor, back to the axis
	}).Revolve()
}

func Build() *solid.Solid {
	block := solid.Box(v3.XYZ(blockX, blockY, blockH), 0).ZeroZ()

	// Grow about the puck's own axis BEFORE placing it, so the chamfer, the
	// bore and the height all grow together and nothing slides. Then flip it
	// mouth-down and drop its flat back flush with the top of the block: the
	// cavity floor lands at blockH - partH*grow = 5.666667 and the bore
	// stands up out of that floor as a post.
	cavity := puck().ScaleUniform(grow).MirrorXY().TopAt(blockH)

	return block.Cut(cavity)
}
