// Package start is the flange before the bolt holes are added: a disc of
// radius 30 and 6 mm thick sitting on the plate (z = 0 .. 6), a central bore
// of radius 8 straight through it, and an annular boss standing on its top
// face from r = 10 out to r = 14, 10 mm tall (z = 6 .. 16).
package start

import "github.com/snowbldr/fluent-sdfx/solid"

const (
	flangeR = 30.0
	flangeH = 6.0 // z = 0 .. 6
	boreR   = 8.0

	bossInnerR = 10.0
	bossOuterR = 14.0
	bossH      = 10.0 // z = 6 .. 16
)

func Build() *solid.Solid {
	flange := solid.Cylinder(flangeH, flangeR, 0).ZeroZ()

	// Annular boss: cut the core out while it is still centred, then stand
	// it on the flange's top face so the cutter never reaches the flange.
	boss := solid.Cylinder(bossH, bossOuterR, 0).
		Cut(solid.Cylinder(bossH+2, bossInnerR, 0)).
		BottomAt(flangeH)

	// Central bore, through the flange only.
	bore := solid.Cylinder(flangeH+2, boreR, 0).BottomAt(-1)

	return flange.Union(boss).Cut(bore)
}
