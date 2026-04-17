package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/vec/v2i"
	"github.com/snowbldr/fluent-sdfx/vec/v3i"
)

func main() {
	kBase := obj.GfBaseParms{
		Size:   v2i.XY(4, 4),
		Magnet: true,
		Hole:   true,
	}
	obj.GfBase(kBase).ToSTL("base_4x4.stl", 300)

	kBody := obj.GfBodyParms{
		Size:  v3i.XYZ(1, 1, 3),
		Hole:  true,
		Empty: true,
	}
	obj.GfBody(kBody).ToSTL("body_1x1x3.stl", 300)

	kBody2 := obj.GfBodyParms{
		Size: v3i.XYZ(1, 2, 1),
	}
	obj.GfBody(kBody2).ToSTL("body_1x2x1.stl", 300)
}
