package validate

import (
	"fmt"
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// surfBisectTol is the distance (mm) to which a detected sign change is
// bisected before a surface position is reported.
const surfBisectTol = 1e-6

// surfDefaultStep is the ray-march step (mm) used when a caller passes 0.
// A surface thinner than this can be stepped over; pass an explicit step
// when measuring features finer than a tenth of a millimetre.
const surfDefaultStep = 0.1

// SurfaceAlong marches from `from` along dir (normalised here) in steps of
// `step` mm (0 means 0.1) up to maxDist mm, until the sign of the field
// changes, then bisects that interval to 1e-6 mm. It returns the distance
// from `from` to that surface and true, or (maxDist, false) if the field
// never changes sign. The starting point may be inside or outside material:
// the crossing found is the first one in either direction of sign.
func SurfaceAlong(s *solid.Solid, from, dir v3.Vec, maxDist, step float64) (float64, bool) {
	if step <= 0 {
		step = surfDefaultStep
	}
	if maxDist <= 0 || dir.Length() == 0 {
		return maxDist, false
	}
	d := dir.Normalize()
	in0 := At(s, from) < 0
	prev := 0.0
	for t := step; ; t += step {
		last := false
		if t >= maxDist {
			t, last = maxDist, true
		}
		if (At(s, from.Add(d.MulScalar(t))) < 0) != in0 {
			return surfBisect(s, from, d, prev, t, in0), true
		}
		if last {
			break
		}
		prev = t
	}
	return maxDist, false
}

// SurfaceR returns the radius (mm) at which a surface is met when marching
// outward from the Z axis at angle deg (degrees counter-clockwise from +X
// seen from +Z) and height z, starting at radius rMin and stopping at rMax.
// Returns (rMax, false) if the field does not change sign over that span.
func SurfaceR(s *solid.Solid, deg, z, rMin, rMax float64) (float64, bool) {
	d, ok := SurfaceAlong(s, Cyl(rMin, deg, z), Cyl(1, deg, 0), rMax-rMin, 0)
	if !ok {
		return rMax, false
	}
	return rMin + d, true
}

// SurfaceZ returns the height (mm) of the first surface met marching down
// (-Z) from (x, y, zTop) to zBottom — the top surface of the solid at that
// plan-view point. Returns (zBottom, false) if no surface is crossed.
func SurfaceZ(s *solid.Solid, x, y, zTop, zBottom float64) (float64, bool) {
	d, ok := SurfaceAlong(s, v3.XYZ(x, y, zTop), v3.Z(-1), zTop-zBottom, 0)
	if !ok {
		return zBottom, false
	}
	return zTop - d, true
}

// Opening measures the clear void span (mm) through p along ±dir: p must be
// in air (At >= 0), and the result is the distance between the two material
// surfaces found within max mm either side. Returns (0, false) if p is
// inside material or if either side finds no surface within max.
func Opening(s *solid.Solid, p, dir v3.Vec, max float64) (float64, bool) {
	if Inside(s, p) {
		return 0, false
	}
	return surfSpan(s, p, dir, max, 0)
}

// WallThickness measures the material run (mm) through p along ±dir: p must
// be inside material (At < 0), and the result is the distance between the
// two surfaces found within max mm either side. Returns (0, false) if p is
// in air or if either side finds no surface within max.
func WallThickness(s *solid.Solid, p, dir v3.Vec, max float64) (float64, bool) {
	if !Inside(s, p) {
		return 0, false
	}
	return surfSpan(s, p, dir, max, 0)
}

// WallReport is the thinnest material run MinWall found: its thickness in
// mm, the interior sample point it was measured through, the axis it was
// measured along, and how many interior points were sampled.
type WallReport struct {
	Min     float64 // thinnest run found, mm; +Inf when nothing was measured
	At      v3.Vec  // the interior grid point it was measured through
	Dir     v3.Vec  // the axis it was measured along (+X, +Y or +Z)
	Samples int     // interior grid points probed (points in air are not counted)
}

// String returns the one-line report, e.g.
// "min wall = 1.983 mm at (0.000, 0.000, 5.000) along Z (412 samples)".
func (r WallReport) String() string {
	if r.Samples == 0 {
		return "min wall: no interior samples in region"
	}
	if math.IsInf(r.Min, 1) {
		return fmt.Sprintf("min wall: no run measured in %d interior samples", r.Samples)
	}
	return fmt.Sprintf("min wall = %.3f mm at (%.3f, %.3f, %.3f) along %s (%d samples)",
		r.Min, r.At.X, r.At.Y, r.At.Z, surfAxisName(r.Dir), r.Samples)
}

// MinWall scans region on a grid of `step` mm and reports the thinnest
// material run found through any interior grid point along each of the three
// axes, searching at most max mm either side. Rays are marched at half the
// grid step (capped at 0.1 mm), so this is a coarse but cheap screening
// check for accidental thin skins, not a true minimum-wall solver.
func MinWall(s *solid.Solid, region v3.Box, step, max float64) WallReport {
	pts := GridPoints(region, step)
	march := step / 2
	if march > surfDefaultStep {
		march = surfDefaultStep
	}
	dirs := [3]v3.Vec{v3.X(1), v3.Y(1), v3.Z(1)}
	type surfHit struct {
		d        float64
		at, dir  v3.Vec
		interior bool
		found    bool
	}
	hits := make([]surfHit, len(pts))
	parallel(len(pts), func(i int) {
		p := pts[i]
		hits[i].d = math.Inf(1)
		if !Inside(s, p) {
			return
		}
		hits[i].interior = true
		for _, d := range dirs {
			w, ok := surfSpan(s, p, d, max, march)
			if ok && w < hits[i].d {
				hits[i].d, hits[i].at, hits[i].dir, hits[i].found = w, p, d, true
			}
		}
	})
	out := WallReport{Min: math.Inf(1)}
	for _, h := range hits {
		if !h.interior {
			continue
		}
		out.Samples++
		if h.found && h.d < out.Min {
			out.Min, out.At, out.Dir = h.d, h.at, h.dir
		}
	}
	return out
}

// RequireMinWall runs MinWall over region and calls t.Errorf if the thinnest
// material run found is below wantMin mm, or if the region contained no
// interior sample at all (which would otherwise pass vacuously).
func RequireMinWall(t testing.TB, s *solid.Solid, region v3.Box, step, max, wantMin float64) {
	t.Helper()
	r := MinWall(s, region, step, max)
	if r.Samples == 0 || r.Min < wantMin {
		t.Errorf("[FAIL] %s (want >= %.3f mm)", r, wantMin)
		return
	}
	t.Logf("[ok] %s (want >= %.3f mm)", r, wantMin)
}

// RequireNear logs "name = got (want w ± tol)" and calls t.Errorf when
// |got-want| > tol — the standard dimension assertion for measurements read
// back out of a solid with SurfaceR, SurfaceZ, Opening or WallThickness.
func RequireNear(t testing.TB, name string, got, want, tol float64) {
	t.Helper()
	if math.Abs(got-want) > tol || math.IsNaN(got) {
		t.Errorf("[FAIL] %s = %.3f (want %.3f ± %.3f)", name, got, want, tol)
		return
	}
	t.Logf("[ok] %s = %.3f (want %.3f ± %.3f)", name, got, want, tol)
}

// surfBisect narrows [lo, hi] — lo on the in0 side of the surface, hi on the
// other — to surfBisectTol and returns the midpoint distance from `from`.
func surfBisect(s *solid.Solid, from, dir v3.Vec, lo, hi float64, in0 bool) float64 {
	for i := 0; i < 200 && hi-lo > surfBisectTol; i++ {
		mid := 0.5 * (lo + hi)
		if (At(s, from.Add(dir.MulScalar(mid))) < 0) == in0 {
			lo = mid
		} else {
			hi = mid
		}
	}
	return 0.5 * (lo + hi)
}

// surfSpan returns the distance between the surfaces found marching from p
// along +dir and -dir, at the given march step (0 = default).
func surfSpan(s *solid.Solid, p, dir v3.Vec, max, step float64) (float64, bool) {
	a, okA := SurfaceAlong(s, p, dir, max, step)
	b, okB := SurfaceAlong(s, p, dir.Neg(), max, step)
	if !okA || !okB {
		return 0, false
	}
	return a + b, true
}

// surfAxisName names an axis direction for a report line.
func surfAxisName(d v3.Vec) string {
	switch {
	case d.X != 0 && d.Y == 0 && d.Z == 0:
		return "X"
	case d.Y != 0 && d.X == 0 && d.Z == 0:
		return "Y"
	case d.Z != 0 && d.X == 0 && d.Y == 0:
		return "Z"
	}
	return fmt.Sprintf("(%.3f, %.3f, %.3f)", d.X, d.Y, d.Z)
}
