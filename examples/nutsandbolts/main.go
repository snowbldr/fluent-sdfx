package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func nutAndBolt(name string, totalLength, shankLength float64) *solid.Solid {
	bolt := obj.Bolt(obj.BoltParms{
		Thread:      name,
		Style:       "hex",
		TotalLength: totalLength,
		ShankLength: shankLength,
	})

	nut := obj.Nut(obj.NutParms{
		Thread: name,
		Style:  "hex",
	}).Translate(v3.XYZ(0, 0, totalLength*1.5))

	return nut.Union(bolt)
}

func main() {
	xOffset := 1.5

	s0 := nutAndBolt("unc_1/4", 2, 0.5).Translate(v3.XYZ(-0.6*xOffset, 0, 0))
	s1 := nutAndBolt("unc_1/2", 2.0, 0.5)
	s2 := nutAndBolt("unc_1", 2.0, 0.5).Translate(v3.XYZ(xOffset, 0, 0))

	solid.UnionAll(s0, s1, s2).STL("nutandbolt.stl", 4.0)
}
