package obj

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/sdf"
	"github.com/snowbldr/fluent-sdfx/solid"
)

// ImportTriMesh builds an SDF3 from triangle mesh data with a BVH using the
// given neighbor/child counts.
func ImportTriMesh(mesh []*sdf.Triangle3, numNeighbors, minChildren, maxChildren int) *solid.Solid {
	return solid.Wrap(obj.ImportTriMesh(mesh, numNeighbors, minChildren, maxChildren))
}

// ImportSTL loads a binary or ASCII STL file and returns an SDF3 for the mesh.
func ImportSTL(path string, numNeighbors, minChildren, maxChildren int) *solid.Solid {
	s, err := obj.ImportSTL(path, numNeighbors, minChildren, maxChildren)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}
