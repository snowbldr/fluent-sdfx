package v3i_test

import (
	"testing"

	"github.com/snowbldr/fluent-sdfx/vec/v3i"
	v3isdf "github.com/snowbldr/sdfx/vec/v3i"
)

func TestConstructors(t *testing.T) {
	if v3i.X(1) != (v3i.Vec{X: 1}) {
		t.Fatal("X")
	}
	if v3i.Y(2) != (v3i.Vec{Y: 2}) {
		t.Fatal("Y")
	}
	if v3i.Z(3) != (v3i.Vec{Z: 3}) {
		t.Fatal("Z")
	}
	if v3i.XY(1, 2) != (v3i.Vec{X: 1, Y: 2}) {
		t.Fatal("XY")
	}
	if v3i.XZ(1, 3) != (v3i.Vec{X: 1, Z: 3}) {
		t.Fatal("XZ")
	}
	if v3i.YZ(2, 3) != (v3i.Vec{Y: 2, Z: 3}) {
		t.Fatal("YZ")
	}
	if v3i.XYZ(1, 2, 3) != (v3i.Vec{X: 1, Y: 2, Z: 3}) {
		t.Fatal("XYZ")
	}
}

func TestRaw(t *testing.T) {
	a := v3i.XYZ(3, 4, 5)
	if r := a.Raw(); r != (v3isdf.Vec{X: 3, Y: 4, Z: 5}) {
		t.Fatalf("Raw: %v", r)
	}
}

func TestArithmetic(t *testing.T) {
	a := v3i.XYZ(1, 2, 3)
	b := v3i.XYZ(4, 5, 6)
	if got := a.Add(b); got != v3i.XYZ(5, 7, 9) {
		t.Fatalf("Add: %v", got)
	}
	if got := a.AddScalar(10); got != v3i.XYZ(11, 12, 13) {
		t.Fatalf("AddScalar: %v", got)
	}
	if got := b.SubScalar(1); got != v3i.XYZ(3, 4, 5) {
		t.Fatalf("SubScalar: %v", got)
	}
}
