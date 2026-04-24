// Package v3i provides a 3D integer vector type with method-chainable arithmetic.
package v3i

import (
	v3isdf "github.com/snowbldr/sdfx/vec/v3i"
)

// Vec is a 3D integer vector with X, Y, Z int components.
type Vec v3isdf.Vec

// X returns Vec{X: x}.
func X(x int) Vec { return Vec{X: x} }

// Y returns Vec{Y: y}.
func Y(y int) Vec { return Vec{Y: y} }

// Z returns Vec{Z: z}.
func Z(z int) Vec { return Vec{Z: z} }

// XY returns Vec{X: x, Y: y}.
func XY(x, y int) Vec { return Vec{X: x, Y: y} }

// XZ returns Vec{X: x, Z: z}.
func XZ(x, z int) Vec { return Vec{X: x, Z: z} }

// YZ returns Vec{Y: y, Z: z}.
func YZ(y, z int) Vec { return Vec{Y: y, Z: z} }

// XYZ returns Vec{X: x, Y: y, Z: z}.
func XYZ(x, y, z int) Vec { return Vec{X: x, Y: y, Z: z} }

// Raw returns the underlying sdfx v3i.Vec.
func (a Vec) Raw() v3isdf.Vec { return v3isdf.Vec(a) }

// Add returns a + b.
func (a Vec) Add(b Vec) Vec { return Vec(v3isdf.Vec(a).Add(v3isdf.Vec(b))) }

// AddScalar adds b to each component.
func (a Vec) AddScalar(b int) Vec { return Vec(v3isdf.Vec(a).AddScalar(b)) }

// SubScalar subtracts b from each component.
func (a Vec) SubScalar(b int) Vec { return Vec(v3isdf.Vec(a).SubScalar(b)) }
