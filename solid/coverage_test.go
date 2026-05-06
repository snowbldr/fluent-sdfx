package solid

import (
	"math"
	"strings"
	"testing"

	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/sdfx/sdf"
	v2sdf "github.com/snowbldr/sdfx/vec/v2"
)

// boxClose reports whether two bounding boxes match within tol on every component.
func boxClose(t *testing.T, got, want Box3, tol float64) bool {
	t.Helper()
	return math.Abs(got.Min.X-want.Min.X) <= tol &&
		math.Abs(got.Min.Y-want.Min.Y) <= tol &&
		math.Abs(got.Min.Z-want.Min.Z) <= tol &&
		math.Abs(got.Max.X-want.Max.X) <= tol &&
		math.Abs(got.Max.Y-want.Max.Y) <= tol &&
		math.Abs(got.Max.Z-want.Max.Z) <= tol
}

// floatClose reports whether |a-b| <= tol.
func floatClose(a, b, tol float64) bool { return math.Abs(a-b) <= tol }

// expectPanic runs fn and reports whether it panicked. msgFragment, if
// non-empty, must appear in the panic message.
func expectPanic(t *testing.T, msgFragment string, fn func()) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic containing %q, got none", msgFragment)
		}
		if msgFragment != "" {
			msg := ""
			switch v := r.(type) {
			case string:
				msg = v
			case error:
				msg = v.Error()
			default:
				msg = ""
			}
			if msg != "" && !strings.Contains(msg, msgFragment) {
				t.Fatalf("panic message %q does not contain %q", msg, msgFragment)
			}
		}
	}()
	fn()
}

// rect2D returns a centered rectangle profile as raw sdf.SDF2 — used to
// avoid the import cycle with the shape package.
func rect2D(w, h float64) sdf.SDF2 {
	return sdf.Box2D(v2sdf.Vec{X: w, Y: h}, 0)
}

// --- Constructors return non-nil with the expected bounds (regression: a
// silent shape change here would break every downstream test). ---

func TestConstructorsReturnExpectedBounds(t *testing.T) {
	t.Run("Box", func(t *testing.T) {
		b := Box(v3.XYZ(10, 6, 4), 0)
		if b == nil {
			t.Fatal("Box returned nil")
		}
		want := NewBox3(v3.XYZ(0, 0, 0), v3.XYZ(10, 6, 4))
		if !boxClose(t, b.Bounds(), want, 1e-9) {
			t.Fatalf("Box(10,6,4).Bounds() = %+v, want %+v", b.Bounds(), want)
		}
	})

	t.Run("Cylinder", func(t *testing.T) {
		c := Cylinder(8, 3, 0)
		if c == nil {
			t.Fatal("Cylinder returned nil")
		}
		bb := c.Bounds()
		// Z extent should be [-4, 4]; XY extent should be [-3, 3].
		if !floatClose(bb.Min.Z, -4, 1e-9) || !floatClose(bb.Max.Z, 4, 1e-9) {
			t.Errorf("Cylinder Z bounds = [%v, %v], want [-4, 4]", bb.Min.Z, bb.Max.Z)
		}
		if !floatClose(bb.Min.X, -3, 1e-9) || !floatClose(bb.Max.X, 3, 1e-9) {
			t.Errorf("Cylinder X bounds = [%v, %v], want [-3, 3]", bb.Min.X, bb.Max.X)
		}
	})

	t.Run("Sphere", func(t *testing.T) {
		s := Sphere(5)
		if s == nil {
			t.Fatal("Sphere returned nil")
		}
		bb := s.Bounds()
		// Sphere bbox should be [-5,5]^3 on every axis.
		if !floatClose(bb.Min.X, -5, 1e-9) || !floatClose(bb.Max.X, 5, 1e-9) {
			t.Errorf("Sphere X bounds = [%v, %v], want [-5, 5]", bb.Min.X, bb.Max.X)
		}
		if !floatClose(bb.Min.Z, -5, 1e-9) || !floatClose(bb.Max.Z, 5, 1e-9) {
			t.Errorf("Sphere Z bounds = [%v, %v], want [-5, 5]", bb.Min.Z, bb.Max.Z)
		}
	})

	t.Run("Cone", func(t *testing.T) {
		// Truncated cone: r0=4, r1=2, height=10.
		c := Cone(10, 4, 2, 0)
		if c == nil {
			t.Fatal("Cone returned nil")
		}
		bb := c.Bounds()
		if !floatClose(bb.Min.Z, -5, 1e-9) || !floatClose(bb.Max.Z, 5, 1e-9) {
			t.Errorf("Cone Z bounds = [%v, %v], want [-5, 5]", bb.Min.Z, bb.Max.Z)
		}
		// Bottom is wider; bbox should encompass r0=4 in X/Y.
		if !floatClose(bb.Min.X, -4, 1e-9) || !floatClose(bb.Max.X, 4, 1e-9) {
			t.Errorf("Cone X bounds = [%v, %v], want [-4, 4]", bb.Min.X, bb.Max.X)
		}
	})
}

