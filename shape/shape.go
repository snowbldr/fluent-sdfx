package shape

import (
	"fmt"
	"math"

	"github.com/snowbldr/fluent-sdfx/plane"
	flrender "github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	"github.com/snowbldr/fluent-sdfx/vec/v2i"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/sdfx/sdf"
	v2sdf "github.com/snowbldr/sdfx/vec/v2"
	v2isdf "github.com/snowbldr/sdfx/vec/v2i"
)

type Shape struct {
	sdf.SDF2
}

// Wrap wraps a raw SDF2 into a Shape.
func Wrap2D(s sdf.SDF2) *Shape {
	return &Shape{s}
}

func v2Slice(pts []v2.Vec) []v2sdf.Vec {
	out := make([]v2sdf.Vec, len(pts))
	for i, p := range pts {
		out[i] = v2sdf.Vec(p)
	}
	return out
}

// --- Constructors ---

// Rect returns an axis-aligned rectangle of the given size, centered at the
// origin. round adds a fillet of that radius on every corner; pass 0 for
// sharp corners. Panics if any size component is negative.
func Rect(size v2.Vec, round float64) *Shape {
	return &Shape{sdf.Box2D(v2sdf.Vec(size), round)}
}

// Circle returns a circle of the given radius, centered at the origin.
// Panics on negative radius.
func Circle(radius float64) *Shape {
	s, _ := sdf.Circle2D(radius)
	return &Shape{s}
}

// Polygon returns a flat (sharp-corner) polygon from the given vertex list.
// Vertices are in 2D (XY) and should be ordered consistently (clockwise or
// counter-clockwise — either works as long as it's consistent). Panics if
// the polygon is degenerate (fewer than 3 distinct vertices) or self-
// intersecting.
func Polygon(pts []v2.Vec) *Shape {
	s, err := sdf.FlatPolygon2D(v2Slice(pts))
	if err != nil {
		panic(err)
	}
	return &Shape{s}
}

// --- Cross-section helpers (3D → 2D) ---

// SliceOf cuts a planar cross-section through s and wraps the result as a
// fluent *Shape, ready for ToDXF/ToSVG/ToPNG and further 2D operations.
func SliceOf(s *solid.Solid, origin, normal v3.Vec) *Shape {
	return Wrap2D(s.Slice2D(origin, normal))
}

// SliceAt cuts a cross-section at the given plane and wraps as *Shape.
func SliceAt(s *solid.Solid, p plane.Plane) *Shape {
	return Wrap2D(s.SliceAt(p))
}

// Bounds returns the shape's 2D axis-aligned bounding box as a fluent v2.Box.
// Use this instead of BoundingBox when you want the fluent box API
// (BoundingBox returns sdfx's raw sdf.Box2).
func (s *Shape) Bounds() Box2 {
	return v2.FromSDF(s.SDF2.BoundingBox())
}

// --- Transform methods ---

func (s *Shape) Translate(v v2.Vec) *Shape {
	return &Shape{sdf.Transform2D(s.SDF2, sdf.Translate2d(v2sdf.Vec(v)))}
}

func (s *Shape) TranslateX(x float64) *Shape     { return s.Translate(v2.X(x)) }
func (s *Shape) TranslateY(y float64) *Shape     { return s.Translate(v2.Y(y)) }
func (s *Shape) TranslateXY(x, y float64) *Shape { return s.Translate(v2.XY(x, y)) }

func (s *Shape) Rotate(angleDeg float64) *Shape {
	return &Shape{sdf.Transform2D(s.SDF2, sdf.Rotate2d(angleDeg*math.Pi/180))}
}

func (s *Shape) Scale(v v2.Vec) *Shape {
	return &Shape{sdf.Transform2D(s.SDF2, sdf.Scale2d(v2sdf.Vec(v)))}
}

// ScaleUniform scales uniformly on both axes. Unlike Scale, distance is preserved.
func (s *Shape) ScaleUniform(k float64) *Shape {
	return &Shape{sdf.ScaleUniform2D(s.SDF2, k)}
}

func (s *Shape) MirrorX() *Shape {
	return &Shape{sdf.Transform2D(s.SDF2, sdf.MirrorX())}
}

