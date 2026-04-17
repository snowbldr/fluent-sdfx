package main

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/units"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

func holes() *shape.Shape {
	l := 1.0
	circleRadius := l * 0.2
	n := 6
	steps := 11

	c := shape.Circle(circleRadius)
	s := c

	dBase := 0.0

	for i := 1; i <= steps; i++ {
		k := i * n
		r := float64(i) * l
		dTheta := units.Tau / float64(k)
		c0 := c.Translate(v2.XY(0, r))

		for j := 0; j < k; j++ {
			c1 := c0.Rotate((dBase + float64(j)*dTheta) * 180 / units.Pi)
			s = s.Union(c1)
		}
		dBase += 0.5 * dTheta
	}

	return s
}

func main() {
	holes().ToDXF("holes.dxf", 300)
}
