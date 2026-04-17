package obj

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
)

// ChamferedCylinder chamfers the top (kt) and bottom (kb) of a cylinder.
func ChamferedCylinder(s *solid.Solid, kb, kt float64) *solid.Solid {
	out, err := obj.ChamferedCylinder(s.SDF3, kb, kt)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(out)
}
