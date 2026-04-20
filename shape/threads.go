package shape

import "github.com/snowbldr/sdfx/sdf"

// AcmeThread returns the 2D profile of an ACME thread of the given radius and pitch.
func AcmeThread(radius, pitch float64) *Shape {
	s, err := sdf.AcmeThread(radius, pitch)
	if err != nil {
		panic(err)
	}
	return &Shape{s}
}

// ISOThread returns the 2D profile of an ISO/UTS metric thread.
// If external is true, returns an external (bolt) thread; otherwise internal (nut).
func ISOThread(radius, pitch float64, external bool) *Shape {
	s, err := sdf.ISOThread(radius, pitch, external)
	if err != nil {
		panic(err)
	}
	return &Shape{s}
}

// ANSIButtressThread returns an ANSI 45/7 buttress thread profile.
func ANSIButtressThread(radius, pitch float64) *Shape {
	s, err := sdf.ANSIButtressThread(radius, pitch)
	if err != nil {
		panic(err)
	}
	return &Shape{s}
}

// PlasticButtressThread returns a plastic-screw-top-style buttress thread profile.
func PlasticButtressThread(radius, pitch float64) *Shape {
	s, err := sdf.PlasticButtressThread(radius, pitch)
	if err != nil {
		panic(err)
	}
	return &Shape{s}
}

// ThreadParams reexposes the sdfx ThreadParameters struct for thread lookup.
type ThreadParams = sdf.ThreadParameters

// ThreadLookup looks up a standard thread by name (e.g. "M4x0.7", "unc_1/4", "npt_1/2").
func ThreadLookup(name string) *ThreadParams {
	t, err := sdf.ThreadLookup(name)
	if err != nil {
		panic(err)
	}
	return t
}
