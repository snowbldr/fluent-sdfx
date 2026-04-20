package shape

import (
	"fmt"

	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	"github.com/snowbldr/sdfx/sdf"
)

// Line2 is a 2D line segment.
type Line2 = v2.Line2

func toSDFLines(m []*Line2) []*sdf.Line2 {
	out := make([]*sdf.Line2, len(m))
	for i, l := range m {
		sl := l.SDF()
		out[i] = &sl
	}
	return out
}

// Mesh2D builds a Shape from a set of 2D line segments.
func Mesh2D(m []*Line2) *Shape {
	s, err := sdf.Mesh2D(toSDFLines(m))
	if err != nil {
		panic(err)
	}
	return &Shape{s}
}

// Mesh2DSlow is a naive-algorithm Mesh2D reference implementation for testing.
func Mesh2DSlow(m []*Line2) *Shape {
	s, err := sdf.Mesh2DSlow(toSDFLines(m))
	if err != nil {
		panic(err)
	}
	return &Shape{s}
}

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
