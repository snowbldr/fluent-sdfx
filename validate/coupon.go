package validate

import (
	"fmt"
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/mesh"
	"github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Physical-iteration helpers: cutting a small printable test piece out of a
// big part, labelling a ladder of variants so the printed parts can be told
// apart by touch, and counting how many separate bodies a solid would
// actually print as.

// couponDimpleOvershoot is how far a dimple cutter is extended proud of the
// face (mm) so its mouth cuts cleanly instead of ending coplanar with the
// surface, where marching cubes can leave a skin. The blind depth measured
// from the face is unaffected.
const couponDimpleOvershoot = 0.05

// Coupon cuts a printable test piece out of s: the intersection of s with
// keep, dropped so its minimum Z is 0 (sitting on the plate). It panics
// with the body count when the result is not a single connected body — a
// coupon sliced through a hole prints as loose parts and measures nothing —
// so check with Components first if the cut is speculative.
func Coupon(s, keep *solid.Solid, cellsPerMM float64) *solid.Solid {
	c := s.Intersect(keep).ZeroZ()
	if n := Components(c, cellsPerMM); n != 1 {
		panic(fmt.Sprintf("validate.Coupon: cut produced %d separate bodies at %.1f cells/mm, want 1 — the keep region slices through a hole or leaves an island; move the cut or print the parts separately", n, cellsPerMM))
	}
	return c
}

// Dimples returns s with n small blind cylinders cut into the face at `at`,
// sunk along -normal to depth `depth` and spaced `spacing` apart along
// `along`, centred on `at`. This is the self-labelling convention for a
// parameter ladder: embossed text will not print legibly at 4 mm, but "the
// one with three dots" is unambiguous in the hand.
//
// normal points out of the face (the direction the dimples face); along is
// the in-face direction the row runs and is used as given, so pass a
// direction that lies in the face. Both are normalised internally. Panics
// on a zero normal or along, or on a non-positive dia or depth.
func Dimples(s *solid.Solid, n int, at, normal, along v3.Vec, dia, depth, spacing float64) *solid.Solid {
	if n <= 0 {
		return s
	}
	if normal.Length() == 0 {
		panic("validate.Dimples: normal must be non-zero (it is the direction the dimpled face points)")
	}
	if along.Length() == 0 {
		panic("validate.Dimples: along must be non-zero (it is the in-face direction the row of dimples runs)")
	}
	if dia <= 0 || depth <= 0 {
		panic(fmt.Sprintf("validate.Dimples: dia and depth must be positive, got dia=%.3f depth=%.3f", dia, depth))
	}
	nrm := normal.Normalize()
	dir := along.Normalize()
	h := depth + couponDimpleOvershoot
	cutters := make([]*solid.Solid, n)
	for i := 0; i < n; i++ {
		offset := (float64(i) - float64(n-1)/2) * spacing
		// Centre the cutter so it spans from couponDimpleOvershoot proud of
		// the face down to exactly depth below it.
		p := at.Add(dir.MulScalar(offset)).Add(nrm.MulScalar((couponDimpleOvershoot - depth) / 2))
		cutters[i] = solid.Cylinder(h, dia/2, 0).RotateToVector(v3.Z(1), nrm).Translate(p)
	}
	return s.Cut(cutters...)
}

// Components counts the connected bodies in s's rendered mesh at
// cellsPerMM: triangles that share a vertex belong to the same body. A part
// that reads as one solid in the field but renders as two bodies will print
// as two loose pieces, which is what the slicer sees.
func Components(s *solid.Solid, cellsPerMM float64) int {
	tris := mesh.CollectTriangles(s, render.NewMarchingCubesOctreeParallel(solid.CellsFor(s, cellsPerMM)))
	return couponComponents(tris)
}

// couponComponents counts vertex-connected components of a triangle mesh by
// union-find over vertices quantised to couponQuantum, so two triangles
// meeting at a shared corner join even if the renderer emitted the corner
// with a last-bit difference.
func couponComponents(tris []mesh.Triangle3) int {
	const couponQuantum = 1e-6 // mm
	type key struct{ x, y, z int64 }
	quant := func(p v3.Vec) key {
		return key{
			int64(math.Round(p.X / couponQuantum)),
			int64(math.Round(p.Y / couponQuantum)),
			int64(math.Round(p.Z / couponQuantum)),
		}
	}

	ids := make(map[key]int, len(tris))
	parent := make([]int, 0, len(tris))
	rank := make([]int, 0, len(tris))
	id := func(p v3.Vec) int {
		k := quant(p)
		if i, ok := ids[k]; ok {
			return i
		}
		i := len(parent)
		ids[k] = i
		parent = append(parent, i)
		rank = append(rank, 0)
		return i
	}
	var find func(i int) int
	find = func(i int) int {
		for parent[i] != i {
			parent[i] = parent[parent[i]] // path halving
			i = parent[i]
		}
		return i
	}
	union := func(a, b int) {
		ra, rb := find(a), find(b)
		if ra == rb {
			return
		}
		if rank[ra] < rank[rb] {
			ra, rb = rb, ra
		}
		parent[rb] = ra
		if rank[ra] == rank[rb] {
			rank[ra]++
		}
	}

	for _, t := range tris {
		a, b, c := id(t[0]), id(t[1]), id(t[2])
		union(a, b)
		union(a, c)
	}
	roots := make(map[int]struct{}, 8)
	for i := range parent {
		roots[find(i)] = struct{}{}
	}
	return len(roots)
}

// RequireOneBody fails the test when s renders as anything other than a
// single connected body at cellsPerMM — the guard against a boolean that
// quietly severed a part, or a cut that left an island.
func RequireOneBody(t testing.TB, s *solid.Solid, cellsPerMM float64) {
	t.Helper()
	n := Components(s, cellsPerMM)
	t.Logf("bodies = %d at %.1f cells/mm", n, cellsPerMM)
	if n != 1 {
		t.Errorf("solid renders as %d separate bodies at %.1f cells/mm, want 1", n, cellsPerMM)
	}
}
