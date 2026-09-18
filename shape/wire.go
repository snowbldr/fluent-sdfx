package shape

import (
	"math"

	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

// WireGroove returns the cross-section of a channel that captures a wire of
// radius r, centered on the origin and opening toward +X.
// The wire seat is a full circle of radius r at the origin; the channel runs
// from there out to x = depth with a flat floor at y = -r (tangent to the
// bottom of the wire) and a roof that is tangent to the circle and rises
// toward +X at (90 - angleDeg) degrees from horizontal. angleDeg = 90 gives
// a parallel-sided slot 2r tall (WireGroove(1,3,90) spans Y -1..1); smaller
// angles flare the mouth (WireGroove(1,3,60) reaches Y 2.887). The roof is
// clamped at y = 3.5*r so a shallow tail angle cannot run away.
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
