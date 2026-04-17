// Package v2 provides a 2D vector type with method-chainable arithmetic.
//
// Use named constructors to avoid Vec{X: ..., Y: ...} boilerplate:
//
//	v2.X(5)      // Vec{X: 5}
//	v2.XY(1, 2)  // Vec{X: 1, Y: 2}
package v2

import (
	"github.com/deadsy/sdfx/vec/conv"
	p2sdf "github.com/deadsy/sdfx/vec/p2"
	v2sdf "github.com/deadsy/sdfx/vec/v2"
	v2isdf "github.com/deadsy/sdfx/vec/v2i"
	"github.com/snowbldr/fluent-sdfx/vec/p2"
	"github.com/snowbldr/fluent-sdfx/vec/v2i"
	"github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Vec is a 2D vector with X, Y float64 components.
type Vec v2sdf.Vec

// VecSet is a slice of Vec.
type VecSet []Vec

// X returns Vec{X: x}.
func X(x float64) Vec { return Vec{X: x} }

// Y returns Vec{Y: y}.
func Y(y float64) Vec { return Vec{Y: y} }

// XY returns Vec{X: x, Y: y}.
func XY(x, y float64) Vec { return Vec{X: x, Y: y} }

// Zero is the zero vector.
var Zero = Vec{}

// Abs returns the component-wise absolute value.
func (a Vec) Abs() Vec { return Vec(v2sdf.Vec(a).Abs()) }

// Add returns a + b.
func (a Vec) Add(b Vec) Vec { return Vec(v2sdf.Vec(a).Add(v2sdf.Vec(b))) }

// AddScalar adds b to each component.
func (a Vec) AddScalar(b float64) Vec { return Vec(v2sdf.Vec(a).AddScalar(b)) }

// Ceil returns the component-wise ceiling.
func (a Vec) Ceil() Vec { return Vec(v2sdf.Vec(a).Ceil()) }

// Clamp clamps each component between b and c.
func (a Vec) Clamp(b, c Vec) Vec {
	return Vec(v2sdf.Vec(a).Clamp(v2sdf.Vec(b), v2sdf.Vec(c)))
}

// Cross returns the 2D cross product (scalar z-component).
func (a Vec) Cross(b Vec) float64 { return v2sdf.Vec(a).Cross(v2sdf.Vec(b)) }

// Div returns component-wise a / b.
func (a Vec) Div(b Vec) Vec { return Vec(v2sdf.Vec(a).Div(v2sdf.Vec(b))) }

// DivScalar divides each component by b.
func (a Vec) DivScalar(b float64) Vec { return Vec(v2sdf.Vec(a).DivScalar(b)) }

// Dot returns the dot product a · b.
func (a Vec) Dot(b Vec) float64 { return v2sdf.Vec(a).Dot(v2sdf.Vec(b)) }

// Equals reports whether a and b are within tolerance of each other.
func (a Vec) Equals(b Vec, tolerance float64) bool {
	return v2sdf.Vec(a).Equals(v2sdf.Vec(b), tolerance)
}

// LTEZero reports whether all components are ≤ 0.
func (a Vec) LTEZero() bool { return v2sdf.Vec(a).LTEZero() }

// LTZero reports whether all components are < 0.
func (a Vec) LTZero() bool { return v2sdf.Vec(a).LTZero() }

// Length returns the Euclidean length.
func (a Vec) Length() float64 { return v2sdf.Vec(a).Length() }

// Length2 returns the squared Euclidean length.
func (a Vec) Length2() float64 { return v2sdf.Vec(a).Length2() }

// Max returns the component-wise max of a and b.
func (a Vec) Max(b Vec) Vec { return Vec(v2sdf.Vec(a).Max(v2sdf.Vec(b))) }

// MaxComponent returns the largest component.
func (a Vec) MaxComponent() float64 { return v2sdf.Vec(a).MaxComponent() }

// Min returns the component-wise min of a and b.
func (a Vec) Min(b Vec) Vec { return Vec(v2sdf.Vec(a).Min(v2sdf.Vec(b))) }

// MinComponent returns the smallest component.
func (a Vec) MinComponent() float64 { return v2sdf.Vec(a).MinComponent() }

// Mul returns component-wise a * b.
func (a Vec) Mul(b Vec) Vec { return Vec(v2sdf.Vec(a).Mul(v2sdf.Vec(b))) }

// MulScalar multiplies each component by b.
func (a Vec) MulScalar(b float64) Vec { return Vec(v2sdf.Vec(a).MulScalar(b)) }

// Neg returns -a.
func (a Vec) Neg() Vec { return Vec(v2sdf.Vec(a).Neg()) }

// Normalize returns the unit vector in the direction of a.
func (a Vec) Normalize() Vec { return Vec(v2sdf.Vec(a).Normalize()) }

// Sub returns a - b.
func (a Vec) Sub(b Vec) Vec { return Vec(v2sdf.Vec(a).Sub(v2sdf.Vec(b))) }

// SubScalar subtracts b from each component.
func (a Vec) SubScalar(b float64) Vec { return Vec(v2sdf.Vec(a).SubScalar(b)) }

// ToV3 extends to 3D with the given Z.
func (a Vec) ToV3(z float64) v3.Vec { return v3.Vec(conv.V2ToV3(v2sdf.Vec(a), z)) }

// ToV2i truncates to integer components.
func (a Vec) ToV2i() v2i.Vec { return v2i.Vec(conv.V2ToV2i(v2sdf.Vec(a))) }

// ToP2 converts cartesian to polar.
func (a Vec) ToP2() p2.Vec { return p2.Vec(conv.V2ToP2(v2sdf.Vec(a))) }

// FromV2i promotes an integer 2D vector to float64.
func FromV2i(a v2i.Vec) Vec { return Vec(conv.V2iToV2(v2isdf.Vec(a))) }

// FromP2 converts polar to cartesian.
func FromP2(a p2.Vec) Vec { return Vec(conv.P2ToV2(p2sdf.Vec(a))) }
