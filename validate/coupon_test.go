package validate

import (
	"math"
	"strings"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func TestComponentsCountsBodies(t *testing.T) {
	cases := []struct {
		name string
		s    *solid.Solid
		want int
	}{
		{"one box", ohBox(10, 10, 10, 0), 1},
		{"two disjoint boxes", ohBox(6, 6, 10, 0).Translate(v3.X(-10)).Union(ohBox(6, 6, 10, 0).Translate(v3.X(10))), 2},
		{"two touching boxes are one body", ohBox(6, 6, 10, 0).Translate(v3.X(-3)).Union(ohBox(6, 6, 10, 0).Translate(v3.X(3))), 1},
		{"three pillars", ohBox(4, 4, 8, 0).Translate(v3.X(-12)).Union(ohBox(4, 4, 8, 0), ohBox(4, 4, 8, 0).Translate(v3.X(12))), 3},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Components(tc.s, 6); got != tc.want {
				t.Errorf("Components = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestRequireOneBody(t *testing.T) {
	RequireOneBody(t, ohBox(10, 10, 10, 0), 6)

	f := &ohFakeTB{}
	split := ohBox(6, 6, 10, 0).Translate(v3.X(-10)).Union(ohBox(6, 6, 10, 0).Translate(v3.X(10)))
	RequireOneBody(f, split, 6)
	if len(f.errors) != 1 || !strings.Contains(f.errors[0], "2 separate bodies") {
		t.Errorf("errors = %v, want one naming 2 bodies", f.errors)
	}
}

func TestCouponCutsAndDrops(t *testing.T) {
	// A 40 mm bar floating between z=5 and z=15; the coupon keeps the
	// middle 10 mm and lands on the plate.
	bar := ohBox(40, 10, 10, 5)
	keep := solid.Box(v3.XYZ(10, 20, 20), 0).Translate(v3.Z(10))
	c := Coupon(bar, keep, 6)

	bb := c.Bounds()
	if math.Abs(bb.Min.Z) > 1e-9 {
		t.Errorf("coupon min Z = %.6f, want 0 (dropped to the plate)", bb.Min.Z)
	}
	size := bb.Size()
	if math.Abs(size.X-10) > 1e-6 || math.Abs(size.Y-10) > 1e-6 || math.Abs(size.Z-10) > 1e-6 {
		t.Errorf("coupon size = %.3f × %.3f × %.3f, want 10 × 10 × 10", size.X, size.Y, size.Z)
	}
	if !Inside(c, v3.XYZ(0, 0, 5)) {
		t.Error("coupon centre is not material")
	}
}

func TestCouponPanicsOnMultipleBodies(t *testing.T) {
	// A keep region spanning two pillars slices out two loose parts.
	pillars := ohBox(4, 4, 10, 0).Translate(v3.X(-8)).Union(ohBox(4, 4, 10, 0).Translate(v3.X(8)))
	keep := solid.Box(v3.XYZ(40, 20, 6), 0).Translate(v3.Z(5))

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("Coupon did not panic on a two-body cut")
		}
		msg, _ := r.(string)
		if !strings.Contains(msg, "2 separate bodies") {
			t.Errorf("panic = %v, want it to name the body count", r)
		}
	}()
	Coupon(pillars, keep, 6)
}

func TestDimplesCutIntoAFace(t *testing.T) {
	block := ohBox(20, 10, 10, 0) // x -10..10, y -5..5, z 0..10
	s := Dimples(block, 3, v3.XYZ(0, 0, 10), v3.Z(1), v3.X(1), 2, 1, 4)

	cases := []struct {
		name   string
		p      v3.Vec
		inside bool
	}{
		{"centre dimple is air", v3.XYZ(0, 0, 9.5), false},
		{"left dimple is air", v3.XYZ(-4, 0, 9.5), false},
		{"right dimple is air", v3.XYZ(4, 0, 9.5), false},
		{"between dimples is material", v3.XYZ(2, 0, 9.5), true},
		{"below the blind depth is material", v3.XYZ(0, 0, 8.5), true},
		{"off the row is material", v3.XYZ(0, 3, 9.5), true},
		{"far from the face is material", v3.XYZ(0, 0, 5), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Inside(s, tc.p); got != tc.inside {
				t.Errorf("Inside(%v) = %v, want %v (sdf %+.3f)", tc.p, got, tc.inside, At(s, tc.p))
			}
		})
	}
}

func TestDimplesOnASideFace(t *testing.T) {
	// The cutter is rotated to the face normal, not just sunk along -Z.
	block := ohBox(20, 10, 10, 0)
	s := Dimples(block, 2, v3.XYZ(10, 0, 5), v3.X(1), v3.Z(1), 2, 1.5, 4)
	if Inside(s, v3.XYZ(9.5, 0, 3)) {
		t.Errorf("lower dimple not cut: sdf %+.3f", At(s, v3.XYZ(9.5, 0, 3)))
	}
	if Inside(s, v3.XYZ(9.5, 0, 7)) {
		t.Errorf("upper dimple not cut: sdf %+.3f", At(s, v3.XYZ(9.5, 0, 7)))
	}
	if !Inside(s, v3.XYZ(9.5, 0, 5)) {
		t.Errorf("material between the dimples was removed: sdf %+.3f", At(s, v3.XYZ(9.5, 0, 5)))
	}
	if !Inside(s, v3.XYZ(8, 0, 3)) {
		t.Errorf("dimple is deeper than asked: sdf %+.3f", At(s, v3.XYZ(8, 0, 3)))
	}
}

func TestDimplesEdgeCases(t *testing.T) {
	block := ohBox(10, 10, 10, 0)
	if got := Dimples(block, 0, v3.XYZ(0, 0, 10), v3.Z(1), v3.X(1), 2, 1, 4); got != block {
		t.Error("Dimples with n=0 should return the solid unchanged")
	}
	for _, tc := range []struct {
		name  string
		call  func()
		match string
	}{
		{"zero normal", func() { Dimples(block, 1, v3.Vec{}, v3.Vec{}, v3.X(1), 2, 1, 4) }, "normal"},
		{"zero along", func() { Dimples(block, 1, v3.Vec{}, v3.Z(1), v3.Vec{}, 2, 1, 4) }, "along"},
		{"zero dia", func() { Dimples(block, 1, v3.Vec{}, v3.Z(1), v3.X(1), 0, 1, 4) }, "dia and depth"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatal("no panic")
				}
				if msg, _ := r.(string); !strings.Contains(msg, tc.match) {
					t.Errorf("panic = %v, want it to mention %q", r, tc.match)
				}
			}()
			tc.call()
		})
	}
}
