package solid

import (
	"math"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v3sdf "github.com/deadsy/sdfx/vec/v3"
	flrender "github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/shape"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

type Solid struct {
	sdf.SDF3
}

func v3Slice(pts []v3.Vec) []v3sdf.Vec {
	out := make([]v3sdf.Vec, len(pts))
	for i, p := range pts {
		out[i] = v3sdf.Vec(p)
	}
	return out
}

// New wraps a raw SDF3 (with optional error) into a Solid.
func New(sdf sdf.SDF3, err error) *Solid {
	if err != nil {
		panic(err)
	}
	return &Solid{sdf}
}

// Wrap wraps a raw SDF3 into a Solid (no error).
func Wrap(sdf sdf.SDF3) *Solid {
	return &Solid{sdf}
}

// BoundingBox returns the solid's 3D axis-aligned bounding box.
func (s *Solid) BoundingBox() Box3 {
	return v3.FromSDF(s.SDF3.BoundingBox())
}

// Evaluate returns the signed distance from p to the solid's surface.
func (s *Solid) Evaluate(p v3.Vec) float64 {
	return s.SDF3.Evaluate(v3sdf.Vec(p))
}

// --- Constructors ---

func Cylinder(height, radius, round float64) *Solid {
	return New(sdf.Cylinder3D(height, radius, round))
}

func Box(size v3.Vec, round float64) *Solid {
	return New(sdf.Box3D(v3sdf.Vec(size), round))
}

func Sphere(radius float64) *Solid {
	return New(sdf.Sphere3D(radius))
}

func Cone(height, r0, r1, round float64) *Solid {
	return New(sdf.Cone3D(height, r0, r1, round))
}

func Extrude(profile *shape.Shape, height float64) *Solid {
	return &Solid{sdf.Extrude3D(profile.SDF2, height)}
}

// Slice cuts a cross-section through a solid, returning a 2D shape.
func Slice(s *Solid, origin, dir v3.Vec) *shape.Shape {
	return shape.Wrap2D(sdf.Slice2D(s.SDF3, v3sdf.Vec(origin), v3sdf.Vec(dir)))
}

func TwistExtrude(profile *shape.Shape, height, twist float64) *Solid {
	return &Solid{sdf.TwistExtrude3D(profile.SDF2, height, twist)}
}

func Screw(profile *shape.Shape, height, start, pitch float64, num int) *Solid {
	return New(sdf.Screw3D(profile.SDF2, height, start, pitch, num))
}

// --- Transform methods ---

func (s *Solid) ZeroZ() *Solid {
	return &Solid{sdf.Transform3D(s.SDF3, sdf.Translate3d(v3sdf.Vec{Z: -s.BoundingBox().Min.Z}))}
}

func (s *Solid) Translate(v v3.Vec) *Solid {
	return &Solid{sdf.Transform3D(s.SDF3, sdf.Translate3d(v3sdf.Vec(v)))}
}

func (s *Solid) RotateX(angleDeg float64) *Solid {
	return &Solid{sdf.Transform3D(s.SDF3, sdf.RotateX(angleDeg*math.Pi/180))}
}

func (s *Solid) RotateY(angleDeg float64) *Solid {
	return &Solid{sdf.Transform3D(s.SDF3, sdf.RotateY(angleDeg*math.Pi/180))}
}

func (s *Solid) RotateZ(angleDeg float64) *Solid {
	return &Solid{sdf.Transform3D(s.SDF3, sdf.RotateZ(angleDeg*math.Pi/180))}
}

func (s *Solid) RotateAxis(axis v3.Vec, angleDeg float64) *Solid {
	return &Solid{sdf.Transform3D(s.SDF3, sdf.Rotate3d(v3sdf.Vec(axis), angleDeg*math.Pi/180))}
}

func (s *Solid) Scale(v v3.Vec) *Solid {
	return &Solid{sdf.Transform3D(s.SDF3, sdf.Scale3d(v3sdf.Vec(v)))}
}

func (s *Solid) Transform(m M44) *Solid {
	return &Solid{sdf.Transform3D(s.SDF3, sdf.M44(m))}
}

// --- Boolean methods ---

// UnionAll combines multiple solids into one.
func UnionAll(solids ...*Solid) *Solid {
	sdf3s := make([]sdf.SDF3, len(solids))
	for i, s := range solids {
		sdf3s[i] = s.SDF3
	}
	return &Solid{sdf.Union3D(sdf3s...)}
}

func (s *Solid) Union(other ...*Solid) *Solid {
	sdf3s := make([]sdf.SDF3, len(other)+1)
	sdf3s[0] = s.SDF3
	for i, o := range other {
		sdf3s[i+1] = o.SDF3
	}
	return &Solid{sdf.Union3D(sdf3s...)}
}

func (s *Solid) Intersect(other *Solid) *Solid {
	return &Solid{sdf.Intersect3D(s.SDF3, other.SDF3)}
}

func (s *Solid) Cut(other ...*Solid) *Solid {
	sdf3s := make([]sdf.SDF3, len(other))
	for i, o := range other {
		sdf3s[i] = o.SDF3
	}
	tool := sdf.Union3D(sdf3s...)
	return &Solid{sdf.Difference3D(s.SDF3, tool)}
}

// --- Modification methods ---

// Shrink returns a new solid inset by the given amount on all surfaces.
func (s *Solid) Shrink(amount float64) *Solid {
	return &Solid{&offsetSDF3{s.SDF3, amount}}
}

// Grow returns a new solid expanded by the given amount on all surfaces.
func (s *Solid) Grow(amount float64) *Solid {
	return &Solid{&offsetSDF3{s.SDF3, -amount}}
}

// Correct scales an SDF's distance values without changing the shape.
// Factor < 1 makes the octree renderer explore more cells, fixing holes
// caused by operations like TwistExtrude3D that overestimate distances.
func (s *Solid) Correct(factor float64) *Solid {
	return &Solid{&correctedSDF3{s.SDF3, factor}}
}

// ToSTL renders the solid to an STL file using the parallel octree renderer.
// meshCells controls resolution (number of cells on the longest axis).
// Optional factor (0-1) controls mesh decimation: 0.5 = keep 50% of triangles.
func (s *Solid) ToSTL(path string, meshCells int, factor ...float64) {
	flrender.ToSTL(s, path, render.NewMarchingCubesOctreeParallel(meshCells), factor...)
}

// To3MF renders the solid to a 3MF file using the parallel octree renderer.
func (s *Solid) To3MF(path string, meshCells int) {
	flrender.To3MF(s, path, meshCells)
}

type offsetSDF3 struct {
	sdf    sdf.SDF3
	offset float64
}

func (o *offsetSDF3) Evaluate(p v3sdf.Vec) float64 {
	return o.sdf.Evaluate(p) + o.offset
}

func (o *offsetSDF3) BoundingBox() sdf.Box3 {
	bb := o.sdf.BoundingBox()
	d := v3sdf.Vec{X: o.offset, Y: o.offset, Z: o.offset}
	return sdf.Box3{Min: bb.Min.Add(d), Max: bb.Max.Sub(d)}
}

type correctedSDF3 struct {
	s      sdf.SDF3
	factor float64
}

func (c *correctedSDF3) Evaluate(p v3sdf.Vec) float64 {
	return c.s.Evaluate(p) * c.factor
}

func (c *correctedSDF3) BoundingBox() sdf.Box3 {
	return c.s.BoundingBox()
}
