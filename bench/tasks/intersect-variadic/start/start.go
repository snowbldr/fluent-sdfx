// Package start is the part before the change: a tube with four thin radial
// webs running straight through it, from the axis all the way out past the
// outer wall and over the tube's full height. Nothing is trimmed yet.
package start

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	tubeOuterR = 12.0
	tubeInnerR = 8.0  // the interior cavity is r < 8
	tubeH      = 20.0 // z = -10 .. +10

	webThick = 1.0  // Y thickness of each web
	webSpanR = 14.0 // webs run from the axis out to r = 14, past the tube
	webCount = 4
)

func Build() *solid.Solid {
	tube := solid.Cylinder(tubeH, tubeOuterR, 0).
		Cut(solid.Cylinder(tubeH+2, tubeInnerR, 0))

	// One web on +X, then four of them at 90 degree spacing.
	web := solid.Box(v3.XYZ(webSpanR, webThick, tubeH), 0).TranslateX(webSpanR / 2)
	webs := web.RotateUnionZ(webCount, solid.RotateZMatrix(360.0/webCount))

	return tube.Union(webs)
}
