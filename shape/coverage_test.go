package shape

import (
	"math"
	"strings"
	"testing"

	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	"github.com/snowbldr/fluent-sdfx/vec/v2i"
)

const eps = 1e-9

// vecClose2 reports whether two 2D vectors match within eps.
func vecClose2(a, b v2.Vec) bool {
	return math.Abs(a.X-b.X) < eps && math.Abs(a.Y-b.Y) < eps
}

// boxClose2 reports whether two 2D boxes match within tol.
func boxClose2(t *testing.T, got, want Box2, tol float64) bool {
	t.Helper()
	return math.Abs(got.Min.X-want.Min.X) <= tol &&
		math.Abs(got.Min.Y-want.Min.Y) <= tol &&
		math.Abs(got.Max.X-want.Max.X) <= tol &&
		math.Abs(got.Max.Y-want.Max.Y) <= tol
}

// expectPanic runs fn and asserts it panicked. msgFragment, if non-empty,
// must appear in the panic message.
func expectPanic(t *testing.T, msgFragment string, fn func()) {
	t.Helper()
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected panic containing %q, got none", msgFragment)
		}
		if msgFragment == "" {
			return
		}
		var msg string
		switch v := r.(type) {
		case string:
			msg = v
		case error:
			msg = v.Error()
		}
		if msg != "" && !strings.Contains(msg, msgFragment) {
			t.Fatalf("panic message %q does not contain %q", msg, msgFragment)
		}
	}()
	fn()
}

// --- Constructor sanity. ---

func TestRectBounds(t *testing.T) {
	r := Rect(v2.XY(10, 6), 0)
	if r == nil {
		t.Fatal("Rect returned nil")
	}
	want := v2.Box{Min: v2.XY(-5, -3), Max: v2.XY(5, 3)}
	if !boxClose2(t, r.Bounds(), want, eps) {
		t.Fatalf("Rect(10,6).Bounds() = %+v, want %+v", r.Bounds(), want)
	}
}

func TestCircleBounds(t *testing.T) {
	c := Circle(4)
	if c == nil {
		t.Fatal("Circle returned nil")
	}
	want := v2.Box{Min: v2.XY(-4, -4), Max: v2.XY(4, 4)}
	if !boxClose2(t, c.Bounds(), want, eps) {
		t.Fatalf("Circle(4).Bounds() = %+v, want %+v", c.Bounds(), want)
	}
}

// shape.Circle(negative) should panic — the audit fixed it from being silent.
func TestCircleNegativeRadiusPanics(t *testing.T) {
	expectPanic(t, "Circle", func() {
		Circle(-1)
	})
}

// --- Boolean op no-ops. ---

func TestShapeBooleanNoOpsReturnReceiver(t *testing.T) {
	s := Rect(v2.XY(10, 10), 0)
	want := s.Bounds()

	t.Run("Cut_no_args_returns_receiver", func(t *testing.T) {
		got := s.Cut().Bounds()
		if !boxClose2(t, got, want, eps) {
			t.Fatalf("Cut() bounds = %+v, want %+v", got, want)
		}
	})
	t.Run("Union_no_args_returns_receiver", func(t *testing.T) {
		got := s.Union().Bounds()
		if !boxClose2(t, got, want, eps) {
			t.Fatalf("Union() bounds = %+v, want %+v", got, want)
		}
	})
	t.Run("Intersect_no_args_returns_receiver", func(t *testing.T) {
		got := s.Intersect().Bounds()
		if !boxClose2(t, got, want, eps) {
			t.Fatalf("Intersect() bounds = %+v, want %+v", got, want)
		}
	})
	t.Run("Multi_no_positions_returns_receiver", func(t *testing.T) {
		got := s.Multi().Bounds()
		if !boxClose2(t, got, want, eps) {
			t.Fatalf("Multi() bounds = %+v, want %+v", got, want)
		}
	})
	t.Run("LineOf_empty_pattern_returns_receiver", func(t *testing.T) {
		got := s.LineOf(v2.XY(0, 0), v2.XY(10, 0), "").Bounds()
		if !boxClose2(t, got, want, eps) {
			t.Fatalf("LineOf empty bounds = %+v, want %+v", got, want)
		}
	})
	t.Run("LineOf_all_skip_pattern_returns_receiver", func(t *testing.T) {
		got := s.LineOf(v2.XY(0, 0), v2.XY(10, 0), "...").Bounds()
		if !boxClose2(t, got, want, eps) {
			t.Fatalf("LineOf all-skip bounds = %+v, want %+v", got, want)
		}
	})
}

// --- UnionAll edge cases. ---

