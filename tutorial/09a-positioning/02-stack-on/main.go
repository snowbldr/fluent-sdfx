// Positioning: stack a cap on top of a body in one chain.
//
// `cap.OnTopOf(body.Top())` is sugar for `cap.Bottom().Above(body.Top())` —
// it places the cap so its bottom anchor lands on the body's top anchor.
// The trailing `.Union()` finalizes the resulting Placement, unioning
// cap into body. Pass a gap (`OnTopOf(body.Top(), 2)`) for clearance.
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	body := solid.Box(v3.XYZ(20, 20, 10), 1).BottomAt(0)
	cap := solid.Sphere(8)

	cap.OnTopOf(body.Top()).Union().STL("out.stl", 5.0)
}
