// Lantern cookbook step 5: cap the lantern and add a finial knob.
//
// All six parts are bare primitives at the top. The single fluent
// expression at the bottom does the entire assembly: Cut the pocket and
// slot ring out of the body, raise onto the polar foot ring, place the
// cap on top, then sit the finial on the cap. Six anchor relations in
// one chain — no bbox math, no Z arithmetic.
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
	slotRadius = bodyRadius - wallThick/2
	slotWidth  = 7.0
	slotHeight = 30.0

	footRadius = 4.0
	footHeight = 4.0
	footRing   = 18.0

	capHeight    = 4.0
	finialRadius = 4.0
)

func main() {
	body := solid.Cylinder(bodyHeight, bodyRadius, 4)
	pocket := solid.Cylinder(pocketDepth, bodyRadius-wallThick, 0)
	slot := solid.Box(v3.XYZ(slotWidth, slotWidth, slotHeight), 1)
	foot := solid.Cylinder(footHeight, footRadius, 0.8)
	cap := solid.Cylinder(capHeight, bodyRadius, 1.5)
	finial := solid.Sphere(finialRadius)

	finial.Bottom().On(
		cap.OnTopOf(
			body.Cut(
				pocket.Top().On(body.Top()).Solid(),
				slot.Top().Below(body.Top(), (pocketDepth-slotHeight)/2).Solid().
					Multi(layout.Polar(slotRadius, slotCount)...),
			).OnTopOf(foot.Multi(layout.Polar(footRing, 4)...).Top()).Union().Top(),
		).Union().Top(),
	).Union().STL("out.stl", 5.0)
}