func (s *Shape) MirrorY() *Shape {
	return &Shape{sdf.Transform2D(s.SDF2, sdf.MirrorY())}
}

func (s *Shape) Transform(m M33) *Shape {
	return &Shape{sdf.Transform2D(s.SDF2, sdf.M33(m))}
}

func (s *Shape) Offset(amount float64) *Shape {
	return &Shape{sdf.Offset2D(s.SDF2, amount)}
}

// Center translates the shape so its bounding box center is at the origin.
func (s *Shape) Center() *Shape {
	return &Shape{sdf.Center2D(s.SDF2)}
}

// CenterAndScale centers on the bounding box then scales uniformly. Distance is preserved.
func (s *Shape) CenterAndScale(k float64) *Shape {
	return &Shape{sdf.CenterAndScale2D(s.SDF2, k)}
}

// --- Boolean methods ---

func (s *Shape) Union(other ...*Shape) *Shape {
	sdf2s := make([]sdf.SDF2, len(other)+1)
	sdf2s[0] = s.SDF2
	for i, o := range other {
		sdf2s[i+1] = o.SDF2
	}
	return &Shape{sdf.Union2D(sdf2s...)}
}

// Add is an alias for Union.
func (s *Shape) Add(other ...*Shape) *Shape {
	return s.Union(other...)
}

func (s *Shape) Cut(other ...*Shape) *Shape {
	sdf2s := make([]sdf.SDF2, len(other))
	for i, o := range other {
		sdf2s[i] = o.SDF2
	}
	return &Shape{sdf.Difference2D(s.SDF2, sdf.Union2D(sdf2s...))}
}

// Difference is an alias for Cut.
func (s *Shape) Difference(other ...*Shape) *Shape {
	return s.Cut(other...)
}

func (s *Shape) Intersect(other ...*Shape) *Shape {
	sdf2s := make([]sdf.SDF2, len(other))
	for i, o := range other {
		sdf2s[i] = o.SDF2
	}
	return &Shape{sdf.Intersect2D(s.SDF2, sdf.Union2D(sdf2s...))}
}

// --- Smooth boolean methods ---

// SmoothUnion blends s with other shapes using a smooth min function.
func (s *Shape) SmoothUnion(min sdf.MinFunc, other ...*Shape) *Shape {
	all := make([]*Shape, 0, len(other)+1)
	all = append(all, s)
	all = append(all, other...)
	return SmoothUnion(min, all...)
}

// SmoothAdd is an alias for SmoothUnion.
func (s *Shape) SmoothAdd(min sdf.MinFunc, other ...*Shape) *Shape {
	return s.SmoothUnion(min, other...)
}

// SmoothCut subtracts the union of tools from s, blended with a smooth max function.
func (s *Shape) SmoothCut(max sdf.MaxFunc, tools ...*Shape) *Shape {
	return SmoothCut(max, s, tools...)
}

// SmoothDifference is an alias for SmoothCut.
func (s *Shape) SmoothDifference(max sdf.MaxFunc, tools ...*Shape) *Shape {
	return SmoothCut(max, s, tools...)
}

// SmoothIntersect intersects s with the union of others, blended with a smooth max function.
func (s *Shape) SmoothIntersect(max sdf.MaxFunc, other ...*Shape) *Shape {
	return SmoothIntersect(max, s, other...)
}

// --- Modification methods ---

// CutLine cuts the shape along a line from point a in direction v.
// The shape to the right of the line remains.
func (s *Shape) CutLine(a, dir v2.Vec) *Shape {
	return &Shape{sdf.Cut2D(s.SDF2, v2sdf.Vec(a), v2sdf.Vec(dir))}
}

// Split cuts the shape along a line through point a in direction dir
// and returns both halves. The first result is the half to the right
// of the directed line; the second is the half to the left.
func (s *Shape) Split(a, dir v2.Vec) (*Shape, *Shape) {
	return s.CutLine(a, dir), s.CutLine(a, dir.Neg())
}

// Elongate stretches the shape by the given amounts along each axis.
func (s *Shape) Elongate(h v2.Vec) *Shape {
	return &Shape{sdf.Elongate2D(s.SDF2, v2sdf.Vec(h))}
}

