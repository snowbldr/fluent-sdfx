package sculpt

import (
	"math"
	"math/rand"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Scenario: a small figure, body and head blended, then sculpted: a nose
// pulled from the front of the face and a stem pulled from the top. The
// figure is then shelled for a hollow casting, which is where a pulled
// field bites unless it has been redistanced.
//
// Frame: +Z up, -Y is the front (the face).
const (
	bodyR  = 12.0
	headR  = 8.0
	headZ  = 14.0
	noseAt = headZ // nose grabbed at the front of the head
	noseBy = 4.0   // pulled forward
	noseR  = 7.0
	stemBy = 3.5
	stemR  = 6.0
)

func figure() *solid.Solid {
	body := solid.Sphere(bodyR)
	head := solid.Sphere(headR).TranslateZ(headZ)
	return solid.SmoothUnion(solid.RoundMin(3), body, head)
}

func sculptedFigure() *solid.Solid {
	f := figure()
	f = Pull(f, v3.XYZ(0, -headR, noseAt), v3.Y(-noseBy), noseR)
	f = Pull(f, v3.Z(headZ+headR), v3.Z(stemBy), stemR)
	return f
}

func TestFigurePullLandsExactly(t *testing.T) {
	plain := figure()
	f := sculptedFigure()

	noseTip := v3.XYZ(0, -headR-noseBy, noseAt)
	stemTip := v3.Z(headZ + headR + stemBy)

	// The grabbed surface points land exactly where asked.
	validate.RequireNear(t, "field at nose tip", validate.At(f, noseTip), 0, 1e-9)
	validate.RequireNear(t, "field at stem tip", validate.At(f, stemTip), 0, 1e-9)
	if !validate.Inside(f, noseTip.Add(v3.Y(0.3))) || validate.Inside(f, noseTip.Add(v3.Y(-0.3))) {
		t.Error("nose tip is not a front-facing surface")
	}
	if !validate.Inside(f, stemTip.Add(v3.Z(-0.3))) || validate.Inside(f, stemTip.Add(v3.Z(0.3))) {
		t.Error("stem tip is not an upward-facing surface")
	}

	// Measured from the head centre, the face is now noseBy farther out.
	headC := v3.Z(headZ)
	d, ok := validate.SurfaceAlong(f, headC, v3.Y(-1), 20, 0.05)
	if !ok {
		t.Fatal("no face surface")
	}
	validate.RequireNear(t, "face to nose tip", d, headR+noseBy, 0.02)

	// Outside both brushes nothing moved: back of head, bottom of body, sides.
	for _, c := range []struct {
		name      string
		from, dir v3.Vec
	}{
		{"back of head", headC, v3.Y(1)},
		{"bottom of body", v3.Zero, v3.Z(-1)},
		{"right of body", v3.Zero, v3.X(1)},
		{"left of head", headC, v3.X(-1)},
	} {
		want, _ := validate.SurfaceAlong(plain, c.from, c.dir, 20, 0.05)
		got, ok := validate.SurfaceAlong(f, c.from, c.dir, 20, 0.05)
		if !ok {
			t.Fatalf("%s: no surface", c.name)
		}
		validate.RequireNear(t, c.name+" unchanged", got, want, 1e-6)
	}

	// The seam between body and head still has its blend.
	seam := v3.XYZ(0, 0, 9.86) // where the two spheres meet
	if validate.At(f, v3.XYZ(0, -7.2, seam.Z)) > 0 {
		t.Error("blend fillet at the neck is gone")
	}
}

func TestFigurePullRefusesToFold(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("a pull of 5 at radius 7 did not panic")
		} else {
			t.Logf("Pull refused: %v", r)
		}
	}()
	Pull(figure(), v3.XYZ(0, -headR, headZ), v3.Y(-5), 7)
}

