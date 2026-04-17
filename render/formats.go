package render

import (
	"fmt"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2sdf "github.com/deadsy/sdfx/vec/v2"
	v2isdf "github.com/deadsy/sdfx/vec/v2i"
	v3sdf "github.com/deadsy/sdfx/vec/v3"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	"github.com/snowbldr/fluent-sdfx/vec/v2i"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Render3 is the sdfx 3D renderer interface. Use it for custom lower-level rendering.
type Render3 = render.Render3

// Render2 is the sdfx 2D renderer interface. Use it for custom lower-level rendering.
type Render2 = render.Render2

// SDF2 is a 2D signed distance function. *shape.Shape satisfies this interface.
type SDF2 interface {
	Evaluate(p v2.Vec) float64
	BoundingBox() v2.Box
}

// SDF3 is a 3D signed distance function. *solid.Solid satisfies this interface.
type SDF3 interface {
	Evaluate(p v3.Vec) float64
	BoundingBox() v3.Box
}

type sdf2Adapter struct{ s SDF2 }

func (a sdf2Adapter) Evaluate(p v2sdf.Vec) float64 { return a.s.Evaluate(v2.Vec(p)) }
func (a sdf2Adapter) BoundingBox() sdf.Box2        { return a.s.BoundingBox().SDF() }

type sdf3Adapter struct{ s SDF3 }

func (a sdf3Adapter) Evaluate(p v3sdf.Vec) float64 { return a.s.Evaluate(v3.Vec(p)) }
func (a sdf3Adapter) BoundingBox() sdf.Box3        { return a.s.BoundingBox().SDF() }

// PNG is a PNG drawing target.
type PNG struct {
	p *render.PNG
}

// DXF is a DXF drawing target.
type DXF struct {
	d *render.DXF
}

// NewPNG returns an empty PNG drawing target sized to bb with the given pixel dimensions.
func NewPNG(name string, bb v2.Box, pixels v2i.Vec) (*PNG, error) {
	p, err := render.NewPNG(name, bb.SDF(), v2isdf.Vec(pixels))
	if err != nil {
		return nil, err
	}
	return &PNG{p}, nil
}

// NewDXF returns an empty DXF drawing target.
func NewDXF(name string) *DXF {
	return &DXF{render.NewDXF(name)}
}

// RenderSDF2 renders a Shape to the PNG.
func (d *PNG) RenderSDF2(s SDF2) { d.p.RenderSDF2(sdf2Adapter{s}) }

// RenderSDF2MinMax renders a Shape with explicit min/max distance thresholds.
func (d *PNG) RenderSDF2MinMax(s SDF2, dmin, dmax float64) {
	d.p.RenderSDF2MinMax(sdf2Adapter{s}, dmin, dmax)
}

// Line draws a line from p0 to p1.
func (d *PNG) Line(p0, p1 v2.Vec) { d.p.Line(v2sdf.Vec(p0), v2sdf.Vec(p1)) }

// Lines draws a polyline through pts.
func (d *PNG) Lines(pts v2.VecSet) {
	out := make(v2sdf.VecSet, len(pts))
	for i, p := range pts {
		out[i] = v2sdf.Vec(p)
	}
	d.p.Lines(out)
}

// Triangle draws a triangle.
func (d *PNG) Triangle(t v2.Triangle2) { d.p.Triangle(t.SDF()) }

// Save writes the PNG to its file.
func (d *PNG) Save() error { return d.p.Save() }

// Box draws a bounding box.
func (d *DXF) Box(b v2.Box) {
	sb := b.SDF()
	d.d.Box(&sb)
}

// Line draws a single line.
func (d *DXF) Line(line *v2.Line2) {
	sl := line.SDF()
	d.d.Line(&sl)
}

// Lines draws a list of lines.
func (d *DXF) Lines(lines []*v2.Line2) {
	out := make([]*sdf.Line2, len(lines))
	for i, l := range lines {
		sl := l.SDF()
		out[i] = &sl
	}
	d.d.Lines(out)
}

// Points draws a set of points with the given radius.
func (d *DXF) Points(pts v2.VecSet, r float64) {
	out := make(v2sdf.VecSet, len(pts))
	for i, p := range pts {
		out[i] = v2sdf.Vec(p)
	}
	d.d.Points(out, r)
}

// Triangle draws a triangle.
func (d *DXF) Triangle(t v2.Triangle2) { d.d.Triangle(t.SDF()) }

// Save writes the DXF to its file.
func (d *DXF) Save() error { return d.d.Save() }

// NewMarchingCubesOctreeParallel returns the default parallel octree renderer for 3D.
func NewMarchingCubesOctreeParallel(meshCells int) Render3 {
	return render.NewMarchingCubesOctreeParallel(meshCells)
}

// NewMarchingSquaresQuadtree returns the default quadtree renderer for 2D.
func NewMarchingSquaresQuadtree(meshCells int) Render2 {
	return render.NewMarchingSquaresQuadtree(meshCells)
}

// NewDualContouring2D returns a dual-contouring 2D renderer.
func NewDualContouring2D(meshCells int) Render2 {
	return render.NewDualContouring2D(meshCells)
}

// To3MF renders an SDF3 to a 3MF file using the parallel octree renderer.
func To3MF(s SDF3, path string, meshCells int) {
	a := sdf3Adapter{s}
	r := render.NewMarchingCubesOctreeParallel(meshCells)
	fmt.Printf("rendering %s (%s)\n", path, r.Info(a))
	render.To3MF(a, path, r)
}

// ToDXF renders an SDF2 to a DXF file using the quadtree marching-squares renderer.
func ToDXF(s SDF2, path string, meshCells int) {
	a := sdf2Adapter{s}
	r := render.NewMarchingSquaresQuadtree(meshCells)
	fmt.Printf("rendering %s (%s)\n", path, r.Info(a))
	render.ToDXF(a, path, r)
}

// ToDXFWith renders an SDF2 to a DXF file using the given 2D renderer.
func ToDXFWith(s SDF2, path string, r Render2) {
	a := sdf2Adapter{s}
	fmt.Printf("rendering %s (%s)\n", path, r.Info(a))
	render.ToDXF(a, path, r)
}

// ToSVG renders an SDF2 to an SVG file using the quadtree marching-squares renderer.
func ToSVG(s SDF2, path string, meshCells int) {
	a := sdf2Adapter{s}
	r := render.NewMarchingSquaresQuadtree(meshCells)
	fmt.Printf("rendering %s (%s)\n", path, r.Info(a))
	render.ToSVG(a, path, r)
}

// ToSVGWith renders an SDF2 to an SVG file using the given 2D renderer.
func ToSVGWith(s SDF2, path string, r Render2) {
	a := sdf2Adapter{s}
	fmt.Printf("rendering %s (%s)\n", path, r.Info(a))
	render.ToSVG(a, path, r)
}

// ToPNG renders an SDF2 to a PNG file sized to the given bounding box and pixel dimensions.
func ToPNG(s SDF2, path string, bb v2.Box, width, height int) {
	p, err := render.NewPNG(path, bb.SDF(), v2isdf.Vec{X: width, Y: height})
	if err != nil {
		panic(err)
	}
	p.RenderSDF2(sdf2Adapter{s})
	if err := p.Save(); err != nil {
		panic(err)
	}
}

// SaveDXF writes a list of line segments to a DXF file.
func SaveDXF(path string, lines []*v2.Line2) {
	out := make([]*sdf.Line2, len(lines))
	for i, l := range lines {
		sl := l.SDF()
		out[i] = &sl
	}
	if err := render.SaveDXF(path, out); err != nil {
		panic(err)
	}
}

// SaveSVG writes a list of line segments to an SVG file.
func SaveSVG(path, lineStyle string, lines []*v2.Line2) {
	out := make([]*sdf.Line2, len(lines))
	for i, l := range lines {
		sl := l.SDF()
		out[i] = &sl
	}
	if err := render.SaveSVG(path, lineStyle, out); err != nil {
		panic(err)
	}
}

// Poly writes a polygon to a DXF file (one line per edge).
func Poly(p *sdf.Polygon, path string) {
	if err := render.Poly(p, path); err != nil {
		panic(err)
	}
}
