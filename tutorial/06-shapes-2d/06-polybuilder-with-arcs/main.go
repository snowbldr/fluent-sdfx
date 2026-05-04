// 2D shapes: a fluent polygon builder with smoothed and chamfered corners.
//
// NewPoly accepts vertices via Add(x, y). Each returns a PolyVertex you can
// modify with .Smooth(r, facets) for a fillet, .Chamfer(c) for a 45° cut,
// or .Arc(r, facets) to replace the prior edge with an arc.
package main

import "github.com/snowbldr/fluent-sdfx/shape"

func main() {
	p := shape.NewPoly()
	p.Add(-15, -10)
	p.Add(15, -10).Smooth(3, 6)
	p.Add(15, 10).Chamfer(4)
	p.Add(-15, 10).Smooth(3, 6)
	p.Close()

	p.Build().Extrude(1).STL("out.stl", 5)
}
