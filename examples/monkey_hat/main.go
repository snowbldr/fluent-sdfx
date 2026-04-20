package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func monkeyWithHat() *solid.Solid {
	// create the SDF from the mesh (a modified Suzanne from Blender with 366 faces)
	monkey := obj.ImportSTL("../../files/monkey.stl", 20, 3, 5)

	// build the hat
	hatHeight := 0.5
	hat := solid.Cylinder(hatHeight, 0.6, 0)
	edge := solid.Cylinder(hatHeight*0.4, 1, 0).Translate(v3.Z(-hatHeight / 2))
	fullHat := hat.Union(edge)

	// put the hat on the monkey
	fullHat = fullHat.Translate(v3.YZ(0.15, 1))
	result := monkey.Union(fullHat)

	// Cache the mesh full SDF3 hierarchy for faster evaluation via voxel interpolation.
	return result.Voxel(64, nil)
}

func main() {
	monkeyWithHat().STL("monkey-out.stl", 1.28)
}
