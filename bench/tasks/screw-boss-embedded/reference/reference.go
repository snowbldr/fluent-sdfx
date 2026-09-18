// Package reference: the tube with an M2 screw boss EMBEDDED in its wall.
//
// The boss is a Ø4.4 cylinder lying horizontally along Y, its axis at
// x = 8.5 (inside the 8..10 wall band) and z = +12. Because the wall is
// curved and the boss is straight, the raw boss pushes out through the
// tube's outer surface as |y| grows — at the ends it reaches r = 12.3, more
// than 2 mm proud. So the boss is INTERSECTED with the tube's outer cylinder
// before it is unioned on: every bit of it that would stand proud is sliced
// off by the parent cylinder, exactly the construction the user ended up
// dictating. What is left adds material only on the inside, thickening the
// wall locally so the Ø2 screw bore has meat around it.
//
// Incident: mcweed U19/U43/U46 "what are those bumps on the outsides of the
// cylinders?... those should be horizontal to the tubes, and embedded in the
// tubes at the corners", "still popping out of the cylinder in a not nice
// way... The extra material needs to be sliced off with a cut from a tube or
// an intersection with a cylinder". Six turns over four days.
package reference

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	outerR = 10.0
	innerR = 8.0
	tubeH  = 30.0

	bossR   = 2.2  // Ø4.4 boss
	bossX   = 8.5  // axis sits in the 8..10 wall band
	bossZ   = 12.0 // near the top of the tube
	bossLen = 12.0 // along Y, centred on y = 0
	boreR   = 1.0  // Ø2 clearance bore for the M2 screw
	boreLen = 16.0 // longer than the boss so both ends exit the body
)

func Build() *solid.Solid {
	tube := solid.Cylinder(tubeH, outerR, 0).Cut(solid.Cylinder(tubeH+2, innerR, 0))

	// The parent cylinder, used as a trimming envelope rather than as body.
	envelope := solid.Cylinder(tubeH, outerR, 0)

	// Boss: axis along Y (Cylinder is Z-axial, so roll it 90 about X).
	boss := solid.Cylinder(bossLen, bossR, 0).RotateX(90).
		Translate(v3.XZ(bossX, bossZ)).
		Intersect(envelope)

	bore := solid.Cylinder(boreLen, boreR, 0).RotateX(90).
		Translate(v3.XZ(bossX, bossZ))

	return tube.Union(boss).Cut(bore)
}
