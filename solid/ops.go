package solid

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
	v3i "github.com/deadsy/sdfx/vec/v3i"
	"github.com/snowbldr/fluent-sdfx/shape"
)

// --- Additional constructors ---

// Capsule returns a cylinder with hemispherical end caps.
func Capsule(height, radius float64) *Solid {
	return New(sdf.Capsule3D(height, radius))
}

// Gyroid returns an infinite gyroid surface with the given period scale per axis.
func Gyroid(scale v3.Vec) *Solid {
	return New(sdf.Gyroid3D(scale))
}

// Revolve creates a solid of revolution by rotating a 2D profile around the Y axis.
func Revolve(profile *shape.Shape) *Solid {
	return New(sdf.Revolve3D(profile))
}

// RevolveAngle creates a partial solid of revolution (theta in degrees).
func RevolveAngle(profile *shape.Shape, angleDeg float64) *Solid {
	return New(sdf.RevolveTheta3D(profile, angleDeg*math.Pi/180))
}

// ExtrudeRounded extrudes a 2D profile with rounded edges.
func ExtrudeRounded(profile *shape.Shape, height, round float64) *Solid {
	return New(sdf.ExtrudeRounded3D(profile, height, round))
}

// ScaleExtrude extrudes a 2D profile while scaling it over the height.
func ScaleExtrude(profile *shape.Shape, height float64, scale v2.Vec) *Solid {
	return &Solid{sdf.ScaleExtrude3D(profile, height, scale)}
}

// ScaleTwistExtrude extrudes a 2D profile while scaling and twisting (radians) over the height.
func ScaleTwistExtrude(profile *shape.Shape, height, twist float64, scale v2.Vec) *Solid {
	return &Solid{sdf.ScaleTwistExtrude3D(profile, height, twist, scale)}
}

// Loft transitions between two 2D profiles over a given height with optional rounding.
func Loft(bottom, top *shape.Shape, height, round float64) *Solid {
	return New(sdf.Loft3D(bottom, top, height, round))
}

// --- Additional transform methods ---

// ScaleUniform scales uniformly on all axes. Unlike Scale, distance is preserved.
func (s *Solid) ScaleUniform(k float64) *Solid {
	return &Solid{sdf.ScaleUniform3D(s, k)}
}

// MirrorXY mirrors across the XY plane (negates Z).
func (s *Solid) MirrorXY() *Solid {
	return &Solid{sdf.Transform3D(s, sdf.MirrorXY())}
}

// MirrorXZ mirrors across the XZ plane (negates Y).
func (s *Solid) MirrorXZ() *Solid {
	return &Solid{sdf.Transform3D(s, sdf.MirrorXZ())}
}

// MirrorYZ mirrors across the YZ plane (negates X).
func (s *Solid) MirrorYZ() *Solid {
	return &Solid{sdf.Transform3D(s, sdf.MirrorYZ())}
}

// --- Additional modification methods ---

// CutPlane cuts the solid along a plane. The solid on the normal side remains.
func (s *Solid) CutPlane(point, normal v3.Vec) *Solid {
	return &Solid{sdf.Cut3D(s, point, normal)}
}

// Elongate stretches the solid by the given amounts along each axis.
func (s *Solid) Elongate(h v3.Vec) *Solid {
	return &Solid{sdf.Elongate3D(s, h)}
}

// Shell hollows out the solid, leaving a shell of the given thickness.
func (s *Solid) Shell(thickness float64) *Solid {
	return New(sdf.Shell3D(s, thickness))
}

// --- Pattern/array methods ---

// Array creates an XYZ grid array of the solid.
func (s *Solid) Array(numX, numY, numZ int, step v3.Vec) *Solid {
	return &Solid{sdf.Array3D(s, v3i.Vec{numX, numY, numZ}, step)}
}

// RotateCopyZ creates N copies of the solid evenly spaced around the Z axis.
func (s *Solid) RotateCopyZ(n int) *Solid {
	return &Solid{sdf.RotateCopy3D(s, n)}
}

// Multi creates a union of the solid at the given positions.
func (s *Solid) Multi(positions []v3.Vec) *Solid {
	return &Solid{sdf.Multi3D(s, positions)}
}

// LineOf creates a union of the solid along a line from p0 to p1.
// The pattern string controls placement: 'x' places a copy, any other char skips.
func (s *Solid) LineOf(p0, p1 v3.Vec, pattern string) *Solid {
	return &Solid{sdf.LineOf3D(s, p0, p1, pattern)}
}

// Orient creates a union of the solid oriented along each direction vector.
// base is the original orientation vector of the solid.
func (s *Solid) Orient(base v3.Vec, directions []v3.Vec) *Solid {
	return &Solid{sdf.Orient3D(s, base, directions)}
}

// RotateUnionZ creates a union of the solid rotated N times by the given step matrix.
// Useful for creating patterns with custom rotation + translation per step.
func (s *Solid) RotateUnionZ(n int, step sdf.M44) *Solid {
	return &Solid{sdf.RotateUnion3D(s, n, step)}
}

// MirrorXeqY mirrors across the X==Y plane (swaps X and Y).
func (s *Solid) MirrorXeqY() *Solid {
	return &Solid{sdf.Transform3D(s, sdf.MirrorXeqY())}
}

// RotateToVector rotates the solid so that the 'from' direction aligns with the 'to' direction.
func (s *Solid) RotateToVector(from, to v3.Vec) *Solid {
	return &Solid{sdf.Transform3D(s, sdf.RotateToVector(from, to))}
}

// Center translates the solid so its bounding box center is at the origin.
func (s *Solid) Center() *Solid {
	bb := s.BoundingBox()
	center := bb.Center()
	return s.Translate(center.Neg())
}
