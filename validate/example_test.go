package validate_test

import (
	"fmt"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Compute validation stats for a 10mm cube. Volume should be ~1000 mm³;
// marching-cubes gives an approximate surface but the mesh is always
// closed, so Watertight is reliable.
func Example() {
	cube := solid.Box(v3.XYZ(10, 10, 10), 0)
	st := validate.Of(cube, 8.0)
	fmt.Printf("watertight=%v vol≈%.0f", st.Watertight, st.Volume)
	// Output: watertight=true vol≈1000
}

// RequireWatertight is a t.Helper-style assertion: drop it in a unit test to
// guard against accidental holes from a refactor that breaks a boolean.
func ExampleRequireWatertight() {
	// In real test code:
	//
	//	func TestBracketSealed(t *testing.T) {
	//		validate.RequireWatertight(t, buildBracket(), 5.0)
	//	}
	t := &testing.T{} // stand-in for godoc
	validate.RequireWatertight(t, solid.Box(v3.XYZ(5, 5, 5), 0), 8.0)
}

// A flat-topped cube has its entire top face overhanging horizontally; for a
// 10mm cube that's ~100 mm² of unsupported ceiling.
func ExampleOverhangArea() {
	cube := solid.Box(v3.XYZ(10, 10, 10), 0)
	st := validate.Of(cube, 8.0)
	fmt.Printf("overhang≈%.0f mm²", st.OverhangArea)
	// Output: overhang≈100 mm²
}
