// Package reference is the intended result: the webs are kept only where
// they are inside the tube's interior cavity AND inside the z band, then
// unioned back onto the tube.
//
// The API point: a true multi-way intersection must be CHAINED,
// webs.Intersect(cavity).Intersect(band). The variadic form intersects
// with the UNION of its arguments (see asfailed).
package reference

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	tubeOuterR = 12.0
	tubeInnerR = 8.0
	tubeH      = 20.0

	webThick = 1.0
	webSpanR = 14.0
	webCount = 4

	bandH    = 10.0 // z = -5 .. +5
	bandSide = 36.0 // wide enough to be a pure z filter
)

func Build() *solid.Solid {
	tube := solid.Cylinder(tubeH, tubeOuterR, 0).
		Cut(solid.Cylinder(tubeH+2, tubeInnerR, 0))

	web := solid.Box(v3.XYZ(webSpanR, webThick, tubeH), 0).TranslateX(webSpanR / 2)
	webs := web.RotateUnionZ(webCount, solid.RotateZMatrix(360.0/webCount))

	cavity := solid.Cylinder(tubeH+2, tubeInnerR, 0)        // the bore, r < 8
	band := solid.Box(v3.XYZ(bandSide, bandSide, bandH), 0) // z = -5 .. +5

	// Chained: webs ∩ cavity ∩ band. Only the web material that is both
	// inside the bore and inside the band survives.
	trimmed := webs.Intersect(cavity).Intersect(band)

	return tube.Union(trimmed)
}
