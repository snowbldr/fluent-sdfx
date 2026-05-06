package obj_test

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v2i "github.com/snowbldr/fluent-sdfx/vec/v2i"
)

// A 1/4-inch UNC hex-head bolt, 2 inches long with a 0.5" smooth shank.
func ExampleBolt() {
	bolt := obj.Bolt(obj.BoltParms{
		Thread:      "unc_1/4",
		Style:       "hex",
		TotalLength: 2.0,
		ShankLength: 0.5,
	})
	// bolt.STL("bolt.stl", 4.0)
	_ = bolt
}

// A 4x4 Gridfinity base plate with magnet mounts and through-holes.
func ExampleGridfinityBase() {
	base := obj.GridfinityBase(obj.GridfinityBaseParms{
		Size:   v2i.XY(4, 4),
		Magnet: true,
		Hole:   true,
	})
	// base.STL("base_4x4.stl", 3.0)
	_ = base
}

// A 70x90 mm panel, 3mm thick, with rounded corners and 3.8mm holes on each edge.
func ExamplePanel3D() {
	panel := obj.Panel3D(obj.PanelParms{
		Size:         v2.XY(70, 90),
		CornerRadius: 5,
		HoleDiameter: 3.8,
		HoleMargin:   [4]float64{5, 5, 5, 5},
		HolePattern:  [4]string{"x", "x", "x", "x"},
		Thickness:    3,
	})
	// panel.STL("panel.stl", 5.0)
	_ = panel
}
