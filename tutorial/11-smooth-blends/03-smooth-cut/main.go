// Smooth blends: SmoothCut fillets the *inside* corner where a tool is
// subtracted from a body. Pair with PolyMax (the smooth-max counterpart of
// PolyMin) for a clean rounded pocket.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	solid.Cylinder(20, 6, 0).
		Top().On(solid.Box(v3.XYZ(20, 20, 14), 1).Top()).
		SmoothCut(solid.PolyMax(2.0)).
		STL("out.stl", 6.0)
}
