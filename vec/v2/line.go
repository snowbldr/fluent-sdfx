package v2

import (
	"github.com/snowbldr/sdfx/sdf"
	v2sdf "github.com/snowbldr/sdfx/vec/v2"
)

// Line2 is a 2D line segment defined by two endpoints.
type Line2 [2]Vec

// SDF returns the sdfx-compatible representation of the line.
// Internal use only.
func (a Line2) SDF() sdf.Line2 {
	return sdf.Line2{v2sdf.Vec(a[0]), v2sdf.Vec(a[1])}
}

// FromSDFLine2 promotes an sdfx Line2 to our Line2.
// Internal use only.
func FromSDFLine2(a sdf.Line2) Line2 {
	return Line2{Vec(a[0]), Vec(a[1])}
}

// Triangle2 is a 2D triangle defined by three vertices.
type Triangle2 [3]Vec

// SDF returns the sdfx-compatible representation of the triangle.
// Internal use only.
func (a Triangle2) SDF() sdf.Triangle2 {
	return sdf.Triangle2{v2sdf.Vec(a[0]), v2sdf.Vec(a[1]), v2sdf.Vec(a[2])}
}

// FromSDFTriangle2 promotes an sdfx Triangle2 to our Triangle2.
// Internal use only.
func FromSDFTriangle2(a sdf.Triangle2) Triangle2 {
	return Triangle2{Vec(a[0]), Vec(a[1]), Vec(a[2])}
}

// BoundingBox returns the axis-aligned bounding box of the line.
func (a *Line2) BoundingBox() Box {
	s := a.SDF()
	b := s.BoundingBox()
	return FromSDF(b)
}

// Degenerate reports whether the line endpoints are within tolerance.
func (a Line2) Degenerate(tolerance float64) bool { return a.SDF().Degenerate(tolerance) }

// Reverse returns the line with endpoints swapped.
func (a Line2) Reverse() Line2 { return Line2{a[1], a[0]} }
