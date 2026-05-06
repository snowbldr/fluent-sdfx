package v2i_test

import (
	"testing"

	"github.com/snowbldr/fluent-sdfx/vec/v2i"
	v2isdf "github.com/snowbldr/sdfx/vec/v2i"
)

func TestConstructors(t *testing.T) {
	if v2i.X(5) != (v2i.Vec{X: 5}) {
		t.Fatal("X")
	}
	if v2i.Y(7) != (v2i.Vec{Y: 7}) {
		t.Fatal("Y")
	}
	if v2i.XY(1, 2) != (v2i.Vec{X: 1, Y: 2}) {
		t.Fatal("XY")
	}
}

func TestRaw(t *testing.T) {
	a := v2i.XY(3, 4)
	if r := a.Raw(); r != (v2isdf.Vec{X: 3, Y: 4}) {
		t.Fatalf("Raw: %v", r)
	}
}

func TestArithmetic(t *testing.T) {
	a := v2i.XY(1, 2)
	b := v2i.XY(3, 4)
	if got := a.Add(b); got != v2i.XY(4, 6) {
		t.Fatalf("Add: %v", got)
	}
	if got := a.AddScalar(5); got != v2i.XY(6, 7) {
		t.Fatalf("AddScalar: %v", got)
	}
	if got := b.SubScalar(1); got != v2i.XY(2, 3) {
		t.Fatalf("SubScalar: %v", got)
	}
}
