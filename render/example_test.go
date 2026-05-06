package render_test

import (
	"github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Use the parallel marching-cubes octree renderer directly. Most callers
// should just use part.STL(path, cellsPerMM) — reach for this when you
// want to share one renderer across multiple outputs or pick a different
// implementation than the default.
func ExampleNewMarchingCubesOctreeParallel() {
	part := solid.Box(v3.XYZ(20, 10, 5), 0)
	r := render.NewMarchingCubesOctreeParallel(solid.CellsFor(part, 5.0))
	// render.ToSTL(part, "box.stl", r)
	_ = r
}

// Render a 2D shape through the marching-squares quadtree renderer, used
// internally by Shape.ToDXF / ToSVG.
func ExampleNewMarchingSquaresQuadtree() {
	r := render.NewMarchingSquaresQuadtree(200)
	// render.ToDXFWith(profile, "out.dxf", r)
	_ = r
}
