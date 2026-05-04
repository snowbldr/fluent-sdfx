// 2D → 3D: Screw revolves a thread profile around Z to make screw threads.
//
// height: total axial length. start: starting angular offset (radians).
// pitch: axial distance per turn. num: number of starts (parallel
// helices); 1 for a normal screw.
//
// shape.AcmeThread is a stock thread profile for the common case; for a
// custom thread, build a closed shape.Polygon and offset it from the axis.
package main

import "github.com/snowbldr/fluent-sdfx/shape"

func main() {
	shape.AcmeThread(5, 3).Screw(20, 0, 3, 1).STL("out.stl", 8.0)
}
