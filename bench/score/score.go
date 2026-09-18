package score

import (
	"fmt"
	"math"
	"runtime/debug"
	"sort"
	"sync"

	"github.com/snowbldr/fluent-sdfx/mesh"
	"github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
)

// ProbeResult is one evaluated probe.
type ProbeResult struct {
	Name   string  `json:"name"`
	Want   bool    `json:"want_inside"`
	Got    bool    `json:"got_inside"`
	SDF    float64 `json:"sdf"`
	Passed bool    `json:"passed"`
}

// Result is the staged score of one candidate.
type Result struct {
	Task     string `json:"task"`
	Compiled bool   `json:"compiled"`
	Built    bool   `json:"built"`
	Error    string `json:"error,omitempty"`

	RefBounds   [2][3]float64 `json:"ref_bounds"`
	CandBounds  [2][3]float64 `json:"cand_bounds"`      // after alignment
	BBoxSizeErr float64       `json:"bbox_size_err_mm"` // max abs diff of bbox size per axis
	BBoxPosErr  float64       `json:"bbox_pos_err_mm"`  // max abs diff of bbox center per axis (after alignment)

	CellMM    float64 `json:"cell_mm"`
	RefVolume float64 `json:"ref_volume_mm3"` // grid-estimated
	IoU       float64 `json:"iou"`
	Missing   float64 `json:"missing"` // reference volume absent from candidate, as fraction of reference volume
	Extra     float64 `json:"extra"`   // candidate volume absent from reference, as fraction of reference volume

	// MaxDev is the symmetric Hausdorff distance between the two voxelised
	// surfaces, in mm, at one-cell resolution. A reference scored against
	// itself reports 0.
	MaxDev   float64    `json:"max_dev_mm"`
	MaxDevAt [3]float64 `json:"max_dev_at"` // cell centre where MaxDev was observed
	MaxDevOn string     `json:"max_dev_on"` // "ref-surface" (far from candidate) or "cand-surface" (far from reference)

	// MaxDevSDF is the older SDF-magnitude measure: the largest |sdf| of one
	// solid sampled at grid points on the other's surface band. Sensitive to
	// sub-cell shifts, but inflated by SDFs that are not true distances.
	// Informational only; not a pass criterion.
	MaxDevSDF   float64    `json:"max_dev_sdf_mm"`
	MaxDevSDFAt [3]float64 `json:"max_dev_sdf_at"`

	Probes      []ProbeResult `json:"probes"`
	ProbesPass  int           `json:"probes_pass"`
	ProbesTotal int           `json:"probes_total"`

	Mesh *MeshResult `json:"mesh,omitempty"`

	Passed  bool     `json:"passed"`
	Reasons []string `json:"reasons,omitempty"`
}

// MeshResult is the optional mesh-stage output.
type MeshResult struct {
	Triangles     int     `json:"triangles"`
	BoundaryEdges int     `json:"boundary_edges"`
	Watertight    bool    `json:"watertight"`
	Components    int     `json:"components"`
	Volume        float64 `json:"volume_mm3"`
}

// Guard runs a Build func and converts a panic into an error.
func Guard(build func() *solid.Solid) (s *solid.Solid, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v\n%s", r, debug.Stack())
		}
	}()
	s = build()
	if s == nil {
		return nil, fmt.Errorf("Build returned nil")
	}
	return s, nil
}

