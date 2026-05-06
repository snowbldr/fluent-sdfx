// Gear cookbook step 4: the polished final gear — chamfered tooth tops,
// shaft bore with set-screw hole, and shrinkage compensation.
package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
)

const shrink = 1.0 / 0.999 // PLA ~0.1%

func main() {
	gear := obj.InvoluteGear(obj.InvoluteGearParms{
		NumberTeeth:      24,
		Module:           2,
		PressureAngleDeg: 20,
		Backlash:         0.05,
		Clearance:        0.2,
		RingWidth:        21,
		Facets:           16,
	}).ExtrudeRounded(4, 0.4) // soft top/bottom edges

	hub := solid.Cylinder(8, 8, 0.5)
	hub.Bottom().On(gear.Bottom()).Union().
		Cut(
			solid.Cylinder(20, 3, 0),                              // shaft bore
			solid.Cylinder(20, 1.25, 0).RotateY(90).TranslateZ(4), // radial set-screw
		).
		ScaleUniform(shrink).
		STL("out.stl", 6.0, 0.5)
}
