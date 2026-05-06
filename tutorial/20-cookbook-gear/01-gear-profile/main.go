// Gear cookbook step 1: a 2D involute gear profile.
//
// obj.InvoluteGear is the workhorse. The pitch diameter equals
// NumberTeeth × Module — so 24 teeth at module 2 gives a 48mm pitch
// circle. Module also sets the tooth size.
//
// RingWidth is the radial thickness of the gear's wall, measured inward
// from the root circle. Set it large (≥ root radius) for a fully solid
// gear; smaller values produce a hollow ring.
package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
)

func main() {
	obj.InvoluteGear(obj.InvoluteGearParms{
		NumberTeeth:      24,
		Module:           2,
		PressureAngleDeg: 20, // 20° is the modern default
		Backlash:         0.05,
		Clearance:        0.2,
		RingWidth:        21, // ≥ root radius → fully solid gear
		Facets:           16,
	}).Extrude(1).STL("out.stl", 5)
}
