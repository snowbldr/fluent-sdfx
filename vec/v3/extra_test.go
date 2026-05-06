package v3_test

import (
	"math"
	"testing"

	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/fluent-sdfx/vec/v3i"
	"github.com/snowbldr/sdfx/sdf"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
)

const eps = 1e-9

func close3(a, b v3.Vec, tol float64) bool {
	return math.Abs(a.X-b.X) <= tol && math.Abs(a.Y-b.Y) <= tol && math.Abs(a.Z-b.Z) <= tol
}

func TestConstructorsAndZero(t *testing.T) {
	if v3.X(1) != (v3.Vec{X: 1}) {
		t.Fatal("X")
	}
	if v3.Y(2) != (v3.Vec{Y: 2}) {
		t.Fatal("Y")
	}
	if v3.Z(3) != (v3.Vec{Z: 3}) {
		t.Fatal("Z")
	}
	if v3.XY(1, 2) != (v3.Vec{X: 1, Y: 2}) {
		t.Fatal("XY")
	}
	if v3.XZ(1, 3) != (v3.Vec{X: 1, Z: 3}) {
		t.Fatal("XZ")
	}
	if v3.YZ(2, 3) != (v3.Vec{Y: 2, Z: 3}) {
		t.Fatal("YZ")
	}
	if v3.XYZ(1, 2, 3) != (v3.Vec{X: 1, Y: 2, Z: 3}) {
		t.Fatal("XYZ")
	}
	if v3.Zero != (v3.Vec{}) {
		t.Fatal("Zero")
	}
}

func TestRawAndArithmetic(t *testing.T) {
	a := v3.XYZ(1, 2, 3)
	b := v3.XYZ(4, 5, 6)
	if r := a.Raw(); r != (v3sdf.Vec{X: 1, Y: 2, Z: 3}) {
		t.Fatalf("Raw")
	}
	if got := a.Add(b); got != v3.XYZ(5, 7, 9) {
		t.Fatalf("Add")
	}
	if got := a.AddScalar(1); got != v3.XYZ(2, 3, 4) {
		t.Fatalf("AddScalar")
	}
	if got := b.Sub(a); got != v3.XYZ(3, 3, 3) {
		t.Fatalf("Sub")
	}
	if got := b.SubScalar(1); got != v3.XYZ(3, 4, 5) {
		t.Fatalf("SubScalar")
	}
	if got := a.Mul(b); got != v3.XYZ(4, 10, 18) {
		t.Fatalf("Mul")
	}
	if got := a.MulScalar(2); got != v3.XYZ(2, 4, 6) {
		t.Fatalf("MulScalar")
	}
	if got := b.Div(a); !close3(got, v3.XYZ(4, 2.5, 2), eps) {
		t.Fatalf("Div: %v", got)
	}
	if got := a.DivScalar(2); got != v3.XYZ(0.5, 1, 1.5) {
		t.Fatalf("DivScalar")
	}
	if got := a.Neg(); got != v3.XYZ(-1, -2, -3) {
		t.Fatalf("Neg")
	}
	if got := v3.XYZ(-1.5, 2.5, -3.5).Abs(); got != v3.XYZ(1.5, 2.5, 3.5) {
		t.Fatalf("Abs")
	}
	if got := v3.XYZ(1.2, 1.7, 2.1).Ceil(); got != v3.XYZ(2, 2, 3) {
		t.Fatalf("Ceil")
	}
}

func TestSinCos(t *testing.T) {
	v := v3.XYZ(0, math.Pi/2, math.Pi)
	s := v.Sin()
	c := v.Cos()
	if math.Abs(s.X) > eps || math.Abs(s.Y-1) > eps || math.Abs(s.Z) > eps {
		t.Fatalf("Sin: %v", s)
	}
	if math.Abs(c.X-1) > eps || math.Abs(c.Y) > eps || math.Abs(c.Z+1) > eps {
		t.Fatalf("Cos: %v", c)
	}
}

func TestDotCrossLengthNormalize(t *testing.T) {
	a := v3.XYZ(1, 0, 0)
	b := v3.XYZ(0, 1, 0)
	if d := a.Dot(b); d != 0 {
		t.Fatalf("Dot: %v", d)
	}
	c := a.Cross(b)
	if !close3(c, v3.XYZ(0, 0, 1), eps) {
		t.Fatalf("Cross: %v", c)
	}
	v := v3.XYZ(2, 3, 6)
	if l := v.Length(); math.Abs(l-7) > eps {
		t.Fatalf("Length: %v", l)
	}
	if l2 := v.Length2(); math.Abs(l2-49) > eps {
		t.Fatalf("Length2: %v", l2)
	}
	n := v.Normalize()
	if math.Abs(n.Length()-1) > eps {
		t.Fatalf("Normalize length: %v", n.Length())
	}
}

func TestEqualsAndComparators(t *testing.T) {
	if !v3.XYZ(1, 2, 3).Equals(v3.XYZ(1+1e-12, 2, 3), 1e-9) {
		t.Fatalf("Equals true expected")
	}
	if v3.XYZ(1, 1, 1).Equals(v3.XYZ(2, 2, 2), 0.1) {
		t.Fatalf("Equals false expected")
	}
	if !v3.XYZ(-1, -1, -1).LTZero() {
		t.Fatalf("LTZero")
	}
	if v3.XYZ(0, 0, 0).LTZero() {
		t.Fatalf("LTZero exclusive of zero")
	}
	if !v3.XYZ(0, 0, 0).LTEZero() {
		t.Fatalf("LTEZero")
	}
}

