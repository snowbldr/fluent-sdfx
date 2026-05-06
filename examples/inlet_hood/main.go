package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%

var baseSize = v3.XYZ(40, 60, 10)
var portSize = v3.XYZ(30, 50, 10)

func outerBase() *solid.Solid {
	return obj.TruncRectPyramid3D(obj.TruncRectPyramidParms{
		Size:         baseSize,
		BaseAngleDeg: 90.0 - 2.0,
		BaseRadius:   baseSize.X * 0.5,
		RoundRadius:  0,
	})
}

func innerBase() *solid.Solid {
	return obj.TruncRectPyramid3D(obj.TruncRectPyramidParms{
		Size:         portSize,
		BaseAngleDeg: 90.0 - 5.0,
		BaseRadius:   portSize.X * 0.5,
		RoundRadius:  0,
	})
}

func hood() *solid.Solid {
	return outerBase().Cut(innerBase())
}

func main() {
	hood().ScaleUniform(shrink).STL("hood.stl", 3.0)
}
