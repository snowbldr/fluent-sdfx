package shape

import (
	"github.com/deadsy/sdfx/sdf"
	"github.com/golang/freetype/truetype"
)

// Font is an opaque truetype font handle for text rendering.
type Font = truetype.Font

// LoadFont loads a truetype font (*.ttf) from disk.
func LoadFont(path string) *Font {
	f, err := sdf.LoadFont(path)
	if err != nil {
		panic(err)
	}
	return f
}

// Text renders a string as a 2D Shape using the given font and total text height.
// The shape is centered horizontally and vertically on the origin (scaled to fit height h).
func Text(font *Font, text string, height float64) *Shape {
	s, err := sdf.Text2D(font, sdf.NewText(text), height)
	if err != nil {
		panic(err)
	}
	return &Shape{s}
}
