// Modifiers: Offset moves every surface along its outward normal by the
// given distance. Positive offset grows the solid, negative shrinks.
//
// Offset is similar to Shrink/Grow but uses the SDF's normal-aligned
// offset (more accurate around curved features).
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	solid.Box(v3.XYZ(20, 20, 8), 0).Offset(2).STL("out.stl", 5.0)
}
