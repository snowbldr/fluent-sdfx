// Bolt assembly cookbook step 5: a washer.
//
// InnerRadius is sized to clear the bolt; OuterRadius gives the load-
// spreading face. Remove > 0 cuts a wedge for a "split" / lock-washer
// look; leave at 0 for a plain washer.
package main

import "github.com/snowbldr/fluent-sdfx/obj"

func main() {
	obj.Washer3D(obj.WasherParms{
		Thickness:   1.5,
		InnerRadius: 4.5, // ~M8 clearance
		OuterRadius: 8.5,
		Remove:      0,
	}).STL("out.stl", 8.0)
}
