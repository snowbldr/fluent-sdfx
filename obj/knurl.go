package obj

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
)

// KnurlParms configures a knurled surface.
type KnurlParms = obj.KnurlParms

// Knurl3D returns a knurled cylindrical surface.
func Knurl3D(p KnurlParms) *solid.Solid {
	s, err := obj.Knurl3D(&p)
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
