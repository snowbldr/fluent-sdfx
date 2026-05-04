// 2D shapes: a 5-pointed star with outer/inner radius and point count.
package main

import "github.com/snowbldr/fluent-sdfx/shape"

func main() {
	shape.Star(15, 7, 5).Extrude(1).STL("out.stl", 5)
}
