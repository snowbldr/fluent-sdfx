package obj

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
)

// FingerButtonParms configures a finger button 2D profile.
type FingerButtonParms = obj.FingerButtonParms

// FingerButton2D returns the 2D profile of a finger button.
func FingerButton2D(p FingerButtonParms) *shape.Shape {
	s, err := obj.FingerButton2D(&p)
	if err != nil {
		panic(err)
	}
	return shape.Wrap2D(s)
}
