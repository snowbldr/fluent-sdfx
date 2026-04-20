package obj

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/sdfx/obj"
)

// WasherParms configures a washer.
type WasherParms = obj.WasherParms

// Washer2D returns the 2D profile of a washer.
func Washer2D(p WasherParms) *shape.Shape {
	s, err := obj.Washer2D(&p)
	if err != nil {
		panic(err)
	}
	return shape.Wrap2D(s)
}

// Washer3D returns a 3D washer.
func Washer3D(p WasherParms) *solid.Solid {
	s, err := obj.Washer3D(&p)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}
