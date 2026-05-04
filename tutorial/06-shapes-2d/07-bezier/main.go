// 2D shapes: a closed bezier curve via NewBezier with slope handles.
//
// HandleFwd(theta, r) sets the forward control-point handle at angle theta
// (degrees) and length r. A symmetric handle on each cardinal vertex
// produces a smooth, rounded blob.
package main

import "github.com/snowbldr/fluent-sdfx/shape"

func main() {
	b := shape.NewBezier()
	b.Add(-15, 0).HandleFwd(90, 6)
	b.Add(0, 12).HandleFwd(0, 6)
	b.Add(15, 0).HandleFwd(-90, 6)
	b.Add(0, -12).HandleFwd(180, 6)
	b.Close()

	b.Build().Extrude(1).STL("out.stl", 5)
}
