package solid

import (
	"math"
	"math/rand"
	"sort"
	"testing"

	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/sdfx/render"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
)

// meshVolume renders s and sums signed tetrahedron volumes (divergence
// theorem) — a coarse but orientation-honest volume for sanity checks.
func meshVolume(s *Solid, cellsPerMM float64) float64 {
	sz := s.Bounds().Size()
	longest := math.Max(sz.X, math.Max(sz.Y, sz.Z))
	cells := int(math.Ceil(longest * cellsPerMM))
	if cells < 30 {
		cells = 30
	}
	tris := render.ToTriangles(s.SDF3, render.NewMarchingCubesOctree(cells))
	vol := 0.0
	for _, t := range tris {
		a, b, c := t[0], t[1], t[2]
		vol += a.Dot(b.Cross(c)) / 6
	}
	return math.Abs(vol)
}

func gridSphere(t *testing.T, cells int) (*Solid, *Solid, GridInfo) {
	t.Helper()
	s := Sphere(10)
	g, info := SampleGridInfo(s, cells)
	return s, g, info
}

// Reconstruction error of a smooth field is bounded by the cell size —
// trilinear interpolation of a sphere's distance is exact along edges and
// sags by at most O(h²/R) inside a cell, far under h.
func TestGridReconstructsSmoothFieldWithinACell(t *testing.T) {
	s, g, info := gridSphere(t, 64)
	r := rand.New(rand.NewSource(7))
	worst := 0.0
	// Inside the source's own box: outside it the field is a deliberately
	// weak lower bound, tested separately.
	for n := 0; n < 20000; n++ {
		p := v3sdf.Vec{X: (r.Float64()*2 - 1) * 10, Y: (r.Float64()*2 - 1) * 10, Z: (r.Float64()*2 - 1) * 10}
		if e := math.Abs(g.SDF3.Evaluate(p) - s.SDF3.Evaluate(p)); e > worst {
			worst = e
		}
	}
	if worst > info.CellSize*0.25 {
		t.Fatalf("worst reconstruction error %.4f mm, cell %.4f mm", worst, info.CellSize)
	}
}

// The renderer's empty-cube skip needs |d| to not exceed the true distance
// to the surface by more than the interpolation error. Against the exact
// sphere field (whose |d| IS the true distance): on the smooth outside the
// linear interpolant of a convex function sits above it by ≤ 3h²/(8R); near
// the centre (the field's medial axis, a cone point) the interpolant
// averages across the kink and can be off by a fraction of a cell, which
// is deep inside and irrelevant to the renderer. Allow the larger of the
// two bounds by depth.
func TestGridOverestimatesOnlyByInterpolationError(t *testing.T) {
	s, g, info := gridSphere(t, 48)
	r := rand.New(rand.NewSource(11))
	h := info.CellSize
	// Chord sag of the interpolant across a cell DIAGONAL (length h√3):
	// (h√3)²/(8R) = 3h²/(8R).
	tolNear := 3*h*h/(8*10) + 1e-5
	for n := 0; n < 50000; n++ {
		p := v3sdf.Vec{X: (r.Float64()*2 - 1) * 14, Y: (r.Float64()*2 - 1) * 14, Z: (r.Float64()*2 - 1) * 14}
		exact := s.SDF3.Evaluate(p)
		got := g.SDF3.Evaluate(p)
		tol := tolNear
		if exact < -2*h { // deep interior: medial-axis kink territory
			tol = 0.5 * h
		}
		if math.Abs(got) > math.Abs(exact)+tol && math.Signbit(got) == math.Signbit(exact) {
			t.Fatalf("at %v: reconstructed |%.5f| exceeds true |%.5f| (tol %.5f)", p, got, exact, tol)
		}
	}
}

// Outside the sampled box the field must be a lower bound on the distance
// to the surface and must grow with distance from the box.
func TestGridOutsideBoxIsAConservativeLowerBound(t *testing.T) {
	s, g, _ := gridSphere(t, 32)
	for _, d := range []float64{1, 5, 20, 100} {
		p := v3sdf.Vec{X: 10 + 2 + d, Y: 0, Z: 0} // sphere r10, box padded ~2 cells
		got, exact := g.SDF3.Evaluate(p), s.SDF3.Evaluate(p)
		if got > exact+1e-6 {
			t.Errorf("outside at x=%.1f: %.4f exceeds true %.4f", p.X, got, exact)
		}
		if got < exact*0.5 {
			t.Errorf("outside at x=%.1f: %.4f is a needlessly weak bound on %.4f", p.X, got, exact)
		}
	}
}

