package shape

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/plane"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/sdfx/sdf"
	v2sdf "github.com/snowbldr/sdfx/vec/v2"
)

type Shape struct {
	sdf.SDF2
}

// Wrap wraps a raw SDF2 into a Shape.
func Wrap2D(s sdf.SDF2) *Shape {
	return &Shape{s}
}

// Bounds returns the shape's 2D axis-aligned bounding box as a fluent v2.Box.
// Use this instead of BoundingBox when you want the fluent box API
// (BoundingBox returns sdfx's raw sdf.Box2).
func (s *Shape) Bounds() Box2 {
	return v2.FromSDF(s.SDF2.BoundingBox())
}

func v2Slice(pts []v2.Vec) []v2sdf.Vec {
	out := make([]v2sdf.Vec, len(pts))
	for i, p := range pts {
		out[i] = v2sdf.Vec(p)
	}
	return out
}

// --- Constructors ---

func Rect(size v2.Vec, round float64) *Shape {
	return &Shape{sdf.Box2D(v2sdf.Vec(size), round)}
}

func Circle(radius float64) *Shape {
	s, _ := sdf.Circle2D(radius)
	return &Shape{s}
}

func Polygon(pts []v2.Vec) *Shape {
	s, err := sdf.FlatPolygon2D(v2Slice(pts))
	if err != nil {
		panic(err)
	}
	return &Shape{s}
}

// --- Transform methods ---

func (s *Shape) Translate(v v2.Vec) *Shape {
	return &Shape{sdf.Transform2D(s.SDF2, sdf.Translate2d(v2sdf.Vec(v)))}
}

func (s *Shape) Rotate(angleDeg float64) *Shape {
	return &Shape{sdf.Transform2D(s.SDF2, sdf.Rotate2d(angleDeg*math.Pi/180))}
}

func (s *Shape) Scale(v v2.Vec) *Shape {
	return &Shape{sdf.Transform2D(s.SDF2, sdf.Scale2d(v2sdf.Vec(v)))}
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

func (s *Shape) Cut(other *Shape) *Shape {
	return &Shape{sdf.Difference2D(s.SDF2, other.SDF2)}
}

// Difference is an alias for Cut.
func (s *Shape) Difference(other *Shape) *Shape {
	return s.Cut(other)
}

func (s *Shape) Intersect(other *Shape) *Shape {
	return &Shape{sdf.Intersect2D(s.SDF2, other.SDF2)}
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
