package shape

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Hole returns a circle of the given DIAMETER, centred on the origin.
// Circle takes a radius, but every hole a person specifies is a diameter —
// a drill size, an M3 clearance hole, a hook hole — so the halving has to
// happen somewhere. Doing it here means it cannot be dropped at the call
// site, which is the arithmetic slip that silently produces a hole twice
// the size it should be.
func Hole(diameter float64) *Shape { return Circle(diameter / 2) }

// Through extrudes the profile into a cutting tool that passes entirely
// through target along dir, standing pad millimetres proud at each end.
//
// This is the fix for trap 11, coincident faces. A cutting tool exactly as
// long as the body leaves a zero-thickness skin where their faces meet: the
// render looks right, the field is ambiguous exactly on the plane, and the
// slicer produces a sealed hole. Sizing the tool from the target's own
// bounding box means the length is never computed by hand and never goes
// stale when a thickness changes.
//
//	plate.Cut(Hole(3.2).TranslateX(12).Through(plate, v3.Z(1), 1))
//
// The profile lies in the plane perpendicular to dir, so for the usual dir
// of ±Z that is the XY plane and a profile positioned before the call lands
// where it was put. The tool is centred on the target along dir only; its
// position across the draw is exactly where the profile put it. For any
// other dir the profile's in-plane offset is rotated with it.
func (s *Shape) Through(target *solid.Solid, dir v3.Vec, pad float64) *solid.Solid {
	d := dir.Normalize()
	box := target.Bounds().Box
	h := box.Size().MulScalar(0.5)
	extent := 2 * (math.Abs(d.X)*h.X + math.Abs(d.Y)*h.Y + math.Abs(d.Z)*h.Z)

	// Extrude is centred on the origin, so the tool already straddles the
	// origin along Z; rotating about the origin keeps it straddling along
	// d. A dir of ±Z needs no rotation at all — the tool is symmetric, and
	// RotateToVector has no axis to turn about when the vectors are
	// parallel.
	tool := s.Extrude(extent + 2*pad)
	if math.Abs(d.Z) < 1-1e-12 {
		tool = tool.RotateToVector(v3.Z(1), d)
	}
	return tool.Translate(d.MulScalar(box.Center().Dot(d)))
}
