package obj

import (
	"github.com/deadsy/sdfx/obj"
	v3sdf "github.com/deadsy/sdfx/vec/v3"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// PanelBoxParms configures a panel box enclosure.
type PanelBoxParms struct {
	Size       v3.Vec  // outer box dimensions (width, height, length)
	Wall       float64 // wall thickness
	Panel      float64 // front/back panel thickness
	Rounding   float64 // radius of corner rounding
	FrontInset float64 // inset depth of box front
	BackInset  float64 // inset depth of box back
	Clearance  float64 // fit clearance (typically 0.05)
	Hole       float64 // diameter of screw holes
	SideTabs   string  // tab pattern b/B (bottom) t/T (top) . (empty)
}

func (p *PanelBoxParms) toSDF() *obj.PanelBoxParms {
	return &obj.PanelBoxParms{
		Size:       v3sdf.Vec(p.Size),
		Wall:       p.Wall,
		Panel:      p.Panel,
		Rounding:   p.Rounding,
		FrontInset: p.FrontInset,
		BackInset:  p.BackInset,
		Clearance:  p.Clearance,
		Hole:       p.Hole,
		SideTabs:   p.SideTabs,
	}
}

// PanelBox3D returns the set of solids making up a panel box enclosure.
func PanelBox3D(p PanelBoxParms) []*solid.Solid {
	parts, err := obj.PanelBox3D(p.toSDF())
	if err != nil {
		panic(err)
	}
	out := make([]*solid.Solid, len(parts))
	for i, s := range parts {
		out[i] = solid.Wrap(s)
	}
	return out
}
