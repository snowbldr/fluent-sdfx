// Package start is the part before the change: a barrel with three internal
// ribs on its bore wall. The ribs run the full height and project 2 mm
// inward from the bore, at 0, 120 and 240 degrees. Nothing has been
// relieved yet.
package start

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	bodyOuterR = 14.0 // barrel outside radius
	boreR      = 8.0  // bore radius
	bodyH      = 20.0 // full height, z = -10..10

	ribDepth = 2.0 // radial: from the bore at r=8 inward to r=6
	ribW     = 3.0 // tangential
	ribCount = 3   // at 0, 120, 240 degrees
)

func Build() *solid.Solid {
	body := solid.Cylinder(bodyH, bodyOuterR, 0).
		Cut(solid.Cylinder(bodyH+2, boreR, 0))

	// One rib, spanning r = 6..8 on the +X side, full height.
	rib := solid.Box(v3.XYZ(ribDepth, ribW, bodyH), 0).
		TranslateX(boreR - ribDepth/2)

	return body.Union(rib.RotateUnionZ(ribCount, solid.RotateZMatrix(360/ribCount)))
}
