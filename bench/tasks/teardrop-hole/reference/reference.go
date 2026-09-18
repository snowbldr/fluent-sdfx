// Package reference: a 20x20x10 block lying flat with a 6 mm horizontal hole
// through it along X, teardropped (45 degree roof, apex up) so the hole prints
// without supports.
package reference

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	blockW = 20.0 // X
	blockD = 20.0 // Y
	blockH = 10.0 // Z
	holeR  = 3.0
)

// Teardrop returns the 2D hole profile: a circle of radius r plus a 45
// degree roof whose apex sits at r*sqrt2 above the centre (+Y in profile).
func Teardrop(r float64) *shape.Shape {
	t := r * math.Sqrt2 / 2 // tangent point offset
	roof := shape.Polygon([]v2.Vec{v2.XY(-t, t), v2.XY(t, t), v2.XY(0, r*math.Sqrt2)})
	return shape.Circle(r).Union(roof)
}

func Build() *solid.Solid {
	block := solid.Box(v3.XYZ(blockW, blockD, blockH), 0)
	// Profile in XY with apex at +Y; extrude along Z, then stand it up so the
	// hole runs along X with the apex pointing +Z.
	hole := Teardrop(holeR).Extrude(blockW + 2).RotateX(90).RotateZ(90)
	return block.Cut(hole)
}
