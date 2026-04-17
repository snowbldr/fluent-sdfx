package obj

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
)

// InvoluteGearParms configures an involute gear 2D profile.
type InvoluteGearParms = obj.InvoluteGearParms

// InvoluteGear returns a 2D involute gear profile.
func InvoluteGear(p InvoluteGearParms) *shape.Shape {
	s, err := obj.InvoluteGear(&p)
	if err != nil {
		panic(err)
	}
	return shape.Wrap2D(s)
}
