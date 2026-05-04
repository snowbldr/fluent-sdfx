// Transforms: RotateX / Y / Z take an angle in degrees. RotateAxis lets
// you rotate around an arbitrary unit vector.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	bar := solid.Box(v3.XYZ(20, 4, 4), 0.5)

	parts := bar.Union(
		bar.RotateZ(60),
		bar.RotateZ(120),
	)
	parts.STL("out.stl", 5.0)
}
