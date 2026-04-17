// Package v3 provides a 3D vector type with method-chainable arithmetic.
//
// Use named constructors to avoid Vec{X: ..., Y: ...} boilerplate:
//
//	v3.X(5)          // Vec{X: 5}
//	v3.YZ(2, 3)      // Vec{Y: 2, Z: 3}
//	v3.XYZ(1, 2, 3)  // Vec{X: 1, Y: 2, Z: 3}
package v3

import (
	"github.com/deadsy/sdfx/vec/conv"
	v3sdf "github.com/deadsy/sdfx/vec/v3"
	v3isdf "github.com/deadsy/sdfx/vec/v3i"
	"github.com/snowbldr/fluent-sdfx/vec/v3i"
)

// Vec is a 3D vector with X, Y, Z float64 components.
type Vec v3sdf.Vec

// VecSet is a slice of Vec.
type VecSet []Vec

// X returns Vec{X: x}.
func X(x float64) Vec { return Vec{X: x} }

// Y returns Vec{Y: y}.
func Y(y float64) Vec { return Vec{Y: y} }

// Z returns Vec{Z: z}.
func Z(z float64) Vec { return Vec{Z: z} }

// XY returns Vec{X: x, Y: y}.
func XY(x, y float64) Vec { return Vec{X: x, Y: y} }

// XZ returns Vec{X: x, Z: z}.
func XZ(x, z float64) Vec { return Vec{X: x, Z: z} }

// YZ returns Vec{Y: y, Z: z}.
func YZ(y, z float64) Vec { return Vec{Y: y, Z: z} }

// XYZ returns Vec{X: x, Y: y, Z: z}.
func XYZ(x, y, z float64) Vec { return Vec{X: x, Y: y, Z: z} }

// Zero is the zero vector.
var Zero = Vec{}

// Abs returns the component-wise absolute value.
func (a Vec) Abs() Vec { return Vec(v3sdf.Vec(a).Abs()) }

// Add returns a + b.
func (a Vec) Add(b Vec) Vec { return Vec(v3sdf.Vec(a).Add(v3sdf.Vec(b))) }

// AddScalar adds b to each component.
func (a Vec) AddScalar(b float64) Vec { return Vec(v3sdf.Vec(a).AddScalar(b)) }

// Ceil returns the component-wise ceiling.
func (a Vec) Ceil() Vec { return Vec(v3sdf.Vec(a).Ceil()) }

// Clamp clamps each component between b and c.
func (a Vec) Clamp(b, c Vec) Vec {
	return Vec(v3sdf.Vec(a).Clamp(v3sdf.Vec(b), v3sdf.Vec(c)))
}

// Cos returns the component-wise cosine.
func (a Vec) Cos() Vec { return Vec(v3sdf.Vec(a).Cos()) }

// Cross returns the cross product a × b.
func (a Vec) Cross(b Vec) Vec { return Vec(v3sdf.Vec(a).Cross(v3sdf.Vec(b))) }

// Div returns component-wise a / b.
func (a Vec) Div(b Vec) Vec { return Vec(v3sdf.Vec(a).Div(v3sdf.Vec(b))) }

// DivScalar divides each component by b.
func (a Vec) DivScalar(b float64) Vec { return Vec(v3sdf.Vec(a).DivScalar(b)) }

// Dot returns the dot product a · b.
func (a Vec) Dot(b Vec) float64 { return v3sdf.Vec(a).Dot(v3sdf.Vec(b)) }

// Equals reports whether a and b are within tolerance of each other.
func (a Vec) Equals(b Vec, tolerance float64) bool {
	return v3sdf.Vec(a).Equals(v3sdf.Vec(b), tolerance)
}

// Get returns the i-th component (0=X, 1=Y, 2=Z).
func (a Vec) Get(i int) float64 { return v3sdf.Vec(a).Get(i) }

// LTEZero reports whether all components are ≤ 0.
func (a Vec) LTEZero() bool { return v3sdf.Vec(a).LTEZero() }

// LTZero reports whether all components are < 0.
func (a Vec) LTZero() bool { return v3sdf.Vec(a).LTZero() }

// Length returns the Euclidean length.
func (a Vec) Length() float64 { return v3sdf.Vec(a).Length() }

// Length2 returns the squared Euclidean length.
func (a Vec) Length2() float64 { return v3sdf.Vec(a).Length2() }

// Max returns the component-wise max of a and b.
func (a Vec) Max(b Vec) Vec { return Vec(v3sdf.Vec(a).Max(v3sdf.Vec(b))) }

// MaxComponent returns the largest component.
func (a Vec) MaxComponent() float64 { return v3sdf.Vec(a).MaxComponent() }

// Min returns the component-wise min of a and b.
func (a Vec) Min(b Vec) Vec { return Vec(v3sdf.Vec(a).Min(v3sdf.Vec(b))) }

// MinComponent returns the smallest component.
func (a Vec) MinComponent() float64 { return v3sdf.Vec(a).MinComponent() }

// Mul returns component-wise a * b.
func (a Vec) Mul(b Vec) Vec { return Vec(v3sdf.Vec(a).Mul(v3sdf.Vec(b))) }

// MulScalar multiplies each component by b.
func (a Vec) MulScalar(b float64) Vec { return Vec(v3sdf.Vec(a).MulScalar(b)) }

// Neg returns -a.
func (a Vec) Neg() Vec { return Vec(v3sdf.Vec(a).Neg()) }

// Normalize returns the unit vector in the direction of a.
func (a Vec) Normalize() Vec { return Vec(v3sdf.Vec(a).Normalize()) }

// Set sets the i-th component (0=X, 1=Y, 2=Z) to val.
func (a *Vec) Set(i int, val float64) { (*v3sdf.Vec)(a).Set(i, val) }

// Sin returns the component-wise sine.
func (a Vec) Sin() Vec { return Vec(v3sdf.Vec(a).Sin()) }

// Sub returns a - b.
func (a Vec) Sub(b Vec) Vec { return Vec(v3sdf.Vec(a).Sub(v3sdf.Vec(b))) }

// SubScalar subtracts b from each component.
func (a Vec) SubScalar(b float64) Vec { return Vec(v3sdf.Vec(a).SubScalar(b)) }

// ToV3i truncates to integer components.
func (a Vec) ToV3i() v3i.Vec { return v3i.Vec(conv.V3ToV3i(v3sdf.Vec(a))) }

// FromV3i promotes an integer 3D vector to float64.
func FromV3i(a v3i.Vec) Vec { return Vec(conv.V3iToV3(v3isdf.Vec(a))) }
