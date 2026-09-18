// Package start: a plain 30 x 20 x 10 block sitting on the build plate
// (x -15..15, y -10..10, z 0..10). No chamfer, no hole, no slot.
package start

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	blockX = 30.0
	blockY = 20.0
	blockZ = 10.0
)

func Build() *solid.Solid {
	return solid.Box(v3.XYZ(blockX, blockY, blockZ), 0).TranslateZ(blockZ / 2)
}
