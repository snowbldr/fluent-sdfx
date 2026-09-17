// Package start: the two housing tubes before the connecting section is
// added. A big battery tube and a smaller airway tube, side by side, with
// nothing between them.
package start

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	bigX   = -15.0 // big tube centre
	bigR   = 8.0
	smallX = 15.0 // small tube centre
	smallR = 5.0
	tubeH  = 20.0 // both, centred on z = 0
)

func Build() *solid.Solid {
	big := solid.Cylinder(tubeH, bigR, 0).Translate(v3.X(bigX))
	small := solid.Cylinder(tubeH, smallR, 0).Translate(v3.X(smallX))
	return big.Union(small)
}
