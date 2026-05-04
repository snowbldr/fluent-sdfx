// 2D shapes: a circle.
package main

import "github.com/snowbldr/fluent-sdfx/shape"

func main() {
	shape.Circle(15).Extrude(1).STL("out.stl", 5)
}
