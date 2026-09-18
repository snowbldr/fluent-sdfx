package validate

import (
	"math"
	"runtime"
	"sort"
	"sync"

	"github.com/snowbldr/fluent-sdfx/mesh"
	"github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
)

// At returns the signed distance field of s at p: negative inside material,
// positive in air, zero on the surface. It takes the fluent vector type, so
// callers never need the sdfx vector package.
//
// The value is exact for primitives and rigid transforms of them, and a
// distance bound for booleans, twists, scales and Lipschitz-corrected
// solids: the sign is always right, but the magnitude can understate the
// true distance. Use At for sign questions and the surface-finding helpers
// (SurfaceAlong, Opening, WallThickness) when you need a measurement.
func At(s *solid.Solid, p v3.Vec) float64 {
	return s.Evaluate(v3sdf.Vec(p))
}

// Inside reports whether p is inside the material of s (At(s, p) < 0).
func Inside(s *solid.Solid, p v3.Vec) bool {
	return At(s, p) < 0
}

// Cyl returns the point at radius r, angle deg (degrees, counter-clockwise
// from +X seen from +Z) and height z: the cylindrical coordinate most
// probes on turned or patterned parts are written in.
func Cyl(r, deg, z float64) v3.Vec {
	a := deg * math.Pi / 180
	return v3.XYZ(r*math.Cos(a), r*math.Sin(a), z)
}

// Gradient returns the unit surface normal direction of the field at p by
// central differences: the direction in which the field increases, which on
// the surface points out of the material. Returns the zero vector where the
// field is flat.
func Gradient(s *solid.Solid, p v3.Vec) v3.Vec {
	const h = 1e-3
	dx := At(s, p.Add(v3.X(h))) - At(s, p.Sub(v3.X(h)))
	dy := At(s, p.Add(v3.Y(h))) - At(s, p.Sub(v3.Y(h)))
	dz := At(s, p.Add(v3.Z(h))) - At(s, p.Sub(v3.Z(h)))
	l := math.Sqrt(dx*dx + dy*dy + dz*dz)
	if l == 0 {
		return v3.Vec{}
	}
	return v3.XYZ(dx/l, dy/l, dz/l)
}

// Surface returns points lying on the surface of s: the deduplicated
// vertices of a marching-cubes mesh at cellsPerMM. These are the sample
// points every surface-to-surface measurement in this package works from.
// They sit on the rendered surface, which is within about half a cell of
// the true one.
func Surface(s *solid.Solid, cellsPerMM float64) []v3.Vec {
	tris := mesh.CollectTriangles(s, render.NewMarchingCubesOctreeParallel(solid.CellsFor(s, cellsPerMM)))
	return Vertices(tris)
}

// Vertices returns the deduplicated vertices of a triangle mesh.
func Vertices(tris []mesh.Triangle3) []v3.Vec {
	seen := make(map[v3.Vec]struct{}, len(tris))
	out := make([]v3.Vec, 0, len(tris))
	for _, t := range tris {
		for _, p := range t {
			if _, ok := seen[p]; ok {
				continue
			}
			seen[p] = struct{}{}
			out = append(out, p)
		}
	}
	return out
}

// Percentile returns the q-th percentile (0..100) of xs by nearest rank on
// a sorted copy; NaN for an empty slice.
func Percentile(xs []float64, q float64) float64 {
	if len(xs) == 0 {
		return math.NaN()
	}
	c := append([]float64(nil), xs...)
	sort.Float64s(c)
	i := int(math.Ceil(q/100*float64(len(c)))) - 1
	if i < 0 {
		i = 0
	}
	if i >= len(c) {
		i = len(c) - 1
	}
	return c[i]
}

// parallel runs fn(i) for i in [0, n) across all CPUs.
func parallel(n int, fn func(i int)) {
	workers := runtime.NumCPU()
	if workers > n {
		workers = n
	}
	if workers <= 1 {
		for i := 0; i < n; i++ {
			fn(i)
		}
		return
	}
	var wg sync.WaitGroup
	next := make(chan int, workers*4)
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range next {
				fn(i)
			}
		}()
	}
	for i := 0; i < n; i++ {
		next <- i
	}
	close(next)
	wg.Wait()
}

// evalAll evaluates s at every point in parallel.
func evalAll(s *solid.Solid, pts []v3.Vec) []float64 {
	out := make([]float64, len(pts))
	parallel(len(pts), func(i int) { out[i] = At(s, pts[i]) })
	return out
}
