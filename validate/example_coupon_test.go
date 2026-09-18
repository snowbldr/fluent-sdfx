package validate_test

import (
	"fmt"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Coupon cuts a printable test piece out of a big part and drops it onto
// the plate, so a 2-hour print becomes a 6-minute one when all you want to
// know is whether the clearance is right.
func ExampleCoupon() {
	bar := solid.Box(v3.XYZ(60, 10, 10), 0).Translate(v3.Z(20)) // floating at z 15..25
	keep := solid.Box(v3.XYZ(12, 20, 20), 0).Translate(v3.Z(20))

	c := validate.Coupon(bar, keep, 6)
	bb := c.Bounds()
	fmt.Printf("%.0f × %.0f × %.0f at z=%.0f", bb.Size().X, bb.Size().Y, bb.Size().Z, bb.Min.Z)
	// Output: 12 × 10 × 10 at z=0
}

// Dimples labels a ladder of variants by touch: print the 0.1 mm-clearance
// coupon with one dot, the 0.2 with two, and they are still tellable apart
// after they come off the plate. Text at this size is not.
func ExampleDimples() {
	coupon := solid.Box(v3.XYZ(20, 10, 10), 0).Translate(v3.Z(5))
	// Three dots in the top face, 2 mm across, 1 mm deep, 4 mm apart.
	marked := validate.Dimples(coupon, 3, v3.XYZ(0, 0, 10), v3.Z(1), v3.X(1), 2, 1, 4)

	fmt.Println(validate.Inside(marked, v3.XYZ(4, 0, 9.5)), validate.Inside(marked, v3.XYZ(2, 0, 9.5)))
	// Output: false true
}

// Components counts what the slicer would see: a solid that renders as two
// bodies prints as two loose parts.
func ExampleComponents() {
	pillars := solid.Box(v3.XYZ(6, 6, 10), 0).Translate(v3.XYZ(-10, 0, 5)).
		Union(solid.Box(v3.XYZ(6, 6, 10), 0).Translate(v3.XYZ(10, 0, 5)))
	fmt.Println(validate.Components(pillars, 6))
	// Output: 2
}

// RequireOneBody guards against a boolean that quietly severed a part.
func ExampleRequireOneBody() {
	// In real test code:
	//
	//	func TestBracketIsOnePiece(t *testing.T) {
	//		validate.RequireOneBody(t, bracket, 8)
	//	}
	t := &testing.T{} // stand-in for godoc
	validate.RequireOneBody(t, solid.Box(v3.XYZ(10, 10, 10), 0), 6)
	fmt.Println("failed:", t.Failed())
	// Output: failed: false
}
