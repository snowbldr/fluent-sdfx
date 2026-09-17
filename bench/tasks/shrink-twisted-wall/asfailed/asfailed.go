// Package asfailed: the inner void made with Shrink(1.0) on the twisted
// solid - the API trap.
//
// Shrink adds a constant to the SDF, so it moves the surface to wherever the
// SDF reads -1.0. On a TwistExtrude'd solid the SDF is deliberately scaled
// by invStretch = 1/sqrt(1 + k^2 * rMax^2) < 1 (a Lipschitz correction that
// stops the octree renderer from skipping cells). With k = pi/20 rad/mm and
// rMax^2 = 34, invStretch = 1/1.3560, so "1.0" reads -1.0 only 1.356 mm into
// the real material. Measured on this part:
//
//	requested wall  1.000 mm
//	actual wall     1.356 mm   (+35.6 percent, 0.356 mm too thick everywhere)
//	top cap         1.356 mm of material left over the bore, which the user
//	                asked to be OPEN at the top
//
// The bore is correspondingly undersized on every face, so whatever the plug
// was hollowed for no longer fits. Nothing in the code looks wrong: the
// constant 1.0 appears exactly once and matches the spec.
//
// Incident: mcweed feedback_fluent_api.md 2026-08-28 ("a nominal 0.37 Shrink
// removed >1mm of real geometry ... offsets on swept solids must be done in
// PROFILE space").
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	profX     = 10.0
	profY     = 6.0
	profRound = 1.0
	plugH     = 20.0
	twistDeg  = 180.0

	wallT  = 1.0
	floorT = 2.0
	zFloor = -plugH/2 + floorT
)

func Build() *solid.Solid {
	plug := shape.Rect(v2.XY(profX, profY), profRound).TwistExtrude(plugH, twistDeg)

	// WRONG: a scalar 3D inset of a twisted solid.
	bore := plug.Shrink(wallT).CutPlane(v3.Z(zFloor), v3.Z(1))

	return plug.Cut(bore)
}
