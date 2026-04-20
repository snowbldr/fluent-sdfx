package shape

import (
	"math"

	"github.com/snowbldr/sdfx/sdf"
)

// GearRackParams defines the parameters for a 2D gear rack.
type GearRackParams struct {
	NumberTeeth      int     // number of rack teeth
	Module           float64 // pitch circle diameter / number of gear teeth
	PressureAngleDeg float64 // pressure angle in degrees
	Backlash         float64 // backlash in units of pitch circumference
	BaseHeight       float64 // height of rack base
}

// GearRack returns the 2D profile for a linear gear rack.
func GearRack(p GearRackParams) *Shape {
	s, err := sdf.GearRack2D(&sdf.GearRackParms{
		NumberTeeth:   p.NumberTeeth,
		Module:        p.Module,
		PressureAngle: p.PressureAngleDeg * math.Pi / 180,
		Backlash:      p.Backlash,
		BaseHeight:    p.BaseHeight,
	})
	if err != nil {
		panic(err)
	}
	return &Shape{s}
}
