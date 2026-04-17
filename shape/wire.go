package shape

import (
	"math"

	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

func WireGroove(r float64, depth float64, angleDeg float64) *Shape {
	alpha := (90.0 - angleDeg) * math.Pi / 180.0

	tx := -r * math.Sin(alpha)
	ty := r * math.Cos(alpha)

	m := math.Tan(alpha)
	c := r / math.Cos(alpha)
	theoreticalExtY := m*depth + c

	maxSafeY := r * 3.5

	// Define the polygon for the flat bottom and sloped/capped roof
	var points []v2.Vec
	points = append(points, v2.Vec{X: 0, Y: 0})      // Anchor in the center
	points = append(points, v2.Vec{X: 0, Y: -r})     // Bottom tangency
	points = append(points, v2.Vec{X: depth, Y: -r}) // Flat floor

	if theoreticalExtY > maxSafeY {
		roofStartX := (maxSafeY - c) / m
		points = append(points, v2.Vec{X: depth, Y: maxSafeY})
		points = append(points, v2.Vec{X: roofStartX, Y: maxSafeY})
		points = append(points, v2.Vec{X: tx, Y: ty})
	} else {
		points = append(points, v2.Vec{X: depth, Y: theoreticalExtY})
		points = append(points, v2.Vec{X: tx, Y: ty})
	}

	poly := Polygon(points)
	wireCircle := Circle(r)

	return poly.Union(wireCircle)
}
