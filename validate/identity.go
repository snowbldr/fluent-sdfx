package validate

import (
	"fmt"
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// IdentityResult is the outcome of comparing two fields point by point: the
// largest absolute difference found (mm), where it was found, how many
// points were compared, and whether the largest difference was within
// tolerance.
type IdentityResult struct {
	Max    float64 // largest |At(a,p) - At(b,p)| over the sample, mm
	At     v3.Vec  // the sample point where Max occurred
	N      int     // points compared
	Passed bool
}

// String returns the one-line report, e.g.
// "[ok] identical: max |Δfield| = 0.000 at (0.000, 0.000, 0.000) over 1331 points".
func (r IdentityResult) String() string {
	tag := "[FAIL]"
	if r.Passed {
		tag = "[ok]"
	}
	return fmt.Sprintf("%s identical: max |Δfield| = %.6f at (%.3f, %.3f, %.3f) over %d points",
		tag, r.Max, r.At.X, r.At.Y, r.At.Z, r.N)
}

// Identical reports whether a and b have the same signed distance field at
// every point of pts within tol mm. It compares field values, not surfaces,
// so it is a refactor guard for solids that should be the same expression,
// not a shape comparison between two different constructions.
func Identical(a, b *solid.Solid, pts []v3.Vec, tol float64) IdentityResult {
	fa, fb := evalAll(a, pts), evalAll(b, pts)
	r := IdentityResult{N: len(pts)}
	for i := range pts {
		d := math.Abs(fa[i] - fb[i])
		if d > r.Max || (i == 0) {
			r.Max, r.At = d, pts[i]
		}
	}
	r.Passed = len(pts) > 0 && r.Max <= tol
	return r
}

// RequireIdentical compares the two fields at pts, logs the report line, and
// calls t.Errorf if any point differs by more than tol mm (or if pts is
// empty, which would otherwise pass vacuously).
func RequireIdentical(t testing.TB, a, b *solid.Solid, pts []v3.Vec, tol float64) {
	t.Helper()
	r := Identical(a, b, pts, tol)
	if !r.Passed {
		if r.N == 0 {
			t.Errorf("[FAIL] identical: no points compared")
			return
		}
		t.Errorf("%s (tol %.6f)", r, tol)
		return
	}
	t.Logf("%s (tol %.6f)", r, tol)
}

// ChangedRegion returns the bounding box of every surface point of either
// solid (sampled at cellsPerMM) at which the other solid's field exceeds tol
// in magnitude — that is, where a and b actually differ — together with the
// number of such points. Returns the zero box and 0 when the two agree
// everywhere, which is how you prove a change was local.
func ChangedRegion(a, b *solid.Solid, cellsPerMM, tol float64) (v3.Box, int) {
	var box v3.Box
	n := 0
	add := func(pts []v3.Vec, other *solid.Solid) {
		f := evalAll(other, pts)
		for i, p := range pts {
			if math.Abs(f[i]) <= tol {
				continue
			}
			if n == 0 {
				box = v3.Box{Min: p, Max: p}
			} else {
				box = box.Include(p)
			}
			n++
		}
	}
	add(Surface(a, cellsPerMM), b)
	add(Surface(b, cellsPerMM), a)
	return box, n
}

// RequireUnchangedOutside calls t.Errorf if any surface point of `before` or
// `after` that lies outside `region` reads further than tol mm from the
// other solid's surface — the assertion that an edit only moved geometry
// inside the region it was supposed to touch.
func RequireUnchangedOutside(t testing.TB, before, after *solid.Solid, region v3.Box, cellsPerMM, tol float64) {
	t.Helper()
	worst, n := 0.0, 0
	var at v3.Vec
	check := func(pts []v3.Vec, other *solid.Solid) {
		f := evalAll(other, pts)
		for i, p := range pts {
			if region.Contains(p) {
				continue
			}
			n++
			if d := math.Abs(f[i]); d > worst {
				worst, at = d, p
			}
		}
	}
	check(Surface(before, cellsPerMM), after)
	check(Surface(after, cellsPerMM), before)
	line := fmt.Sprintf("unchanged outside region: max |field| = %.4f at (%.3f, %.3f, %.3f) over %d surface points (tol %.4f)",
		worst, at.X, at.Y, at.Z, n, tol)
	if worst > tol {
		t.Errorf("[FAIL] %s", line)
		return
	}
	t.Logf("[ok] %s", line)
}

// RequireBounds calls t.Errorf if s's bounding box corners differ from want
// by more than tol mm in any component. The box compared is the SDF's own
// bounding box (s.Bounds()), which can be larger than the visible geometry
// for solids built with WithBounds or generous primitives.
func RequireBounds(t testing.TB, s *solid.Solid, want v3.Box, tol float64) {
	t.Helper()
	got := s.Bounds().Box
	line := fmt.Sprintf("bounds = (%.3f, %.3f, %.3f)..(%.3f, %.3f, %.3f) (want (%.3f, %.3f, %.3f)..(%.3f, %.3f, %.3f) ± %.3f)",
		got.Min.X, got.Min.Y, got.Min.Z, got.Max.X, got.Max.Y, got.Max.Z,
		want.Min.X, want.Min.Y, want.Min.Z, want.Max.X, want.Max.Y, want.Max.Z, tol)
	if !got.Equals(want, tol) {
		t.Errorf("[FAIL] %s", line)
		return
	}
	t.Logf("[ok] %s", line)
}

// RequireStandsOn calls t.Errorf unless the minimum Z of s's bounding box is
// within tol mm of z — the print-plate assertion, normally z = 0, that
// catches a part left floating or sunk after a transform.
func RequireStandsOn(t testing.TB, s *solid.Solid, z, tol float64) {
	t.Helper()
	RequireNear(t, "min Z", s.Bounds().Min.Z, z, tol)
}

// GridPoints returns a regular grid of points filling box at spacing step
// mm, starting at box.Min and including box.Max on any axis whose size is a
// whole multiple of step. Panics if step is not positive.
func GridPoints(box v3.Box, step float64) []v3.Vec {
	if step <= 0 {
		panic(fmt.Sprintf("validate.GridPoints: step must be > 0, got %v", step))
	}
	xs := identAxis(box.Min.X, box.Max.X, step)
	ys := identAxis(box.Min.Y, box.Max.Y, step)
	zs := identAxis(box.Min.Z, box.Max.Z, step)
	out := make([]v3.Vec, 0, len(xs)*len(ys)*len(zs))
	for _, x := range xs {
		for _, y := range ys {
			for _, z := range zs {
				out = append(out, v3.XYZ(x, y, z))
			}
		}
	}
	return out
}

// identAxis returns the sample coordinates from lo to hi at the given step,
// always including lo and never exceeding hi (to within rounding).
func identAxis(lo, hi, step float64) []float64 {
	if hi < lo {
		lo, hi = hi, lo
	}
	n := int(math.Floor((hi-lo)/step + 1e-9))
	out := make([]float64, 0, n+1)
	for i := 0; i <= n; i++ {
		out = append(out, lo+float64(i)*step)
	}
	return out
}
