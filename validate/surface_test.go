package validate

import (
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// surfTube is a tube of outer radius 5, bore radius 3 (so a 2 mm wall) and
// 20 mm tall, centred on the origin.
func surfTube() *solid.Solid {
	return solid.Cylinder(20, 5, 0).Cut(solid.Cylinder(30, 3, 0))
}

func TestSurfaceAlong(t *testing.T) {
	cube := probeCube() // |x|,|y|,|z| < 5
	tests := []struct {
		name     string
		from     v3.Vec
		dir      v3.Vec
		maxDist  float64
		want     float64
		wantFind bool
	}{
		{"out through +X face", v3.Vec{}, v3.X(1), 10, 5, true},
		{"out through -Y face", v3.Vec{}, v3.Y(-1), 10, 5, true},
		{"out through +Z face", v3.Vec{}, v3.Z(1), 10, 5, true},
		{"in from outside", v3.XYZ(10, 0, 0), v3.X(-1), 10, 5, true},
		{"unnormalised direction", v3.Vec{}, v3.X(7), 10, 5, true},
		{"too short to reach", v3.Vec{}, v3.X(1), 3, 3, false},
		{"parallel miss", v3.XYZ(0, 0, 20), v3.X(1), 10, 10, false},
		{"zero direction", v3.Vec{}, v3.Vec{}, 10, 10, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := SurfaceAlong(cube, tc.from, tc.dir, tc.maxDist, 0)
			if ok != tc.wantFind {
				t.Fatalf("found = %v, want %v (got %.6f)", ok, tc.wantFind, got)
			}
			if math.Abs(got-tc.want) > 1e-5 {
				t.Errorf("distance = %.6f, want %.6f", got, tc.want)
			}
		})
	}
}

// A step coarser than the feature still finds the crossing, because the
// crossing is detected by sign change over the step and then bisected.
func TestSurfaceAlongStep(t *testing.T) {
	cube := probeCube()
	for _, step := range []float64{0.5, 0.1, 0.01} {
		got, ok := SurfaceAlong(cube, v3.Vec{}, v3.X(1), 10, step)
		if !ok || math.Abs(got-5) > 1e-5 {
			t.Errorf("step %.2f: got (%.6f, %v), want (5, true)", step, got, ok)
		}
	}
}

// SurfaceR must march radially in the Cyl convention: deg counter-clockwise
// from +X seen from +Z. A 10 x 2 x 10 slab settles it — 5 mm out along X,
// 1 mm out along Y.
func TestSurfaceR(t *testing.T) {
	slab := solid.Box(v3.XYZ(10, 2, 10), 0)
	tube := surfTube()
	tests := []struct {
		name     string
		s        *solid.Solid
		deg, z   float64
		rMin     float64
		rMax     float64
		want     float64
		wantFind bool
	}{
		{"slab along +X", slab, 0, 0, 0, 20, 5, true},
		{"slab along +Y", slab, 90, 0, 0, 20, 1, true},
		{"slab along -X", slab, 180, 0, 0, 20, 5, true},
		{"slab along -Y", slab, 270, 0, 0, 20, 1, true},
		{"tube bore wall", tube, 0, 0, 0, 20, 3, true},
		{"tube outer wall", tube, 37, 0, 4, 20, 5, true},
		{"above the tube", tube, 0, 15, 0, 20, 20, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := SurfaceR(tc.s, tc.deg, tc.z, tc.rMin, tc.rMax)
			if ok != tc.wantFind {
				t.Fatalf("found = %v, want %v (got r = %.6f)", ok, tc.wantFind, got)
			}
			if math.Abs(got-tc.want) > 1e-5 {
				t.Errorf("radius = %.6f, want %.6f", got, tc.want)
			}
		})
	}
}

func TestSurfaceZ(t *testing.T) {
	// A 10 mm cube (top at z = 5) with a 20 mm tall, 2 mm radius post on top
	// (spanning z = 0..20 before the union, so the post top is z = 20).
	part := probeCube().Union(solid.Cylinder(20, 2, 0).TranslateZ(10))
	tests := []struct {
		name     string
		x, y     float64
		want     float64
		wantFind bool
	}{
		{"on the post", 0, 0, 20, true},
		{"on the cube top", 4, 4, 5, true},
		{"past both", 20, 20, -50, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := SurfaceZ(part, tc.x, tc.y, 50, -50)
			if ok != tc.wantFind {
				t.Fatalf("found = %v, want %v (got z = %.6f)", ok, tc.wantFind, got)
			}
			if math.Abs(got-tc.want) > 1e-5 {
				t.Errorf("z = %.6f, want %.6f", got, tc.want)
			}
		})
	}
}

