// Package asfailed is the "just a flat groove" reading the model kept
// producing: a plain SQUARE (rectangular) groove, 2 deep and 3 tall, with no
// roof slope at all.
//
// It gets the depth and the height right and the shape wrong. Printed
// tube-axis-vertical the square groove has a flat ceiling spanning the full
// 2 mm of depth — an unsupported bridge, exactly the "pointy bit" / flat
// ceiling the user was trying to get rid of (mcweed I-4, U333-U338: "No, I
// don't think you're getting it. you just made a different bad version of
// what we had").
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

const (
	outerR = 10.0
	innerR = 6.0
	tubeH  = 20.0

	grooveDepth = 2.0
	grooveTall  = 3.0
	overcut     = 1.0
)

func Build() *solid.Solid {
	tube := solid.Cylinder(tubeH, outerR, 0).Cut(solid.Cylinder(tubeH+2, innerR, 0))

	// Square profile in (r, z): flat floor, vertical back wall the full 3
	// tall, flat roof.
	groove := shape.Polygon([]v2.Vec{
		v2.XY(innerR-overcut, -grooveTall/2),
		v2.XY(innerR+grooveDepth, -grooveTall/2),
		v2.XY(innerR+grooveDepth, grooveTall/2),
		v2.XY(innerR-overcut, grooveTall/2),
	}).Revolve()

	return tube.Cut(groove)
}
