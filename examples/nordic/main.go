package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

var baseThickness = 3.0
var pillarHeight = 15.0

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%

// nRF52DK
// https://www.nordicsemi.com/Software-and-tools/Development-Kits/nRF52-DK

func nRF52dkStandoffs() *solid.Solid {
	zOfs := 0.5 * (pillarHeight + baseThickness)

	// standoffs with screw holes
	k := obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4,
	}
	positions0 := []v3.Vec{
		v3.XYZ(550.0*units.Mil, 300.0*units.Mil, zOfs),
		v3.XYZ(2600.0*units.Mil, 1600.0*units.Mil, zOfs),
		v3.XYZ(2600.0*units.Mil, 500.0*units.Mil, zOfs),
		v3.XYZ(3800.0*units.Mil, 300.0*units.Mil, zOfs),
	}
	s0 := obj.Standoff3D(k).Multi(positions0...)

	// standoffs with support stubs
	k.HoleDepth = -2.0
	positions1 := []v3.Vec{v3.XYZ(600.0*units.Mil, 2200.0*units.Mil, zOfs)}
	s1 := obj.Standoff3D(k).Multi(positions1...)

	return s0.Union(s1)
}

func nRF52dk() *solid.Solid {
	baseX := 120.0
	baseY := 64.0
	pcbX := 102.0
	pcbY := 63.5

	base2d := obj.Panel2D(obj.PanelParms{
		Size:         v2.XY(baseX, baseY),
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	})

	c1 := shape.Rect(v2.XY(53.0, 35.0), 3.0).Translate(v2.XY(-22.0, 1.0))
	c2 := shape.Rect(v2.XY(20.0, 40.0), 3.0).Translate(v2.XY(37.0, 3.0))

	s2 := base2d.Cut(c1.Union(c2)).Extrude(baseThickness).
		Translate(v3.XY(0.5*pcbX, pcbY-0.5*baseY))

	s3 := nRF52dkStandoffs()
	return solid.SmoothUnion(solid.PolyMin(3.0), s2, s3)
}

// nRF52833DK

func nRF52833dkStandoffs() *solid.Solid {
	zOfs := 0.5 * (pillarHeight + baseThickness)

	k := obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 6.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4,
	}
	positions0 := []v3.Vec{
		v3.XYZ(550.0*units.Mil, 300.0*units.Mil, zOfs),
		v3.XYZ(2600.0*units.Mil, 500.0*units.Mil, zOfs),
		v3.XYZ(2600.0*units.Mil, 1600.0*units.Mil, zOfs),
		v3.XYZ(5050.0*units.Mil, 1825.0*units.Mil, zOfs),
	}
	s0 := obj.Standoff3D(k).Multi(positions0...)

	k.HoleDepth = -2.0
	positions1 := []v3.Vec{
		v3.XYZ(600.0*units.Mil, 2200.0*units.Mil, zOfs),
		v3.XYZ(3550.0*units.Mil, 2200.0*units.Mil, zOfs),
		v3.XYZ(3800.0*units.Mil, 300.0*units.Mil, zOfs),
	}
	s1 := obj.Standoff3D(k).Multi(positions1...)

	return s0.Union(s1)
}

func nRF52833dk() *solid.Solid {
	baseX := 154.0
	baseY := 64.0
	pcbX := 136.53
	pcbY := 63.50

	base2d := obj.Panel2D(obj.PanelParms{
		Size:         v2.XY(baseX, baseY),
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
		HolePattern:  [4]string{"x", "x", "x", "x"},
	})

	c1 := shape.Rect(v2.XY(53.0, 35.0), 3.0).Translate(v2.X(-40.0))
	c2 := shape.Rect(v2.XY(40.0, 35.0), 3.0).Translate(v2.X(32.0))

	s2 := base2d.Cut(c1.Union(c2)).Extrude(baseThickness).
		Translate(v3.XY(0.5*pcbX, pcbY-0.5*baseY))

	s3 := nRF52833dkStandoffs()
	return solid.SmoothUnion(solid.PolyMin(3.0), s2, s3)
}

func main() {
	nRF52dk().ScaleUniform(shrink).STL("nrf52dk.stl", 3.0)
	nRF52833dk().ScaleUniform(shrink).STL("nrf52833dk.stl", 3.0)
}
