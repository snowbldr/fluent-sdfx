package solid

import (
	"github.com/snowbldr/sdfx/sdf"
)

// Blend controls how unions/differences/intersections are combined.
// Use ChamferMin, RoundMin, ExpMin, PowMin, PolyMin for smooth unions.
// Use PolyMax for smooth differences/intersections.
type MinFunc = sdf.MinFunc
type MaxFunc = sdf.MaxFunc

// RoundMin returns a smooth min function that rounds the inner edge with radius k.
func RoundMin(k float64) MinFunc { return sdf.RoundMin(k) }

// ChamferMin returns a min function that produces a 45° chamfer of size k at the union seam.
func ChamferMin(k float64) MinFunc { return sdf.ChamferMin(k) }

// ExpMin returns an exponential smooth min with strength k.
func ExpMin(k float64) MinFunc { return sdf.ExpMin(k) }

// PowMin returns a power-based smooth min with exponent k.
func PowMin(k float64) MinFunc { return sdf.PowMin(k) }

// PolyMin returns a polynomial smooth min with parameter k.
func PolyMin(k float64) MinFunc { return sdf.PolyMin(k) }

// PolyMax returns a polynomial smooth max with parameter k.
func PolyMax(k float64) MaxFunc { return sdf.PolyMax(k) }

// SmoothUnion returns a union of solids blended with the given MinFunc.
func SmoothUnion(min MinFunc, solids ...*Solid) *Solid {
	sdf3s := make([]sdf.SDF3, len(solids))
	for i, s := range solids {
		sdf3s[i] = s.SDF3
	}
	u := sdf.Union3D(sdf3s...)
	u.(*sdf.UnionSDF3).SetMin(min)
	return &Solid{u}
}

// SmoothDifference subtracts the union of tools from s, blended with the given MaxFunc.
func SmoothDifference(max MaxFunc, s *Solid, tools ...*Solid) *Solid {
	sdf3s := make([]sdf.SDF3, len(tools))
	for i, t := range tools {
		sdf3s[i] = t.SDF3
	}
	d := sdf.Difference3D(s.SDF3, sdf.Union3D(sdf3s...))
	d.(*sdf.DifferenceSDF3).SetMax(max)
	return &Solid{d}
}

// SmoothCut is an alias for SmoothDifference.
func SmoothCut(max MaxFunc, s *Solid, tools ...*Solid) *Solid {
	return SmoothDifference(max, s, tools...)
}

// SmoothAdd is an alias for SmoothUnion.
func SmoothAdd(min MinFunc, solids ...*Solid) *Solid {
	return SmoothUnion(min, solids...)
}

// SmoothIntersection intersects s with the union of others, blended with the given MaxFunc.
func SmoothIntersection(max MaxFunc, s *Solid, others ...*Solid) *Solid {
	sdf3s := make([]sdf.SDF3, len(others))
	for i, o := range others {
		sdf3s[i] = o.SDF3
	}
	i := sdf.Intersect3D(s.SDF3, sdf.Union3D(sdf3s...))
	i.(*sdf.IntersectionSDF3).SetMax(max)
	return &Solid{i}
}
