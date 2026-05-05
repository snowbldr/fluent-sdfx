// Lantern cookbook step 2: hollow out the tea light pocket.
//
// `body.Cut(pocket.Top().On(body.Top()).Solid())` — the chain's subject
// is the body (so the body is kept), and the argument is the pocket
// positioned so its top aligns with the body's top. No bbox math.
package main

import "github.com/snowbldr/fluent-sdfx/solid"

const (
	bodyHeight  = 50.0
	bodyRadius  = 25.0
	wallThick   = 5.0
	pocketDepth = 40.0
)

func main() {
	body := solid.Cylinder(bodyHeight, bodyRadius, 4)
	pocket := solid.Cylinder(pocketDepth, bodyRadius-wallThick, 0)

	body.Cut(pocket.Top().On(body.Top()).Solid()).STL("out.stl", 4.0)
}
