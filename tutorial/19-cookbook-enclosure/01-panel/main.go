// Enclosure cookbook step 1: a panel-mounted front face.
//
// obj.Panel3D handles the rounded corners and per-edge mounting holes.
// We'll add internal structure in the next steps.
package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

func main() {
	panel := obj.Panel3D(obj.PanelParms{
		Size:         v2.XY(80, 50),
		CornerRadius: 4,
		HoleDiameter: 3.2, // M3 clearance
		HoleMargin:   [4]float64{6, 6, 6, 6},
		HolePattern:  [4]string{"x.x", "x", "x.x", "x"},
		Thickness:    3,
	})

	panel.STL("out.stl", 6.0)
}
