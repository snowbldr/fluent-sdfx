// Parametric helpers: a 2x1 Gridfinity base.
//
// Size is in Gridfinity units (42mm per cell). Magnet/Hole flags add the
// standard 6mm magnet pockets and M3 mounting holes.
package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	v2i "github.com/snowbldr/fluent-sdfx/vec/v2i"
)

func main() {
	obj.GfBase(obj.GfBaseParms{
		Size:   v2i.XY(2, 1),
		Magnet: true,
		Hole:   true,
	}).STL("out.stl", 5.0)
}
