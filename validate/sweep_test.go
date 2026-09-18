package validate

import (
	"math"
	"strings"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const sweepCPM = 3.0

// sweepPeg is a 2 mm cube on a 10 mm arm: it sits at (10, 0, 0) and sweeps a
// circle of radius 10 about Z.
func sweepPeg() *solid.Solid {
	return solid.Box(v3.XYZ(2, 2, 2), 0).TranslateX(10)
}

// sweepBlock is a 6 mm cube centred at (30, 0, 0): the peg's head-on obstacle.
func sweepBlock() *solid.Solid {
	return solid.Box(v3.XYZ(6, 6, 6), 0).TranslateX(30)
}

// sweepWall is a tall 4 mm-thick wall spanning x and z, centred at y = atY.
func sweepWall(atY float64) *solid.Solid {
	return solid.Box(v3.XYZ(40, 4, 40), 0).TranslateY(atY)
}

func sweepNear(t *testing.T, name string, got, want, tol float64) {
	t.Helper()
	if math.Abs(got-want) > tol {
		t.Errorf("%s = %.4f, want %.4f ± %.4f", name, got, want, tol)
	}
}

// TestScrewPosesSignConvention pins the convention down on a known point:
// +X must swing to +Y at 90° (counter-clockwise seen from +Z) while a
// positive lead lifts it by lead*deg/360.
func TestScrewPosesSignConvention(t *testing.T) {
	const lead = 8.0
	poses := ScrewPoses(90, 360, lead)
	if len(poses) != 5 {
		t.Fatalf("ScrewPoses(90, 360, ...) gave %d poses, want 5 (0, 90, 180, 270, 360)", len(poses))
	}
	p := v3.XYZ(10, 0, 0)
	tests := []struct {
		deg     float64
		idx     int
		want    v3.Vec
		comment string
	}{
		{0, 0, v3.XYZ(10, 0, 0), "identity"},
		{90, 1, v3.XYZ(0, 10, lead*0.25), "+X swings to +Y, rising a quarter lead"},
		{180, 2, v3.XYZ(-10, 0, lead*0.5), "half turn, half a lead up"},
		{270, 3, v3.XYZ(0, -10, lead*0.75), "three quarters"},
		{360, 4, v3.XYZ(10, 0, lead), "full turn, one whole lead up"},
	}
	for _, tc := range tests {
		got := poses[tc.idx].MulPosition(p)
		if !got.Equals(tc.want, 1e-9) {
			t.Errorf("pose at %.0f° (%s): (10,0,0) -> (%.4f, %.4f, %.4f), want (%.4f, %.4f, %.4f)",
				tc.deg, tc.comment, got.X, got.Y, got.Z, tc.want.X, tc.want.Y, tc.want.Z)
		}
	}

	// A negative lead drives the body down as it turns counter-clockwise.
	down := ScrewPoses(90, 90, -lead)[1].MulPosition(p)
	sweepNear(t, "negative lead z at 90°", down.Z, -lead*0.25, 1e-9)
}

// TestScrewPosesOnRealGeometry settles the same convention through a solid:
// the peg's material must actually be found at +Y after a quarter turn.
func TestScrewPosesOnRealGeometry(t *testing.T) {
	const lead = 8.0
	peg := sweepPeg()
	quarter := ScrewPoses(90, 90, lead)[1]

	pts := Surface(peg, sweepCPM)
	var minY, maxY, minZ, maxZ = math.Inf(1), math.Inf(-1), math.Inf(1), math.Inf(-1)
	for _, p := range pts {
		q := quarter.MulPosition(p)
		minY, maxY = math.Min(minY, q.Y), math.Max(maxY, q.Y)
		minZ, maxZ = math.Min(minZ, q.Z), math.Max(maxZ, q.Z)
	}
	// The 2 mm cube centred at (10,0,0) lands centred at (0, 10, lead/4).
	sweepNear(t, "min Y after a quarter turn", minY, 9, 0.2)
	sweepNear(t, "max Y after a quarter turn", maxY, 11, 0.2)
	sweepNear(t, "min Z after a quarter turn", minZ, lead/4-1, 0.2)
	sweepNear(t, "max Z after a quarter turn", maxZ, lead/4+1, 0.2)
}

func TestSweepSteps(t *testing.T) {
	tests := []struct {
		name           string
		from, to, step float64
		want           []float64
	}{
		{"exact multiple", 0, 90, 30, []float64{0, 30, 60, 90}},
		{"short final step", 0, 100, 30, []float64{0, 30, 60, 90, 100}},
		{"single pose", 0, 0, 10, []float64{0}},
		{"offset range", 2, 5, 1, []float64{2, 3, 4, 5}},
		{"backwards", 0, -3, 1, []float64{0, -1, -2, -3}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := sweepSteps(tc.from, tc.to, tc.step, "test")
			if len(got) != len(tc.want) {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
			for i := range got {
				sweepNear(t, "step", got[i], tc.want[i], 1e-9)
			}
		})
	}
}

func TestPoseGeneratorsPanicOnBadInput(t *testing.T) {
	tests := []struct {
		name string
		fn   func()
	}{
		{"zero step", func() { ScrewPoses(0, 90, 1) }},
		{"negative step", func() { LinearPoses(v3.X(1), 0, 10, -1) }},
		{"zero dir", func() { LinearPoses(v3.Vec{}, 0, 10, 1) }},
		{"zero axis", func() { RotationPoses(v3.Vec{}, 0, 90, 10) }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Error("want a panic, got none")
				}
			}()
			tc.fn()
		})
	}
}

