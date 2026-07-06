package render

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"

	"github.com/snowbldr/fluent-sdfx/render/meshopt"
	"github.com/snowbldr/sdfx/render"
	"github.com/snowbldr/sdfx/sdf"
)

// DefaultDecimateMaxError is the default cap, in model units (mm), on how far
// decimation may move the surface from the rendered mesh. Well under a
// typical 3D-print layer height, so decimated and full meshes print
// identically. Planar regions collapse at ~zero error regardless, so this
// only limits simplification of curved surfaces.
const DefaultDecimateMaxError = 0.05

// ToSTL renders an SDF3 to an STL file atomically.
// Optional decimate values: decimate[0] (0-1) enables decimation and sets the
// per-pass triangle-removal aggressiveness (passes repeat to a fixed point,
// so the error budget — not this ratio — determines the final size).
// decimate[1] is the maximum surface deviation in model units (mm) that
// decimation may introduce (default DefaultDecimateMaxError).
func ToSTL(s SDF3, path string, r render.Render3, decimate ...float64) {
	fmt.Printf("rendering %s (%s)\n", path, r.Info(s))

	mesh := render.ToTriangles(s, r)
	fmt.Printf("  %d triangles", len(mesh))

	removeRatio := 0.0
	if len(decimate) > 0 {
		removeRatio = decimate[0]
	}
	maxError := DefaultDecimateMaxError
	if len(decimate) > 1 {
		maxError = decimate[1]
	}
	if removeRatio > 0 && removeRatio < 1 {
		mesh = DecimateMesh(mesh, 1-removeRatio, maxError)
	}
	fmt.Println()

	dir := filepath.Dir(path)
	if dir == "" {
		dir = "."
	}
	tmp, err := os.CreateTemp(dir, ".stl-tmp-*")
	if err != nil {
		panic(fmt.Errorf("ToSTL: create temp file in %s: %w", dir, err))
	}
	tmpPath := tmp.Name()
	tmp.Close()

	if err := render.SaveSTL(tmpPath, mesh); err != nil {
		os.Remove(tmpPath)
		panic(fmt.Errorf("ToSTL: write %s: %w", tmpPath, err))
	}

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		panic(fmt.Errorf("ToSTL: rename to %s: %w", path, err))
	}
}

// DecimateMesh simplifies a triangle mesh down to keepRatio of its original
// size, moving no surface point more than maxErrorMM (model units) from the
// input mesh. The triangle target is abandoned early rather than exceed
// maxErrorMM, so planar regions still collapse fully while curved detail is
// preserved.
//
// Marching-cubes output carries sub-hundredth-mm ripple on flat faces, so a
// tight budget can stall mid-face (each collapse eats a little budget until
// the accumulated error hits the cap). To keep flats collapsing to a handful
// of triangles at any budget, decimation runs in two passes: after the first,
// the dominant planes are detected from the largest surviving triangles,
// every vertex within maxErrorMM of a plane is snapped exactly onto it
// (correcting the render noise rather than adding error), and the now truly
// coplanar mesh is simplified again.
func DecimateMesh(mesh []*sdf.Triangle3, keepRatio, maxErrorMM float64) []*sdf.Triangle3 {
	n := len(mesh)
	if n == 0 {
		return mesh
	}
	verts := make([]float32, n*9)
	for i, t := range mesh {
		off := i * 9
		for j := 0; j < 3; j++ {
			verts[off+j*3+0] = float32(t[j].X)
			verts[off+j*3+1] = float32(t[j].Y)
			verts[off+j*3+2] = float32(t[j].Z)
		}
	}

	// meshoptimizer's error is relative to the mesh extent (longest
	// bounding-box axis) — convert the absolute mm budget to that scale.
	extent := meshExtent(mesh)
	relError := 1.0
	if extent > 0 {
		relError = maxErrorMM / extent
	}

	out, count, achieved := simplifyToFixedPoint(verts, n, keepRatio, float32(relError))

	if planes := detectPlanes(out, count, extent); len(planes) > 0 {
		snapped := snapToPlanes(out, count, planes, maxErrorMM)
		if snapped > 0 {
			before := count
			var achieved2 float32
			out, count, achieved2 = simplifyToFixedPoint(out[:count*9], count, keepRatio, float32(relError))
			if achieved2 > achieved {
				achieved = achieved2
			}
			fmt.Printf("\n  planar snap: %d planes, %d vertices snapped, %d → %d triangles",
				len(planes), snapped, before, count)
		}
	}

	result := make([]*sdf.Triangle3, count)
	for i := 0; i < count; i++ {
		off := i * 9
		tri := sdf.Triangle3{}
		for j := 0; j < 3; j++ {
			tri[j].X = float64(out[off+j*3+0])
			tri[j].Y = float64(out[off+j*3+1])
			tri[j].Z = float64(out[off+j*3+2])
		}
		result[i] = &tri
	}
	fmt.Printf("\n  → %d after decimation (kept %.1f%% of %d, deviation ≤ %.3fmm of %.3fmm budget)",
		count, 100*float64(count)/float64(n), n, float64(achieved)*extent, maxErrorMM)
	return result
}

