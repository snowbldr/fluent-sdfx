// Package reference: a ring wall with three radial port holes at 120 degree
// spacing, the first one at -150 degrees from +X (so the set sits at -150,
// -30 and +90).
//
// The order that makes rotational patterning safe: build the port on the
// canonical +X axis (0 degrees), PATTERN it there, and only then rotate the
// finished set to the wanted start angle. RotateCopyZ/RotateCopyAtZ fold the
// query domain into one sector of width 360/n centred on 0 degrees, so a
// source that has already been rotated out of that sector is at the mercy of
// the fold.
package reference

import (
	"github.com/snowbldr/fluent-sdfx/solid"
)

const (
	ringOuterR = 14.0
	ringInnerR = 9.0
	ringH      = 12.0 // z = -6 .. +6

	portR     = 2.0
	portCount = 3
	portFirst = -150.0 // degrees from +X

	portInnerR = ringInnerR - 2 // cutter starts inside the bore
	portOuterR = ringOuterR + 2 // and ends clear of the outer wall
	portLen    = portOuterR - portInnerR
	portMidR   = (portOuterR + portInnerR) / 2
)

func Build() *solid.Solid {
	ring := solid.Cylinder(ringH, ringOuterR, 0).
		Cut(solid.Cylinder(ringH+2, ringInnerR, 0))

	// One port cutter, radial, on the +X axis: canonical 0 degrees.
	port := solid.Cylinder(portLen, portR, 0).RotateY(90).TranslateX(portMidR)

	// Pattern FIRST while the source is canonical, THEN aim the whole set.
	ports := port.RotateCopyZ(portCount).RotateZ(portFirst)

	return ring.Cut(ports)
}
