package solid

import (
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// PyramidFrustum is an exact square frustum, half-width halfBase at z=0
// to halfTop at z=h (a Box trimmed by four planes).
//
// Use this rather than a Loft of two squares. Loft lerps the two distance
// fields, and since a smaller square is not an offset of a larger one, the
// result's corners pull in by a few tenths of a millimetre — enough to
// matter against a 0.3 mm clearance. The plane-trimmed construction here
// is exact at every height.
func PyramidFrustum(halfBase, halfTop, h float64) *Solid {
	b := Box(v3.XYZ(2*halfBase, 2*halfBase, h), 0).BottomAt(0)
	d := halfBase - halfTop
	return b.
		CutPlane(v3.XYZ(halfBase, 0, 0), v3.XYZ(-h, 0, -d)).
		CutPlane(v3.XYZ(-halfBase, 0, 0), v3.XYZ(h, 0, -d)).
		CutPlane(v3.XYZ(0, halfBase, 0), v3.XYZ(0, -h, -d)).
		CutPlane(v3.XYZ(0, -halfBase, 0), v3.XYZ(0, h, -d))
}
