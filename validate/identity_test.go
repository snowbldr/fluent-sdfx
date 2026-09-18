package validate

import (
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// identBore is the 20 mm cube with a 4 mm blind bore sunk 5 mm into its top
// face (the cube's top is z = 10, the bore floor is z = 5).
func identBore() (before, after *solid.Solid) {
	before = solid.Box(v3.XYZ(20, 20, 20), 0)
	return before, before.Cut(solid.Cylinder(10, 2, 0).TranslateZ(10))
}

func TestIdentical(t *testing.T) {
	cube := solid.Box(v3.XYZ(10, 10, 10), 0)
	same := solid.Box(v3.XYZ(10, 10, 10), 0)
	bigger := solid.Box(v3.XYZ(10.5, 10, 10), 0)
	pts := GridPoints(v3.Box{Min: v3.XYZ(-8, -8, -8), Max: v3.XYZ(8, 8, 8)}, 2)
	tests := []struct {
		name     string
		a, b     *solid.Solid
		tol      float64
		wantPass bool
		wantMax  float64
	}{
		{"same expression", cube, cube, 1e-9, true, 0},
		{"rebuilt the same way", cube, same, 1e-9, true, 0},
		{"0.25 mm bigger in X", cube, bigger, 1e-9, false, 0.25},
		{"0.25 mm bigger, loose tolerance", cube, bigger, 0.3, true, 0.25},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := Identical(tc.a, tc.b, pts, tc.tol)
			if r.Passed != tc.wantPass {
				t.Errorf("Passed = %v, want %v (%s)", r.Passed, tc.wantPass, r)
			}
			if math.Abs(r.Max-tc.wantMax) > 1e-6 {
				t.Errorf("Max = %.6f, want %.6f", r.Max, tc.wantMax)
			}
			if r.N != len(pts) {
				t.Errorf("N = %d, want %d", r.N, len(pts))
			}
		})
	}
	if r := Identical(cube, cube, nil, 1e-9); r.Passed {
		t.Errorf("an empty point set must not pass vacuously: %s", r)
	}
}

func TestRequireIdentical(t *testing.T) {
	cube := solid.Box(v3.XYZ(10, 10, 10), 0)
	bigger := solid.Box(v3.XYZ(11, 10, 10), 0)
	pts := GridPoints(v3.Box{Min: v3.XYZ(-4, -4, -4), Max: v3.XYZ(4, 4, 4)}, 2)
	tests := []struct {
		name     string
		b        *solid.Solid
		pts      []v3.Vec
		wantFail bool
	}{
		{"identical", cube, pts, false},
		{"different", bigger, pts, true},
		{"no points", cube, nil, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := &probeTB{}
			RequireIdentical(rec, cube, tc.b, tc.pts, 1e-9)
			if failed := len(rec.errs) > 0; failed != tc.wantFail {
				t.Errorf("failed = %v, want %v (logs %v errs %v)", failed, tc.wantFail, rec.logs, rec.errs)
			}
		})
	}
}

func TestChangedRegion(t *testing.T) {
	before, after := identBore()
	t.Run("identical solids report no change", func(t *testing.T) {
		box, n := ChangedRegion(before, before, 2, 0.01)
		if n != 0 || box != (v3.Box{}) {
			t.Errorf("ChangedRegion of a solid with itself = (%v, %d), want (zero box, 0)", box, n)
		}
	})
	t.Run("a blind bore changes only its own neighbourhood", func(t *testing.T) {
		box, n := ChangedRegion(before, after, 2, 0.05)
		if n == 0 {
			t.Fatalf("a 4 mm bore was cut but no change was found")
		}
		// The bore is r <= 2 about the Z axis, from z = 5 to the top face.
		if box.Min.X < -3 || box.Max.X > 3 || box.Min.Y < -3 || box.Max.Y > 3 {
			t.Errorf("changed region %v spreads beyond the bore in X/Y", box)
		}
		if box.Min.Z < 4 || box.Max.Z > 10.5 {
			t.Errorf("changed region %v spreads beyond z = 5..10", box)
		}
	})
}

