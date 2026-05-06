package solid

import (
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/plane"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/sdfx/sdf"
	v2sdf "github.com/snowbldr/sdfx/vec/v2"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
)

// --- Translate variants ---

func TestTranslateVariants(t *testing.T) {
	s := Box(v3.XYZ(2, 2, 2), 0)
	cases := []struct {
		name string
		got  v3.Vec
		want v3.Vec
	}{
		{"TranslateX", s.TranslateX(5).Bounds().Center(), v3.XYZ(5, 0, 0)},
		{"TranslateY", s.TranslateY(7).Bounds().Center(), v3.XYZ(0, 7, 0)},
		{"TranslateZ", s.TranslateZ(-3).Bounds().Center(), v3.XYZ(0, 0, -3)},
		{"TranslateXY", s.TranslateXY(1, 2).Bounds().Center(), v3.XYZ(1, 2, 0)},
		{"TranslateXZ", s.TranslateXZ(1, 2).Bounds().Center(), v3.XYZ(1, 0, 2)},
		{"TranslateYZ", s.TranslateYZ(1, 2).Bounds().Center(), v3.XYZ(0, 1, 2)},
		{"TranslateXYZ", s.TranslateXYZ(1, 2, 3).Bounds().Center(), v3.XYZ(1, 2, 3)},
		{"Translate", s.Translate(v3.XYZ(4, 5, 6)).Bounds().Center(), v3.XYZ(4, 5, 6)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if !vecClose(c.got, c.want) {
				t.Fatalf("%s center = %+v, want %+v", c.name, c.got, c.want)
			}
		})
	}
}

// --- Rotate variants ---

func TestRotateVariants(t *testing.T) {
	t.Run("RotateX_swaps_Y_Z", func(t *testing.T) {
		// 90° about X swaps Y and Z extents.
		s := Box(v3.XYZ(2, 4, 6), 0)
		r := s.RotateX(90)
		size := r.Bounds().Size()
		// X stays 2; Y and Z swap.
		if math.Abs(size.X-2) > 1e-6 || math.Abs(size.Y-6) > 1e-6 || math.Abs(size.Z-4) > 1e-6 {
			t.Fatalf("RotateX(90) size = %+v, want X=2 Y=6 Z=4", size)
		}
	})
	t.Run("RotateY_swaps_X_Z", func(t *testing.T) {
		s := Box(v3.XYZ(2, 4, 6), 0)
		r := s.RotateY(90)
		size := r.Bounds().Size()
		if math.Abs(size.X-6) > 1e-6 || math.Abs(size.Y-4) > 1e-6 || math.Abs(size.Z-2) > 1e-6 {
			t.Fatalf("RotateY(90) size = %+v, want X=6 Y=4 Z=2", size)
		}
	})
	t.Run("RotateZ_swaps_X_Y", func(t *testing.T) {
		s := Box(v3.XYZ(2, 4, 6), 0)
		r := s.RotateZ(90)
		size := r.Bounds().Size()
		if math.Abs(size.X-4) > 1e-6 || math.Abs(size.Y-2) > 1e-6 || math.Abs(size.Z-6) > 1e-6 {
			t.Fatalf("RotateZ(90) size = %+v, want X=4 Y=2 Z=6", size)
		}
	})
	t.Run("RotateAxis_X_axis_equals_RotateX", func(t *testing.T) {
		s := Box(v3.XYZ(2, 4, 6), 0)
		a := s.RotateAxis(v3.XYZ(1, 0, 0), 90).Bounds().Size()
		b := s.RotateX(90).Bounds().Size()
		if !vecClose(a, b) {
			t.Fatalf("RotateAxis(X,90) = %+v, RotateX(90) = %+v", a, b)
		}
	})
	t.Run("RotateToVector_Z_to_X", func(t *testing.T) {
		// Rotating a Z-tall box from +Z to +X should make it X-long.
		s := Box(v3.XYZ(2, 2, 10), 0)
		r := s.RotateToVector(v3.XYZ(0, 0, 1), v3.XYZ(1, 0, 0))
		size := r.Bounds().Size()
		if math.Abs(size.X-10) > 1e-6 || math.Abs(size.Z-2) > 1e-6 {
			t.Fatalf("RotateToVector(Z,X) size = %+v, want X=10 Z=2", size)
		}
	})
}

// --- Mirror variants ---

func TestMirrorVariants(t *testing.T) {
	// Translate the box so the mirror has a measurable effect.
	s := Box(v3.XYZ(2, 2, 2), 0).Translate(v3.XYZ(3, 5, 7))
	t.Run("MirrorXY_negates_Z", func(t *testing.T) {
		c := s.MirrorXY().Bounds().Center()
		if !vecClose(c, v3.XYZ(3, 5, -7)) {
			t.Fatalf("MirrorXY center = %+v, want (3,5,-7)", c)
		}
	})
	t.Run("MirrorXZ_negates_Y", func(t *testing.T) {
		c := s.MirrorXZ().Bounds().Center()
		if !vecClose(c, v3.XYZ(3, -5, 7)) {
			t.Fatalf("MirrorXZ center = %+v, want (3,-5,7)", c)
		}
	})
	t.Run("MirrorYZ_negates_X", func(t *testing.T) {
		c := s.MirrorYZ().Bounds().Center()
		if !vecClose(c, v3.XYZ(-3, 5, 7)) {
			t.Fatalf("MirrorYZ center = %+v, want (-3,5,7)", c)
		}
	})
	t.Run("MirrorXeqY_swaps_X_Y", func(t *testing.T) {
		c := s.MirrorXeqY().Bounds().Center()
		if !vecClose(c, v3.XYZ(5, 3, 7)) {
			t.Fatalf("MirrorXeqY center = %+v, want (5,3,7)", c)
		}
	})
}

