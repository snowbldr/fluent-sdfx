// SweepHelix showcase: sweeps a variety of 2D profiles along helical paths
// with and without flatEnds. Renders each to its own STL plus a combined
// layout file so you can eyeball every variant side-by-side.
package main

import (
	"fmt"

	"github.com/snowbldr/fluent-sdfx/mesh"
	flrender "github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const cellsPerMM = 8.0

type variant struct {
	name                  string
	profile               *shape.Shape
	radius, turns, height float64
	flatEnds              bool
}

func variants() []variant {
	return []variant{
		// Small square profile — classic sweep-a-box-along-a-helix use case.
		{"sq1_flat", shape.Rect(v2.XY(1, 1), 0), 5, 3, 15, true},
		{"sq1_open", shape.Rect(v2.XY(1, 1), 0), 5, 3, 15, false},

		// Circular cross-section: makes a tube-like spring shape.
		{"circle_flat", shape.Circle(1.0), 6, 4, 20, true},
		{"circle_open", shape.Circle(1.0), 6, 4, 20, false},

		// Radially wide, axially thin — like a flat ribbon wrapped around.
		{"ribbon_flat", shape.Rect(v2.XY(3, 0.6), 0), 8, 3, 12, true},

		// 5-point star — spiky profile corkscrews up the helix.
		{"star_flat", shape.Star(1.5, 0.7, 5), 6, 3, 15, true},

		// Hexagonal profile — nut-like cross-section swept into a coil.
		{"hex_flat", shape.Hexagon(1.2), 6, 3, 14, true},

		// Triangular profile — points outward as it winds.
		{"tri_flat", shape.Triangle(1.3), 6, 3, 14, true},

		// Cross/plus profile — four-armed cross-section.
		{"cross_flat", shape.Cross(2.0, 0.6), 7, 3, 15, true},
	}
}

func main() {
	vs := variants()

	var all []*mesh.Triangle3
	xOffset := 0.0
	prevR := 0.0

	for i, v := range vs {
		s := solid.SweepHelix(v.profile, v.radius, v.turns, v.height, v.flatEnds)

		individual := fmt.Sprintf("helix_%s.stl", v.name)
		s.STL(individual, cellsPerMM)

		bb := s.BoundingBox()
		r := bb.Max.X
		cells := solid.CellsFor(s, cellsPerMM)
		tris := mesh.CollectTriangles(s, flrender.NewMarchingCubesOctreeParallel(cells))

		if i > 0 {
			gap := 3.0
			if space := max(prevR, r) * 0.5; space > gap {
				gap = space
			}
			xOffset += gap
		}
		for j := range tris {
			t := tris[j]
			t[0] = t[0].Add(v3.X(xOffset + r))
			t[1] = t[1].Add(v3.X(xOffset + r))
			t[2] = t[2].Add(v3.X(xOffset + r))
			all = append(all, &t)
		}
		xOffset += r * 2
		prevR = r
	}

	fmt.Printf("writing combined layout: %d triangles\n", len(all))
	mesh.SaveSTL("helix_all.stl", all)
}
