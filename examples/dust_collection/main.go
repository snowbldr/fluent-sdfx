package main

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
)

// dust deputy tapered pipe
const ddOuterDiameter = 51.0
const ddLength = 39.0

var ddTaper = units.DtoR(2.0)

// vacuum hose 2.5" male fitting
const vhOuterDiameter = 58.0
const vhClearance = 0.6

var vhTaper = units.DtoR(0.4)

const wallThickness = 4.0

// dust deputy (female), 2.5" vacuum (female)
func dustDeputyToVacuumFF() *solid.Solid {
	const t = wallThickness
	const transitionLength = 15.0
	const vhLength = 30.0

	r0 := ddOuterDiameter * 0.5
	r1 := r0 - ddLength*math.Tan(ddTaper)
	r3 := (vhOuterDiameter + vhClearance) * 0.5
	r2 := r3 - (vhLength * math.Tan(vhTaper))

	h0 := 0.0
	h1 := h0 + ddLength
	h2 := h1 + transitionLength
	h3 := h2 + vhLength

	p := shape.NewPoly()
	p.Add(r0+t, h0)
	p.Add(r1+t, h1).Smooth(t, 4)
	p.Add(r2+t, h2).Smooth(t, 4)
	p.Add(r3+t, h3)
	p.Add(r3, h3)
	p.Add(r2, h2).Smooth(t, 4)
	p.Add(r1, h1).Smooth(t, 4)
	p.Add(r0, h0)

	return solid.Revolve(p.Build())
}

// 2.5" vacuum (male) to pipe (male)
func vacuumToPipeMM(name string) *solid.Solid {
	k := obj.PipeLookup(name, "mm")

	t := wallThickness
	transitionLength := 15.0

	r0 := k.Outer
	r1 := vhOuterDiameter * 0.5

	h0 := 0.0
	h1 := h0 + 35.0
	h2 := h1 + transitionLength
	h3 := h2 + 20.0

	p := shape.NewPoly()
	p.Add(r0, h0)
	p.Add(r0, h1).Smooth(t, 4)
	p.Add(r1, h2).Smooth(t, 4)
	p.Add(r1, h3)
	p.Add(r1-t, h3)
	p.Add(r1-t, h2).Smooth(t, 4)
	p.Add(r0-t, h1).Smooth(t, 4)
	p.Add(r0-t, h0)

	return solid.Revolve(p.Build())
}

// dust deputy (female) to pipe (male)
func dustDeputyToPipeFM(name string) *solid.Solid {
	k := obj.PipeLookup(name, "mm")

	t := wallThickness
	transitionLength := 15.0

	r0 := k.Outer
	r2 := (ddOuterDiameter * 0.5) + t
	r1 := r2 - ddLength*math.Tan(ddTaper)

	h0 := 0.0
	h1 := h0 + 35.0
	h2 := h1 + transitionLength
	h3 := h2 + ddLength

	p := shape.NewPoly()
	p.Add(r0, h0)
	p.Add(r0, h1).Smooth(t, 4)
	p.Add(r1, h2).Smooth(t, 4)
	p.Add(r2, h3)
	p.Add(r2-t, h3)
	p.Add(r1-t, h2).Smooth(t, 4)
	p.Add(r0-t, h1).Smooth(t, 4)
	p.Add(r0-t, h0)

	return solid.Revolve(p.Build())
}

func main() {
	dustDeputyToVacuumFF().ToSTL("fdd_fvh25.stl", 150)
	vacuumToPipeMM("sch40:2").ToSTL("mvh25_mpvc.stl", 150)
	dustDeputyToPipeFM("sch40:2").ToSTL("fdd_mpvc.stl", 150)
}