// --- Scale, ScaleUniform, Transform ---

func TestScaleScaleUniformTransform(t *testing.T) {
	s := Box(v3.XYZ(4, 4, 4), 0)
	t.Run("Scale", func(t *testing.T) {
		size := s.Scale(v3.XYZ(2, 3, 0.5)).Bounds().Size()
		if math.Abs(size.X-8) > 1e-6 || math.Abs(size.Y-12) > 1e-6 || math.Abs(size.Z-2) > 1e-6 {
			t.Fatalf("Scale(2,3,0.5) size = %+v, want (8,12,2)", size)
		}
	})
	t.Run("ScaleUniform", func(t *testing.T) {
		size := s.ScaleUniform(2).Bounds().Size()
		if math.Abs(size.X-8) > 1e-6 || math.Abs(size.Y-8) > 1e-6 || math.Abs(size.Z-8) > 1e-6 {
			t.Fatalf("ScaleUniform(2) size = %+v, want (8,8,8)", size)
		}
	})
	t.Run("Transform_Translate3d", func(t *testing.T) {
		m := Translate3d(v3.XYZ(1, 2, 3))
		c := s.Transform(m).Bounds().Center()
		if !vecClose(c, v3.XYZ(1, 2, 3)) {
			t.Fatalf("Transform(Translate3d) center = %+v, want (1,2,3)", c)
		}
	})
}

// --- Center / ZeroZ ---

func TestCenterAndZeroZ(t *testing.T) {
	t.Run("Center_moves_to_origin", func(t *testing.T) {
		s := Box(v3.XYZ(2, 2, 2), 0).Translate(v3.XYZ(7, 8, 9))
		c := s.Center().Bounds().Center()
		if !vecClose(c, v3.XYZ(0, 0, 0)) {
			t.Fatalf("Center() bounds center = %+v, want origin", c)
		}
	})
	t.Run("ZeroZ_puts_bottom_at_z0", func(t *testing.T) {
		s := Box(v3.XYZ(2, 2, 4), 0).Translate(v3.XYZ(0, 0, 7))
		z := s.ZeroZ().Bounds().Min.Z
		if math.Abs(z) > 1e-6 {
			t.Fatalf("ZeroZ bottom Z = %v, want 0", z)
		}
	})
}

// --- BottomAt/TopAt/LeftAt/RightAt/FrontAt/BackAt + CenterAt ---

func TestAxisAtSetters(t *testing.T) {
	s := Box(v3.XYZ(4, 6, 8), 0)
	t.Run("TopAt", func(t *testing.T) {
		got := s.TopAt(10).Bounds().Max.Z
		if math.Abs(got-10) > 1e-6 {
			t.Fatalf("TopAt(10) Max.Z = %v, want 10", got)
		}
	})
	t.Run("LeftAt", func(t *testing.T) {
		got := s.LeftAt(-2).Bounds().Min.X
		if math.Abs(got+2) > 1e-6 {
			t.Fatalf("LeftAt(-2) Min.X = %v, want -2", got)
		}
	})
	t.Run("RightAt", func(t *testing.T) {
		got := s.RightAt(5).Bounds().Max.X
		if math.Abs(got-5) > 1e-6 {
			t.Fatalf("RightAt(5) Max.X = %v, want 5", got)
		}
	})
	t.Run("FrontAt", func(t *testing.T) {
		got := s.FrontAt(-1).Bounds().Min.Y
		if math.Abs(got+1) > 1e-6 {
			t.Fatalf("FrontAt(-1) Min.Y = %v, want -1", got)
		}
	})
	t.Run("BackAt", func(t *testing.T) {
		got := s.BackAt(4).Bounds().Max.Y
		if math.Abs(got-4) > 1e-6 {
			t.Fatalf("BackAt(4) Max.Y = %v, want 4", got)
		}
	})
	t.Run("CenterAt", func(t *testing.T) {
		got := s.CenterAt(v3.XYZ(1, 2, 3)).Bounds().Center()
		if !vecClose(got, v3.XYZ(1, 2, 3)) {
			t.Fatalf("CenterAt center = %+v, want (1,2,3)", got)
		}
	})
}

// --- Anchor sugar verbs (UnderneathOf, LeftOf/RightOf/InFrontOf/BehindOf, Inside) ---

