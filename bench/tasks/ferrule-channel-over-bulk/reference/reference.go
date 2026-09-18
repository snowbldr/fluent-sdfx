// Package reference: exactly the groove the wire needs and no more. The
// ferrule is 2.4 across, so the groove is 2.4 wide and 2.4 deep, cut into
// the top face along y = 0, running from the hole out to the +X edge. The
// block keeps its full thickness everywhere else.
//
// Incident: mcweed U35 -- "we don't need that huge box of material right?
// Just a small box where the ferrule is". The model had cut a 4.5 x 14 lane
// where a 2.4 groove was all that was wanted.
package reference

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	blockW = 30.0
	blockD = 20.0
	blockH = 8.0

	holeR = 1.5

	grooveW = 2.4 // Y, the ferrule diameter
	grooveD = 2.4 // Z, measured down from the top face: floor at z = 5.6
	over    = 1.0 // cutter overshoot past the +X edge
)

func Build() *solid.Solid {
	block := solid.Box(v3.XYZ(blockW, blockD, blockH), 0).TranslateZ(blockH / 2)
	hole := solid.Cylinder(blockH+2, holeR, 0).TranslateZ(blockH / 2)

	// Groove: x from the hole (x = 0) out through the +X edge, y +-1.2,
	// z 5.6..8. Cut tool overshoots the edge and the top face.
	grooveLen := blockW/2 + over
	groove := solid.Box(v3.XYZ(grooveLen, grooveW, grooveD+over), 0).
		Translate(v3.XYZ(grooveLen/2, 0, blockH-grooveD+(grooveD+over)/2))

	return block.Cut(hole, groove)
}