// Cache wraps the shape in a distance-value cache, trading memory for
// faster repeated evaluation (useful for slow-to-evaluate SDFs like text or meshes).
func (s *Shape) Cache() *Shape {
	return &Shape{sdf.Cache2D(s.SDF2)}
}

// --- Pattern/array methods ---

// Array creates an XY grid array of the shape.
func (s *Shape) Array(numX, numY int, step v2.Vec) *Shape {
	return &Shape{sdf.Array2D(s.SDF2, v2isdf.Vec{X: numX, Y: numY}, v2sdf.Vec(step))}
}

// SmoothArray creates an XY grid array using min for blending adjacent copies.
// Pair with solid.PolyMin / solid.RoundMin etc.
func (s *Shape) SmoothArray(numX, numY int, step v2.Vec, min sdf.MinFunc) *Shape {
	arr := sdf.Array2D(s.SDF2, v2isdf.Vec{X: numX, Y: numY}, v2sdf.Vec(step))
	arr.(*sdf.ArraySDF2).SetMin(min)
	return &Shape{arr}
}

// RotateCopy creates N copies of the shape evenly spaced in a full circle.
func (s *Shape) RotateCopy(n int) *Shape {
	return &Shape{sdf.RotateCopy2D(s.SDF2, n)}
}

// RotateUnion creates a union of the shape rotated N times by the given step matrix.
func (s *Shape) RotateUnion(n int, step M33) *Shape {
	return &Shape{sdf.RotateUnion2D(s.SDF2, n, sdf.M33(step))}
}

// SmoothRotateUnion creates N rotated copies blended with min.
func (s *Shape) SmoothRotateUnion(n int, step M33, min sdf.MinFunc) *Shape {
	ru := sdf.RotateUnion2D(s.SDF2, n, sdf.M33(step))
	ru.(*sdf.RotateUnionSDF2).SetMin(min)
	return &Shape{ru}
}

// Multi creates a union of the shape at the given positions. Variadic so
// you can write `hole.Multi(v2.XY(5, 0), v2.XY(-5, 0))` directly; pass a
// slice with `hole.Multi(positions...)`.
func (s *Shape) Multi(positions ...v2.Vec) *Shape {
	return &Shape{sdf.Multi2D(s.SDF2, v2Slice(positions))}
}

// LineOf creates a union of the shape along a line from p0 to p1.
// The pattern string controls placement: 'x' places a copy, any other char skips.
func (s *Shape) LineOf(p0, p1 v2.Vec, pattern string) *Shape {
	return &Shape{sdf.LineOf2D(s.SDF2, v2sdf.Vec(p0), v2sdf.Vec(p1), pattern)}
}

// --- 2D → 3D methods ---

// Extrude linearly extrudes the shape to a solid of the given height.
func (s *Shape) Extrude(height float64) *solid.Solid {
	return solid.Extrude(s.SDF2, height)
}

// ExtrudeRounded extrudes the shape to a solid with rounded top and bottom edges.
func (s *Shape) ExtrudeRounded(height, round float64) *solid.Solid {
	return solid.ExtrudeRounded(s.SDF2, height, round)
}

// TwistExtrude extrudes while rotating the profile about the Z axis over the height.
// twist is the total rotation in radians.
func (s *Shape) TwistExtrude(height, twist float64) *solid.Solid {
	return solid.TwistExtrude(s.SDF2, height, twist)
}

// ScaleExtrude extrudes while scaling the profile linearly over the height.
func (s *Shape) ScaleExtrude(height float64, scale v2.Vec) *solid.Solid {
	return solid.ScaleExtrude(s.SDF2, height, scale)
}

// ScaleTwistExtrude extrudes while scaling and twisting (radians) the profile over the height.
func (s *Shape) ScaleTwistExtrude(height, twist float64, scale v2.Vec) *solid.Solid {
	return solid.ScaleTwistExtrude(s.SDF2, height, twist, scale)
}

// Revolve rotates the shape around the Y axis to form a solid of revolution.
func (s *Shape) Revolve() *solid.Solid {
	return solid.Revolve(s.SDF2)
}