func TestAnchorSugarVerbs(t *testing.T) {
	body := Box(v3.XYZ(10, 10, 10), 0)
	tool := Box(v3.XYZ(2, 2, 2), 0)

	t.Run("UnderneathOf", func(t *testing.T) {
		moved := tool.UnderneathOf(body.Bottom(), 1).Solid()
		// tool's top should sit 1 below body's bottom (z=-5).
		want := -5.0 - 1.0
		if math.Abs(moved.Top().Point.Z-want) > 1e-6 {
			t.Fatalf("UnderneathOf top.Z = %v, want %v", moved.Top().Point.Z, want)
		}
	})
	t.Run("LeftOf_sugar", func(t *testing.T) {
		moved := tool.LeftOf(body.Left(), 0.5).Solid()
		// tool's right at body's left - 0.5 = -5.5.
		if math.Abs(moved.Right().Point.X+5.5) > 1e-6 {
			t.Fatalf("LeftOf right.X = %v, want -5.5", moved.Right().Point.X)
		}
	})
	t.Run("RightOf_sugar", func(t *testing.T) {
		moved := tool.RightOf(body.Right(), 0.5).Solid()
		if math.Abs(moved.Left().Point.X-5.5) > 1e-6 {
			t.Fatalf("RightOf left.X = %v, want 5.5", moved.Left().Point.X)
		}
	})
	t.Run("InFrontOf_sugar", func(t *testing.T) {
		moved := tool.InFrontOf(body.Front(), 0.5).Solid()
		if math.Abs(moved.Back().Point.Y+5.5) > 1e-6 {
			t.Fatalf("InFrontOf back.Y = %v, want -5.5", moved.Back().Point.Y)
		}
	})
	t.Run("BehindOf_sugar", func(t *testing.T) {
		moved := tool.BehindOf(body.Back(), 0.5).Solid()
		if math.Abs(moved.Front().Point.Y-5.5) > 1e-6 {
			t.Fatalf("BehindOf front.Y = %v, want 5.5", moved.Front().Point.Y)
		}
	})
	t.Run("Inside", func(t *testing.T) {
		moved := tool.Translate(v3.XYZ(50, 50, 50)).Inside(body).Solid()
		c := moved.Bounds().Center()
		if !vecClose(c, v3.XYZ(0, 0, 0)) {
			t.Fatalf("Inside center = %+v, want origin", c)
		}
	})
}

// --- AnchoredSolid placement methods ---

func TestAnchoredSolidPlacement(t *testing.T) {
	a := Box(v3.XYZ(2, 2, 2), 0)
	b := Box(v3.XYZ(10, 10, 10), 0)

	t.Run("Above_Below_RightOf_LeftOf_Behind_InFrontOf", func(t *testing.T) {
		_ = a.Bottom().Above(b.Top(), 0)
		_ = a.Top().Below(b.Bottom(), 0)
		_ = a.Left().RightOf(b.Right(), 0)
		_ = a.Right().LeftOf(b.Left(), 0)
		_ = a.Front().Behind(b.Back(), 0)
		_ = a.Back().InFrontOf(b.Front(), 0)
	})
	t.Run("AtX_AtY_At", func(t *testing.T) {
		got := a.Right().AtX(10).Bounds().Max.X
		if math.Abs(got-10) > 1e-6 {
			t.Fatalf("AtX(10) Max.X = %v, want 10", got)
		}
		got2 := a.Front().AtY(-3).Bounds().Min.Y
		if math.Abs(got2+3) > 1e-6 {
			t.Fatalf("AtY(-3) Min.Y = %v, want -3", got2)
		}
		got3 := a.AnchorAt(0, 0, 0).At(v3.XYZ(1, 2, 3)).Bounds().Center()
		if !vecClose(got3, v3.XYZ(1, 2, 3)) {
			t.Fatalf("At(1,2,3) center = %+v, want (1,2,3)", got3)
		}
	})
	t.Run("ShiftX_ShiftY_ShiftZ", func(t *testing.T) {
		base := a.Top()
		shifted := base.ShiftX(1).ShiftY(2).ShiftZ(3)
		want := base.Point.Add(v3.XYZ(1, 2, 3))
		if !vecClose(shifted.Point, want) {
			t.Fatalf("ShiftXYZ point = %+v, want %+v", shifted.Point, want)
		}
	})
}

// --- Placement boolean variants ---

func TestPlacementBooleans(t *testing.T) {
	body := Box(v3.XYZ(10, 10, 10), 0)
	tool := Box(v3.XYZ(2, 2, 2), 0)
	min := PolyMin(0.5)
	max := PolyMax(0.5)

	t.Run("Union_Add", func(t *testing.T) {
		p := tool.Bottom().Above(body.Top())
		_ = p.Union()
		_ = p.Add()
	})
	t.Run("SmoothUnion_SmoothAdd", func(t *testing.T) {
		p := tool.Bottom().Above(body.Top())
		_ = p.SmoothUnion(min)
		_ = p.SmoothAdd(min)
	})
	t.Run("Difference_Cut", func(t *testing.T) {
		p := tool.Bottom().Above(body.Top())
		_ = p.Difference()
		_ = p.Cut()
	})
	t.Run("SmoothDifference_SmoothCut", func(t *testing.T) {
		p := tool.Bottom().Above(body.Top())
		_ = p.SmoothDifference(max)
		_ = p.SmoothCut(max)
	})
	t.Run("Intersect_SmoothIntersect", func(t *testing.T) {
		p := tool.Inside(body)
		_ = p.Intersect()
		_ = p.SmoothIntersect(max)
	})
}

