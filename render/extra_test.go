package render_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v2i "github.com/snowbldr/fluent-sdfx/vec/v2i"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/sdfx/sdf"
)

// --- ToSTL panic paths ---
//
// ToSTL panics on filesystem errors (per package convention). The atomic
// writer first creates a temp file in the *destination directory*, so a
// path whose parent doesn't exist will trigger the panic in
// os.CreateTemp.

func TestToSTLPanicsOnBadPath(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on bad path")
		}
	}()
	cube := solid.Box(v3.XYZ(2, 2, 2), 0)
	r := render.NewMarchingCubesOctreeParallel(solid.CellsFor(cube, 2))
	render.ToSTL(cube, "/this/path/should/not/exist/badness/foo.stl", r)
}

// ToSTL with decimation=0 should keep all triangles (decimation disabled).
func TestToSTLNoDecimation(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nodec.stl")
	cube := solid.Box(v3.XYZ(4, 4, 4), 0)
	r := render.NewMarchingCubesOctreeParallel(solid.CellsFor(cube, 2))
	render.ToSTL(cube, path, r, 0)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("STL not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("STL file empty")
	}
}

// ToSTL with decimation outside the (0,1) range falls through unchanged
// (covers the !(>0 && <1) branch).
func TestToSTLDecimationOutOfRange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "fullkeep.stl")
	cube := solid.Box(v3.XYZ(4, 4, 4), 0)
	r := render.NewMarchingCubesOctreeParallel(solid.CellsFor(cube, 2))
	render.ToSTL(cube, path, r, 1.0) // not in (0,1) — disables
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("STL not created: %v", err)
	}
}

// --- To3MF ---

func TestTo3MF(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cube.3mf")
	cube := solid.Box(v3.XYZ(4, 4, 4), 0)
	render.To3MF(cube, path, solid.CellsFor(cube, 2))
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("3MF not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("3MF file empty")
	}
}

// --- ToDXF / ToSVG (low resolution, write to tmp) ---

func TestToDXF(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "circ.dxf")
	c := shape.Circle(5)
	render.ToDXF(c, path, 32)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("DXF not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("DXF file empty")
	}
}

func TestToDXFWith(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "circ-with.dxf")
	c := shape.Circle(5)
	r := render.NewMarchingSquaresQuadtree(32)
	render.ToDXFWith(c, path, r)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("DXF not created: %v", err)
	}
}

func TestToSVG(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "circ.svg")
	c := shape.Circle(5)
	render.ToSVG(c, path, 32)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("SVG not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("SVG file empty")
	}
}

func TestToSVGWith(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "circ-with.svg")
	c := shape.Circle(5)
	r := render.NewMarchingSquaresQuadtree(32)
	render.ToSVGWith(c, path, r)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("SVG not created: %v", err)
	}
}

// --- ToPNG ---

func TestToPNG(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "circ.png")
	c := shape.Circle(5)
	bb := c.Bounds()
	render.ToPNG(c, path, bb, 64, 64)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("PNG not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("PNG file empty")
	}
}

// ToPNG panics on a path that can't be written (parent dir doesn't exist).
func TestToPNGPanicsOnBadPath(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on bad path")
		}
	}()
	c := shape.Circle(2)
	render.ToPNG(c, "/this/path/should/not/exist/foo.png", c.Bounds(), 8, 8)
}

// --- SaveDXF / SaveSVG free functions ---

func TestSaveDXF(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lines.dxf")
	a := v2.Line2{v2.XY(0, 0), v2.XY(10, 0)}
	b := v2.Line2{v2.XY(10, 0), v2.XY(10, 10)}
	render.SaveDXF(path, []*v2.Line2{&a, &b})
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("DXF not created: %v", err)
	}
}

func TestSaveDXFPanicsOnBadPath(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on bad path")
		}
	}()
	a := v2.Line2{v2.XY(0, 0), v2.XY(1, 0)}
	render.SaveDXF("/this/path/should/not/exist/lines.dxf", []*v2.Line2{&a})
}

func TestSaveSVG(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lines.svg")
	a := v2.Line2{v2.XY(0, 0), v2.XY(10, 0)}
	b := v2.Line2{v2.XY(10, 0), v2.XY(10, 10)}
	render.SaveSVG(path, "stroke:black", []*v2.Line2{&a, &b})
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("SVG not created: %v", err)
	}
}

