// Bolt assembly cookbook step 2: bolt with a smooth shank.
//
// ShankLength carves an unthreaded section out of TotalLength immediately
// below the head. Useful when the bolt rides through a smooth hole and
// you don't want the threads chewing the bore.
package main

import "github.com/snowbldr/fluent-sdfx/obj"

const thread = "M8x1.25"

func main() {
	bolt := obj.Bolt(obj.BoltParms{
		Thread:      thread,
		Style:       "hex",
		TotalLength: 30,
		ShankLength: 8,
	})

	bolt.STL("out.stl", 6.0)
}
