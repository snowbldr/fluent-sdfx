// Package reference: the core with a half-round vertical boss on its +X
// side carrying an AXIAL glass-bead thermistor: two lead holes, one at each
// end of the bead, 5 mm apart along the boss, each running radially into
// the core. (Hidden spec: see answers.md.)
package reference

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	coreH = 30.0
	coreR = 10.0

	bossR   = 4.0 // half-round rib radius
	bossLen = 8.0 // along Z
	holeR   = 0.4 // 0.8 mm lead holes
	holeSep = 5.0 // centre to centre, along Z (axial bead: leads exit each end)
	holeLen = 8.0 // from boss surface into the core
)

func Build() *solid.Solid {
	core := solid.Cylinder(coreH, coreR, 0)
	boss := solid.Cylinder(bossLen, bossR, 0).Translate(v3.X(coreR))
	hole := solid.Cylinder(holeLen, holeR, 0).RotateY(90).Translate(v3.X(coreR + bossR - holeLen/2))
	holes := hole.Translate(v3.Z(holeSep / 2)).Union(hole.Translate(v3.Z(-holeSep / 2)))
	return core.Union(boss).Cut(holes)
}
