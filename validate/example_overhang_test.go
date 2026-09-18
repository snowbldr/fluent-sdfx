package validate_test

import (
	"fmt"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// exampleT is a stem carrying a wider cross-bar: 5 mm of bar hangs
// unsupported on each side, the classic part that needs supports.
func exampleT() *solid.Solid {
	stem := solid.Box(v3.XYZ(6, 6, 10), 0).Translate(v3.Z(5))
	bar := solid.Box(v3.XYZ(16, 6, 4), 0).Translate(v3.Z(12))
	return stem.Union(bar)
}

// Overhangs audits the solid in the frame it is in: minimum Z is the build
// plate. The cross-bar underside is a flat 90° ceiling with air beneath it,
// so it is reported; the part's own bottom face is on the plate and is not.
func ExampleOverhangs() {
	r := validate.Overhangs(exampleT(), 8, validate.OverhangOptions{}) // default 45°
	fmt.Printf("unsupported=%v worst=%.0f°", r.Area > 1, r.WorstDeg)
	// Output: unsupported=true worst=90°
}

// A per-height threshold exempts a band whose angle is there by
// construction — a helical flank, or here the cross-bar underside — without
// blinding the audit everywhere else.
func ExampleOverhangOptions_threshold() {
	r := validate.Overhangs(exampleT(), 8, validate.OverhangOptions{
		Threshold: func(z float64) float64 {
			if z > 9 {
				return 90 // the bar underside is allowed to be flat
			}
			return 45
		},
	})
	fmt.Printf("faces=%d area=%.1f mm²", r.Faces, r.Area)
	// Output: faces=0 area=0.0 mm²
}

// RequireOverhangArea is the one-liner form: it logs the report and fails
// the test when the unsupported area exceeds the budget in mm².
func ExampleRequireOverhangArea() {
	// In real test code:
	//
	//	func TestBracketPrints(t *testing.T) {
	//		validate.RequireOverhangArea(t, bracket, 8, validate.OverhangOptions{MaxDeg: 50}, 1.0)
	//	}
	t := &testing.T{} // stand-in for godoc
	box := solid.Box(v3.XYZ(10, 10, 10), 0).Translate(v3.Z(5))
	validate.RequireOverhangArea(t, box, 6, validate.OverhangOptions{}, 0.5)
	fmt.Println("failed:", t.Failed())
	// Output: failed: false
}
