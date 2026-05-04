// Cross-sections: plane.At lets you slice on any oriented plane, not just
// axis-aligned ones. Here a sphere with a torus through it, sliced
// diagonally.
package main

import (
	"github.com/snowbldr/fluent-sdfx/plane"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	shape.SliceAt(
		solid.Sphere(10).Cut(solid.Torus(8, 2.5).RotateX(90)),
		plane.At(v3.Zero, v3.XYZ(1, 1, 1).Normalize()),
	).Extrude(1).STL("out.stl", 8.0)
}
