package v2

import (
	"math"
	"testing"
)

func finite2(xs ...float64) bool {
	for _, x := range xs {
		if math.IsNaN(x) || math.IsInf(x, 0) {
			return false
		}
	}
	return true
}

// FuzzAddSubRoundTrip checks (a+b)-b ≈ a within fp tolerance.
func FuzzAddSubRoundTrip(f *testing.F) {
	f.Add(1.0, 2.0, 3.0, 4.0)
	f.Add(0.0, 0.0, 0.0, 0.0)
	f.Add(-1e9, 1e9, 1e-9, -1e-9)
	f.Fuzz(func(t *testing.T, ax, ay, bx, by float64) {
		if !finite2(ax, ay, bx, by) {
			t.Skip()
		}
		a := XY(ax, ay)
		b := XY(bx, by)
		got := a.Add(b).Sub(b)
		mag := math.Abs(ax) + math.Abs(ay) + math.Abs(bx) + math.Abs(by)
		tol := 1e-9 * (1 + mag)
		if math.Abs(got.X-a.X) > tol || math.Abs(got.Y-a.Y) > tol {
			t.Fatalf("(a+b)-b = %+v, want %+v (tol=%v)", got, a, tol)
		}
	})
}

// FuzzDotCommutative checks a·b == b·a (exact: same fp ops).
func FuzzDotCommutative(f *testing.F) {
	f.Add(1.0, 2.0, 3.0, 4.0)
	f.Add(0.0, 0.0, 1.0, -1.0)
	f.Fuzz(func(t *testing.T, ax, ay, bx, by float64) {
		if !finite2(ax, ay, bx, by) {
			t.Skip()
		}
		a := XY(ax, ay)
		b := XY(bx, by)
		ab := a.Dot(b)
		ba := b.Dot(a)
		if ab != ba {
			t.Fatalf("Dot not commutative: a·b=%v, b·a=%v", ab, ba)
		}
	})
}

// FuzzCrossAntiCommutative checks a×b ≈ -(b×a) for the 2D scalar cross.
// Equality holds algebraically; we allow a small relative tolerance because
// the two expressions perform the subtraction in different orders, which
// can differ by 1 ULP under IEEE-754 rounding.
func FuzzCrossAntiCommutative(f *testing.F) {
	f.Add(1.0, 0.0, 0.0, 1.0)
	f.Add(2.0, 3.0, -4.0, 5.0)
	f.Fuzz(func(t *testing.T, ax, ay, bx, by float64) {
		if !finite2(ax, ay, bx, by) {
			t.Skip()
		}
		a := XY(ax, ay)
		b := XY(bx, by)
		ab := a.Cross(b)
		ba := b.Cross(a)
		mag := math.Abs(ax*by) + math.Abs(ay*bx)
		tol := 1e-12 * (1 + mag)
		if math.Abs(ab+ba) > tol {
			t.Fatalf("Cross not anti-commutative: a×b=%v, -(b×a)=%v (diff=%v, tol=%v)",
				ab, -ba, ab+ba, tol)
		}
	})
}

// FuzzLengthNonNegative checks Length() never returns NaN/negative for
// finite input, and matches sqrt(Length2()).
func FuzzLengthNonNegative(f *testing.F) {
	f.Add(3.0, 4.0)
	f.Add(0.0, 0.0)
	f.Add(-1e150, 1e150)
	f.Fuzz(func(t *testing.T, x, y float64) {
		if !finite2(x, y) {
			t.Skip()
		}
		v := XY(x, y)
		l := v.Length()
		if math.IsNaN(l) || l < 0 {
			t.Fatalf("Length(%+v) = %v, want non-negative", v, l)
		}
		l2 := v.Length2()
		if math.IsNaN(l2) || l2 < 0 {
			t.Fatalf("Length2(%+v) = %v, want non-negative", v, l2)
		}
		// l^2 ≈ l2 (allow some fp slop, skip on overflow).
		if !math.IsInf(l2, 0) {
			diff := math.Abs(l*l - l2)
			tol := 1e-9 * (1 + l2)
			if diff > tol {
				t.Fatalf("Length^2=%v, Length2=%v (diff=%v, tol=%v)", l*l, l2, diff, tol)
			}
		}
	})
}

