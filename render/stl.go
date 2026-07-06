package render

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

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
// Optional decimate values: decimate[0] (0-1) is the fraction of triangles to
// remove — 0.1 removes 10% (keeps 90%); 0.9 removes 90% (keeps 10%); 0
// disables decimation. decimate[1] is the maximum surface deviation in model
// units (mm) that decimation may introduce (default DefaultDecimateMaxError);
// the triangle target is abandoned early rather than exceed it.
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

	target := int(float64(n) * keepRatio)
	out, count, achieved := meshopt.Simplify(verts, n, target, float32(relError))

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
	fmt.Printf(" → %d after decimation (kept %.0f%%, deviation ≤ %.3fmm of %.3fmm budget)",
		count, 100*float64(count)/float64(n), float64(achieved)*extent, maxErrorMM)
	return result
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
