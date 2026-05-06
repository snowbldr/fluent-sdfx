package obj

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/sdfx/obj"
)

// KnurlParms configures a knurled cylindrical surface. ThetaDeg is the
// helix angle in degrees, matching the rest of the public API.
type KnurlParms struct {
	Length   float64 // length of cylinder
	Radius   float64 // radius of cylinder
	Pitch    float64 // knurl pitch
	Height   float64 // knurl height
	ThetaDeg float64 // knurl helix angle in degrees (typically ~45)
}

// Knurl3D returns a knurled cylindrical surface.
func Knurl3D(p KnurlParms) *solid.Solid {
	sdfParms := obj.KnurlParms{
		Length: p.Length,
		Radius: p.Radius,
		Pitch:  p.Pitch,
		Height: p.Height,
		Theta:  p.ThetaDeg * math.Pi / 180,
	}
	s, err := obj.Knurl3D(&sdfParms)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}

// KnurledHead3D returns a knurled head (e.g. thumb screw).
func KnurledHead3D(r, h, pitch float64) *solid.Solid {
	s, err := obj.KnurledHead3D(r, h, pitch)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}
