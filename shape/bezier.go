package shape

import (
	"math"

	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	"github.com/snowbldr/sdfx/sdf"
	v2sdf "github.com/snowbldr/sdfx/vec/v2"
)

// Bezier is a fluent wrapper over sdf.Bezier for building bezier curves.
type Bezier struct {
	b *sdf.Bezier
}

// NewBezier starts a new bezier curve builder.
func NewBezier() *Bezier {
	return &Bezier{b: sdf.NewBezier()}
}

// Add appends an endpoint vertex at (x,y).
func (b *Bezier) Add(x, y float64) *BezierVertex {
	return &BezierVertex{v: b.b.Add(x, y)}
}

// AddV2 appends an endpoint vertex at v.
func (b *Bezier) AddV2(v v2.Vec) *BezierVertex {
	return &BezierVertex{v: b.b.AddV2(v2sdf.Vec(v))}
}

// Close marks the curve as closed (end connects to start).
func (b *Bezier) Close() *Bezier {
	b.b.Close()
	return b
}

// Polygon converts the bezier curve into a (sampled) sdf.Polygon.
func (b *Bezier) Polygon() *sdf.Polygon {
	p, err := b.b.Polygon()
	if err != nil {
		panic(err)
	}
	return p
}

// Build returns a Shape from the bezier curve using a polygon SDF.
func (b *Bezier) Build() *Shape {
	s, err := b.b.Mesh2D()
	if err != nil {
		panic(err)
	}
	return &Shape{s}
}

// Vertices returns the resolved vertices of the bezier curve (after sampling).
func (b *Bezier) Vertices() []v2.Vec {
	vs := b.Polygon().Vertices()
	out := make([]v2.Vec, len(vs))
	for i, v := range vs {
		out[i] = v2.Vec(v)
	}
	return out
}

// BezierVertex wraps sdf.BezierVertex with fluent, degree-based handle modifiers.
type BezierVertex struct {
	v *sdf.BezierVertex
}

// Mid marks this vertex as a midpoint (control point between endpoints).
func (v *BezierVertex) Mid() *BezierVertex {
	v.v.Mid()
	return v
}

// HandleFwd sets the forward slope handle (thetaDeg in degrees, r in length).
func (v *BezierVertex) HandleFwd(thetaDeg, r float64) *BezierVertex {
	v.v.HandleFwd(thetaDeg*math.Pi/180, r)
	return v
}

// HandleRev sets the reverse slope handle (thetaDeg in degrees, r in length).
func (v *BezierVertex) HandleRev(thetaDeg, r float64) *BezierVertex {
	v.v.HandleRev(thetaDeg*math.Pi/180, r)
	return v
}

// Handle sets both slope handles: forward at theta + r_fwd, reverse at (theta+180) + r_rev.
func (v *BezierVertex) Handle(thetaDeg, fwd, rev float64) *BezierVertex {
	v.v.Handle(thetaDeg*math.Pi/180, fwd, rev)
	return v
}