// --- Smooth booleans (free functions) ---

func TestSmoothFreeFunctions(t *testing.T) {
	a := Box(v3.XYZ(4, 4, 4), 0)
	b := Sphere(3).TranslateX(3)
	min := PolyMin(0.5)
	max := PolyMax(0.5)

	t.Run("SmoothUnion", func(t *testing.T) {
		got := SmoothUnion(min, a, b)
		if got == nil {
			t.Fatal("SmoothUnion returned nil")
		}
	})
	t.Run("SmoothAdd_alias", func(t *testing.T) {
		got := SmoothAdd(min, a, b)
		if got == nil {
			t.Fatal("SmoothAdd returned nil")
		}
	})
	t.Run("SmoothDifference", func(t *testing.T) {
		got := SmoothDifference(max, a, b)
		if got == nil {
			t.Fatal("SmoothDifference returned nil")
		}
	})
	t.Run("SmoothCut_alias", func(t *testing.T) {
		got := SmoothCut(max, a, b)
		if got == nil {
			t.Fatal("SmoothCut returned nil")
		}
	})
	t.Run("SmoothIntersection", func(t *testing.T) {
		got := SmoothIntersection(max, a, b)
		if got == nil {
			t.Fatal("SmoothIntersection returned nil")
		}
	})
}

// --- Smooth booleans (method counterparts) ---

func TestSmoothMethods(t *testing.T) {
	a := Box(v3.XYZ(4, 4, 4), 0)
	b := Sphere(3).TranslateX(3)
	min := PolyMin(0.5)
	max := PolyMax(0.5)

	if a.SmoothUnion(min, b) == nil {
		t.Fatal("SmoothUnion method nil")
	}
	if a.SmoothAdd(min, b) == nil {
		t.Fatal("SmoothAdd method nil")
	}
	if a.SmoothCut(max, b) == nil {
		t.Fatal("SmoothCut method nil")
	}
	if a.SmoothDifference(max, b) == nil {
		t.Fatal("SmoothDifference method nil")
	}
	if a.SmoothIntersect(max, b) == nil {
		t.Fatal("SmoothIntersect method nil")
	}
}

// --- Min/Max blend functions ---

func TestBlendFunctions(t *testing.T) {
	if PolyMin(1.0) == nil {
		t.Fatal("PolyMin nil")
	}
	if PolyMax(1.0) == nil {
		t.Fatal("PolyMax nil")
	}
	if RoundMin(1.0) == nil {
		t.Fatal("RoundMin nil")
	}
	if ChamferMin(1.0) == nil {
		t.Fatal("ChamferMin nil")
	}
	if ExpMin(1.0) == nil {
		t.Fatal("ExpMin nil")
	}
	if PowMin(1.0) == nil {
		t.Fatal("PowMin nil")
	}
}

// --- Shell, Offset, Grow, Shrink, Elongate ---

func TestShellOffsetGrowShrinkElongate(t *testing.T) {
	s := Box(v3.XYZ(10, 10, 10), 0)

	t.Run("Shell", func(t *testing.T) {
		if s.Shell(0.5) == nil {
			t.Fatal("Shell nil")
		}
	})
	t.Run("Offset", func(t *testing.T) {
		if s.Offset(0.5) == nil {
			t.Fatal("Offset nil")
		}
	})
	t.Run("Grow_expands_bbox", func(t *testing.T) {
		got := s.Grow(2).Bounds().Size()
		// offsetSDF3 expands by amount on each side: original size + 2*amount.
		want := v3.XYZ(14, 14, 14)
		if !vecClose(got, want) {
			t.Fatalf("Grow(2) size = %+v, want %+v", got, want)
		}
	})
	t.Run("Shrink_contracts_bbox", func(t *testing.T) {
		got := s.Shrink(2).Bounds().Size()
		want := v3.XYZ(6, 6, 6)
		if !vecClose(got, want) {
			t.Fatalf("Shrink(2) size = %+v, want %+v", got, want)
		}
	})
	t.Run("Elongate", func(t *testing.T) {
		// Elongate stretches by h on each axis.
		got := s.Elongate(v3.XYZ(2, 0, 0)).Bounds().Size()
		// Should be larger in X.
		if got.X <= 10 {
			t.Fatalf("Elongate X size = %v, expected > 10", got.X)
		}
	})
}

// --- Extrude family ---

