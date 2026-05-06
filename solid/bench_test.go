package solid_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/snowbldr/fluent-sdfx/layout"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// BenchmarkBoxConstruct measures the cost of constructing a single Box (no rendering).
func BenchmarkBoxConstruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = solid.Box(v3.XYZ(20, 20, 20), 0)
	}
}

// BenchmarkSolidUnion measures union construction over N parts.
// Run with -benchmem to spot allocation regressions; if a refactor
// introduces N² behavior, the n=100 case will balloon vs the n=10 case.
func BenchmarkSolidUnion(b *testing.B) {
	for _, n := range []int{2, 10, 100} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			parts := make([]*solid.Solid, n)
			for i := range parts {
				parts[i] = solid.Box(v3.XYZ(10, 10, 10), 0).TranslateX(float64(i * 15))
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = solid.UnionAll(parts...)
			}
		})
	}
}

// BenchmarkSolidCutMulti measures the typical "drill 4 holes via .Multi(layout.RectCorners(...)...)" pipeline.
func BenchmarkSolidCutMulti(b *testing.B) {
	body := solid.Box(v3.XYZ(80, 50, 10), 0)
	for i := 0; i < b.N; i++ {
		hole := solid.Cylinder(20, 2.5, 0).Multi(layout.RectCorners(70, 40)...)
		_ = body.Cut(hole)
	}
}

// BenchmarkBounds measures Bounds() on a moderately complex composite solid.
func BenchmarkBounds(b *testing.B) {
	s := solid.Box(v3.XYZ(80, 50, 10), 0).
		Cut(solid.Cylinder(20, 2.5, 0).Multi(layout.RectCorners(70, 40)...)).
		Union(solid.Sphere(5).TranslateZ(10))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Bounds()
	}
}

// BenchmarkSTL measures a full STL render at a moderate density.
// Skipped under -short because this is the dominant benchmark cost.
func BenchmarkSTL(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping STL render benchmark in -short mode")
	}
	dir := b.TempDir()
	s := solid.Box(v3.XYZ(20, 20, 20), 1).
		Cut(solid.Cylinder(25, 4, 0).Multi(layout.RectCorners(12, 12)...))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path := filepath.Join(dir, fmt.Sprintf("bench-%d.stl", i))
		s.STL(path, 2.0)
		_ = os.Remove(path)
	}
}
