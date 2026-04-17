package shape

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	v2sdf "github.com/deadsy/sdfx/vec/v2"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

type Shape struct {
	sdf.SDF2
}

// Wrap wraps a raw SDF2 into a Shape.
func Wrap2D(s sdf.SDF2) *Shape {
	return &Shape{s}
}

// BoundingBox returns the shape's 2D axis-aligned bounding box.
func (s *Shape) BoundingBox() Box2 {
	return v2.FromSDF(s.SDF2.BoundingBox())
}

// Evaluate returns the signed distance from p to the shape's surface.
func (s *Shape) Evaluate(p v2.Vec) float64 {
	return s.SDF2.Evaluate(v2sdf.Vec(p))
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

func (s *Shape) Cut(other *Shape) *Shape {
	return &Shape{sdf.Difference2D(s.SDF2, other.SDF2)}
}

func (s *Shape) Intersect(other *Shape) *Shape {
	return &Shape{sdf.Intersect2D(s.SDF2, other.SDF2)}
}
