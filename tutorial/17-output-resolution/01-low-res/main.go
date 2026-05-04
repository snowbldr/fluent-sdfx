// Output: cellsPerMM=0.5 — fast preview render. Visible faceting on
// curved surfaces, but writes in milliseconds.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	solid.Box(v3.XYZ(20, 20, 20), 4).
		Cut(solid.Sphere(11)).
		STL("out.stl", 0.5)
}
