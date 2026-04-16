package shape

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v2i "github.com/deadsy/sdfx/vec/v2i"
)

// --- Additional constructors ---

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

// --- Additional transform methods ---

// ScaleUniform scales uniformly on both axes. Unlike Scale, distance is preserved.
func (s *Shape) ScaleUniform(k float64) *Shape {
	return &Shape{sdf.ScaleUniform2D(s, k)}
}

// Center translates the shape so its bounding box center is at the origin.
func (s *Shape) Center() *Shape {
	return &Shape{sdf.Center2D(s)}
}

// CenterAndScale centers on the bounding box then scales uniformly. Distance is preserved.
func (s *Shape) CenterAndScale(k float64) *Shape {
	return &Shape{sdf.CenterAndScale2D(s, k)}
}

// --- Additional modification methods ---

// CutLine cuts the shape along a line from point a in direction v.
// The shape to the right of the line remains.
func (s *Shape) CutLine(a, dir v2.Vec) *Shape {
	return &Shape{sdf.Cut2D(s, a, dir)}
}

// Elongate stretches the shape by the given amounts along each axis.
func (s *Shape) Elongate(h v2.Vec) *Shape {
	return &Shape{sdf.Elongate2D(s, h)}
}

// --- Pattern/array methods ---

// Array creates an XY grid array of the shape.
func (s *Shape) Array(numX, numY int, step v2.Vec) *Shape {
	return &Shape{sdf.Array2D(s, v2i.Vec{numX, numY}, step)}
}

// RotateCopy creates N copies of the shape evenly spaced in a full circle.
func (s *Shape) RotateCopy(n int) *Shape {
	return &Shape{sdf.RotateCopy2D(s, n)}
}

// Multi creates a union of the shape at the given positions.
func (s *Shape) Multi(positions []v2.Vec) *Shape {
	return &Shape{sdf.Multi2D(s, positions)}
}

// LineOf creates a union of the shape along a line from p0 to p1.
// The pattern string controls placement: 'x' places a copy, any other char skips.
func (s *Shape) LineOf(p0, p1 v2.Vec, pattern string) *Shape {
	return &Shape{sdf.LineOf2D(s, p0, p1, pattern)}
}

// RotateUnion creates a union of the shape rotated N times by the given step matrix.
func (s *Shape) RotateUnion(n int, step sdf.M33) *Shape {
	return &Shape{sdf.RotateUnion2D(s, n, step)}
}
