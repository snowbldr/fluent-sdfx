// Bolt assembly cookbook step 3: try a different thread standard.
//
// "Knurl" head style produces a knurled cap instead of hex; useful for
// printable thumbscrews. Thread "M5x0.8" is a common metric small bolt.
package main

import "github.com/snowbldr/fluent-sdfx/obj"

func main() {
	obj.Bolt(obj.BoltParms{
		Thread:      "M5x0.8",
		Style:       "knurl",
		TotalLength: 20,
		ShankLength: 4,
	}).STL("out.stl", 8.0)
}
