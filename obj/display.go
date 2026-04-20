package obj

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	"github.com/snowbldr/sdfx/obj"
	v2sdf "github.com/snowbldr/sdfx/vec/v2"
)

// DisplayParms configures a display bezel.
type DisplayParms struct {
	Window          v2.Vec  // display window
	Rounding        float64 // window corner rounding
	Supports        v2.Vec  // support positions
	SupportHeight   float64 // support height
	SupportDiameter float64 // support diameter
	HoleDiameter    float64 // support hole diameter
	Offset          v2.Vec  // offset between window and supports
	Thickness       float64 // panel thickness
	Countersunk     bool    // counter sink screws on panel face
}

func (p *DisplayParms) toSDF() *obj.DisplayParms {
	return &obj.DisplayParms{
		Window:          v2sdf.Vec(p.Window),
		Rounding:        p.Rounding,
		Supports:        v2sdf.Vec(p.Supports),
		SupportHeight:   p.SupportHeight,
		SupportDiameter: p.SupportDiameter,
		HoleDiameter:    p.HoleDiameter,
		Offset:          v2sdf.Vec(p.Offset),
		Thickness:       p.Thickness,
		Countersunk:     p.Countersunk,
	}
}

// Display returns a 3D display bezel. If negative is true the cutouts are returned.
func Display(p DisplayParms, negative bool) *solid.Solid {
	s, err := obj.Display(p.toSDF(), negative)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}