func TestExtrudeFamily(t *testing.T) {
	prof := rect2D(4, 4)
	t.Run("Extrude", func(t *testing.T) {
		s := Extrude(prof, 6)
		if math.Abs(s.Bounds().Size().Z-6) > 1e-6 {
			t.Fatalf("Extrude Z size = %v, want 6", s.Bounds().Size().Z)
		}
	})
	t.Run("ExtrudeRounded", func(t *testing.T) {
		s := ExtrudeRounded(prof, 6, 0.5)
		if s == nil {
			t.Fatal("ExtrudeRounded nil")
		}
	})
	t.Run("TwistExtrude", func(t *testing.T) {
		s := TwistExtrude(prof, 6, 90)
		if math.Abs(s.Bounds().Size().Z-6) > 1e-6 {
			t.Fatalf("TwistExtrude Z size = %v, want 6", s.Bounds().Size().Z)
		}
	})
	t.Run("ScaleExtrude", func(t *testing.T) {
		s := ScaleExtrude(prof, 6, v2.XY(0.5, 0.5))
		if s == nil {
			t.Fatal("ScaleExtrude nil")
		}
	})
	t.Run("ScaleTwistExtrude", func(t *testing.T) {
		s := ScaleTwistExtrude(prof, 6, 45, v2.XY(0.8, 0.8))
		if s == nil {
			t.Fatal("ScaleTwistExtrude nil")
		}
	})
	t.Run("Loft", func(t *testing.T) {
		bottom := rect2D(4, 4)
		top := rect2D(2, 2)
		s := Loft(bottom, top, 6, 0)
		if s == nil {
			t.Fatal("Loft nil")
		}
	})
}

// --- Revolve / Screw / SweepHelix ---

func TestRevolveScrewSweepHelix(t *testing.T) {
	// Revolve a profile offset from origin so it forms a torus-like solid.
	prof := sdf.Transform2D(rect2D(2, 2), sdf.Translate2d(v2sdf.Vec{X: 0, Y: 5}))

	t.Run("Revolve", func(t *testing.T) {
		s := Revolve(prof)
		if s == nil {
			t.Fatal("Revolve nil")
		}
	})
	t.Run("RevolveAngle", func(t *testing.T) {
		s := RevolveAngle(prof, 180)
		if s == nil {
			t.Fatal("RevolveAngle nil")
		}
	})
	t.Run("Screw", func(t *testing.T) {
		s := Screw(rect2D(1, 1), 8, 0, 2, 1)
		if s == nil {
			t.Fatal("Screw nil")
		}
	})
	t.Run("SweepHelix_round_ends", func(t *testing.T) {
		s := SweepHelix(rect2D(1, 1), 5, 2, 8, false)
		if s == nil {
			t.Fatal("SweepHelix nil")
		}
	})
	t.Run("SweepHelix_flat_ends", func(t *testing.T) {
		s := SweepHelix(rect2D(1, 1), 5, 2, 8, true)
		if s == nil {
			t.Fatal("SweepHelix flatEnds nil")
		}
		// Trigger Evaluate via Bounds + a sample.
		_ = s.SDF3.Evaluate(v3sdf.Vec{X: 5, Y: 0, Z: 0})
		_ = s.BoundingBox()
	})
}

// --- Multi, LineOf, Array, RotateCopyZ, RotateUnion variants, Orient ---

func TestPatternsAndArrays(t *testing.T) {
	s := Box(v3.XYZ(2, 2, 2), 0)

	t.Run("Multi_with_positions", func(t *testing.T) {
		got := s.Multi(v3.XYZ(0, 0, 0), v3.XYZ(10, 0, 0)).Bounds()
		if got.Size().X < 12 {
			t.Fatalf("Multi X size = %v, expected >= 12", got.Size().X)
		}
	})
	t.Run("LineOf_with_xx_pattern", func(t *testing.T) {
		got := s.LineOf(v3.XYZ(0, 0, 0), v3.XYZ(10, 0, 0), "xx")
		if got == nil {
			t.Fatal("LineOf nil")
		}
	})
	t.Run("Array", func(t *testing.T) {
		got := s.Array(2, 1, 1, v3.XYZ(5, 0, 0))
		if got == nil {
			t.Fatal("Array nil")
		}
	})
	t.Run("SmoothArray", func(t *testing.T) {
		got := s.SmoothArray(2, 1, 1, v3.XYZ(5, 0, 0), PolyMin(0.5))
		if got == nil {
			t.Fatal("SmoothArray nil")
		}
	})
	t.Run("RotateCopyZ", func(t *testing.T) {
		got := s.TranslateX(5).RotateCopyZ(4)
		if got == nil {
			t.Fatal("RotateCopyZ nil")
		}
	})
	t.Run("RotateUnionZ", func(t *testing.T) {
		got := s.TranslateX(5).RotateUnionZ(4, RotateZMatrix(90))
		if got == nil {
			t.Fatal("RotateUnionZ nil")
		}
	})
	t.Run("SmoothRotateUnionZ", func(t *testing.T) {
		got := s.TranslateX(5).SmoothRotateUnionZ(4, RotateZMatrix(90), PolyMin(0.5))
		if got == nil {
			t.Fatal("SmoothRotateUnionZ nil")
		}
	})
	t.Run("Orient", func(t *testing.T) {
		got := s.Orient(v3.XYZ(0, 0, 1), []v3.Vec{v3.XYZ(1, 0, 0), v3.XYZ(0, 1, 0)})
		if got == nil {
			t.Fatal("Orient nil")
		}
	})
}

// --- Cross-section: Slice, Slice2D, SliceAt, Bounds, BoundingBox ---

