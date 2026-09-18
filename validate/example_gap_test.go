package validate_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// gapPin is a Ø8 x 20 pin; gapSocket is a Ø20 x 20 block bored Ø8.5, so the
// pin has 0.25 mm of clearance all round.
func gapPin() *solid.Solid    { return solid.Cylinder(20, 4, 0) }
func gapSocket() *solid.Solid { return solid.Cylinder(20, 10, 0).Cut(solid.Cylinder(22, 4.25, 0)) }

// Clearance reads the second solid's field at every surface point of the
// first: how close does the pin come to the wall of its socket?
func ExampleClearance() {
	g := validate.Clearance(gapPin(), gapSocket(), 6.0)
	fmt.Println("pin never touches the socket:", g.Min > 0)
	fmt.Println("at least 0.2 mm of air everywhere:", g.Min >= 0.2)
	// Output:
	// pin never touches the socket: true
	// at least 0.2 mm of air everywhere: true
}

// ClearanceIn asks the same question about one region only — here just the
// nose of the pin, ignoring the rest of the part.
func ExampleClearanceIn() {
	nose := v3.Box{Min: v3.XYZ(-10, -10, 8), Max: v3.XYZ(10, 10, 10)}
	g := validate.ClearanceIn(gapPin(), gapSocket(), nose, 6.0)
	full := validate.Clearance(gapPin(), gapSocket(), 6.0)
	fmt.Println("nose sampled fewer points than the whole pin:", g.N > 0 && g.N < full.N)
	fmt.Println("nose still has 0.2 mm of air:", g.Min >= 0.2)
	// Output:
	// nose sampled fewer points than the whole pin: true
	// nose still has 0.2 mm of air: true
}

// RequireClearance is the one-liner that replaces a hand-written loop over
// probe points: it logs the whole gap report and fails on the worst reading.
func ExampleRequireClearance() {
	// In real test code:
	//
	//	func TestPinFitsSocket(t *testing.T) {
	//		validate.RequireClearance(t, pin(), socket(), 6.0, 0.2)
	//		validate.RequireNoContact(t, pin(), socket(), 6.0)
	//	}
	t := &testing.T{} // stand-in for godoc
	validate.RequireClearance(t, gapPin(), gapSocket(), 6.0, 0.2)
	validate.RequireNoContact(t, gapPin(), gapSocket(), 6.0)
}

// RequireInterference is the anti-flip partner: it proves the clearance check
// would really have caught a collision. An oversized pin must jam.
func ExampleRequireInterference() {
	// In real test code:
	//
	//	func TestOversizePinJams(t *testing.T) {
	//		validate.RequireInterference(t, solid.Cylinder(20, 5, 0), socket(), 6.0, 0.5)
	//	}
	t := &testing.T{} // stand-in for godoc
	validate.RequireInterference(t, solid.Cylinder(20, 5, 0), gapSocket(), 6.0, 0.5)
}

// Congruent is the refactor guard: the rewritten builder must land on the
// same surface as the old one, checked in both directions.
func ExampleCongruent() {
	old := solid.Box(v3.XYZ(10, 10, 10), 0)
	rebuilt := solid.Box(v3.XYZ(10, 10, 10), 0).Intersect(solid.Box(v3.XYZ(30, 30, 30), 0))
	grown := solid.Box(v3.XYZ(11, 10, 10), 0)

	_, _, sameOK := validate.Congruent(old, rebuilt, 5.0, 0.05)
	_, _, grownOK := validate.Congruent(old, grown, 5.0, 0.05)
	fmt.Println("rebuilt:", sameOK, " grown 1mm:", grownOK)
	// Output:
	// rebuilt: true  grown 1mm: false
}

// DeviationSTL checks an exported mesh against the solid it came from — the
// question to ask after decimating a part for a slicer.
//
// There is no Output: comment because Solid.STL prints its own render
// progress; the example asserts by panicking instead.
func ExampleDeviationSTL() {
	dir, err := os.MkdirTemp("", "validate-example")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	part := solid.Box(v3.XYZ(10, 10, 10), 0)
	path := filepath.Join(dir, "part.stl")
	part.STL(path, 5.0, 0.5) // render, then decimate

	g, err := validate.DeviationSTL(path, part)
	if err != nil {
		panic(err)
	}
	const budget = 0.1 // mm; well under a print layer height
	if g.Max > budget {
		panic("decimation moved the surface: " + g.String())
	}
}
