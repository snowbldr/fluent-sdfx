// Text & 2D output: render a string with a TrueType font, write as PNG.
//
// shape.Text(font, text, height) returns a *Shape. height sets the cap
// height; the resulting shape is centred on the origin.
package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/tutorial/internal/tutorialfont"
)

func main() {
	// The PNG renderer shows the SDF as a heatmap, so the actual file
	// for docs is the extruded STL produced alongside.
	shape.Text(tutorialfont.Load(), "fluent-sdfx", 10).
		ToPNG("out.png", 1200, 400)
}
