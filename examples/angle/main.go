package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/units"
)

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%

func main() {
	const l = 1.25 * units.MillimetresPerInch
	const t = 0.125 * units.MillimetresPerInch
	const r = 0.125 * units.MillimetresPerInch

	k := obj.AngleParams{
		X:          obj.AngleLeg{Length: l, Thickness: t},
		Y:          obj.AngleLeg{Length: l, Thickness: t},
		RootRadius: r,
		Length:     12 * units.MillimetresPerInch,
	}

	obj.Angle3D(k).ScaleUniform(shrink).ToSTL("angle.stl", 300)
}
