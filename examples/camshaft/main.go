package main

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func camshaft() *solid.Solid {
	// build the shaft from an SoR
	const l0 = 13.0 / 16.0
	const r0 = (5.0 / 16.0) / 2.0
	const l1 = (3.0/32.0)*2.0 + (5.0/16.0)*2.0 + (11.0 / 16.0) + (3.0/16.0)*4.0
	const r1 = (13.0 / 32.0) / 2.0
	const l2 = 1.0 / 2.0
	const r2 = (5.0 / 16.0) / 2.0
	const l3 = 3.0 / 8.0
	r3 := r2 - l3*math.Tan(10.0*math.Pi/180)
	const l4 = 1.0 / 4.0

	p := shape.NewPoly()
	p.Add(0, 0)
	p.Add(r0, 0).Rel()
	p.Add(0, l0).Rel()
	p.Add(r1-r0, 0).Rel()
	p.Add(0, l1).Rel()
	p.Add(r2-r1, 0).Rel()
	p.Add(0, l2).Rel()
	p.Add(r3-r2, l3).Rel()
	p.Add(0, l4).Rel()
	p.Add(-r3, 0).Rel()

	shaft := p.Build().Revolve()

	// make the cams
	const valveDiameter = 0.25
	const rockerRatio = 1.0
	const lift = valveDiameter * rockerRatio * 0.25
	const camDiameter = 5.0 / 8.0
	const camWidth = 3.0 / 16.0
	const k = 1.05
	inletThetaDeg := -110.0

	inlet2d := shape.MakeThreeArcCam(lift, 115.0, camDiameter, k)
	exhaust2d := shape.MakeThreeArcCam(lift, 125.0, camDiameter, k)

	zOfs := (13.0 / 16.0) + (3.0 / 32.0) + (camWidth / 2.0)
	ex4 := exhaust2d.Extrude(camWidth).Translate(v3.Z(zOfs))

	zOfs += (5.0 / 16.0) + camWidth
	in3 := inlet2d.Extrude(camWidth).Translate(v3.Z(zOfs)).RotateZ(inletThetaDeg)

	zOfs += (11.0 / 16.0) + camWidth
	in2 := inlet2d.Extrude(camWidth).Translate(v3.Z(zOfs)).RotateZ(inletThetaDeg + 180)

	zOfs += (5.0 / 16.0) + camWidth
	ex1 := exhaust2d.Extrude(camWidth).Translate(v3.Z(zOfs)).RotateZ(180)

	return solid.UnionAll(shaft, ex1, in2, in3, ex4)
}

func main() {
	camshaft().STL("camshaft.stl", 4.0)
}
