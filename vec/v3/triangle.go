package v3

import (
	"github.com/snowbldr/sdfx/sdf"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
)

// Triangle3 is a 3D triangle defined by three vertices.
type Triangle3 [3]Vec

// SDF returns the sdfx-compatible representation of the triangle.
// Internal use only.
func (t Triangle3) SDF() sdf.Triangle3 {
	return sdf.Triangle3{v3sdf.Vec(t[0]), v3sdf.Vec(t[1]), v3sdf.Vec(t[2])}
}

// FromSDFTriangle3 promotes an sdfx Triangle3 to our Triangle3.
// Internal use only.
func FromSDFTriangle3(t sdf.Triangle3) Triangle3 {
	return Triangle3{Vec(t[0]), Vec(t[1]), Vec(t[2])}
}

// Normal returns the triangle's unit normal vector.
func (t *Triangle3) Normal() Vec {
	s := t.SDF()
	return Vec(s.Normal())
}

// Degenerate reports whether the triangle is degenerate (two or more vertices within tolerance).
func (t *Triangle3) Degenerate(tolerance float64) bool {
	s := t.SDF()
	return s.Degenerate(tolerance)
}
