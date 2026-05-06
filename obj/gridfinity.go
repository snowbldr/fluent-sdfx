package obj

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v2i "github.com/snowbldr/fluent-sdfx/vec/v2i"
	v3i "github.com/snowbldr/fluent-sdfx/vec/v3i"
	"github.com/snowbldr/sdfx/obj"
	v2isdf "github.com/snowbldr/sdfx/vec/v2i"
	v3isdf "github.com/snowbldr/sdfx/vec/v3i"
)

// GridfinityBaseParms configures a Gridfinity base.
type GridfinityBaseParms struct {
	Size   v2i.Vec // size of base in Gridfinity units
	Magnet bool    // add magnet mounts
	Hole   bool    // add mounting holes
}

func (p *GridfinityBaseParms) toSDF() *obj.GfBaseParms {
	return &obj.GfBaseParms{
		Size:   v2isdf.Vec(p.Size),
		Magnet: p.Magnet,
		Hole:   p.Hole,
	}
}

// GridfinityBodyParms configures a Gridfinity body (bin / container).
type GridfinityBodyParms struct {
	Size  v3i.Vec // size of body in Gridfinity units
	Empty bool    // return an empty container
	Hole  bool    // add through holes to the body
}

func (p *GridfinityBodyParms) toSDF() *obj.GfBodyParms {
	return &obj.GfBodyParms{
		Size:  v3isdf.Vec(p.Size),
		Empty: p.Empty,
		Hole:  p.Hole,
	}
}

// GridfinityBase returns a Gridfinity base plate of the given grid size.
func GridfinityBase(p GridfinityBaseParms) *solid.Solid {
	return solid.Wrap(obj.GfBase(p.toSDF()))
}

// GridfinityBody returns a Gridfinity bin / container of the given grid size.
func GridfinityBody(p GridfinityBodyParms) *solid.Solid {
	return solid.Wrap(obj.GfBody(p.toSDF()))
}
