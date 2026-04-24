package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

var k0 = obj.GenevaParms{
	NumSectors:     6,
	CenterDistance: 50.0,
	DriverRadius:   20.0,
	DrivenRadius:   40.0,
	PinRadius:      2.5,
	Clearance:      0.1,
}

var k1 = obj.GenevaParms{
	NumSectors:     10,
	CenterDistance: 45.0,
	DriverRadius:   12.0,
	DrivenRadius:   45.0,
	PinRadius:      2.0,
	Clearance:      0.1,
}

func main() {
	_ = k1
	k := k0

	sDriver, sDriven := obj.Geneva2D(k)

	wheelHeight := 5.0                 // height of wheels
	holeRadius := 3.25                 // radius of center hole
	hubRadius := 10.0                  // hub radius for driven wheel
	baseRadius := 1.5 * k.DriverRadius // radius of base for driver wheel

	// extrude the driver wheel, add base, drill center hole
	driver := sDriver.Extrude(wheelHeight).
		Translate(v3.Z(wheelHeight / 2))
	base := solid.Cylinder(wheelHeight, baseRadius, 0).
		Translate(v3.Z(-wheelHeight / 2))
	driver = driver.Union(base)
	hole := solid.Cylinder(2*wheelHeight, holeRadius, 0)
	driver = driver.Cut(hole)

	// extrude the driven wheel, add hub, drill center hole
	driven := sDriven.Extrude(wheelHeight).
		Translate(v3.Z(-wheelHeight / 2))
	hub := solid.Cylinder(wheelHeight, hubRadius, 0).
		Translate(v3.Z(wheelHeight / 2))
	driven = driven.Union(hub).Cut(hole)

	const cellsPerMM = 3.0
	driver.STL("driver.stl", cellsPerMM)
	driven.STL("driven.stl", cellsPerMM)

	driver = driver.Translate(v3.X(-0.8 * k.DrivenRadius))
	driven = driven.Translate(v3.X(k.DrivenRadius))
	driver.Union(driven).STL("geneva.stl", cellsPerMM)
}
