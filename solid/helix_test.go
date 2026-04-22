package solid_test

import (
	"testing"

	"github.com/snowbldr/fluent-sdfx/mesh"
	flrender "github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

// TestSweepHelixFlatEndsWatertight verifies that flatEnds=true produces a
// closed mesh for a variety of profiles. A hollow end cap would leave
// boundary edges on the open face.
func TestSweepHelixFlatEndsWatertight(t *testing.T) {
	const cellsPerMM = 8.0

	cases := []struct {
		name                  string
		profile               *shape.Shape
		radius, turns, height float64
	}{
		{"square", shape.Rect(v2.XY(1, 1), 0), 5, 3, 15},
		{"circle", shape.Circle(1.0), 6, 4, 20},
		{"tall_rect", shape.Rect(v2.XY(1, 3), 0), 8, 2, 12},
		{"ribbon", shape.Rect(v2.XY(3, 0.6), 0), 8, 3, 12},
		{"short", shape.Rect(v2.XY(1, 1), 0), 5, 1, 5},
		{"iso_thread", shape.ISOThread(5, 2, true), 5, 4, 12},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := solid.SweepHelix(tc.profile, tc.radius, tc.turns, tc.height, true)
			cells := solid.CellsFor(s, cellsPerMM)
			tris := mesh.CollectTriangles(s, flrender.NewMarchingCubesOctreeParallel(cells))
			if len(tris) == 0 {
				t.Fatalf("no triangles rendered")
			}
			ok, boundary := mesh.IsWatertight(tris)
			if !ok {
				t.Errorf("mesh not watertight: %d boundary edges", boundary)
			}
		})
	}
}
