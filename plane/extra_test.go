package plane_test

import (
	"testing"

	"github.com/snowbldr/fluent-sdfx/plane"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func TestAxisNormals(t *testing.T) {
	if plane.X != v3.XYZ(1, 0, 0) {
		t.Fatalf("X: %v", plane.X)
	}
	if plane.Y != v3.XYZ(0, 1, 0) {
		t.Fatalf("Y: %v", plane.Y)
	}
	if plane.Z != v3.XYZ(0, 0, 1) {
		t.Fatalf("Z: %v", plane.Z)
	}
}

func TestAt(t *testing.T) {
	p := plane.At(v3.XYZ(1, 2, 3), v3.XYZ(0, 0, 1))
	if p.Origin != v3.XYZ(1, 2, 3) || p.Normal != v3.XYZ(0, 0, 1) {
		t.Fatalf("At: %+v", p)
	}
}

func TestAtAxes(t *testing.T) {
	if p := plane.AtX(5); p.Origin != v3.X(5) || p.Normal != plane.X {
		t.Fatalf("AtX: %+v", p)
	}
	if p := plane.AtY(5); p.Origin != v3.Y(5) || p.Normal != plane.Y {
		t.Fatalf("AtY: %+v", p)
	}
	if p := plane.AtZ(5); p.Origin != v3.Z(5) || p.Normal != plane.Z {
		t.Fatalf("AtZ: %+v", p)
	}
}

func TestStandardPlanes(t *testing.T) {
	if plane.XY.Origin != v3.Z(0) || plane.XY.Normal != plane.Z {
		t.Fatalf("XY: %+v", plane.XY)
	}
	if plane.XZ.Origin != v3.Y(0) || plane.XZ.Normal != plane.Y {
		t.Fatalf("XZ: %+v", plane.XZ)
	}
	if plane.YZ.Origin != v3.X(0) || plane.YZ.Normal != plane.X {
		t.Fatalf("YZ: %+v", plane.YZ)
	}
}
