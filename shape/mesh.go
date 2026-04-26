package shape

import (
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
