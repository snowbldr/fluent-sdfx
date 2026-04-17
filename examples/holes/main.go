package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%

// testHoles returns a panel with various holes for test fitting.
func testHoles() *solid.Solid {
	const xInc = 15
	const yInc = 15
	const rInc = 0.1
	const nX = 5
	const nY = 8

	xOfs := 0.0
	yOfs := 0.0
	r := 1.5

	circles := make([]*shape.Shape, 0, nX*nY)
	for j := 0; j < nY; j++ {
		for k := 0; k < nX; k++ {
			circles = append(circles, shape.Circle(r).Translate(v2.XY(xOfs, yOfs)))
			r += rInc
			xOfs += xInc
		}
		xOfs = 0.0
		yOfs += yInc
	}

	holes := circles[0].Union(circles[1:]...)
	xOfs = -float64(nX-1) * xInc * 0.5
	yOfs = -float64(nY-1) * yInc * 0.5
	holes = holes.Translate(v2.XY(xOfs, yOfs))

	// make a panel
	k := obj.PanelParms{
		Size:         v2.XY((nX+1)*xInc, (nY+1)*yInc),
		CornerRadius: xInc * 0.2,
	}
	panel := obj.Panel2D(k).Cut(holes)

	return solid.Extrude(panel, 3)
}

func main() {
	testHoles().ScaleUniform(shrink).ToSTL("test_holes.stl", 300)
}