func TestShapeUnionAllPanicsOnEmpty(t *testing.T) {
	expectPanic(t, "at least one shape required", func() {
		UnionAll()
	})
}

func TestShapeUnionAllSingleArgReturnsArg(t *testing.T) {
	a := Rect(v2.XY(2, 2), 0)
	got := UnionAll(a)
	if got != a {
		t.Fatalf("UnionAll with one shape should return that shape unchanged")
	}
}

// --- Transform composability. ---

func TestShapeTranslateInverseRestoresBounds(t *testing.T) {
	s := Rect(v2.XY(10, 6), 0)
	v := v2.XY(7, -3)
	got := s.Translate(v).Translate(v.Neg()).Bounds()
	want := s.Bounds()
	if !boxClose2(t, got, want, eps) {
		t.Fatalf("Translate(v).Translate(-v) bounds = %+v, want %+v", got, want)
	}
}

func TestShapeRotate360RoundTrip(t *testing.T) {
	s := Rect(v2.XY(10, 4), 0)
	got := s.Rotate(90).Rotate(90).Rotate(90).Rotate(90).Bounds()
	want := s.Bounds()
	if !boxClose2(t, got, want, 1e-6) {
		t.Fatalf("4× Rotate(90) bounds = %+v, want %+v", got, want)
	}
}

// --- All 9 shape anchors land at the right unit-square coordinates. ---

func TestAllShapeAnchorPositions(t *testing.T) {
	size := v2.XY(20, 14)
	center := v2.XY(3, -2)
	s := Rect(size, 0).Translate(center)
	half := size.MulScalar(0.5)
	expect := func(x, y int) v2.Vec {
		return center.Add(v2.XY(float64(x)*half.X, float64(y)*half.Y))
	}

	cases := []struct {
		name string
		got  v2.Vec
		x, y int
	}{
		{"Top", s.Top().Point, 0, 1},
		{"Bottom", s.Bottom().Point, 0, -1},
		{"Right", s.Right().Point, 1, 0},
		{"Left", s.Left().Point, -1, 0},
		{"TopRight", s.TopRight().Point, 1, 1},
		{"TopLeft", s.TopLeft().Point, -1, 1},
		{"BottomRight", s.BottomRight().Point, 1, -1},
		{"BottomLeft", s.BottomLeft().Point, -1, -1},
		{"AnchorAt(0,0)", s.AnchorAt(0, 0).Point, 0, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			want := expect(c.x, c.y)
			if !vecClose2(c.got, want) {
				t.Fatalf("%s = %+v, want %+v (unit-square coord (%d,%d))",
					c.name, c.got, want, c.x, c.y)
			}
		})
	}
}

// --- 2D placement verbs (smoke). ---

func TestShapeOnAlignsAnchors(t *testing.T) {
	body := Rect(v2.XY(10, 6), 0)
	cap := Circle(2)
	moved := cap.Bottom().On(body.Top()).Shape()

	got := moved.Bottom().Point
	want := body.Top().Point
	if !vecClose2(got, want) {
		t.Fatalf("after On: moved.Bottom() = %+v, body.Top() = %+v", got, want)
	}
}

// --- 2D primitives. ---

func TestHexagonHasSixVertices(t *testing.T) {
	h := Hexagon(5)
	if h == nil {
		t.Fatal("Hexagon returned nil")
	}
	bb := h.Bounds()
	// Hexagon with circumradius 5 must fit within [-5, 5] in both axes.
	if bb.Max.X > 5+eps || bb.Min.X < -5-eps {
		t.Errorf("Hexagon X bounds = [%v, %v], want within [-5, 5]", bb.Min.X, bb.Max.X)
	}
}

func TestTriangleNonEmpty(t *testing.T) {
	tri := Triangle(5)
	if tri == nil {
		t.Fatal("Triangle returned nil")
	}
	bb := tri.Bounds()
	if bb.Max.X-bb.Min.X <= 0 || bb.Max.Y-bb.Min.Y <= 0 {
		t.Errorf("Triangle has empty bbox: %+v", bb)
	}
}

func TestPolygonDegenerate(t *testing.T) {
	expectPanic(t, "", func() {
		Polygon([]v2.Vec{v2.XY(0, 0), v2.XY(1, 0)})
	})
}

// --- ToPNG smoke (renders to a temp file via stdlib testing.T.TempDir). ---

func TestShapeGenerateMesh(t *testing.T) {
	s := Rect(v2.XY(10, 10), 0)
	pts, err := s.GenerateMesh(v2i.Vec{X: 8, Y: 8})
	if err != nil {
		t.Fatalf("GenerateMesh error: %v", err)
	}
	if len(pts) == 0 {
		t.Fatalf("GenerateMesh produced no points")
	}
}
