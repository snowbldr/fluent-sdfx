package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
)

const wallThickness = 1.0

func carveinside(path string) *solid.Solid {
	// create the SDF from the mesh
	// WARNING: It will only work on non-intersecting closed-surface(s) meshes.
	// Pass negative value for inside: Shrink moves surface inward by wallThickness.
	return obj.ImportSTL(path, 20, 3, 5).Shrink(wallThickness)
}

func main() {
	carveinside("../../files/teapot.stl").ToSTL("inside-carved-out.stl", 300)
}