func TestRequireUnchangedOutside(t *testing.T) {
	before, after := identBore()
	tests := []struct {
		name     string
		region   v3.Box
		wantFail bool
	}{
		{"region covers the bore", v3.Box{Min: v3.XYZ(-4, -4, 4), Max: v3.XYZ(4, 4, 11)}, false},
		{"region misses the bore", v3.Box{Min: v3.XYZ(-9, -9, -9), Max: v3.XYZ(-5, -5, -5)}, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := &probeTB{}
			RequireUnchangedOutside(rec, before, after, tc.region, 2, 0.2)
			if failed := len(rec.errs) > 0; failed != tc.wantFail {
				t.Errorf("failed = %v, want %v (logs %v errs %v)", failed, tc.wantFail, rec.logs, rec.errs)
			}
		})
	}
}

func TestRequireBounds(t *testing.T) {
	cube := solid.Box(v3.XYZ(10, 10, 10), 0)
	tests := []struct {
		name     string
		want     v3.Box
		tol      float64
		wantFail bool
	}{
		{"exact", v3.Box{Min: v3.XYZ(-5, -5, -5), Max: v3.XYZ(5, 5, 5)}, 1e-9, false},
		{"within tolerance", v3.Box{Min: v3.XYZ(-5.01, -5, -5), Max: v3.XYZ(5, 5, 5)}, 0.05, false},
		{"wrong size", v3.Box{Min: v3.XYZ(-6, -5, -5), Max: v3.XYZ(5, 5, 5)}, 1e-9, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := &probeTB{}
			RequireBounds(rec, cube, tc.want, tc.tol)
			if failed := len(rec.errs) > 0; failed != tc.wantFail {
				t.Errorf("failed = %v, want %v (logs %v errs %v)", failed, tc.wantFail, rec.logs, rec.errs)
			}
		})
	}
}

func TestRequireStandsOn(t *testing.T) {
	cube := solid.Box(v3.XYZ(10, 10, 10), 0)
	tests := []struct {
		name     string
		s        *solid.Solid
		z        float64
		wantFail bool
	}{
		{"zeroed to the plate", cube.ZeroZ(), 0, false},
		{"still centred on the origin", cube, 0, true},
		{"centred, asked for -5", cube, -5, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := &probeTB{}
			RequireStandsOn(rec, tc.s, tc.z, 1e-6)
			if failed := len(rec.errs) > 0; failed != tc.wantFail {
				t.Errorf("failed = %v, want %v (logs %v errs %v)", failed, tc.wantFail, rec.logs, rec.errs)
			}
		})
	}
}

func TestGridPoints(t *testing.T) {
	tests := []struct {
		name  string
		box   v3.Box
		step  float64
		wantN int
	}{
		{"unit cube at 1 mm", v3.Box{Min: v3.Vec{}, Max: v3.XYZ(2, 2, 2)}, 1, 27},
		{"max face off the grid", v3.Box{Min: v3.Vec{}, Max: v3.XYZ(2, 2, 2)}, 0.75, 27},
		{"step larger than the box", v3.Box{Min: v3.Vec{}, Max: v3.XYZ(1, 1, 1)}, 5, 1},
		{"flat in Z", v3.Box{Min: v3.XYZ(0, 0, 3), Max: v3.XYZ(1, 1, 3)}, 0.5, 9},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pts := GridPoints(tc.box, tc.step)
			if len(pts) != tc.wantN {
				t.Errorf("len = %d, want %d", len(pts), tc.wantN)
			}
			for _, p := range pts {
				if p.X < tc.box.Min.X-1e-9 || p.X > tc.box.Max.X+1e-9 ||
					p.Y < tc.box.Min.Y-1e-9 || p.Y > tc.box.Max.Y+1e-9 ||
					p.Z < tc.box.Min.Z-1e-9 || p.Z > tc.box.Max.Z+1e-9 {
					t.Fatalf("point %v outside the box %v", p, tc.box)
				}
			}
			if len(pts) > 0 && pts[0] != tc.box.Min {
				t.Errorf("first point = %v, want box.Min %v", pts[0], tc.box.Min)
			}
		})
	}
	t.Run("non-positive step panics", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Errorf("GridPoints with step 0 must panic")
			}
		}()
		GridPoints(v3.Box{Max: v3.XYZ(1, 1, 1)}, 0)
	})
}

func TestIdentityResultString(t *testing.T) {
	r := IdentityResult{Max: 0, At: v3.Vec{}, N: 1331, Passed: true}
	want := "[ok] identical: max |Δfield| = 0.000000 at (0.000, 0.000, 0.000) over 1331 points"
	if got := r.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
	r.Passed = false
	if got := r.String(); got[:6] != "[FAIL]" {
		t.Errorf("failing report = %q, want a [FAIL] prefix", got)
	}
}
