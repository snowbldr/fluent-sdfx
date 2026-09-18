package solid

import (
	"fmt"
	"math"

	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// PyramidFrustum is an exact square frustum sitting on z=0: half-width
// halfBase at z=0 tapering (or widening) linearly to halfTop at z=h, made
// as a Box trimmed by four planes. Either half-width may be the larger.
// Panics if either half-width or h is not positive.
//
// Use this rather than a Loft of two squares. Loft lerps the two distance
// fields, and since a smaller square is not an offset of a larger one, the
// result's corners pull in by a few tenths of a millimetre — enough to
// matter against a 0.3 mm clearance. The plane-trimmed construction here
// is exact at every height.
func PyramidFrustum(halfBase, halfTop, h float64) *Solid {
	if halfBase <= 0 || halfTop <= 0 || h <= 0 {
		panic(fmt.Sprintf("solid.PyramidFrustum: halfBase %g, halfTop %g and h %g must all be positive", halfBase, halfTop, h))
	}
	// The blank must be as wide as the wider end, or a widening frustum is
	// clipped to a prism of the base width.
	wide := math.Max(halfBase, halfTop)
	b := Box(v3.XYZ(2*wide, 2*wide, h), 0).BottomAt(0)
	// Each plane passes through (halfBase, 0, 0) and (halfTop, 0, h); its
	// normal (-h, 0, -d) always has a negative outward component, so
	// CutPlane keeps the interior whichever end is wider.
	d := halfBase - halfTop
	return b.
		CutPlane(v3.XYZ(halfBase, 0, 0), v3.XYZ(-h, 0, -d)).
		CutPlane(v3.XYZ(-halfBase, 0, 0), v3.XYZ(h, 0, -d)).
		CutPlane(v3.XYZ(0, halfBase, 0), v3.XYZ(0, -h, -d)).
		CutPlane(v3.XYZ(0, -halfBase, 0), v3.XYZ(0, h, -d))
}
