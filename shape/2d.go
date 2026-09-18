package shape

import (
	"math"

	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

// Star returns a star polygon with the given number of points, tips exactly
// at radius outer and valleys exactly at radius inner, centered on the
// origin. The first tip points along +X and the tips are spaced 360/points
// degrees apart. Corners are sharp; round them with Offset if you want, and
// remember Offset changes the size.
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
	return Polygon(pts)
}

// Hexagon returns a regular hexagon centered on the origin, with radius as
// the circumradius (center to vertex) and a vertex on +X.
// Vertices sit every 60 degrees starting on +X, so flats face +Y and -Y and
// the across-flats distance is radius*sqrt(3): Hexagon(10) spans X -10..10
// and Y -8.660..8.660.
func Hexagon(radius float64) *Shape {
	var pts []v2.Vec
	for i := 0; i < 6; i++ {
		angle := float64(i) * (2 * math.Pi / 6)
		pts = append(pts, v2.Vec{X: radius * math.Cos(angle), Y: radius * math.Sin(angle)})
	}
	return Polygon(pts)
}

// Triangle returns an equilateral triangle centered on the origin and
// inscribed in a circle of the given radius, with its apex on +Y.
// The base is horizontal at y = -radius/2 and the base corners are at
// x = +/-radius*sqrt(3)/2, so Triangle(10) spans X -8.660..8.660 and
// Y -5..10.
func Triangle(radius float64) *Shape {
	var pts []v2.Vec
	for i := 0; i < 3; i++ {
		angle := float64(i)*(2*math.Pi/3) + (math.Pi / 2)
		pts = append(pts, v2.Vec{X: radius * math.Cos(angle), Y: radius * math.Sin(angle)})
	}
	return Polygon(pts)
}

// Cross returns a plus sign centered on the origin: two bars of length
// width and the given thickness, one along X and one along Y.
// The bounding box is width by width, so Cross(10, 2) spans -5..5 on both
// axes with arms 2 wide.
func Cross(width, thickness float64) *Shape {
	vBar := Rect(v2.Vec{X: thickness, Y: width}, 0)
	hBar := Rect(v2.Vec{X: width, Y: thickness}, 0)
	return vBar.Union(hBar)
}
