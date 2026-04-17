package obj

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
)

// KeywayParameters configures a keyway.
type KeywayParameters = obj.KeywayParameters

// Keyway2D returns the 2D cross-section of a shaft with a keyway.
func Keyway2D(p KeywayParameters) *shape.Shape {
	s, err := obj.Keyway2D(&p)
	if err != nil {
		panic(err)
	}
	return shape.Wrap2D(s)
}

// Keyway3D returns a 3D shaft with a keyway.
func Keyway3D(p KeywayParameters) *solid.Solid {
	s, err := obj.Keyway3D(&p)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}
