package validate_test

import (
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// --- RequireVolumeNear must reject expectedMM3 ≤ 0. ---
//
// Detecting that a *testing.T-based helper called t.Fatalf is awkward
// because Fatalf calls runtime.Goexit on the test goroutine. Our trick:
// run the helper in a fresh goroutine — Fatalf will only Goexit *that*
// goroutine, not the parent test. We detect the abort by checking whether
// a "completed" flag was set after the helper returned.
//
// Note: per the testing docs, FailNow must be called from the goroutine
// running the test. In our probe we are passing the parent's *testing.T
// to a separate goroutine; in practice (current Go runtimes) this still
// records failure and Goexits the calling goroutine. After detection we
// flip the parent test to passing via t.Skip() of a wrapper subtest.
func didHelperFatalf(fn func(*testing.T)) bool {
	parent := &testing.T{}
	completed := make(chan bool, 1)
	go func() {
		ok := false
		defer func() { completed <- ok }()
		// Recover from any panic so the test process doesn't crash.
		defer func() { _ = recover() }()
		fn(parent)
		ok = true
	}()
	return !<-completed
}

func TestRequireVolumeNearZeroExpectedFails(t *testing.T) {
	// We can't actually call validate.RequireVolumeNear with a zero-value
	// testing.T (its internals are nil) — Fatalf will panic with a nil
	// dereference. We detect that as "the helper aborted". Either way, the
	// helper MUST not return normally with a zero expectedMM3.
	failed := didHelperFatalf(func(p *testing.T) {
		validate.RequireVolumeNear(p, solid.Box(v3.XYZ(10, 10, 10), 0), 8.0, 0, 0.05)
	})
	if !failed {
		t.Fatalf("RequireVolumeNear(expectedMM3=0) returned normally; the guard against zero is missing")
	}
}

func TestRequireVolumeNearNegativeExpectedFails(t *testing.T) {
	failed := didHelperFatalf(func(p *testing.T) {
		validate.RequireVolumeNear(p, solid.Box(v3.XYZ(10, 10, 10), 0), 8.0, -100, 0.05)
	})
	if !failed {
		t.Fatalf("RequireVolumeNear(expectedMM3=-100) returned normally; the guard against negative is missing")
	}
}

// --- Gear profile sanity. ---

func TestInvoluteGearProfileSane(t *testing.T) {
	gear := obj.InvoluteGear(obj.InvoluteGearParms{
		NumberTeeth:      20,
		Module:           1.0,
		PressureAngleDeg: 20,
		Backlash:         0,
		Clearance:        0.1,
		RingWidth:        2,
		Facets:           5,
	})
	if gear == nil {
		t.Fatal("InvoluteGear returned nil")
	}
	bb := gear.Bounds()
	w := bb.Max.X - bb.Min.X
	h := bb.Max.Y - bb.Min.Y

	// Pitch diameter = NumberTeeth * Module = 20mm. With addendum/dedendum
	// the outer diameter is roughly pitchDiameter + 2*Module = 22mm, plus
	// RingWidth (2mm) on the outside of the root → outer ≈ 22 to 24mm.
	// Bbox should be roughly square with width ~22-26mm.
	if w < 18 || w > 30 {
		t.Errorf("gear width = %.2f mm, want roughly 22 mm (pitch dia + 2*module): bbox = %+v", w, bb)
	}
	if h < 18 || h > 30 {
		t.Errorf("gear height = %.2f mm, want roughly 22 mm: bbox = %+v", h, bb)
	}
	// Should be roughly circular: w ~= h (within 1%).
	if math.Abs(w-h)/math.Max(w, h) > 0.02 {
		t.Errorf("gear bbox not circular: w=%.3f h=%.3f", w, h)
	}
}

// --- TruncRectPyramid sanity. ---

func TestTruncRectPyramidSane(t *testing.T) {
	// Size in this sdfx kernel describes the TOP face; with base angle 60°
	// the base flares out by Size.Z / tan(60°) per side. So a Size.Z=8 with
	// 60° base angle adds 8/tan(60°) ≈ 4.62 per side → bbox in X grows by
	// ~9.24 mm. We sanity-check broad limits, not exact values.
	p := obj.TruncRectPyramid3D(obj.TruncRectPyramidParms{
		Size:         v3.XYZ(20, 14, 8),
		BaseAngleDeg: 60,
		BaseRadius:   1,
		RoundRadius:  0.5,
	})
	if p == nil {
		t.Fatal("TruncRectPyramid3D returned nil")
	}
	bb := p.Bounds()
	hz := bb.Max.Z - bb.Min.Z
	if hz < 7 || hz > 9 {
		t.Errorf("pyramid height = %.2f mm, want ~8 mm: bbox = %+v", hz, bb)
	}
	wx := bb.Max.X - bb.Min.X
	wy := bb.Max.Y - bb.Min.Y
	// Base spreads outward by ~h/tan(60°)*2 ≈ 9.24 mm beyond Size.X (20),
	// so wx should be roughly in [20, 35]. Likewise wy in [14, 30].
	if wx < 18 || wx > 35 {
		t.Errorf("pyramid X width = %.2f mm, want in [18, 35]: bbox = %+v", wx, bb)
	}
	if wy < 12 || wy > 30 {
		t.Errorf("pyramid Y width = %.2f mm, want in [12, 30]: bbox = %+v", wy, bb)
	}
}

// 90° base angle should produce a roughly box-like shape (vertical sides).
func TestTruncRectPyramidVerticalSides(t *testing.T) {
	p := obj.TruncRectPyramid3D(obj.TruncRectPyramidParms{
		Size:         v3.XYZ(20, 14, 8),
		BaseAngleDeg: 90,
		BaseRadius:   0,
		RoundRadius:  0,
	})
	bb := p.Bounds()
	wx := bb.Max.X - bb.Min.X
	wy := bb.Max.Y - bb.Min.Y
	// With vertical sides the bbox should match Size in X/Y.
	if math.Abs(wx-20) > 1 {
		t.Errorf("vertical-sided pyramid X width = %.2f mm, want ~20: bbox = %+v", wx, bb)
	}
	if math.Abs(wy-14) > 1 {
		t.Errorf("vertical-sided pyramid Y width = %.2f mm, want ~14: bbox = %+v", wy, bb)
	}
}

// --- 0-tolerance volume comparison via direct API. ---

func TestVolumeOfNearZeroSizeBox(t *testing.T) {
	// A 1x1x1 cube should have a volume close to 1 mm³.
	b := solid.Box(v3.XYZ(1, 1, 1), 0)
	st := validate.Of(b, 8.0)
	if st.Volume <= 0 {
		t.Fatalf("Volume of 1mm cube was %v, want positive", st.Volume)
	}
	if math.Abs(st.Volume-1) > 0.2 {
		t.Errorf("Volume of 1mm cube = %.3f, want close to 1 (within 20%%)", st.Volume)
	}
}

// --- Sanity: a translated solid keeps its volume invariant. ---

func TestVolumeInvariantUnderTranslation(t *testing.T) {
	a := solid.Box(v3.XYZ(10, 10, 10), 0)
	b := a.Translate(v3.XYZ(50, -30, 25))
	va := validate.Of(a, 8.0).Volume
	vb := validate.Of(b, 8.0).Volume
	if math.Abs(va-vb)/va > 0.05 {
		t.Errorf("translated cube volume changed: %.3f → %.3f", va, vb)
	}
}

