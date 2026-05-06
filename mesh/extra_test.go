package mesh_test

import (
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/snowbldr/fluent-sdfx/mesh"
	"github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// --- Triangle3 / Triangle2 / Line2 type-alias smoke. ---
//
// These are aliases to the v2/v3 vector types, so a value made with the
// underlying type is assignable to the alias and round-trips through
// SDF()/FromSDF.

func TestTriangle3AliasRoundTrip(t *testing.T) {
	tri := mesh.Triangle3{v3.XYZ(0, 0, 0), v3.XYZ(1, 0, 0), v3.XYZ(0, 1, 0)}
	st := tri.SDF()
	back := v3.FromSDFTriangle3(st)
	for i := 0; i < 3; i++ {
		if back[i] != tri[i] {
			t.Errorf("Triangle3 vertex %d: got %v want %v", i, back[i], tri[i])
		}
	}
}

func TestTriangle2AliasRoundTrip(t *testing.T) {
	tri := mesh.Triangle2{v2.XY(0, 0), v2.XY(1, 0), v2.XY(0, 1)}
	st := tri.SDF()
	back := v2.FromSDFTriangle2(st)
	for i := 0; i < 3; i++ {
		if back[i] != tri[i] {
			t.Errorf("Triangle2 vertex %d: got %v want %v", i, back[i], tri[i])
		}
	}
}

func TestLine2AliasRoundTrip(t *testing.T) {
	l := mesh.Line2{v2.XY(0, 0), v2.XY(3, 4)}
	sl := l.SDF()
	back := v2.FromSDFLine2(sl)
	if back != l {
		t.Errorf("Line2: got %v want %v", back, l)
	}
}

// --- TriangleI.ToTriangle2: index lookup ---

func TestTriangleIToTriangle2(t *testing.T) {
	pts := v2.VecSet{v2.XY(0, 0), v2.XY(10, 0), v2.XY(0, 10), v2.XY(5, 5)}
	tri := mesh.TriangleI{0, 1, 2}
	got := tri.ToTriangle2(pts)
	want := mesh.Triangle2{pts[0], pts[1], pts[2]}
	for i := 0; i < 3; i++ {
		if got[i] != want[i] {
			t.Errorf("ToTriangle2 vertex %d: got %v want %v", i, got[i], want[i])
		}
	}
}

// --- TriangleISet.Canonical returns canonical-form copies. ---

func TestTriangleISetCanonical(t *testing.T) {
	in := mesh.TriangleISet{{5, 2, 7}, {9, 1, 4}}
	out := in.Canonical()
	if len(out) != len(in) {
		t.Fatalf("Canonical length mismatch: got %d want %d", len(out), len(in))
	}
	for i, tri := range out {
		min := tri[0]
		for _, v := range tri {
			if v < min {
				min = v
			}
		}
		if tri[0] != min {
			t.Errorf("Canonical[%d] first = %d, want smallest = %d", i, tri[0], min)
		}
	}
}

// --- ToTriangles vs CollectTriangles agreement on additional primitives. ---

func TestToCollectTrianglesAgreeSphere(t *testing.T) {
	s := solid.Sphere(4)
	cells := solid.CellsFor(s, 4)
	pTris := mesh.ToTriangles(s, render.NewMarchingCubesOctreeParallel(cells))
	vTris := mesh.CollectTriangles(s, render.NewMarchingCubesOctreeParallel(cells))
	if len(pTris) != len(vTris) {
		t.Errorf("Sphere ToTriangles n=%d, CollectTriangles n=%d", len(pTris), len(vTris))
	}
	if len(pTris) == 0 {
		t.Fatal("no triangles for sphere")
	}
}

func TestToCollectTrianglesAgreeCylinder(t *testing.T) {
	s := solid.Cylinder(8, 3, 0)
	cells := solid.CellsFor(s, 4)
	pTris := mesh.ToTriangles(s, render.NewMarchingCubesOctreeParallel(cells))
	vTris := mesh.CollectTriangles(s, render.NewMarchingCubesOctreeParallel(cells))
	if len(pTris) != len(vTris) {
		t.Errorf("Cylinder ToTriangles n=%d, CollectTriangles n=%d", len(pTris), len(vTris))
	}
	if len(pTris) == 0 {
		t.Fatal("no triangles for cylinder")
	}
}

// --- IsWatertight on a cylinder (closed, capped). ---

func TestCylinderIsWatertight(t *testing.T) {
	s := solid.Cylinder(8, 3, 0)
	cells := solid.CellsFor(s, 5)
	tris := mesh.CollectTriangles(s, render.NewMarchingCubesOctreeParallel(cells))
	if len(tris) == 0 {
		t.Fatal("no triangles for cylinder")
	}
	ok, n := mesh.IsWatertight(tris)
	if !ok {
		t.Fatalf("cylinder mesh not watertight: %d boundary edges (%d triangles)", n, len(tris))
	}
}

// --- Open mesh has non-zero boundary edges. ---
//
// Construct a single triangle by hand: it has 3 boundary edges (each edge is
// shared by only one triangle), so CountBoundaryEdges should return 3.

func TestCountBoundaryEdgesOpenMesh(t *testing.T) {
	tri := mesh.Triangle3{v3.XYZ(0, 0, 0), v3.XYZ(1, 0, 0), v3.XYZ(0, 1, 0)}
	tris := []mesh.Triangle3{tri}
	n := mesh.CountBoundaryEdges(tris)
	if n != 3 {
		t.Errorf("single triangle CountBoundaryEdges = %d, want 3", n)
	}
	ok, n2 := mesh.IsWatertight(tris)
	if ok {
		t.Errorf("single triangle should not be watertight")
	}
	if n2 != 3 {
		t.Errorf("IsWatertight returned boundary count = %d, want 3", n2)
	}
}

// Two triangles sharing one edge → 4 boundary edges.
func TestCountBoundaryEdgesTwoTrianglesSharingEdge(t *testing.T) {
	a := mesh.Triangle3{v3.XYZ(0, 0, 0), v3.XYZ(1, 0, 0), v3.XYZ(0, 1, 0)}
	b := mesh.Triangle3{v3.XYZ(1, 0, 0), v3.XYZ(1, 1, 0), v3.XYZ(0, 1, 0)}
	tris := []mesh.Triangle3{a, b}
	if n := mesh.CountBoundaryEdges(tris); n != 4 {
		t.Errorf("CountBoundaryEdges = %d, want 4 for two triangles sharing one edge", n)
	}
}

// --- SaveSTL to tmp file. Verify size > 0 and parse the binary STL header
// to recover the triangle count, asserting it matches what we wrote. ---

func TestSaveSTLBinary(t *testing.T) {
	// Three loose triangles — easy to count.
	tris := []*mesh.Triangle3{
		{v3.XYZ(0, 0, 0), v3.XYZ(1, 0, 0), v3.XYZ(0, 1, 0)},
		{v3.XYZ(0, 0, 0), v3.XYZ(2, 0, 0), v3.XYZ(0, 2, 0)},
		{v3.XYZ(0, 0, 0), v3.XYZ(3, 0, 0), v3.XYZ(0, 3, 0)},
	}
	dir, err := os.MkdirTemp("", "mesh-savestl-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	path := filepath.Join(dir, "tris.stl")

	mesh.SaveSTL(path, tris)

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("STL file is empty")
	}

	// Parse the binary STL: 80-byte header, then uint32 little-endian tri count.
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	hdr := make([]byte, 80)
	if _, err := f.Read(hdr); err != nil {
		t.Fatalf("read header: %v", err)
	}
	var n uint32
	if err := binary.Read(f, binary.LittleEndian, &n); err != nil {
		t.Fatalf("read count: %v", err)
	}
	if int(n) != len(tris) {
		t.Errorf("STL triangle count = %d, want %d", n, len(tris))
	}
}

// --- Delaunay smoke: square (4 pts → 2 tris already covered) and a
// 5-point set. Also exercise Delaunay2dSlow on the same inputs. ---

func TestDelaunay2dSlowQuad(t *testing.T) {
	pts := v2.VecSet{v2.XY(0, 0), v2.XY(10, 0), v2.XY(10, 10), v2.XY(0, 10)}
	tris := mesh.Delaunay2dSlow(pts)
	// The "slow" reference algorithm doesn't dedupe to a minimal triangulation
	// the way the fast Delaunay2d does — it just needs to return a valid set
	// of triangles covering the convex hull. Assert non-empty.
	if len(tris) == 0 {
		t.Errorf("Delaunay2dSlow quad → 0 tris, want > 0")
	}
}

func TestDelaunay2dPentagon(t *testing.T) {
	pts := v2.VecSet{
		v2.XY(0, 0), v2.XY(10, 0), v2.XY(10, 10), v2.XY(0, 10), v2.XY(5, 5),
	}
	tris := mesh.Delaunay2d(pts)
	if len(tris) == 0 {
		t.Fatal("no triangles for 5-point set")
	}
}

// Delaunay2dE on degenerate input (empty set) should return an error rather
// than panicking. Exercises the non-panic error-return path.
func TestDelaunay2dEEmpty(t *testing.T) {
	_, err := mesh.Delaunay2dE(v2.VecSet{})
	if err == nil {
		t.Error("expected error from Delaunay2dE on empty input")
		return
	}
	// Touch errors.Is so the linter is happy with the import.
	_ = errors.Is(err, err)
}

// --- VertexToLine: closed connects last to first. ---

func TestVertexToLineClosesLoop(t *testing.T) {
	pts := []v2.Vec{v2.XY(0, 0), v2.XY(1, 0), v2.XY(1, 1), v2.XY(0, 1)}
	closed := mesh.VertexToLine(pts, true)
	if len(closed) != 4 {
		t.Fatalf("closed loop of 4 verts → %d lines, want 4", len(closed))
	}
	last := closed[3]
	if last[0] != pts[3] || last[1] != pts[0] {
		t.Errorf("closing edge = %v, want [%v %v]", last, pts[3], pts[0])
	}
}

// --- Conversion helper at the package boundary. ---

func TestSDFLinesConverts(t *testing.T) {
	a := mesh.Line2{v2.XY(0, 0), v2.XY(1, 1)}
	b := mesh.Line2{v2.XY(2, 2), v2.XY(3, 3)}
	in := []*mesh.Line2{&a, &b}
	out := mesh.SDFLines(in)
	if len(out) != len(in) {
		t.Fatalf("SDFLines len = %d, want %d", len(out), len(in))
	}
	if out[0] == nil || out[1] == nil {
		t.Fatal("SDFLines returned nil entry")
	}
	if out[0][0].X != 0 || out[0][1].Y != 1 {
		t.Errorf("SDFLines[0] roundtrip mismatch: %v", out[0])
	}
}
