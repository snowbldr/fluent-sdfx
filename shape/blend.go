package shape

import "github.com/snowbldr/sdfx/sdf"

// SmoothUnion blends shapes with a smooth min function (see solid.RoundMin
// etc for Mins). Panics with zero shapes; returns the lone shape unchanged
// when called with a single shape.
func SmoothUnion(min sdf.MinFunc, shapes ...*Shape) *Shape {
	if len(shapes) == 0 {
		panic("shape.SmoothUnion: at least one shape required")
	}
	if len(shapes) == 1 {
		return shapes[0]
	}
	sdf2s := make([]sdf.SDF2, len(shapes))
	for i, s := range shapes {
		sdf2s[i] = s.SDF2
	}
	u := sdf.Union2D(sdf2s...)
	u.(*sdf.UnionSDF2).SetMin(min)
	return &Shape{u}
}

// SmoothAdd is an alias for SmoothUnion.
func SmoothAdd(min sdf.MinFunc, shapes ...*Shape) *Shape {
	return SmoothUnion(min, shapes...)
}

// SmoothCut subtracts the union of tools from s, blended with a smooth max
// function. With no tools, returns s unchanged.
func SmoothCut(max sdf.MaxFunc, s *Shape, tools ...*Shape) *Shape {
	if len(tools) == 0 {
		return s
	}
	sdf2s := make([]sdf.SDF2, len(tools))
	for i, t := range tools {
		sdf2s[i] = t.SDF2
	}
	d := sdf.Difference2D(s.SDF2, sdf.Union2D(sdf2s...))
	d.(*sdf.DifferenceSDF2).SetMax(max)
	return &Shape{d}
}

// SmoothIntersect intersects s with the union of others, blended with a
// smooth max function. With no others, returns s unchanged.
func SmoothIntersect(max sdf.MaxFunc, s *Shape, others ...*Shape) *Shape {
	if len(others) == 0 {
		return s
	}
	sdf2s := make([]sdf.SDF2, len(others))
	for i, o := range others {
		sdf2s[i] = o.SDF2
	}
	i := sdf.Intersect2D(s.SDF2, sdf.Union2D(sdf2s...))
	i.(*sdf.IntersectionSDF2).SetMax(max)
	return &Shape{i}
}
