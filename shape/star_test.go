package shape_test

import (
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/shape"
	v2sdf "github.com/snowbldr/sdfx/vec/v2"
)

// Star's tips must be exactly at outer and its valleys exactly at inner.
// An earlier version offset the polygon by outer/10 "to round the corners",
// so Star(10, 5, 5) reached 11; that is a bug in a constructor whose
// arguments are radii.
func TestStarTipsAndValleysAreExact(t *testing.T) {
	s := shape.Star(10, 5, 5)
	b := s.Bounds()
	if math.Abs(b.Max.X-10) > 1e-9 {
		t.Fatalf("tip on +X at %.6f, want 10", b.Max.X)
	}
	at := func(r, deg float64) float64 {
		a := deg * math.Pi / 180
		return s.Evaluate(v2sdf.Vec{X: r * math.Cos(a), Y: r * math.Sin(a)})
	}
	if at(9.99, 0) >= 0 || at(10.01, 0) <= 0 {
		t.Errorf("tip is not at r=10: sdf(9.99)=%.4f sdf(10.01)=%.4f", at(9.99, 0), at(10.01, 0))
	}
	// Valleys sit halfway between tips, at 36° for a five-point star.
	if at(4.9, 36) >= 0 || at(5.1, 36) <= 0 {
		t.Errorf("valley is not at r=5: sdf(4.9)=%.4f sdf(5.1)=%.4f", at(4.9, 36), at(5.1, 36))
	}
}
