package shape

import (
	"github.com/deadsy/sdfx/sdf"
	v2sdf "github.com/deadsy/sdfx/vec/v2"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

// Poly is a fluent wrapper over sdf.Polygon that supports smoothed,
// chamfered, arc-filleted, and polar-relative vertices.
type Poly struct {
	p *sdf.Polygon
}

// NewPoly starts a new polygon builder.
func NewPoly() *Poly {
	return &Poly{p: sdf.NewPolygon()}
}

// Add appends a cartesian vertex and returns a handle for chained modifiers.
func (p *Poly) Add(x, y float64) *PolyVertex {
	return &PolyVertex{v: p.p.Add(x, y)}
}

// AddV2 appends a cartesian vertex from a v2.Vec.
func (p *Poly) AddV2(v v2.Vec) *PolyVertex {
	return &PolyVertex{v: p.p.AddV2(v2sdf.Vec(v))}
}

// AddV2Set appends a sequence of cartesian vertices.
func (p *Poly) AddV2Set(pts []v2.Vec) *Poly {
	p.p.AddV2Set(v2Slice(pts))
	return p
}

// Drop removes the last vertex.
func (p *Poly) Drop() *Poly {
	p.p.Drop()
	return p
}

// Close marks the polygon as closed (open polygons have their start/end joined automatically for SDF use).
func (p *Poly) Close() *Poly {
	p.p.Close()
	return p
}

// Reverse reverses the vertex order.
func (p *Poly) Reverse() *Poly {
	p.p.Reverse()
	return p
}

// Vertices returns the resolved cartesian vertices (after smoothing, arcs, etc).
func (p *Poly) Vertices() []v2.Vec {
	vs := p.p.Vertices()
	out := make([]v2.Vec, len(vs))
	for i, v := range vs {
		out[i] = v2.Vec(v)
	}
	return out
}

// Raw exposes the underlying sdf.Polygon for advanced use.
func (p *Poly) Raw() *sdf.Polygon {
	return p.p
}

// Mesh2D materializes the polygon as a Shape (SDF2).
func (p *Poly) Mesh2D() *Shape {
	s, err := p.p.Mesh2D()
	if err != nil {
		panic(err)
	}
	return &Shape{s}
}

// Build returns a Shape using the cache-friendly FlatPolygon2D backend.
func (p *Poly) Build() *Shape {
	return &Shape{mustFlat(sdf.FlatPolygon2D(p.p.Vertices()))}
}

func mustFlat(s sdf.SDF2, err error) sdf.SDF2 {
	if err != nil {
		panic(err)
	}
	return s
}

// PolyVertex wraps sdf.PolygonVertex with fluent modifiers that return the vertex for chaining.
type PolyVertex struct {
	v *sdf.PolygonVertex
}

// Rel marks this vertex as relative to the previous one.
func (v *PolyVertex) Rel() *PolyVertex {
	v.v.Rel()
	return v
}

// Polar marks this vertex as polar (x=radius, y=angleRadians). Note: sdfx uses radians here.
func (v *PolyVertex) Polar() *PolyVertex {
	v.v.Polar()
	return v
}

// Smooth rounds the corner at this vertex with the given radius and number of facets.
func (v *PolyVertex) Smooth(radius float64, facets int) *PolyVertex {
	v.v.Smooth(radius, facets)
	return v
}

// Chamfer cuts the corner at this vertex with a 45° chamfer of the given size.
func (v *PolyVertex) Chamfer(size float64) *PolyVertex {
	v.v.Chamfer(size)
	return v
}

// Arc replaces the straight edge from the previous vertex with an arc of the given radius.
// Positive radius curves outward (convex); negative curves inward (concave).
func (v *PolyVertex) Arc(radius float64, facets int) *PolyVertex {
	v.v.Arc(radius, facets)
	return v
}
