package shape

import (
	"github.com/deadsy/sdfx/sdf"
	v2sdf "github.com/deadsy/sdfx/vec/v2"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

// M33 is a 3x3 homogeneous transformation matrix for 2D.
type M33 sdf.M33

// NewM33 returns a 3x3 matrix from a row-major 9-element array.
func NewM33(x [9]float64) M33 { return M33(sdf.NewM33(x)) }

// Translate2d returns a 3x3 translation matrix.
func Translate2d(v v2.Vec) M33 { return M33(sdf.Translate2d(v2sdf.Vec(v))) }

// Rotate2d returns a 3x3 rotation matrix around the origin (degrees).
func Rotate2d(deg float64) M33 { return M33(sdf.Rotate2d(deg * sdf.Pi / 180)) }

// Scale2d returns a 3x3 scale matrix.
func Scale2d(v v2.Vec) M33 { return M33(sdf.Scale2d(v2sdf.Vec(v))) }

// Identity2d returns the 3x3 identity matrix.
func Identity2d() M33 { return M33(sdf.Identity2d()) }

// MirrorX returns a 3x3 matrix mirroring across the X axis.
func MirrorX() M33 { return M33(sdf.MirrorX()) }

// MirrorY returns a 3x3 matrix mirroring across the Y axis.
func MirrorY() M33 { return M33(sdf.MirrorY()) }

// Add returns a + b.
func (a M33) Add(b M33) M33 { return M33(sdf.M33(a).Add(sdf.M33(b))) }

// Mul returns a * b.
func (a M33) Mul(b M33) M33 { return M33(sdf.M33(a).Mul(sdf.M33(b))) }

// MulScalar multiplies every entry by k.
func (a M33) MulScalar(k float64) M33 { return M33(sdf.M33(a).MulScalar(k)) }

// Determinant returns the matrix determinant.
func (a M33) Determinant() float64 { return sdf.M33(a).Determinant() }

// Inverse returns the inverse matrix.
func (a M33) Inverse() M33 { return M33(sdf.M33(a).Inverse()) }

// Equals reports whether a and b are within tolerance of each other.
func (a M33) Equals(b M33, tolerance float64) bool {
	return sdf.M33(a).Equals(sdf.M33(b), tolerance)
}

// Values returns the matrix as a 9-element array.
func (a M33) Values() [9]float64 { return sdf.M33(a).Values() }

// MulPosition transforms a position vector.
func (a M33) MulPosition(v v2.Vec) v2.Vec {
	return v2.Vec(sdf.M33(a).MulPosition(v2sdf.Vec(v)))
}

// MulBox transforms a bounding box.
func (a M33) MulBox(box Box2) Box2 {
	return v2.FromSDF(sdf.M33(a).MulBox(box.SDF()))
}
