// Lantern cookbook step 1: a rounded cylindrical body sitting flush on
// the build plate.
//
// `BottomAt(0)` lands the cylinder so its bottom face is exactly at z=0
// — no Translate / ZeroZ math, just a statement about where the bottom
// goes.
//
// The pattern across this cookbook: declare each part as a bare
// primitive at the top, then position and combine them in one fluent
// expression at the bottom.
package main

import "github.com/snowbldr/fluent-sdfx/solid"

const (
	bodyHeight = 50.0
	bodyRadius = 25.0
)

func main() {
	body := solid.Cylinder(bodyHeight, bodyRadius, 4)

	body.BottomAt(0).STL("out.stl", 4.0)
}
