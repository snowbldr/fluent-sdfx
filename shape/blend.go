package shape

import "github.com/snowbldr/sdfx/sdf"

// SmoothUnion blends shapes with a smooth min function (see solid.RoundMin etc for Mins).
func SmoothUnion(min sdf.MinFunc, shapes ...*Shape) *Shape {
	sdf2s := make([]sdf.SDF2, len(shapes))
	for i, s := range shapes {
		sdf2s[i] = s.SDF2
	}
	u := sdf.Union2D(sdf2s...)
	u.(*sdf.UnionSDF2).SetMin(min)
	return &Shape{u}
}

// SmoothCut returns a - b blended with a smooth max function.
func SmoothCut(max sdf.MaxFunc, a, b *Shape) *Shape {
	d := sdf.Difference2D(a.SDF2, b.SDF2)
	d.(*sdf.DifferenceSDF2).SetMax(max)
	return &Shape{d}
}

// SmoothIntersect returns a ∩ b blended with a smooth max function.
func SmoothIntersect(max sdf.MaxFunc, a, b *Shape) *Shape {
	i := sdf.Intersect2D(a.SDF2, b.SDF2)
	i.(*sdf.IntersectionSDF2).SetMax(max)
	return &Shape{i}
}
