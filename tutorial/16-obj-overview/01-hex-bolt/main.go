// Parametric helpers: a hex-head bolt via obj.Bolt.
//
// Thread is one of the standard names in sdfx's thread table — "M8x1.25",
// "unc_1/4", etc. See sdf.ThreadLookup or shape.ThreadLookup.
package main

import "github.com/snowbldr/fluent-sdfx/obj"

func main() {
	obj.Bolt(obj.BoltParms{
		Thread:      "M10x1.5",
		Style:       "hex",
		TotalLength: 30,
		ShankLength: 5,
	}).STL("out.stl", 6.0)
}
