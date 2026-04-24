package main

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const panelHoleDiameter = 4.0
const baseHoleDiameter = 4.0

const panelThickness = 3.0
const panelWidth = 190
const panelHeight = 75

const mountWidth = 16.0

const holeMargin = 0.5 * (mountWidth + panelThickness)

// pam8302: 2.5W Class D Audio Amplifier
func pam8302(thickness float64) *solid.Solid {
	const pillarHeight = 4.5

	xOfs := 0.4 * units.MillimetresPerInch * 0.5
	zOfs := 0.5 * (thickness + pillarHeight)
	return obj.Standoff3D(obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 4.5,
		HoleDepth:      pillarHeight,
		HoleDiameter:   2.0,
	}).Multi([]v3.Vec{v3.XYZ(xOfs, 0, zOfs), v3.XYZ(-xOfs, 0, zOfs)})
}

func display0(thickness float64, negative bool) *solid.Solid {
	return obj.Display(obj.DisplayParms{
		Window:          v2.XY(60, 45),
		Rounding:        2.0,
		Supports:        v2.XY(76.08, 44.0),
		SupportHeight:   4.0,
		SupportDiameter: 5.0,
		HoleDiameter:    3.0,
		Offset:          v2.XY(2.5, 0),
		Thickness:       thickness,
		Countersunk:     true,
	}, negative)
}

func display1(thickness float64, negative bool) *solid.Solid {
	return obj.Display(obj.DisplayParms{
		Window:          v2.XY(26, 14),
		Rounding:        1,
		Supports:        v2.XY(23.5, 23.8),
		SupportHeight:   2.1,
		SupportDiameter: 4.5,
		HoleDiameter:    2.5,
		Offset:          v2.XY(0, -2.0),
		Thickness:       thickness,
		Countersunk:     true,
	}, negative)
}

func speakerGrille(thickness float64, negative bool) *solid.Solid {
	const grilleRadius = 77.5 * 0.5

	if negative {
		return obj.CircleGrille3D(obj.CircleGrilleParms{
			HoleDiameter:      4.0,
			GrilleDiameter:    2.0 * grilleRadius,
			RadialSpacing:     0.5,
			TangentialSpacing: 0.5,
			Thickness:         thickness,
		})
	}

	return obj.Washer3D(obj.WasherParms{
		Thickness:   thickness,
		InnerRadius: grilleRadius,
		OuterRadius: grilleRadius + thickness,
	}).Translate(v3.XYZ(0, 0, thickness))
}

// pcbMount0 mounts the adafruit half breadboard with the rpi-pico.
func pcbMount0() *solid.Solid {
	const width = 60.0
	const length = 90.0
	const margin = 5.0
	const height = 10.0
	const thickness = 3.0

	panel := obj.Panel3D(obj.PanelParms{
		Size:         v2.XY(width, length),
		CornerRadius: margin,
		HoleDiameter: 3.8,
		HoleMargin:   [4]float64{margin, margin, margin, margin},
		HolePattern:  [4]string{"x", "x", "x", "x"},
		Thickness:    thickness,
		Ridge:        v2.XY(width-3.5*margin, length-3.5*margin),
	})

	zOfs := 0.5 * (height + thickness)
	yOfs := 0.5 * 2.9 * units.MillimetresPerInch
	standoffs := obj.Standoff3D(obj.StandoffParms{
		PillarHeight:   height,
		PillarDiameter: 8,
		HoleDepth:      height,
		HoleDiameter:   2.4,
	}).Multi([]v3.Vec{v3.XYZ(0, -yOfs, zOfs), v3.XYZ(0, yOfs, zOfs)})

	return panel.Union(standoffs)
}

// pcbMount1 mounts a pcb with the SDR frontend.
func pcbMount1() *solid.Solid {
	const width = 70.0
	const length = 110.0
	const margin = 5.0
	const height = 10.0
	const thickness = 3.0

	panel := obj.Panel3D(obj.PanelParms{
		Size:         v2.XY(width, length),
		CornerRadius: margin,
		HoleDiameter: 3.8,
		HoleMargin:   [4]float64{margin, margin, margin, margin},
		HolePattern:  [4]string{"x", "x", "x", "x"},
		Thickness:    thickness,
		Ridge:        v2.XY(width-3.5*margin, length-3.5*margin),
	})

	zOfs := 0.5 * (height + thickness)
	yOfs := 0.5 * 3.4 * units.MillimetresPerInch
	standoffs := obj.Standoff3D(obj.StandoffParms{
		PillarHeight:   height,
		PillarDiameter: 8,
		HoleDepth:      height,
		HoleDiameter:   2.4,
	}).Multi([]v3.Vec{v3.XYZ(0, -yOfs, zOfs), v3.XYZ(0, yOfs, zOfs)})

	return panel.Union(standoffs)
}

