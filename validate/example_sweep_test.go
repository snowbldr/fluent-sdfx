package validate_test

import (
	"fmt"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// sweepHandle is a 2 mm peg on a 10 mm arm: it sweeps a circle of radius 10
// about Z. sweepGuard is a wall 20 mm out in +Y.
func sweepHandle() *solid.Solid {
	return solid.Box(v3.XYZ(2, 2, 2), 0).TranslateX(10)
}
func sweepGuard(atY float64) *solid.Solid {
	return solid.Box(v3.XYZ(40, 4, 40), 0).TranslateY(atY)
}

// Sweep transforms the moving body's surface through every pose and reads the
// fixed body's field: one call answers "does it clear all the way round?".
func ExampleSweep() {
	r := validate.Sweep(sweepHandle(), sweepGuard(20), validate.RotationPoses(v3.Z(1), 0, 345, 15), 3.0)
	fmt.Println("clears the whole turn:", r.Min > 0)
	fmt.Println("worst case is the 90° pose:", r.PoseIdx == 6)
	// Output:
	// clears the whole turn: true
	// worst case is the 90° pose: true
}

// ScrewPoses is the motion of a thread backing out: counter-clockwise seen
// from +Z, rising lead*deg/360 as it turns.
func ExampleScrewPoses() {
	poses := validate.ScrewPoses(90, 360, 8) // 90° steps, one turn, 8 mm lead
	p := poses[1].MulPosition(v3.XYZ(10, 0, 0))
	fmt.Printf("after 90°, (10, 0, 0) is at (%.1f, %.1f, %.1f)\n", p.X, p.Y, p.Z)
	// Output:
	// after 90°, (10, 0, 0) is at (0.0, 10.0, 2.0)
}

// Compose covers a whole tolerance band at once: every mounting offset
// combined with every step of the motion.
func ExampleCompose() {
	turn := validate.RotationPoses(v3.Z(1), 0, 345, 15)
	slop := validate.LinearPoses(v3.Y(1), -0.4, 0.4, 0.4)
	poses := validate.Compose(turn, slop) // turn first, then the mounting offset
	fmt.Printf("%d poses = %d angles x %d offsets\n", len(poses), len(turn), len(slop))
	// Output:
	// 72 poses = 24 angles x 3 offsets
}

// RequireSweepClear guards the motion in one line; RequireSweepBlocked is the
// anti-flip partner that proves the check can actually fail.
func ExampleRequireSweepClear() {
	// In real test code:
	//
	//	func TestHandleSwingsFree(t *testing.T) {
	//		poses := validate.RotationPoses(v3.Z(1), 0, 345, 15)
	//		validate.RequireSweepClear(t, handle(), guard(), poses, 3.0, 1.0)
	//		validate.RequireSweepBlocked(t, handle(), stop(), poses, 3.0, 0.5)
	//	}
	t := &testing.T{} // stand-in for godoc
	poses := validate.RotationPoses(v3.Z(1), 0, 345, 15)
	validate.RequireSweepClear(t, sweepHandle(), sweepGuard(20), poses, 3.0, 1.0)
	validate.RequireSweepBlocked(t, sweepHandle(), sweepGuard(10), poses, 3.0, 0.5)
}
