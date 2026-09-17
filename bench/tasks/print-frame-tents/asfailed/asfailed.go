// Package asfailed: the model tented the pass-throughs for the USE frame.
// It read "45 degree gable so it prints without support", took up to mean
// model +Z because that is up as the bracket hangs on the wall, and built
// the ridge at (y = 9, z = 33) with the slopes coming off the corners at
// z = 29. In the use frame that looks exactly right, and every casual check
// passes: the opening is still 8 mm square, the roof is at 45 degrees, no
// face is horizontal when you look at the render with Z up.
//
// It is useless. The part prints on the plate's -Y face, so print-up is
// model +Y: the real ceiling of each hole is the face at y = 13, which this
// version leaves flat -- an 8 x 5 mm horizontal overhang inside a hole you
// cannot reach -- while the gable it did cut is a sideways notch in the
// wall that spans nothing.
//
// Incident: mined1 #7 (U18) -- "you have some funny bits on the housing,
// angled roofs where it's not appropriate. the two halves of the body will
// print lying down". The model: "an assumption I carried through five
// revisions without asking".
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plateW = 30.0
	plateT = 6.0
	plateH = 40.0

	shelfD  = 20.0
	shelfT  = 5.0
	shelfZ0 = 10.0

	wallT = 5.0
	wallD = 20.0
	wallH = 20.0

	mountD = 4.0
	mountX = 5.0
	mountY = 12.0

	over = 2.0

	holeSq = 8.0
	holeY  = 9.0
	holeZ  = 25.0
	rise   = holeSq / 2
)

const wallX = plateW/2 - wallT/2

// tent is the same house-shaped prism, but the extra RotateX(90) swings the
// apex from model +Y round to model +Z: a gable for the way the bracket is
// used, not for the way it is printed.
func tent() *solid.Solid {
	h := holeSq / 2
	profile := shape.Polygon([]v2.Vec{
		v2.XY(-h, -h),
		v2.XY(h, -h),
		v2.XY(h, h),
		v2.XY(0, h+rise),
		v2.XY(-h, h),
	})
	return profile.Extrude(wallT + 2*over).
		RotateY(90).
		RotateX(90).
		Translate(v3.XYZ(wallX, holeY, holeZ))
}

func Build() *solid.Solid {
	plate := solid.Box(v3.XYZ(plateW, plateT, plateH), 0).
		Translate(v3.XYZ(0, -plateT/2, plateH/2))

	shelf := solid.Box(v3.XYZ(plateW, shelfD, shelfT), 0).
		Translate(v3.XYZ(0, shelfD/2, shelfZ0+shelfT/2))

	wall := solid.Box(v3.XYZ(wallT, wallD, wallH), 0).
		Translate(v3.XYZ(wallX, wallD/2, shelfZ0+shelfT+wallH/2))

	mount := solid.Cylinder(shelfT+2*over, mountD/2, 0).
		Translate(v3.XYZ(mountX, mountY, shelfZ0+shelfT/2))

	t := tent()

	body := plate.Union(shelf, wall, wall.MirrorYZ())
	return body.Cut(mount, mount.MirrorYZ(), t, t.MirrorYZ())
}
