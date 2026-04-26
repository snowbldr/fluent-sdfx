package solid

import (
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/sdfx/sdf"
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
