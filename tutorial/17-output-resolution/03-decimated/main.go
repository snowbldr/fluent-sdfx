// Output: high-res render decimated to 25% triangle count. Visually
// similar to the full mesh but a fraction of the file size.
//
// The trailing argument is the keep-fraction passed through meshoptimizer.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	solid.Box(v3.XYZ(20, 20, 20), 4).
		Cut(solid.Sphere(11)).
		STL("out.stl", 8.0, 0.25)
}