// --- Boolean op no-op behavior (regressions from earlier audit). ---

func TestBooleanNoOpsReturnReceiver(t *testing.T) {
	s := Box(v3.XYZ(10, 10, 10), 0)
	want := s.Bounds()

	t.Run("Cut_no_args_returns_receiver", func(t *testing.T) {
		got := s.Cut().Bounds()
		if !boxClose(t, got, want, 1e-9) {
			t.Fatalf("Cut() bounds changed: got %+v want %+v", got, want)
		}
	})
	t.Run("Union_no_args_returns_receiver", func(t *testing.T) {
		got := s.Union().Bounds()
		if !boxClose(t, got, want, 1e-9) {
			t.Fatalf("Union() bounds changed: got %+v want %+v", got, want)
		}
	})
	t.Run("Intersect_no_args_returns_receiver", func(t *testing.T) {
		got := s.Intersect().Bounds()
		if !boxClose(t, got, want, 1e-9) {
			t.Fatalf("Intersect() bounds changed: got %+v want %+v", got, want)
		}
	})
	t.Run("Multi_no_positions_returns_receiver", func(t *testing.T) {
		got := s.Multi().Bounds()
		if !boxClose(t, got, want, 1e-9) {
			t.Fatalf("Multi() bounds changed: got %+v want %+v", got, want)
		}
	})
	t.Run("LineOf_empty_pattern_returns_receiver", func(t *testing.T) {
		got := s.LineOf(v3.XYZ(0, 0, 0), v3.XYZ(10, 0, 0), "").Bounds()
		if !boxClose(t, got, want, 1e-9) {
			t.Fatalf("LineOf empty pattern bounds changed: got %+v want %+v", got, want)
		}
	})
	t.Run("LineOf_all_skip_pattern_returns_receiver", func(t *testing.T) {
		got := s.LineOf(v3.XYZ(0, 0, 0), v3.XYZ(10, 0, 0), "...").Bounds()
		if !boxClose(t, got, want, 1e-9) {
			t.Fatalf("LineOf all-skip pattern bounds changed: got %+v want %+v", got, want)
		}
	})
}

// --- Smooth boolean edge cases. ---

func TestSmoothUnionPanicsOnEmpty(t *testing.T) {
	expectPanic(t, "at least one solid required", func() {
		SmoothUnion(RoundMin(1.0))
	})
}

func TestSmoothUnionSingleArgReturnsArg(t *testing.T) {
	a := Box(v3.XYZ(4, 4, 4), 0)
	got := SmoothUnion(RoundMin(1.0), a)
	if got != a {
		t.Fatalf("SmoothUnion with one solid should return that solid unchanged")
	}
}

func TestUnionAllPanicsOnEmpty(t *testing.T) {
	expectPanic(t, "at least one solid required", func() {
		UnionAll()
	})
}

func TestUnionAllSingleArgReturnsArg(t *testing.T) {
	a := Box(v3.XYZ(4, 4, 4), 0)
	got := UnionAll(a)
	if got != a {
		t.Fatalf("UnionAll with one solid should return that solid unchanged")
	}
}

func TestSmoothDifferenceNoToolsReturnsReceiver(t *testing.T) {
	s := Box(v3.XYZ(10, 10, 10), 0)
	got := SmoothDifference(PolyMax(1.0), s)
	if got != s {
		t.Fatalf("SmoothDifference with no tools should return s unchanged")
	}
}

func TestSmoothIntersectionNoOthersReturnsReceiver(t *testing.T) {
	s := Box(v3.XYZ(10, 10, 10), 0)
	got := SmoothIntersection(PolyMax(1.0), s)
	if got != s {
		t.Fatalf("SmoothIntersection with no others should return s unchanged")
	}
}

// --- Helix sweep edge case. ---

func TestSweepHelixZeroTurnsPanics(t *testing.T) {
	expectPanic(t, "turns must be > 0", func() {
		SweepHelix(rect2D(1, 1), 5, 0, 10, false)
	})
}

func TestSweepHelixZeroHeightPanics(t *testing.T) {
	expectPanic(t, "height must be > 0", func() {
		SweepHelix(rect2D(1, 1), 5, 1, 0, false)
	})
}

// --- TwistExtrude: 90° twist should NOT be interpreted as 14 turns
// (sdfx historically expects radians; if we forget the deg→rad conversion the
// resulting solid wraps around many times, ballooning the X/Y bounding box). ---

