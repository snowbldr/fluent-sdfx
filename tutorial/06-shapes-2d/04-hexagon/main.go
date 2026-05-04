// 2D shapes: a regular hexagon by inscribed radius.
package main

import "github.com/snowbldr/fluent-sdfx/shape"

func main() {
	shape.Hexagon(12).Extrude(1).STL("out.stl", 5)
}
