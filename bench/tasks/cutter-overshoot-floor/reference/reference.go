// Package reference is the intended result: six radial grooves in the TOP
// face only, 1 mm deep and 2 mm wide, running from radius 5 out to the rim.
// The 3 mm of disc below each groove is untouched and the bottom face is
// still a solid, unbroken disc.
//
// The rule the cutter obeys: overshoot only in directions where there is
// nothing to protect. Radially the cutter runs past the rim (r = 17) so the
// groove breaks out cleanly; vertically it is pinned with BottomAt at
// z = discH - grooveDeep and extends UPWARDS past the top face, so the floor
// under the groove cannot be touched no matter how tall the cutter is.
package reference

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	discR = 15.0
	discH = 4.0 // z = 0 .. 4

	grooveCount  = 6
	grooveW      = 2.0 // across the groove
	grooveDeep   = 1.0 // into the top face: z = 3 .. 4
	grooveStartR = 5.0 // grooves begin 5 mm out from the centre
	grooveOverR  = 2.0 // radial overshoot past the rim
	grooveOverZ  = 2.0 // vertical overshoot, ABOVE the top face only

	grooveLen = discR + grooveOverR - grooveStartR
)

func Build() *solid.Solid {
	disc := solid.Cylinder(discH, discR, 0).ZeroZ()

	// One groove cutter on +X: long enough to break out of the rim, and
	// bottomed exactly grooveDeep below the top face.
	cutter := solid.Box(v3.XYZ(grooveLen, grooveW, grooveDeep+grooveOverZ), 0).
		TranslateX(grooveStartR + grooveLen/2).
		BottomAt(discH - grooveDeep)

	return disc.Cut(cutter.RotateCopyZ(grooveCount))
}
