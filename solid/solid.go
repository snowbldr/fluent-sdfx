package solid

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/plane"
	flrender "github.com/snowbldr/fluent-sdfx/render"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/sdfx/render"
	"github.com/snowbldr/sdfx/sdf"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
	v3isdf "github.com/snowbldr/sdfx/vec/v3i"
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

// UnionAll combines multiple solids into one.
func UnionAll(solids ...*Solid) *Solid {
	sdf3s := make([]sdf.SDF3, len(solids))
	for i, s := range solids {
		sdf3s[i] = s.SDF3
	}
	return &Solid{sdf.Union3D(sdf3s...)}
}

// Bounds returns the solid's 3D axis-aligned bounding box as a fluent v3.Box.
// Use this instead of BoundingBox when you want the fluent box API
// (BoundingBox returns sdfx's raw sdf.Box3).
func (s *Solid) Bounds() Box3 {
	return wrapBox3(v3.FromSDF(s.SDF3.BoundingBox()))
}

// --- Transform methods ---

func (s *Solid) ZeroZ() *Solid {
	return &Solid{sdf.Transform3D(s.SDF3, sdf.Translate3d(v3sdf.Vec{Z: -s.SDF3.BoundingBox().Min.Z}))}
}

func (s *Solid) Translate(v v3.Vec) *Solid {
	return &Solid{sdf.Transform3D(s.SDF3, sdf.Translate3d(v3sdf.Vec(v)))}
}

func (s *Solid) TranslateX(x float64) *Solid         { return s.Translate(v3.X(x)) }
func (s *Solid) TranslateY(y float64) *Solid         { return s.Translate(v3.Y(y)) }
func (s *Solid) TranslateZ(z float64) *Solid         { return s.Translate(v3.Z(z)) }
func (s *Solid) TranslateXY(x, y float64) *Solid     { return s.Translate(v3.XY(x, y)) }
func (s *Solid) TranslateXZ(x, z float64) *Solid     { return s.Translate(v3.XZ(x, z)) }
func (s *Solid) TranslateYZ(y, z float64) *Solid     { return s.Translate(v3.YZ(y, z)) }
func (s *Solid) TranslateXYZ(x, y, z float64) *Solid { return s.Translate(v3.XYZ(x, y, z)) }

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

// RotateToVector rotates the solid so that the 'from' direction aligns with the 'to' direction.
func (s *Solid) RotateToVector(from, to v3.Vec) *Solid {
	return &Solid{sdf.Transform3D(s.SDF3, sdf.RotateToVector(v3sdf.Vec(from), v3sdf.Vec(to)))}
}

func (s *Solid) Scale(v v3.Vec) *Solid {
	return &Solid{sdf.Transform3D(s.SDF3, sdf.Scale3d(v3sdf.Vec(v)))}
}

// ScaleUniform scales uniformly on all axes. Unlike Scale, distance is preserved.
func (s *Solid) ScaleUniform(k float64) *Solid {
	return &Solid{sdf.ScaleUniform3D(s.SDF3, k)}
}

func (s *Solid) Transform(m M44) *Solid {
	return &Solid{sdf.Transform3D(s.SDF3, sdf.M44(m))}
}

// MirrorXY mirrors across the XY plane (negates Z).
func (s *Solid) MirrorXY() *Solid {
	return &Solid{sdf.Transform3D(s.SDF3, sdf.MirrorXY())}
}

// MirrorXZ mirrors across the XZ plane (negates Y).
func (s *Solid) MirrorXZ() *Solid {
	return &Solid{sdf.Transform3D(s.SDF3, sdf.MirrorXZ())}
}

// MirrorYZ mirrors across the YZ plane (negates X).
func (s *Solid) MirrorYZ() *Solid {
	return &Solid{sdf.Transform3D(s.SDF3, sdf.MirrorYZ())}
}

// MirrorXeqY mirrors across the X==Y plane (swaps X and Y).
func (s *Solid) MirrorXeqY() *Solid {
	return &Solid{sdf.Transform3D(s.SDF3, sdf.MirrorXeqY())}
}

// Center translates the solid so its bounding box center is at the origin.
func (s *Solid) Center() *Solid {
	bb := s.Bounds()
	center := bb.Center()
	return s.Translate(center.Neg())
}

// --- Boolean methods ---

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

