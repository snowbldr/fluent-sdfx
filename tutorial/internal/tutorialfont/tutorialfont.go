// Package tutorialfont embeds a TrueType font (Computer Modern, cmr10) so
// every tutorial step that needs text can run with no external assets.
//
// The font file is BaKoMa Computer Modern by Basil K. Malyshev, freely
// redistributable. Substitute any other .ttf if you prefer.
package tutorialfont

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/snowbldr/fluent-sdfx/shape"
)

//go:embed cmr10.ttf
var cmr10Bytes []byte

// Load returns a *shape.Font for the embedded cmr10 font. It writes the
// font to a temp file (sdfx's LoadFont takes a path) and loads it from there.
func Load() *shape.Font {
	tmp, err := os.CreateTemp("", "fluent-sdfx-*.ttf")
	if err != nil {
		panic(err)
	}
	defer tmp.Close()
	if _, err := tmp.Write(cmr10Bytes); err != nil {
		panic(err)
	}
	return shape.LoadFont(filepath.Clean(tmp.Name()))
}
