package shape

import (
	"fmt"
	"math"
	"strings"
	"testing"

	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v2sdf "github.com/snowbldr/sdfx/vec/v2"
)

// rotate2D returns p rotated CCW by angleDeg about the origin.
func rotate2D(p v2.Vec, angleDeg float64) v2.Vec {
	a := angleDeg * math.Pi / 180
	c, s := math.Cos(a), math.Sin(a)
	return v2.XY(c*p.X-s*p.Y, s*p.X+c*p.Y)
}

// assertNRotSymmetric samples the SDF at many points around the origin
// and verifies that f(p) ≈ f(rotate_360/n(p)) — i.e., the union has true
// N-fold rotational symmetry. The fold-with-wrong-canonical bug
// manifests as missing chunks in some copies, which breaks this
// invariant immediately.
func assertNRotSymmetric(t *testing.T, s *Shape, n, samplesPerRing int, label string) {
	t.Helper()
	step := 360.0 / float64(n)
	radii := []float64{0.5, 2.0, 4.0, 5.5, 7.0}
	const tol = 1e-6
	for _, r := range radii {
		for i := 0; i < samplesPerRing; i++ {
			theta := 2 * math.Pi * float64(i) / float64(samplesPerRing)
			p := v2.XY(r*math.Cos(theta), r*math.Sin(theta))
			pr := rotate2D(p, step)
			f0 := s.SDF2.Evaluate(v2sdf.Vec(p))
			f1 := s.SDF2.Evaluate(v2sdf.Vec(pr))
			if math.Abs(f0-f1) > tol {
				t.Fatalf("%s: not %d-fold symmetric: r=%.2f θ=%.1f° f0=%v f1=%v Δ=%v",
					label, n, r, theta*180/math.Pi, f0, f1, f0-f1)
			}
		}
	}
}

func TestRotateCopyAutoRecenter(t *testing.T) {
	// A small disk offset along +X — small angular extent, fits cleanly
	// inside any sector for n in {2..8}. Rotating to several offsets
	// inside (and outside) the canonical sector exercises the wrap line
	// at ±180° as well as off-axis cases.
	src := Circle(1).TranslateX(5)
	for _, n := range []int{2, 3, 4, 6, 8} {
		sectorWidth := 360.0 / float64(n)
		for _, frac := range []float64{0, 0.25, 0.5, 0.75, 1.0, 2.5, -1.5} {
			angle := frac * sectorWidth
			res := src.Rotate(angle).RotateCopy(n)
			assertNRotSymmetric(t, res, n, 16, fmt.Sprintf("n=%d at=%.1f°", n, angle))
		}
	}
}

func TestRotateCopyEqualsRotateUnion(t *testing.T) {
	// For a canonical source that fits cleanly inside the sector, the
	// folded RotateCopy and the explicit RotateUnion (with the matching
	// step) must produce the same distance field — proves the fix is
	// correct, not just plausible-looking.
	src := Circle(1).TranslateX(5)
	const tol = 1e-9
	for _, n := range []int{3, 4, 6} {
		a := src.RotateCopy(n)
		b := src.RotateUnion(n, Rotate2d(360.0/float64(n)))
		for x := -8.0; x <= 8.0; x += 1.5 {
			for y := -8.0; y <= 8.0; y += 1.5 {
				p := v2sdf.Vec{X: x, Y: y}
				da := a.SDF2.Evaluate(p)
				db := b.SDF2.Evaluate(p)
				if d := math.Abs(da - db); d > tol {
					t.Fatalf("n=%d p=%v: RotateCopy=%v RotateUnion=%v Δ=%v", n, p, da, db, d)
				}
			}
		}
	}
}

func TestRotateCopyWrapStraddle(t *testing.T) {
	// Regression for the SawTooth-wrap failure: a source rotated to an
	// angle where its bbox straddles the ±180° period boundary. Pre-fix,
	// `src.Rotate(-150).RotateCopy(3)` clipped the copy whose bbox
	// crossed the wrap line; with auto-recenter it produces three
	// intact copies.
	src := Circle(1).TranslateX(5)
	for _, atDeg := range []float64{-150, -180, 180, 179.5, -179.5, 60, -60} {
		res := src.Rotate(atDeg).RotateCopy(3)
		assertNRotSymmetric(t, res, 3, 24, fmt.Sprintf("wrap@%.1f°", atDeg))
	}
}

