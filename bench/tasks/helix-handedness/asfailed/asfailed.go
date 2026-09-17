// Package asfailed: the model got the handedness backwards and produced a
// LEFT-hand rib — the same core, the same pitch, the same 2 turns, the same
// starting phase at z = 0, but the rib winds the other way.
//
// Mirroring about the XZ plane (flipping Y) is exactly what a model does
// when it "fixes" a helix it thinks is the wrong way round, or what falls
// out of building the helix from a negated angle. Nothing about the part's
// bounding box, height, mass or turn count reveals it: only the position of
// the rib along the path does, which is what the probes check.
//
// Incident family: orientation/handedness assumptions that survive every
// sanity check short of looking at the part (mcweed helical airway ribs).
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
)

const (
	coreR = 8.0
	coreH = 24.0

	ribProfileR = 1.0
	ribTurns    = 2.0
	ribHeight   = 24.0
)

func Build() *solid.Solid {
	core := solid.Cylinder(coreH, coreR, 0)

	// WRONG: mirrored about XZ, so the angle RETREATS as z rises — a
	// left-hand thread.
	rib := shape.Circle(ribProfileR).SweepHelix(coreR, ribTurns, ribHeight, false).MirrorXZ()

	return core.Union(rib)
}
