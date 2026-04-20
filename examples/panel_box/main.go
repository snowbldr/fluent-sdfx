package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	s := obj.PanelBox3D(obj.PanelBoxParms{
		Size:       v3.XYZ(50.0, 40.0, 60.0),
		Wall:       2.5,     // wall thickness
		Panel:      3.0,     // panel thickness
		Rounding:   5.0,     // outer corner rounding
		FrontInset: 2.0,     // inset for front panel
		BackInset:  2.0,     // inset for pack panel
		Hole:       3.4,     // #6 screw
		SideTabs:   "TbtbT", // tab pattern
	})
	s[0].STL("panel.stl", 3.0)
	s[1].STL("top.stl", 3.0)
	s[2].STL("bottom.stl", 3.0)
}
