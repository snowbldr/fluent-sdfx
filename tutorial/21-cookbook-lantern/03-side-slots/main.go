// Lantern cookbook step 3: punch decorative slots through the wall.
//
// Body, pocket, and slot are bare primitives at the top. The assembly
// chain at the bottom does all the work: place the pocket on the seated
// body and Cut, then Cut the polar slot ring out of the result.
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
	slotZ      = bodyHeight - pocketDepth/2 // mid-pocket
)

func main() {
	body := solid.Cylinder(bodyHeight, bodyRadius, 4)
	pocket := solid.Cylinder(pocketDepth, bodyRadius-wallThick, 0)
	slot := solid.Box(v3.XYZ(slotWidth, slotWidth, slotHeight), 1)

	pocket.Top().On(body.BottomAt(0).Top()).Cut().
		Cut(slot.TranslateZ(slotZ).Multi(layout.Polar(slotRadius, slotCount)...)).
		STL("out.stl", 5.0)
}
