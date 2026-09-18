// Package reference: the two tubes joined by a slab whose two flat side
// faces are the EXTERNAL TANGENT PLANES common to both cylinders. Because
// the tubes are different sizes the slab is a tapered wedge, not a box, and
// it meets each tube tangentially, so there is no step anywhere along the
// join — "flush", "smooth", "sleek".
//
// Geometry. For circles centred at (bigX, 0) r = bigR and (smallX, 0)
// r = smallR, the upper external tangent line { p : n.p = c } satisfies
// n.C1 = c - bigR and n.C2 = c - smallR, so n.(C2 - C1) = bigR - smallR and
//
//	nx = (bigR - smallR) / d,  ny = sqrt(1 - nx^2),  c = nx*bigX + bigR
//
// with d the centre distance. The tangent point on each circle is
// Ci + ri*n. The slab's end faces are the vertical chords through those
// tangent points, so each end is a chord of its own tube: the slab's corners
// land exactly on the cylinder and the rest of the end face is buried inside
// it. No sliver, no step, one welded body.
//
// Incident: mcweed U46 "the rectangular section needs to be flush with the
// tubes to create a sleek look... we'll need this to be an angled piece
// since we have two different size tubes".
package reference

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	bigX   = -15.0
	bigR   = 8.0
	smallX = 15.0
	smallR = 5.0
	tubeH  = 20.0
)

func Build() *solid.Solid {
	big := solid.Cylinder(tubeH, bigR, 0).Translate(v3.X(bigX))
	small := solid.Cylinder(tubeH, smallR, 0).Translate(v3.X(smallX))

	// Unit normal of the upper external tangent line.
	d := smallX - bigX
	nx := (bigR - smallR) / d
	ny := math.Sqrt(1 - nx*nx)

	// Tangent points: centre + radius * normal.
	bigT := v2.XY(bigX+bigR*nx, bigR*ny)
	smallT := v2.XY(smallX+smallR*nx, smallR*ny)

	slab := shape.Polygon([]v2.Vec{
		bigT,
		smallT,
		v2.XY(smallT.X, -smallT.Y),
		v2.XY(bigT.X, -bigT.Y),
	}).Extrude(tubeH)

	return big.Union(small, slab)
}
