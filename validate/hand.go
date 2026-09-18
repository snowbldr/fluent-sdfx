// Handedness of helical geometry, decided by measurement rather than by
// remembering which way a builder wound its profile.
//
// Right-hand is defined geometrically: the helix ADVANCES ALONG +Z WHILE
// TURNING COUNTER-CLOCKWISE SEEN FROM +Z — a standard right-hand thread,
// z = pitch * theta / 360. Left-hand is its mirror image, z = -pitch *
// theta / 360.
//
// The test is screw invariance. A helical solid of axial pitch p is unchanged
// by the screw transform "rotate about +Z by 360*d/p, then translate +d along
// +Z" if it is right-hand, and by the same transform with the rotation negated
// if it is left-hand. So push the solid's own surface points through both
// transforms and read its field: the true hand leaves them on the surface
// (field ~ 0), the wrong hand throws them into air or material.

package validate

import (
	"fmt"
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Hand is the handedness of a helical feature about the +Z axis.
type Hand int

const (
	// Unknown means the two screw-invariance errors were too close to call.
	Unknown Hand = iota
	// RightHand rises along +Z while turning counter-clockwise seen from +Z.
	RightHand
	// LeftHand rises along +Z while turning clockwise seen from +Z.
	LeftHand
)

// String returns "right-hand", "left-hand" or "unknown".
func (h Hand) String() string {
	switch h {
	case RightHand:
		return "right-hand"
	case LeftHand:
		return "left-hand"
	default:
		return "unknown"
	}
}

// HandResult reports which hand won the screw-invariance test and by how
// much: RightErr and LeftErr are the 95th percentiles of |field| (mm) after
// applying the right-hand and left-hand screw transform to N surface points.
// The smaller error is the true hand; a genuine helix gives one near zero and
// one on the order of the feature's own size.
type HandResult struct {
	Hand              Hand
	RightErr, LeftErr float64 // mm, p95 |field| after the right/left screw transform
	N                 int     // surface points used (end bands excluded)
}

// String renders the verdict and both screw-invariance errors as one
// test-log line.
func (r HandResult) String() string {
	return fmt.Sprintf("hand=%s rightErr=%+.3f leftErr=%+.3f n=%d", r.Hand, r.RightErr, r.LeftErr, r.N)
}

// HelixHand decides the hand of a helical solid of the given axial pitch (mm
// of rise per full turn about Z, always positive) by screw invariance: it
// samples the surface at cellsPerMM, advances each point by d = pitch/8 along
// +Z while rotating it +360*d/pitch (right-hand) and -360*d/pitch (left-hand),
// and returns the hand whose p95 |field| is smaller, with both errors. Points
// within d of the solid's top and bottom surface extremes are dropped, since
// the cut ends of any finite helix are not screw-invariant either way. The
// verdict is Unknown when the two errors are within 10% of each other.
func HelixHand(s *solid.Solid, pitch, cellsPerMM float64) HandResult {
	if !(pitch > 0) {
		panic(fmt.Sprintf("validate.HelixHand: pitch must be > 0 mm per turn, got %g", pitch))
	}
	d := pitch / 8
	pts := handTrimEnds(Surface(s, cellsPerMM), d)
	r := HandResult{N: len(pts)}
	if len(pts) == 0 {
		return r
	}
	r.RightErr = handScrewErr(s, pts, pitch, d, RightHand)
	r.LeftErr = handScrewErr(s, pts, pitch, d, LeftHand)
	switch {
	case r.RightErr < 0.9*r.LeftErr:
		r.Hand = RightHand
	case r.LeftErr < 0.9*r.RightErr:
		r.Hand = LeftHand
	default:
		r.Hand = Unknown
	}
	return r
}

// RequireHelixHand fails the test unless HelixHand decides s winds the way
// want does — the guard against a mirrored thread shipping as a right-hand one.
func RequireHelixHand(t testing.TB, s *solid.Solid, pitch, cellsPerMM float64, want Hand) {
	t.Helper()
	r := HelixHand(s, pitch, cellsPerMM)
	t.Logf("%s (want %s, pitch %.3f mm/turn)", r, want, pitch)
	if r.N == 0 {
		t.Errorf("helix hand: no usable surface points at cellsPerMM=%g (is the solid shorter than pitch/4?)", cellsPerMM)
		return
	}
	if r.Hand != want {
		t.Errorf("helix hand = %s, want %s (rightErr=%+.3f leftErr=%+.3f over %d points)",
			r.Hand, want, r.RightErr, r.LeftErr, r.N)
	}
}

// --- unexported helpers ---

// handScrewPose returns the screw transform of the given hand that advances
// +d along Z: rotate about +Z by +360*d/pitch for RightHand, -360*d/pitch for
// LeftHand, then translate +d along Z.
func handScrewPose(pitch, d float64, h Hand) solid.M44 {
	deg := 360 * d / pitch
	if h == LeftHand {
		deg = -deg
	}
	return solid.Translate3d(v3.Z(d)).Mul(solid.RotateZMatrix(deg))
}

// handScrewErr is the p95 of |field of s| after screwing every point of pts
// by one step of the given hand.
func handScrewErr(s *solid.Solid, pts []v3.Vec, pitch, d float64, h Hand) float64 {
	m := handScrewPose(pitch, d, h)
	moved := make([]v3.Vec, len(pts))
	for i, p := range pts {
		moved[i] = m.MulPosition(p)
	}
	vals := evalAll(s, moved)
	for i := range vals {
		vals[i] = math.Abs(vals[i])
	}
	return Percentile(vals, 95)
}

// handTrimEnds drops the points within d of the lowest and highest Z in pts:
// the cut ends of a finite helix break screw invariance for both hands and
// would only dilute the comparison.
func handTrimEnds(pts []v3.Vec, d float64) []v3.Vec {
	if len(pts) == 0 {
		return pts
	}
	lo, hi := math.Inf(1), math.Inf(-1)
	for _, p := range pts {
		lo, hi = math.Min(lo, p.Z), math.Max(hi, p.Z)
	}
	d = math.Abs(d)
	out := make([]v3.Vec, 0, len(pts))
	for _, p := range pts {
		if p.Z >= lo+d && p.Z <= hi-d {
			out = append(out, p)
		}
	}
	return out
}
