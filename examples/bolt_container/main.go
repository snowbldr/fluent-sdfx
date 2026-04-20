package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const hexRadius = 40.0
const hexHeight = 20.0
const screwRadius = hexRadius * 0.7
const threadPitch = screwRadius / 5.0
const screwLength = 40.0
const tolerance = 0.5

const baseThickness = 4.0

func boltContainer() *solid.Solid {
	// build hex head
	hex := obj.HexHead3D(hexRadius, hexHeight, "tb")

	// build the screw portion
	r := screwRadius - tolerance
	l := screwLength
	thread := shape.ISOThread(r, threadPitch, true)
	screw := solid.Screw(thread, l, 0, threadPitch, 1)

	// chamfer the thread
	screwShape := obj.ChamferedCylinder(screw, 0, 0.25).Translate(v3.Z(l / 2))

	// build the internal cavity
	cavR := screwRadius * 0.75
	cavL := screwLength + hexHeight
	round := screwRadius * 0.1
	ofs := (cavL / 2) - (hexHeight / 2) + baseThickness
	cavity := solid.Cylinder(cavL, cavR, round).Translate(v3.Z(ofs))

	return hex.Union(screwShape).Cut(cavity)
}

func main() {
	boltContainer().STL("container.stl", 2.0)
}
