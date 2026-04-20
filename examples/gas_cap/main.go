package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const capRadius = 56.0 / 2.0
const capHeight = 28.0
const capThickness = 4.0
const threadPitch = 6.0
const holeRadius = 0.0 // 33.0 / 2.0

const threadDiameter = 48.5 // just right
const threadRadius = threadDiameter / 2.0

func capOuter() *solid.Solid {
	return obj.KnurledHead3D(capRadius, capHeight, capRadius*0.25)
}

func capInner() *solid.Solid {
	thread := shape.PlasticButtressThread(threadRadius, threadPitch)
	screw := solid.Screw(thread, capHeight, 0, threadPitch, 1)
	return screw.Translate(v3.Z(-capThickness))
}

func capHole() *solid.Solid {
	if holeRadius == 0 {
		return nil
	}
	return solid.Cylinder(capHeight, holeRadius, 0)
}

func gasCap() *solid.Solid {
	inner := capInner()
	if hole := capHole(); hole != nil {
		inner = inner.Union(hole)
	}
	return capOuter().Cut(inner)
}

func main() {
	gasCap().STL("cap.stl", 2.0)
}