// simplifyToFixedPoint runs meshopt simplification repeatedly until the error
// budget — not the per-pass triangle quota — is what stops it. A single pass
// stops at its target count even when further zero-error collapses remain
// (visible as a half-simplified, half-dense mesh: zero-error collapses are
// processed in memory order, a spatial sweep). Each pass targets
// count*keepRatio of the previous pass; iteration ends when a pass stops
// shrinking the mesh by more than 2%.
func simplifyToFixedPoint(verts []float32, count int, keepRatio float64, relError float32) ([]float32, int, float32) {
	out, achieved := verts, float32(0)
	for pass := 0; pass < 16; pass++ {
		target := int(float64(count) * keepRatio)
		out2, count2, err2 := meshopt.Simplify(out[:count*9], count, target, relError)
		if err2 > achieved {
			achieved = err2
		}
		shrunk := count2 < count*98/100
		out, count = out2, count2
		if !shrunk {
			break
		}
	}
	return out, count, achieved
}

// maxDetectedPlanes bounds the plane list so snapping stays O(verts·planes).
const maxDetectedPlanes = 32

// detectPlanes finds the dominant planes of a simplified triangle soup by
// clustering its largest triangles (whose planes are reliable — render noise
// barely tilts a triangle with mm-scale edges). Returns [nx ny nz d] rows
// with unit normals, largest total area first.
func detectPlanes(verts []float32, count int, extent float64) [][4]float64 {
	minArea := extent * extent * 1e-4 // 0.01% of extent² — ~1.7mm² on a 130mm part
	type cluster struct {
		n    [3]float64 // area-weighted normal sum (renormalized on read)
		d    float64    // area-weighted offset sum
		area float64
	}
	var clusters []*cluster
	for i := 0; i < count; i++ {
		off := i * 9
		ax, ay, az := float64(verts[off+0]), float64(verts[off+1]), float64(verts[off+2])
		ux, uy, uz := float64(verts[off+3])-ax, float64(verts[off+4])-ay, float64(verts[off+5])-az
		vx, vy, vz := float64(verts[off+6])-ax, float64(verts[off+7])-ay, float64(verts[off+8])-az
		cx := uy*vz - uz*vy
		cy := uz*vx - ux*vz
		cz := ux*vy - uy*vx
		twoArea := math.Sqrt(cx*cx + cy*cy + cz*cz)
		area := twoArea / 2
		if area < minArea {
			continue
		}
		nx, ny, nz := cx/twoArea, cy/twoArea, cz/twoArea
		d := nx*ax + ny*ay + nz*az
		matched := false
		for _, c := range clusters {
			cl := math.Sqrt(c.n[0]*c.n[0] + c.n[1]*c.n[1] + c.n[2]*c.n[2])
			dot := (nx*c.n[0] + ny*c.n[1] + nz*c.n[2]) / cl
			if dot > 0.99999 && math.Abs(d-c.d/c.area) < 0.02 { // ~0.25°, 0.02mm
				c.n[0] += nx * area
				c.n[1] += ny * area
				c.n[2] += nz * area
				c.d += d * area
				c.area += area
				matched = true
				break
			}
		}
		if !matched && len(clusters) < maxDetectedPlanes*4 {
			clusters = append(clusters, &cluster{n: [3]float64{nx * area, ny * area, nz * area}, d: d * area, area: area})
		}
	}
	sort.Slice(clusters, func(i, j int) bool { return clusters[i].area > clusters[j].area })
	if len(clusters) > maxDetectedPlanes {
		clusters = clusters[:maxDetectedPlanes]
	}
	planes := make([][4]float64, 0, len(clusters))
	for _, c := range clusters {
		l := math.Sqrt(c.n[0]*c.n[0] + c.n[1]*c.n[1] + c.n[2]*c.n[2])
		if l == 0 {
			continue
		}
		planes = append(planes, [4]float64{c.n[0] / l, c.n[1] / l, c.n[2] / l, c.d / c.area})
	}
	return planes
}

// snapToPlanes projects every vertex lying within eps of a detected plane
// exactly onto it (two rounds, so vertices near an edge shared by two planes
// converge onto the edge line). The rule is position-only — identical for
// every copy of a vertex in the soup — so no welding is needed and no cracks
// open. Movement is bounded by eps per plane, within the decimation budget.
// Returns the number of vertex copies moved.
func snapToPlanes(verts []float32, count int, planes [][4]float64, eps float64) int {
	moved := 0
	for i := 0; i < count*3; i++ {
		off := i * 3
		x, y, z := float64(verts[off]), float64(verts[off+1]), float64(verts[off+2])
		didMove := false
		for round := 0; round < 2; round++ {
			for _, p := range planes {
				dist := p[0]*x + p[1]*y + p[2]*z - p[3]
				if dist != 0 && math.Abs(dist) < eps {
					x -= p[0] * dist
					y -= p[1] * dist
					z -= p[2] * dist
					didMove = true
				}
			}
		}
		if didMove {
			verts[off] = float32(x)
			verts[off+1] = float32(y)
			verts[off+2] = float32(z)
			moved++
		}
	}
	return moved
}

// meshExtent returns the longest bounding-box axis of the mesh — the same
// scale meshoptimizer normalizes its error values by.
func meshExtent(mesh []*sdf.Triangle3) float64 {
	inf := math.Inf(1)
	minX, minY, minZ := inf, inf, inf
	maxX, maxY, maxZ := -inf, -inf, -inf
	for _, t := range mesh {
		for j := 0; j < 3; j++ {
			minX = math.Min(minX, t[j].X)
			minY = math.Min(minY, t[j].Y)
			minZ = math.Min(minZ, t[j].Z)
			maxX = math.Max(maxX, t[j].X)
			maxY = math.Max(maxY, t[j].Y)
			maxZ = math.Max(maxZ, t[j].Z)
		}
	}
	return math.Max(maxX-minX, math.Max(maxY-minY, maxZ-minZ))
}
