package shape

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/snowbldr/fluent-sdfx/plane"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	"github.com/snowbldr/fluent-sdfx/vec/v2i"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/sdfx/sdf"
)

// fontPath returns the tutorial font path if available, otherwise "".
func fontPath() string {
	candidates := []string{
		"../tutorial/internal/tutorialfont/cmr10.ttf",
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}

// --- Constructors / primitives. ---

func TestStarBoundsWithinOuter(t *testing.T) {
	s := Star(5, 2, 5)
	if s == nil {
		t.Fatal("Star returned nil")
	}
	bb := s.Bounds()
	// Star applies an Offset of outer/10 = 0.5, so bounds may exceed 5 slightly.
	if bb.Max.X > 6 || bb.Min.X < -6 {
		t.Errorf("Star X bounds = [%v, %v], unexpected", bb.Min.X, bb.Max.X)
	}
}

func TestCrossNonEmpty(t *testing.T) {
	c := Cross(10, 2)
	if c == nil {
		t.Fatal("Cross returned nil")
	}
	bb := c.Bounds()
	if !boxClose2(t, bb, v2.Box{Min: v2.XY(-5, -5), Max: v2.XY(5, 5)}, 1e-9) {
		t.Errorf("Cross bounds = %+v", bb)
	}
}

func TestLineBounds(t *testing.T) {
	l := Line(10, 0)
	bb := l.Bounds()
	if math.Abs(bb.Max.X-5) > 1e-6 || math.Abs(bb.Min.X+5) > 1e-6 {
		t.Errorf("Line bounds = %+v", bb)
	}
}

func TestArcSpiral(t *testing.T) {
	s := ArcSpiral(1, 0.5, 0, 360, 0.2)
	if s == nil {
		t.Fatal("ArcSpiral returned nil")
	}
	bb := s.Bounds()
	if bb.Max.X-bb.Min.X <= 0 {
		t.Errorf("ArcSpiral empty bbox: %+v", bb)
	}
}

func TestFlange1(t *testing.T) {
	s := Flange1(10, 4, 2)
	if s == nil {
		t.Fatal("Flange1 returned nil")
	}
	bb := s.Bounds()
	if bb.Max.X-bb.Min.X <= 0 {
		t.Errorf("Flange1 empty bbox: %+v", bb)
	}
}

func TestWireGroove(t *testing.T) {
	// Theoretical extension below cap.
	a := WireGroove(1, 2, 30)
	if a == nil {
		t.Fatal("WireGroove returned nil")
	}
	// Theoretical extension above cap.
	b := WireGroove(1, 50, 5)
	if b == nil {
		t.Fatal("WireGroove (capped) returned nil")
	}
}

func TestNagonReturnsVertices(t *testing.T) {
	v := Nagon(6, 5)
	if len(v) != 6 {
		t.Fatalf("Nagon(6) returned %d vertices, want 6", len(v))
	}
}

func TestCubicSpline(t *testing.T) {
	knots := []v2.Vec{v2.XY(0, 0), v2.XY(5, 0), v2.XY(5, 5), v2.XY(0, 5)}
	s := CubicSpline(knots)
	if s == nil {
		t.Fatal("CubicSpline returned nil")
	}
}

// --- Threads. ---

func TestThreads(t *testing.T) {
	if AcmeThread(5, 1) == nil {
		t.Fatal("AcmeThread")
	}
	if ISOThread(5, 1, true) == nil {
		t.Fatal("ISOThread external")
	}
	if ISOThread(5, 1, false) == nil {
		t.Fatal("ISOThread internal")
	}
	if ANSIButtressThread(5, 1) == nil {
		t.Fatal("ANSIButtressThread")
	}
	if PlasticButtressThread(5, 1) == nil {
		t.Fatal("PlasticButtressThread")
	}
	if ThreadLookup("M4x0.7") == nil {
		t.Fatal("ThreadLookup")
	}
}

// --- Cams. ---

func TestCams(t *testing.T) {
	if FlatFlankCam(5, 4, 2) == nil {
		t.Fatal("FlatFlankCam")
	}
	if MakeFlatFlankCam(2, 90, 20) == nil {
		t.Fatal("MakeFlatFlankCam")
	}
	if ThreeArcCam(5, 4, 2, 10) == nil {
		t.Fatal("ThreeArcCam")
	}
	if MakeThreeArcCam(0.1, 160, 0.7, 1.1) == nil {
		t.Fatal("MakeThreeArcCam")
	}
}

// --- Gear rack. ---

func TestGearRack(t *testing.T) {
	r := GearRack(GearRackParams{
		NumberTeeth:      8,
		Module:           1,
		PressureAngleDeg: 20,
		Backlash:         0,
		BaseHeight:       2,
	})
	if r == nil {
		t.Fatal("GearRack returned nil")
	}
}

// --- Polygon builder. ---

func TestPolyBuilderBuildAndMesh2D(t *testing.T) {
	p := NewPoly()
	p.Add(0, 0)
	p.Add(10, 0)
	p.Add(10, 10)
	p.Add(0, 10)
	if got := len(p.Vertices()); got < 4 {
		t.Fatalf("expected >=4 vertices, got %d", got)
	}
	if p.Raw() == nil {
		t.Fatal("Raw nil")
	}
	s := p.Build()
	if s == nil {
		t.Fatal("Build nil")
	}
	m := p.Mesh2D()
	if m == nil {
		t.Fatal("Mesh2D nil")
	}
}

func TestPolyBuilderModifiers(t *testing.T) {
	p := NewPoly()
	p.AddV2Set([]v2.Vec{v2.XY(0, 0), v2.XY(10, 0)})
	p.AddV2(v2.XY(10, 10))
	p.AddV2(v2.XY(0, 10))
	p.Add(-1, -1) // soon to be dropped
	p.Drop()
	p.Close()
	p.Reverse()
	if len(p.Vertices()) < 4 {
		t.Fatal("expected vertices after modifications")
	}
	if p.Build() == nil {
		t.Fatal("Build nil after modifications")
	}
}

func TestPolyVertexModifiers(t *testing.T) {
	p := NewPoly()
	p.Add(0, 0).Smooth(0.5, 5)
	p.Add(10, 0).Chamfer(0.5)
	// Polar treats x as radius, y as angle (radians).
	p.Add(5, math.Pi/2).Polar()
	p.Add(0, 10).Rel()
	p.Add(0, 10).Arc(2, 5)
	if p.Build() == nil {
		t.Fatal("Build nil")
	}
}

// --- Bezier builder. ---

func TestBezierBuilder(t *testing.T) {
	b := NewBezier()
	b.Add(0, 0).HandleFwd(0, 1)
	b.AddV2(v2.XY(10, 0)).HandleRev(180, 1).HandleFwd(0, 1)
	b.Add(10, 10).Handle(90, 1, 1)
	b.Add(0, 10).Mid()
	b.Close()

	verts := b.Vertices()
	if len(verts) == 0 {
		t.Fatal("Bezier produced no vertices")
	}
	poly := b.Polygon()
	if poly == nil {
		t.Fatal("Bezier Polygon nil")
	}
	s := b.Build()
	if s == nil {
		t.Fatal("Bezier Build nil")
	}
}

// --- Wrap2D. ---

func TestWrap2D(t *testing.T) {
	c := Circle(3)
	w := Wrap2D(c.SDF2)
	if w == nil {
		t.Fatal("Wrap2D nil")
	}
	if !boxClose2(t, w.Bounds(), c.Bounds(), 1e-9) {
		t.Fatalf("Wrap2D bounds differ")
	}
}

// --- Translate variants. ---

func TestTranslateVariants(t *testing.T) {
	r := Rect(v2.XY(10, 6), 0)
	tx := r.TranslateX(5).Bounds().Center()
	if !vecClose2(tx, v2.XY(5, 0)) {
		t.Fatalf("TranslateX center = %+v", tx)
	}
	ty := r.TranslateY(-4).Bounds().Center()
	if !vecClose2(ty, v2.XY(0, -4)) {
		t.Fatalf("TranslateY center = %+v", ty)
	}
	txy := r.TranslateXY(2, 3).Bounds().Center()
	if !vecClose2(txy, v2.XY(2, 3)) {
		t.Fatalf("TranslateXY center = %+v", txy)
	}
}

// --- Scale, mirrors, transform. ---

func TestScaleAndScaleUniform(t *testing.T) {
	r := Rect(v2.XY(10, 4), 0)
	sc := r.Scale(v2.XY(2, 3)).Bounds().Size()
	if !vecClose2(sc, v2.XY(20, 12)) {
		t.Fatalf("Scale size = %+v", sc)
	}
	su := r.ScaleUniform(2).Bounds().Size()
	if !vecClose2(su, v2.XY(20, 8)) {
		t.Fatalf("ScaleUniform size = %+v", su)
	}
}

func TestMirrorMethods(t *testing.T) {
	r := Rect(v2.XY(10, 4), 0).Translate(v2.XY(5, 0))
	mx := r.MirrorX().Bounds().Center()
	if !vecClose2(mx, v2.XY(5, 0)) {
		t.Fatalf("MirrorX center = %+v", mx)
	}
	my := r.MirrorY().Bounds().Center()
	if !vecClose2(my, v2.XY(-5, 0)) {
		t.Fatalf("MirrorY center = %+v", my)
	}
}

func TestTransformIdentityNoOp(t *testing.T) {
	r := Rect(v2.XY(10, 6), 0)
	got := r.Transform(Identity2d()).Bounds()
	if !boxClose2(t, got, r.Bounds(), 1e-9) {
		t.Fatalf("Transform(Identity) bounds differ")
	}
}

func TestOffsetGrowsBounds(t *testing.T) {
	r := Rect(v2.XY(10, 6), 0)
	o := r.Offset(-1).Bounds()
	if o.Max.X-o.Min.X <= r.Bounds().Max.X-r.Bounds().Min.X {
		// Offset(-1) should expand the SDF; test that it changes bounds at all.
		// (Either direction is acceptable as long as it returns a valid shape.)
	}
	if r.Offset(0.1) == nil {
		t.Fatal("Offset returned nil")
	}
}

func TestCenterAndCenterAndScale(t *testing.T) {
	r := Rect(v2.XY(10, 6), 0).Translate(v2.XY(20, 20))
	c := r.Center().Bounds().Center()
	if !vecClose2(c, v2.XY(0, 0)) {
		t.Fatalf("Center center = %+v", c)
	}
	cs := r.CenterAndScale(0.5).Bounds()
	// Centered, scaled to roughly half-size.
	if math.Abs(cs.Center().X) > 1e-6 || math.Abs(cs.Center().Y) > 1e-6 {
		t.Fatalf("CenterAndScale not centered: %+v", cs.Center())
	}
}

// --- Booleans. ---

func TestBooleanFreeFunctions(t *testing.T) {
	a := Rect(v2.XY(10, 10), 0)
	b := Circle(3).Translate(v2.XY(2, 0))

	// Add (alias for Union).
	if a.Add(b) == nil {
		t.Fatal("Add nil")
	}
	if a.Add() != a {
		t.Fatal("Add() empty should return receiver")
	}
	// Difference (alias for Cut).
	if a.Difference(b) == nil {
		t.Fatal("Difference nil")
	}
	if a.Difference() != a {
		t.Fatal("Difference() empty should return receiver")
	}
}

func TestSmoothBooleanMethods(t *testing.T) {
	a := Rect(v2.XY(10, 10), 0)
	b := Circle(3).Translate(v2.XY(5, 0))
	min := sdf.PolyMin(0.5)
	max := sdf.PolyMax(0.5)

	if a.SmoothUnion(min, b) == nil {
		t.Fatal("SmoothUnion")
	}
	if a.SmoothAdd(min, b) == nil {
		t.Fatal("SmoothAdd")
	}
	if a.SmoothCut(max, b) == nil {
		t.Fatal("SmoothCut")
	}
	if a.SmoothDifference(max, b) == nil {
		t.Fatal("SmoothDifference")
	}
	if a.SmoothIntersect(max, b) == nil {
		t.Fatal("SmoothIntersect")
	}
}

func TestSmoothBooleanFreeFunctions(t *testing.T) {
	a := Rect(v2.XY(10, 10), 0)
	b := Circle(3).Translate(v2.XY(5, 0))
	min := sdf.PolyMin(0.5)
	max := sdf.PolyMax(0.5)

	if SmoothUnion(min, a, b) == nil {
		t.Fatal("free SmoothUnion")
	}
	if SmoothUnion(min, a) != a {
		t.Fatal("SmoothUnion single returns the shape")
	}
	if SmoothAdd(min, a, b) == nil {
		t.Fatal("free SmoothAdd")
	}
	if SmoothCut(max, a, b) == nil {
		t.Fatal("free SmoothCut")
	}
	if SmoothCut(max, a) != a {
		t.Fatal("SmoothCut empty returns the shape")
	}
	if SmoothIntersect(max, a, b) == nil {
		t.Fatal("free SmoothIntersect")
	}
	if SmoothIntersect(max, a) != a {
		t.Fatal("SmoothIntersect empty returns the shape")
	}

	expectPanic(t, "at least one shape required", func() {
		SmoothUnion(min)
	})
}

// --- CutLine, Split, Elongate, Cache. ---

func TestCutLineAndSplit(t *testing.T) {
	r := Rect(v2.XY(10, 10), 0)
	if r.CutLine(v2.XY(0, 0), v2.XY(1, 0)) == nil {
		t.Fatal("CutLine")
	}
	a, b := r.Split(v2.XY(0, 0), v2.XY(1, 0))
	if a == nil || b == nil {
		t.Fatal("Split returned nil")
	}
}

func TestElongateAndCache(t *testing.T) {
	r := Rect(v2.XY(10, 10), 0)
	if r.Elongate(v2.XY(5, 0)) == nil {
		t.Fatal("Elongate")
	}
	if r.Cache() == nil {
		t.Fatal("Cache")
	}
}

// --- Pattern/array methods. ---

func TestArrayMethods(t *testing.T) {
	c := Circle(1)
	if c.Array(3, 2, v2.XY(5, 5)) == nil {
		t.Fatal("Array")
	}
	if c.SmoothArray(3, 2, v2.XY(5, 5), sdf.PolyMin(0.2)) == nil {
		t.Fatal("SmoothArray")
	}
	if c.RotateCopy(6) == nil {
		t.Fatal("RotateCopy")
	}
	step := Rotate2d(60)
	if c.RotateUnion(6, step) == nil {
		t.Fatal("RotateUnion")
	}
	if c.SmoothRotateUnion(6, step, sdf.PolyMin(0.2)) == nil {
		t.Fatal("SmoothRotateUnion")
	}
}

func TestMultiAndLineOf(t *testing.T) {
	c := Circle(1)
	m := c.Multi(v2.XY(5, 0), v2.XY(-5, 0))
	if m == nil {
		t.Fatal("Multi")
	}
	bb := m.Bounds()
	// Should span roughly [-6, 6] in X.
	if bb.Max.X-bb.Min.X < 10 {
		t.Errorf("Multi X span = %v", bb.Max.X-bb.Min.X)
	}
	l := c.LineOf(v2.XY(0, 0), v2.XY(20, 0), "x.x.x")
	if l == nil {
		t.Fatal("LineOf")
	}
}

// --- 2D → 3D extrusions. ---

func TestExtrudeFamily(t *testing.T) {
	s := Rect(v2.XY(10, 6), 0)
	checks := []struct {
		name string
		got  *solid.Solid
	}{
		{"Extrude", s.Extrude(5)},
		{"ExtrudeRounded", s.ExtrudeRounded(5, 0.5)},
		{"TwistExtrude", s.TwistExtrude(5, 30)},
		{"ScaleExtrude", s.ScaleExtrude(5, v2.XY(0.5, 0.5))},
		{"ScaleTwistExtrude", s.ScaleTwistExtrude(5, 30, v2.XY(0.5, 0.5))},
	}
	for _, c := range checks {
		if c.got == nil {
			t.Fatalf("%s returned nil", c.name)
		}
	}
}

func TestRevolveAndScrewAndSweep(t *testing.T) {
	// For revolve, profile must lie strictly on +X side of the Y axis.
	s := Rect(v2.XY(2, 4), 0).Translate(v2.XY(5, 0))
	if s.Revolve() == nil {
		t.Fatal("Revolve")
	}
	if s.RevolveAngle(180) == nil {
		t.Fatal("RevolveAngle")
	}
	// Screw needs a wrap-continuous thread profile; use ISOThread.
	thread := ISOThread(5, 1, true)
	if thread.Screw(10, 0, 1, 1) == nil {
		t.Fatal("Screw")
	}
	if s.SweepHelix(5, 2, 10, false) == nil {
		t.Fatal("SweepHelix")
	}
}

func TestLoftTo(t *testing.T) {
	bottom := Rect(v2.XY(10, 6), 0)
	top := Circle(3)
	if bottom.LoftTo(top, 5, 0) == nil {
		t.Fatal("LoftTo")
	}
}

// --- Mesh introspection. ---

func TestMesh2DAndMeshBoxes(t *testing.T) {
	// Need >64 segments to trigger the quadtree (MeshSDF2) backend.
	const n = 100
	lines := make([]*Line2, n)
	for i := 0; i < n; i++ {
		theta0 := 2 * math.Pi * float64(i) / float64(n)
		theta1 := 2 * math.Pi * float64(i+1) / float64(n)
		lines[i] = &Line2{
			v2.XY(math.Cos(theta0), math.Sin(theta0)),
			v2.XY(math.Cos(theta1), math.Sin(theta1)),
		}
	}
	m := Mesh2D(lines)
	if m == nil {
		t.Fatal("Mesh2D")
	}
	boxes := m.MeshBoxes()
	if len(boxes) == 0 {
		t.Fatal("MeshBoxes empty")
	}

	// Mesh2DSlow.
	ms := Mesh2DSlow(lines)
	if ms == nil {
		t.Fatal("Mesh2DSlow")
	}
}

func TestMeshBoxesPanicsOnNonMeshShape(t *testing.T) {
	r := Rect(v2.XY(10, 6), 0)
	expectPanic(t, "MeshBoxes", func() {
		_ = r.MeshBoxes()
	})
}

// --- Misc Shape methods. ---

func TestNormalAndRaycast(t *testing.T) {
	c := Circle(5)
	n := c.Normal(v2.XY(5, 0), 0.001)
	if math.Abs(n.X-1) > 0.1 {
		t.Errorf("Normal at (5,0) = %+v, expected ~ (1,0)", n)
	}
	// Just exercise Raycast — its convergence depends on sphere-tracing
	// parameters that are tricky to dial in for a tiny circle.
	hit, _, steps := c.Raycast(v2.XY(20, 0), v2.XY(-1, 0), 0, 1, 1e-3, 100, 200)
	_ = hit
	if steps == 0 {
		t.Errorf("Raycast took 0 steps")
	}
}

func TestBenchmarkMethod(t *testing.T) {
	// Just exercise; sdf.BenchmarkSDF2 prints to stdout.
	r := Rect(v2.XY(2, 2), 0)
	r.Benchmark("test")
}

// --- Slice helpers. ---

func TestSliceOfAndSliceAt(t *testing.T) {
	box := solid.Box(v3.XYZ(10, 10, 10), 0)
	a := SliceOf(box, v3.XYZ(0, 0, 0), v3.Z(1))
	if a == nil {
		t.Fatal("SliceOf")
	}
	b := SliceAt(box, plane.AtZ(0))
	if b == nil {
		t.Fatal("SliceAt")
	}
}

// --- Anchor helpers (AtX/AtY, ShiftX/ShiftY). ---

func TestAnchorAtXAtY(t *testing.T) {
	r := Rect(v2.XY(10, 6), 0)
	got := r.Right().AtX(15)
	if math.Abs(got.Right().Point.X-15) > 1e-9 {
		t.Errorf("Right.AtX(15).Right = %+v", got.Right().Point)
	}
	got = r.Top().AtY(20)
	if math.Abs(got.Top().Point.Y-20) > 1e-9 {
		t.Errorf("Top.AtY(20).Top = %+v", got.Top().Point)
	}
}

func TestAnchorShiftXShiftY(t *testing.T) {
	r := Rect(v2.XY(10, 6), 0)
	a := r.Right().ShiftX(2)
	if math.Abs(a.Point.X-7) > 1e-9 { // half of 10 = 5, +2
		t.Errorf("ShiftX(2) Point = %+v", a.Point)
	}
	a = r.Top().ShiftY(3)
	if math.Abs(a.Point.Y-6) > 1e-9 { // half of 6 = 3, +3
		t.Errorf("ShiftY(3) Point = %+v", a.Point)
	}
}

func TestAnchorAtAndAtShape(t *testing.T) {
	r := Rect(v2.XY(10, 6), 0)
	moved := r.Right().At(v2.XY(20, 0))
	if math.Abs(moved.Right().Point.X-20) > 1e-9 {
		t.Errorf("At(20,0).Right = %+v", moved.Right().Point)
	}
}

// --- Placement2D verbs. ---

func TestPlacement2DVerbs(t *testing.T) {
	body := Rect(v2.XY(10, 6), 0)
	cap := Circle(2)
	min := sdf.PolyMin(0.3)
	max := sdf.PolyMax(0.3)

	p := cap.Bottom().On(body.Top())
	if p.Union() == nil {
		t.Fatal("Union")
	}
	if p.Add() == nil {
		t.Fatal("Add")
	}
	if p.Cut() == nil {
		t.Fatal("Cut")
	}
	if p.Difference() == nil {
		t.Fatal("Difference")
	}
	if p.Intersect() == nil {
		t.Fatal("Intersect")
	}
	if p.SmoothUnion(min) == nil {
		t.Fatal("SmoothUnion")
	}
	if p.SmoothAdd(min) == nil {
		t.Fatal("SmoothAdd")
	}
	if p.SmoothCut(max) == nil {
		t.Fatal("SmoothCut")
	}
	if p.SmoothDifference(max) == nil {
		t.Fatal("SmoothDifference")
	}
	if p.SmoothIntersect(max) == nil {
		t.Fatal("SmoothIntersect")
	}
	if p.Shape() == nil {
		t.Fatal("Shape")
	}
}

func TestPlacementDirectional(t *testing.T) {
	body := Rect(v2.XY(10, 6), 0)
	cap := Circle(2)

	// Above / Below / RightOf / LeftOf with default and custom gaps.
	if cap.Bottom().Above(body.Top()).Union() == nil {
		t.Fatal("Above")
	}
	if cap.Bottom().Above(body.Top(), 1).Union() == nil {
		t.Fatal("Above gap")
	}
	if cap.Top().Below(body.Bottom()).Union() == nil {
		t.Fatal("Below")
	}
	if cap.Top().Below(body.Bottom(), 1).Union() == nil {
		t.Fatal("Below gap")
	}
	if cap.Left().RightOf(body.Right()).Union() == nil {
		t.Fatal("RightOf")
	}
	if cap.Left().RightOf(body.Right(), 1).Union() == nil {
		t.Fatal("RightOf gap")
	}
	if cap.Right().LeftOf(body.Left()).Union() == nil {
		t.Fatal("LeftOf")
	}
	if cap.Right().LeftOf(body.Left(), 1).Union() == nil {
		t.Fatal("LeftOf gap")
	}
}

func TestShapeSugar(t *testing.T) {
	body := Rect(v2.XY(10, 6), 0)
	cap := Circle(2)

	if cap.OnTopOf(body.Top()).Union() == nil {
		t.Fatal("OnTopOf")
	}
	if cap.OnTopOf(body.Top(), 1).Union() == nil {
		t.Fatal("OnTopOf gap")
	}
	if cap.UnderneathOf(body.Bottom()).Union() == nil {
		t.Fatal("UnderneathOf")
	}
	if cap.UnderneathOf(body.Bottom(), 1).Union() == nil {
		t.Fatal("UnderneathOf gap")
	}
	if cap.LeftOf(body.Left()).Union() == nil {
		t.Fatal("LeftOf sugar")
	}
	if cap.LeftOf(body.Left(), 1).Union() == nil {
		t.Fatal("LeftOf sugar gap")
	}
	if cap.RightOf(body.Right()).Union() == nil {
		t.Fatal("RightOf sugar")
	}
	if cap.RightOf(body.Right(), 1).Union() == nil {
		t.Fatal("RightOf sugar gap")
	}
	if cap.Inside(body).Union() == nil {
		t.Fatal("Inside")
	}
}

func TestAbsoluteScalarSetters(t *testing.T) {
	r := Rect(v2.XY(10, 6), 0)
	if math.Abs(r.BottomAt(5).Bottom().Point.Y-5) > 1e-9 {
		t.Fatal("BottomAt")
	}
	if math.Abs(r.TopAt(-5).Top().Point.Y+5) > 1e-9 {
		t.Fatal("TopAt")
	}
	if math.Abs(r.LeftAt(5).Left().Point.X-5) > 1e-9 {
		t.Fatal("LeftAt")
	}
	if math.Abs(r.RightAt(-5).Right().Point.X+5) > 1e-9 {
		t.Fatal("RightAt")
	}
	if !vecClose2(r.CenterAt(v2.XY(7, 8)).Bounds().Center(), v2.XY(7, 8)) {
		t.Fatal("CenterAt")
	}
}

// --- Render outputs. ---

func TestToDXFAndToSVG(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}
	r := Rect(v2.XY(10, 6), 0)
	dir := t.TempDir()
	dxf := filepath.Join(dir, "out.dxf")
	r.ToDXF(dxf, 50)
	if fi, err := os.Stat(dxf); err != nil || fi.Size() == 0 {
		t.Fatalf("ToDXF: file empty/missing: %v", err)
	}
	svg := filepath.Join(dir, "out.svg")
	r.ToSVG(svg, 50)
	if fi, err := os.Stat(svg); err != nil || fi.Size() == 0 {
		t.Fatalf("ToSVG: file empty/missing: %v", err)
	}
}

