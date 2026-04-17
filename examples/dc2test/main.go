package main

import (
	flrender "github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/shape"
)

func main() {
	f := shape.LoadFont("../../files/cmr10.ttf")
	s0 := shape.Text(f, "hi!", 10.0)
	flrender.ToDXFWith(s0, "output.dxf", flrender.NewDualContouring2D(50))
}
