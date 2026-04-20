package main

import (
	"log"
	"math"

	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
	"github.com/snowbldr/fluent-sdfx/vec/p2"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// 608 bearing
const bearingOuterOD = 22.0  // outer diameter of outer race
const bearingInnerOD = 12.1  // outer diameter of inner race
const bearingInnerID = 8.0   // inner diameter of inner race
const bearingThickness = 7.0 // bearing thickness

// Adjust clearance to give good interference fits for the bearings and spin caps.
const clearance = 0.0

// Return an N petal bezier flower.
func flower(n int, r0, r1, r2 float64) *shape.Shape {
	theta := units.Tau / float64(n)
	b := shape.NewBezier()

	k0 := v2.X(r1).Add(v2.FromP2(p2.RT(r0, -135*math.Pi/180)))
	k1 := v2.X(r1).Add(v2.FromP2(p2.RT(r0, -45*math.Pi/180)))
	k2 := v2.X(r1).Add(v2.FromP2(p2.RT(r0, 45*math.Pi/180)))
	k3 := v2.X(r1).Add(v2.FromP2(p2.RT(r0, 135*math.Pi/180)))
	k4 := v2.FromP2(p2.RT(r2, theta/2))

	m := shape.Rotate2d(theta * 180 / math.Pi)

	for i := 0; i < n; i++ {
		ofs := float64(i) * theta
		ofsDeg := ofs * 180 / math.Pi
		thetaHalfDeg := theta / 2 * 180 / math.Pi

		b.AddV2(k0).Handle(ofsDeg-45, r0/2, r0/2)
		b.AddV2(k1).Handle(ofsDeg+45, r0/2, r0/2)
		b.AddV2(k2).Handle(ofsDeg+135, r0/2, r0/2)
		b.AddV2(k3).Handle(ofsDeg+225, r0/2, r0/2)
		b.AddV2(k4).Handle(ofsDeg+thetaHalfDeg+90, r2/1.5, r2/1.5)

		k0 = m.MulPosition(k0)
		k1 = m.MulPosition(k1)
		k2 = m.MulPosition(k2)
		k3 = m.MulPosition(k3)
		k4 = m.MulPosition(k4)
	}

	b.Close()
	return b.Build()
}

func body1() (*solid.Solid, error) {
	n := 3
	t := bearingThickness
	r := bearingOuterOD / 2

	r0 := r + 4.0
	r1 := 45.0 - r0
	r2 := r + 4.0

	// body
	f := flower(n, r0, r1, r2)
	s1 := solid.ExtrudeRounded(f, t, t/4.0)

	// periphery holes
	s2 := obj.BoltCircle3D(t, r+clearance, r1, n)
	// center hole
	s3 := solid.Cylinder(t, r+clearance, 0)
	return s1.Cut(s2.Union(s3)), nil
}

func body2() (*solid.Solid, error) {
	t := bearingThickness
	r := bearingOuterOD / 2
	r0 := r + 4.0

	// build the arm
	p := shape.NewPoly()
	p.Add(r, -t/2)
	p.Add(r0, -t/2)
	p.Add(r0, t/2)
	p.Add(r, t/2)
	armShape := shape.Polygon(p.Vertices())
	arm := solid.RevolveAngle(armShape, 270).
		Translate(v3.X(-1.5 * r0))

	// create 6 arms
	arms := arm.RotateUnionZ(6, solid.RotateZMatrix(60))

	// add the center
	body := solid.Cylinder(t, r0, 0)
	body = body.Union(arms)

	// remove the center hole
	hole := solid.Cylinder(t, r, 0)
	return body.Cut(hole), nil
}

// Basic spin cap with variable pin size.
func spincap(pinR, pinL float64) *solid.Solid {
	t := 3.0  // thickness of the spin cap
	st := 1.0 // spacer thickness

	r0 := bearingOuterOD / 2
	r1 := bearingInnerOD / 2

	p := shape.NewPoly()
	p.Add(0, -t-st)
	p.Add(r0, -t-st).Smooth(t/1.5, 6)
	p.Add(r0, -st)
	p.Add(r1, -st)
	p.Add(r1, 0)
	p.Add(pinR, 0)
	p.Add(pinR, pinL)
	p.Add(0, pinL)

	return solid.Revolve(shape.Polygon(p.Vertices()))
}

// Push to fit spincap for single spinner.
func spincapSingle() *solid.Solid {
	gap := 1.0
	r := (bearingInnerID / 2) - clearance
	l := (bearingThickness - gap) / 2
	return spincap(r, l)
}

// Threaded spincap for double spinners.
func spincapDouble(male bool) (*solid.Solid, error) {
	r := (bearingInnerID / 2) - clearance
	threadR := r * 0.8
	threadPitch := 1.0
	threadTolerance := 0.25
	l := bearingThickness

	if male {
		// Add an external screw thread.
		t := shape.ISOThread(threadR-threadTolerance, threadPitch, true)
		screw := solid.Screw(t, bearingThickness, 0, threadPitch, 1)
		screwC := obj.ChamferedCylinder(screw, 0, 0.5).Translate(v3.Z(1.5 * l))
		sc := spincap(r, l+0.5)
		return sc.Union(screwC), nil
	}
	// Add an internal screw thread.
	t := shape.ISOThread(threadR, threadPitch, false)
	screw := solid.Screw(t, bearingThickness, 0, threadPitch, 1).
		Translate(v3.Z(l * 0.5))
	sc := spincap(r, l-0.5)
	return sc.Cut(screw), nil
}

// Inner washer for double spinner.
func spincapWasher() *solid.Solid {
	return obj.Washer3D(obj.WasherParms{
		Thickness:   1.0,
		InnerRadius: (bearingInnerID / 2) * 1.05,
		OuterRadius: (bearingOuterOD + bearingInnerID) / 4,
	})
}

func main() {
	b1, err := body1()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	b1.STL("body1.stl", 3.0)

	b2, err := body2()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	b2.STL("body2.stl", 3.0)

	spincapSingle().STL("cap_single.stl", 1.5)

	scdm, err := spincapDouble(true)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	scdm.STL("cap_double_male.stl", 1.5)

	scdf, err := spincapDouble(false)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	scdf.STL("cap_double_female.stl", 1.5)

	spincapWasher().STL("washer.stl", 1.5)
}
