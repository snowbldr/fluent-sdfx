// Text & 2D output: render a string and write to DXF (vector format used
// by laser cutters and CAD software).
//
// Cache the SDF for evaluation speedup — Text shapes are expensive to
// evaluate point-by-point because the SDF samples bezier outlines. We
// produce both a DXF (the actual output) and an extruded STL alongside
// so the screenshot pipeline has a 3D render to capture.
package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/tutorial/internal/tutorialfont"
)

func main() {
	text := shape.Text(tutorialfont.Load(), "fluent-sdfx", 10).Cache()
	text.ToDXF("out.dxf", 600)
	text.Extrude(1).STL("out.stl", 8.0)
}
