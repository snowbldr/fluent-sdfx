// Lantern cookbook step 3: punch decorative slots through the wall.
//
// One Cut on the body, two args: the positioned pocket and a polar ring of
// slots. Each slot is positioned relative to the body's top — no absolute
// Z constants, so the body stays centred at origin.
package main

import (
	"github.com/snowbldr/fluent-sdfx/layout"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	bodyHeight  = 50.0
	bodyRadius  = 25.0
	wallThick   = 5.0
	pocketDepth = 40.0

	slotCount  = 8
	slotRadius = bodyRadius - wallThick/2 // sit centred in the wall
	slotWidth  = 7.0
	slotHeight = 30.0
)

func main() {
	body := solid.Cylinder(bodyHeight, bodyRadius, 4)
	pocket := solid.Cylinder(pocketDepth, bodyRadius-wallThick, 0)
	slot := solid.Box(v3.XYZ(slotWidth, slotWidth, slotHeight), 1)

	body.Cut(
		pocket.Top().On(body.Top()).Solid(),
		slot.Top().Below(body.Top(), (pocketDepth-slotHeight)/2).Solid().
			Multi(layout.Polar(slotRadius, slotCount)...),
	).STL("out.stl", 5.0)
}