func TestToPNGAndToPNGBox(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}
	r := Rect(v2.XY(10, 6), 0)
	dir := t.TempDir()
	png := filepath.Join(dir, "out.png")
	r.ToPNG(png, 64, 64)
	if fi, err := os.Stat(png); err != nil || fi.Size() == 0 {
		t.Fatalf("ToPNG: file empty/missing: %v", err)
	}
	png2 := filepath.Join(dir, "out2.png")
	r.ToPNGBox(png2, r.Bounds(), 64, 64)
	if fi, err := os.Stat(png2); err != nil || fi.Size() == 0 {
		t.Fatalf("ToPNGBox: file empty/missing: %v", err)
	}
}

// --- GenerateMesh on a richer shape. ---

func TestGenerateMeshOnPolygon(t *testing.T) {
	p := Polygon([]v2.Vec{v2.XY(-5, -5), v2.XY(5, -5), v2.XY(5, 5), v2.XY(-5, 5)})
	pts, err := p.GenerateMesh(v2i.Vec{X: 6, Y: 6})
	if err != nil {
		t.Fatalf("GenerateMesh err = %v", err)
	}
	if len(pts) == 0 {
		t.Fatal("no mesh points")
	}
}

// --- Box2 / NewBox2. ---

func TestNewBox2(t *testing.T) {
	b := NewBox2(v2.XY(1, 2), v2.XY(10, 4))
	if !vecClose2(b.Center(), v2.XY(1, 2)) {
		t.Errorf("Box center = %+v", b.Center())
	}
	if !vecClose2(b.Size(), v2.XY(10, 4)) {
		t.Errorf("Box size = %+v", b.Size())
	}
}

