package main

import (
	"github.com/snowbldr/fluent-sdfx/mesh"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	side := 30.0

	a := v3.XYZ(0, 0, 0)
	b := v3.XYZ(side, 0, 0)
	c := v3.XYZ(0, side, 0)
	d := v3.XYZ(0, 0, side)

	t1 := &mesh.Triangle3{a, b, d}
	t2 := &mesh.Triangle3{a, c, b}
	t3 := &mesh.Triangle3{a, d, c}
	t4 := &mesh.Triangle3{b, c, d}

	mesh.SaveSTL("simple.stl", []*mesh.Triangle3{t1, t2, t3, t4})
}
