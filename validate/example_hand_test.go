package validate_test

import (
	"fmt"
	"testing"

	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
)

// handThread is a 1 mm wire wound on a 5 mm radius, 4 turns over 16 mm: a
// 4 mm pitch. shape.SweepHelix winds right-hand — it rises along +Z while
// turning counter-clockwise seen from +Z.
func handThread() *solid.Solid { return shape.Circle(1).SweepHelix(5, 4, 16, false) }

// HelixHand decides handedness by measurement, so a mirrored part cannot slip
// through a review as "obviously still right-hand".
func ExampleHelixHand() {
	right := validate.HelixHand(handThread(), 4.0, 4.0)
	left := validate.HelixHand(handThread().MirrorXZ(), 4.0, 4.0)
	fmt.Println("as built: ", right.Hand)
	fmt.Println("mirrored: ", left.Hand)
	// Output:
	// as built:  right-hand
	// mirrored:  left-hand
}

// RequireHelixHand is the one-liner guard for a threaded part: it logs both
// screw-invariance errors and fails if the thread winds the wrong way.
func ExampleRequireHelixHand() {
	// In real test code:
	//
	//	func TestCapThreadIsRightHand(t *testing.T) {
	//		validate.RequireHelixHand(t, capThread(), 4.0, 6.0, validate.RightHand)
	//	}
	t := &testing.T{} // stand-in for godoc
	validate.RequireHelixHand(t, handThread(), 4.0, 4.0, validate.RightHand)
}
