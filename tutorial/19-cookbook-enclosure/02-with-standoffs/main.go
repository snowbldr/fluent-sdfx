// Enclosure cookbook step 2: panel with PCB standoffs at the corners.
//
// `layout.RectCorners(w, d)` returns the 4 corner positions of a rectangle
// centered on the origin — drop them into the variadic `.Multi(...)` and
// you get one standoff at each. `.OnTopOf(panel.Top())` flushes the
// standoffs to the panel's back face without bbox math.
package main

import (
	"github.com/snowbldr/fluent-sdfx/layout"
	"github.com/snowbldr/fluent-sdfx/obj"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

func main() {
	standoff := obj.Standoff3D(obj.StandoffParms{
		PillarHeight:   10,
		PillarDiameter: 5,
		HoleDepth:      8,
		HoleDiameter:   2.5,
		NumberWebs:     0,
	})

	panel := obj.Panel3D(obj.PanelParms{
		Size:         v2.XY(80, 50),
		CornerRadius: 4,
		Thickness:    3,
	})

	standoff.OnTopOf(panel.Top()).Solid().
		Multi(layout.RectCorners(60, 36)...).
		Union(panel).
		STL("out.stl", 6.0)
}
