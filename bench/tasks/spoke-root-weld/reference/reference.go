// Package reference: the hub with the spoke fully welded into it.
//
// The spoke is a flat slab (3 thick in Y) meeting a CURVED hub wall. A slab
// whose root face is placed ON the hub radius only touches the hub on the
// midplane (y = 0); at the slab edges (y = +/-1.5) the hub surface has
// already fallen back to x = sqrt(8^2 - 1.5^2) = 7.858, so a wedge-shaped
// notch is left open at the root corners. The fix used in the incident is to
// start the slab INSIDE the hub, past the chord it spans, so the whole root
// face is buried in hub material and the weld is solid across the full face.
//
// Incident: mcweed I-12 / U304-U308 ("the gusset seems like it's still not
// fully attached to the center pole. I see a little gap"). The model's
// midplane section "proved" the weld was sound; the user zoomed the model's
// own render and found the notch at the edges. The fix was
// rootR = sqrt(hubR^2 - (ws/2)^2) - 0.35.
package reference

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	hubR = 8.0  // hub outer radius
	hubH = 10.0 // hub height, z 0..10

	spokeT   = 3.0  // spoke thickness in Y
	spokeOut = 10.0 // how far the spoke sticks out past the hub surface, along +X
	spokeH   = 10.0 // full hub height

	// Root plane: buried inside the hub, past the chord the slab spans at
	// y = +/-spokeT/2 (that chord sits at x = 7.858). 6.0 is comfortably
	// inside it, so every point of the root face is hub material.
	spokeRootX = 6.0
	spokeTipX  = hubR + spokeOut // 18.0
	spokeLen   = spokeTipX - spokeRootX
)

func Build() *solid.Solid {
	hub := solid.Cylinder(hubH, hubR, 0).TranslateZ(hubH / 2)

	// Slab from x = 6 (inside the hub) out to x = 18 (10 past the surface).
	spoke := solid.Box(v3.XYZ(spokeLen, spokeT, spokeH), 0).
		Translate(v3.XYZ((spokeRootX+spokeTipX)/2, 0, spokeH/2))

	return hub.Union(spoke)
}
