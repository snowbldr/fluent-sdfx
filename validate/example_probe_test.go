package validate_test

import (
	"fmt"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// A probe table says, in one place, what is meant to be material and what is
// meant to be air. Negative readings are inside the part. The magnitudes are
// distance bounds after a boolean — the last reading is 3, not the true 2 mm
// to the top face — so probe tables assert sign, and measure with
// SurfaceZ/SurfaceR.
func ExampleProbes() {
	tube := solid.Cylinder(20, 5, 0).Cut(solid.Cylinder(30, 3, 0))
	for _, r := range validate.Probes(tube, []validate.Probe{
		validate.In("wall at mid height", validate.Cyl(4, 0, 0)),
		validate.Out("bore is clear", validate.Cyl(1, 0, 0)),
		validate.Out("above the tube", v3.XYZ(0, 0, 12)),
	}) {
		fmt.Println(r)
	}
	// Output:
	// [ok] wall at mid height (4.000, 0.000, 0.000) want=inside got=inside sdf=-1.000
	// [ok] bore is clear (1.000, 0.000, 0.000) want=outside got=outside sdf=+2.000
	// [ok] above the tube (0.000, 0.000, 12.000) want=outside got=outside sdf=+3.000
}

// The one-liner a part's test wants: a named table, one report line each and
// a summary, with a t.Errorf per wrong point.
func ExampleRequireProbes() {
	t := &testing.T{} // in real code this is the test's own *testing.T
	plate := solid.Box(v3.XYZ(20, 20, 4), 0)
	validate.RequireProbes(t, plate, []validate.Probe{
		validate.In("centre of the plate", v3.Vec{}),
		validate.Out("clear of the edge", v3.XYZ(11, 0, 0)),
		// 0.5 mm of material must remain under the pocket floor:
		{Name: "pocket floor", P: v3.XYZ(0, 0, -1.5), Inside: true, Margin: 0.5},
	})
	fmt.Println("failed:", t.Failed())
	// Output: failed: false
}

// Probes across an assembly: the map key selects the part, and result names
// come back as "part:point".
func ExampleProbesIn() {
	parts := map[string]*solid.Solid{
		"base": solid.Box(v3.XYZ(20, 20, 4), 0),
		"pin":  solid.Cylinder(10, 2, 0).TranslateZ(7),
	}
	for _, r := range validate.ProbesIn(parts, []validate.NamedProbe{
		{Name: "centre", Solid: "base", P: v3.Vec{}, Inside: true},
		{Name: "mid shank", Solid: "pin", P: v3.XYZ(0, 0, 7), Inside: true},
	}) {
		fmt.Println(r)
	}
	// Output:
	// [ok] base:centre (0.000, 0.000, 0.000) want=inside got=inside sdf=-2.000
	// [ok] pin:mid shank (0.000, 0.000, 7.000) want=inside got=inside sdf=-2.000
}
