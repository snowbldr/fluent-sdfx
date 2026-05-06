package mesh_test

import (
	"testing"

	"github.com/snowbldr/fluent-sdfx/mesh"
	"github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// BenchmarkCountBoundaryEdges measures boundary-edge counting on a
// precomputed mesh. Watertight check is the main consumer of this in
// validate.Of, so a regression here propagates everywhere.
func BenchmarkCountBoundaryEdges(b *testing.B) {
	body := solid.Box(v3.XYZ(20, 20, 20), 1)
	hole := solid.Cylinder(25, 2.5, 0).
		TranslateXY(7, 7).Union(
		solid.Cylinder(25, 2.5, 0).TranslateXY(-7, 7),
		solid.Cylinder(25, 2.5, 0).TranslateXY(7, -7),
		solid.Cylinder(25, 2.5, 0).TranslateXY(-7, -7),
	)
	s := body.Cut(hole)
	tris := mesh.CollectTriangles(s, render.NewMarchingCubesOctreeParallel(solid.CellsFor(s, 4.0)))
	b.Logf("mesh triangle count: %d", len(tris))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mesh.CountBoundaryEdges(tris)
	}
}
