// Package start is the housing before the change: a cylindrical housing,
// outer radius 20, bore radius 16 (a 4 mm wall), 24 tall, z 0..24, with
//
//   - a boss on the inside of the wall: 10 wide tangentially (Y), 6 deep
//     radially, 12 tall (z 6..18), centred on the angle-0 axis (+X). Its
//     outer face sits flush against the bore wall at r = 16 and it projects
//     inward to r = 10. It is drawn 2 mm oversize into the wall so the
//     union has no coincident face.
//   - a 5 mm hole drilled radially on that same angle-0 axis at z = 12,
//     from outside the wall all the way through the wall and the boss.
package start

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	outerR = 20.0 // outside of the housing wall
	wallT  = 4.0  // wall thickness
	boreR  = outerR - wallT
	height = 24.0 // z 0..24

	bossOuterR = 16.0 // radius of the boss's outer face
	bossInnerR = 10.0 // radius of the boss's inner face
	bossW      = 10.0 // tangential, along Y at angle 0
	bossH      = 12.0 // along Z
	bossZ      = 12.0 // centre height

	holeD = 5.0  // radial hole diameter
	holeZ = 12.0 // radial hole height

	over = 2.0 // overshoot past each face
)

func Build() *solid.Solid {
	housing := solid.Cylinder(height, outerR, 0).TranslateZ(height / 2)
	bore := solid.Cylinder(height+2*over, boreR, 0).TranslateZ(height / 2)

	// Boss: a box from the boss's inner face out past the bore, so the
	// union buries its outer end in the wall.
	bossThick := bossOuterR + over - bossInnerR
	boss := solid.Box(v3.XYZ(bossThick, bossW, bossH), 0).
		Translate(v3.XYZ(bossInnerR+bossThick/2, 0, bossZ))

	// Radial hole along +X: from clear of the boss's inner face out to
	// clear of the housing's outer face.
	holeLen := (outerR + over) - (bossInnerR - over)
	hole := solid.Cylinder(holeLen, holeD/2, 0).
		RotateY(90).
		Translate(v3.XYZ(bossInnerR-over+holeLen/2, 0, holeZ))

	return housing.Cut(bore).Union(boss).Cut(hole)
}
