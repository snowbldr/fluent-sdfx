package solid

import (
	"math"
	"runtime"
	"sync"

	"github.com/snowbldr/sdfx/sdf"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
)

// GridField is a solid's signed distance sampled once onto a dense,
// axis-aligned grid and reconstructed by trilinear interpolation. It
// exists to make an expensive field cheap: a mesh-backed SDF costs a
// spatial-index query per evaluation, and a marching-cubes render of a mold
// block asks for billions of them. Sampled once, every later evaluation is
// eight array reads.
//
// It replaces sdf.VoxelSDF3, which stored the same samples in a Go map (24
// bytes of key and a hash per corner — a 512³ grid ran to ~10 GB) and
// filled it single-threaded. This stores float32 in a flat slice (512³ is
// 537 MB, 256³ is 67 MB) and fills in parallel.
//
// Three properties matter for what sits on top of it:
//
//   - It is the plain trilinear interpolant of exact samples — no global
//     or per-voxel rescaling. VoxelSDF3 scales every value by 1/√3 to
//     keep the octree renderer's empty-cube skip (|d(centre)| ≥ the
//     cube's half-diagonal, no margin) provably safe; that makes every
//     grid-backed Offset or Shell a third too small, which is worse than
//     the problem it solves. A per-voxel bound was tried and is worse in
//     a subtler way: the bound is loose exactly on the field's medial
//     axis (the centre of a sphere, the mid-plane of a slab), where each
//     axis's corner difference is a full cell, so deep-interior values
//     were scaled by 1/√3 too and a 20 mm shell came out 30% thick.
//
//     What the unscaled interpolant actually does to the renderer: a
//     convex combination of 1-Lipschitz samples overestimates the true
//     distance by at most the distance to the nearest sample, and in
//     practice by h²/(8R) on smooth regions and a small fraction of h
//     within a cell of a sharp edge. That is the same order as the
//     renderer's own vertex placement, and it can only mis-skip a cube
//     whose face the surface grazes within that margin — a sliver at
//     most, as with sdfx's own inexact fields (smooth booleans, offsets).
//
//   - Outside the grid it returns a lower bound on the distance to the
//     surface: the larger of (distance to the sampled box + the smallest
//     boundary value) and (the boundary value at the nearest point −
//     distance to it). Continuous across the boundary, exact on the
//     faces, weak far out in the corner regions. Sample the region you
//     will read (SampleGridBounds).
//
//   - Magnitudes are exact up to interpolation error everywhere, deep
//     interior included, so Offset, Shell and clearances built on a
//     GridField mean what they say.
//
//   - Features smaller than a cell are lost. That is inherent to any grid;
//     choose cells for the smallest feature you need, not the part size.
type GridField struct {
	vals    []float32 // (nx+1)*(ny+1)*(nz+1) corner samples, x fastest
	nx, ny  int
	nz      int
	bb      sdf.Box3 // the sampled box (source bounds plus padding)
	h       v3sdf.Vec
	minEdge float64 // smallest sample on the box's surface (≥ 0 with padding)
	src     sdf.Box3
}

// GridInfo reports what a GridField cost.
type GridInfo struct {
	Cells    [3]int  // voxels per axis
	CellSize float64 // mm per voxel (uniform)
	Samples  int     // corner samples stored
	Bytes    int     // memory held by the samples
}

// SampleGrid samples s onto a GridField with `cells` voxels along its
// longest axis (the other axes scale to keep voxels cubic), padded by two
// voxels on every side so the surface never touches the sampled boundary.
// The fill runs across all CPUs.
//
// Choose the sampled region for where the field will be READ, not just
// where the surface is: outside the grid the field is only a lower bound
// (see the type comment), and a boolean partner that reads the pattern's
// field far from it — a cottle around a pattern, say — will place its
// own surface slightly wrong wherever that bound is weaker than its own
// distance. SampleGridBounds takes the region explicitly.
func (s *Solid) SampleGrid(cells int) *Solid {
	g, _ := SampleGridInfo(s, cells)
	return g
}

