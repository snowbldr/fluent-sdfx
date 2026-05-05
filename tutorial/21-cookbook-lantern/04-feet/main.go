// Lantern cookbook step 4: 4 small feet under the body.
//
// Cut(pocket, slots) carves the body, then OnTopOf(...).Union() raises
// the carved body onto a polar ring of feet. Every relation is anchor-
// named and relative — no bbox math, no absolute Z constants.
//
// 4 feet (rather than 3) because `Polar` with an even count is bbox-
// symmetric — the feet array's bbox top centre sits on the world Z axis,
// so the lantern lands centred on the feet.
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
)

func main() {
	body := solid.Cylinder(bodyHeight, bodyRadius, 4)
	pocket := solid.Cylinder(pocketDepth, bodyRadius-wallThick, 0)
	slot := solid.Box(v3.XYZ(slotWidth, slotWidth, slotHeight), 1)
	foot := solid.Cylinder(footHeight, footRadius, 0.8)

	body.Cut(
		pocket.Top().On(body.Top()).Solid(),
		slot.Top().Below(body.Top(), (pocketDepth-slotHeight)/2).Solid().
			Multi(layout.Polar(slotRadius, slotCount)...),
	).OnTopOf(foot.Multi(layout.Polar(footRing, 4)...).Top()).Union().STL("out.stl", 5.0)
}