// Compare scores cand against ref under spec. Both must be non-nil.
func Compare(ref, cand *solid.Solid, spec Spec) Result {
	spec.Defaults()
	r := Result{Task: spec.Name, Compiled: true, Built: true}

	rb := ref.Bounds()
	if spec.Align == "center" {
		cb := cand.Bounds()
		cand = cand.Translate(rb.Center().Sub(cb.Center()))
	}
	cb := cand.Bounds()
	r.RefBounds = boxArr(rb.Box)
	r.CandBounds = boxArr(cb.Box)
	r.BBoxSizeErr = maxAbs3(rb.Size().Sub(cb.Size()))
	r.BBoxPosErr = maxAbs3(rb.Center().Sub(cb.Center()))

	// Grid over the union bbox with a one-cell margin.
	ub := rb.Box.Extend(cb.Box)
	size := ub.Size()
	longest := math.Max(size.X, math.Max(size.Y, size.Z))
	cell := longest / float64(spec.Grid)
	r.CellMM = cell
	// Sample lattice: one full cell of margin on every side, and cell centres
	// offset by an irrational fraction of a cell so they never land exactly
	// on a face plane that sits a whole number of cells from the bbox minimum
	// (where a coincident-face SDF is exactly 0 and the sign is ambiguous).
	ub = v3.Box{Min: ub.Min.Sub(v3.XYZ(cell, cell, cell)), Max: ub.Max.Add(v3.XYZ(cell, cell, cell))}
	size = ub.Size()
	nx := int(math.Ceil(size.X / cell))
	ny := int(math.Ceil(size.Y / cell))
	nz := int(math.Ceil(size.Z / cell))
	const latticeOffset = 0.381966 // 1 - 1/phi

	// Occupancy grids, filled by z-slab in parallel. Each goroutine writes
	// only its own slab, so no locking is needed.
	n := nx * ny * nz
	occA := make([]bool, n)
	occB := make([]bool, n)
	sdfDev := make([]float64, nz) // per-slab max |other sdf| on a surface band (informational)
	sdfAt := make([]v3sdf.Vec, nz)
	var wg sync.WaitGroup
	sem := make(chan struct{}, 8)
	for iz := 0; iz < nz; iz++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(iz int) {
			defer wg.Done()
			defer func() { <-sem }()
			z := ub.Min.Z + (float64(iz)+latticeOffset)*cell
			for iy := 0; iy < ny; iy++ {
				y := ub.Min.Y + (float64(iy)+latticeOffset)*cell
				for ix := 0; ix < nx; ix++ {
					x := ub.Min.X + (float64(ix)+latticeOffset)*cell
					p := v3sdf.Vec{X: x, Y: y, Z: z}
					a := ref.Evaluate(p)
					b := cand.Evaluate(p)
					i := (iz*ny+iy)*nx + ix
					occA[i] = a < 0
					occB[i] = b < 0
					if math.Abs(a) < cell*0.5 {
						if d := math.Abs(b); d > sdfDev[iz] {
							sdfDev[iz], sdfAt[iz] = d, p
						}
					}
					if math.Abs(b) < cell*0.5 {
						if d := math.Abs(a); d > sdfDev[iz] {
							sdfDev[iz], sdfAt[iz] = d, p
						}
					}
				}
			}
		}(iz)
	}
	wg.Wait()

	var inter, union, refOnly, candOnly int
	for i := 0; i < n; i++ {
		switch {
		case occA[i] && occB[i]:
			inter++
			union++
		case occA[i]:
			refOnly++
			union++
		case occB[i]:
			candOnly++
			union++
		}
	}
	refCells := inter + refOnly
	r.RefVolume = float64(refCells) * cell * cell * cell
	if union > 0 {
		r.IoU = float64(inter) / float64(union)
	}
	if refCells > 0 {
		r.Missing = float64(refOnly) / float64(refCells)
		r.Extra = float64(candOnly) / float64(refCells)
	}

	// Surface deviation: symmetric Hausdorff distance between the two
	// voxelised surfaces (cells whose occupancy differs from a 6-neighbour),
	// via a chamfer distance transform. This depends only on occupancy, so
	// SDFs that are not true distances (twists, scales, loose CSG bounds,
	// internal membranes) cannot inflate it. Resolution is one cell.
	surfA := surface(occA, nx, ny, nz)
	surfB := surface(occB, nx, ny, nz)
	dA, iA := hausdorff(surfA, surfB, nx, ny, nz) // farthest A-surface cell from B's surface
	dB, iB := hausdorff(surfB, surfA, nx, ny, nz)
	if dA >= dB {
		r.MaxDev, r.MaxDevOn = dA*cell, "ref-surface"
		r.MaxDevAt = cellCenter(iA, nx, ny, ub.Min, cell)
	} else {
		r.MaxDev, r.MaxDevOn = dB*cell, "cand-surface"
		r.MaxDevAt = cellCenter(iB, nx, ny, ub.Min, cell)
	}
	for iz := range sdfDev {
		if sdfDev[iz] > r.MaxDevSDF {
			r.MaxDevSDF = sdfDev[iz]
			r.MaxDevSDFAt = [3]float64{sdfAt[iz].X, sdfAt[iz].Y, sdfAt[iz].Z}
		}
	}

	// Probes are given in the reference frame; the candidate has already
	// been aligned into it.
	for _, p := range spec.Probes {
		v := cand.Evaluate(v3sdf.Vec{X: p.P[0], Y: p.P[1], Z: p.P[2]})
		got := v < 0
		pr := ProbeResult{Name: p.Name, Want: p.Inside, Got: got, SDF: v, Passed: got == p.Inside}
		if pr.Passed {
			r.ProbesPass++
		}
		r.Probes = append(r.Probes, pr)
	}
	r.ProbesTotal = len(spec.Probes)

	if spec.MeshCellsPerMM > 0 {
		r.Mesh = MeshStage(cand, spec.MeshCellsPerMM)
	}

	// Verdict.
	if r.IoU < spec.Pass.IoUMin {
		r.Reasons = append(r.Reasons, fmt.Sprintf("iou %.3f < %.3f", r.IoU, spec.Pass.IoUMin))
	}
	if r.MaxDev > spec.Pass.MaxDevMM {
		r.Reasons = append(r.Reasons, fmt.Sprintf("max_dev %.2fmm > %.2fmm", r.MaxDev, spec.Pass.MaxDevMM))
	}
	if r.Extra > spec.Pass.ExtraMax {
		r.Reasons = append(r.Reasons, fmt.Sprintf("extra volume %.1f%% > %.1f%%", r.Extra*100, spec.Pass.ExtraMax*100))
	}
	if r.ProbesPass < r.ProbesTotal {
		var failed []string
		for _, p := range r.Probes {
			if !p.Passed {
				failed = append(failed, p.Name)
			}
		}
		sort.Strings(failed)
		r.Reasons = append(r.Reasons, fmt.Sprintf("probes failed: %v", failed))
	}
	if r.Mesh != nil && (!r.Mesh.Watertight || r.Mesh.Components != 1) {
		r.Reasons = append(r.Reasons, fmt.Sprintf("mesh: watertight=%v components=%d", r.Mesh.Watertight, r.Mesh.Components))
	}
	r.Passed = len(r.Reasons) == 0
	return r
}

