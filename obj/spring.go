package obj

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
)

// SpringParms configures a flat planar spring.
type SpringParms = obj.SpringParms

// SpringLength returns the overall length of the spring defined by p.
func SpringLength(p SpringParms) float64 {
	return p.SpringLength()
}

// Spring2D returns the 2D profile of a flat planar spring.
func Spring2D(p SpringParms) *shape.Shape {
	s, err := p.Spring2D()
	if err != nil {
		panic(err)
	}
	return shape.Wrap2D(s)
}

// Spring3D returns a 3D extruded flat planar spring.
func Spring3D(p SpringParms) *solid.Solid {
	s, err := p.Spring3D()
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}
