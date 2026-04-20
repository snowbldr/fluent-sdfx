package obj

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/sdfx/obj"
)

// Hex2D returns a 2D regular hexagon with optional corner rounding.
func Hex2D(radius, round float64) *shape.Shape {
	s, err := obj.Hex2D(radius, round)
	if err != nil {
		panic(err)
	}
	return shape.Wrap2D(s)
}

// Hex3D returns a 3D hex prism.
func Hex3D(radius, height, round float64) *solid.Solid {
	s, err := obj.Hex3D(radius, height, round)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}

// HexHead3D returns a hex head with chamfered or rounded top/bottom.
// round is "tb" (top+bottom), "t" (top), "b" (bottom), or "" (none).
func HexHead3D(radius, height float64, round string) *solid.Solid {
	s, err := obj.HexHead3D(radius, height, round)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}
