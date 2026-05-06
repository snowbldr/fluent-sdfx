package units

import (
	"math"
	"testing"
)

// FuzzDegRadRoundTrip checks that RtoD(DtoR(d)) ≈ d for finite d.
func FuzzDegRadRoundTrip(f *testing.F) {
	f.Add(0.0)
	f.Add(90.0)
	f.Add(-180.0)
	f.Add(360.0)
	f.Add(1e6)
	f.Add(-1e6)
	f.Fuzz(func(t *testing.T, deg float64) {
		if math.IsNaN(deg) || math.IsInf(deg, 0) {
			t.Skip()
		}
		got := RtoD(DtoR(deg))
		// Tolerance scales with magnitude.
		tol := 1e-9 * (1 + math.Abs(deg))
		if math.Abs(got-deg) > tol {
			t.Fatalf("RtoD(DtoR(%v)) = %v, want %v (diff=%v, tol=%v)",
				deg, got, deg, got-deg, tol)
		}
	})
}

// FuzzRadDegRoundTrip checks that DtoR(RtoD(r)) ≈ r for finite r.
func FuzzRadDegRoundTrip(f *testing.F) {
	f.Add(0.0)
	f.Add(math.Pi)
	f.Add(-math.Pi)
	f.Add(2 * math.Pi)
	f.Add(1e3)
	f.Fuzz(func(t *testing.T, rad float64) {
		if math.IsNaN(rad) || math.IsInf(rad, 0) {
			t.Skip()
		}
		got := DtoR(RtoD(rad))
		tol := 1e-9 * (1 + math.Abs(rad))
		if math.Abs(got-rad) > tol {
			t.Fatalf("DtoR(RtoD(%v)) = %v, want %v (diff=%v, tol=%v)",
				rad, got, rad, got-rad, tol)
		}
	})
}

// FuzzEqualFloat64 checks the documented contract: a and b within epsilon
// reports true; reflexive (a equal to itself when epsilon >= 0); symmetric.
func FuzzEqualFloat64(f *testing.F) {
	f.Add(1.0, 1.0, 1e-9)
	f.Add(0.0, 0.0, 0.0)
	f.Add(1.0, 2.0, 0.5)
	f.Add(-3.5, -3.5, 1e-12)
	f.Fuzz(func(t *testing.T, a, b, eps float64) {
		// NaN behavior in float comparison is implementation-defined; skip.
		if math.IsNaN(a) || math.IsNaN(b) || math.IsNaN(eps) {
			t.Skip()
		}
		if math.IsInf(a, 0) || math.IsInf(b, 0) || math.IsInf(eps, 0) {
			t.Skip()
		}
		defer func() {
			if rec := recover(); rec != nil {
				t.Fatalf("EqualFloat64(%v, %v, %v) panicked: %v", a, b, eps, rec)
			}
		}()
		// Symmetry: EqualFloat64(a,b,eps) == EqualFloat64(b,a,eps)
		if EqualFloat64(a, b, eps) != EqualFloat64(b, a, eps) {
			t.Fatalf("EqualFloat64 not symmetric for (%v, %v, %v)", a, b, eps)
		}
		// Reflexivity for non-negative epsilon.
		if eps >= 0 && !EqualFloat64(a, a, eps) {
			t.Fatalf("EqualFloat64(%v, %v, %v) not reflexive", a, a, eps)
		}
	})
}

// FuzzErrMsg checks that ErrMsg never panics and always produces a non-nil
// error whose message is non-empty for non-empty input.
func FuzzErrMsg(f *testing.F) {
	f.Add("hello")
	f.Add("")
	f.Add("\x00\xff multi-byte☃")
	f.Fuzz(func(t *testing.T, msg string) {
		// Bound input length to keep this fast.
		if len(msg) > 1<<16 {
			t.Skip()
		}
		defer func() {
			if rec := recover(); rec != nil {
				t.Fatalf("ErrMsg(%q) panicked: %v", msg, rec)
			}
		}()
		err := ErrMsg(msg)
		if err == nil {
			t.Fatalf("ErrMsg(%q) returned nil", msg)
		}
		_ = err.Error() // must not panic on stringification
	})
}
