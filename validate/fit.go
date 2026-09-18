package validate

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/snowbldr/fluent-sdfx/mesh"
	"github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Assembly inspection: write the STLs a human opens to look at a fit, and
// print the numbers a model reads instead of looking.

// fitMaxVolumeSamples caps the grid-occupancy volume estimate (Manifest) so
// a big part at a fine cellsPerMM cannot turn a summary table into a
// multi-minute render. Past the cap the grid is coarsened uniformly and the
// estimate gets correspondingly rougher.
const fitMaxVolumeSamples = 4e6

// Part is a named solid in an assembly, as passed to FitCheck and Manifest.
type Part struct {
	Name  string
	Solid *solid.Solid
}

// FitCheck writes an inspection set for an assembly and reports every
// pairwise clearance. It renders <dir>/<name>_assembly.stl (the union of
// all parts) and, unless normal is the zero vector, <dir>/<name>_cutaway.stl
// (that union cut by the plane through point, keeping the side the normal
// points to). The returned map is keyed "a>b" and holds Clearance(a, b):
// b's field sampled on a's surface, so negative values mean a penetrates b.
// Parts must have unique non-empty names.
func FitCheck(dir, name string, parts []Part, point, normal v3.Vec, cellsPerMM float64) (map[string]Gap, error) {
	if len(parts) == 0 {
		return nil, fmt.Errorf("validate.FitCheck: no parts")
	}
	seen := make(map[string]struct{}, len(parts))
	solids := make([]*solid.Solid, len(parts))
	for i, p := range parts {
		if p.Name == "" {
			return nil, fmt.Errorf("validate.FitCheck: part %d has no name", i)
		}
		if p.Solid == nil {
			return nil, fmt.Errorf("validate.FitCheck: part %q has a nil solid", p.Name)
		}
		if _, dup := seen[p.Name]; dup {
			return nil, fmt.Errorf("validate.FitCheck: duplicate part name %q", p.Name)
		}
		seen[p.Name] = struct{}{}
		solids[i] = p.Solid
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("validate.FitCheck: %w", err)
	}

	assembly := solid.UnionAll(solids...)
	if err := fitSaveSTL(filepath.Join(dir, name+"_assembly.stl"), assembly, cellsPerMM); err != nil {
		return nil, err
	}
	if normal.Length() > 0 {
		cut := assembly.CutPlane(point, normal)
		if err := fitSaveSTL(filepath.Join(dir, name+"_cutaway.stl"), cut, cellsPerMM); err != nil {
			return nil, err
		}
	}

	gaps := make(map[string]Gap, len(parts)*(len(parts)-1))
	for i, a := range parts {
		for j, b := range parts {
			if i == j {
				continue
			}
			gaps[a.Name+">"+b.Name] = Clearance(a.Solid, b.Solid, cellsPerMM)
		}
	}
	return gaps, nil
}

// fitSaveSTL renders s at cellsPerMM and writes it to path, converting the
// mesh writer's panic into an error so FitCheck can report it.
func fitSaveSTL(path string, s *solid.Solid, cellsPerMM float64) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("validate: writing %s: %v", path, r)
		}
	}()
	tris := mesh.ToTriangles(s, render.NewMarchingCubesOctreeParallel(solid.CellsFor(s, cellsPerMM)))
	mesh.SaveSTL(path, tris)
	return nil
}

// Manifest summarises parts as a markdown table — name, bounding-box min and
// max, size, volume and body count — for a build log, a PR comment, or a
// model deciding whether the part it just wrote is the size it meant. The
// volume is estimated by grid occupancy (the fraction of cell centres at
// cellsPerMM whose field is negative), so it is approximate and cheap, not
// the mesh volume; bodies is Components at the same density.
func Manifest(parts []Part, cellsPerMM float64) string {
	var b strings.Builder
	b.WriteString("| name | min | max | size | volume mm³ | bodies |\n")
	b.WriteString("| --- | --- | --- | --- | --- | --- |\n")
	for _, p := range parts {
		if p.Solid == nil {
			fmt.Fprintf(&b, "| %s | — | — | — | — | — |\n", p.Name)
			continue
		}
		bb := p.Solid.Bounds()
		size := bb.Size()
		fmt.Fprintf(&b, "| %s | (%.3f, %.3f, %.3f) | (%.3f, %.3f, %.3f) | %.3f × %.3f × %.3f | %.1f | %d |\n",
			p.Name,
			bb.Min.X, bb.Min.Y, bb.Min.Z,
			bb.Max.X, bb.Max.Y, bb.Max.Z,
			size.X, size.Y, size.Z,
			fitVolume(p.Solid, cellsPerMM),
			Components(p.Solid, cellsPerMM))
	}
	return b.String()
}

// fitVolume estimates the volume of s in mm³ by counting the cell centres
// of a regular grid over its bounding box (spacing 1/cellsPerMM) that read
// inside the material. The grid is coarsened uniformly if it would exceed
// fitMaxVolumeSamples points.
func fitVolume(s *solid.Solid, cellsPerMM float64) float64 {
	bb := s.Bounds()
	size := bb.Size()
	if size.X <= 0 || size.Y <= 0 || size.Z <= 0 || cellsPerMM <= 0 {
		return 0
	}
	step := 1 / cellsPerMM
	if total := size.X * size.Y * size.Z / (step * step * step); total > fitMaxVolumeSamples {
		step *= math.Cbrt(total / fitMaxVolumeSamples)
	}
	nx := int(math.Max(1, math.Floor(size.X/step)))
	ny := int(math.Max(1, math.Floor(size.Y/step)))
	nz := int(math.Max(1, math.Floor(size.Z/step)))
	// Cell centres of an nx×ny×nz grid that exactly tiles the bounding box.
	dx, dy, dz := size.X/float64(nx), size.Y/float64(ny), size.Z/float64(nz)
	cellVol := dx * dy * dz

	pts := make([]v3.Vec, nx*ny)
	var vol float64
	for k := 0; k < nz; k++ {
		z := bb.Min.Z + (float64(k)+0.5)*dz
		for i := 0; i < nx; i++ {
			x := bb.Min.X + (float64(i)+0.5)*dx
			for j := 0; j < ny; j++ {
				pts[i*ny+j] = v3.XYZ(x, bb.Min.Y+(float64(j)+0.5)*dy, z)
			}
		}
		for _, d := range evalAll(s, pts) {
			if d < 0 {
				vol += cellVol
			}
		}
	}
	return vol
}
