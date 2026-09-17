// Package asfailed is the incident construction: the model wrote the
// variadic call webs.Intersect(cavity, band) expecting
// webs ∩ cavity ∩ band. In fluent-sdfx the variadic arguments are UNIONED
// first, so it actually means webs ∩ (cavity ∪ band).
//
// Incident: mcweed 2026-08-26 — "an agent used webs.Intersect(interiors,
// band) expecting webs ∩ interiors ∩ band and silently lost the cavity
// confinement". The damage here is the same two ways:
//   - web material at r < 8 survives OUTSIDE the band (spokes at |z| > 5);
//   - web material inside the band survives OUTSIDE the cavity, so fins
//     stick out past the tube's outer wall at r = 12 .. 14.
package asfailed

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

	bandH    = 10.0
	bandSide = 36.0
)

func Build() *solid.Solid {
	tube := solid.Cylinder(tubeH, tubeOuterR, 0).
		Cut(solid.Cylinder(tubeH+2, tubeInnerR, 0))

	web := solid.Box(v3.XYZ(webSpanR, webThick, tubeH), 0).TranslateX(webSpanR / 2)
	webs := web.RotateUnionZ(webCount, solid.RotateZMatrix(360.0/webCount))

	cavity := solid.Cylinder(tubeH+2, tubeInnerR, 0)
	band := solid.Box(v3.XYZ(bandSide, bandSide, bandH), 0)

	// WRONG: variadic Intersect means webs ∩ (cavity ∪ band).
	trimmed := webs.Intersect(cavity, band)

	return tube.Union(trimmed)
}
