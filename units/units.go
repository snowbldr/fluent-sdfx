// Package units re-exports common numeric constants and helpers from sdfx.
package units

import "github.com/deadsy/sdfx/sdf"

// Math constants.
const (
	Pi  = sdf.Pi
	Tau = sdf.Tau
)

// Length conversions.
const (
	// Mil is 1/1000 of an inch in mm.
	Mil = sdf.Mil
	// MillimetresPerInch converts inches to millimetres.
	MillimetresPerInch = sdf.MillimetresPerInch
	// InchesPerMillimetre converts millimetres to inches.
	InchesPerMillimetre = sdf.InchesPerMillimetre
)

// DtoR converts degrees to radians.
func DtoR(degrees float64) float64 { return sdf.DtoR(degrees) }

// RtoD converts radians to degrees.
func RtoD(radians float64) float64 { return radians * 180 / sdf.Pi }

// EqualFloat64 reports whether a and b are equal within epsilon.
func EqualFloat64(a, b, epsilon float64) bool { return sdf.EqualFloat64(a, b, epsilon) }

// ErrMsg creates an error with the given message plus a stack frame hint.
func ErrMsg(msg string) error { return sdf.ErrMsg(msg) }
