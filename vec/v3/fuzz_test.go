package v3

import (
	"math"
	"testing"
)

func finite3(xs ...float64) bool {
	for _, x := range xs {
		if math.IsNaN(x) || math.IsInf(x, 0) {
			return false
		}
	}
	return true
}

// FuzzAddSubRoundTrip checks (a+b)-b ≈ a within fp tolerance.
func FuzzAddSubRoundTrip(f *testing.F) {
	f.Add(1.0, 2.0, 3.0, 4.0, 5.0, 6.0)
	f.Add(0.0, 0.0, 0.0, 0.0, 0.0, 0.0)
	f.Add(-1e9, 1e9, 1e-9, -1e-9, 7.0, -7.0)
	f.Fuzz(func(t *testing.T, ax, ay, az, bx, by, bz float64) {
		if !finite3(ax, ay, az, bx, by, bz) {
			t.Skip()
		}
		a := XYZ(ax, ay, az)
		b := XYZ(bx, by, bz)
		got := a.Add(b).Sub(b)
		mag := math.Abs(ax) + math.Abs(ay) + math.Abs(az) +
			math.Abs(bx) + math.Abs(by) + math.Abs(bz)
		tol := 1e-9 * (1 + mag)
		if math.Abs(got.X-a.X) > tol ||
			math.Abs(got.Y-a.Y) > tol ||
			math.Abs(got.Z-a.Z) > tol {
			t.Fatalf("(a+b)-b = %+v, want %+v (tol=%v)", got, a, tol)
		}
	})
}

// FuzzDotCommutative checks a·b == b·a exactly.
func FuzzDotCommutative(f *testing.F) {
	f.Add(1.0, 2.0, 3.0, 4.0, 5.0, 6.0)
	f.Add(0.0, 0.0, 0.0, 1.0, -1.0, 1.0)
	f.Fuzz(func(t *testing.T, ax, ay, az, bx, by, bz float64) {
		if !finite3(ax, ay, az, bx, by, bz) {
			t.Skip()
		}
		a := XYZ(ax, ay, az)
		b := XYZ(bx, by, bz)
		if a.Dot(b) != b.Dot(a) {
			t.Fatalf("Dot not commutative: a·b=%v, b·a=%v", a.Dot(b), b.Dot(a))
		}
	})
}

// FuzzCrossPerpendicular checks (a × b) · a == 0 and (a × b) · b == 0
// within fp tolerance — the defining property of the cross product.
func FuzzCrossPerpendicular(f *testing.F) {
	f.Add(1.0, 0.0, 0.0, 0.0, 1.0, 0.0)
	f.Add(2.0, 3.0, 4.0, -1.0, 5.0, 7.0)
	f.Add(0.0, 0.0, 0.0, 1.0, 2.0, 3.0)
	f.Fuzz(func(t *testing.T, ax, ay, az, bx, by, bz float64) {
		if !finite3(ax, ay, az, bx, by, bz) {
			t.Skip()
		}
		a := XYZ(ax, ay, az)
		b := XYZ(bx, by, bz)
		c := a.Cross(b)
		// Magnitudes that bound the rounding error.
		la := a.Length()
		lb := b.Length()
		lc := c.Length()
		if math.IsInf(la, 0) || math.IsInf(lb, 0) || math.IsInf(lc, 0) {
			t.Skip()
		}
		// Generous tolerance: dot ~ |c| * |a| * eps.
		tol := 1e-9 * (1 + lc*(la+lb))
		da := c.Dot(a)
		db := c.Dot(b)
		if math.Abs(da) > tol {
			t.Fatalf("(a×b)·a = %v, want ~0 (tol=%v) for a=%+v b=%+v", da, tol, a, b)
		}
		if math.Abs(db) > tol {
			t.Fatalf("(a×b)·b = %v, want ~0 (tol=%v) for a=%+v b=%+v", db, tol, a, b)
		}
	})
}

// FuzzCrossAntiCommutative checks a × b ≈ -(b × a) component-wise.
// Algebraic equality holds; we allow a small tolerance because the two
// expressions perform their subtractions in different orders, which can
// differ by 1 ULP under IEEE-754 rounding.
func FuzzCrossAntiCommutative(f *testing.F) {
	f.Add(1.0, 2.0, 3.0, 4.0, 5.0, 6.0)
	f.Fuzz(func(t *testing.T, ax, ay, az, bx, by, bz float64) {
		if !finite3(ax, ay, az, bx, by, bz) {
			t.Skip()
		}
		a := XYZ(ax, ay, az)
		b := XYZ(bx, by, bz)
		ab := a.Cross(b)
		ba := b.Cross(a)
		// Per-component magnitude bounds the rounding error.
		magX := math.Abs(ay*bz) + math.Abs(az*by)
		magY := math.Abs(az*bx) + math.Abs(ax*bz)
		magZ := math.Abs(ax*by) + math.Abs(ay*bx)
		tolX := 1e-12 * (1 + magX)
		tolY := 1e-12 * (1 + magY)
		tolZ := 1e-12 * (1 + magZ)
		if math.Abs(ab.X+ba.X) > tolX ||
			math.Abs(ab.Y+ba.Y) > tolY ||
			math.Abs(ab.Z+ba.Z) > tolZ {
			t.Fatalf("a×b=%+v, -(b×a)=%+v", ab, Vec{X: -ba.X, Y: -ba.Y, Z: -ba.Z})
		}
	})
}

