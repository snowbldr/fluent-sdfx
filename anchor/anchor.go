// Package anchor finds points on a solid's surface by how they were
// found, so that a later step can name "where the ray from the front of
// the head meets the face" instead of a coordinate that was true once.
//
// A body anchor (Top, BottomFrontLeft, ...) is a point on the bounding box.
// A surface anchor is a point on the surface itself, with the surface
// normal there, resolved against the solid it is given every time it is
// asked. Change an earlier dimension and re-run, and the anchor moves with
// the surface; that is the whole point. Two ways to find one:
//
//   - Ray: the first place a ray crosses the surface.
//   - Nearest: the surface point nearest a given point.
//
// Both work on any field the library makes, including the distance bounds
// that warps and blends leave (sculpt.Pull, a twist, a smooth union): the
// search never trusts a magnitude to be a distance. Both panic when there
// is no answer (a ray that misses), because an anchor that silently falls
// back to somewhere else is the failure this package exists to prevent.
//
// Status: experimental, like sculpt. Not in docs/llms.txt yet.
package anchor

import (
	"fmt"
	"math"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Point is a point on a surface and the outward unit normal there.
type Point struct {
	Pos    v3.Vec `json:"pos"`    // on the surface
	Normal v3.Vec `json:"normal"` // outward unit normal, out of the material
}

// Offset returns the point d along the normal: positive d is outside the
// material, negative is inside it. Offset(-1) on a wall is a point 1 mm
// into the wall; Offset(2) on the inside face of a hollow is a point 2 mm
// into the void, which is how a fill seed is anchored to a surface.
func (p Point) Offset(d float64) v3.Vec { return p.Pos.Add(p.Normal.MulScalar(d)) }

// Ray returns the first place the ray from `from` along `dir` crosses the
// surface of s, from outside or from inside. A ray that starts on the
// surface returns its start. Panics when the ray leaves the solid's bounds
// without crossing, or when dir is zero.
func Ray(s *solid.Solid, from, dir v3.Vec) Point {
	if dir.Length() == 0 {
		panic("anchor.Ray: dir is zero; give the direction the ray travels")
	}
	d := dir.Normalize()
	bb := s.Bounds().Box
	diag := bb.Size().Length()
	eps := tol(diag)
	// Far enough to cross the whole box from wherever the ray starts.
	maxT := from.Sub(bb.Center()).Length() + diag
	// The field may be a bound much smaller than the distance, so a step is
	// never shorter than this; a feature thinner than it can be stepped over.
	minStep := diag / 2000

	f0 := eval(s, from)
	if math.Abs(f0) <= eps {
		return at(s, from, diag)
	}
	inside := f0 < 0
	t, f := 0.0, f0
	for t < maxT {
		tn := t + math.Max(math.Abs(f), minStep)
		fn := eval(s, from.Add(d.MulScalar(tn)))
		if (fn < 0) != inside || math.Abs(fn) <= eps {
			// Bisect the crossing between t and tn.
			lo, hi := t, tn
			for i := 0; i < 60 && hi-lo > eps*1e-3; i++ {
				mid := (lo + hi) / 2
				if (eval(s, from.Add(d.MulScalar(mid))) < 0) == inside {
					lo = mid
				} else {
					hi = mid
				}
			}
			return at(s, from.Add(d.MulScalar((lo+hi)/2)), diag)
		}
		t, f = tn, fn
	}
	panic(fmt.Sprintf("anchor.Ray: the ray from %s along %s misses the surface; aim it at the solid or start it nearer", fmtVec(from), fmtVec(d)))
}

// Toward is Ray aimed from `from` at the point `to`.
func Toward(s *solid.Solid, from, to v3.Vec) Point {
	if to.Sub(from).Length() == 0 {
		panic("anchor.Toward: from and to are the same point")
	}
	return Ray(s, from, to.Sub(from))
}

// Nearest returns the surface point of s nearest `near`. On an exact
// distance field that is one projection; on a bound (after a warp or a
// blend) the projection lands on the surface but not at the nearest point,
// and a few tangential steps walk it there. On a surface with several
// equally near points, which one comes back is not specified.
func Nearest(s *solid.Solid, near v3.Vec) Point {
	diag := s.Bounds().Box.Size().Length()
	eps := tol(diag)
	q := project(s, near, diag)
	for i := 0; i < 32; i++ {
		n := unit(grad(s, q, diag))
		v := near.Sub(q)
		tan := v.Sub(n.MulScalar(v.Dot(n)))
		if tan.Length() <= eps {
			break
		}
		// Half the tangential step: on a curved surface the full step
		// overshoots and the walk oscillates about the foot.
		q = project(s, q.Add(tan.MulScalar(0.5)), diag)
	}
	return at(s, q, diag)
}

// project walks p onto the surface along the field's gradient by Newton
// steps, which do not care how the field is scaled.
func project(s *solid.Solid, p v3.Vec, diag float64) v3.Vec {
	eps := tol(diag)
	for i := 0; i < 64; i++ {
		f := eval(s, p)
		if math.Abs(f) <= eps*1e-2 {
			return p
		}
		g := grad(s, p, diag)
		g2 := g.Length2()
		if g2 == 0 {
			// The centre of a sphere: every direction is as near. Step off
			// it a little and look again; three tries along the axes.
			if i < 3 {
				p = p.Add([]v3.Vec{v3.Z(-1), v3.X(1), v3.Y(1)}[i].MulScalar(diag * 1e-3))
				continue
			}
			panic(fmt.Sprintf("anchor: the field is flat at %s, so there is no direction to the surface", fmtVec(p)))
		}
		step := g.MulScalar(f / g2)
		// Never jump more than a tenth of the part in one step: a gradient
		// that is nearly flat would throw the point across it.
		if l := step.Length(); l > diag/10 {
			step = step.MulScalar(diag / 10 / l)
		}
		p = p.Sub(step)
	}
	return p
}

func at(s *solid.Solid, p v3.Vec, diag float64) Point {
	return Point{Pos: p, Normal: unit(grad(s, p, diag))}
}

func eval(s *solid.Solid, p v3.Vec) float64 { return s.Evaluate(p.Raw()) }

// grad is the central-difference gradient, unnormalised so a Newton step
// can use its length.
func grad(s *solid.Solid, p v3.Vec, diag float64) v3.Vec {
	h := tol(diag) * 10
	return v3.XYZ(
		eval(s, p.Add(v3.X(h)))-eval(s, p.Sub(v3.X(h))),
		eval(s, p.Add(v3.Y(h)))-eval(s, p.Sub(v3.Y(h))),
		eval(s, p.Add(v3.Z(h)))-eval(s, p.Sub(v3.Z(h))),
	).DivScalar(2 * h)
}

func unit(v v3.Vec) v3.Vec {
	if l := v.Length(); l > 0 {
		return v.DivScalar(l)
	}
	return v
}

// tol is the surface tolerance: a millionth of the part, never below a
// nanometre.
func tol(diag float64) float64 { return math.Max(diag*1e-6, 1e-9) }

func fmtVec(v v3.Vec) string { return fmt.Sprintf("(%.4g, %.4g, %.4g)", v.X, v.Y, v.Z) }
