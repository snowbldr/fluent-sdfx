// Bolt assembly cookbook step 1: a single hex bolt.
//
// obj.Bolt builds the entire fastener — head, shank, threads — in one call.
// Every later step in this cookbook adds one more part beside it.
package main

import "github.com/snowbldr/fluent-sdfx/obj"

const thread = "M8x1.25"

func main() {
	bolt := obj.Bolt(obj.BoltParms{
		Thread:      thread,
		Style:       "hex",
		TotalLength: 25,
		ShankLength: 0,
	})

	bolt.STL("out.stl", 6.0)
}