func TestRotateCopyAtPlacesCopies(t *testing.T) {
	// RotateCopyAt(n, atDeg) treats the source as canonical-centered
	// and rotates the union to atDeg. Sampling on the expected copy
	// centers (at angle atDeg + k·360/n, distance 5) must hit the
	// disk; sampling halfway between adjacent copies must miss it.
	src := Circle(1).TranslateX(5) // canonical: bbox center on +X
	const r = 5.0
	for _, n := range []int{3, 4, 6} {
		for _, atDeg := range []float64{0, 30, -150, 90, 359} {
			res := src.RotateCopyAt(n, atDeg)
			for k := 0; k < n; k++ {
				angle := atDeg + 360.0/float64(n)*float64(k)
				p := rotate2D(v2.XY(r, 0), angle)
				if v := res.SDF2.Evaluate(v2sdf.Vec(p)); v > 1e-9 {
					t.Errorf("n=%d at=%v k=%d: expected inside, got f(%v)=%v",
						n, atDeg, k, p, v)
				}
				midAngle := angle + 180.0/float64(n)
				pm := rotate2D(v2.XY(r, 0), midAngle)
				if v := res.SDF2.Evaluate(v2sdf.Vec(pm)); v < 0.5 {
					t.Errorf("n=%d at=%v k=%d: expected outside at mid, got f(%v)=%v",
						n, atDeg, k, pm, v)
				}
			}
		}
	}
}

func TestRotateCopyWideSourcePanics(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for wide source")
		}
		msg := panicMsg(r)
		if !strings.Contains(msg, "RotateUnion") {
			t.Errorf("panic should point at RotateUnion: %v", msg)
		}
		if !strings.Contains(msg, "exceeds sector") {
			t.Errorf("panic should explain the cause: %v", msg)
		}
	}()
	// A 2×2 box centered at (0.5, 0) crosses the −X half-plane, giving
	// a bbox angular extent ≈ 233° — far wider than the 90° sector
	// for n=4. Recentering can't help, so the wrapper must panic
	// rather than emit silently-wrong geometry.
	src := Rect(v2.XY(2, 2), 0).TranslateX(0.5)
	_ = src.RotateCopy(4)
}

func TestRotateCopyAtWideSourcePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for wide source")
		}
	}()
	// Offset Rect — bbox is off-axis enough to escape the
	// rotationally-symmetric fast path, but its angular extent at
	// (0.5, 0) crosses ±π and exceeds 90° (n=4 sector).
	src := Rect(v2.XY(2, 2), 0).TranslateX(0.5)
	_ = src.RotateCopyAt(4, 30)
}

func TestRotateCopyCenteredOnOrigin(t *testing.T) {
	// A bbox-centered-on-origin source (a plain Circle): atan2 of the
	// bbox center is meaningless, so isNearOrigin2D triggers the fast
	// path. The result must be valid and 6-fold symmetric.
	res := Circle(2).RotateCopy(6)
	if res == nil {
		t.Fatal("nil")
	}
	assertNRotSymmetric(t, res, 6, 24, "centered")
}

func TestRotateCopyN1IsIdentity(t *testing.T) {
	// n=1 means "one copy in a full circle" — a pass-through. The
	// extent check must not trip on a wide source for n=1.
	src := Rect(v2.XY(4, 2), 0).TranslateX(3)
	if src.RotateCopy(1) == nil {
		t.Fatal("nil")
	}
}

func TestRotateCopyZeroOrNegativePanics(t *testing.T) {
	for _, n := range []int{0, -1, -3} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("RotateCopy(%d) did not panic", n)
				}
			}()
			_ = Circle(1).RotateCopy(n)
		}()
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("RotateCopyAt(%d, 0) did not panic", n)
				}
			}()
			_ = Circle(1).RotateCopyAt(n, 0)
		}()
	}
}

func TestRotateCopySectorWedgeWithCaps(t *testing.T) {
	// Regression for the user's failure: a ~50° wedge with cap-circle
	// bulges (≈64° bbox angular extent), rotated to -150°, RotateCopy(3)
	// must produce three intact 3-fold-symmetric copies. Pre-fix, the
	// fold's ±180° wrap line cut a chunk out of the copy whose bbox
	// straddled it; post-fix, auto-recenter rotates to canonical first.
	//
	// The wedge here doesn't include the origin (offset start at X=4)
	// so the bbox is angularly tight — a vertex-at-origin source would
	// have bbox angular extent ≈180° and trip the unfixable guard,
	// which is the documented constraint for bbox-based estimation.
	const halfDeg = 25.0
	const inner = 8.0
	const outer = 10.0
	const cap = 0.5
	tipUp := v2.XY(outer*math.Cos(halfDeg*math.Pi/180), outer*math.Sin(halfDeg*math.Pi/180))
	tipDn := v2.XY(tipUp.X, -tipUp.Y)
	wedge := Polygon([]v2.Vec{
		v2.XY(inner, 0),
		tipUp,
		tipDn,
	})
	capUp := Circle(cap).Translate(tipUp)
	capDn := Circle(cap).Translate(tipDn)
	profile := wedge.Union(capUp, capDn)
	res := profile.Rotate(-150).RotateCopy(3)
	assertNRotSymmetric(t, res, 3, 32, "wedge-with-caps@-150°")
}

func panicMsg(r any) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
