package render

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/snowbldr/fluent-sdfx/render/meshopt"
	"github.com/snowbldr/sdfx/render"
	"github.com/snowbldr/sdfx/sdf"
)

// ToSTL renders an SDF3 to an STL file atomically.
// Optional decimate (0-1) is the fraction of triangles to remove:
// 0.1 removes 10% (keeps 90%); 0.9 removes 90% (keeps 10%). 0 disables decimation.
func ToSTL(s SDF3, path string, r render.Render3, decimate ...float64) {
	a := sdf3Adapter{s}
	fmt.Printf("rendering %s (%s)\n", path, r.Info(a))

	mesh := render.ToTriangles(a, r)
	fmt.Printf("  %d triangles", len(mesh))

	removeRatio := 0.0
	if len(decimate) > 0 {
		removeRatio = decimate[0]
	}
	if removeRatio > 0 && removeRatio < 1 {
		mesh = decimateMesh(mesh, 1-removeRatio)
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

// decimateMesh simplifies a triangle mesh down to keepRatio of its original size.
func decimateMesh(mesh []*sdf.Triangle3, keepRatio float64) []*sdf.Triangle3 {
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

	target := int(float64(n) * keepRatio)
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
	fmt.Printf(" → %d after decimation (kept %.0f%%)", count, keepRatio*100)
	return result
}
