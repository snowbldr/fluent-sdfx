package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/vec/v2i"
	"github.com/snowbldr/fluent-sdfx/vec/v3i"
)

func main() {
	obj.GridfinityBase(obj.GridfinityBaseParms{
		Size:   v2i.XY(4, 4),
		Magnet: true,
		Hole:   true,
	}).STL("base_4x4.stl", 3.0)

	obj.GridfinityBody(obj.GridfinityBodyParms{
		Size:  v3i.XYZ(1, 1, 3),
		Hole:  true,
		Empty: true,
	}).STL("body_1x1x3.stl", 3.0)

	obj.GridfinityBody(obj.GridfinityBodyParms{
		Size: v3i.XYZ(1, 2, 1),
	}).STL("body_1x2x1.stl", 3.0)
}