func TestCrossSectionAndBounds(t *testing.T) {
	s := Sphere(5)

	t.Run("Slice_free_function", func(t *testing.T) {
		got := Slice(s, v3.XYZ(0, 0, 0), v3.XYZ(0, 0, 1))
		if got == nil {
			t.Fatal("Slice nil")
		}
	})
	t.Run("Slice2D_method", func(t *testing.T) {
		got := s.Slice2D(v3.XYZ(0, 0, 0), v3.XYZ(0, 0, 1))
		if got == nil {
			t.Fatal("Slice2D nil")
		}
	})
	t.Run("SliceAt_method", func(t *testing.T) {
		got := s.SliceAt(plane.AtZ(0))
		if got == nil {
			t.Fatal("SliceAt nil")
		}
	})
	t.Run("Bounds_and_BoundingBox_agree", func(t *testing.T) {
		bb := s.BoundingBox()
		bounds := s.Bounds()
		if math.Abs(bb.Min.X-bounds.Min.X) > 1e-9 || math.Abs(bb.Max.Z-bounds.Max.Z) > 1e-9 {
			t.Fatalf("Bounds %+v ≠ BoundingBox %+v", bounds, bb)
		}
	})
}

// --- Wrap, New, UnionAll multi-arg ---

func TestWrapNewUnionAll(t *testing.T) {
	raw, err := sdf.Sphere3D(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("Wrap", func(t *testing.T) {
		s := Wrap(raw)
		if s == nil {
			t.Fatal("Wrap nil")
		}
	})
	t.Run("Wrap_nil_panics", func(t *testing.T) {
		expectPanic(t, "nil sdf.SDF3", func() { Wrap(nil) })
	})
	t.Run("New_with_error_panics", func(t *testing.T) {
		expectPanic(t, "", func() { _ = New(nil, errString("boom")) })
	})
	t.Run("UnionAll_two_args", func(t *testing.T) {
		got := UnionAll(Box(v3.XYZ(2, 2, 2), 0), Box(v3.XYZ(2, 2, 2), 0).TranslateX(10))
		if got == nil {
			t.Fatal("UnionAll nil")
		}
	})
}

type errString string

func (e errString) Error() string { return string(e) }

// --- Constructors not previously tested: Capsule, Torus, Gyroid, Cone variants ---

func TestUntestedConstructors(t *testing.T) {
	t.Run("Capsule", func(t *testing.T) {
		c := Capsule(10, 2)
		if c == nil {
			t.Fatal("Capsule nil")
		}
	})
	t.Run("Torus", func(t *testing.T) {
		tor := Torus(10, 2)
		bb := tor.Bounds()
		// Expect Z extent ±2, XY extent ±12.
		if math.Abs(bb.Max.Z-2) > 1e-9 || math.Abs(bb.Min.Z+2) > 1e-9 {
			t.Fatalf("Torus Z bounds = [%v, %v], want [-2, 2]", bb.Min.Z, bb.Max.Z)
		}
		if math.Abs(bb.Max.X-12) > 1e-9 || math.Abs(bb.Min.X+12) > 1e-9 {
			t.Fatalf("Torus X bounds = [%v, %v], want [-12, 12]", bb.Min.X, bb.Max.X)
		}
		// Evaluate at center should be negative (inside is far from tube): actually
		// the center of a torus is OUTSIDE the tube (distance majorR=10 minus minorR=2).
		d := tor.SDF3.Evaluate(v3sdf.Vec{X: 0, Y: 0, Z: 0})
		if d < 0 {
			t.Fatalf("Torus center should be outside (d>0), got %v", d)
		}
		// Evaluate at a point on the tube center should be inside.
		d2 := tor.SDF3.Evaluate(v3sdf.Vec{X: 10, Y: 0, Z: 0})
		if d2 > 0 {
			t.Fatalf("Torus tube-center should be inside (d<0), got %v", d2)
		}
	})
	t.Run("Gyroid", func(t *testing.T) {
		g := Gyroid(v3.XYZ(1, 1, 1))
		if g == nil {
			t.Fatal("Gyroid nil")
		}
	})
	t.Run("Cone_with_round", func(t *testing.T) {
		c := Cone(10, 4, 2, 0.5)
		if c == nil {
			t.Fatal("Cone(round) nil")
		}
	})
	t.Run("Cone_to_tip", func(t *testing.T) {
		c := Cone(10, 4, 0, 0)
		if c == nil {
			t.Fatal("Cone tip nil")
		}
	})
}

// --- Other Solid ops: Cut/Difference, CutPlane, Split, Correct ---

func TestCutPlaneSplitCorrect(t *testing.T) {
	s := Box(v3.XYZ(10, 10, 10), 0)

	t.Run("CutPlane", func(t *testing.T) {
		half := s.CutPlane(v3.XYZ(0, 0, 0), v3.XYZ(0, 0, 1))
		// Should keep the +Z half: bottom should be near 0.
		if half.Bounds().Min.Z < -1 {
			t.Fatalf("CutPlane Min.Z = %v, expected near 0", half.Bounds().Min.Z)
		}
	})
	t.Run("Split", func(t *testing.T) {
		a, b := s.Split(plane.AtZ(0))
		if a == nil || b == nil {
			t.Fatal("Split returned nil halves")
		}
	})
	t.Run("Correct", func(t *testing.T) {
		got := s.Correct(0.5)
		if got == nil {
			t.Fatal("Correct nil")
		}
	})
	t.Run("Difference_alias", func(t *testing.T) {
		other := Box(v3.XYZ(2, 2, 2), 0)
		got := s.Difference(other)
		if got == nil {
			t.Fatal("Difference nil")
		}
	})
}

// --- offsetSDF3 / correctedSDF3 / torusSDF3 Evaluate + BoundingBox direct ---

func TestInternalSDF3Evaluate(t *testing.T) {
	s := Box(v3.XYZ(10, 10, 10), 0)

	t.Run("offsetSDF3_Evaluate", func(t *testing.T) {
		o := s.Grow(1)
		_ = o.SDF3.Evaluate(v3sdf.Vec{X: 0, Y: 0, Z: 0})
		_ = o.SDF3.BoundingBox()
	})
	t.Run("correctedSDF3_Evaluate", func(t *testing.T) {
		c := s.Correct(0.5)
		_ = c.SDF3.Evaluate(v3sdf.Vec{X: 0, Y: 0, Z: 0})
		_ = c.SDF3.BoundingBox()
	})
	t.Run("torusSDF3_BoundingBox", func(t *testing.T) {
		tor := Torus(5, 1)
		bb := tor.SDF3.BoundingBox()
		if bb.Max.X != 6 {
			t.Fatalf("Torus bb.Max.X = %v, want 6", bb.Max.X)
		}
	})
}

// --- Render path: CellsFor + SetMinCells ---

func TestRenderHelpers(t *testing.T) {
	s := Box(v3.XYZ(20, 10, 5), 0)
	t.Run("CellsFor_default", func(t *testing.T) {
		got := CellsFor(s, 5.0)
		// Longest axis 20 * 5 = 100.
		if got != 100 {
			t.Fatalf("CellsFor = %v, want 100", got)
		}
	})
	t.Run("CellsFor_floored", func(t *testing.T) {
		// Tiny part: 0.1 mm box at 1 cell/mm = 1 cell, floored to MinCells=32.
		tiny := Box(v3.XYZ(0.1, 0.1, 0.1), 0)
		got := CellsFor(tiny, 1.0)
		if got != MinCells {
			t.Fatalf("tiny CellsFor = %v, want MinCells=%v", got, MinCells)
		}
	})
	t.Run("SetMinCells_round_trip", func(t *testing.T) {
		old := MinCells
		defer SetMinCells(old)
		SetMinCells(64)
		// After setting, minCells() should report 64.
		if minCells() != 64 {
			t.Fatalf("minCells() = %v, want 64", minCells())
		}
		// And tiny CellsFor should now floor at 64.
		tiny := Box(v3.XYZ(0.1, 0.1, 0.1), 0)
		got := CellsFor(tiny, 1.0)
		if got != 64 {
			t.Fatalf("tiny CellsFor after SetMinCells(64) = %v, want 64", got)
		}
	})
}

// --- Sampling/utility methods: Normal, Raycast, Voxel, Benchmark ---

func TestSamplingMethods(t *testing.T) {
	s := Sphere(5)

	t.Run("Normal_at_surface_points_outward", func(t *testing.T) {
		n := s.Normal(v3.XYZ(5, 0, 0), 1e-3)
		// Should be roughly (1,0,0).
		if n.X < 0.9 {
			t.Fatalf("Normal at (5,0,0) X = %v, want ~1", n.X)
		}
	})
	t.Run("Raycast_hits_surface", func(t *testing.T) {
		_, dist, _ := s.Raycast(v3.XYZ(20, 0, 0), v3.XYZ(-1, 0, 0), 0, 0.5, 1e-3, 100, 100)
		// Should hit ~15 mm in.
		if dist <= 0 || dist > 20 {
			t.Fatalf("Raycast distance = %v, want positive ≤ 20", dist)
		}
	})
	t.Run("Voxel", func(t *testing.T) {
		v := s.Voxel(8, nil)
		if v == nil {
			t.Fatal("Voxel nil")
		}
	})
}

// --- Mesh / MeshSlow ---

func TestMeshConstructors(t *testing.T) {
	// Single triangle isn't a closed solid, but Mesh3D should still build the SDF.
	tris := []*v3.Triangle3{
		{v3.XYZ(0, 0, 0), v3.XYZ(1, 0, 0), v3.XYZ(0, 1, 0)},
		{v3.XYZ(0, 0, 0), v3.XYZ(0, 1, 0), v3.XYZ(0, 0, 1)},
		{v3.XYZ(0, 0, 0), v3.XYZ(0, 0, 1), v3.XYZ(1, 0, 0)},
		{v3.XYZ(1, 0, 0), v3.XYZ(0, 1, 0), v3.XYZ(0, 0, 1)},
	}
	t.Run("Mesh", func(t *testing.T) {
		s := Mesh(tris)
		if s == nil {
			t.Fatal("Mesh nil")
		}
	})
	t.Run("MeshSlow", func(t *testing.T) {
		s := MeshSlow(tris)
		if s == nil {
			t.Fatal("MeshSlow nil")
		}
	})
}

// --- Matrix package surface ---

func TestMatrixSurface(t *testing.T) {
	t.Run("Constructors", func(t *testing.T) {
		_ = NewM44([16]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1})
		_ = Translate3d(v3.XYZ(1, 2, 3))
		_ = RotateXMatrix(45)
		_ = RotateYMatrix(45)
		_ = RotateZMatrix(45)
		_ = Rotate3dMatrix(v3.XYZ(0, 0, 1), 30)
		_ = RotateToVector(v3.XYZ(1, 0, 0), v3.XYZ(0, 1, 0))
		_ = Scale3d(v3.XYZ(2, 2, 2))
		_ = MirrorXY()
		_ = MirrorXZ()
		_ = MirrorYZ()
		_ = MirrorXeqY()
	})
	t.Run("Identity_and_inverse", func(t *testing.T) {
		id := Identity3d()
		if math.Abs(id.Determinant()-1) > 1e-9 {
			t.Fatalf("Identity determinant = %v, want 1", id.Determinant())
		}
		// Inverse of identity is identity.
		inv := id.Inverse()
		if !id.Equals(inv, 1e-9) {
			t.Fatalf("Inverse of identity not identity")
		}
	})
	t.Run("Mul_and_MulPosition", func(t *testing.T) {
		m := Translate3d(v3.XYZ(1, 2, 3)).Mul(Identity3d())
		got := m.MulPosition(v3.XYZ(0, 0, 0))
		if !vecClose(got, v3.XYZ(1, 2, 3)) {
			t.Fatalf("MulPosition = %+v, want (1,2,3)", got)
		}
	})
	t.Run("Values_round_trip", func(t *testing.T) {
		v := Identity3d().Values()
		if v[0] != 1 || v[5] != 1 || v[10] != 1 || v[15] != 1 {
			t.Fatalf("Identity Values diagonal = (%v,%v,%v,%v), want all 1", v[0], v[5], v[10], v[15])
		}
	})
	t.Run("MulBox", func(t *testing.T) {
		b := NewBox3(v3.XYZ(0, 0, 0), v3.XYZ(2, 2, 2))
		out := Translate3d(v3.XYZ(5, 0, 0)).MulBox(b)
		if !vecClose(out.Center(), v3.XYZ(5, 0, 0)) {
			t.Fatalf("MulBox center = %+v, want (5,0,0)", out.Center())
		}
	})
}

