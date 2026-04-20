package obj

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v2i "github.com/snowbldr/fluent-sdfx/vec/v2i"
	v3i "github.com/snowbldr/fluent-sdfx/vec/v3i"
	"github.com/snowbldr/sdfx/obj"
	v2isdf "github.com/snowbldr/sdfx/vec/v2i"
	v3isdf "github.com/snowbldr/sdfx/vec/v3i"
)

// GfBaseParms configures a Gridfinity base.
type GfBaseParms struct {
	Size   v2i.Vec // size of base in gridfinity units
	Magnet bool    // add magnet mounts
	Hole   bool    // add mounting holes
}

func (p *GfBaseParms) toSDF() *obj.GfBaseParms {
	return &obj.GfBaseParms{
		Size:   v2isdf.Vec(p.Size),
		Magnet: p.Magnet,
		Hole:   p.Hole,
	}
}

// GfBodyParms configures a Gridfinity body.
type GfBodyParms struct {
	Size  v3i.Vec // size of body in gridfinity units
	Empty bool    // return an empty container
	Hole  bool    // add through holes to the body
}

func (p *GfBodyParms) toSDF() *obj.GfBodyParms {
	return &obj.GfBodyParms{
		Size:  v3isdf.Vec(p.Size),
		Empty: p.Empty,
		Hole:  p.Hole,
	}
}

// GfBase returns a Gridfinity base plate.
func GfBase(p GfBaseParms) *solid.Solid {
	return solid.Wrap(obj.GfBase(p.toSDF()))
}

// GfBody returns a Gridfinity body.
func GfBody(p GfBodyParms) *solid.Solid {
	return solid.Wrap(obj.GfBody(p.toSDF()))
}
