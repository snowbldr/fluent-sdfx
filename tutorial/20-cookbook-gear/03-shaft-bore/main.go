// Gear cookbook step 3: drill a shaft bore through gear and hub.
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
	hub.Bottom().On(gear.Bottom()).Union().
		Cut(solid.Cylinder(20, 3, 0)).
		STL("out.stl", 5.0)
}
