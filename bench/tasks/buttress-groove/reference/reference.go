// Package reference: the tube with a BUTTRESS-thread groove cut around the
// inside of the bore at mid height.
//
// The r-z profile the user described as "an L with a 45 degree little roof,
// like a staple where one leg is bent at 45 degrees":
//
//	r=6 (bore)                                 r=8 (back)
//	  z=+1.5  o                                              <- roof breaks
//	           \                                                out at the
//	            \   45 deg roof, up and IN toward the bore       bore
//	             \
//	  z=-0.5      \                                   o        <- top of the
//	                                                  |           back wall
//	  z=-1.5  o----------------------------------------o       <- flat floor
//	                    flat floor, 2 mm deep          ^ vertical back wall,
//	                                                     1 tall
//
// Opening at the bore is the full 3 tall; the back wall is 3 - 2 = 1 tall
// because the 45 degree roof eats exactly the 2 mm of depth. Printed
// tube-axis-vertical this has a flat floor and no flat ceiling anywhere in
// the groove zone, which is the whole point.
//
// Incident: mcweed I-4 (U333-U341), four rounds and three rebuilds to get
// this profile; resolved only when the model named it a buttress thread.
package reference

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

const (
	outerR = 10.0
	innerR = 6.0
	tubeH  = 20.0

	grooveDepth = 2.0                      // radial, bore out into the wall
	grooveTall  = 3.0                      // opening height at the bore
	grooveBackR = innerR + grooveDepth     // r = 8, the deep back wall
	grooveZ0    = -grooveTall / 2          // flat floor, z = -1.5
	backTall    = grooveTall - grooveDepth // 45 deg roof eats 2 of the 3

	// overcut pulls the profile inward past the bore so the cutter's inner
	// face exits the body rather than coinciding with the bore surface
	// (coincident faces make phantom SDF membranes, mcweed U299).
	overcut = 1.0
)

func Build() *solid.Solid {
	tube := solid.Cylinder(tubeH, outerR, 0).Cut(solid.Cylinder(tubeH+2, innerR, 0))

	// Profile in (r, z). Revolve() spins a 2D shape about the Z axis: shape
	// X becomes radius, shape Y becomes Z. The roof lies on the 45 degree
	// line z = (grooveZ0 + backTall) + (grooveBackR - r), extended inward to
	// r = innerR - overcut so the cut face is continuous through the bore.
	roofInnerZ := grooveZ0 + backTall + (grooveBackR - (innerR - overcut))
	groove := shape.Polygon([]v2.Vec{
		v2.XY(innerR-overcut, grooveZ0),       // floor, inside the bore
		v2.XY(grooveBackR, grooveZ0),          // floor, at the back wall
		v2.XY(grooveBackR, grooveZ0+backTall), // top of the vertical back wall
		v2.XY(innerR-overcut, roofInnerZ),     // 45 deg roof, up and in
	}).Revolve()

	return tube.Cut(groove)
}
