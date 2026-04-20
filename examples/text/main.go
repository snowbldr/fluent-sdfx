package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
)

func main() {
	f := shape.LoadFont("../../files/cmr10.ttf")

	s0 := shape.Text(f, "SDFX!\nHello,\nWorld!", 10.0)

	// cache the sdf for evaluation speedup
	cached := s0.Cache()

	cached.ToDXF("shape.dxf", 600)
	cached.ToSVG("shape.svg", 600)

	solid.ExtrudeRounded(cached, 1.0, 0.2).STL("shape.stl", 6.0)
}