// The reconstruction is continuous and its gradient stays near 1 away
// from the solid's sharp features: within a cell of an edge or of the
// field's medial axis the interpolant averages across a kink and can
// inflate to at most √3; elsewhere it must track the unit-gradient field.
func TestGridGradientStaysNearUnity(t *testing.T) {
	b := Box(v3.XYZ(10, 6, 4), 0)
	g, info := SampleGridInfo(b, 40)
	h := info.CellSize
	r := rand.New(rand.NewSource(3))
	step := h * 0.37
	worstSmooth, worstAll := 0.0, 0.0
	// Sample inside the padded grid (box half 5×3×2, pad 2h): the outside
	// bound is a different function with its own, separately tested,
	// behaviour.
	for n := 0; n < 30000; n++ {
		p := v3sdf.Vec{X: (r.Float64()*2 - 1) * (5 + 1.5*h), Y: (r.Float64()*2 - 1) * (3 + 1.5*h), Z: (r.Float64()*2 - 1) * (2 + 1.5*h)}
		dir := v3sdf.Vec{X: r.Float64()*2 - 1, Y: r.Float64()*2 - 1, Z: r.Float64()*2 - 1}
		if dir.Length() == 0 {
			continue
		}
		q := p.Add(dir.Normalize().MulScalar(step))
		ratio := math.Abs(g.SDF3.Evaluate(q)-g.SDF3.Evaluate(p)) / step
		if ratio > worstAll {
			worstAll = ratio
		}
		// "Smooth": more than three cells from the surface, and — inside,
		// where the box's distance field kinks wherever the nearest face
		// changes — as far from those medial surfaces, i.e. the two
		// nearest faces not within 3h of tying. Outside, the distance to a
		// convex body is C¹, so the surface margin suffices. Three cells
		// because the interpolant mixes partials from different edges of a
		// voxel: at distance m·h from an edge the field's direction turns
		// by ~1/m across a cell and the norm inflates by ~1/cos(1/2m) —
		// 1.06 at 1.5h, 1.014 at 3h.
		e := b.SDF3.Evaluate(p)
		smooth := e > 3*h
		if e < -3*h {
			d := []float64{5 - math.Abs(p.X), 3 - math.Abs(p.Y), 2 - math.Abs(p.Z)}
			sort.Float64s(d)
			smooth = d[1]-d[0] > 3*h
		}
		if smooth && ratio > worstSmooth {
			worstSmooth = ratio
		}
	}
	if worstSmooth > 1.03 {
		t.Fatalf("gradient inflation %.3f away from features", worstSmooth)
	}
	if worstAll > math.Sqrt(3)+0.05 {
		t.Fatalf("gradient inflation %.3f exceeds the trilinear bound √3", worstAll)
	}
}

// A grid sampled over a caller-chosen region is accurate everywhere in
// that region, including far from the surface — the case a boolean
// partner needs.
func TestGridBoundsCoversTheFarField(t *testing.T) {
	s := Sphere(10)
	big := NewBox3(v3.XYZ(0, 0, 0), v3.XYZ(80, 80, 80))
	g, info := SampleGridBoundsInfo(s, big, 128)
	for _, p := range []v3sdf.Vec{{X: 30, Y: 0, Z: 0}, {X: 25, Y: 25, Z: 0}, {X: 20, Y: 20, Z: 20}, {X: -35, Y: 10, Z: -5}} {
		exact, got := s.SDF3.Evaluate(p), g.SDF3.Evaluate(p)
		if math.Abs(got-exact) > info.CellSize*0.25 {
			t.Errorf("at %v: %.4f vs exact %.4f (cell %.3f)", p, got, exact, info.CellSize)
		}
	}
	if b := g.Bounds(); b.Max.X > 10.001 {
		t.Errorf("reported bounds must stay the source's, got %v", b)
	}
}

