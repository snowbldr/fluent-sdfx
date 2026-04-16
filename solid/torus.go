package solid

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

// Torus creates a torus (donut) centered at the origin in the XY plane.
// majorR is the distance from the center to the center of the tube.
// minorR is the radius of the tube cross-section.
func Torus(majorR, minorR float64) *Solid {
	return &Solid{&torusSDF3{majorR: majorR, minorR: minorR}}
}

type torusSDF3 struct {
	majorR float64
	minorR float64
}

func (t *torusSDF3) Evaluate(p v3.Vec) float64 {
	q := math.Sqrt(p.X*p.X+p.Y*p.Y) - t.majorR
	return math.Sqrt(q*q+p.Z*p.Z) - t.minorR
}

func (t *torusSDF3) BoundingBox() sdf.Box3 {
	r := t.majorR + t.minorR
	return sdf.Box3{
		Min: v3.Vec{X: -r, Y: -r, Z: -t.minorR},
		Max: v3.Vec{X: r, Y: r, Z: t.minorR},
	}
}
