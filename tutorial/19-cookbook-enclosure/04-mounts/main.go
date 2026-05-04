// Enclosure cookbook step 4: add internal screw mounts to a single shell.
//
// Build the shell as outer-box minus an inner cavity. Anchor verbs
// (`BottomAt`, `OnTopOf`) keep the cavity flush with the inner floor and
// drop the standoffs into place without manual Z math; `layout.RectCorners`
// gives the 4 corner positions for the standoff array.
package main

import (
	"github.com/snowbldr/fluent-sdfx/layout"
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	const w, h, d = 80.0, 50.0, 30.0
	const wall = 2.0

	shell := solid.Box(v3.XYZ(w, h, d), 3)
	// Cavity sits on the inner floor and pokes 1mm above the top —
	// that opens the top cleanly without leaving a thin lip.
	cavity := solid.Box(v3.XYZ(w-2*wall, h-2*wall, d-wall+1), 2).
		BottomAt(shell.Bottom().Point.Z + wall)

	standoff := obj.Standoff3D(obj.StandoffParms{
		PillarHeight:   d - 2*wall,
		PillarDiameter: 6,
		HoleDepth:      d - 2*wall - 2,
		HoleDiameter:   2.5,
		NumberWebs:     4,
		WebHeight:      6,
		WebDiameter:    10,
		WebWidth:       1,
	})

	shell.Cut(cavity).
		Union(standoff.Multi(layout.RectCorners(w-16, h-16)...)).
		STL("out.stl", 4.0)
}
