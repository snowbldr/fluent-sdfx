package shape

import (
	"math"

	v2 "github.com/deadsy/sdfx/vec/v2"
)

func Star(outer, inner float64, points int) *Shape {
	var pts []v2.Vec
	for i := 0; i < points*2; i++ {
		r := outer
		if i%2 != 0 {
			r = inner
		}
		angle := float64(i) * math.Pi / float64(points)
		pts = append(pts, v2.Vec{X: r * math.Cos(angle), Y: r * math.Sin(angle)})
	}
	return Polygon(pts).Offset(outer / 10)
}

func Hexagon(radius float64) *Shape {
	var pts []v2.Vec
	for i := 0; i < 6; i++ {
		angle := float64(i) * (2 * math.Pi / 6)
		pts = append(pts, v2.Vec{X: radius * math.Cos(angle), Y: radius * math.Sin(angle)})
	}
	return Polygon(pts)
}

func Triangle(radius float64) *Shape {
	var pts []v2.Vec
	for i := 0; i < 3; i++ {
		angle := float64(i)*(2*math.Pi/3) + (math.Pi / 2)
		pts = append(pts, v2.Vec{X: radius * math.Cos(angle), Y: radius * math.Sin(angle)})
	}
	return Polygon(pts)
}

func Cross(width, thickness float64) *Shape {
	vBar := Rect(v2.Vec{X: thickness, Y: width}, 0)
	hBar := Rect(v2.Vec{X: width, Y: thickness}, 0)
	return vBar.Union(hBar)
}
