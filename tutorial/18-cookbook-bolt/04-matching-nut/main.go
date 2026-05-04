// Bolt assembly cookbook step 4: a matching nut.
//
// obj.Nut takes a thread name and a style. The internal thread is sized
// to mate with a bolt of the same name; Tolerance adds a positive
// clearance to the internal thread radius for a printable fit.
package main

import "github.com/snowbldr/fluent-sdfx/obj"

const thread = "M8x1.25"

func main() {
	obj.Nut(obj.NutParms{
		Thread:    thread,
		Style:     "hex",
		Tolerance: 0.1,
	}).STL("out.stl", 6.0)
}
