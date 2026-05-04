// Booleans: Cut subtracts the tool from the body (set difference).
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	solid.Box(v3.XYZ(20, 20, 20), 1).
		Cut(solid.Sphere(11).TranslateX(11)).
		STL("out.stl", 4.0)
}