// SampleGridInfo is SampleGrid and the grid's cost.
func SampleGridInfo(s *Solid, cells int) (*Solid, GridInfo) {
	src := s.SDF3.BoundingBox()
	size := src.Size()
	longest := math.Max(size.X, math.Max(size.Y, size.Z))
	if cells < 2 {
		cells = 2
	}
	h := longest / float64(cells)
	pad := v3sdf.Vec{X: 2 * h, Y: 2 * h, Z: 2 * h}
	return sampleGrid(s, sdf.Box3{Min: src.Min.Sub(pad), Max: src.Max.Add(pad)}, h)
}

// SampleGridBounds samples s over the region bb — which should cover
// everywhere the field will be evaluated, and must contain s's bounds —
// with `cells` voxels along bb's longest axis.
func (s *Solid) SampleGridBounds(bb Box3, cells int) *Solid {
	g, _ := SampleGridBoundsInfo(s, bb, cells)
	return g
}

// SampleGridBoundsInfo is SampleGridBounds and the grid's cost.
func SampleGridBoundsInfo(s *Solid, bb Box3, cells int) (*Solid, GridInfo) {
	region := sdf.Box3{Min: v3sdf.Vec(bb.Min), Max: v3sdf.Vec(bb.Max)}
	src := s.SDF3.BoundingBox()
	// The region must contain the source, or the surface would cross the
	// sampled boundary and the outside bound would be a lie.
	region.Min = v3sdf.Vec{X: math.Min(region.Min.X, src.Min.X), Y: math.Min(region.Min.Y, src.Min.Y), Z: math.Min(region.Min.Z, src.Min.Z)}
	region.Max = v3sdf.Vec{X: math.Max(region.Max.X, src.Max.X), Y: math.Max(region.Max.Y, src.Max.Y), Z: math.Max(region.Max.Z, src.Max.Z)}
	size := region.Size()
	longest := math.Max(size.X, math.Max(size.Y, size.Z))
	if cells < 2 {
		cells = 2
	}
	return sampleGrid(s, region, longest/float64(cells))
}

// sampleGrid samples s over region at cell size h. The region is expanded
// outward to a whole number of cells so that h is exact on every axis.
func sampleGrid(s *Solid, region sdf.Box3, h float64) (*Solid, GridInfo) {
	src := s.SDF3.BoundingBox()
	size := region.Size()
	nx := int(math.Ceil(size.X/h - 1e-9))
	ny := int(math.Ceil(size.Y/h - 1e-9))
	nz := int(math.Ceil(size.Z/h - 1e-9))
	full := v3sdf.Vec{X: float64(nx) * h, Y: float64(ny) * h, Z: float64(nz) * h}
	c := region.Center()
	bb := sdf.Box3{Min: c.Sub(full.MulScalar(0.5)), Max: c.Add(full.MulScalar(0.5))}

	g := &GridField{
		vals: make([]float32, (nx+1)*(ny+1)*(nz+1)),
		nx:   nx, ny: ny, nz: nz,
		bb:  bb,
		h:   v3sdf.Vec{X: h, Y: h, Z: h},
		src: src,
	}
	g.fill(s.SDF3)
	g.minEdge = g.boundaryMin()

	info := GridInfo{
		Cells:    [3]int{nx, ny, nz},
		CellSize: h,
		Samples:  len(g.vals),
		Bytes:    len(g.vals) * 4,
	}
	return Wrap(g), info
}

// fill evaluates the source at every corner, one z-slab per goroutine.
func (g *GridField) fill(src sdf.SDF3) {
	nx1, ny1 := g.nx+1, g.ny+1
	workers := runtime.NumCPU()
	var wg sync.WaitGroup
	slabs := make(chan int, g.nz+1)
	for k := 0; k <= g.nz; k++ {
		slabs <- k
	}
	close(slabs)
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := range slabs {
				z := g.bb.Min.Z + float64(k)*g.h.Z
				base := k * nx1 * ny1
				for j := 0; j <= g.ny; j++ {
					y := g.bb.Min.Y + float64(j)*g.h.Y
					row := base + j*nx1
					for i := 0; i <= g.nx; i++ {
						x := g.bb.Min.X + float64(i)*g.h.X
						g.vals[row+i] = float32(src.Evaluate(v3sdf.Vec{X: x, Y: y, Z: z}))
					}
				}
			}
		}()
	}
	wg.Wait()
}

