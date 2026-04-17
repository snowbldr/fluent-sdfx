package solid

import (
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func toSDFTris(triangles []*v3.Triangle3) []*sdf.Triangle3 {
	out := make([]*sdf.Triangle3, len(triangles))
	for i, t := range triangles {
		st := t.SDF()
		out[i] = &st
	}
	return out
}

// Mesh builds an SDF from a triangle mesh (fast, r-tree accelerated).
func Mesh(triangles []*v3.Triangle3) *Solid {
	s, err := sdf.Mesh3D(toSDFTris(triangles))
	if err != nil {
		panic(err)
	}
	return &Solid{s}
}

// MeshSlow builds an SDF from a triangle mesh using the slow, brute-force method.
func MeshSlow(triangles []*v3.Triangle3) *Solid {
	s, err := sdf.Mesh3DSlow(toSDFTris(triangles))
	if err != nil {
		panic(err)
	}
	return &Solid{s}
}

// Voxel samples the solid into a voxel grid and returns an SDF3 that trilinearly interpolates it.
// meshCells is the resolution on the longest axis; progress receives 0-1 sampling progress (may be nil).
func (s *Solid) Voxel(meshCells int, progress chan float64) *Solid {
	return &Solid{sdf.NewVoxelSDF3(s.SDF3, meshCells, progress)}
}
