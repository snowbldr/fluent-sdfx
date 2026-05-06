package shape_test

import (
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/layout"
	"github.com/snowbldr/fluent-sdfx/shape"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

// BenchmarkPolygon measures construction of a 100-vertex polygon.
func BenchmarkPolygon(b *testing.B) {
	const n = 100
	pts := make([]v2.Vec, n)
	for i := 0; i < n; i++ {
		theta := 2 * math.Pi * float64(i) / float64(n)
		pts[i] = v2.XY(10*math.Cos(theta), 10*math.Sin(theta))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = shape.Polygon(pts)
	}
}

// BenchmarkRectCornersMulti measures the typical "circle.Multi(layout.RectCorners2(...))" pipeline.
func BenchmarkRectCornersMulti(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = shape.Circle(1).Multi(layout.RectCorners2(80, 50)...)
	}
}