func TestOpeningAndWallThickness(t *testing.T) {
	// Two 10 mm cubes either side of a 3 mm gap in X.
	pair := probeCube().TranslateX(-6.5).Union(probeCube().TranslateX(6.5))
	tube := surfTube()
	tests := []struct {
		name     string
		fn       func() (float64, bool)
		want     float64
		wantFind bool
	}{
		{"gap between two cubes", func() (float64, bool) { return Opening(pair, v3.Vec{}, v3.X(1), 20) }, 3, true},
		{"bore of the tube", func() (float64, bool) { return Opening(tube, v3.Vec{}, v3.X(1), 20) }, 6, true},
		{"opening from inside material", func() (float64, bool) { return Opening(pair, v3.XYZ(-6.5, 0, 0), v3.X(1), 20) }, 0, false},
		{"opening with no far wall", func() (float64, bool) { return Opening(pair, v3.Vec{}, v3.Y(1), 20) }, 0, false},
		{"tube wall radially", func() (float64, bool) { return WallThickness(tube, v3.XYZ(4, 0, 0), v3.X(1), 20) }, 2, true},
		{"tube wall axially", func() (float64, bool) { return WallThickness(tube, v3.XYZ(4, 0, 0), v3.Z(1), 30) }, 20, true},
		{"cube wall through the middle", func() (float64, bool) { return WallThickness(probeCube(), v3.Vec{}, v3.X(1), 20) }, 10, true},
		{"wall from air", func() (float64, bool) { return WallThickness(tube, v3.Vec{}, v3.X(1), 20) }, 0, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := tc.fn()
			if ok != tc.wantFind {
				t.Fatalf("found = %v, want %v (got %.6f)", ok, tc.wantFind, got)
			}
			if math.Abs(got-tc.want) > 1e-5 {
				t.Errorf("span = %.6f, want %.6f", got, tc.want)
			}
		})
	}
}

func TestMinWall(t *testing.T) {
	tube := surfTube() // 2 mm wall
	region := v3.Box{Min: v3.XYZ(-5, -5, -1), Max: v3.XYZ(5, 5, 1)}
	r := MinWall(tube, region, 1, 20)
	if r.Samples == 0 {
		t.Fatalf("no interior samples: %s", r)
	}
	if math.Abs(r.Min-2) > 0.01 {
		t.Errorf("min wall = %.4f, want 2.000 (%s)", r.Min, r)
	}
	if got := math.Hypot(r.At.X, r.At.Y); math.Abs(got-4) > 1.01 {
		t.Errorf("thinnest run reported at radius %.3f, want inside the 3..5 wall (%s)", got, r)
	}
	// An all-air region yields no samples at all.
	empty := MinWall(tube, v3.Box{Min: v3.XYZ(-1, -1, -1), Max: v3.XYZ(1, 1, 1)}, 1, 20)
	if empty.Samples != 0 || !math.IsInf(empty.Min, 1) {
		t.Errorf("all-air region: %s, want no samples and +Inf", empty)
	}
}

func TestRequireMinWall(t *testing.T) {
	tube := surfTube()
	region := v3.Box{Min: v3.XYZ(-5, -5, -1), Max: v3.XYZ(5, 5, 1)}
	tests := []struct {
		name     string
		region   v3.Box
		wantMin  float64
		wantFail bool
	}{
		{"2 mm wall passes a 1.5 mm budget", region, 1.5, false},
		{"2 mm wall fails a 3 mm budget", region, 3, true},
		{"empty region always fails", v3.Box{Min: v3.XYZ(-1, -1, -1), Max: v3.XYZ(1, 1, 1)}, 0.1, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := &probeTB{}
			RequireMinWall(rec, tube, tc.region, 1, 20, tc.wantMin)
			if got := len(rec.errs) > 0; got != tc.wantFail {
				t.Errorf("failed = %v, want %v (logs %v errs %v)", got, tc.wantFail, rec.logs, rec.errs)
			}
		})
	}
}

func TestWallReportString(t *testing.T) {
	r := WallReport{Min: 1.983, At: v3.XYZ(0, 0, 5), Dir: v3.Z(1), Samples: 412}
	want := "min wall = 1.983 mm at (0.000, 0.000, 5.000) along Z (412 samples)"
	if got := r.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
	if got := (WallReport{Min: math.Inf(1)}).String(); got != "min wall: no interior samples in region" {
		t.Errorf("empty report = %q", got)
	}
}

func TestRequireNear(t *testing.T) {
	tests := []struct {
		name           string
		got, want, tol float64
		wantFail       bool
	}{
		{"exact", 5, 5, 0.01, false},
		{"within tolerance", 5.005, 5, 0.01, false},
		{"outside tolerance", 5.02, 5, 0.01, true},
		{"NaN never passes", math.NaN(), 5, 0.01, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := &probeTB{}
			RequireNear(rec, "bore dia", tc.got, tc.want, tc.tol)
			if failed := len(rec.errs) > 0; failed != tc.wantFail {
				t.Errorf("failed = %v, want %v (logs %v errs %v)", failed, tc.wantFail, rec.logs, rec.errs)
			}
			all := append(append([]string{}, rec.logs...), rec.errs...)
			if !rec.has("bore dia = ", all) || !rec.has("(want 5.000 ± 0.010)", all) {
				t.Errorf("report line not in the standard form: %v", all)
			}
		})
	}
}

func TestSurfAxisName(t *testing.T) {
	tests := []struct {
		d    v3.Vec
		want string
	}{
		{v3.X(1), "X"},
		{v3.Y(-1), "Y"},
		{v3.Z(1), "Z"},
		{v3.XYZ(1, 1, 0), "(1.000, 1.000, 0.000)"},
	}
	for _, tc := range tests {
		if got := surfAxisName(tc.d); got != tc.want {
			t.Errorf("surfAxisName(%v) = %q, want %q", tc.d, got, tc.want)
		}
	}
}
