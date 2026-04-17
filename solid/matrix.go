package solid

import (
	"github.com/deadsy/sdfx/sdf"
	v3sdf "github.com/deadsy/sdfx/vec/v3"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// M44 is a 4x4 homogeneous transformation matrix.
type M44 sdf.M44

// NewM44 returns a 4x4 matrix from a row-major 16-element array.
func NewM44(x [16]float64) M44 { return M44(sdf.NewM44(x)) }

// Translate3d returns a 4x4 translation matrix.
func Translate3d(v v3.Vec) M44 { return M44(sdf.Translate3d(v3sdf.Vec(v))) }

// RotateXMatrix returns a 4x4 rotation matrix around the X axis (degrees).
func RotateXMatrix(deg float64) M44 { return M44(sdf.RotateX(deg * sdf.Pi / 180)) }

// RotateYMatrix returns a 4x4 rotation matrix around the Y axis (degrees).
func RotateYMatrix(deg float64) M44 { return M44(sdf.RotateY(deg * sdf.Pi / 180)) }

// RotateZMatrix returns a 4x4 rotation matrix around the Z axis (degrees).
func RotateZMatrix(deg float64) M44 { return M44(sdf.RotateZ(deg * sdf.Pi / 180)) }

// Rotate3dMatrix returns a 4x4 rotation matrix around an arbitrary axis (degrees).
func Rotate3dMatrix(axis v3.Vec, deg float64) M44 {
	return M44(sdf.Rotate3d(v3sdf.Vec(axis), deg*sdf.Pi/180))
}

// RotateToVector returns a 4x4 matrix rotating a onto b.
func RotateToVector(a, b v3.Vec) M44 {
	return M44(sdf.RotateToVector(v3sdf.Vec(a), v3sdf.Vec(b)))
}

// Scale3d returns a 4x4 scale matrix.
func Scale3d(v v3.Vec) M44 { return M44(sdf.Scale3d(v3sdf.Vec(v))) }

// Identity3d returns the 4x4 identity matrix.
func Identity3d() M44 { return M44(sdf.Identity3d()) }

// MirrorXY returns a 4x4 mirror matrix about the XY plane.
func MirrorXY() M44 { return M44(sdf.MirrorXY()) }

// MirrorXZ returns a 4x4 mirror matrix about the XZ plane.
func MirrorXZ() M44 { return M44(sdf.MirrorXZ()) }

// MirrorYZ returns a 4x4 mirror matrix about the YZ plane.
func MirrorYZ() M44 { return M44(sdf.MirrorYZ()) }

// MirrorXeqY returns a 4x4 mirror matrix about the X=Y plane.
func MirrorXeqY() M44 { return M44(sdf.MirrorXeqY()) }

// Mul returns a * b.
func (a M44) Mul(b M44) M44 { return M44(sdf.M44(a).Mul(sdf.M44(b))) }

// Determinant returns the matrix determinant.
func (a M44) Determinant() float64 { return sdf.M44(a).Determinant() }

// Inverse returns the inverse matrix.
func (a M44) Inverse() M44 { return M44(sdf.M44(a).Inverse()) }

// Equals reports whether a and b are within tolerance of each other.
func (a M44) Equals(b M44, tolerance float64) bool {
	return sdf.M44(a).Equals(sdf.M44(b), tolerance)
}

// Values returns the matrix as a 16-element array.
func (a M44) Values() [16]float64 { return sdf.M44(a).Values() }

// MulPosition transforms a position vector.
func (a M44) MulPosition(v v3.Vec) v3.Vec {
	return v3.Vec(sdf.M44(a).MulPosition(v3sdf.Vec(v)))
}

// MulBox transforms a bounding box.
func (a M44) MulBox(box Box3) Box3 {
	return v3.FromSDF(sdf.M44(a).MulBox(box.SDF()))
}
