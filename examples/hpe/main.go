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

func ap723hSupport() *solid.Solid {
	const w0 = 20
	const l0 = 60
	const h0 = 6

	b0 := solid.Box(v3.XYZ(w0, l0, h0), 0)

	const h1 = 3.7
	const l1 = 47
	b1 := solid.Box(v3.XYZ(w0, l1, h1), 0).
		Translate(v3.YZ(0.5*(l1-l0), 0.5*(h1-h0)))

	hole := solid.Cylinder(h0-h1, 1.2, 0).
		Translate(v3.XYZ(0.5*w0-3.0, l1-0.5*l0, 0.5*h1))

	return b0.Cut(b1.Union(hole))
}

func ap723hStandoffs() *solid.Solid {
	s := obj.Standoff3D(obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 10.0,
		HoleDepth:      10.0,
		HoleDiameter:   4.0,
	}).Union(obj.Standoff3D(obj.StandoffParms{
		PillarHeight:   pillarHeight + 2.0,
		PillarDiameter: 5.5,
		HoleDepth:      10.0,
		HoleDiameter:   2.4,
	}))

	zOfs := 0.5 * (pillarHeight + baseThickness)
	return s.Multi(v3.XYZ(0, 0, zOfs), v3.XYZ(103.0, 0, zOfs), v3.XYZ(103.0, 152.0, zOfs), v3.XYZ(0, 152.0, zOfs))
}

func ap723hMount() *solid.Solid {
	pcbX := 102.5
	pcbY := 152.0
	baseX := pcbX + 20.0
	baseY := pcbY + 20.0

	base2d := obj.Panel2D(obj.PanelParms{
		Size:         v2.XY(baseX, baseY),
		CornerRadius: 5.0,
	})

	c1 := shape.Rect(v2.XY(baseX-35, baseY-35), 5.0)

	s2 := base2d.Cut(c1).Extrude(baseThickness).
		Translate(v3.XY(0.5*pcbX, 0.5*pcbY))

	// reinforcing ribs
	const ribHeight = 5.0
	zOfs := 0.5 * (ribHeight + baseThickness)
	r0 := solid.Box(v3.XYZ(3.0, 0.9*pcbY, ribHeight), 0).
		Translate(v3.YZ(0.5*pcbY, zOfs))
	r1 := r0.Translate(v3.X(pcbX))
	r2 := solid.Box(v3.XYZ(0.8*pcbX, 3.0, ribHeight), 0).
		Translate(v3.XZ(0.5*pcbX, zOfs))
	r3 := r2.Translate(v3.Y(pcbY))

	s2 = s2.Union(r0, r1, r2, r3)

	s3 := ap723hStandoffs()
	return solid.SmoothUnion(solid.PolyMin(3.0), s2, s3)
}

func ap725Standoffs() *solid.Solid {
	zOfs := 0.5 * (pillarHeight + baseThickness)
	return obj.Standoff3D(obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 138.0 * units.Mil * 2.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4,
	}).Multi(v3.XYZ(0, 0, zOfs),
		v3.XYZ(5984.255*units.Mil, 0, zOfs),
		v3.XYZ(5551.185*units.Mil, 4704.72*units.Mil, zOfs),
		v3.XYZ(433.071*units.Mil, 4704.72*units.Mil, zOfs),
		v3.XYZ(2700.795*units.Mil, 5389.76*units.Mil, zOfs),
		v3.XYZ(3714.565*units.Mil, 1708.66*units.Mil, zOfs))
}

func ap725Mount() *solid.Solid {
	baseX := 165.0
	baseY := 150.0
	pcbX := 5984.255 * units.Mil
	pcbY := 5389.76 * units.Mil

	base2d := obj.Panel2D(obj.PanelParms{
		Size:         v2.XY(baseX, baseY),
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
	})

	c1 := shape.Rect(v2.XY(100, 65.0), 3.0).Translate(v2.Y(20))
	c2 := shape.Rect(v2.XY(135, 25.0), 3.0).Translate(v2.Y(-50))

	s2 := base2d.Cut(c1.Union(c2)).Extrude(baseThickness).
		Translate(v3.XY(0.5*pcbX, 0.5*pcbY))

	const ribHeight = 5.0
	zOfs := 0.5 * (ribHeight + baseThickness)
	r0 := solid.Box(v3.XYZ(3.0, 0.9*pcbY, ribHeight), 0).
		Translate(v3.YZ(0.5*pcbY, zOfs))
	r1 := r0.Translate(v3.X(pcbX))
	r2 := solid.Box(v3.XYZ(0.9*pcbX, 3.0, ribHeight), 0).
		Translate(v3.XZ(0.5*pcbX, zOfs))

	s2 = s2.Union(r0, r1, r2)
	s3 := ap725Standoffs()
	return solid.SmoothUnion(solid.PolyMin(3.0), s2, s3)
}

const holeSquare = 6102.36 * units.Mil

func ap745Standoffs() *solid.Solid {
	zOfs := 0.5 * (pillarHeight + baseThickness)
	return obj.Standoff3D(obj.StandoffParms{
		PillarHeight:   pillarHeight,
		PillarDiameter: 138.0 * units.Mil * 2.0,
		HoleDepth:      10.0,
		HoleDiameter:   2.4,
	}).Multi(v3.XYZ(0, holeSquare, zOfs),
		v3.XYZ(0, 0, zOfs),
		v3.XYZ(holeSquare, holeSquare, zOfs),
		v3.XYZ(holeSquare, 0, zOfs),
		v3.XYZ(3937.01*units.Mil, 7047.24*units.Mil, zOfs),
		v3.XYZ(1240.16*units.Mil, 5570.87*units.Mil, zOfs),
		v3.XYZ(2648.46*units.Mil, 3485.15*units.Mil, zOfs),
		v3.XYZ(3693.46*units.Mil, 610.15*units.Mil, zOfs))
}

func ap745Mount() *solid.Solid {
	baseX := 170.0
	baseY := 190.0
	pcbX := 6102.36 * units.Mil
	pcbY := 7047.24 * units.Mil

	base2d := obj.Panel2D(obj.PanelParms{
		Size:         v2.XY(baseX, baseY),
		CornerRadius: 5.0,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5.0, 5.0, 5.0, 5.0},
	})

	c1 := shape.Rect(v2.XY(140, 50.0), 3.0).Translate(v2.Y(-37))
	c2 := shape.Rect(v2.XY(90.0, 50.0), 3.0).Translate(v2.XY(15, 45))

	s2 := base2d.Cut(c1.Union(c2)).Extrude(baseThickness).
		Translate(v3.XY(0.5*pcbX, 0.5*pcbY))

	const ribHeight = 5.0
	zOfs := 0.5 * (ribHeight + baseThickness)
	r0 := solid.Box(v3.XYZ(3.0, 0.75*pcbY, ribHeight), 0).
		Translate(v3.YZ(0.5*pcbY-12.0, zOfs))
	r1 := r0.Translate(v3.X(holeSquare))
	s2 = s2.Union(r0, r1)

	s3 := ap745Standoffs()
	return solid.SmoothUnion(solid.PolyMin(3.0), s2, s3)
}

func main() {
	ap725Mount().STL("ap725.stl", 5.0)
	ap745Mount().STL("ap745.stl", 5.0)
	ap723hMount().STL("ap723h.stl", 5.0)
	ap723hSupport().STL("ap723h_support.stl", 5.0)
}
