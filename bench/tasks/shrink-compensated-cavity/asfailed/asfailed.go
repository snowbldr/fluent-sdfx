// Package asfailed is the incident construction: the two slips that ride
// along with "compensate for 4% shrinkage".
//
//  1. The compensation was applied as "4% bigger" — a factor of 1.04 — when
//     the factor that undoes a 4% shrink is 1/0.96 = 1.0416666... The cavity
//     wall lands at r = 10.400 instead of 10.416667 and the cavity floor at
//     z = 5.680 instead of 5.666667. Every fired puck comes out 0.16% small:
//     0.033 mm on the diameter, which is exactly the class of error that
//     showed up at U278 as "caliper 88.0 vs CAD 87.16" and became the 1%
//     silicone-to-plaster oversize process rule.
//  2. The 1 mm chamfer was written into the cavity at its nominal 1 mm
//     instead of being grown with the rest of the puck, because it was
//     thought of as a finishing feature rather than as part geometry. The
//     chamfer is a face of the part like any other; it shrinks with the part,
//     so it has to grow with the cavity.
//
// Both errors are a few hundredths of a millimetre — under the scorer's grid
// cell, invisible to IoU and to max_dev. Only the intent probes, which
// evaluate the exact SDF at a point, separate this from the reference.
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	blockX = 34.0
	blockY = 34.0
	blockH = 14.0

	partR     = 10.0
	partH     = 8.0
	boreR     = 3.0
	boreDepth = 5.0
	chamfer   = 1.0 // WRONG: left at nominal, never grown

	grow = 1.04 // WRONG: 1 + shrink instead of 1 / (1 - shrink)
)

// Grown dimensions, written out the way they were written in the incident:
// the numbers that were thought about got the factor, the chamfer did not.
const (
	cavR     = partR * grow     // 10.400 (should be 10.416667)
	cavH     = partH * grow     // 8.320  (should be 8.333333)
	postR    = boreR * grow     // 3.120  (should be 3.125)
	postDeep = boreDepth * grow // 5.200  (should be 5.208333)
)

func cavityShape() *solid.Solid {
	return shape.Polygon([]v2.Vec{
		v2.XY(0, 0),
		v2.XY(cavR, 0),
		v2.XY(cavR, cavH-chamfer), // chamfer still 1.0 mm tall
		v2.XY(cavR-chamfer, cavH), // and still 1.0 mm deep radially
		v2.XY(postR, cavH),
		v2.XY(postR, cavH-postDeep),
		v2.XY(0, cavH-postDeep),
	}).Revolve()
}

func Build() *solid.Solid {
	block := solid.Box(v3.XYZ(blockX, blockY, blockH), 0).ZeroZ()

	// Orientation and placement are right; only the numbers are wrong.
	cavity := cavityShape().MirrorXY().TopAt(blockH)

	return block.Cut(cavity)
}
