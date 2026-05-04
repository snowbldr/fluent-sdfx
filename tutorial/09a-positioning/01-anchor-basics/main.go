// Positioning: every Solid has 27 named anchors on its bounding box —
// 6 face centers (Top, Bottom, Left, Right, Front, Back), 12 edge
// midpoints (TopRight, BottomFront, ...) and 8 corners
// (TopFrontRight, ...). Pair them with At/AtX/AtY/AtZ to land an
// anchor at a literal coordinate without doing bbox math by hand.
//
// Here BottomAt(0) drops the cylinder so its bottom face sits on z=0
// — equivalent to ZeroZ but reads naturally and composes with the
// rest of the positioning verbs.
package main

import "github.com/snowbldr/fluent-sdfx/solid"

func main() {
	solid.Cylinder(20, 10, 1).BottomAt(0).STL("out.stl", 5.0)
}
