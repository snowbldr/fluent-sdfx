// Package start is the part before the change: a plain 24x24x20 block
// sitting on z=0, printed on its -Z face. Nothing has been cut into it yet.
package start

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	blockW = 24.0 // X
	blockD = 24.0 // Y
	blockH = 20.0 // Z, block spans z 0..20
)

func Build() *solid.Solid {
	return solid.Box(v3.XYZ(blockW, blockD, blockH), 0).TranslateZ(blockH / 2)
}
