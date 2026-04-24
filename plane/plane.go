// Package plane provides 3D plane helpers for slicing operations.
//
// Axis normals for common slice directions:
//
//	plane.X, plane.Y, plane.Z
//
// Plane values at a specific offset along an axis:
//
//	plane.AtZ(10)       // horizontal plane at z=10
//	plane.AtX(0)        // YZ plane at origin
//	plane.At(origin, n) // arbitrary plane
package plane

import (
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Axis-aligned unit normals. Useful as the `normal` argument to Slice.
var (
	X = v3.XYZ(1, 0, 0)
	Y = v3.XYZ(0, 1, 0)
	Z = v3.XYZ(0, 0, 1)
)

// Plane is an oriented plane defined by a point on the plane and a normal.
type Plane struct {
	Origin v3.Vec
	Normal v3.Vec
}

// At returns a Plane passing through origin with the given normal.
func At(origin, normal v3.Vec) Plane {
	return Plane{Origin: origin, Normal: normal}
}

// AtX returns the YZ plane at x.
func AtX(x float64) Plane {
	return Plane{Origin: v3.X(x), Normal: X}
}

// AtY returns the XZ plane at y.
func AtY(y float64) Plane {
	return Plane{Origin: v3.Y(y), Normal: Y}
}

// AtZ returns the XY plane at z.
func AtZ(z float64) Plane {
	return Plane{Origin: v3.Z(z), Normal: Z}
}

// XY returns the XY plane at the origin (normal = Z).
var XY = AtZ(0)

// XZ returns the XZ plane at the origin (normal = Y).
var XZ = AtY(0)

// YZ returns the YZ plane at the origin (normal = X).
var YZ = AtX(0)
