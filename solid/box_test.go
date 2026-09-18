package solid_test

import (
	"strings"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Box documents a panic when round exceeds half the smallest dimension.
// It used to accept it silently and produce a box whose fillets overlap.
func TestBoxRoundTooLargePanics(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Box(10x10x10, round 6) did not panic")
		}
		if !strings.Contains(r.(string), "round") {
			t.Fatalf("unexpected panic message: %v", r)
		}
	}()
	solid.Box(v3.XYZ(10, 10, 10), 6)
}

func TestBoxRoundAtHalfIsAllowed(t *testing.T) {
	s := solid.Box(v3.XYZ(10, 20, 30), 5) // exactly half the smallest side: a rounded-end pill, fine
	if b := s.Bounds(); b.Max.X != 5 {
		t.Fatalf("bounds %v", b)
	}
}
