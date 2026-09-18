// Package start is the part before the change: a 30x30x4 plate with a plain
// cylindrical boss standing on its top face, centred.
package start

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plateW = 30.0 // X
	plateD = 30.0 // Y
	plateH = 4.0  // Z, plate spans z 0..4

	bossR = 4.0 // boss radius
	bossH = 6.0 // boss height, z 4..10
)

func Build() *solid.Solid {
	plate := solid.Box(v3.XYZ(plateW, plateD, plateH), 0).TranslateZ(plateH / 2)
	boss := solid.Cylinder(bossH, bossR, 0).TranslateZ(plateH + bossH/2)
	return plate.Union(boss)
}
