// Transforms: Translate moves a solid by a vector. The axis-specific
// helpers TranslateX / Y / Z avoid constructing a v3.Vec for the common case.
package main

import "github.com/snowbldr/fluent-sdfx/solid"

func main() {
	solid.Sphere(5).
		Union(
			solid.Sphere(5).TranslateX(12),
			solid.Sphere(5).TranslateY(12),
		).STL("out.stl", 5.0)
}
