// Gear cookbook step 2: extrude the gear to thickness and add a centre hub.
//
// The hub strengthens the gear around the shaft and keeps it from sagging
// when printed in plastic. `hub.OnTopOf(gear.Bottom())` flushes the hub's
// bottom face to the gear's bottom face, so the hub fuses into the gear
// from below and protrudes above as a raised boss.
package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
)

func main() {
	gear := obj.InvoluteGear(obj.InvoluteGearParms{
		NumberTeeth:      24,
		Module:           2,
		PressureAngleDeg: 20,
		Backlash:         0.05,
		Clearance:        0.2,
		RingWidth:        21,
		Facets:           16,
	}).Extrude(4)

	hub := solid.Cylinder(8, 8, 0.5)
	hub.Bottom().On(gear.Bottom()).Union().STL("out.stl", 5.0)
}
