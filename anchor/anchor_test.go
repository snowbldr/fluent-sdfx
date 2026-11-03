package anchor

import (
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/sculpt"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func TestRayOnSphereIsExact(t *testing.T) {
	s := solid.Sphere(10)
	p := Ray(s, v3.XYZ(0, -40, 0), v3.Y(1))
	if !p.Pos.Equals(v3.Y(-10), 1e-6) {
		t.Errorf("hit %v, want (0, -10, 0)", p.Pos)
	}
	if !p.Normal.Equals(v3.Y(-1), 1e-4) {
		t.Errorf("normal %v, want (0, -1, 0)", p.Normal)
	}
	// From inside, the first crossing is the way out.
	p = Ray(s, v3.Zero, v3.Z(1))
	if !p.Pos.Equals(v3.Z(10), 1e-6) {
		t.Errorf("from inside: hit %v, want (0, 0, 10)", p.Pos)
	}
	// Toward a point is the same ray.
	p = Toward(s, v3.XYZ(30, 0, 0), v3.Zero)
	if !p.Pos.Equals(v3.X(10), 1e-6) {
		t.Errorf("toward: hit %v, want (10, 0, 0)", p.Pos)
	}
}

func TestRayMissPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("a ray pointing away from the solid did not panic")
		}
	}()
	Ray(solid.Sphere(10), v3.XYZ(0, -40, 0), v3.Y(-1))
}

// A pulled solid's field is a bound scaled well below the distance; the
// ray must still land on the surface, not short of it.
func TestRayOnABound(t *testing.T) {
	s := sculpt.Pull(solid.Sphere(10), v3.Y(-10), v3.Y(-5), 10)
	p := Ray(s, v3.XYZ(0, -40, 0), v3.Y(1))
	validate.RequireNear(t, "pulled tip", p.Pos.Y, -15, 1e-4)
	validate.RequireNear(t, "field at the hit", validate.At(s, p.Pos), 0, 1e-6)
}

func TestNearest(t *testing.T) {
	s := solid.Sphere(10)
	p := Nearest(s, v3.XYZ(20, 20, 0))
	want := v3.XYZ(1, 1, 0).Normalize().MulScalar(10)
	if !p.Pos.Equals(want, 1e-6) {
		t.Errorf("nearest %v, want %v", p.Pos, want)
	}

	// On a bound, the projection alone lands somewhere on the surface; the
	// tangential walk has to bring it to the foot of the perpendicular.
	pulled := sculpt.Pull(solid.Sphere(10), v3.Y(-10), v3.XYZ(0, -4, 0), 9)
	near := v3.XYZ(9, -14, 3)
	q := Nearest(pulled, near)
	validate.RequireNear(t, "field at the nearest point", validate.At(pulled, q.Pos), 0, 1e-6)
	v := near.Sub(q.Pos)
	tan := v.Sub(q.Normal.MulScalar(v.Dot(q.Normal))).Length()
	t.Logf("nearest on the pulled sphere: %v, %.3f mm away, tangential residual %.2g", q.Pos, v.Length(), tan)
	if tan > 1e-3 {
		t.Errorf("the line to the nearest point is %.3g off the normal", tan)
	}
	// No surface point is nearer, checked on a ring of samples around it.
	for i := 0; i < 64; i++ {
		a := 2 * math.Pi * float64(i) / 64
		dir := v3.XYZ(math.Cos(a), 0, math.Sin(a))
		other := Ray(pulled, near, near.Sub(q.Pos).Neg().Add(dir.MulScalar(2)))
		if other.Pos.Sub(near).Length() < v.Length()-1e-6 {
			t.Errorf("surface point %v is nearer (%.4f < %.4f)", other.Pos, other.Pos.Sub(near).Length(), v.Length())
		}
	}
}

func TestOffset(t *testing.T) {
	p := Ray(solid.Sphere(10), v3.Zero, v3.X(1))
	if !p.Offset(2).Equals(v3.X(12), 1e-6) || !p.Offset(-2).Equals(v3.X(8), 1e-6) {
		t.Errorf("offsets %v and %v", p.Offset(2), p.Offset(-2))
	}
}

// The Phase 2.5 scenario, in Go: a figure whose head radius changes from
// 8 to 10. A nose pulled at the point a ray anchors still comes out of the
// face, a stem placed on head.Top still sits on top, a boss placed on the
// plate's top still sits flush. A pull at the coordinates that were the
// face when the head was 8 would now be inside the head and pull nothing.
func TestReflowFollowsAnchors(t *testing.T) {
	const nose = 4.0
	build := func(headR float64) (*solid.Solid, Point, *solid.Solid, *solid.Solid) {
		body := solid.Sphere(12)
		head := solid.Sphere(headR).TranslateZ(14)
		figure := solid.SmoothUnion(solid.RoundMin(3), body, head)
		face := Ray(figure, v3.XYZ(0, -40, 14), v3.Y(1))
		pulled := sculpt.Pull(figure, face.Pos, face.Normal.MulScalar(nose), 7)
		stem := solid.Cylinder(4, 1, 0).Bottom().On(head.Top()).Solid()
		return pulled, face, stem, head
	}
	for _, r := range []float64{8, 10} {
		pulled, face, stem, head := build(r)
		validate.RequireNear(t, "face found on the head", face.Pos.Y, -r, 1e-4)
		tip := face.Offset(nose)
		validate.RequireNear(t, "nose tip on the surface", validate.At(pulled, tip), 0, 1e-6)
		t.Logf("head %.0f: face at y=%.3f, nose tip at y=%.3f", r, face.Pos.Y, tip.Y)
		validate.RequireNear(t, "stem bottom on the head top", stem.Bounds().Box.Min.Z, head.Bounds().Box.Max.Z, 1e-9)
		validate.RequireNear(t, "stem top", stem.Bounds().Box.Max.Z, 14+r+4, 1e-9)
	}
	// The stale coordinate: the face of the r=8 head is inside the r=10 head.
	stale := v3.XYZ(0, -8, 14)
	big, _, _, _ := build(10)
	if !validate.Inside(big, stale) {
		t.Error("the r=8 face point should be inside the r=10 head; the scenario is not testing anything")
	}

	for _, thick := range []float64{4, 6} {
		plate := solid.Box(v3.XYZ(40, 30, thick), 0)
		boss := solid.Cylinder(6, 4, 0).Bottom().On(plate.Top()).Solid()
		validate.RequireNear(t, "boss flush on the plate", boss.Bounds().Box.Min.Z, plate.Bounds().Box.Max.Z, 1e-9)
	}
}