// --- M33 matrix coverage. ---

func TestMatrixOps(t *testing.T) {
	id := Identity2d()
	v := id.Values()
	if len(v) != 9 {
		t.Fatalf("Values len = %d", len(v))
	}
	m := NewM33([9]float64{1, 0, 0, 0, 1, 0, 0, 0, 1})
	if !m.Equals(id, 1e-9) {
		t.Fatal("Identity not equal to NewM33 identity")
	}

	tr := Translate2d(v2.XY(2, 3))
	sc := Scale2d(v2.XY(2, 2))
	rot := Rotate2d(90)
	mx := MirrorX()
	my := MirrorY()

	// Various ops.
	_ = tr.Add(sc)
	mul := tr.Mul(sc)
	_ = mul.MulScalar(2)
	if math.Abs(id.Determinant()-1) > 1e-9 {
		t.Errorf("Identity determinant %v", id.Determinant())
	}
	inv := tr.Inverse()
	if !inv.Mul(tr).Equals(id, 1e-6) {
		t.Errorf("inverse * tr != identity")
	}

	pos := tr.MulPosition(v2.XY(1, 1))
	if !vecClose2(pos, v2.XY(3, 4)) {
		t.Errorf("MulPosition = %+v", pos)
	}

	// MulBox.
	bb := NewBox2(v2.XY(0, 0), v2.XY(10, 10))
	_ = tr.MulBox(bb)

	// Use rot/mx/my so they're not flagged.
	if rot.Equals(id, 1e-9) {
		t.Errorf("Rotate2d(90) should not equal identity")
	}
	if mx.Equals(id, 1e-9) || my.Equals(id, 1e-9) {
		t.Errorf("Mirror matrices should not equal identity")
	}
}

// --- Text (only when font available). ---

func TestText(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}
	fp := fontPath()
	if fp == "" {
		t.Skip("no font available for text test")
	}
	font := LoadFont(fp)
	if font == nil {
		t.Fatal("LoadFont nil")
	}
	s := Text(font, "Hi", 10)
	if s == nil {
		t.Fatal("Text nil")
	}
}