func TestTwistExtrude90DegMatchesRotated(t *testing.T) {
	// Use a tall, thin profile so a twist visibly enlarges the XY footprint —
	// but a 90° twist should stay within sqrt(profile_x²+profile_y²)/2 of origin.
	profile := rect2D(10, 2)
	twisted := TwistExtrude(profile, 5, 90).Bounds()

	// A correct 90° twist: starts at +/- (5, 1), rotates to +/- (1, 5).
	// Bounding box spans roughly [-5, 5] in both X and Y. A 14-turn (radians)
	// mistake would still fit inside this box (since twist is sweep, not scale),
	// but the X/Y extent of a 90° rotation linearly interpolating between
	// 0° and 90° should produce a bbox no larger than the rectangle of size
	// (10, 10) — i.e., the diagonal of the rotation envelope. Here we assert
	// the bbox doesn't blow far beyond that envelope.
	if twisted.Max.X > 5.5 {
		t.Errorf("TwistExtrude 90° X.Max = %v exceeds linear-twist bound 5.5 (maybe deg/rad mix-up)",
			twisted.Max.X)
	}
	if twisted.Min.X < -5.5 {
		t.Errorf("TwistExtrude 90° X.Min = %v exceeds linear-twist bound -5.5 (maybe deg/rad mix-up)",
			twisted.Min.X)
	}
	if twisted.Max.Y > 5.5 {
		t.Errorf("TwistExtrude 90° Y.Max = %v exceeds linear-twist bound 5.5 (maybe deg/rad mix-up)",
			twisted.Max.Y)
	}
}

// --- Transform composability: Translate(v).Translate(-v) should restore bounds. ---

func TestTranslateInverseRestoresBounds(t *testing.T) {
	s := Box(v3.XYZ(10, 10, 10), 0)
	v := v3.XYZ(7, -3, 11)
	got := s.Translate(v).Translate(v.Neg()).Bounds()
	want := s.Bounds()
	if !boxClose(t, got, want, 1e-9) {
		t.Fatalf("Translate(v).Translate(-v) bounds = %+v, want %+v", got, want)
	}
}

// 4× RotateZ(90) should round-trip an asymmetric solid.
func TestRotateZ360RoundTrip(t *testing.T) {
	s := Box(v3.XYZ(10, 4, 6), 0)
	got := s.RotateZ(90).RotateZ(90).RotateZ(90).RotateZ(90).Bounds()
	want := s.Bounds()
	if !boxClose(t, got, want, 1e-6) {
		t.Fatalf("4× RotateZ(90) bounds = %+v, want %+v", got, want)
	}
}

// --- Anchors: every one of the 27 should land on the right unit-cube
// coordinate of the bbox. ---

func TestAllSolidAnchorPositions(t *testing.T) {
	// A non-cube box translated so the center is non-trivial. This catches
	// any bug that confuses min/max/center.
	size := v3.XYZ(20, 14, 10)
	center := v3.XYZ(3, -2, 5)
	s := Box(size, 0).Translate(center)

	half := size.MulScalar(0.5)
	expect := func(x, y, z int) v3.Vec {
		return center.Add(v3.XYZ(float64(x)*half.X, float64(y)*half.Y, float64(z)*half.Z))
	}

	cases := []struct {
		name    string
		got     v3.Vec
		x, y, z int
	}{
		// faces
		{"Top", s.Top().Point, 0, 0, 1},
		{"Bottom", s.Bottom().Point, 0, 0, -1},
		{"Right", s.Right().Point, 1, 0, 0},
		{"Left", s.Left().Point, -1, 0, 0},
		{"Back", s.Back().Point, 0, 1, 0},
		{"Front", s.Front().Point, 0, -1, 0},
		// edges
		{"TopRight", s.TopRight().Point, 1, 0, 1},
		{"TopLeft", s.TopLeft().Point, -1, 0, 1},
		{"TopFront", s.TopFront().Point, 0, -1, 1},
		{"TopBack", s.TopBack().Point, 0, 1, 1},
		{"BottomRight", s.BottomRight().Point, 1, 0, -1},
		{"BottomLeft", s.BottomLeft().Point, -1, 0, -1},
		{"BottomFront", s.BottomFront().Point, 0, -1, -1},
		{"BottomBack", s.BottomBack().Point, 0, 1, -1},
		{"FrontRight", s.FrontRight().Point, 1, -1, 0},
		{"FrontLeft", s.FrontLeft().Point, -1, -1, 0},
		{"BackRight", s.BackRight().Point, 1, 1, 0},
		{"BackLeft", s.BackLeft().Point, -1, 1, 0},
		// corners
		{"TopFrontRight", s.TopFrontRight().Point, 1, -1, 1},
		{"TopFrontLeft", s.TopFrontLeft().Point, -1, -1, 1},
		{"TopBackRight", s.TopBackRight().Point, 1, 1, 1},
		{"TopBackLeft", s.TopBackLeft().Point, -1, 1, 1},
		{"BottomFrontRight", s.BottomFrontRight().Point, 1, -1, -1},
		{"BottomFrontLeft", s.BottomFrontLeft().Point, -1, -1, -1},
		{"BottomBackRight", s.BottomBackRight().Point, 1, 1, -1},
		{"BottomBackLeft", s.BottomBackLeft().Point, -1, 1, -1},
		// center via AnchorAt(0,0,0)
		{"AnchorAt(0,0,0)", s.AnchorAt(0, 0, 0).Point, 0, 0, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			want := expect(c.x, c.y, c.z)
			if !vecClose(c.got, want) {
				t.Fatalf("%s = %+v, want %+v (unit-cube coord (%d,%d,%d))",
					c.name, c.got, want, c.x, c.y, c.z)
			}
		})
	}
}
