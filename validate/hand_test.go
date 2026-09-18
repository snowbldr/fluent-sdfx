package validate

import (
	"strings"
	"testing"

	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	handCPM   = 4.0
	handPitch = 4.0 // 16 mm of height over 4 turns
)

// handHelix is a 1 mm round wire wound on a 5 mm radius, 4 turns in 16 mm.
func handHelix() *solid.Solid {
	return shape.Circle(1).SweepHelix(5, 4, 16, false)
}

// TestSweepHelixIsRightHandByConstruction settles the library's convention
// before HelixHand is trusted to report it: a point on the helix at angle
// theta must sit at z = pitch*theta/360 (it rises as it turns counter-
// clockwise seen from +Z), so screwing a surface point forward by
// (+360*d/pitch about Z, +d along Z) must leave it on the surface.
func TestSweepHelixIsRightHandByConstruction(t *testing.T) {
	h := handHelix()
	const d = handPitch / 8
	pts := handTrimEnds(Surface(h, handCPM), d)
	if len(pts) == 0 {
		t.Fatal("no surface points")
	}

	right := handScrewErr(h, pts, handPitch, d, RightHand)
	left := handScrewErr(h, pts, handPitch, d, LeftHand)
	if right > 0.02 {
		t.Errorf("rising-with-counter-clockwise screw error = %.4f mm, want ~0: "+
			"SweepHelix is not right-hand by the geometric definition", right)
	}
	if left < 0.2 {
		t.Errorf("mirrored screw error = %.4f mm, want clearly non-zero", left)
	}

	// The same statement, read straight off the geometry: material sits at
	// (r, theta, pitch*theta/360) and air sits at the mirrored height
	// -pitch*theta/360. Only quarter-turn angles are probed: this helix has 4
	// turns 4 mm apart carrying a 1 mm wire, so at 45° or 180° the mirrored
	// height lands on (or within the wire radius of) a NEIGHBOURING turn and
	// says nothing about handedness. At 90° and 270° it is 2 mm clear of
	// every turn.
	for _, deg := range []float64{90, 270} {
		z := handPitch * deg / 360
		if !Inside(h, Cyl(5, deg, z)) {
			t.Errorf("no material at %.0f°, z=%+.3f (the right-hand track)", deg, z)
		}
		if Inside(h, Cyl(5, deg, -z)) {
			t.Errorf("material found at %.0f°, z=%+.3f (the left-hand track)", deg, -z)
		}
	}
}

func TestHelixHand(t *testing.T) {
	h := handHelix()
	tests := []struct {
		name string
		s    *solid.Solid
		want Hand
	}{
		{"SweepHelix default", h, RightHand},
		{"mirrored in XZ", h.MirrorXZ(), LeftHand},
		{"mirrored in YZ", h.Transform(solid.MirrorYZ()), LeftHand},
		{"mirrored twice is right-hand again", h.MirrorXZ().Transform(solid.MirrorYZ()), RightHand},
		{"rotated 37° about Z is still right-hand", h.RotateZ(37), RightHand},
		// A body of revolution is screw-invariant both ways: undecidable.
		{"plain cylinder", solid.Cylinder(16, 6, 0), Unknown},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := HelixHand(tc.s, handPitch, handCPM)
			if r.N == 0 {
				t.Fatal("no usable surface points")
			}
			if r.Hand != tc.want {
				t.Errorf("Hand = %s, want %s (%s)", r.Hand, tc.want, r)
			}
			if r.RightErr < 0 || r.LeftErr < 0 {
				t.Errorf("errors must be absolute: %s", r)
			}
		})
	}
}

// TestHelixHandNegative is the anti-flip case: a left-hand helix must not be
// reported right-hand, and the winning error must be far smaller, not just
// marginally so.
func TestHelixHandNegative(t *testing.T) {
	left := handHelix().MirrorXZ()
	r := HelixHand(left, handPitch, handCPM)
	if r.Hand == RightHand {
		t.Fatalf("mirrored helix reported right-hand: %s", r)
	}
	if r.LeftErr > 0.02 {
		t.Errorf("LeftErr = %.4f mm, want ~0 for a true left-hand helix (%s)", r.LeftErr, r)
	}
	if r.RightErr < 10*r.LeftErr+0.1 {
		t.Errorf("the verdict is not decisive: %s", r)
	}
}

func TestHelixHandWrongPitch(t *testing.T) {
	// Told the wrong pitch, neither screw transform is an invariance, so the
	// call must not confidently pick a hand.
	r := HelixHand(handHelix(), handPitch*1.5, handCPM)
	if r.RightErr < 0.05 && r.LeftErr < 0.05 {
		t.Errorf("both errors small at the wrong pitch: %s", r)
	}
}

func TestHelixHandPanicsOnBadPitch(t *testing.T) {
	for _, pitch := range []float64{0, -4} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("HelixHand(..., %g, ...) did not panic", pitch)
				}
			}()
			HelixHand(handHelix(), pitch, handCPM)
		}()
	}
}

func TestRequireHelixHand(t *testing.T) {
	h := handHelix()
	// Positive: both hands asserted correctly.
	RequireHelixHand(t, h, handPitch, handCPM, RightHand)
	RequireHelixHand(t, h.MirrorXZ(), handPitch, handCPM, LeftHand)

	// Negative, via the pure form (a real Require* failure would fail this test).
	if got := HelixHand(h, handPitch, handCPM).Hand; got == LeftHand {
		t.Errorf("Hand = %s, expected RequireHelixHand(..., LeftHand) to fail here", got)
	}
}

func TestHandString(t *testing.T) {
	tests := []struct {
		h    Hand
		want string
	}{
		{RightHand, "right-hand"},
		{LeftHand, "left-hand"},
		{Unknown, "unknown"},
		{Hand(99), "unknown"},
	}
	for _, tc := range tests {
		if got := tc.h.String(); got != tc.want {
			t.Errorf("Hand(%d).String() = %q, want %q", int(tc.h), got, tc.want)
		}
	}

	r := HandResult{Hand: RightHand, RightErr: 0.001, LeftErr: 1.25, N: 12}
	want := "hand=right-hand rightErr=+0.001 leftErr=+1.250 n=12"
	if got := r.String(); got != want {
		t.Errorf("HandResult.String() = %q, want %q", got, want)
	}
}

func TestHandTrimEnds(t *testing.T) {
	pts := []v3.Vec{v3.Z(0), v3.Z(0.5), v3.Z(5), v3.Z(9.5), v3.Z(10)}
	got := handTrimEnds(pts, 1)
	if len(got) != 1 || got[0].Z != 5 {
		t.Errorf("handTrimEnds = %v, want just z=5", got)
	}
	if n := len(handTrimEnds(nil, 1)); n != 0 {
		t.Errorf("handTrimEnds(nil) = %d points, want 0", n)
	}
	// Trimming more than the whole span leaves nothing, and HelixHand says so.
	if n := len(handTrimEnds(pts, 20)); n != 0 {
		t.Errorf("over-trim left %d points, want 0", n)
	}
}

func TestHelixHandNoUsablePoints(t *testing.T) {
	// A flat disc is thinner than pitch/8, so every point is trimmed away.
	r := HelixHand(solid.Cylinder(0.2, 5, 0), 8, 4)
	if r.N != 0 || r.Hand != Unknown {
		t.Errorf("got %s, want an empty Unknown result", r)
	}
	if !strings.Contains(r.String(), "hand=unknown") {
		t.Errorf("String() = %q", r.String())
	}
}