// MeshStage renders cand and reports watertightness, volume and the number
// of connected components (vertex-shared triangle islands).
func MeshStage(cand *solid.Solid, cellsPerMM float64) *MeshResult {
	tris := mesh.CollectTriangles(cand, render.NewMarchingCubesOctreeParallel(solid.CellsFor(cand, cellsPerMM)))
	bnd := mesh.CountBoundaryEdges(tris)
	return &MeshResult{
		Triangles:     len(tris),
		BoundaryEdges: bnd,
		Watertight:    bnd == 0,
		Components:    Components(tris),
		Volume:        volume(tris),
	}
}

// Components counts connected components of a triangle soup by shared
// (quantized) vertices using union-find.
func Components(tris []mesh.Triangle3) int {
	if len(tris) == 0 {
		return 0
	}
	const q = 1e-6
	key := func(v v3.Vec) [3]int64 {
		return [3]int64{int64(math.Round(v.X / q)), int64(math.Round(v.Y / q)), int64(math.Round(v.Z / q))}
	}
	ids := map[[3]int64]int{}
	parent := []int{}
	find := func(i int) int {
		for parent[i] != i {
			parent[i] = parent[parent[i]]
			i = parent[i]
		}
		return i
	}
	union := func(a, b int) {
		ra, rb := find(a), find(b)
		if ra != rb {
			parent[ra] = rb
		}
	}
	id := func(v v3.Vec) int {
		k := key(v)
		if i, ok := ids[k]; ok {
			return i
		}
		i := len(parent)
		ids[k] = i
		parent = append(parent, i)
		return i
	}
	for _, t := range tris {
		a, b, c := id(t[0]), id(t[1]), id(t[2])
		union(a, b)
		union(b, c)
	}
	roots := map[int]struct{}{}
	for i := range parent {
		roots[find(i)] = struct{}{}
	}
	return len(roots)
}

func volume(tris []mesh.Triangle3) float64 {
	var v float64
	for _, t := range tris {
		v += t[0].Dot(t[1].Cross(t[2])) / 6.0
	}
	return math.Abs(v)
}

func boxArr(b v3.Box) [2][3]float64 {
	return [2][3]float64{{b.Min.X, b.Min.Y, b.Min.Z}, {b.Max.X, b.Max.Y, b.Max.Z}}
}

func maxAbs3(v v3.Vec) float64 {
	return math.Max(math.Abs(v.X), math.Max(math.Abs(v.Y), math.Abs(v.Z)))
}