func picoRxBezel(thickness float64) *solid.Solid {
	const ridgeWidth = panelWidth - (2.0 * mountWidth) - 1.0

	panel := obj.Panel3D(obj.PanelParms{
		Size:         v2.XY(panelWidth, panelHeight),
		CornerRadius: 5.0,
		HoleDiameter: panelHoleDiameter,
		HoleMargin:   [4]float64{holeMargin, holeMargin, holeMargin, holeMargin},
		HolePattern:  [4]string{"x", "x", "x", "x"},
		Thickness:    thickness,
		Ridge:        v2.XY(ridgeWidth, 0),
	})

	re := obj.KeyedHole3D(obj.KeyedHoleParms{
		Diameter:  9.2,
		KeySize:   0.9,
		NumKeys:   2,
		Thickness: thickness,
	})

	pb := solid.Box(v3.XYZ(13.2, 10.8, thickness), 0)
	xOfs := 22.0
	pb0 := pb.Translate(v3.XYZ(xOfs, 0, 0))
	pb1 := pb.Translate(v3.XYZ(-xOfs, 0, 0))

	d1n := display1(thickness, true)
	d1p := display1(thickness, false)
	yOfs := 13.0
	xOfs = 47.0
	d1n = d1n.Translate(v3.XYZ(xOfs, yOfs, 0))
	d1p = d1p.Translate(v3.XYZ(xOfs, yOfs, 0))

	yOfs = -17.0
	input := re.Union(pb0, pb1).Translate(v3.XYZ(xOfs, yOfs, 0))

	d0n := display0(thickness, true)
	d0p := display0(thickness, false)
	yOfs = 0.0
	xOfs = -35.0
	d0n = d0n.Translate(v3.XYZ(xOfs, yOfs, 0))
	d0p = d0p.Translate(v3.XYZ(xOfs, yOfs, 0))

	return panel.Union(d0p, d1p).Cut(input, d0n, d1n)
}

func twoHoles(thickness, diameter, distance float64) *solid.Solid {
	h := solid.Cylinder(thickness, 0.5*diameter, 0)
	xOfs := 0.5 * distance
	return h.Translate(v3.XYZ(xOfs, 0, 0)).Union(h.Translate(v3.XYZ(-xOfs, 0, 0)))
}

func sideMount(thickness float64, lhs bool) *solid.Solid {
	const mHeight = 95.0
	const mLength = 125.0
	const mSlope = 75.0
	const mRound0 = 2.0

	mRound2 := mRound0 + thickness

	d := mSlope / math.Sqrt(2)
	bh := mHeight - 2.0*mRound2
	bl := mLength - 2.0*mRound2
	bw := 2.0 * (mountWidth - mRound2)

	p := shape.NewPoly()
	p.Add(-bh*0.5, bh*0.5)
	p.Add(0, -bh).Rel()
	p.Add(bl, 0).Rel()
	p.Add(0, bh-d).Rel()
	p.Add(-d, d).Rel()
	b0 := p.Build()
	box := b0.Extrude(bw)

	inner := box.Offset(mRound0)
	outer := box.Offset(mRound2)
	s := outer.Cut(inner)

	s = s.CutPlane(v3.XYZ(0, 0, 0), v3.XYZ(0, 0, -1))
	s = s.Translate(v3.XYZ(0, 0, mountWidth-0.5*thickness))

	// base holes
	baseHoles := twoHoles(thickness, baseHoleDiameter, 0.7*mLength).RotateX(-90)
	xOfs := 0.5 * (mLength - mHeight)
	yOfs := 0.5 * (mHeight - thickness)
	zOfs := 0.5 * mountWidth
	baseHoles = baseHoles.Translate(v3.XYZ(xOfs, -yOfs, zOfs))

	// panel holes
	const panelHoleDistance = panelHeight - 2.0*holeMargin
	panelHoles := twoHoles(thickness, panelHoleDiameter, panelHoleDistance).
		RotateX(-90).
		RotateZ(-45)
	delta := (mRound0 + 0.5*thickness) / math.Sqrt(2)
	xOfs = bl - 0.5*(bh+d) + delta
	yOfs = 0.5*(bh-d) + delta
	zOfs = 0.5 * mountWidth
	panelHoles = panelHoles.Translate(v3.XYZ(xOfs, yOfs, zOfs))

	s = s.Cut(baseHoles.Union(panelHoles))

	if lhs {
		s = s.MirrorXZ()
	}

	return s
}

func rhsMount(thickness float64) *solid.Solid {
	rhs := sideMount(thickness, false)
	sn := speakerGrille(thickness, true)
	sp := speakerGrille(thickness, false)
	amp := pam8302(thickness).Translate(v3.XYZ(55, -10, 0))
	return rhs.Union(sp, amp).Cut(sn)
}

func lhsMount(thickness float64) *solid.Solid {
	return sideMount(thickness, true)
}

func main() {
	pcbMount0().STL("pcb_mount0.stl", 5.0)
	pcbMount1().STL("pcb_mount1.stl", 5.0)
	picoRxBezel(panelThickness).STL("bezel.stl", 5.0)
	rhsMount(panelThickness).STL("rhs.stl", 5.0)
	lhsMount(panelThickness).STL("lhs.stl", 5.0)
}
