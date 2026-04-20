package main

import (
	"log"
	"math"

	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func test1() error {
	s0 := shape.Rect(v2.XY(0.8, 1.2), 0.05)
	s1 := solid.RevolveAngle(s0, 225)
	s1.STL("test1.stl", 2.0)
	return nil
}

func test2() error {
	s0 := shape.Rect(v2.XY(0.8, 1.2), 0.1)
	s1 := solid.Extrude(s0, 0.3)
	s1.STL("test2.stl", 2.0)
	return nil
}

func test3() error {
	s0 := shape.Circle(0.1).Translate(v2.X(1))
	s1 := solid.Revolve(s0)
	s1.STL("test3.stl", 2.0)
	return nil
}

func test4() error {
	s0 := shape.Rect(v2.XY(0.2, 0.4), 0.05).Translate(v2.X(1))
	s1 := solid.RevolveAngle(s0, 270)
	s1.STL("test4.stl", 2.0)
	return nil
}

func test5() error {
	s0 := shape.Rect(v2.XY(0.2, 0.4), 0.05).
		Rotate(45).
		Translate(v2.X(1))
	s1 := solid.RevolveAngle(s0, 315)
	s1.STL("test5.stl", 2.0)
	return nil
}

func test6() error {
	s0 := solid.Sphere(0.5)
	d := 0.4
	s1 := s0.Translate(v3.Y(d))
	s2 := s0.Translate(v3.Y(-d))
	s3 := solid.SmoothUnion(solid.PolyMin(0.1), s1, s2)
	s3.STL("test6.stl", 2.0)
	return nil
}

func test7() error {
	s0 := solid.Box(v3.XYZ(0.8, 0.8, 0.05), 0)
	s1 := s0.RotateAxis(v3.X(1), 60)
	s2 := solid.SmoothUnion(solid.PolyMin(0.1), s0, s1)
	s3 := s2.RotateAxis(v3.Z(1), -30)
	s3.STL("test7.stl", 2.0)
	return nil
}

func test9() error {
	solid.Sphere(10.0).STL("test9.stl", 2.0)
	return nil
}

func test10() error {
	s0 := solid.Box(v3.XYZ(0.8, 0.8, 0.05), 0)
	s1 := s0.RotateAxis(v3.X(1), 60)
	s := solid.SmoothUnion(solid.PolyMin(0.1), s0, s1)
	s.STL("test10.stl", 2.0)
	return nil
}

func test11() error {
	solid.Capsule(3.0, 1.4).STL("test11.stl", 2.0)
	return nil
}

func test12() error {
	k := 0.1
	points := []v2.Vec{v2.XY(0, -k), v2.XY(k, k), v2.XY(-k, k)}
	s0 := shape.Polygon(points).Translate(v2.X(0.8))
	s1 := solid.RevolveAngle(s0, 360)
	s1.STL("test12.stl", 2.0)
	return nil
}

func test13() error {
	k := 0.4
	s0 := shape.Polygon([]v2.Vec{v2.XY(k, -k), v2.XY(k, k), v2.XY(-k, k), v2.XY(-k, -k)}).
		Translate(v2.X(0.8))
	s1 := solid.RevolveAngle(s0, 270)
	s1.STL("test13.stl", 2.0)
	return nil
}

func test14() error {
	// size
	a := 0.3
	b := 0.7
	// rotation
	theta := 30.0
	c := math.Cos(theta * math.Pi / 180)
	si := math.Sin(theta * math.Pi / 180)
	// translate
	j := 10.0
	k := 2.0

	points := []v2.Vec{v2.XY(j+c*a-si*b, k+si*a+c*b), v2.XY(j-c*a-si*b, k-si*a+c*b), v2.XY(j-c*a+si*b, k-si*a-c*b), v2.XY(j+c*a+si*b, k+si*a-c*b)}
	s0 := shape.Polygon(points)
	s1 := solid.RevolveAngle(s0, 300)
	s1.STL("test14.stl", 2.0)
	return nil
}

func test15() error {
	a := 1.0
	b := 1.0
	theta := 0.0
	j := 3.0
	k := 0.0

	points := []v2.Vec{v2.XY(0, -b), v2.XY(a, b), v2.XY(-a, b)}
	s0 := shape.Polygon(points).Rotate(theta).Translate(v2.XY(j, k))
	s1 := solid.RevolveAngle(s0, 300).RotateAxis(v3.Z(1), 30)
	s1.STL("test15.stl", 2.0)
	return nil
}

func test16() error {
	a0 := 1.3
	b0 := 0.4
	a1 := 1.3
	b1 := 1.3
	c := 0.8
	theta := 20.0
	j := 4.0
	k := 0.0

	points := []v2.Vec{v2.XY(b0, -c), v2.XY(a0, c), v2.XY(-a1, c), v2.XY(-b1, -c)}
	s0 := shape.Polygon(points).Rotate(theta).Translate(v2.XY(j, k))
	s1 := solid.RevolveAngle(s0, 300).RotateAxis(v3.Z(1), 30)
	s1.STL("test16.stl", 2.0)
	return nil
}

func test17() error {
	a := 1.3
	b := 0.4
	j := 3.0
	k := 0.0

	points := []v2.Vec{v2.XY(a, 0), v2.XY(-a, b), v2.XY(-a, -b)}
	s0 := shape.Polygon(points).Translate(v2.XY(j, k))
	s1 := solid.RevolveAngle(s0, 300).RotateAxis(v3.Z(1), 30)
	s1.STL("test17.stl", 2.0)
	return nil
}

func test18() error {
	r0 := 10.0
	r1 := 8.0
	r2 := 7.5
	r3 := 9.0

	h0 := 4.0
	h1 := 6.0
	h2 := 5.5
	h3 := 3.5
	h4 := 1.0

	points := []v2.Vec{v2.XY(0, 0), v2.XY(r0, 0), v2.XY(r0, h0), v2.XY(r1, h1), v2.XY(r2, h2), v2.XY(r3, h3), v2.XY(r3, h4), v2.XY(0, h4)}
	s0 := shape.Polygon(points)
	s1 := solid.RevolveAngle(s0, 300).RotateAxis(v3.Z(1), 30)
	s1.STL("test18.stl", 2.0)
	return nil
}

func test19() error {
	r := 2.0
	k := 1.9
	s0 := shape.Circle(r)
	s := s0.SmoothArray(3, 7, v2.XY(k*r, k*r), solid.PolyMin(0.8))
	s2 := solid.Extrude(s, 1.0)
	s2.STL("test19.stl", 2.0)
	return nil
}

func test20() error {
	r := 4.0
	d := 20.0
	s0 := shape.Circle(r).Translate(v2.X(d))
	ru := s0.SmoothRotateUnion(5, shape.Rotate2d(20), solid.PolyMin(1.2))
	s1 := solid.Extrude(ru, 10.0)
	s1.STL("test20.stl", 2.0)
	return nil
}

func test21() error {
	r := 2.0
	k := 1.9
	s0 := solid.Sphere(r)
	s := s0.SmoothArray(3, 7, 5, v3.XYZ(k*r, k*r, k*r), solid.PolyMin(0.8))
	s.STL("test21.stl", 2.0)
	return nil
}

func test22() error {
	r := 4.0
	d := 20.0
	s0 := solid.Sphere(r).Translate(v3.X(d))
	ru := s0.SmoothRotateUnionZ(5, solid.Rotate3dMatrix(v3.Z(1), 20), solid.PolyMin(1.2))
	ru.STL("test22.stl", 2.0)
	return nil
}

func test26() error {
	solid.Cylinder(5, 2, 1).STL("test26.stl", 2.0)
	return nil
}

func test27() error {
	r := 5.0
	posn := v3.VecSet{v3.XYZ(2*r, 2*r, 0), v3.XYZ(-r, r, 0), v3.XYZ(r, -r, 0), v3.XYZ(-r, -r, 0), v3.XYZ(0, 0, 0)}
	cyl := solid.Cylinder(3, 1, 0)
	s := cyl.Multi(posn)
	s.STL("test27.stl", 2.0)
	return nil
}

func test28() error {
	solid.Cone(20, 12, 8, 2).STL("test28.stl", 2.0)
	return nil
}

func test29() error {
	s0 := shape.Line(10, 3)
	s1 := solid.Extrude(s0, 4)
	s1.STL("test29.stl", 2.0)
	return nil
}

func test30() error {
	s0 := shape.Line(10, 3).CutLine(v2.X(4), v2.XY(1, 1))
	s1 := solid.Extrude(s0, 4)
	s1.STL("test30.stl", 2.0)
	return nil
}

func test31() error {
	s := obj.CounterSunkHole3D(30, 2)
	s.STL("test31.stl", 2.0)
	return nil
}

func test32() error {
	s0 := shape.MakeFlatFlankCam(0.094, 2.0*57.5, 0.625)
	s1 := solid.Extrude(s0, 0.1)
	s1.STL("cam0.stl", 2.0)
	return nil
}

func test33() error {
	s0 := shape.ThreeArcCam(30, 20, 5, 50000)
	s1 := solid.Extrude(s0, 4)
	s1.STL("cam1.stl", 2.0)
	return nil
}

func test34() error {
	s0 := shape.MakeThreeArcCam(0.1, 2.0*80, 0.7, 1.1)
	s1 := solid.Extrude(s0, 0.1)
	s1.STL("cam2.stl", 2.0)
	return nil
}

func test35() error {
	r := 7.0
	d := 20.0
	s0 := shape.Line(r, 1.0).Translate(v2.X(d)).RotateCopy(15)
	s1 := solid.Extrude(s0, 10.0)
	s1.STL("rotate_copy.stl", 2.0)
	return nil
}

func test36() error {
	sDriver, sDriven := obj.Geneva2D(obj.GenevaParms{
		NumSectors:     6,
		CenterDistance: 100,
		DriverRadius:   40,
		DrivenRadius:   80,
		PinRadius:      5,
		Clearance:      0.5,
	})
	solid.Extrude(sDriver, 10).STL("driver.stl", 2.0)
	solid.Extrude(sDriven, 10).STL("driven.stl", 2.0)
	return nil
}

func test37() error {
	r := 5.0
	p := 2.0
	isoThread := shape.ISOThread(r, p, true)
	s := solid.Screw(isoThread, 50, 0, p, 1)
	s.STL("screw.stl", 4.0)
	return nil
}

func test39() error {
	s0 := shape.Flange1(30, 20, 10)
	s1 := solid.Extrude(s0, 5)
	s1.STL("flange.stl", 2.0)
	return nil
}

func test40() error {
	d := 30.0
	wall := 5.0
	s0 := solid.Box(v3.XYZ(d, d, d), wall/2)
	s1 := solid.Box(v3.XYZ(d-wall, d-wall, d), wall/2).
		Translate(v3.Z(wall / 2))
	s := solid.SmoothDifference(solid.PolyMax(2), s0, s1)
	s.STL("rounded_box.stl", 2.0)
	return nil
}

func test41() error {
	s0 := solid.Cylinder(20.0, 5.0, 0)
	s1 := solid.Slice(s0, v3.Zero, v3.YZ(1, 1))
	s2 := solid.Revolve(s1)
	s2.STL("ellipsoid_egg.stl", 2.0)
	return nil
}

func test43() error {
	s0 := shape.Line(10, 3).CutLine(v2.X(4), v2.XY(1, 1))
	s1 := solid.ExtrudeRounded(s0, 4, 1)
	s1.STL("cut2d.stl", 3.0)
	return nil
}

func test44() error {
	r := 100.0
	s0 := shape.Polygon(shape.Nagon(5, r))
	s1 := shape.Circle(r / 2)
	s2 := solid.Loft(s1, s0, 200.0, 20.0)
	s2.STL("loft.stl", 3.0)
	return nil
}

func test49() error {
	s0 := shape.Circle(0.8)
	s0.ToDXF("circle_2d.dxf", 50)
	return nil
}

func test50() error {
	obj.Washer3D(obj.WasherParms{
		Thickness:   10,
		InnerRadius: 40,
		OuterRadius: 50,
		Remove:      0.3,
	}).STL("washer.stl", 3.0)
	return nil
}

func test51() error {
	s := obj.StdPipe3D("sch40:1", "mm", 100)
	s.STL("standard_pipe.stl", 3.0)
	return nil
}

type testFunc func() error

var testFuncs = []testFunc{
	test1,
	test2,
	test3,
	test4,
	test5,
	test6,
	test7,
	test9,
	test10,
	test11,
	test12,
	test13,
	test14,
	test15,
	test16,
	test17,
	test18,
	test19,
	test20,
	test21,
	test22,
	test26,
	test27,
	test28,
	test29,
	test30,
	test31,
	test32,
	test33,
	test34,
	test35,
	test36,
	test37,
	test39,
	test40,
	test41,
	test43,
	test44,
	test49,
	test50,
	test51,
}

func main() {
	for i, test := range testFuncs {
		err := test()
		if err != nil {
			log.Fatalf("error with testFuncs[%d]: %s\n", i, err)
		}
	}
}
