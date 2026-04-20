package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/units"
)

// Tolerance: Measured in mm. Typically 0.0 to 0.4. Larger is looser.
// Smaller is tighter.
const mmTolerance = 0.3
const inchTolerance = mmTolerance / units.MillimetresPerInch

// Quality: mesh resolution in cells per millimeter.
const quality = 2.0

func inch() {
	obj.Bolt(obj.BoltParms{
		Thread:      "unc_5/8",
		Style:       "knurl",
		Tolerance:   inchTolerance,
		TotalLength: 2.0,
		ShankLength: 0.5,
	}).ScaleUniform(units.MillimetresPerInch).STL("inch_bolt.stl", quality)

	obj.Nut(obj.NutParms{
		Thread:    "unc_5/8",
		Style:     "knurl",
		Tolerance: inchTolerance,
	}).ScaleUniform(units.MillimetresPerInch).STL("inch_nut.stl", quality)
}

func metric() {
	obj.Bolt(obj.BoltParms{
		Thread:      "M16x2",
		Style:       "hex",
		Tolerance:   mmTolerance,
		TotalLength: 50.0,
		ShankLength: 10.0,
	}).STL("metric_bolt.stl", quality)

	obj.Nut(obj.NutParms{
		Thread:    "M16x2",
		Style:     "hex",
		Tolerance: mmTolerance,
	}).STL("metric_nut.stl", quality)
}

func main() {
	inch()
	metric()
}
