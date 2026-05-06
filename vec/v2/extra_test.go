package v2_test

import (
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/vec/p2"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	"github.com/snowbldr/fluent-sdfx/vec/v2i"
	"github.com/snowbldr/sdfx/sdf"
	v2sdf "github.com/snowbldr/sdfx/vec/v2"
)

const eps = 1e-9

func close2(a, b v2.Vec, tol float64) bool {
	return math.Abs(a.X-b.X) <= tol && math.Abs(a.Y-b.Y) <= tol
}

func TestConstructorsAndZero(t *testing.T) {
	if got := v2.X(5); got != (v2.Vec{X: 5}) {
		t.Fatalf("X: %v", got)
	}
	if got := v2.Y(7); got != (v2.Vec{Y: 7}) {
		t.Fatalf("Y: %v", got)
	}
	if got := v2.XY(1, 2); got != (v2.Vec{X: 1, Y: 2}) {
		t.Fatalf("XY: %v", got)
	}
	if v2.Zero != (v2.Vec{}) {
		t.Fatalf("Zero")
	}
}

func TestRawAndArithmetic(t *testing.T) {
	a := v2.XY(2, 3)
	b := v2.XY(4, 5)
	if r := a.Raw(); r != (v2sdf.Vec{X: 2, Y: 3}) {
		t.Fatalf("Raw: %v", r)
	}
	if got := a.Add(b); got != v2.XY(6, 8) {
		t.Fatalf("Add: %v", got)
	}
	if got := a.AddScalar(1); got != v2.XY(3, 4) {
		t.Fatalf("AddScalar: %v", got)
	}
	if got := b.Sub(a); got != v2.XY(2, 2) {
		t.Fatalf("Sub: %v", got)
	}
	if got := b.SubScalar(1); got != v2.XY(3, 4) {
		t.Fatalf("SubScalar: %v", got)
	}
	if got := a.Mul(b); got != v2.XY(8, 15) {
		t.Fatalf("Mul: %v", got)
	}
	if got := a.MulScalar(2); got != v2.XY(4, 6) {
		t.Fatalf("MulScalar: %v", got)
	}
	if got := b.Div(a); got != v2.XY(2, 5.0/3.0) {
		t.Fatalf("Div: %v", got)
	}
	if got := a.DivScalar(2); got != v2.XY(1, 1.5) {
		t.Fatalf("DivScalar: %v", got)
	}
	if got := a.Neg(); got != v2.XY(-2, -3) {
		t.Fatalf("Neg: %v", got)
	}
	if got := v2.XY(-1.5, 2.5).Abs(); got != v2.XY(1.5, 2.5) {
		t.Fatalf("Abs: %v", got)
	}
	if got := v2.XY(1.2, 1.7).Ceil(); got != v2.XY(2, 2) {
		t.Fatalf("Ceil: %v", got)
	}
}

func TestDotCrossLengthNormalize(t *testing.T) {
	a := v2.XY(3, 4)
	b := v2.XY(1, 0)
	if d := a.Dot(b); d != 3 {
		t.Fatalf("Dot: %v", d)
	}
	if c := a.Cross(b); c != -4 {
		t.Fatalf("Cross: %v", c)
	}
	if l := a.Length(); math.Abs(l-5) > eps {
		t.Fatalf("Length: %v", l)
	}
	if l2 := a.Length2(); math.Abs(l2-25) > eps {
		t.Fatalf("Length2: %v", l2)
	}
	n := a.Normalize()
	if math.Abs(n.Length()-1) > eps {
		t.Fatalf("Normalize length: %v", n.Length())
	}
}

func TestEqualsAndComparators(t *testing.T) {
	a := v2.XY(1, 2)
	b := v2.XY(1+1e-12, 2-1e-12)
	if !a.Equals(b, 1e-9) {
		t.Fatalf("Equals true expected")
	}
	if v2.XY(1, 1).Equals(v2.XY(2, 2), 0.1) {
		t.Fatalf("Equals false expected")
	}
	if !v2.XY(-1, -1).LTZero() {
		t.Fatalf("LTZero")
	}
	if v2.XY(0, 0).LTZero() {
		t.Fatalf("LTZero exclusive of zero")
	}
	if !v2.XY(0, 0).LTEZero() {
		t.Fatalf("LTEZero")
	}
}

func TestMinMaxComponentsClamp(t *testing.T) {
	a := v2.XY(1, 5)
	b := v2.XY(3, 2)
	if got := a.Min(b); got != v2.XY(1, 2) {
		t.Fatalf("Min: %v", got)
	}
	if got := a.Max(b); got != v2.XY(3, 5) {
		t.Fatalf("Max: %v", got)
	}
	if got := a.MinComponent(); got != 1 {
		t.Fatalf("MinComponent: %v", got)
	}
	if got := a.MaxComponent(); got != 5 {
		t.Fatalf("MaxComponent: %v", got)
	}
	c := v2.XY(0, 0)
	d := v2.XY(2, 4)
	if got := v2.XY(-1, 5).Clamp(c, d); got != v2.XY(0, 4) {
		t.Fatalf("Clamp: %v", got)
	}
}

