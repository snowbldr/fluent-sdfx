// Bolt assembly cookbook step 6: the full assembly — bolt, washer, nut.
//
// Three concerns in this step:
//  1. Position each part at its correct Z coordinate.
//  2. Use TotalLength as the canonical reference for stack-up math.
//  3. Lay out the parts side-by-side in X for an exploded-view render.
package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
)

const (
	thread      = "M8x1.25"
	totalLength = 30.0
	shankLength = 6.0
)

func main() {
	bolt := obj.Bolt(obj.BoltParms{
		Thread:      thread,
		Style:       "hex",
		TotalLength: totalLength,
		ShankLength: shankLength,
	})

	washer := obj.Washer3D(obj.WasherParms{
		Thickness:   1.5,
		InnerRadius: 4.5,
		OuterRadius: 8.5,
	})

	nut := obj.Nut(obj.NutParms{
		Thread:    thread,
		Style:     "hex",
		Tolerance: 0.1,
	})

	// Exploded layout along X. An exploded view is one of the few cases
	// where literal `Translate` reads better than an anchor verb — each
	// part keeps its natural Z, only X is shifted.
	parts := solid.UnionAll(
		bolt.TranslateX(-20),
		washer,
		nut.TranslateX(20),
	)
	parts.STL("out.stl", 6.0)
}
