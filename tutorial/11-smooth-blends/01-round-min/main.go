// Smooth blends: SmoothUnion with a RoundMin blend function fillets the
// junction between two solids with a circular radius.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	solid.Box(v3.XYZ(20, 20, 2), 1).
		UnderneathOf(solid.Cylinder(10, 4, 0).Bottom()).
		SmoothAdd(solid.RoundMin(1.25)).
		STL("out.stl", 5.0)
}
