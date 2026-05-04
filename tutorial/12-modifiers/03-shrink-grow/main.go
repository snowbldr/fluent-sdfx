// Modifiers: Shrink/Grow uniformly inset or expand all surfaces by a
// scalar. A common use is press-fit clearance — grow the negative tool
// or shrink the positive part.
//
// Shown here as a side-by-side preview: post on the left, cleared hole
// on the right.
package main

import "github.com/snowbldr/fluent-sdfx/solid"

func main() {
	solid.Cylinder(15, 5, 0).TranslateX(-8).
		Union(solid.Cylinder(20, 5, 0).Grow(0.15).TranslateX(8)).
		STL("out.stl", 6.0)
}