func (s *Solid) Intersect(other ...*Solid) *Solid {
	sdf3s := make([]sdf.SDF3, len(other))
	for i, o := range other {
		sdf3s[i] = o.SDF3
	}
	return &Solid{sdf.Intersect3D(s.SDF3, sdf.Union3D(sdf3s...))}
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

// --- Smooth boolean methods ---

// SmoothUnion blends s with other solids using the given MinFunc.
func (s *Solid) SmoothUnion(min MinFunc, other ...*Solid) *Solid {
	all := make([]*Solid, 0, len(other)+1)
	all = append(all, s)
	all = append(all, other...)
	return SmoothUnion(min, all...)
}

// SmoothAdd is an alias for SmoothUnion.
func (s *Solid) SmoothAdd(min MinFunc, other ...*Solid) *Solid {
	return s.SmoothUnion(min, other...)
}

// SmoothCut subtracts the union of tools from s, blended with the given MaxFunc.
func (s *Solid) SmoothCut(max MaxFunc, tools ...*Solid) *Solid {
	return SmoothDifference(max, s, tools...)
}

// SmoothDifference is an alias for SmoothCut.
func (s *Solid) SmoothDifference(max MaxFunc, tools ...*Solid) *Solid {
	return SmoothDifference(max, s, tools...)
}

// SmoothIntersect intersects s with the union of others, blended with the given MaxFunc.
func (s *Solid) SmoothIntersect(max MaxFunc, other ...*Solid) *Solid {
	return SmoothIntersection(max, s, other...)
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

// CutPlane cuts the solid along a plane. The solid on the normal side remains.
func (s *Solid) CutPlane(point, normal v3.Vec) *Solid {
	return &Solid{sdf.Cut3D(s.SDF3, v3sdf.Vec(point), v3sdf.Vec(normal))}
}

// Split cuts the solid along a plane and returns both halves.
// The first result is the half on the plane's normal side; the second
// is the opposite side.
func (s *Solid) Split(p plane.Plane) (*Solid, *Solid) {
	return s.CutPlane(p.Origin, p.Normal), s.CutPlane(p.Origin, p.Normal.Neg())
}

// Elongate stretches the solid by the given amounts along each axis.
func (s *Solid) Elongate(h v3.Vec) *Solid {
	return &Solid{sdf.Elongate3D(s.SDF3, v3sdf.Vec(h))}
}

// Shell hollows out the solid, leaving a shell of the given thickness.
func (s *Solid) Shell(thickness float64) *Solid {
	return New(sdf.Shell3D(s.SDF3, thickness))
}

// Offset expands (positive) or contracts (negative) the solid by the given
// distance along its surface normal.
func (s *Solid) Offset(distance float64) *Solid {
	return Wrap(sdf.Offset3D(s.SDF3, distance))
}

// --- Pattern/array methods ---

// Array creates an XYZ grid array of the solid.
func (s *Solid) Array(numX, numY, numZ int, step v3.Vec) *Solid {
	return &Solid{sdf.Array3D(s.SDF3, v3isdf.Vec{X: numX, Y: numY, Z: numZ}, v3sdf.Vec(step))}
}

// SmoothArray creates an XYZ grid array using min for blending adjacent copies.
// Pair with PolyMin / RoundMin etc.
func (s *Solid) SmoothArray(numX, numY, numZ int, step v3.Vec, min sdf.MinFunc) *Solid {
	arr := sdf.Array3D(s.SDF3, v3isdf.Vec{X: numX, Y: numY, Z: numZ}, v3sdf.Vec(step))
	arr.(*sdf.ArraySDF3).SetMin(min)
	return &Solid{arr}
}

// RotateCopyZ creates N copies of the solid evenly spaced around the Z axis.
func (s *Solid) RotateCopyZ(n int) *Solid {
	return &Solid{sdf.RotateCopy3D(s.SDF3, n)}
}

// RotateUnionZ creates a union of the solid rotated N times by the given step matrix.
// Useful for creating patterns with custom rotation + translation per step.
func (s *Solid) RotateUnionZ(n int, step M44) *Solid {
	return &Solid{sdf.RotateUnion3D(s.SDF3, n, sdf.M44(step))}
}

// SmoothRotateUnionZ creates N rotated copies blended with min.
func (s *Solid) SmoothRotateUnionZ(n int, step M44, min sdf.MinFunc) *Solid {
	ru := sdf.RotateUnion3D(s.SDF3, n, sdf.M44(step))
	ru.(*sdf.RotateUnionSDF3).SetMin(min)
	return &Solid{ru}
}

// Multi creates a union of the solid at the given positions.
func (s *Solid) Multi(positions []v3.Vec) *Solid {
	return &Solid{sdf.Multi3D(s.SDF3, v3Slice(positions))}
}

// LineOf creates a union of the solid along a line from p0 to p1.
// The pattern string controls placement: 'x' places a copy, any other char skips.
func (s *Solid) LineOf(p0, p1 v3.Vec, pattern string) *Solid {
	return &Solid{sdf.LineOf3D(s.SDF3, v3sdf.Vec(p0), v3sdf.Vec(p1), pattern)}
}

// Orient creates a union of the solid oriented along each direction vector.
// base is the original orientation vector of the solid.
func (s *Solid) Orient(base v3.Vec, directions []v3.Vec) *Solid {
	return &Solid{sdf.Orient3D(s.SDF3, v3sdf.Vec(base), v3Slice(directions))}
}

// --- Sampling methods ---

// Voxel samples the solid into a voxel grid and returns an SDF3 that trilinearly interpolates it.
// meshCells is the resolution on the longest axis; progress receives 0-1 sampling progress (may be nil).
func (s *Solid) Voxel(meshCells int, progress chan float64) *Solid {
	return &Solid{sdf.NewVoxelSDF3(s.SDF3, meshCells, progress)}
}

// Benchmark reports the evaluation speed of the solid's SDF3.
func (s *Solid) Benchmark(description string) {
	sdf.BenchmarkSDF3(description, s.SDF3)
}

// Normal returns the surface normal at point p, computed via finite
// differences with sample step eps. p does not need to lie on the surface.
func (s *Solid) Normal(p v3.Vec, eps float64) v3.Vec {
	return v3.Vec(sdf.Normal3(s.SDF3, v3sdf.Vec(p), eps))
}

// Raycast sphere-traces a ray from `from` in direction `dir` until it hits
// the solid's surface. scaleAndSigmoid > 0 enables sigmoid scaling for
// non-Lipschitz SDFs, stepScale < 1 trades evaluations for precision,
// epsilon is the surface-hit threshold, maxDist caps the trace length, and
// maxSteps caps iteration count. Returns the collision point, the distance
// traveled (negative if no hit), and the number of steps taken.
func (s *Solid) Raycast(from, dir v3.Vec, scaleAndSigmoid, stepScale, epsilon, maxDist float64, maxSteps int) (v3.Vec, float64, int) {
	hit, t, steps := sdf.Raycast3(s.SDF3, v3sdf.Vec(from), v3sdf.Vec(dir), scaleAndSigmoid, stepScale, epsilon, maxDist, maxSteps)
	return v3.Vec(hit), t, steps
}

// --- Render methods ---

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
