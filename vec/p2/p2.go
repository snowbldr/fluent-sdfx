// Package p2 provides a 2D polar vector type.
package p2

import (
	p2sdf "github.com/deadsy/sdfx/vec/p2"
)

// Vec is a 2D polar vector with R (radius) and Theta (radians) components.
type Vec p2sdf.Vec

// R returns Vec{R: r}.
func R(r float64) Vec { return Vec{R: r} }

// T returns Vec{Theta: t}.
func T(t float64) Vec { return Vec{Theta: t} }

// RT returns Vec{R: r, Theta: t}.
func RT(r, t float64) Vec { return Vec{R: r, Theta: t} }
