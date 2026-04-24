package solid

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/plane"
	flrender "github.com/snowbldr/fluent-sdfx/render"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/sdfx/render"
	"github.com/snowbldr/sdfx/sdf"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
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

// Bounds returns the solid's 3D axis-aligned bounding box as a fluent v3.Box.
// Use this instead of BoundingBox when you want the fluent box API
// (BoundingBox returns sdfx's raw sdf.Box3).
func (s *Solid) Bounds() Box3 {
	return v3.FromSDF(s.SDF3.BoundingBox())
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

// Extrude linearly extrudes a 2D profile to a solid of the given height.
//
// The profile argument accepts any sdf.SDF2 — notably *shape.Shape, which
// embeds sdf.SDF2. This keeps the package free of a back-reference to
// shape so the shape package can import solid and attach fluent methods.
func Extrude(profile sdf.SDF2, height float64) *Solid {
	return &Solid{sdf.Extrude3D(profile, height)}
}

// Slice cuts a planar cross-section through a solid and returns it as a
// raw sdf.SDF2. Wrap with shape.Of for the full *shape.Shape fluent API,
// or use shape.Slice / solid.SliceShape helpers.
func Slice(s *Solid, origin, normal v3.Vec) sdf.SDF2 {
	return sdf.Slice2D(s.SDF3, v3sdf.Vec(origin), v3sdf.Vec(normal))
}

func TwistExtrude(profile sdf.SDF2, height, twist float64) *Solid {
	return &Solid{sdf.TwistExtrude3D(profile, height, twist)}
}

func Screw(profile sdf.SDF2, height, start, pitch float64, num int) *Solid {
	return New(sdf.Screw3D(profile, height, start, pitch, num))
}

// --- Transform methods ---

func (s *Solid) ZeroZ() *Solid {
	return &Solid{sdf.Transform3D(s.SDF3, sdf.Translate3d(v3sdf.Vec{Z: -s.SDF3.BoundingBox().Min.Z}))}
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

// Add is an alias for Union.
func (s *Solid) Add(other ...*Solid) *Solid {
	return s.Union(other...)
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

// Difference is an alias for Cut.
func (s *Solid) Difference(other ...*Solid) *Solid {
	return s.Cut(other...)
}

// --- 3D → 2D cross-section methods ---

// Slice2D cuts a planar cross-section through the solid and returns it as a
// raw sdf.SDF2. Wrap with shape.Wrap2D (or use shape.SliceOf) for the fluent
// *shape.Shape API.
func (s *Solid) Slice2D(origin, normal v3.Vec) sdf.SDF2 {
	return sdf.Slice2D(s.SDF3, v3sdf.Vec(origin), v3sdf.Vec(normal))
}

// SliceAt cuts a cross-section at the given plane. Wrap with shape.Wrap2D
// (or use shape.SliceAt) for the fluent *shape.Shape API.
func (s *Solid) SliceAt(p plane.Plane) sdf.SDF2 {
	return s.Slice2D(p.Origin, p.Normal)
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

// MinCells is the floor applied by CellsFor so sub-mm parts (e.g. a 1 mm
// sphere at 3 cells/mm, which would otherwise give 3 cells and render as
// an empty mesh) still produce a recognizable shape. Raise it for more
// sub-mm detail, lower it (or set to 1) for raw density behavior.
var MinCells = 32

// STL renders the solid to an STL file using the parallel octree renderer.
//
// cellsPerMM is a mesh density — cells per millimeter along the longest
// bounding-box axis — so the same value gives proportional detail across
// parts of different sizes. Examples:
//
//	500 mm enclosure, preview        → 0.2   (100 cells)
//	500 mm enclosure, final          → 2.0   (1000 cells)
//	50 mm bracket, typical           → 5.0   (250 cells)
//	10 mm gear, detailed             → 20.0  (200 cells)
//	1 mm sphere                      → 12.0  (floored to MinCells=32)
//
// Render time scales roughly with cells³, so halving cellsPerMM is ~8×
// faster — drop it low for iteration, crank it for final output. Tiny
// parts are floored at MinCells cells (see MinCells) so they don't
// render as empty.
//
// Optional decimate (0-1) is the fraction of triangles to remove:
// 0.1 removes 10% (keeps 90%); 0.9 removes 90% (keeps 10%). 0 disables decimation.
func (s *Solid) STL(path string, cellsPerMM float64, decimate ...float64) {
	flrender.ToSTL(s, path, render.NewMarchingCubesOctreeParallel(CellsFor(s, cellsPerMM)), decimate...)
}

// ThreeMF renders the solid to a 3MF file using the parallel octree renderer.
//
// cellsPerMM is a mesh density — cells per millimeter along the longest
// bounding-box axis — so the same value gives proportional detail across
// parts of different sizes. See STL for a table of typical values and
// notes on the MinCells floor.
func (s *Solid) ThreeMF(path string, cellsPerMM float64) {
	flrender.To3MF(s, path, CellsFor(s, cellsPerMM))
}

// MF3 is an alias for ThreeMF, for users who expect to autocomplete from "3".
func (s *Solid) MF3(path string, cellsPerMM float64) {
	s.ThreeMF(path, cellsPerMM)
}

// CellsFor returns the cell count along the longest bounding-box axis
// for a given mesh density, floored at MinCells.
//
// The result is ceil(longestAxis_mm * cellsPerMM), clamped up to MinCells
// so sub-mm parts don't collapse to zero or single-cell renders that
// marching cubes emits as empty meshes.
func CellsFor(s *Solid, cellsPerMM float64) int {
	size := s.SDF3.BoundingBox().Size()
	longest := math.Max(math.Max(size.X, size.Y), size.Z)
	cells := int(math.Ceil(longest * cellsPerMM))
	if cells < MinCells {
		cells = MinCells
	}
	return cells
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
