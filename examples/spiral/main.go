package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
)

func main() {
	s := shape.ArcSpiral(1.0, 20.0, 45, 8*360, 1.0)
	s.ToDXF("spiral.dxf", 400)
}