// boundaryMin is the smallest sample on the six faces of the grid.
func (g *GridField) boundaryMin() float64 {
	nx1, ny1 := g.nx+1, g.ny+1
	m := math.Inf(1)
	at := func(i, j, k int) {
		if v := float64(g.vals[k*nx1*ny1+j*nx1+i]); v < m {
			m = v
		}
	}
	for k := 0; k <= g.nz; k++ {
		for j := 0; j <= g.ny; j++ {
			at(0, j, k)
			at(g.nx, j, k)
		}
		for i := 0; i <= g.nx; i++ {
			at(i, 0, k)
			at(i, g.ny, k)
		}
	}
	for j := 0; j <= g.ny; j++ {
		for i := 0; i <= g.nx; i++ {
			at(i, j, 0)
			at(i, j, g.nz)
		}
	}
	if m < 0 {
		m = 0
	}
	return m
}

func (g *GridField) at(i, j, k int) float64 {
	return float64(g.vals[k*(g.nx+1)*(g.ny+1)+j*(g.nx+1)+i])
}

// Evaluate reconstructs the field at p. See the type comment for the
// guarantees.
func (g *GridField) Evaluate(p v3sdf.Vec) float64 {
	// Outside the sampled box: a lower bound on the distance to the
	// surface — the distance to the box, plus the least distance any
	// boundary point has to the surface.
	q := p
	outside := false
	if p.X < g.bb.Min.X {
		q.X = g.bb.Min.X
		outside = true
	} else if p.X > g.bb.Max.X {
		q.X = g.bb.Max.X
		outside = true
	}
	if p.Y < g.bb.Min.Y {
		q.Y = g.bb.Min.Y
		outside = true
	} else if p.Y > g.bb.Max.Y {
		q.Y = g.bb.Max.Y
		outside = true
	}
	if p.Z < g.bb.Min.Z {
		q.Z = g.bb.Min.Z
		outside = true
	} else if p.Z > g.bb.Max.Z {
		q.Z = g.bb.Max.Z
		outside = true
	}
	if outside {
		// Two valid lower bounds on the distance to the surface, take the
		// larger. Any path from p to the surface crosses the boundary, so
		// d(p,S) ≥ d(p,box) + min over the boundary; and by the triangle
		// inequality d(p,S) ≥ d(q,S) − |p−q| for the clamped point q. The
		// first is weak away from the box's face centres, the second is
		// tight there and keeps the field continuous across the boundary.
		out := p.Sub(q).Length()
		return math.Max(out+g.minEdge, g.interp(q)-out)
	}
	return g.interp(p)
}

// interp is the trilinear reconstruction at a point inside the grid.
func (g *GridField) interp(p v3sdf.Vec) float64 {

	// Voxel index and the fractional position inside it; the far faces of
	// the last voxels belong to those voxels, not to a phantom one beyond.
	fx := (p.X - g.bb.Min.X) / g.h.X
	fy := (p.Y - g.bb.Min.Y) / g.h.Y
	fz := (p.Z - g.bb.Min.Z) / g.h.Z
	i, j, k := int(fx), int(fy), int(fz)
	if i >= g.nx {
		i = g.nx - 1
	}
	if j >= g.ny {
		j = g.ny - 1
	}
	if k >= g.nz {
		k = g.nz - 1
	}
	dx, dy, dz := fx-float64(i), fy-float64(j), fz-float64(k)

	c000 := g.at(i, j, k)
	c100 := g.at(i+1, j, k)
	c010 := g.at(i, j+1, k)
	c110 := g.at(i+1, j+1, k)
	c001 := g.at(i, j, k+1)
	c101 := g.at(i+1, j, k+1)
	c011 := g.at(i, j+1, k+1)
	c111 := g.at(i+1, j+1, k+1)

	c00 := c000*(1-dx) + c100*dx
	c10 := c010*(1-dx) + c110*dx
	c01 := c001*(1-dx) + c101*dx
	c11 := c011*(1-dx) + c111*dx
	c0 := c00*(1-dy) + c10*dy
	c1 := c01*(1-dy) + c11*dy
	return c0*(1-dz) + c1*dz
}

// BoundingBox is the source solid's box: the padding is an implementation
// detail, and reporting it would grow every neighbour that sizes itself
// off Bounds().
func (g *GridField) BoundingBox() sdf.Box3 { return g.src }
