// Enclosure cookbook step 5: the full enclosure with a tray shell, four
// PCB standoffs, and a separate screw-on lid.
//
// `body` and `lid` are kept named because they're the two conceptual
// pieces of the assembly — the final Union lays them out exploded along Y
// using `BehindOf` to put the lid one panel-width plus 20mm behind the
// body without measuring anything by hand.
package main

import (
	"github.com/snowbldr/fluent-sdfx/layout"
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	const w, h, d = 80.0, 50.0, 30.0
	const wall = 2.0
	const sx, sy = (w - 16) / 2, (h - 16) / 2

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

	shell := solid.Box(v3.XYZ(w, h, d), 3)
	cavity := solid.Box(v3.XYZ(w-2*wall, h-2*wall, d-wall+1), 2).
		BottomAt(shell.Bottom().Point.Z + wall)

	body := shell.Cut(cavity).
		Union(standoff.Multi(layout.RectCorners(w-16, h-16)...))

	// Lid: a panel sized to match the shell, with mounting holes that
	// align with the standoff screws.
	lid := obj.Panel3D(obj.PanelParms{
		Size:         v2.XY(w, h),
		CornerRadius: 3,
		HoleDiameter: 3.2,
		HoleMargin:   [4]float64{(h-2*sy)/2 - 0.001, (w-2*sx)/2 - 0.001, (h-2*sy)/2 - 0.001, (w-2*sx)/2 - 0.001},
		HolePattern:  [4]string{"x.x", "x", "x.x", "x"},
		Thickness:    wall,
	})

	lid.BehindOf(body.Back(), 20).Union().STL("out.stl", 4.0)
}
