package validate_test

import (
	"testing"

	"github.com/snowbldr/fluent-sdfx/mesh"
	"github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// benchSolid is a moderate composite (cube with 4 holes) used as input
// for mesh-based metrics. Re-rendered for each benchmark setup but the
// rendered triangle list is reused across the b.N inner loop.
func benchSolid() *solid.Solid {
	body := solid.Box(v3.XYZ(20, 20, 20), 1)
	hole := solid.Cylinder(25, 2.5, 0).
		TranslateXY(7, 7).Union(
		solid.Cylinder(25, 2.5, 0).TranslateXY(-7, 7),
		solid.Cylinder(25, 2.5, 0).TranslateXY(7, -7),
		solid.Cylinder(25, 2.5, 0).TranslateXY(-7, -7),
	)
	return body.Cut(hole)
}

// benchTris renders benchSolid once at cellsPerMM=4 and returns the
// triangle list. ~10K triangles for the box-with-4-holes geometry.
func benchTris(b *testing.B) []mesh.Triangle3 {
	b.Helper()
	s := benchSolid()
	return mesh.CollectTriangles(s, render.NewMarchingCubesOctreeParallel(solid.CellsFor(s, 4.0)))
}

// BenchmarkVolume measures the per-iteration cost of the volume integral
// on a precomputed triangle mesh.
func BenchmarkVolume(b *testing.B) {
	tris := benchTris(b)
	b.Logf("mesh triangle count: %d", len(tris))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validate.Volume(tris)
	}
}

// BenchmarkOverhangArea measures the cost of overhang detection at the
// FDM 45° threshold on a precomputed mesh.
func BenchmarkOverhangArea(b *testing.B) {
	tris := benchTris(b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validate.OverhangArea(tris, 45.0)
	}
}

// BenchmarkOf measures the full validate pipeline (render + every metric)
// at cellsPerMM=4. This is the dominant cost in this package; gate behind
// -short so quick CI passes still finish fast.
func BenchmarkOf(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping full pipeline benchmark in -short mode")
	}
	s := benchSolid()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validate.Of(s, 4.0)
	}
}