// FuzzNormalize checks Normalize never panics and that the result has unit
// length when the input is non-zero. Zero-vector behavior is just probed
// for a non-panic.
func FuzzNormalize(f *testing.F) {
	f.Add(3.0, 4.0)
	f.Add(0.0, 0.0)        // zero vector — must not panic
	f.Add(1e-300, 1e-300)  // tiny
	f.Add(1e150, 1e150)    // large
	f.Fuzz(func(t *testing.T, x, y float64) {
		if !finite2(x, y) {
			t.Skip()
		}
		v := XY(x, y)
		defer func() {
			if rec := recover(); rec != nil {
				t.Fatalf("Normalize(%+v) panicked: %v", v, rec)
			}
		}()
		n := v.Normalize()
		// If input was non-zero and not so large/small that we underflow/overflow,
		// |n| ≈ 1.
		l := v.Length()
		if l == 0 || math.IsInf(l, 0) || l < 1e-300 {
			return
		}
		nl := n.Length()
		if math.IsNaN(nl) {
			t.Fatalf("Normalize produced NaN length for %+v", v)
		}
		if math.Abs(nl-1) > 1e-9 {
			t.Fatalf("Normalize(%+v) length = %v, want ~1", v, nl)
		}
	})
}

// FuzzScalarOps checks AddScalar/SubScalar/MulScalar/DivScalar invariants.
func FuzzScalarOps(f *testing.F) {
	f.Add(1.0, 2.0, 3.0)
	f.Add(0.0, 0.0, 1.0)
	f.Add(-5.0, 7.0, -2.0)
	f.Fuzz(func(t *testing.T, x, y, s float64) {
		if !finite2(x, y, s) {
			t.Skip()
		}
		v := XY(x, y)
		// (v + s) - s ≈ v
		got := v.AddScalar(s).SubScalar(s)
		mag := math.Abs(x) + math.Abs(y) + math.Abs(s)
		tol := 1e-9 * (1 + mag)
		if math.Abs(got.X-v.X) > tol || math.Abs(got.Y-v.Y) > tol {
			t.Fatalf("(v+s)-s = %+v, want %+v", got, v)
		}
		// (v * s) / s ≈ v when s != 0 and not extreme.
		if s != 0 && math.Abs(s) > 1e-150 && math.Abs(s) < 1e150 {
			rt := v.MulScalar(s).DivScalar(s)
			tol2 := 1e-9 * (1 + math.Abs(x) + math.Abs(y))
			if math.Abs(rt.X-v.X) > tol2 || math.Abs(rt.Y-v.Y) > tol2 {
				t.Fatalf("(v*s)/s = %+v, want %+v (s=%v)", rt, v, s)
			}
		}
	})
}

// FuzzNeg checks -(-v) == v exactly.
func FuzzNeg(f *testing.F) {
	f.Add(1.0, 2.0)
	f.Add(0.0, 0.0)
	f.Fuzz(func(t *testing.T, x, y float64) {
		if !finite2(x, y) {
			t.Skip()
		}
		v := XY(x, y)
		nn := v.Neg().Neg()
		if nn.X != v.X || nn.Y != v.Y {
			t.Fatalf("Neg.Neg(%+v) = %+v", v, nn)
		}
	})
}

// FuzzMinMax checks Min/Max produce per-component min/max and that
// Min(a,b) <= Max(a,b) component-wise.
func FuzzMinMax(f *testing.F) {
	f.Add(1.0, 2.0, 3.0, 4.0)
	f.Add(-1.0, 5.0, 5.0, -1.0)
	f.Fuzz(func(t *testing.T, ax, ay, bx, by float64) {
		if !finite2(ax, ay, bx, by) {
			t.Skip()
		}
		a := XY(ax, ay)
		b := XY(bx, by)
		mn := a.Min(b)
		mx := a.Max(b)
		if mn.X > mx.X || mn.Y > mx.Y {
			t.Fatalf("Min %+v > Max %+v", mn, mx)
		}
		if mn.X != math.Min(ax, bx) || mn.Y != math.Min(ay, by) {
			t.Fatalf("Min wrong: got %+v", mn)
		}
		if mx.X != math.Max(ax, bx) || mx.Y != math.Max(ay, by) {
			t.Fatalf("Max wrong: got %+v", mx)
		}
	})
}
