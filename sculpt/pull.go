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
// whose field is a safe lower bound on the distance, scaled down by a
// constant so it renders without holes. The magnitudes are therefore wrong
// everywhere, not just near the pull: a Shell or Offset built straight on
// a pulled solid comes out several times too thick. Sample the result and
// Redistance before anything that reads magnitudes.
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
	bb := s.Bounds().Box
	bb = bb.Extend(bb.Translate(by))
	return solid.Wrap(&pullSDF{
		s:     s.SDF3,
		dest:  at.Add(by).Raw(),
		by:    by.Raw(),
		r2:    radius * radius,
		scale: 1 / (1 + stretch),
		bb:    bb.SDF(),
	})
}

type pullSDF struct {
	s     sdf.SDF3
	dest  v3sdf.Vec // where the grabbed point lands
	by    v3sdf.Vec
	r2    float64
	scale float64
	bb    sdf.Box3
}

func (p *pullSDF) BoundingBox() sdf.Box3 { return p.bb }

// Evaluate reads the source field at the point that was pulled here. The
// falloff is centred on the landing point, so the landing point reads the
// source exactly at the grabbed point, and the warp's inverse is closed
// form.
func (p *pullSDF) Evaluate(q v3sdf.Vec) float64 {
	t2 := q.Sub(p.dest).Length2() / p.r2
	if t2 >= 1 {
		return p.s.Evaluate(q) * p.scale
	}
	w := (1 - t2) * (1 - t2)
	return p.s.Evaluate(q.Sub(p.by.MulScalar(w))) * p.scale
}
