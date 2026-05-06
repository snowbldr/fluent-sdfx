package p2_test

import (
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/vec/p2"
	p2sdf "github.com/snowbldr/sdfx/vec/p2"
)

func TestConstructors(t *testing.T) {
	if got := p2.R(2); got != (p2.Vec{R: 2}) {
		t.Fatalf("R: %v", got)
	}
	if got := p2.T(math.Pi); got != (p2.Vec{Theta: math.Pi}) {
		t.Fatalf("T: %v", got)
	}
	if got := p2.RT(3, math.Pi/2); got != (p2.Vec{R: 3, Theta: math.Pi / 2}) {
		t.Fatalf("RT: %v", got)
	}
}

func TestRaw(t *testing.T) {
	a := p2.RT(2, 1)
	if r := a.Raw(); r != (p2sdf.Vec{R: 2, Theta: 1}) {
		t.Fatalf("Raw: %v", r)
	}
}
