// Package reference: a cylindrical core with a single-start RIGHT-HAND
// helical rib wrapped around it.
//
// Handedness, checked by probing and not assumed: shape.SweepHelix produces
// a RIGHT-hand helix as built. Sweeping a circle at radius 8 for 2 turns
// over 24 mm puts rib material at angle 0 deg at z = 0, +90 deg at z = +3,
// 180 deg at z = +6, 270 deg at z = +9 — the angle ADVANCES
// counter-clockwise (seen from +Z) as z rises, which is the right-hand
// screw rule, i.e. clockwise as the thread runs away from you downwards.
//
// SweepHelix(radius, turns, height, flatEnds) sweeps the profile with the
// profile's own plane radial: the swept circle's centre line sits at
// `radius`, so with radius = coreR the 2 mm rib stands 1 mm proud of the
// core. The result is centred on z = 0. flatEnds=false clips the sweep to
// z = -12 .. +12 so the rib ends flush with the core faces; flatEnds=true
// lets the profile run past the end planes by up to its radius (to +-13).
package reference

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
)

const (
	coreR = 8.0  // cylinder core radius
	coreH = 24.0 // z = -12 .. +12

	ribProfileR = 1.0  // 2 mm diameter round rib
	ribTurns    = 2.0  // 2 turns over the 24 mm height
	ribHeight   = 24.0 // => 12 mm pitch, one turn per 12 mm
)

func Build() *solid.Solid {
	core := solid.Cylinder(coreH, coreR, 0)

	// Right-hand, single start: SweepHelix's native sense.
	rib := shape.Circle(ribProfileR).SweepHelix(coreR, ribTurns, ribHeight, false)

	return core.Union(rib)
}
