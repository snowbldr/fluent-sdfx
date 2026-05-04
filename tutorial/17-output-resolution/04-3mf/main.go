// Output: ThreeMF (.3mf) is the modern alternative to STL — a zipped
// XML format with explicit units, colour, and metadata. Same density
// parameter as STL.
//
// We also write an STL alongside so the screenshot pipeline picks one up.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	part := solid.Box(v3.XYZ(20, 20, 20), 4).Cut(solid.Sphere(11))
	part.ThreeMF("out.3mf", 5.0)
	part.STL("out.stl", 5.0)
}