// The pulled field is a bound near the brushes, and the test says so in
// numbers: a shell built on it is several times too thick through a pull,
// the same shell on the redistanced grid is the stated thickness, and
// through the untouched body, far from any brush, both are right.
func TestFigureShellNeedsRedistance(t *testing.T) {
	f := sculptedFigure()
	const wall = 1.6

	g := Sample(f, 128)
	h := g.CellSize()
	t.Logf("cell %.3f mm, %v cells", h, g.Cells())
	if n := len(g.Components()); n != 1 {
		t.Fatalf("sculpted figure is %d bodies", n)
	}
	r := g.Redistance()
	rs := r.Solid()

	// Gradient magnitude of a true distance field is 1. Sample it a
	// millimetre inside the surface: deeper, a thin feature's medial axis
	// kinks the field and central differences read low there by design.
	for _, p := range []v3.Vec{
		v3.XYZ(0, -headR-noseBy+1, noseAt), // just inside the nose tip
		v3.XYZ(0, -headR+1, noseAt+2),      // the cheek the nose was pulled from
		v3.Z(headZ + headR + stemBy - 1),   // just inside the stem tip
		v3.XYZ(0, 0, -6),                   // deep in the body, far from any brush
	} {
		m := gradMag(rs, p, h)
		if math.Abs(m-1) > 0.12 {
			t.Errorf("|∇| at %v = %.3f, want 1 ± 0.12", p, m)
		}
	}

	rawHollow := f.Shell(wall)
	goodHollow := rs.Shell(wall)
	for _, c := range []struct {
		name   string
		p, dir v3.Vec
		pulled bool
	}{
		{"through the nose", v3.XYZ(0, -headR-noseBy, noseAt), v3.Y(1), true},
		{"through the body side", v3.X(-bodyR), v3.X(1), false},
		{"through the stem", v3.Z(headZ + headR + stemBy), v3.Z(1), true},
	} {
		good, ok := validate.WallThickness(goodHollow, c.p, c.dir, 10)
		if !ok {
			t.Errorf("%s: redistanced wall not measurable", c.name)
			continue
		}
		raw, _ := validate.WallThickness(rawHollow, c.p, c.dir, 10)
		t.Logf("%s: stated %.2f, raw pulled %.3f, redistanced %.3f", c.name, wall, raw, good)
		validate.RequireNear(t, c.name+" wall", good, wall, 0.25)
		switch {
		case c.pulled && raw < 2*wall:
			t.Errorf("%s: the raw pulled shell (%.3f) is not the trap it should be (expected > %.1f)", c.name, raw, 2*wall)
		case !c.pulled:
			validate.RequireNear(t, c.name+" raw wall, far from a brush", raw, wall, 0.05)
		}
	}
}

// The pulled field must stay a bound: never steeper than 1 (or a sphere
// tracer steps through the surface) and never more than the distance to
// the redistanced surface. Checked on random pairs of nearby points and
// against the redistanced grid, after a chain of strokes.
func TestPullFieldIsABound(t *testing.T) {
	f := sculptedFigure()
	for i := 0; i < 6; i++ {
		a := float64(i)
		f = Pull(f, v3.XYZ(10*math.Cos(a), 10*math.Sin(a), 2*a-4), v3.XYZ(math.Cos(a), math.Sin(a), 0.2).MulScalar(2), 5)
	}
	r := Sample(f, 128).Redistance().Solid()
	bb := f.Bounds()
	rng := rand.New(rand.NewSource(3))
	steep, over := 0.0, 0.0
	for i := 0; i < 20000; i++ {
		s := bb.Size()
		p := bb.Min.Add(v3.XYZ(rng.Float64()*s.X, rng.Float64()*s.Y, rng.Float64()*s.Z))
		q := p.Add(v3.XYZ(rng.NormFloat64(), rng.NormFloat64(), rng.NormFloat64()).MulScalar(0.5))
		fp := validate.At(f, p)
		if l := math.Abs(fp-validate.At(f, q)) / p.Sub(q).Length(); l > steep {
			steep = l
		}
		if d := math.Abs(fp) - math.Abs(validate.At(r, p)); d > over {
			over = d
		}
	}
	h := bb.Size().MaxComponent() / 128
	t.Logf("steepest slope %.4f; most above the true distance %.3f mm (a cell is %.3f)", steep, over, h)
	if steep > 1+1e-6 {
		t.Errorf("the field is steeper than 1 (%.4f)", steep)
	}
	if over > h {
		t.Errorf("the field overstates the distance by %.3f mm", over)
	}
}

// gradMag is the magnitude of the field's gradient at p by central
// differences one cell wide.
func gradMag(s *solid.Solid, p v3.Vec, eps float64) float64 {
	dx := validate.At(s, p.Add(v3.X(eps))) - validate.At(s, p.Sub(v3.X(eps)))
	dy := validate.At(s, p.Add(v3.Y(eps))) - validate.At(s, p.Sub(v3.Y(eps)))
	dz := validate.At(s, p.Add(v3.Z(eps))) - validate.At(s, p.Sub(v3.Z(eps)))
	return math.Sqrt(dx*dx+dy*dy+dz*dz) / (2 * eps)
}

// Repeated strokes must not inflate the bounds: only the brush's ball
// moves, so the bounds grow to cover the pulled tips and no further. A
// mesh or grid sized from the bounds would otherwise coarsen with every
// stroke.
func TestPullBoundsStayTight(t *testing.T) {
	f := figure()
	before := f.Bounds().Size()
	for i := 0; i < 12; i++ {
		f = Pull(f, v3.XYZ(0, -headR, headZ), v3.Y(-0.5), noseR) // a dab in the same spot
	}
	after := f.Bounds().Size()
	t.Logf("bounds before %v, after twelve dabs %v", before, after)
	if after.X > before.X+1e-9 || after.Z > before.Z+1e-9 {
		t.Errorf("bounds grew across X or Z (%v to %v) for pulls along -Y inside the form", before, after)
	}
	// The pulls only ever go out along -Y, and only within the brush ball.
	if grow := after.Y - before.Y; grow > noseR {
		t.Errorf("bounds grew %.2f mm along Y, more than the brush radius %.1f", grow, noseR)
	}
}
