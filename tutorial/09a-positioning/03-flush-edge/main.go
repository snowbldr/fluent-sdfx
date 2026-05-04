// Positioning: place a tab so its right face is flush with the host's
// right face, inset slightly, with its bottom on the floor.
//
// `tab.Right().On(box.Right()).Solid()` aligns the tab's right anchor
// onto the box's right anchor and returns the moved Solid (escaping
// the Placement when we don't want a boolean here). From there standard
// chained transforms — TranslateX for the inset, BottomAt for the floor.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	box := solid.Box(v3.XYZ(40, 30, 20), 1).BottomAt(0)
	tab := solid.Box(v3.XYZ(6, 6, 4), 0)

	tab.Right().On(box.Right()).Solid().
		TranslateX(-2).
		BottomAt(2).
		Union(box).STL("out.stl", 8.0)
}
