package obj

import (
	"github.com/deadsy/sdfx/obj"
	v2sdf "github.com/deadsy/sdfx/vec/v2"
	v3sdf "github.com/deadsy/sdfx/vec/v3"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// PanelParms configures a 2D/3D panel with holes.
type PanelParms struct {
	Size         v2.Vec     // size of the panel
	CornerRadius float64    // radius of rounded corners
	HoleDiameter float64    // diameter of panel holes
	HoleMargin   [4]float64 // hole margins for top, right, bottom, left
	HolePattern  [4]string  // hole pattern for top, right, bottom, left
	Thickness    float64    // panel thickness (3d only)
	Ridge        v2.Vec     // add side ridges for reinforcing (3d only)
}

func (p *PanelParms) toSDF() *obj.PanelParms {
	return &obj.PanelParms{
		Size:         v2sdf.Vec(p.Size),
		CornerRadius: p.CornerRadius,
		HoleDiameter: p.HoleDiameter,
		HoleMargin:   p.HoleMargin,
		HolePattern:  p.HolePattern,
		Thickness:    p.Thickness,
		Ridge:        v2sdf.Vec(p.Ridge),
	}
}

// EuroRackParms configures a Eurorack module panel.
type EuroRackParms = obj.EuroRackParms

// PanelHoleParms configures a through-panel hole with indent/orientation.
type PanelHoleParms struct {
	Diameter    float64 // hole diameter
	Thickness   float64 // panel thickness
	Indent      v3.Vec  // indent size
	Offset      float64 // indent offset from main axis
	Orientation float64 // orientation of indent, 0 == x-axis
}

func (p *PanelHoleParms) toSDF() *obj.PanelHoleParms {
	return &obj.PanelHoleParms{
		Diameter:    p.Diameter,
		Thickness:   p.Thickness,
		Indent:      v3sdf.Vec(p.Indent),
		Offset:      p.Offset,
		Orientation: p.Orientation,
	}
}

// Panel2D returns a 2D panel profile with mounting holes.
func Panel2D(p PanelParms) *shape.Shape {
	s, err := obj.Panel2D(p.toSDF())
	if err != nil {
		panic(err)
	}
	return shape.Wrap2D(s)
}

// Panel3D returns a 3D panel with mounting holes.
func Panel3D(p PanelParms) *solid.Solid {
	s, err := obj.Panel3D(p.toSDF())
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}

// EuroRackPanel2D returns a 2D Eurorack module panel.
func EuroRackPanel2D(p EuroRackParms) *shape.Shape {
	s, err := obj.EuroRackPanel2D(&p)
	if err != nil {
		panic(err)
	}
	return shape.Wrap2D(s)
}

// EuroRackPanel3D returns a 3D Eurorack module panel.
func EuroRackPanel3D(p EuroRackParms) *solid.Solid {
	s, err := obj.EuroRackPanel3D(&p)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}

// PanelHole3D returns a 3D panel hole with optional indent.
func PanelHole3D(p PanelHoleParms) *solid.Solid {
	s, err := obj.PanelHole3D(p.toSDF())
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}
