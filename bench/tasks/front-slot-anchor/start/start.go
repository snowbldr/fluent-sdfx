// Package start is the part before the change: a plain rectangular block,
// 30 x 20 x 10 mm, centred on the origin. Its FRONT face (the low-Y face)
// is at y = -10; its back face is at y = +10.
package start

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	blockW = 30.0 // X
	blockD = 20.0 // Y: front face y = -10, back face y = +10
	blockH = 10.0 // Z
)

func Build() *solid.Solid {
	return solid.Box(v3.XYZ(blockW, blockD, blockH), 0)
}
