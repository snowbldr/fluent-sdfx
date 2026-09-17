// Package reference: the boss is 0.4 mm shorter and that is the ONLY
// change. Straight walls, sharp top edge, sharp root -- exactly as it was,
// just 5.6 mm tall instead of 6.0.
//
// Incident: mcweed U222-U223 -- the model added a taper to a boss nobody
// asked for ("the taper... means that it's not going to fit into the mold",
// then "no taper ffs"). The user wants minimal deltas to a thing that
// works and strongly dislikes unrequested cleverness.
package reference

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plateW = 30.0
	plateD = 30.0
	plateH = 4.0

	bossR = 4.0
	bossH = 5.6 // was 6.0: 0.4 off the height, nothing else
)

func Build() *solid.Solid {
	plate := solid.Box(v3.XYZ(plateW, plateD, plateH), 0).TranslateZ(plateH / 2)
	boss := solid.Cylinder(bossH, bossR, 0).TranslateZ(plateH + bossH/2)
	return plate.Union(boss)
}
