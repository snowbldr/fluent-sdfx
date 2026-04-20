package main

import (
	"errors"
	"log"

	"github.com/snowbldr/fluent-sdfx/mesh"
	"github.com/snowbldr/fluent-sdfx/obj"
	flrender "github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func gyroidCube() *solid.Solid {
	l := 100.0
	k := l * 0.2

	gyroid := solid.Gyroid(v3.XYZ(k, k, k))
	box := solid.Box(v3.XYZ(l, l, l), 0)
	return box.Intersect(gyroid)
}

func gyroidSurface() *solid.Solid {
	l := 60.0
	k := l * 0.5

	s := solid.Gyroid(v3.XYZ(k, k, k)).Shell(k * 0.025)
	box := solid.Box(v3.XYZ(l, l, l), 0)
	s = box.Intersect(s)

	sphere := solid.Sphere(k * 0.15)
	d := l * 0.5
	s0 := sphere.Translate(v3.XYZ(d, d, d))
	s1 := sphere.Translate(v3.XYZ(-d, -d, -d))

	return s.Cut(s0.Union(s1))
}

func gyroidTeapot(cyclesPerSide int) (*solid.Solid, error) {
	if cyclesPerSide < 1 {
		return nil, errors.New("cycles per side should not be <= 0")
	}

	teapot := obj.ImportSTL("../../files/teapot.stl", 20, 3, 5)

	min := teapot.BoundingBox().Min
	max := teapot.BoundingBox().Max

	kX := (max.X - min.X) / float64(cyclesPerSide)
	kY := (max.Y - min.Y) / float64(cyclesPerSide)
	kZ := (max.Z - min.Z) / float64(cyclesPerSide)

	gyroid := solid.Gyroid(v3.XYZ(kX, kY, kZ))
	return teapot.Intersect(gyroid), nil
}

func main() {
	gyroidCube().STL("gyroid_cube.stl", 3.0)
	gyroidSurface().STL("gyroid_surface.stl", 1.5)

	s2, err := gyroidTeapot(10)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	s2.STL("gyroid_teapot.stl", 2.0)

	// Rendering to triangles and then saving the STL.
	m2 := mesh.ToTriangles(s2, flrender.NewMarchingCubesOctreeParallel(200))
	mesh.SaveSTL("gyroid_teapot_mesh.stl", m2)
}
