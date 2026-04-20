package obj

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/sdfx/obj"
)

// AngleLeg describes one leg of an angle bracket.
type AngleLeg = obj.AngleLeg

// AngleParams configures an angle bracket (L-profile).
type AngleParams = obj.AngleParms

// Angle2D returns the 2D cross-section of an L-profile angle bracket.
func Angle2D(p AngleParams) *shape.Shape {
	s, err := obj.Angle2D(&p)
	if err != nil {
		panic(err)
	}
	return shape.Wrap2D(s)
}

// Angle3D returns an extruded L-profile angle bracket.
func Angle3D(p AngleParams) *solid.Solid {
	s, err := obj.Angle3D(&p)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}
