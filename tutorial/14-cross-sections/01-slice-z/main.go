// Cross-sections: SliceAt cuts a 2D cross-section through a 3D solid.
//
// The result is a *Shape, ready for the usual 2D operations. Here we slice
// horizontally at z=0, then re-extrude the slice 1mm thick so we can
// render it back as a 3D part.
package main

import (
	"github.com/snowbldr/fluent-sdfx/plane"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	// A 3D solid with internal structure: a hollow box with two bosses.
	body := solid.Box(v3.XYZ(20, 20, 20), 1).Shell(1.5).
		Union(
			solid.Cylinder(20, 3, 0).TranslateXY(6, 6),
			solid.Cylinder(20, 3, 0).TranslateXY(-6, -6),
		)

	shape.SliceAt(body, plane.AtZ(0)).Extrude(1).STL("out.stl", 8.0)
}