// surface marks cells whose occupancy differs from any 6-neighbour.
func surface(occ []bool, nx, ny, nz int) []bool {
	out := make([]bool, len(occ))
	idx := func(x, y, z int) int { return (z*ny+y)*nx + x }
	for z := 0; z < nz; z++ {
		for y := 0; y < ny; y++ {
			for x := 0; x < nx; x++ {
				i := idx(x, y, z)
				v := occ[i]
				if (x > 0 && occ[idx(x-1, y, z)] != v) || (x < nx-1 && occ[idx(x+1, y, z)] != v) ||
					(y > 0 && occ[idx(x, y-1, z)] != v) || (y < ny-1 && occ[idx(x, y+1, z)] != v) ||
					(z > 0 && occ[idx(x, y, z-1)] != v) || (z < nz-1 && occ[idx(x, y, z+1)] != v) {
					out[i] = true
				}
			}
		}
	}
	return out
}

// hausdorff returns the largest chamfer distance (in cells) from any cell of
// from to the nearest cell of to, and the index of that cell. A 3-4-5-style
// two-pass chamfer transform with weights 1, sqrt2, sqrt3.
func hausdorff(from, to []bool, nx, ny, nz int) (float64, int) {
	const inf = float32(1e9)
	n := nx * ny * nz
	dt := make([]float32, n)
	for i := range dt {
		if to[i] {
			dt[i] = 0
		} else {
			dt[i] = inf
		}
	}
	idx := func(x, y, z int) int { return (z*ny+y)*nx + x }
	w := func(dx, dy, dz int) float32 {
		switch dx*dx + dy*dy + dz*dz {
		case 1:
			return 1
		case 2:
			return 1.41421356
		default:
			return 1.73205081
		}
	}
	// Forward pass: neighbours "before" in scan order.
	for z := 0; z < nz; z++ {
		for y := 0; y < ny; y++ {
			for x := 0; x < nx; x++ {
				i := idx(x, y, z)
				d := dt[i]
				if d == 0 {
					continue
				}
				for dz := -1; dz <= 0; dz++ {
					for dy := -1; dy <= 1; dy++ {
						for dx := -1; dx <= 1; dx++ {
							if dz == 0 && (dy > 0 || (dy == 0 && dx >= 0)) {
								continue
							}
							xx, yy, zz := x+dx, y+dy, z+dz
							if xx < 0 || yy < 0 || zz < 0 || xx >= nx || yy >= ny || zz >= nz {
								continue
							}
							if v := dt[idx(xx, yy, zz)] + w(dx, dy, dz); v < d {
								d = v
							}
						}
					}
				}
				dt[i] = d
			}
		}
	}
	// Backward pass.
	for z := nz - 1; z >= 0; z-- {
		for y := ny - 1; y >= 0; y-- {
			for x := nx - 1; x >= 0; x-- {
				i := idx(x, y, z)
				d := dt[i]
				if d == 0 {
					continue
				}
				for dz := 0; dz <= 1; dz++ {
					for dy := -1; dy <= 1; dy++ {
						for dx := -1; dx <= 1; dx++ {
							if dz == 0 && (dy < 0 || (dy == 0 && dx <= 0)) {
								continue
							}
							xx, yy, zz := x+dx, y+dy, z+dz
							if xx < 0 || yy < 0 || zz < 0 || xx >= nx || yy >= ny || zz >= nz {
								continue
							}
							if v := dt[idx(xx, yy, zz)] + w(dx, dy, dz); v < d {
								d = v
							}
						}
					}
				}
				dt[i] = d
			}
		}
	}
	best, at := float32(0), -1
	for i := 0; i < n; i++ {
		if from[i] && dt[i] > best {
			best, at = dt[i], i
		}
	}
	if best >= inf {
		return 0, -1 // one side has no surface at all
	}
	return float64(best), at
}

func cellCenter(i, nx, ny int, min v3.Vec, cell float64) [3]float64 {
	if i < 0 {
		return [3]float64{}
	}
	x := i % nx
	y := (i / nx) % ny
	z := i / (nx * ny)
	const latticeOffset = 0.381966
	return [3]float64{min.X + (float64(x)+latticeOffset)*cell, min.Y + (float64(y)+latticeOffset)*cell, min.Z + (float64(z)+latticeOffset)*cell}
}