func TestMinMaxComponentsClamp(t *testing.T) {
	a := v3.XYZ(1, 5, 2)
	b := v3.XYZ(3, 2, 4)
	if got := a.Min(b); got != v3.XYZ(1, 2, 2) {
		t.Fatalf("Min: %v", got)
	}
	if got := a.Max(b); got != v3.XYZ(3, 5, 4) {
		t.Fatalf("Max: %v", got)
	}
	if got := a.MinComponent(); got != 1 {
		t.Fatalf("MinComponent")
	}
	if got := a.MaxComponent(); got != 5 {
		t.Fatalf("MaxComponent")
	}
	if got := v3.XYZ(-1, 5, 0).Clamp(v3.XYZ(0, 0, 0), v3.XYZ(2, 4, 1)); got != v3.XYZ(0, 4, 0) {
		t.Fatalf("Clamp: %v", got)
	}
}

func TestGetSet(t *testing.T) {
	v := v3.XYZ(1, 2, 3)
	if v.Get(0) != 1 || v.Get(1) != 2 || v.Get(2) != 3 {
		t.Fatalf("Get")
	}
	v.Set(0, 10)
	v.Set(1, 20)
	v.Set(2, 30)
	if v != v3.XYZ(10, 20, 30) {
		t.Fatalf("Set: %v", v)
	}
}

func TestConversions(t *testing.T) {
	a := v3.XYZ(1.7, 2.3, -3.9)
	if got := a.ToV3i(); got != v3i.XYZ(1, 2, -3) {
		t.Fatalf("ToV3i: %v", got)
	}
	if got := v3.FromV3i(v3i.XYZ(1, 2, 3)); got != v3.XYZ(1, 2, 3) {
		t.Fatalf("FromV3i: %v", got)
	}
}

func TestBoxBasics(t *testing.T) {
	b := v3.NewBox(v3.XYZ(0, 0, 0), v3.XYZ(2, 4, 6))
	if !close3(b.Center(), v3.Zero, eps) {
		t.Fatalf("Center: %v", b.Center())
	}
	if !close3(b.Size(), v3.XYZ(2, 4, 6), eps) {
		t.Fatalf("Size: %v", b.Size())
	}
	if !b.Contains(v3.XYZ(0.5, 1, 2)) {
		t.Fatalf("Contains")
	}
	if !b.Equals(b, eps) {
		t.Fatalf("Equals")
	}
	tr := b.Translate(v3.XYZ(1, 1, 1))
	if !close3(tr.Center(), v3.XYZ(1, 1, 1), eps) {
		t.Fatalf("Translate")
	}
	enl := b.Enlarge(v3.XYZ(1, 1, 1))
	if !close3(enl.Size(), v3.XYZ(3, 5, 7), eps) {
		t.Fatalf("Enlarge")
	}
	ext := b.Extend(v3.NewBox(v3.XYZ(10, 10, 10), v3.XYZ(2, 2, 2)))
	if !ext.Contains(v3.XYZ(10, 10, 10)) {
		t.Fatalf("Extend")
	}
	inc := b.Include(v3.XYZ(10, 0, 0))
	if !inc.Contains(v3.XYZ(10, 0, 0)) {
		t.Fatalf("Include")
	}
	sc := b.ScaleAboutCenter(2)
	if !close3(sc.Size(), v3.XYZ(4, 8, 12), eps) {
		t.Fatalf("ScaleAboutCenter: %v", sc.Size())
	}
	cu := b.Cube()
	if math.Abs(cu.Size().X-cu.Size().Y) > eps || math.Abs(cu.Size().Y-cu.Size().Z) > eps {
		t.Fatalf("Cube: %v", cu.Size())
	}
	verts := b.Vertices()
	if len(verts) != 8 {
		t.Fatalf("Vertices len: %d", len(verts))
	}
	if a := b.Anchor(1, 1, 1); !close3(a, v3.XYZ(1, 2, 3), eps) {
		t.Fatalf("Anchor: %v", a)
	}
	r := b.Random()
	if !b.Contains(r) {
		t.Fatalf("Random outside box: %v", r)
	}
	rs := b.RandomSet(5)
	if len(rs) != 5 {
		t.Fatalf("RandomSet len")
	}
}

func TestBoxSDFFromSDF(t *testing.T) {
	b := v3.NewBox(v3.XYZ(0, 0, 0), v3.XYZ(2, 2, 2))
	sdfBox := b.SDF()
	if got := v3.FromSDF(sdfBox); !got.Equals(b, eps) {
		t.Fatalf("FromSDF/SDF roundtrip")
	}
	_ = v3.FromSDF(sdf.Box3{Min: v3sdf.Vec{X: -1}, Max: v3sdf.Vec{X: 1, Y: 1, Z: 1}})
}

func TestTriangle3(t *testing.T) {
	tr := v3.Triangle3{v3.XYZ(0, 0, 0), v3.XYZ(1, 0, 0), v3.XYZ(0, 1, 0)}
	n := tr.Normal()
	if !close3(n, v3.XYZ(0, 0, 1), eps) {
		t.Fatalf("Normal: %v", n)
	}
	if tr.Degenerate(eps) {
		t.Fatalf("Degenerate false expected")
	}
	deg := v3.Triangle3{v3.XYZ(0, 0, 0), v3.XYZ(0, 0, 0), v3.XYZ(1, 1, 1)}
	if !deg.Degenerate(eps) {
		t.Fatalf("Degenerate true expected")
	}
	_ = tr.SDF()
	_ = v3.FromSDFTriangle3(tr.SDF())
}