func TestLinearAndRotationPoses(t *testing.T) {
	// LinearPoses normalises dir: 5 mm along (2,0,0) is 5 mm, not 10.
	lp := LinearPoses(v3.XYZ(2, 0, 0), 0, 5, 5)
	if len(lp) != 2 {
		t.Fatalf("got %d poses, want 2", len(lp))
	}
	got := lp[1].MulPosition(v3.Vec{})
	if !got.Equals(v3.XYZ(5, 0, 0), 1e-9) {
		t.Errorf("origin translated to (%.3f, %.3f, %.3f), want (5, 0, 0)", got.X, got.Y, got.Z)
	}

	// RotationPoses about +X is counter-clockwise seen from +X: +Y -> +Z.
	rp := RotationPoses(v3.X(1), 0, 90, 90)
	gotR := rp[1].MulPosition(v3.Y(1))
	if !gotR.Equals(v3.Z(1), 1e-9) {
		t.Errorf("+Y rotated 90° about +X -> (%.3f, %.3f, %.3f), want (0, 0, 1)", gotR.X, gotR.Y, gotR.Z)
	}
}

func TestCompose(t *testing.T) {
	lift := LinearPoses(v3.Z(1), 1, 3, 1)     // 3 poses: +1, +2, +3 in Z
	shift := LinearPoses(v3.X(1), 10, 20, 10) // 2 poses: +10, +20 in X
	c := Compose(lift, shift)
	if len(c) != 6 {
		t.Fatalf("Compose gave %d poses, want 6", len(c))
	}
	// a-major order: c[0] = lift[0] then shift[0].
	got := c[0].MulPosition(v3.Vec{})
	if !got.Equals(v3.XYZ(10, 0, 1), 1e-9) {
		t.Errorf("c[0] moved the origin to (%.3f, %.3f, %.3f), want (10, 0, 1)", got.X, got.Y, got.Z)
	}
	got = c[5].MulPosition(v3.Vec{})
	if !got.Equals(v3.XYZ(20, 0, 3), 1e-9) {
		t.Errorf("c[5] moved the origin to (%.3f, %.3f, %.3f), want (20, 0, 3)", got.X, got.Y, got.Z)
	}

	// Composition order matters: rotate-then-translate is not the reverse.
	rot := RotationPoses(v3.Z(1), 90, 90, 90)
	tr := LinearPoses(v3.X(1), 10, 10, 10)
	rotFirst := Compose(rot, tr)[0].MulPosition(v3.XYZ(1, 0, 0)) // -> (0,1,0) -> (10,1,0)
	trFirst := Compose(tr, rot)[0].MulPosition(v3.XYZ(1, 0, 0))  // -> (11,0,0) -> (0,11,0)
	if !rotFirst.Equals(v3.XYZ(10, 1, 0), 1e-9) {
		t.Errorf("rotate then translate -> (%.3f, %.3f, %.3f), want (10, 1, 0)", rotFirst.X, rotFirst.Y, rotFirst.Z)
	}
	if !trFirst.Equals(v3.XYZ(0, 11, 0), 1e-9) {
		t.Errorf("translate then rotate -> (%.3f, %.3f, %.3f), want (0, 11, 0)", trFirst.X, trFirst.Y, trFirst.Z)
	}

	// An empty family is the identity.
	if n := len(Compose(nil, shift)); n != len(shift) {
		t.Errorf("Compose(nil, shift) gave %d poses, want %d", n, len(shift))
	}
	if n := len(Compose(shift, nil)); n != len(shift) {
		t.Errorf("Compose(shift, nil) gave %d poses, want %d", n, len(shift))
	}
}

