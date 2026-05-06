package obj

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/sdfx/obj"
)

// InvoluteGearParms configures an involute gear 2D profile. PressureAngleDeg
// is in degrees (the typical engineering value is 20°), matching the rest of
// the public API; the underlying sdfx kernel uses radians internally.
type InvoluteGearParms struct {
	NumberTeeth      int     // number of gear teeth
	Module           float64 // pitch circle diameter / number of gear teeth
	PressureAngleDeg float64 // gear pressure angle in degrees (typically 20)
	Backlash         float64 // backlash expressed as per-tooth distance at pitch circumference
	Clearance        float64 // additional root clearance
	RingWidth        float64 // width of ring wall (from root circle)
	Facets           int     // number of facets for involute flank
}

// InvoluteGear returns a 2D involute gear profile.
func InvoluteGear(p InvoluteGearParms) *shape.Shape {
	sdfParms := obj.InvoluteGearParms{
		NumberTeeth:   p.NumberTeeth,
		Module:        p.Module,
		PressureAngle: p.PressureAngleDeg * math.Pi / 180,
		Backlash:      p.Backlash,
		Clearance:     p.Clearance,
		RingWidth:     p.RingWidth,
		Facets:        p.Facets,
	}
	s, err := obj.InvoluteGear(&sdfParms)
	if err != nil {
		panic(err)
	}
	return shape.Wrap2D(s)
}
