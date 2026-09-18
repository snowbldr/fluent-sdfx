// Motion clearance: transform one body's surface through a family of poses
// and read the other body's field at every transformed point.
//
// A static clearance check only proves the assembly fits where it is drawn.
// Sweep proves it still fits everywhere it moves — the lid opening, the
// screw backing out, the bore shifting within its tolerance band.
//
// The pose generators are all right-handed about +Z: a positive angle turns
// counter-clockwise seen from +Z (it maps +X onto +Y), which is what
// solid.RotateZMatrix does.

package validate

import (
	"fmt"
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Pose is a rigid transform applied to the moving body before its surface
// points are read against the fixed body: a solid.M44, applied with
// MulPosition, in world coordinates.
type Pose = solid.M44

// ScrewPoses returns the poses of a screw motion about +Z: at each angle
// deg = 0, step, 2*step ... total (degrees, always including total, with a
// short final step if needed) the body is rotated counter-clockwise seen from +Z by deg and
// translated along +Z by lead*deg/360 mm. A positive lead therefore rises as
// it turns counter-clockwise — a right-hand screw being unscrewed upward —
// and a negative lead drives it down. step and total must be positive.
func ScrewPoses(step, total, lead float64) []Pose {
	degs := sweepSteps(0, total, step, "ScrewPoses")
	poses := make([]Pose, len(degs))
	for i, deg := range degs {
		poses[i] = solid.Translate3d(v3.Z(lead * deg / 360)).Mul(solid.RotateZMatrix(deg))
	}
	return poses
}

// LinearPoses returns pure translations of the moving body along the unit
// direction of dir, by from, from+step, ... to millimetres. dir need not be
// normalised; from may exceed to, in which case the travel runs backwards.
func LinearPoses(dir v3.Vec, from, to, step float64) []Pose {
	if dir.Length() == 0 {
		panic("validate.LinearPoses: dir must be a non-zero vector")
	}
	u := dir.Normalize()
	ds := sweepSteps(from, to, step, "LinearPoses")
	poses := make([]Pose, len(ds))
	for i, d := range ds {
		poses[i] = solid.Translate3d(u.MulScalar(d))
	}
	return poses
}

// RotationPoses returns rotations of the moving body about the line through
// the origin along axis, by from, from+step, ... to degrees, counter-clockwise
// seen from the +axis end. Translate the body so its hinge sits on the origin
// before sweeping, or compose these with LinearPoses.
func RotationPoses(axis v3.Vec, from, to, step float64) []Pose {
	if axis.Length() == 0 {
		panic("validate.RotationPoses: axis must be a non-zero vector")
	}
	degs := sweepSteps(from, to, step, "RotationPoses")
	poses := make([]Pose, len(degs))
	for i, deg := range degs {
		poses[i] = solid.Rotate3dMatrix(axis, deg)
	}
	return poses
}

// Compose returns the cartesian product of two pose families: len(a)*len(b)
// poses, each applying an a pose FIRST and then a b pose (b.Mul(a)), in
// a-major order — "for every bore shift, the whole screw travel". An empty
// family is treated as the identity, so Compose(nil, b) returns b.
func Compose(a, b []Pose) []Pose {
	if len(a) == 0 {
		return append([]Pose(nil), b...)
	}
	if len(b) == 0 {
		return append([]Pose(nil), a...)
	}
	out := make([]Pose, 0, len(a)*len(b))
	for _, pa := range a {
		for _, pb := range b {
			out = append(out, pb.Mul(pa))
		}
	}
	return out
}

// SweepResult reports the worst clearance found over a whole motion: Min is
// the fixed body's signed field (mm) at the worst transformed surface point
// of the moving body, so negative means the bodies collide there.
type SweepResult struct {
	Min     float64 // mm, signed: negative = collision
	PoseIdx int     // index into the pose slice where Min occurred
	At      v3.Vec  // the moving-body point, already transformed by that pose
	Poses   int     // poses evaluated
	Points  int     // moving-body surface points per pose
}

// String renders the sweep as one test-log line: the worst clearance, which
// pose and world point produced it, and how much was sampled.
func (r SweepResult) String() string {
	if r.Poses == 0 || r.Points == 0 {
		return fmt.Sprintf("sweep: nothing sampled (%d poses x %d points)", r.Poses, r.Points)
	}
	return fmt.Sprintf("sweep min=%+.3f at pose %d/%d (%.3f, %.3f, %.3f) over %d points",
		r.Min, r.PoseIdx, r.Poses, r.At.X, r.At.Y, r.At.Z, r.Points)
}

// Sweep samples moving's surface at cellsPerMM once, transforms those points
// by each pose, and reads fixed's signed field at every transformed point.
// The result's Min is the worst clearance over the whole motion in mm
// (negative = the bodies interfere at that pose).
func Sweep(moving, fixed *solid.Solid, poses []Pose, cellsPerMM float64) SweepResult {
	pts := Surface(moving, cellsPerMM)
	res := SweepResult{Min: math.Inf(1), Poses: len(poses), Points: len(pts)}
	if len(pts) == 0 || len(poses) == 0 {
		res.Min = 0
		return res
	}
	type best struct {
		min float64
		at  v3.Vec
		idx int
	}
	bests := make([]best, len(poses))
	parallel(len(poses), func(i int) {
		b := best{min: math.Inf(1), idx: i}
		m := poses[i]
		for _, p := range pts {
			q := m.MulPosition(p)
			if v := At(fixed, q); v < b.min {
				b.min, b.at = v, q
			}
		}
		bests[i] = b
	})
	for _, b := range bests {
		if b.min < res.Min {
			res.Min, res.At, res.PoseIdx = b.min, b.at, b.idx
		}
	}
	return res
}

// RequireSweepClear fails the test unless the moving body stays at least
// minClear mm away from the fixed body at every pose (r.Min >= minClear).
func RequireSweepClear(t testing.TB, moving, fixed *solid.Solid, poses []Pose, cellsPerMM, minClear float64) {
	t.Helper()
	r := Sweep(moving, fixed, poses, cellsPerMM)
	t.Logf("%s (want min >= %+.3f)", r, minClear)
	if r.Poses == 0 || r.Points == 0 {
		t.Errorf("sweep: nothing to sample (%d poses, %d surface points at cellsPerMM=%g)", r.Poses, r.Points, cellsPerMM)
		return
	}
	if r.Min < minClear {
		t.Errorf("sweep clearance = %+.3f mm at pose %d (%.3f, %.3f, %.3f), want >= %+.3f mm",
			r.Min, r.PoseIdx, r.At.X, r.At.Y, r.At.Z, minClear)
	}
}

// RequireSweepBlocked fails the test unless the moving body collides with the
// fixed body by at least atLeast mm somewhere in the motion (r.Min <= -atLeast).
// This is the anti-flip guard: it proves a RequireSweepClear over the same
// poses would actually catch an obstruction rather than pass by missing it.
func RequireSweepBlocked(t testing.TB, moving, fixed *solid.Solid, poses []Pose, cellsPerMM, atLeast float64) {
	t.Helper()
	r := Sweep(moving, fixed, poses, cellsPerMM)
	t.Logf("%s (want min <= %+.3f)", r, -atLeast)
	if r.Poses == 0 || r.Points == 0 {
		t.Errorf("sweep: nothing to sample (%d poses, %d surface points at cellsPerMM=%g)", r.Poses, r.Points, cellsPerMM)
		return
	}
	if r.Min > -atLeast {
		t.Errorf("expected the sweep to be blocked by at least %.3f mm, got min = %+.3f mm at pose %d (%.3f, %.3f, %.3f)",
			atLeast, r.Min, r.PoseIdx, r.At.X, r.At.Y, r.At.Z)
	}
}

// --- unexported helpers ---

// sweepSteps returns from, from+step, ... always ending exactly at to (a
// short final step is appended when the span is not a whole multiple of
// step), walking backwards when to < from. step must be positive.
func sweepSteps(from, to, step float64, who string) []float64 {
	if !(step > 0) {
		panic(fmt.Sprintf("validate.%s: step must be > 0, got %g", who, step))
	}
	span := to - from
	dir := 1.0
	if span < 0 {
		dir, span = -1.0, -span
	}
	n := int(math.Floor(span/step + 1e-9))
	out := make([]float64, 0, n+2)
	for i := 0; i <= n; i++ {
		out = append(out, from+dir*float64(i)*step)
	}
	if span-float64(n)*step > 1e-9*math.Max(1, span) {
		out = append(out, to)
	}
	return out
}
