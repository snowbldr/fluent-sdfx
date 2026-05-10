package solid

import (
	"testing"

	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// --- AnchoredSolid.Attach* (receiver = base, argument = part that moves) ---

func TestAttachAlignsAnchors(t *testing.T) {
	body := Cylinder(10, 5, 0)
	cap := Sphere(3)
	// Receiver-stays form: body's top is the base; cap's bottom moves to it.
	moved := body.Top().Attach(cap.Bottom()).Solid()
	got := moved.Bottom().Point
	want := body.Top().Point
	if !vecClose(got, want) {
		t.Fatalf("Attach: moved.Bottom() = %+v, want body.Top() %+v", got, want)
	}
}

func TestAttachIsInverseOfOn(t *testing.T) {
	// For every Attach* there is an equivalent inverted On/RightOf/etc.
	// — a.AttachX(b, gap) ≡ b.X(a, gap). Verify by anchor-point parity
	// across all 7 verbs.
	a := Box(v3.XYZ(4, 4, 4), 0).Translate(v3.XYZ(10, 5, 2))
	b := Sphere(2)
	const gap = 1.5
	cases := []struct {
		name string
		got  Placement
		want Placement
	}{
		{"Attach", a.Right().Attach(b.Left()), b.Left().On(a.Right())},
		{"AttachAbove", a.Top().AttachAbove(b.Bottom(), gap), b.Bottom().Above(a.Top(), gap)},
		{"AttachBelow", a.Bottom().AttachBelow(b.Top(), gap), b.Top().Below(a.Bottom(), gap)},
		{"AttachRight", a.Right().AttachRight(b.Left(), gap), b.Left().RightOf(a.Right(), gap)},
		{"AttachLeft", a.Left().AttachLeft(b.Right(), gap), b.Right().LeftOf(a.Left(), gap)},
		{"AttachBehind", a.Back().AttachBehind(b.Front(), gap), b.Front().Behind(a.Back(), gap)},
		{"AttachInFront", a.Front().AttachInFront(b.Back(), gap), b.Back().InFrontOf(a.Front(), gap)},
	}
	for _, c := range cases {
		gotMoved := c.got.Moved.Bounds().Center()
		wantMoved := c.want.Moved.Bounds().Center()
		if !vecClose(gotMoved, wantMoved) {
			t.Errorf("%s: moved center %+v, want %+v", c.name, gotMoved, wantMoved)
		}
		// Base should always be a (the receiver of Attach*) — that's
		// what makes this readable for chained assemblies.
		if c.got.Base != a {
			t.Errorf("%s: Base = %p, want a (%p)", c.name, c.got.Base, a)
		}
	}
}

func TestAttachAboveGap(t *testing.T) {
	plate := Box(v3.XYZ(20, 20, 2), 0)
	post := Cylinder(5, 1, 0)
	moved := plate.Top().AttachAbove(post.Bottom(), 3).Solid()
	// Plate top is at z=1; post bottom should be at z=1+3=4.
	if got := moved.Bottom().Point.Z; got != 4 {
		t.Fatalf("post bottom Z = %v, want 4", got)
	}
}

func TestAttachRightGap(t *testing.T) {
	a := Box(v3.XYZ(4, 4, 4), 0) // bbox right at x=2
	b := Sphere(1)
	moved := a.Right().AttachRight(b.Left(), 0.5).Solid()
	// b.Left should land at a.Right + 0.5 = 2 + 0.5 = 2.5.
	if got := moved.Left().Point.X; got != 2.5 {
		t.Fatalf("b.Left X = %v, want 2.5", got)
	}
}

// --- Solid.Attach* sugar (defaults touching faces) ---

func TestAttachOnTopMatchesAnchor(t *testing.T) {
	// s.AttachOnTop(part, g) ≡ s.Top().AttachAbove(part.Bottom(), g)
	s := Box(v3.XYZ(10, 10, 4), 0)
	p := Sphere(2)
	got := s.AttachOnTop(p, 1).Solid().Bottom().Point
	want := s.Top().AttachAbove(p.Bottom(), 1).Solid().Bottom().Point
	if !vecClose(got, want) {
		t.Fatalf("AttachOnTop ≠ Top().AttachAbove: %+v vs %+v", got, want)
	}
}

func TestAttachSugarLandsOnExpectedFace(t *testing.T) {
	// Each AttachX sugar should put the part flush with the touching
	// face of s, then offset by gap along the natural axis.
	s := Box(v3.XYZ(10, 10, 4), 0)
	p := Box(v3.XYZ(2, 2, 2), 0)
	const g = 0.25
	cases := []struct {
		name string
		moved *Solid
		face string // which anchor of moved should match
		want v3.Vec
	}{
		{"AttachOnTop", s.AttachOnTop(p, g).Solid(), "bottom", v3.XYZ(0, 0, 2+g)},
		{"AttachUnderneath", s.AttachUnderneath(p, g).Solid(), "top", v3.XYZ(0, 0, -2-g)},
		{"AttachRight", s.AttachRight(p, g).Solid(), "left", v3.XYZ(5+g, 0, 0)},
		{"AttachLeft", s.AttachLeft(p, g).Solid(), "right", v3.XYZ(-5-g, 0, 0)},
		{"AttachBehind", s.AttachBehind(p, g).Solid(), "front", v3.XYZ(0, 5+g, 0)},
		{"AttachInFront", s.AttachInFront(p, g).Solid(), "back", v3.XYZ(0, -5-g, 0)},
	}
	for _, c := range cases {
		var got v3.Vec
		switch c.face {
		case "bottom":
			got = c.moved.Bottom().Point
		case "top":
			got = c.moved.Top().Point
		case "left":
			got = c.moved.Left().Point
		case "right":
			got = c.moved.Right().Point
		case "front":
			got = c.moved.Front().Point
		case "back":
			got = c.moved.Back().Point
		}
		if !vecClose(got, c.want) {
			t.Errorf("%s: %s anchor = %+v, want %+v", c.name, c.face, got, c.want)
		}
	}
}

func TestAttachUnionMatchesOn(t *testing.T) {
	// The whole point of Attach* is that .Union() composes the same
	// way as the inverted On chain. Sample the resulting SDF on a few
	// points to verify they're geometrically identical.
	a := Box(v3.XYZ(4, 4, 4), 0)
	b := Sphere(1.5)
	const gap = 0.5
	left := a.Top().AttachAbove(b.Bottom(), gap).Union()
	right := b.Bottom().Above(a.Top(), gap).Union()
	for _, p := range []v3.Vec{
		v3.XYZ(0, 0, 0), v3.XYZ(0, 0, 3), v3.XYZ(2, 0, 2.5), v3.XYZ(-3, -3, -3),
	} {
		la := left.Evaluate(p.Raw())
		lb := right.Evaluate(p.Raw())
		if abs(la-lb) > 1e-9 {
			t.Errorf("at %+v: left %v, right %v", p, la, lb)
		}
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
