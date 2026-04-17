package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	k := obj.PanelBoxParms{
		Size:       v3.XYZ(50.0, 40.0, 60.0),
		Wall:       2.5,     // wall thickness
		Panel:      3.0,     // panel thickness
		Rounding:   5.0,     // outer corner rounding
		FrontInset: 2.0,     // inset for front panel
		BackInset:  2.0,     // inset for pack panel
		Hole:       3.4,     // #6 screw
		SideTabs:   "TbtbT", // tab pattern
	}
	s := obj.PanelBox3D(k)
	s[0].ToSTL("panel.stl", 300)
	s[1].ToSTL("top.stl", 300)
	s[2].ToSTL("bottom.stl", 300)
}