func TestSaveSVGPanicsOnBadPath(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on bad path")
		}
	}()
	a := v2.Line2{v2.XY(0, 0), v2.XY(1, 0)}
	render.SaveSVG("/this/path/should/not/exist/lines.svg", "stroke:black", []*v2.Line2{&a})
}

// --- Poly: write a polygon to DXF ---

func TestPoly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "poly.dxf")
	p := sdf.NewPolygon()
	p.Add(0, 0)
	p.Add(10, 0)
	p.Add(10, 10)
	p.Add(0, 10)
	p.Close()
	render.Poly(p, path)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("DXF not created: %v", err)
	}
}

// --- PNG drawing-target API ---
//
// NewPNG returns a target. Exercise its methods. Save() writes the file.

func TestPNGDrawingTarget(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "draw.png")
	bb := v2.NewBox(v2.XY(0, 0), v2.XY(20, 20))
	p, err := render.NewPNG(path, bb, v2i.XY(64, 64))
	if err != nil {
		t.Fatalf("NewPNG: %v", err)
	}
	if p == nil {
		t.Fatal("NewPNG returned nil")
	}
	// Render an SDF2 onto the canvas.
	c := shape.Circle(5)
	p.RenderSDF2(c)
	p.RenderSDF2MinMax(c, -2, 2)
	// Draw geometry primitives.
	p.Line(v2.XY(-5, -5), v2.XY(5, 5))
	p.Lines(v2.VecSet{v2.XY(-5, 0), v2.XY(0, 5), v2.XY(5, 0)})
	p.Triangle(v2.Triangle2{v2.XY(0, 0), v2.XY(3, 0), v2.XY(0, 3)})
	if err := p.Save(); err != nil {
		t.Fatalf("PNG.Save: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("PNG not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("PNG file empty")
	}
}

// NewPNG with degenerate pixel dims (zero box) returns an error from
// sdf.NewMap2; this exercises the error-propagation branch of NewPNG.
func TestNewPNGBadInputs(t *testing.T) {
	// Zero-size bounding box → NewMap2 errors.
	bb := v2.NewBox(v2.XY(0, 0), v2.XY(0, 0))
	_, err := render.NewPNG("/tmp/should-not-be-written.png", bb, v2i.XY(8, 8))
	if err == nil {
		t.Error("expected error from NewPNG with zero-size bounding box")
	}
}

// PNG.Save should fail (return error) when the file path can't be opened.
func TestPNGSaveBadPath(t *testing.T) {
	bb := v2.NewBox(v2.XY(0, 0), v2.XY(2, 2))
	p, err := render.NewPNG("/this/path/should/not/exist/draw.png", bb, v2i.XY(8, 8))
	if err != nil {
		// Some implementations might error at construction — that's also fine.
		return
	}
	if err := p.Save(); err == nil {
		t.Error("expected error from PNG.Save with bad path")
	}
}

// --- DXF drawing-target API ---

func TestDXFDrawingTarget(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "draw.dxf")
	d := render.NewDXF(path)
	if d == nil {
		t.Fatal("NewDXF returned nil")
	}
	bb := v2.NewBox(v2.XY(0, 0), v2.XY(20, 20))
	d.Box(bb)
	a := v2.Line2{v2.XY(0, 0), v2.XY(10, 0)}
	d.Line(&a)
	b := v2.Line2{v2.XY(10, 0), v2.XY(10, 10)}
	c := v2.Line2{v2.XY(10, 10), v2.XY(0, 10)}
	d.Lines([]*v2.Line2{&b, &c})
	d.Points(v2.VecSet{v2.XY(1, 1), v2.XY(2, 2)}, 0.25)
	d.Triangle(v2.Triangle2{v2.XY(0, 0), v2.XY(3, 0), v2.XY(0, 3)})
	if err := d.Save(); err != nil {
		t.Fatalf("DXF.Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("DXF not created: %v", err)
	}
}

// --- Renderer constructor smoke (the wrapped ones). ---

func TestRendererConstructorsNonNil(t *testing.T) {
	if r := render.NewMarchingCubesOctreeParallel(16); r == nil {
		t.Error("NewMarchingCubesOctreeParallel returned nil")
	}
	if r := render.NewMarchingSquaresQuadtree(16); r == nil {
		t.Error("NewMarchingSquaresQuadtree returned nil")
	}
	if r := render.NewDualContouring2D(16); r == nil {
		t.Error("NewDualContouring2D returned nil")
	}
}
