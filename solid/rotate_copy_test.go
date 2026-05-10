package solid

import (
	"fmt"
	"math"
	"strings"
	"testing"

	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
)

// rotateZ3D returns p rotated CCW about Z by angleDeg.
func rotateZ3D(p v3.Vec, angleDeg float64) v3.Vec {
	a := angleDeg * math.Pi / 180
	c, s := math.Cos(a), math.Sin(a)
	return v3.XYZ(c*p.X-s*p.Y, s*p.X+c*p.Y, p.Z)
}

// assertNRotSymmetricZ samples the SDF at many points around the Z axis
// and verifies f(p) ≈ f(rotateZ_360/n(p)) — true N-fold symmetry. The
// fold-with-wrong-canonical bug shows up as missing chunks in some
// copies, which trips this invariant immediately.
func assertNRotSymmetricZ(t *testing.T, s *Solid, n, samplesPerRing int, label string) {
	t.Helper()
	step := 360.0 / float64(n)
	radii := []float64{0.5, 2.0, 4.0, 5.5, 7.0}
	zs := []float64{-1.0, 0.0, 1.0}
	const tol = 1e-6
	for _, z := range zs {
		for _, r := range radii {
			for i := 0; i < samplesPerRing; i++ {
				theta := 2 * math.Pi * float64(i) / float64(samplesPerRing)
				p := v3.XYZ(r*math.Cos(theta), r*math.Sin(theta), z)
				pr := rotateZ3D(p, step)
				f0 := s.SDF3.Evaluate(v3sdf.Vec(p))
				f1 := s.SDF3.Evaluate(v3sdf.Vec(pr))
				if math.Abs(f0-f1) > tol {
					t.Fatalf("%s: not %d-fold symmetric: r=%.2f θ=%.1f° z=%.1f f0=%v f1=%v Δ=%v",
						label, n, r, theta*180/math.Pi, z, f0, f1, f0-f1)
				}
			}
		}
	}
}

func TestRotateCopyZAutoRecenter(t *testing.T) {
	src := Sphere(1).TranslateX(5)
	for _, n := range []int{2, 3, 4, 6, 8} {
		sectorWidth := 360.0 / float64(n)
		for _, frac := range []float64{0, 0.25, 0.5, 0.75, 1.0, 2.5, -1.5} {
			angle := frac * sectorWidth
			res := src.RotateZ(angle).RotateCopyZ(n)
			assertNRotSymmetricZ(t, res, n, 16, fmt.Sprintf("n=%d at=%.1f°", n, angle))
		}
	}
}

func TestRotateCopyZEqualsRotateUnionZ(t *testing.T) {
	src := Sphere(1).TranslateX(5)
	const tol = 1e-9
	for _, n := range []int{3, 4, 6} {
		a := src.RotateCopyZ(n)
		b := src.RotateUnionZ(n, RotateZMatrix(360.0/float64(n)))
		for x := -8.0; x <= 8.0; x += 2.0 {
			for y := -8.0; y <= 8.0; y += 2.0 {
				for z := -2.0; z <= 2.0; z += 1.0 {
					p := v3sdf.Vec{X: x, Y: y, Z: z}
					da := a.SDF3.Evaluate(p)
					db := b.SDF3.Evaluate(p)
					if d := math.Abs(da - db); d > tol {
						t.Fatalf("n=%d p=%v: RotateCopyZ=%v RotateUnionZ=%v Δ=%v", n, p, da, db, d)
					}
				}
			}
		}
	}
}

func TestRotateCopyZWrapStraddle(t *testing.T) {
	// Regression for the SawTooth-wrap failure in 3D: pre-fix,
	// `src.RotateZ(-150).RotateCopyZ(3)` clipped the copy whose
	// XY bbox straddled the ±180° period boundary.
	src := Sphere(1).TranslateX(5)
	for _, atDeg := range []float64{-150, -180, 180, 179.5, -179.5, 60, -60} {
		res := src.RotateZ(atDeg).RotateCopyZ(3)
		assertNRotSymmetricZ(t, res, 3, 24, fmt.Sprintf("wrap@%.1f°", atDeg))
	}
}

func TestRotateCopyAtZPlacesCopies(t *testing.T) {
	src := Sphere(1).TranslateX(5)
	const r = 5.0
	for _, n := range []int{3, 4, 6} {
		for _, atDeg := range []float64{0, 30, -150, 90, 359} {
			res := src.RotateCopyAtZ(n, atDeg)
			for k := 0; k < n; k++ {
				angle := atDeg + 360.0/float64(n)*float64(k)
				p := rotateZ3D(v3.XYZ(r, 0, 0), angle)
				if v := res.SDF3.Evaluate(v3sdf.Vec(p)); v > 1e-9 {
					t.Errorf("n=%d at=%v k=%d: expected inside, got f(%v)=%v",
						n, atDeg, k, p, v)
				}
				midAngle := angle + 180.0/float64(n)
				pm := rotateZ3D(v3.XYZ(r, 0, 0), midAngle)
				if v := res.SDF3.Evaluate(v3sdf.Vec(pm)); v < 0.5 {
					t.Errorf("n=%d at=%v k=%d: expected outside at mid, got f(%v)=%v",
						n, atDeg, k, pm, v)
				}
			}
		}
	}
}

func TestRotateCopyZWideSourcePanics(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for wide source")
		}
		msg := panicMsgZ(r)
		if !strings.Contains(msg, "RotateUnionZ") {
			t.Errorf("panic should point at RotateUnionZ: %v", msg)
		}
		if !strings.Contains(msg, "exceeds sector") {
			t.Errorf("panic should explain the cause: %v", msg)
		}
	}()
	// 2×2×2 box at (0.5, 0, 0): XY bbox crosses −X half-plane, angular
	// extent ≈ 233° — far wider than the 90° sector for n=4.
	src := Box(v3.XYZ(2, 2, 2), 0).TranslateX(0.5)
	_ = src.RotateCopyZ(4)
}

func TestRotateCopyAtZWideSourcePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for wide source")
		}
	}()
	src := Box(v3.XYZ(2, 2, 2), 0).TranslateX(0.5)
	_ = src.RotateCopyAtZ(4, 30)
}

func TestRotateCopyZCenteredOnAxis(t *testing.T) {
	// Cylinder (rotationally symmetric about Z): bbox center on the Z
	// axis triggers the fast path. Result must be valid and N-fold
	// symmetric.
	res := Cylinder(2, 3, 0).RotateCopyZ(6)
	if res == nil {
		t.Fatal("nil")
	}
	assertNRotSymmetricZ(t, res, 6, 24, "centered")
}

func TestRotateCopyZN1IsIdentity(t *testing.T) {
	src := Box(v3.XYZ(4, 2, 1), 0).TranslateX(3)
	if src.RotateCopyZ(1) == nil {
		t.Fatal("nil")
	}
}

func TestRotateCopyZZeroOrNegativePanics(t *testing.T) {
	for _, n := range []int{0, -1, -3} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("RotateCopyZ(%d) did not panic", n)
				}
			}()
			_ = Sphere(1).RotateCopyZ(n)
		}()
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("RotateCopyAtZ(%d, 0) did not panic", n)
				}
			}()
			_ = Sphere(1).RotateCopyAtZ(n, 0)
		}()
	}
}

func panicMsgZ(r any) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
