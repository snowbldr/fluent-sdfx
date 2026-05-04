// Lantern cookbook step 2: hollow out the tea light pocket.
//
// `pocket.Top().On(body.BottomAt(0).Top())` aligns the pocket's top
// anchor with the seated body's top anchor, returning a Placement;
// `.Cut()` is a Placement finalizer that subtracts the placed pocket
// from the body. Parts at the top, assembly at the bottom — one fluent
// expression, anchor-relational, with no bbox math.
package main

import "github.com/snowbldr/fluent-sdfx/solid"

const (
	bodyHeight  = 50.0
	bodyRadius  = 25.0
	wallThick   = 5.0
	pocketDepth = 40.0
)

func main() {
	body := solid.Cylinder(bodyHeight, bodyRadius, 4)
	pocket := solid.Cylinder(pocketDepth, bodyRadius-wallThick, 0)

	pocket.Top().On(body.BottomAt(0).Top()).Cut().STL("out.stl", 4.0)
}
