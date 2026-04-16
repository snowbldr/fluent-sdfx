package render

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	"github.com/snowbldr/fluent-sdfx/render/meshopt"
)

// ToSTL renders an SDF3 to an STL file atomically.
// Optional factor (0-1) controls mesh decimation: 0.5 = keep 50% of triangles.
func ToSTL(s sdf.SDF3, path string, r render.Render3, factor ...float64) {
	fmt.Printf("rendering %s (%s)\n", path, r.Info(s))

	mesh := render.ToTriangles(s, r)
	fmt.Printf("  %d triangles", len(mesh))

	decimateFactor := 0.0
	if len(factor) > 0 {
		decimateFactor = factor[0]
	}
	if decimateFactor > 0 && decimateFactor < 1 {
		mesh = decimateMesh(mesh, 1-decimateFactor)
	}
	fmt.Println()

	dir := filepath.Dir(path)
	if dir == "" {
		dir = "."
	}
	tmp, err := os.CreateTemp(dir, ".stl-tmp-*")
	if err != nil {
		fmt.Printf("error creating temp file: %s\n", err)
		return
	}
	tmpPath := tmp.Name()
	tmp.Close()

	if err := render.SaveSTL(tmpPath, mesh); err != nil {
		fmt.Printf("error writing STL: %s\n", err)
		os.Remove(tmpPath)
		return
	}

	if err := os.Rename(tmpPath, path); err != nil {
		fmt.Printf("error renaming STL: %s\n", err)
		os.Remove(tmpPath)
		return
	}
}

func decimateMesh(mesh []*sdf.Triangle3, factor float64) []*sdf.Triangle3 {
	n := len(mesh)
	verts := make([]float32, n*9)
	for i, t := range mesh {
		off := i * 9
		verts[off+0] = float32(t[0].X)
		verts[off+1] = float32(t[0].Y)
		verts[off+2] = float32(t[0].Z)
		verts[off+3] = float32(t[1].X)
		verts[off+4] = float32(t[1].Y)
		verts[off+5] = float32(t[1].Z)
		verts[off+6] = float32(t[2].X)
		verts[off+7] = float32(t[2].Y)
		verts[off+8] = float32(t[2].Z)
	}

	target := int(float64(n) * factor)
	out, count := meshopt.Simplify(verts, n, target, 0.1)

	result := make([]*sdf.Triangle3, count)
	for i := 0; i < count; i++ {
		off := i * 9
		tri := sdf.Triangle3{
			{X: float64(out[off+0]), Y: float64(out[off+1]), Z: float64(out[off+2])},
			{X: float64(out[off+3]), Y: float64(out[off+4]), Z: float64(out[off+5])},
			{X: float64(out[off+6]), Y: float64(out[off+7]), Z: float64(out[off+8])},
		}
		result[i] = &tri
	}
	fmt.Printf(" → %d after decimation (%.0f%%)", count, factor*100)
	return result
}
