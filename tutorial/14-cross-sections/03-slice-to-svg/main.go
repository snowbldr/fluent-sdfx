// Cross-sections: a slice can be exported directly to SVG/DXF for
// laser cutting or 2D documentation. This step writes both the SVG and
// an extruded STL so the screenshot pipeline picks one up.
package main

import (
	"github.com/snowbldr/fluent-sdfx/plane"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	// A bracket-like part: an L-profile with a boss on one wing.
	slice := shape.SliceAt(
		solid.Box(v3.XYZ(20, 4, 14), 0.5).
			Union(
				solid.Box(v3.XYZ(4, 14, 14), 0.5).TranslateY(7),
				solid.Cylinder(14, 3, 0.5).TranslateXY(-7, 0),
			),
		plane.AtZ(0),
	)

	slice.ToSVG("bracket-section.svg", 200)
	slice.Extrude(1).STL("out.stl", 8.0)
}
