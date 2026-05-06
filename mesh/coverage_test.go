package mesh_test

import (
	"testing"

	"github.com/snowbldr/fluent-sdfx/mesh"
	"github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// --- A closed cube must be watertight: zero boundary edges. ---

func TestCubeIsWatertight(t *testing.T) {
	cube := solid.Box(v3.XYZ(10, 10, 10), 0)
	cells := solid.CellsFor(cube, 6)
	tris := mesh.CollectTriangles(cube, render.NewMarchingCubesOctreeParallel(cells))
	if len(tris) == 0 {
		t.Fatal("no triangles produced for a 10mm cube; renderer regression")
	}
	ok, n := mesh.IsWatertight(tris)
	if !ok {
		t.Fatalf("cube mesh is not watertight: %d boundary edges (%d triangles)", n, len(tris))
	}
	if mesh.CountBoundaryEdges(tris) != 0 {
		t.Errorf("CountBoundaryEdges = %d, want 0 for sealed cube", mesh.CountBoundaryEdges(tris))
	}
}

// Sphere should also be closed.
func TestSphereIsWatertight(t *testing.T) {
	s := solid.Sphere(5)
	cells := solid.CellsFor(s, 6)
	tris := mesh.CollectTriangles(s, render.NewMarchingCubesOctreeParallel(cells))
	if len(tris) == 0 {
		t.Fatal("no triangles produced for sphere")
	}
	ok, n := mesh.IsWatertight(tris)
	if !ok {
		t.Fatalf("sphere mesh is not watertight: %d boundary edges (%d triangles)", n, len(tris))
	}
}

// --- ToTriangles vs. CollectTriangles: same set, different value/pointer types. ---

func TestToTrianglesAndCollectTrianglesAgree(t *testing.T) {
	s := solid.Box(v3.XYZ(8, 8, 8), 0)
	cells := solid.CellsFor(s, 4)
	r := render.NewMarchingCubesOctreeParallel(cells)
	pTris := mesh.ToTriangles(s, r)
	r2 := render.NewMarchingCubesOctreeParallel(cells)
	vTris := mesh.CollectTriangles(s, r2)
	if len(pTris) != len(vTris) {
		t.Errorf("ToTriangles n=%d, CollectTriangles n=%d — should match", len(pTris), len(vTris))
	}
}

// --- VertexToLine helper. ---

func TestVertexToLineOpenAndClosed(t *testing.T) {
	// A 3-vertex polyline: open → 2 segments, closed → 3 segments.
	pts := []v2.Vec{v2.XY(0, 0), v2.XY(10, 0), v2.XY(10, 10)}

	open := mesh.VertexToLine(pts, false)
	if len(open) != 2 {
		t.Errorf("open polyline of 3 vertices → %d lines, want 2", len(open))
	}

	closed := mesh.VertexToLine(pts, true)
	if len(closed) != 3 {
		t.Errorf("closed polyline of 3 vertices → %d lines, want 3", len(closed))
	}
}

// --- TriangleI canonicalization. ---

func TestTriangleICanonical(t *testing.T) {
	tri := mesh.TriangleI{5, 2, 7}
	tri.Canonical()
	// Smallest index should be first.
	if tri[0] != 2 {
		t.Errorf("Canonical: first index = %d, want 2 (smallest)", tri[0])
	}
}

func TestTriangleISetEquals(t *testing.T) {
	a := mesh.TriangleISet{{0, 1, 2}, {1, 2, 3}}
	b := mesh.TriangleISet{{1, 2, 3}, {0, 1, 2}}
	if !a.Equals(b) {
		t.Errorf("Equals should ignore triangle order")
	}

	c := mesh.TriangleISet{{0, 1, 2}, {1, 2, 4}}
	if a.Equals(c) {
		t.Errorf("Equals should detect different triangles")
	}
}

// --- Delaunay smoke: 4 points → 2 triangles forming a quad. ---

func TestDelaunay2dQuad(t *testing.T) {
	pts := v2.VecSet{
		v2.XY(0, 0), v2.XY(10, 0), v2.XY(10, 10), v2.XY(0, 10),
	}
	tris := mesh.Delaunay2d(pts)
	if len(tris) != 2 {
		t.Errorf("Delaunay2d of 4-point quad → %d triangles, want 2", len(tris))
	}
}
