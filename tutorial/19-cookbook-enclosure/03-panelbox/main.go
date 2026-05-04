// Enclosure cookbook step 3: switch to obj.PanelBox3D for a full
// panel-and-shell assembly.
//
// PanelBox3D returns a slice of 4 solids: [front-panel, body-shell,
// back-panel, ...screws]. We lay them out side-by-side for an exploded
// preview.
package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	parts := obj.PanelBox3D(obj.PanelBoxParms{
		Size:       v3.XYZ(80, 50, 30),
		Wall:       2,
		Panel:      3,
		Rounding:   3,
		FrontInset: 2,
		BackInset:  2,
		Clearance:  0.05,
		Hole:       3.2,
		SideTabs:   "T.B",
	})

	// Lay out each piece along Y for an exploded view.
	exploded := make([]*solid.Solid, len(parts))
	for i, p := range parts {
		exploded[i] = p.TranslateY(float64(i-len(parts)/2) * 60)
	}
	solid.UnionAll(exploded...).STL("out.stl", 4.0)
}
