package validate

import (
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// 10 mm cube: V = 1000 mm³, surface area = 600 mm². Tolerances are loose
// because marching cubes only matches the exact value as cells → ∞.
func TestCubeMetrics(t *testing.T) {
	const side = 10.0
	cube := solid.Box(v3.XYZ(side, side, side), 0)
	st := Of(cube, 8.0)
	if !st.Watertight {
		t.Fatalf("cube should be watertight; got %d boundary edges", st.BoundaryEdges)
	}
	if rel := math.Abs(st.Volume-1000) / 1000; rel > 0.05 {
		t.Errorf("volume = %.2f mm³, want 1000 ± 5%% (off by %.2f%%)", st.Volume, rel*100)
	}
	if rel := math.Abs(st.SurfaceArea-600) / 600; rel > 0.05 {
		t.Errorf("surface area = %.2f mm², want 600 ± 5%% (off by %.2f%%)", st.SurfaceArea, rel*100)
	}
}

// A sphere of r = 10: V = 4π/3 · r³ ≈ 4188.79 mm³.
func TestSphereVolume(t *testing.T) {
	want := 4.0 / 3.0 * math.Pi * 1000 // r = 10
	sphere := solid.Sphere(10)
	st := Of(sphere, 8.0)
	if !st.Watertight {
		t.Fatalf("sphere should be watertight")
	}
	if rel := math.Abs(st.Volume-want) / want; rel > 0.02 {
		t.Errorf("sphere volume = %.2f, want %.2f ± 2%% (off by %.2f%%)", st.Volume, want, rel*100)
	}
}

// A flat-roof box has a 100% horizontal-ceiling top — entire top face is a
// 90° overhang. For a 10mm cube the top is 100 mm².
func TestOverhangFlatTop(t *testing.T) {
	cube := solid.Box(v3.XYZ(10, 10, 10), 0)
	st := Of(cube, 8.0)
	// Both top (100 mm² overhang at 90°) and bottom face look identical to
	// the test (n.z = -1 vs +1) — only the top counts. Z-normal triangles
	// face down only on the bottom face; the test should report ~100 mm².
	if rel := math.Abs(st.OverhangArea-100) / 100; rel > 0.05 {
		t.Errorf("flat-bottom overhang = %.2f mm², want ~100 ± 5%% (off by %.2f%%)", st.OverhangArea, rel*100)
	}
}

// A sphere has hemispherical overhang area: half its surface points
// downward beyond 0°, but only the bottom cap (below 45° from vertical)
// counts at the FDM threshold. Specifically for a unit sphere, the area of
// the spherical cap from polar angle 0 to 45° from -Z axis is 2π·r²·(1−cos45°)
// = 2π · r² · (1 − √2/2). For r = 10, that's ≈ 184.0 mm².
func TestOverhangSphere(t *testing.T) {
	want := 2 * math.Pi * 100 * (1 - math.Cos(math.Pi/4))
	sphere := solid.Sphere(10)
	st := Of(sphere, 12.0)
	if rel := math.Abs(st.OverhangArea-want) / want; rel > 0.10 {
		t.Errorf("sphere overhang area = %.2f mm², want %.2f ± 10%% (off by %.2f%%)", st.OverhangArea, want, rel*100)
	}
}

// RequireWatertight should pass for a sealed solid and fail for an open one.
// We can't construct a non-watertight solid via the public API easily; instead
// just exercise the happy path here.
func TestRequireWatertightPasses(t *testing.T) {
	RequireWatertight(t, solid.Box(v3.XYZ(5, 5, 5), 0), 8.0)
}

func TestRequireVolumeNearPasses(t *testing.T) {
	RequireVolumeNear(t, solid.Box(v3.XYZ(10, 10, 10), 0), 8.0, 1000.0, 0.05)
}

// RequireMaxOverhang with a generous tolerance should pass — the cube's only
// downward area is its flat bottom face (~100 mm²), well above 0 but a
// 200 mm² tolerance accepts it.
func TestRequireMaxOverhangPasses(t *testing.T) {
	RequireMaxOverhang(t, solid.Box(v3.XYZ(10, 10, 10), 0), 8.0, 45.0, 200.0)
}
