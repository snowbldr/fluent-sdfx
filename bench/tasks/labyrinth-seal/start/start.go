// Package start: the two mould-ring plates that meet at a seam, before the
// labyrinth seal is added. Plate A is the lower half (its mating face is its
// TOP face, at z = 0); plate B is the upper half (its mating face is its
// BOTTOM face, also at z = 0).
//
// They are laid out side by side, B shifted +40 in X, so both halves are in
// one solid and can be scored together in a shared frame — the same way they
// sit on a build plate, not assembled.
package start

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plateLen = 30.0 // X, the run of the seam
	plateW   = 20.0 // Y
	plateT   = 5.0  // Z
	layoutX  = 40.0 // plate B's offset in X so the two sit side by side
)

func Build() *solid.Solid {
	plateA := solid.Box(v3.XYZ(plateLen, plateW, plateT), 0).Translate(v3.Z(-plateT / 2))
	plateB := solid.Box(v3.XYZ(plateLen, plateW, plateT), 0).Translate(v3.XZ(layoutX, plateT/2))
	return plateA.Union(plateB)
}
