// Project setup: a minimal main.go that proves your fluent-sdfx install works.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	solid.Box(v3.XYZ(20, 20, 20), 2).
		Cut(solid.Sphere(11)).
		STL("out.stl", 3.0)
}