// FuzzLengthNonNegative checks Length is non-negative and matches Length2.
func FuzzLengthNonNegative(f *testing.F) {
	f.Add(1.0, 2.0, 2.0)
	f.Add(0.0, 0.0, 0.0)
	f.Add(-1e150, 1e150, 0.0)
	f.Fuzz(func(t *testing.T, x, y, z float64) {
		if !finite3(x, y, z) {
			t.Skip()
		}
		v := XYZ(x, y, z)
		l := v.Length()
		l2 := v.Length2()
		if math.IsNaN(l) || l < 0 {
			t.Fatalf("Length(%+v) = %v", v, l)
		}
		if math.IsNaN(l2) || l2 < 0 {
			t.Fatalf("Length2(%+v) = %v", v, l2)
		}
		if !math.IsInf(l2, 0) {
			tol := 1e-9 * (1 + l2)
			if math.Abs(l*l-l2) > tol {
				t.Fatalf("Length^2=%v, Length2=%v", l*l, l2)
			}
		}
	})
}

// FuzzNormalize checks Normalize never panics and produces a unit-length
// vector for non-degenerate inputs.
func FuzzNormalize(f *testing.F) {
	f.Add(1.0, 2.0, 2.0)
	f.Add(0.0, 0.0, 0.0)
	f.Add(1e-300, 1e-300, 1e-300)
	f.Add(1e150, 1e150, 1e150)
	f.Fuzz(func(t *testing.T, x, y, z float64) {
		if !finite3(x, y, z) {
			t.Skip()
		}
		v := XYZ(x, y, z)
		defer func() {
			if rec := recover(); rec != nil {
				t.Fatalf("Normalize(%+v) panicked: %v", v, rec)
			}
		}()
		n := v.Normalize()
		l := v.Length()
		if l == 0 || math.IsInf(l, 0) || l < 1e-300 {
			return
		}
		nl := n.Length()
		if math.IsNaN(nl) {
			t.Fatalf("Normalize -> NaN length for %+v", v)
		}
		if math.Abs(nl-1) > 1e-9 {
			t.Fatalf("Normalize(%+v) |n|=%v, want ~1", v, nl)
		}
	})
}

// FuzzScalarOps checks AddScalar/SubScalar and MulScalar/DivScalar
// round-trip within fp tolerance.
func FuzzScalarOps(f *testing.F) {
	f.Add(1.0, 2.0, 3.0, 4.0)
	f.Add(0.0, 0.0, 0.0, 1.0)
	f.Add(-5.0, 7.0, -1.0, -2.0)
	f.Fuzz(func(t *testing.T, x, y, z, s float64) {
		if !finite3(x, y, z, s) {
			t.Skip()
		}
		v := XYZ(x, y, z)
		got := v.AddScalar(s).SubScalar(s)
		mag := math.Abs(x) + math.Abs(y) + math.Abs(z) + math.Abs(s)
		tol := 1e-9 * (1 + mag)
		if math.Abs(got.X-v.X) > tol ||
			math.Abs(got.Y-v.Y) > tol ||
			math.Abs(got.Z-v.Z) > tol {
			t.Fatalf("(v+s)-s = %+v, want %+v", got, v)
		}
		if s != 0 && math.Abs(s) > 1e-150 && math.Abs(s) < 1e150 {
			rt := v.MulScalar(s).DivScalar(s)
			tol2 := 1e-9 * (1 + math.Abs(x) + math.Abs(y) + math.Abs(z))
			if math.Abs(rt.X-v.X) > tol2 ||
				math.Abs(rt.Y-v.Y) > tol2 ||
				math.Abs(rt.Z-v.Z) > tol2 {
				t.Fatalf("(v*s)/s = %+v, want %+v (s=%v)", rt, v, s)
			}
		}
	})
}

// FuzzNeg checks -(-v) == v exactly.
func FuzzNeg(f *testing.F) {
	f.Add(1.0, 2.0, 3.0)
	f.Add(0.0, 0.0, 0.0)
	f.Fuzz(func(t *testing.T, x, y, z float64) {
		if !finite3(x, y, z) {
			t.Skip()
		}
		v := XYZ(x, y, z)
		nn := v.Neg().Neg()
		if nn.X != v.X || nn.Y != v.Y || nn.Z != v.Z {
			t.Fatalf("Neg.Neg(%+v) = %+v", v, nn)
		}
	})
}

// FuzzGetSet exercises Get/Set without going out of bounds; out-of-bounds
// indices are documented as invalid so we skip them.
func FuzzGetSet(f *testing.F) {
	f.Add(1.0, 2.0, 3.0, 0, 9.0)
	f.Add(0.0, 0.0, 0.0, 2, -7.0)
	f.Fuzz(func(t *testing.T, x, y, z float64, idx int, val float64) {
		if !finite3(x, y, z, val) {
			t.Skip()
		}
		if idx < 0 || idx > 2 {
			t.Skip()
		}
		v := XYZ(x, y, z)
		defer func() {
			if rec := recover(); rec != nil {
				t.Fatalf("Get/Set(%+v, idx=%d, val=%v) panicked: %v", v, idx, val, rec)
			}
		}()
		_ = v.Get(idx)
		v.Set(idx, val)
		if v.Get(idx) != val {
			t.Fatalf("Set/Get round-trip failed: idx=%d val=%v got=%v", idx, val, v.Get(idx))
		}
	})
}