// --- Box3 chainable methods ---

func TestBox3Chainables(t *testing.T) {
	b := NewBox3(v3.XYZ(0, 0, 0), v3.XYZ(2, 4, 6))
	t.Run("Solid", func(t *testing.T) {
		s := b.Solid()
		if s == nil {
			t.Fatal("Box3.Solid nil")
		}
		if !vecClose(s.Bounds().Size(), v3.XYZ(2, 4, 6)) {
			t.Fatalf("Box3.Solid size = %+v, want (2,4,6)", s.Bounds().Size())
		}
	})
	t.Run("Cube", func(t *testing.T) {
		c := b.Cube()
		// Should be cube of largest dimension.
		s := c.Size()
		if s.X != s.Y || s.Y != s.Z {
			t.Fatalf("Cube size = %+v, want isotropic", s)
		}
	})
	t.Run("Enlarge", func(t *testing.T) {
		e := b.Enlarge(v3.XYZ(1, 1, 1))
		if e.Size().X <= b.Size().X {
			t.Fatalf("Enlarge did not grow")
		}
	})
	t.Run("Extend", func(t *testing.T) {
		o := NewBox3(v3.XYZ(20, 0, 0), v3.XYZ(2, 2, 2))
		e := b.Extend(o)
		if e.Size().X <= b.Size().X {
			t.Fatalf("Extend did not grow X")
		}
	})
	t.Run("Include", func(t *testing.T) {
		i := b.Include(v3.XYZ(20, 0, 0))
		if i.Max.X < 20 {
			t.Fatalf("Include did not extend Max.X to 20")
		}
	})
	t.Run("ScaleAboutCenter", func(t *testing.T) {
		sc := b.ScaleAboutCenter(2)
		if !vecClose(sc.Size(), v3.XYZ(4, 8, 12)) {
			t.Fatalf("ScaleAboutCenter(2) size = %+v, want (4,8,12)", sc.Size())
		}
	})
	t.Run("Translate", func(t *testing.T) {
		tr := b.Translate(v3.XYZ(10, 0, 0))
		if !vecClose(tr.Center(), v3.XYZ(10, 0, 0)) {
			t.Fatalf("Translate center = %+v, want (10,0,0)", tr.Center())
		}
	})
	t.Run("Equals", func(t *testing.T) {
		if !b.Equals(b, 1e-9) {
			t.Fatal("Box3 self.Equals = false")
		}
	})
}

// --- Helpers: v3Slice ---

func TestV3Slice(t *testing.T) {
	out := v3Slice([]v3.Vec{v3.XYZ(1, 2, 3), v3.XYZ(4, 5, 6)})
	if len(out) != 2 {
		t.Fatalf("v3Slice len = %v, want 2", len(out))
	}
	if out[0].X != 1 || out[1].Z != 6 {
		t.Fatalf("v3Slice values wrong: %+v", out)
	}
}
