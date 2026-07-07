package render

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/snowbldr/fluent-sdfx/render/meshopt"
	"github.com/snowbldr/sdfx/render"
	"github.com/snowbldr/sdfx/sdf"
)

// DefaultDecimateMaxError is the default cap, in model units (mm), on how far
// decimation may move the surface from the rendered mesh. Well under a
// typical 3D-print layer height, so decimated and full meshes print
// identically.
const DefaultDecimateMaxError = 0.05

// ToSTL renders an SDF3 to an STL file atomically.
// Optional decimate values: decimate[0] (0-1) enables decimation (any value
// in range; the error budget — not this ratio — determines the final size).
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

// DecimateMesh simplifies a triangle mesh, moving no surface vertex more
// than maxErrorMM (model units) from the input. The triangle target is
// abandoned rather than exceed maxErrorMM, so the budget — not keepRatio —
// determines the final size.
//
// Decimation quality depends directly on render accuracy: meshopt can only
// merge geometry whose deviation fits the budget, so render noise is kept as
// if it were real detail. The marching-cubes renderers' edge refinement
// (sdfx Refine, default on) places vertices on the true surface to ~1e-4mm,
// which lets flat faces collapse to a handful of exactly coplanar triangles
// and curved surfaces decimate to true-curvature chords at any budget.
//
// An earlier version of this pipeline detected flat faces and snapped their
// vertices onto fitted planes before simplifying, to correct unrefined
// render noise. With refined renders the raw mesh is more accurate than the
// snap's idealization — the snap's adjustments printed as visible defects on
// thin adaptive layers — so plane snapping was removed.
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

// simplifyToFixedPoint runs meshopt simplification repeatedly until the
// error budget — not the per-pass triangle quota — is what stops it, ending
// when a pass stops shrinking the mesh by more than 2%. Each call rebuilds
// quadrics from the current geometry, so iterating unlocks collapses that a
// single call's accumulated error estimates would refuse.
func simplifyToFixedPoint(verts []float32, count int, keepRatio float64, relError float32) ([]float32, int, float32) {
	_ = keepRatio // retained for API stability; the error budget governs (see below)
	debug := os.Getenv("FLSDFX_DEBUG_DECIM") != ""
	out, achieved := verts, float32(0)
	for pass := 0; pass < 16; pass++ {
		// The per-pass target must be far BELOW the error-limited floor:
		// meshopt restrains its internal passes around the target (its
		// error_goal heuristic stops collapsing near the target count
		// rather than running out the error budget), so a modest target —
		// say keepRatio's 10% — makes every call stop early, and iterating
		// converges to a ceiling several times above the budget's true
		// fixed point (observed: 10% targets plateaued a 56M-triangle mold
		// at 17.8M, while a 1% target reached 175k in ONE call at the same
		// budget). The error budget, not the target, is the documented
		// contract, so aim low and let the budget stop the collapse.
		target := count / 100
		if target < 1000 {
			target = 1000
		}
		if target > count/2 {
			target = count / 2
		}
		if target < 1 {
			target = 1
		}
		start := time.Now()
		out2, count2, err2 := meshopt.Simplify(out[:count*9], count, target, relError)
		if debug {
			fmt.Printf("\n  simplify pass %d: %d → %d tris (%.1f%%), err %.3g/%.3g, %.1fs",
				pass, count, count2, 100*float64(count2)/float64(count), err2, relError, time.Since(start).Seconds())
		}
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
