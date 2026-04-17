package main

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
)

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%

const steps = 20

func sprue(r, l, k float64) *solid.Solid {
	a0 := math.Pi * r * r
	h0 := math.Pow(k/a0, 2)
	dh := l / steps

	p := shape.NewPoly()
	p.Add(0, 0)
	for h := 0.0; h <= l; h += dh {
		a := k / math.Sqrt(h+h0)
		rr := math.Sqrt(a / math.Pi)
		p.Add(rr, h)
	}
	p.Add(0, l)

	return solid.Revolve(p.Build())
}

func main() {
	sprue(20, 100, 3000).ScaleUniform(shrink).ToSTL("sprue.stl", 300)
}