func TestConversions(t *testing.T) {
	a := v2.XY(1.7, 2.3)
	if got := a.ToV3(5); got.X != 1.7 || got.Y != 2.3 || got.Z != 5 {
		t.Fatalf("ToV3: %v", got)
	}
	if got := a.ToV2i(); got != v2i.XY(1, 2) {
		t.Fatalf("ToV2i: %v", got)
	}
	p := v2.XY(1, 0).ToP2()
	if math.Abs(p.R-1) > eps || math.Abs(p.Theta) > eps {
		t.Fatalf("ToP2: %v", p)
	}
	if got := v2.FromV2i(v2i.XY(3, 4)); got != v2.XY(3, 4) {
		t.Fatalf("FromV2i: %v", got)
	}
	if got := v2.FromP2(p2.RT(2, 0)); !close2(got, v2.XY(2, 0), eps) {
		t.Fatalf("FromP2: %v", got)
	}
}

func TestBoxBasics(t *testing.T) {
	b := v2.NewBox(v2.XY(0, 0), v2.XY(2, 4))
	if !close2(b.Center(), v2.Zero, eps) {
		t.Fatalf("Center: %v", b.Center())
	}
	if !close2(b.Size(), v2.XY(2, 4), eps) {
		t.Fatalf("Size: %v", b.Size())
	}
	if !b.Contains(v2.XY(0.5, 1)) {
		t.Fatalf("Contains")
	}
	if !b.Equals(b, eps) {
		t.Fatalf("Equals")
	}
	b2 := b.Translate(v2.XY(1, 1))
	if !close2(b2.Center(), v2.XY(1, 1), eps) {
		t.Fatalf("Translate: %v", b2.Center())
	}
	enl := b.Enlarge(v2.XY(1, 1))
	if !close2(enl.Size(), v2.XY(3, 5), eps) {
		t.Fatalf("Enlarge: %v", enl.Size())
	}
	ext := b.Extend(v2.NewBox(v2.XY(10, 10), v2.XY(2, 2)))
	if !ext.Contains(v2.XY(10, 10)) {
		t.Fatalf("Extend")
	}
	inc := b.Include(v2.XY(10, 0))
	if !inc.Contains(v2.XY(10, 0)) {
		t.Fatalf("Include")
	}
	sc := b.ScaleAboutCenter(2)
	if !close2(sc.Size(), v2.XY(4, 8), eps) {
		t.Fatalf("ScaleAboutCenter: %v", sc.Size())
	}
	sq := b.Square()
	if math.Abs(sq.Size().X-sq.Size().Y) > eps {
		t.Fatalf("Square: %v", sq.Size())
	}
	verts := b.Vertices()
	if len(verts) != 4 {
		t.Fatalf("Vertices len: %d", len(verts))
	}
	if a := b.Anchor(1, 1); !close2(a, v2.XY(1, 2), eps) {
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
	b := v2.NewBox(v2.XY(0, 0), v2.XY(2, 2))
	sdfBox := b.SDF()
	if got := v2.FromSDF(sdfBox); !got.Equals(b, eps) {
		t.Fatalf("FromSDF/SDF roundtrip: %v vs %v", got, b)
	}
	// Direct sdf.Box2 use as coverage of FromSDF too.
	_ = v2.FromSDF(sdf.Box2{Min: v2sdf.Vec{X: -1}, Max: v2sdf.Vec{X: 1, Y: 1}})
}

func TestLine2(t *testing.T) {
	l := v2.Line2{v2.XY(0, 0), v2.XY(3, 4)}
	bb := l.BoundingBox()
	if !close2(bb.Min, v2.XY(0, 0), eps) || !close2(bb.Max, v2.XY(3, 4), eps) {
		t.Fatalf("BoundingBox: %v", bb)
	}
	if l.Degenerate(eps) {
		t.Fatalf("Degenerate false expected")
	}
	deg := v2.Line2{v2.XY(1, 1), v2.XY(1, 1)}
	if !deg.Degenerate(eps) {
		t.Fatalf("Degenerate true expected")
	}
	r := l.Reverse()
	if r[0] != l[1] || r[1] != l[0] {
		t.Fatalf("Reverse: %v", r)
	}
	_ = l.SDF()
	_ = v2.FromSDFLine2(l.SDF())
}

func TestTriangle2(t *testing.T) {
	tr := v2.Triangle2{v2.XY(0, 0), v2.XY(1, 0), v2.XY(0, 1)}
	_ = tr.SDF()
	_ = v2.FromSDFTriangle2(tr.SDF())
}
