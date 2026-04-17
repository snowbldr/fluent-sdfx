package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

var module = (5.0 / 8.0) / 20.0
var pa = units.DtoR(20.0)
var h = 0.15
var numberTeeth = 20

func gear() *solid.Solid {
	k := obj.InvoluteGearParms{
		NumberTeeth:   numberTeeth,
		Module:        module,
		PressureAngle: pa,
		RingWidth:     0.05,
		Facets:        7,
	}
	return solid.Extrude(obj.InvoluteGear(k), h)
}

func rack() *solid.Solid {
	k := shape.GearRackParams{
		NumberTeeth:      11,
		Module:           module,
		PressureAngleDeg: 20.0,
		BaseHeight:       0.025,
	}
	rack2d := shape.GearRack(k)
	return solid.Extrude(rack2d, h)
}

func main() {
	g := gear()
	r := rack()

	g = g.RotateAxis(v3.XYZ(0, 0, 1), 180.0/float64(numberTeeth)).
		Translate(v3.XYZ(0, 0.39, 0))

	r.Union(g).ToSTL("gear.stl", 200)
}
