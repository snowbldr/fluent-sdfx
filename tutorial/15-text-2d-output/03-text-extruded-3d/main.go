// Text & 2D output: extrude a text profile into 3D with rounded edges.
package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/tutorial/internal/tutorialfont"
)

func main() {
	shape.Text(tutorialfont.Load(), "fluent\nsdfx", 10).
		Cache().
		ExtrudeRounded(2, 0.4).
		STL("out.stl", 10.0)
}
