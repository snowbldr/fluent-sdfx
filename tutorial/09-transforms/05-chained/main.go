// Transforms: every transform returns a new *Solid, so they chain.
// The whole part below is a single expression.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	solid.
		Box(v3.XYZ(20, 6, 6), 0.5).
		RotateZ(30).
		RotateX(15).
		TranslateZ(8).
		ScaleUniform(0.9).
		STL("out.stl", 5.0)
}
