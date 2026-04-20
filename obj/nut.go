package obj

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/sdfx/obj"
)

// NutParms configures a nut.
type NutParms = obj.NutParms

// ThreadedCylinderParms configures a threaded cylinder.
type ThreadedCylinderParms = obj.ThreadedCylinderParms

// Nut returns a 3D nut.
func Nut(p NutParms) *solid.Solid {
	s, err := obj.Nut(&p)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}

// ThreadedCylinder returns a 3D threaded cylinder from the given parameters.
func ThreadedCylinder(p ThreadedCylinderParms) *solid.Solid {
	s, err := p.Object()
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}
