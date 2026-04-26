package shape

import (
	"math"

	"github.com/snowbldr/sdfx/sdf"
)

// Line returns a line segment from (-l/2,0) to (l/2,0) with optional rounding.
func Line(length, round float64) *Shape {
	return &Shape{sdf.Line2D(length, round)}
}

// ArcSpiral returns a 2D Archimedean spiral (r = a + k*theta).
// start/end are angles in degrees, d is the offset distance (half-thickness).
func ArcSpiral(a, k, startDeg, endDeg, d float64) *Shape {
	s, err := sdf.ArcSpiral2D(a, k, startDeg*math.Pi/180, endDeg*math.Pi/180, d)
	if err != nil {
		panic(err)
	}
	return &Shape{s}
}
