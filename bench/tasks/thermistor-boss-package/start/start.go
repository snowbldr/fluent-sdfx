// Package start: a plain ceramic core, a vertical cylinder, before the
// thermistor boss is added.
package start

import "github.com/snowbldr/fluent-sdfx/solid"

const (
	coreH = 30.0
	coreR = 10.0
)

func Build() *solid.Solid {
	return solid.Cylinder(coreH, coreR, 0)
}
