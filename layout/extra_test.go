package layout_test

import (
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/layout"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func TestSetGridLimit_RoundTrip(t *testing.T) {
	// Default limit is 1M; lower it, then attempt a Grid that exceeds the
	// new cap, then restore.
	layout.SetGridLimit(100)
	defer layout.SetGridLimit(1_000_000)

	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic for nx*ny exceeding the new limit")
		}
	}()
	_ = layout.Grid(1, 1, 11, 11) // 121 > 100
}

func TestGrid_ExceedsLimit(t *testing.T) {
	// At the default 1M limit, 1001x1001 = ~1.002M should panic.
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic for nx*ny > gridLimit")
		}
	}()
	_ = layout.Grid(1, 1, 1001, 1001)
}

func TestGrid2_ExceedsLimit(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic for 2D nx*ny > gridLimit")
		}
	}()
	_ = layout.Grid2(1, 1, 1001, 1001)
}

func TestPolar2_NegativeN(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic on Polar2 with n<0")
		}
	}()
	_ = layout.Polar2(5, -1)
}

func TestPolar_NegativeN(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic on Polar with n<0")
		}
	}()
	_ = layout.Polar(5, -1)
}

func TestPolarArc2_AllVariants(t *testing.T) {
	// n == 1 lands at startDeg.
	one := layout.PolarArc2(10, 1, 90, 360)
	if len(one) != 1 || math.Abs(one[0].X) > 1e-9 || math.Abs(one[0].Y-10) > 1e-9 {
		t.Fatalf("PolarArc2(10, 1, 90, 360) = %+v, want a single point at (0, 10)", one)
	}
	// n == 4 over 270° starting at 0: angles 0, 90, 180, 270.
	four := layout.PolarArc2(10, 4, 0, 270)
	if len(four) != 4 {
		t.Fatalf("PolarArc2(10, 4, 0, 270) returned %d points", len(four))
	}
	want := []v2.Vec{v2.XY(10, 0), v2.XY(0, 10), v2.XY(-10, 0), v2.XY(0, -10)}
	for i, p := range four {
		if math.Abs(p.X-want[i].X) > 1e-9 || math.Abs(p.Y-want[i].Y) > 1e-9 {
			t.Fatalf("PolarArc2 point %d = %+v, want %+v", i, p, want[i])
		}
	}
}

func TestLine2_SinglePoint(t *testing.T) {
	pts := layout.Line2(v2.XY(1, 2), v2.XY(7, 9), 1)
	if len(pts) != 1 || pts[0] != (v2.XY(1, 2)) {
		t.Fatalf("Line2 with n=1 should be just p0; got %+v", pts)
	}
}

// Cross-package sanity: spreading layout.Polar through .Multi-like usage
// stays at z=0 (a property the docs claim).
func TestPolar_AllAtZeroZ(t *testing.T) {
	pts := layout.Polar(7.5, 16)
	for i, p := range pts {
		if p.Z != 0 {
			t.Fatalf("Polar point %d has Z=%v, want 0", i, p.Z)
		}
	}
}

// Compile-time check that Polar's output type plugs into a Multi-style API.
var _ = []v3.Vec(layout.Polar(1, 1))
