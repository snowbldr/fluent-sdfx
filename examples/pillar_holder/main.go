package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

var shrink = 1.0 / 0.999 // PLA ~0.1%

var wallThickness = 2.5
var wallHeight = 15.0
var pillarWidth = 33.0
var pillarRadius = 4.0
var feetWidth = 6.0
var baseThickness = 3.0

func base() *solid.Solid {
	w := pillarWidth + 2.0*(feetWidth+wallThickness)
	h := pillarWidth + 2.0*wallThickness
	r := pillarRadius + wallThickness
	base2d := shape.Rect(v2.XY(w, h), r)
	return solid.Extrude(base2d, baseThickness)
}

func wall(w, r float64) *solid.Solid {
	base := shape.Rect(v2.XY(w, w), r)
	s := solid.Extrude(base, wallHeight)
	ofs := 0.5 * (wallHeight - baseThickness)
	return s.Translate(v3.XYZ(0, 0, ofs))
}

func holder() *solid.Solid {
	b := base()
	outer := wall(pillarWidth+2.0*wallThickness, pillarRadius+wallThickness)
	inner := wall(pillarWidth, pillarRadius)
	return b.Union(outer).Cut(inner)
}

func main() {
	holder().ScaleUniform(shrink).ToSTL("holder.stl", 300)
}
