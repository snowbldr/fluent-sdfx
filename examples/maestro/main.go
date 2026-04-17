package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%

func servoControllerMount() *solid.Solid {
	// standoff
	k0 := obj.StandoffParms{
		PillarHeight:   0.5 * units.MillimetresPerInch,
		PillarDiameter: 5,
		HoleDepth:      10,
		HoleDiameter:   2.4, // #4 screw
	}
	standoff := obj.Standoff3D(k0)

	// standoffs
	h0 := v3.XYZ(-0.45, -0.8, 0.25).MulScalar(units.MillimetresPerInch)
	h1 := v3.XYZ(0.05, 0.8, 0.25).MulScalar(units.MillimetresPerInch)
	standoffs := standoff.Multi([]v3.Vec{h0, h1})

	// base
	k1 := obj.PanelParms{
		Size:         v2.XY(1.1, 1.8).MulScalar(units.MillimetresPerInch),
		CornerRadius: 2,
		HoleDiameter: 2.4, // #4 screw
		HoleMargin:   [4]float64{4, 4, 4, 4},
		HolePattern:  [4]string{"x", "x", ".x", ""},
		Thickness:    3,
	}
	base := obj.Panel3D(k1)

	return base.Union(standoffs)
}

func main() {
	servoControllerMount().ScaleUniform(shrink).ToSTL("mm18.stl", 300)
}
