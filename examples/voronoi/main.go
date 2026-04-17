package main

import (
	"log"

	"github.com/snowbldr/fluent-sdfx/mesh"
	"github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/shape"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	"github.com/snowbldr/fluent-sdfx/vec/v2i"
)

func main() {
	// random set of vertices
	b := shape.NewBox2(v2.XY(0, 0), v2.XY(20, 20))
	pts := b.RandomSet(20)
	pixels := v2i.XY(800, 800)
	k := 1.5
	path := "voronoi.png"

	// use a 0-radius circle as a point
	point := shape.Circle(0.0)

	// build an SDF for the points
	var s0 *shape.Shape
	for i := range pts {
		p := point.Translate(pts[i])
		if s0 == nil {
			s0 = p
		} else {
			s0 = s0.Union(p)
		}
	}

	bb := s0.BoundingBox().ScaleAboutCenter(k)
	log.Printf("rendering %s (%dx%d)\n", path, pixels.X, pixels.Y)
	d, err := render.NewPNG(path, bb, pixels)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	d.RenderSDF2(s0)

	ts := mesh.Delaunay2d(pts)
	for _, t := range ts {
		d.Triangle(t.ToTriangle2(pts))
	}

	d.Save()
}
