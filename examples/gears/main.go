package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

var module = (5.0 / 8.0) / 20.0
var paDeg = 20.0
var h = 0.15
var numberTeeth = 20

func gear() *solid.Solid {
	return obj.InvoluteGear(obj.InvoluteGearParms{
		NumberTeeth:      numberTeeth,
		Module:           module,
		PressureAngleDeg: paDeg,
		RingWidth:        0.05,
		Facets:            7,
	}).Extrude(h)
}

func rack() *solid.Solid {
	return shape.GearRack(shape.GearRackParams{
		NumberTeeth:      11,
		Module:           module,
		PressureAngleDeg: 20.0,
		BaseHeight:       0.025,
	}).Extrude(h)
}

func main() {
	g := gear()
	r := rack()

	g = g.RotateAxis(v3.XYZ(0, 0, 1), 180.0/float64(numberTeeth)).
		Translate(v3.XYZ(0, 0.39, 0))

	r.Union(g).STL("gear.stl", 2.0)
}
