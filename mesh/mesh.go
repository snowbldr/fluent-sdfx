// Package mesh provides triangle-mesh utilities and types that complement the
// SDF-based shape and solid packages.
package mesh

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/sdfx/render"
	"github.com/snowbldr/sdfx/sdf"
	v2sdf "github.com/snowbldr/sdfx/vec/v2"
)

// Triangle3 is a 3D triangle (three v3.Vec vertices).
type Triangle3 = v3.Triangle3

// Triangle2 is a 2D triangle (three v2.Vec vertices).
type Triangle2 = v2.Triangle2

// Line2 is a 2D line segment defined by two v2.Vec endpoints.
type Line2 = v2.Line2

// TriangleI is a triangle referencing three indices into a vertex list.
type TriangleI [3]int

// TriangleISet is a set of triangles defined by vertex indices.
type TriangleISet []TriangleI

// ToTriangle2 returns the triangle formed by indexing into pts.
func (t TriangleI) ToTriangle2(pts v2.VecSet) Triangle2 {
	return Triangle2{pts[t[0]], pts[t[1]], pts[t[2]]}
}

// Canonical rotates t so its smallest index is first.
func (t *TriangleI) Canonical() {
	r := render.TriangleI(*t)
	r.Canonical()
	*t = TriangleI(r)
}

// Canonical returns a copy of ts with each triangle in canonical form.
func (ts TriangleISet) Canonical() []TriangleI {
	r := make(render.TriangleISet, len(ts))
	for i, t := range ts {
		r[i] = render.TriangleI(t)
	}
	r = r.Canonical()
	out := make([]TriangleI, len(r))
	for i, t := range r {
		out[i] = TriangleI(t)
	}
	return out
}

// Equals reports whether two triangle sets are identical after canonicalization.
func (ts TriangleISet) Equals(s TriangleISet) bool {
	a := make(render.TriangleISet, len(ts))
	for i, t := range ts {
		a[i] = render.TriangleI(t)
	}
	b := make(render.TriangleISet, len(s))
	for i, t := range s {
		b[i] = render.TriangleI(t)
	}
	return a.Equals(b)
}

func fromSDFTriangleISet(ts render.TriangleISet) TriangleISet {
	out := make(TriangleISet, len(ts))
	for i, t := range ts {
		out[i] = TriangleI(t)
	}
	return out
}

func toSDFTris(tris []*Triangle3) []*sdf.Triangle3 {
	out := make([]*sdf.Triangle3, len(tris))
	for i, t := range tris {
		st := t.SDF()
		out[i] = &st
	}
	return out
}

func fromSDFTris(tris []*sdf.Triangle3) []*Triangle3 {
	out := make([]*Triangle3, len(tris))
	for i, t := range tris {
		mt := v3.FromSDFTriangle3(*t)
		out[i] = &mt
	}
	return out
}

func fromSDFTrisValue(tris []sdf.Triangle3) []Triangle3 {
	out := make([]Triangle3, len(tris))
	for i, t := range tris {
		out[i] = v3.FromSDFTriangle3(t)
	}
	return out
}

func toSDFLines(lines []*Line2) []*sdf.Line2 {
	out := make([]*sdf.Line2, len(lines))
	for i, l := range lines {
		sl := l.SDF()
		out[i] = &sl
	}
	return out
}

// SDFLines converts our []*Line2 to sdfx []*sdf.Line2 at package boundaries.
func SDFLines(lines []*Line2) []*sdf.Line2 { return toSDFLines(lines) }

// ToTriangles renders a solid to a triangle mesh using the given renderer.
func ToTriangles(s *solid.Solid, r render.Render3) []*Triangle3 {
	return fromSDFTris(render.ToTriangles(s.SDF3, r))
}

// CollectTriangles renders a solid to a triangle mesh using the given renderer.
// Unlike ToTriangles, this returns value-type triangles.
func CollectTriangles(s *solid.Solid, r render.Render3) []Triangle3 {
	return fromSDFTrisValue(render.CollectTriangles(s.SDF3, r))
}

// CountBoundaryEdges counts edges that only belong to one triangle — a non-zero
// count indicates a non-closed mesh.
func CountBoundaryEdges(tris []Triangle3) int {
	sdfTris := make([]sdf.Triangle3, len(tris))
	for i, t := range tris {
		sdfTris[i] = t.SDF()
	}
	return render.CountBoundaryEdges(sdfTris)
}

// IsWatertight reports whether the mesh is a closed surface: every edge is
// shared by exactly two triangles. Returns true and zero boundary edges for
// a sealed solid.
func IsWatertight(tris []Triangle3) (bool, int) {
	n := CountBoundaryEdges(tris)
	return n == 0, n
}

// SaveSTL writes a triangle mesh to an STL file.
func SaveSTL(path string, mesh []*Triangle3) {
	if err := render.SaveSTL(path, toSDFTris(mesh)); err != nil {
		panic(err)
	}
}

func v2Slice(pts []v2.Vec) []v2sdf.Vec {
	out := make([]v2sdf.Vec, len(pts))
	for i, p := range pts {
		out[i] = v2sdf.Vec(p)
	}
	return out
}

// Delaunay2d computes a Delaunay triangulation for a set of 2D points.
func Delaunay2d(vs v2.VecSet) TriangleISet {
	t, err := Delaunay2dE(vs)
	if err != nil {
		panic(err)
	}
	return t
}

// Delaunay2dE computes a Delaunay triangulation for a set of 2D points,
// returning an error instead of panicking.
func Delaunay2dE(vs v2.VecSet) (TriangleISet, error) {
	t, err := render.Delaunay2d(v2Slice(vs))
	if err != nil {
		return nil, err
	}
	return fromSDFTriangleISet(t), nil
}

// Delaunay2dSlow is a naive-algorithm reference implementation for testing.
func Delaunay2dSlow(vs v2.VecSet) TriangleISet {
	t, err := render.Delaunay2dSlow(v2Slice(vs))
	if err != nil {
		panic(err)
	}
	return fromSDFTriangleISet(t)
}

// VertexToLine connects a sequence of vertices with line segments.
// If closed is true, the last vertex connects back to the first.
func VertexToLine(vs []v2.Vec, closed bool) []*Line2 {
	sdfLines := sdf.VertexToLine(v2Slice(vs), closed)
	out := make([]*Line2, len(sdfLines))
	for i, l := range sdfLines {
		ml := v2.FromSDFLine2(*l)
		out[i] = &ml
	}
	return out
}
