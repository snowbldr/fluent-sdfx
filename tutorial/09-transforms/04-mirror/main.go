// Transforms: Mirror operations reflect a solid across an axis plane.
// Useful for making symmetric parts from one half.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	// One asymmetric half, then mirror it through the YZ plane.
	half := solid.Box(v3.XYZ(8, 14, 6), 0.5).TranslateX(5)
	half.Union(half.MirrorYZ()).STL("out.stl", 5.0)
}