func TestSweepClearanceAndCollision(t *testing.T) {
	peg := sweepPeg()
	full := RotationPoses(v3.Z(1), 0, 345, 15)
	tests := []struct {
		name    string
		fixed   *solid.Solid
		poses   []Pose
		wantMin float64 // worst clearance over the motion
	}{
		// Peg reaches y = 11 at 90°; the wall's near face is at y = 18.
		{"full turn clears a distant wall", sweepWall(20), full, 7},
		// Wall near face at y = 8; the peg's centre passes y = 10 -> 2 mm in.
		{"full turn hits a close wall", sweepWall(10), full, -2},
		// Stopping at 45° puts the peg's far corner at y = 11sin45 + cos45 =
		// 8.485; the wall's near face is at y = 12.
		{"quarter turn stops short", sweepWall(14), RotationPoses(v3.Z(1), 0, 45, 15), 3.515},
		// Straight-line travel toward a 6 mm block whose near face is at x = 27.
		{"linear travel clears", sweepBlock(), LinearPoses(v3.X(1), 0, 8, 2), 8},
		{"linear travel rams", sweepBlock(), LinearPoses(v3.X(1), 0, 20, 2), -2},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := Sweep(peg, tc.fixed, tc.poses, sweepCPM)
			if r.Poses != len(tc.poses) || r.Points == 0 {
				t.Fatalf("sampled %d poses x %d points", r.Poses, r.Points)
			}
			if math.Abs(r.Min-tc.wantMin) > 0.3 {
				t.Errorf("Min = %+.3f, want %+.3f (%s)", r.Min, tc.wantMin, r)
			}
			if r.PoseIdx < 0 || r.PoseIdx >= len(tc.poses) {
				t.Errorf("PoseIdx = %d, out of range", r.PoseIdx)
			}
		})
	}
}

func TestSweepWorstPoseIsTheClosestApproach(t *testing.T) {
	// The peg passes closest to a +Y wall at 90°, pose index 6 of 15° steps.
	r := Sweep(sweepPeg(), sweepWall(20), RotationPoses(v3.Z(1), 0, 345, 15), sweepCPM)
	if r.PoseIdx != 6 {
		t.Errorf("PoseIdx = %d (%.0f°), want 6 (90°): %s", r.PoseIdx, float64(r.PoseIdx)*15, r)
	}
	if math.Abs(r.At.Y-11) > 0.3 {
		t.Errorf("At = (%.3f, %.3f, %.3f), want the peg's +Y face at y≈11", r.At.X, r.At.Y, r.At.Z)
	}
}

func TestSweepResultString(t *testing.T) {
	r := SweepResult{Min: -0.25, PoseIdx: 3, At: v3.XYZ(1, 2, 3), Poses: 10, Points: 400}
	want := "sweep min=-0.250 at pose 3/10 (1.000, 2.000, 3.000) over 400 points"
	if got := r.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
	empty := SweepResult{}
	if !strings.Contains(empty.String(), "nothing sampled") {
		t.Errorf("empty String() = %q, want the nothing-sampled message", empty.String())
	}
}

func TestSweepWithNoPoses(t *testing.T) {
	r := Sweep(sweepPeg(), sweepWall(20), nil, sweepCPM)
	if r.Poses != 0 || r.Min != 0 {
		t.Errorf("empty sweep = %+v, want zero poses and Min 0", r)
	}
}

func TestRequireSweepClearAndBlocked(t *testing.T) {
	peg := sweepPeg()
	full := RotationPoses(v3.Z(1), 0, 345, 15)

	// Positive: 7 mm of clearance really is at least 5 mm.
	RequireSweepClear(t, peg, sweepWall(20), full, sweepCPM, 5)
	// Positive: the close wall really does block the turn by at least 1 mm.
	RequireSweepBlocked(t, peg, sweepWall(10), full, sweepCPM, 1)

	// Negative, via the pure form (a real Require* failure would fail this test).
	if r := Sweep(peg, sweepWall(20), full, sweepCPM); r.Min >= 9 {
		t.Errorf("min = %+.3f, expected RequireSweepClear(..., 9) to fail", r.Min)
	}
	if r := Sweep(peg, sweepWall(20), full, sweepCPM); r.Min <= 0 {
		t.Errorf("min = %+.3f, expected RequireSweepBlocked to fail on a clear sweep", r.Min)
	}
}

// TestSweepComposedFamily is the "shift × screw" case from the spec: the peg
// turns, and on top of that its mounting may sit anywhere in a ±0.4 mm
// tolerance band in world Y. The sweep must cover every combination.
func TestSweepComposedFamily(t *testing.T) {
	turn := RotationPoses(v3.Z(1), 0, 345, 15)
	shift := LinearPoses(v3.Y(1), -0.4, 0.4, 0.4)
	poses := Compose(turn, shift) // turn first, then the mounting offset
	if len(poses) != len(turn)*len(shift) {
		t.Fatalf("composed to %d poses, want %d", len(poses), len(turn)*len(shift))
	}
	r := Sweep(sweepPeg(), sweepWall(20), poses, sweepCPM)
	// Worst case: at 90° the peg reaches y = 11, plus 0.4 of offset, against
	// the wall face at y = 18.
	if math.Abs(r.Min-6.6) > 0.3 {
		t.Errorf("Min = %+.3f, want +6.600 (%s)", r.Min, r)
	}
}
