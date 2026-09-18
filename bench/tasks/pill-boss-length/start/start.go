// Package start is the part as it existed before the change was requested.
// A flat ring with three stadium-section ("pill") bosses standing on its
// top face, each boss tangent to the ring (long axis along the ring's
// circumference, short axis radial).
package start

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

	bossLen = 6.0  // along the ring circumference (tangential)
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
