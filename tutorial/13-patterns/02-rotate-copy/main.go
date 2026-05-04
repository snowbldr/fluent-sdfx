// Patterns: RotateCopyZ duplicates a solid n times around the Z axis,
// evenly spaced. The receiver should sit off-axis so the copies don't
// overlap themselves.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	solid.Box(v3.XYZ(2, 4, 6), 0.5).
		TranslateX(10).
		RotateCopyZ(8).
		STL("out.stl", 6.0)
}
