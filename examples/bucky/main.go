package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

type edge [2]int

var edges = []edge{
	{1, 4}, {4, 8}, {8, 10}, {10, 6}, {6, 1},
	{0, 4}, {4, 9}, {9, 11}, {11, 6}, {6, 0},
	{3, 7}, {7, 10}, {10, 8}, {8, 5}, {5, 3},
	{5, 9}, {9, 11}, {11, 7}, {7, 2}, {2, 5},
	{0, 1}, {1, 9}, {9, 5}, {5, 8}, {8, 0},
	{4, 9}, {9, 3}, {3, 2}, {2, 8}, {8, 4},
	{1, 0}, {0, 10}, {10, 7}, {7, 11}, {11, 1},
	{2, 3}, {3, 11}, {11, 6}, {6, 10}, {10, 2},
	{0, 10}, {10, 2}, {2, 5}, {5, 4}, {4, 0},
	{1, 4}, {4, 5}, {5, 3}, {3, 11}, {11, 1},
	{7, 2}, {2, 8}, {8, 0}, {0, 6}, {6, 7},
	{1, 6}, {6, 7}, {7, 3}, {3, 9}, {9, 1},
}

const phi = 1.618033988749895

var vertex = []v3.Vec{
	v3.XYZ(1, phi, 0),
	v3.XYZ(-1, phi, 0),
	v3.XYZ(1, -phi, 0),
	v3.XYZ(-1, -phi, 0),
	v3.XYZ(0, 1, phi),
	v3.XYZ(0, -1, phi),
	v3.XYZ(0, 1, -phi),
	v3.XYZ(0, -1, -phi),
	v3.XYZ(phi, 0, 1),
	v3.XYZ(-phi, 0, 1),
	v3.XYZ(phi, 0, -1),
	v3.XYZ(-phi, 0, -1),
}

func icosahedron() *solid.Solid {
	r0 := phi * 0.05
	r1 := r0 * 2.0

	k := obj.ArrowParms{
		Axis:  [2]float64{0, r0},
		Head:  [2]float64{0, r1},
		Tail:  [2]float64{0, r1},
		Style: "b.",
	}

	var bb *solid.Solid
	for _, e := range edges {
		head := vertex[e[0]]
		tail := vertex[e[1]]
		arrow := obj.DirectedArrow3D(k, head, tail)
		if bb == nil {
			bb = arrow
		} else {
			bb = bb.Union(arrow)
		}
	}

	return bb
}

func main() {
	s := icosahedron()
	s.STL("icosahedron.stl", 3.0)
	s.ThreeMF("icosahedron.3mf", 3.0)
}