// RevolveAngle creates a partial revolution sweeping angleDeg degrees around the Y axis.
func (s *Shape) RevolveAngle(angleDeg float64) *solid.Solid {
	return solid.RevolveAngle(s.SDF2, angleDeg)
}

// Screw sweeps the shape helically to produce a threaded solid.
// height is the axial length, start is the starting z, pitch is the axial
// distance per turn, and num is the number of starts (parallel helices).
func (s *Shape) Screw(height, start, pitch float64, num int) *solid.Solid {
	return solid.Screw(s.SDF2, height, start, pitch, num)
}

// SweepHelix sweeps the shape along a helix of the given radius.
// See solid.SweepHelix for the flatEnds parameter semantics.
func (s *Shape) SweepHelix(radius, turns, height float64, flatEnds bool) *solid.Solid {
	return solid.SweepHelix(s.SDF2, radius, turns, height, flatEnds)
}

// LoftTo transitions from this shape (bottom) to top over the given height with optional rounding.
func (s *Shape) LoftTo(top *Shape, height, round float64) *solid.Solid {
	return solid.Loft(s.SDF2, top.SDF2, height, round)
}

// --- Mesh introspection ---

// MeshBoxes returns the acceleration-structure boxes of a mesh-backed Shape.
// Panics if the Shape wasn't constructed from Mesh2D / Mesh2DSlow.
func (s *Shape) MeshBoxes() []Box2 {
	ms, ok := s.SDF2.(*sdf.MeshSDF2)
	if !ok {
		panic(fmt.Sprintf("MeshBoxes: shape is %T, not a mesh SDF", s.SDF2))
	}
	boxes := ms.Boxes()
	out := make([]Box2, len(boxes))
	for i, b := range boxes {
		out[i] = v2.FromSDF(*b)
	}
	return out
}

// Benchmark reports the evaluation speed of the shape's SDF2.
func (s *Shape) Benchmark(description string) {
	sdf.BenchmarkSDF2(description, s.SDF2)
}

// Normal returns the surface normal at point p, computed via finite
// differences with sample step eps. p does not need to lie on the surface.
func (s *Shape) Normal(p v2.Vec, eps float64) v2.Vec {
	return v2.Vec(sdf.Normal2(s.SDF2, v2sdf.Vec(p), eps))
}

// Raycast sphere-traces a ray from `from` in direction `dir` until it hits
// the shape's surface. See Solid.Raycast for parameter semantics.
func (s *Shape) Raycast(from, dir v2.Vec, scaleAndSigmoid, stepScale, epsilon, maxDist float64, maxSteps int) (v2.Vec, float64, int) {
	hit, t, steps := sdf.Raycast2(s.SDF2, v2sdf.Vec(from), v2sdf.Vec(dir), scaleAndSigmoid, stepScale, epsilon, maxDist, maxSteps)
	return v2.Vec(hit), t, steps
}

// GenerateMesh samples the shape's interior on the given integer grid and
// returns the set of mesh points. Useful for surface-point extraction.
func (s *Shape) GenerateMesh(grid v2i.Vec) (v2.VecSet, error) {
	pts, err := sdf.GenerateMesh2D(s.SDF2, v2isdf.Vec(grid))
	if err != nil {
		return nil, err
	}
	out := make(v2.VecSet, len(pts))
	for i, p := range pts {
		out[i] = v2.Vec(p)
	}
	return out, nil
}

// --- Render methods ---

// ToDXF renders the shape to a DXF file using the quadtree marching-squares renderer.
func (s *Shape) ToDXF(path string, meshCells int) {
	flrender.ToDXF(s, path, meshCells)
}

// ToSVG renders the shape to an SVG file using the quadtree marching-squares renderer.
func (s *Shape) ToSVG(path string, meshCells int) {
	flrender.ToSVG(s, path, meshCells)
}

// ToPNG rasterizes the shape to a PNG file of the given pixel dimensions,
// centered on the shape's bounding box.
func (s *Shape) ToPNG(path string, width, height int) {
	flrender.ToPNG(s, path, s.Bounds(), width, height)
}

// ToPNGBox rasterizes the shape to a PNG file for the given bounding box and pixel dimensions.
func (s *Shape) ToPNGBox(path string, bb Box2, width, height int) {
	flrender.ToPNG(s, path, bb, width, height)
}
