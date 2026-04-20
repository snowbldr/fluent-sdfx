package shape

import (
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	"github.com/snowbldr/sdfx/sdf"
)

// CubicSpline returns a 2D SDF approximation of a closed cubic spline through the given knots.
func CubicSpline(knots []v2.Vec) *Shape {
	s, err := sdf.CubicSpline2D(v2Slice(knots))
	if err != nil {
		panic(err)
	}
	return &Shape{s}
}

// Nagon returns the vertices of a regular n-gon inscribed in a circle of the given radius.
func Nagon(n int, radius float64) []v2.Vec {
	vs := sdf.Nagon(n, radius)
	out := make([]v2.Vec, len(vs))
	for i, v := range vs {
		out[i] = v2.Vec(v)
	}
	return out
}
