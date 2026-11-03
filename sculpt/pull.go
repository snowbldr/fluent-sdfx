package sculpt

import (
	"fmt"
	"math"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/sdfx/sdf"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
)

// pullSlope is the steepest slope of the falloff (1−t²)², at t = 1/√3. It
// bounds how much the warp stretches space, and so both the fold limit and
// the field scaling.
var pullSlope = 8 / (3 * math.Sqrt(3))

// Pull grabs the surface at `at` and drags it by `by`, as a sculpting brush
// does: the point at `at` lands exactly at at+by, the surface within
// `radius` of the landing point follows with a smooth falloff, and nothing
// farther away moves. A nose pulled from a face, a stem from the top of a
// gourd, a dent pushed in with a negative pull.
//
// The pull cannot be longer than about 0.65 × radius, or the surface
// would fold through itself; Pull panics past that limit. Pull in steps,
// or widen the radius.
//
// The result is a solid whose surface is exactly where it should be and
// whose field is a safe lower bound on the distance. Inside the brush the
// field is scaled down by a constant so it renders without holes, and
// near the brush it is held below the distance to the brush, so the
// magnitudes are wrong there: a Shell or Offset built straight on a
// pulled solid comes out several times too thick through the pull. Far
// from the brush (farther from it than from the surface) the field is the
// source's, unchanged, so a chain of pulls does not compound into a field
// too small everywhere for the mesher to skip empty space. Sample the
// result and Redistance before anything that reads magnitudes.
func Pull(s *solid.Solid, at, by v3.Vec, radius float64) *solid.Solid {
	if radius <= 0 {
		panic(fmt.Sprintf("sculpt.Pull: radius must be positive, got %g", radius))
	}
	l := by.Length()
	if l == 0 {
		return s
	}
	stretch := l * pullSlope / radius
	if stretch >= 1 {
		panic(fmt.Sprintf("sculpt.Pull: a pull of %.3g mm with radius %.3g folds the surface through itself; keep the pull under %.3g mm (0.65 × radius) or widen the radius", l, radius, radius/pullSlope))
	}
	// Only the ball of the brush around the landing point changes, so the
	// result lies within the source's bounds plus that ball. Extending by
	// the whole box moved by `by` grew the bounds by the full pull every
	// time, and a mesh or grid sized from the bounds got coarser with each
	// stroke.
	dest := at.Add(by)
	reach := v3.XYZ(radius, radius, radius)
	bb := s.Bounds().Box.Extend(v3.Box{Min: dest.Sub(reach), Max: dest.Add(reach)})
	return solid.Wrap(&pullSDF{
		s:     s.SDF3,
		dest:  dest.Raw(),
		by:    by.Raw(),
		r:     radius,
		r2:    radius * radius,
		scale: 1 / (1 + stretch),
		bb:    bb.SDF(),
	})
}

type pullSDF struct {
	s     sdf.SDF3
	dest  v3sdf.Vec // where the grabbed point lands
	by    v3sdf.Vec
	r, r2 float64
	scale float64
	bb    sdf.Box3
}

func (p *pullSDF) BoundingBox() sdf.Box3 { return p.bb }

// Evaluate reads the source field at the point that was pulled here. The
// falloff is centred on the landing point, so the landing point reads the
// source exactly at the grabbed point, and the warp's inverse is closed
// form.
//
// Outside the brush's ball the surface is the source's, and the only
// surface nearer than the source says is the pulled part inside the ball.
// A path from q to it crosses the ball at some b, where the inside field
// c·s(b) is already a bound, so the distance is at least
// c·|s(q)| + (1-c)·δ, δ being q's distance to the ball. The field outside
// is the smaller of that and |s(q)|: continuous with the inside at the
// ball, 1-Lipschitz (a min of two), and exactly the source once the ball
// is farther away than the surface.
func (p *pullSDF) Evaluate(q v3sdf.Vec) float64 {
	d2 := q.Sub(p.dest).Length2()
	if d2 >= p.r2 {
		f := p.s.Evaluate(q)
		a := math.Abs(f)
		if b := p.scale*a + (1-p.scale)*(math.Sqrt(d2)-p.r); b < a {
			a = b
		}
		return math.Copysign(a, f)
	}
	t2 := d2 / p.r2
	w := (1 - t2) * (1 - t2)
	return p.s.Evaluate(q.Sub(p.by.MulScalar(w))) * p.scale
}
