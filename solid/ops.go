package solid

import (
	"math"

	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/sdfx/sdf"
	v2sdf "github.com/snowbldr/sdfx/vec/v2"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
)

// Capsule returns a cylinder with hemispherical end caps.
func Capsule(height, radius float64) *Solid {
	return New(sdf.Capsule3D(height, radius))
}

// Gyroid returns an infinite gyroid surface with the given period scale per axis.
func Gyroid(scale v3.Vec) *Solid {
	return New(sdf.Gyroid3D(v3sdf.Vec(scale)))
}

// Revolve creates a solid of revolution by rotating a 2D profile around the Y axis.
func Revolve(profile sdf.SDF2) *Solid {
	return New(sdf.Revolve3D(profile))
}

// RevolveAngle creates a partial solid of revolution (theta in degrees).
func RevolveAngle(profile sdf.SDF2, angleDeg float64) *Solid {
	return New(sdf.RevolveTheta3D(profile, angleDeg*math.Pi/180))
}

// ExtrudeRounded extrudes a 2D profile with rounded edges.
func ExtrudeRounded(profile sdf.SDF2, height, round float64) *Solid {
	return New(sdf.ExtrudeRounded3D(profile, height, round))
}

// ScaleExtrude extrudes a 2D profile while scaling it over the height.
func ScaleExtrude(profile sdf.SDF2, height float64, scale v2.Vec) *Solid {
	return &Solid{sdf.ScaleExtrude3D(profile, height, v2sdf.Vec(scale))}
}

// ScaleTwistExtrude extrudes a 2D profile while scaling and twisting it
// twistDeg degrees over the height.
func ScaleTwistExtrude(profile sdf.SDF2, height, twistDeg float64, scale v2.Vec) *Solid {
	return &Solid{sdf.ScaleTwistExtrude3D(profile, height, twistDeg*math.Pi/180, v2sdf.Vec(scale))}
}

// Loft transitions between two 2D profiles over a given height with optional rounding.
func Loft(bottom, top sdf.SDF2, height, round float64) *Solid {
	return New(sdf.Loft3D(bottom, top, height, round))
}
