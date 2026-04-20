package obj

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/sdfx/obj"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
)

// ArrowParms configures a 3D arrow.
type ArrowParms = obj.ArrowParms

// Arrow3D returns a 3D arrow aligned with the Z axis.
func Arrow3D(p ArrowParms) *solid.Solid {
	s, err := obj.Arrow3D(&p)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}

// Axes3D returns XYZ coordinate axes with arrowheads.
func Axes3D(p0, p1 v3.Vec) *solid.Solid {
	s, err := obj.Axes3D(v3sdf.Vec(p0), v3sdf.Vec(p1))
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}

// DirectedArrow3D returns an arrow from head to tail.
func DirectedArrow3D(p ArrowParms, head, tail v3.Vec) *solid.Solid {
	s, err := obj.DirectedArrow3D(&p, v3sdf.Vec(head), v3sdf.Vec(tail))
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}