// Sharp features: near a box's edges and corners the interpolant may
// overestimate the distance by a small fraction of a cell (the field
// there is the distance to an edge, convex with a tiny curvature radius),
// never by more.
func TestGridBoxCornersStayConservative(t *testing.T) {
	b := Box(v3.XYZ(10, 6, 4), 0)
	g, info := SampleGridInfo(b, 40)
	r := rand.New(rand.NewSource(5))
	for n := 0; n < 50000; n++ {
		p := v3sdf.Vec{X: (r.Float64()*2 - 1) * 7, Y: (r.Float64()*2 - 1) * 5, Z: (r.Float64()*2 - 1) * 4}
		exact, got := b.SDF3.Evaluate(p), g.SDF3.Evaluate(p)
		if math.Signbit(got) == math.Signbit(exact) && math.Abs(got) > math.Abs(exact)+0.25*info.CellSize {
			t.Fatalf("at %v: |%.5f| > true |%.5f|", p, got, exact)
		}
	}
}

// The grid reports the source's bounds, not its padded sample box, so
// neighbours that size themselves off Bounds() do not grow.
func TestGridReportsSourceBounds(t *testing.T) {
	s, g, _ := gridSphere(t, 16)
	a, b := s.Bounds(), g.Bounds()
	if a.Min != b.Min || a.Max != b.Max {
		t.Fatalf("bounds %v != source %v", b, a)
	}
}

// Two fills of the same solid are bit-identical: the parallel fill writes
// each sample exactly once.
func TestGridFillIsDeterministic(t *testing.T) {
	s := Sphere(10).Union(Box(v3.XYZ(4, 4, 30), 0))
	a := s.SampleGrid(48).SDF3.(*GridField)
	b := s.SampleGrid(48).SDF3.(*GridField)
	if len(a.vals) != len(b.vals) {
		t.Fatal("different sample counts")
	}
	for i := range a.vals {
		if a.vals[i] != b.vals[i] {
			t.Fatalf("sample %d differs: %v vs %v", i, a.vals[i], b.vals[i])
		}
	}
}

// A boolean on a grid-backed solid renders to about the right volume —
// the field is usable as an ordinary operand, not just for lookups.
func TestGridSurvivesBooleanAndRender(t *testing.T) {
	g := Sphere(10).SampleGrid(64)
	cut := Box(v3.XYZ(30, 30, 30), 0).Cut(g)
	want := 30*30*30 - 4.0/3*math.Pi*1000
	got := meshVolume(cut, 3)
	if math.Abs(got-want)/want > 0.02 {
		t.Fatalf("box-minus-grid-sphere volume %.1f, want ≈%.1f", got, want)
	}
}

func TestGridInfo(t *testing.T) {
	_, _, info := gridSphere(t, 50)
	if info.CellSize <= 0 || info.Samples != (info.Cells[0]+1)*(info.Cells[1]+1)*(info.Cells[2]+1) || info.Bytes != info.Samples*4 {
		t.Fatalf("inconsistent info %+v", info)
	}
	if info.Cells[0] != 54 { // 50 along the longest axis + 2 cells of padding each side
		t.Errorf("cells along x = %d, want 54", info.Cells[0])
	}
}

// Across the grid boundary the field is continuous and still a lower
// bound — including at the sampled box's corners, where the naive
// "distance to box + minimum boundary value" bound jumps.
func TestGridOutsideBoundIsContinuousAtCorners(t *testing.T) {
	s, g, info := gridSphere(t, 32)
	h := info.CellSize
	edge := 10 + 2*h // padded box half-size
	in := v3sdf.Vec{X: edge - 0.01, Y: edge - 0.01, Z: edge - 0.01}
	out := v3sdf.Vec{X: edge + 0.01, Y: edge + 0.01, Z: edge + 0.01}
	vi, vo := g.SDF3.Evaluate(in), g.SDF3.Evaluate(out)
	if math.Abs(vi-vo) > 0.05 {
		t.Fatalf("jump across the grid corner: inside %.4f, outside %.4f", vi, vo)
	}
	if vo > s.SDF3.Evaluate(out)+1e-6 {
		t.Fatalf("outside bound %.4f exceeds the true distance %.4f", vo, s.SDF3.Evaluate(out))
	}
}
