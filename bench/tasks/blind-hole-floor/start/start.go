// Package start: a 20x20x3 base plate (z -3..0) with a plain solid boss
// standing on it (radius 5, 7.5 tall, z 0..7.5). Nothing is drilled yet.
package start

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plateX = 20.0
	plateY = 20.0
	plateT = 3.0 // z -3..0

	bossR = 5.0
	bossH = 7.5 // z 0..7.5
)

func Build() *solid.Solid {
	plate := solid.Box(v3.XYZ(plateX, plateY, plateT), 0).TranslateZ(-plateT / 2)
	boss := solid.Cylinder(bossH, bossR, 0).TranslateZ(bossH / 2)
	return plate.Union(boss)
}
