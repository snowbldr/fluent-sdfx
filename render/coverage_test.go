package render_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// --- ToSTL smoke: writes a real binary STL to a temp file. ---

func TestToSTLWritesNonEmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cube.stl")

	cube := solid.Box(v3.XYZ(5, 5, 5), 0)
	cells := solid.CellsFor(cube, 4)
	render.ToSTL(cube, path, render.NewMarchingCubesOctreeParallel(cells))

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("STL file not created at %s: %v", path, err)
	}
	if info.Size() == 0 {
		t.Errorf("STL file at %s is empty", path)
	}
	// Binary STL header is 80 bytes + 4-byte tri count + 50 bytes/triangle.
	// At minimum we expect well over the 84-byte header for a closed cube.
	if info.Size() < 200 {
		t.Errorf("STL file size = %d bytes, want > 200 for a cube mesh", info.Size())
	}
}

// ToSTL with decimation should still produce a valid file.
func TestToSTLWithDecimation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "decimated.stl")

	s := solid.Sphere(5)
	cells := solid.CellsFor(s, 6)
	render.ToSTL(s, path, render.NewMarchingCubesOctreeParallel(cells), 0.5) // remove 50%

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("STL file not created at %s: %v", path, err)
	}
	if info.Size() == 0 {
		t.Errorf("decimated STL file at %s is empty", path)
	}
}

// Solid.STL should round-trip to the parallel octree renderer too.
func TestSolidSTLConvenience(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sphere.stl")
	solid.Sphere(3).STL(path, 4)

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("STL file not created at %s: %v", path, err)
	}
	if info.Size() == 0 {
		t.Errorf("STL file at %s is empty", path)
	}
}

// --- 2D renderers smoke. ---

func TestNewMarchingCubesOctreeParallelReturnsNonNil(t *testing.T) {
	if r := render.NewMarchingCubesOctreeParallel(32); r == nil {
		t.Errorf("NewMarchingCubesOctreeParallel(32) returned nil")
	}
}

func TestNewMarchingSquaresQuadtreeReturnsNonNil(t *testing.T) {
	if r := render.NewMarchingSquaresQuadtree(32); r == nil {
		t.Errorf("NewMarchingSquaresQuadtree(32) returned nil")
	}
}

func TestNewDualContouring2DReturnsNonNil(t *testing.T) {
	if r := render.NewDualContouring2D(32); r == nil {
		t.Errorf("NewDualContouring2D(32) returned nil")
	}
}
