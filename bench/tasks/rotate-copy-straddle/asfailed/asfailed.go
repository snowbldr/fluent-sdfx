// Package asfailed is the incident construction, in the form that still
// bites today: rotate the single port to its start angle FIRST, then ask for
// the rotational copies at that same angle.
//
// Incident: mcweed airway port profile — "rotating to -150 degrees BEFORE
// RotateCopy(3) put the cap circles across the -180 degree period boundary,
// dropping a chunk of each port."
//
// RotateCopyAtZ documents that the source is assumed to sit already in the
// canonical sector (centred on 0 degrees); it folds there and rotates the
// union afterwards. Handed a source that is already at -150 degrees, the
// folded domain (-60 .. +60 degrees) never reaches the port material at all,
// so the cutter evaluates as empty and ALL THREE ports silently vanish: the
// ring comes out solid, with no warning and no panic.
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/solid"
)

const (
	ringOuterR = 14.0
	ringInnerR = 9.0
	ringH      = 12.0

	portR     = 2.0
	portCount = 3
	portFirst = -150.0

	portInnerR = ringInnerR - 2
	portOuterR = ringOuterR + 2
	portLen    = portOuterR - portInnerR
	portMidR   = (portOuterR + portInnerR) / 2
)

func Build() *solid.Solid {
	ring := solid.Cylinder(ringH, ringOuterR, 0).
		Cut(solid.Cylinder(ringH+2, ringInnerR, 0))

	port := solid.Cylinder(portLen, portR, 0).RotateY(90).TranslateX(portMidR)

	// WRONG: aim first, pattern second. The source no longer lies in the
	// fold's canonical sector, so the pattern drops it.
	ports := port.RotateZ(portFirst).RotateCopyAtZ(portCount, portFirst)

	return ring.Cut(ports)
}
