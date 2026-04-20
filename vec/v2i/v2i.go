// Package v2i provides a 2D integer vector type with method-chainable arithmetic.
package v2i

import (
	v2isdf "github.com/snowbldr/sdfx/vec/v2i"
)

// Vec is a 2D integer vector with X, Y int components.
type Vec v2isdf.Vec

// X returns Vec{X: x}.
func X(x int) Vec { return Vec{X: x} }

// Y returns Vec{Y: y}.
func Y(y int) Vec { return Vec{Y: y} }

// XY returns Vec{X: x, Y: y}.
func XY(x, y int) Vec { return Vec{X: x, Y: y} }

// Add returns a + b.
func (a Vec) Add(b Vec) Vec { return Vec(v2isdf.Vec(a).Add(v2isdf.Vec(b))) }

// AddScalar adds b to each component.
func (a Vec) AddScalar(b int) Vec { return Vec(v2isdf.Vec(a).AddScalar(b)) }

// SubScalar subtracts b from each component.
func (a Vec) SubScalar(b int) Vec { return Vec(v2isdf.Vec(a).SubScalar(b)) }
