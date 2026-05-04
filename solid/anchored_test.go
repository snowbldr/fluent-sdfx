package solid

import (
	"math"
	"testing"

	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const eps = 1e-9

func vecClose(a, b v3.Vec) bool {
	return math.Abs(a.X-b.X) < eps && math.Abs(a.Y-b.Y) < eps && math.Abs(a.Z-b.Z) < eps
}

func TestTopAnchorOfCylinder(t *testing.T) {
	got := Cylinder(10, 5, 0).Top().Point
	want := v3.XYZ(0, 0, 5)
	if !vecClose(got, want) {
		t.Fatalf("Cylinder(10,5,0).Top().Point = %+v, want %+v", got, want)
	}
}

func TestRightAnchorAfterTranslate(t *testing.T) {
	got := Box(v3.XYZ(20, 10, 6), 0).TranslateX(15).Right().Point
	want := v3.XYZ(25, 0, 0)
	if !vecClose(got, want) {
		t.Fatalf("Box(20,10,6).TranslateX(15).Right().Point = %+v, want %+v", got, want)
	}
}

func TestOnAlignsAnchors(t *testing.T) {
	body := Cylinder(10, 5, 0)
	cap := Sphere(3)
	moved := cap.Bottom().On(body.Top()).Solid()

	gotBottom := moved.Bottom().Point
	wantTop := body.Top().Point
	if !vecClose(gotBottom, wantTop) {
		t.Fatalf("after On: moved.Bottom() = %+v, body.Top() = %+v", gotBottom, wantTop)
	}
}

func TestOnTopOfEquivalence(t *testing.T) {
	a := Box(v3.XYZ(4, 4, 4), 0)
	b := Cylinder(10, 5, 0)

	got := a.OnTopOf(b.Top(), 2).Solid().Bottom().Point
	want := a.Bottom().Above(b.Top(), 2).Solid().Bottom().Point
	if !vecClose(got, want) {
		t.Fatalf("OnTopOf(target,2) ≠ Bottom().Above(target,2): got %+v, want %+v", got, want)
	}
}

func TestBottomAtPreservesXY(t *testing.T) {
	s := Box(v3.XYZ(4, 6, 8), 0).Translate(v3.XYZ(7, -3, 11))
	moved := s.BottomAt(0)
	c := moved.Bounds().Center()
	if math.Abs(c.X-7) > eps || math.Abs(c.Y+3) > eps {
		t.Fatalf("BottomAt should not change X/Y: center = %+v, want X=7, Y=-3", c)
	}
	if math.Abs(moved.Bottom().Point.Z) > eps {
		t.Fatalf("BottomAt(0): bottom Z = %v, want 0", moved.Bottom().Point.Z)
	}
}

func TestPlacementCutEqualsManualCut(t *testing.T) {
	body := Box(v3.XYZ(20, 20, 20), 0)
	drill := Cylinder(30, 2, 0)

	auto := drill.OnTopOf(body.Top()).Cut()
	manual := body.Cut(drill.Bottom().On(body.Top()).Solid())

	a := auto.Bounds()
	m := manual.Bounds()
	if !a.Equals(m, 1e-9) {
		t.Fatalf("Placement Cut bounds differ from manual: auto=%+v manual=%+v", a, m)
	}
}
