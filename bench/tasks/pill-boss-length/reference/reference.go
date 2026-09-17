// Package reference is the intended result: each boss 0.4 mm shorter at
// each END of its long (tangential) axis, rounded ends kept, width and
// height untouched.
// A flat ring with three stadium-section ("pill") bosses standing on its
// top face, each boss tangent to the ring (long axis along the ring's
// circumference, short axis radial).
package reference

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	ringOuterR = 15.0
	ringInnerR = 10.0
	ringH      = 4.0

	bossLen = 5.2  // was 6.0: 0.4 off each end
	bossW   = 2.0  // radial
	bossH   = 3.0  // Z, standing on the ring's top face
	bossR   = 12.5 // radius of the boss centre line
)

func Build() *solid.Solid {
	ring := solid.Cylinder(ringH, ringOuterR, 0).
		Cut(solid.Cylinder(ringH+2, ringInnerR, 0))

	// Stadium: long axis along X, fully rounded ends.
	boss := shape.Rect(v2.XY(bossLen, bossW), bossW/2).
		Extrude(bossH).
		Translate(v3.XYZ(0, bossR, ringH/2+bossH/2))

	return ring.Union(boss.RotateUnionZ(3, solid.RotateZMatrix(120)))
}
