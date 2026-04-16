package shape

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
)

type Shape struct {
	sdf.SDF2
}

// Wrap wraps a raw SDF2 into a Shape.
func Wrap2D(s sdf.SDF2) *Shape {
	return &Shape{s}
}

// --- Constructors ---

func Rect(size v2.Vec, round float64) *Shape {
	return &Shape{sdf.Box2D(size, round)}
}

func Circle(radius float64) *Shape {
	s, _ := sdf.Circle2D(radius)
	return &Shape{s}
}

func Polygon(pts []v2.Vec) *Shape {
	s, err := sdf.FlatPolygon2D(pts)
	if err != nil {
		panic(err)
	}
	return &Shape{s}
}

// --- Transform methods ---

func (s *Shape) Translate(v v2.Vec) *Shape {
	return &Shape{sdf.Transform2D(s, sdf.Translate2d(v))}
}

func (s *Shape) Rotate(angleDeg float64) *Shape {
	return &Shape{sdf.Transform2D(s, sdf.Rotate2d(angleDeg*math.Pi/180))}
}

func (s *Shape) Scale(v v2.Vec) *Shape {
	return &Shape{sdf.Transform2D(s, sdf.Scale2d(v))}
}

func (s *Shape) MirrorX() *Shape {
	return &Shape{sdf.Transform2D(s, sdf.MirrorX())}
}

func (s *Shape) MirrorY() *Shape {
	return &Shape{sdf.Transform2D(s, sdf.MirrorY())}
}

func (s *Shape) Transform(m sdf.M33) *Shape {
	return &Shape{sdf.Transform2D(s, m)}
}

func (s *Shape) Offset(amount float64) *Shape {
	return &Shape{sdf.Offset2D(s, amount)}
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
