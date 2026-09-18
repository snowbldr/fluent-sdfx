package validate_test

import (
	"fmt"
	"os"

	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// FitCheck writes the two STLs you actually open — the whole assembly and a
// cutaway through it — and hands back every pairwise clearance, so one call
// covers both "does it look right" and "does it fit".
func ExampleFitCheck() {
	block := solid.Box(v3.XYZ(20, 20, 10), 0).Cut(solid.Cylinder(12, 5, 0))
	pin := solid.Cylinder(8, 4.8, 0)

	dir, _ := os.MkdirTemp("", "fit")
	defer os.RemoveAll(dir)

	gaps, err := validate.FitCheck(dir, "pinfit", []validate.Part{
		{Name: "block", Solid: block},
		{Name: "pin", Solid: pin},
	}, v3.Vec{}, v3.Y(1), 8) // cutaway: the plane through the origin, keeping +Y
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("clear:", gaps["pin>block"].Min > 0)
	// Output: clear: true
}

// Manifest is the summary line a build log (or a model) reads instead of
// opening the STL: size, volume and how many separate bodies each part is.
func ExampleManifest() {
	fmt.Print(validate.Manifest([]validate.Part{
		{Name: "cube", Solid: solid.Box(v3.XYZ(10, 10, 10), 0)},
	}, 4))
	// Output:
	// | name | min | max | size | volume mm³ | bodies |
	// | --- | --- | --- | --- | --- | --- |
	// | cube | (-5.000, -5.000, -5.000) | (5.000, 5.000, 5.000) | 10.000 × 10.000 × 10.000 | 1000.0 | 1 |
}
