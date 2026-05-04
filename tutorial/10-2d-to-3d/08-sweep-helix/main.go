// 2D → 3D: SweepHelix sweeps a 2D profile along a helical path.
//
// radius is the helix radius; turns the number of revolutions; height the
// total axial length. flatEnds=true caps the sweep with flat planes
// perpendicular to Z.
package main

import "github.com/snowbldr/fluent-sdfx/shape"

func main() {
	shape.Circle(1.5).SweepHelix(8, 4, 30, true).STL("out.stl", 6.0)
}
