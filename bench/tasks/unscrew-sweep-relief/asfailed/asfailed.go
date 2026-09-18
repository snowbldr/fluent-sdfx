// Package asfailed: the relief cut as a flat band at the lug's STARTING
// height.
//
// The model checked the lug against the ribs where the lug is now -- 5 tall,
// centred on z=0 -- grew that by the 0.3 clearance and cut z = -2.8..2.8
// out of all three ribs. At the start pose the check is perfect: the lug
// sits in that band and clears every rib by 0.3. Render it, section it,
// probe it at t=0 and it all agrees.
//
// It is wrong in both directions as soon as the insert moves, because the
// rotation and the 6 mm rise are coupled:
//
//   - the rib at 120 deg is met a sixth of a turn out, 1 mm up;
//   - the rib at 240 deg half a turn out, 3 mm up;
//   - the rib at 0 deg five sixths of a turn out, 5 mm up.
//
// So the band leaves rib material standing exactly where the risen lug has
// to pass (the rib at 0 deg is untouched above z=2.8 and the lug arrives
// there at z=5), and it throws away rib material below the path that was
// never in the way -- the ribs are what holds the insert, and this is the
// half of them nearest the bottom.
//
// Incident: mcweed I-7 (U533-U541), the shoulder and then the wedge flare
// that jammed the fingers -- the same class of miss twice, resolved only
// when the pin was swept through two unscrew turns instead of checked in
// place.
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	bodyOuterR = 14.0
	boreR      = 8.0
	bodyH      = 20.0

	ribDepth = 2.0
	ribW     = 3.0
	ribCount = 3

	lugH      = 5.0 // the lug's height, centred on z = 0 at the start pose
	clearance = 0.3

	// WRONG: the band is placed at the lug's starting height and never
	// follows the 6 mm rise that comes with the turn.
	bandH = lugH + 2*clearance // 5.6, centred on z = 0
)

func Build() *solid.Solid {
	body := solid.Cylinder(bodyH, bodyOuterR, 0).
		Cut(solid.Cylinder(bodyH+2, boreR, 0))

	rib := solid.Box(v3.XYZ(ribDepth, ribW, bodyH), 0).
		TranslateX(boreR - ribDepth/2)
	ribs := rib.RotateUnionZ(ribCount, solid.RotateZMatrix(360/ribCount))

	// A flat band through every rib, full radial depth. Only the ribs are
	// cut, so this differs from the reference in nothing but where the
	// relief sits.
	band := solid.Cylinder(bandH, bodyOuterR+1, 0)

	return body.Union(ribs.Cut(band))
}
