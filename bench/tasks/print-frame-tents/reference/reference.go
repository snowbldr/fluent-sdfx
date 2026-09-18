// Package reference: the two cable pass-throughs, tented for the PRINT
// frame.
//
// The bracket is modelled and used with Z up, but it prints lying on its
// back: the plate's -Y face goes on the build plate. So in the print frame
// model +Y is up and model +-Z is horizontal. The pass-throughs run along X
// through the 5 mm side walls, which stay vertical walls in the print frame,
// so each hole is a horizontal tunnel with a real ceiling.
//
// That ceiling is the face of the hole at high model Y. The gable therefore
// lives in the YZ cross-section: 45 degree slopes from the two upper corners
// of the 8 mm square opening, at (z = 21, y = 13) and (z = 29, y = 13), up
// to a ridge at (z = 25, y = 17) which runs along X, the hole's own axis.
// Every ceiling face is then at 45 degrees in the print frame and the hole
// prints support-free. The 8 mm square clear opening is untouched; the tent
// sits on top of it.
//
// Incident: mined1 #7 (U18) -- "you have some funny bits on the housing,
// angled roofs where it's not appropriate. the two halves of the body will
// print lying down"; the model carried a standing-print assumption through
// five revisions and gabled everything for the wrong frame. mined3: the
// model/print/use frame mapping was "THE critical lesson of this session".
package reference

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

	holeSq = 8.0        // clear square opening: 8 in Y by 8 in Z
	holeY  = 9.0        // centre, model Y
	holeZ  = 25.0       // centre, model Z
	rise   = holeSq / 2 // 45 degrees across 8 mm of opening: 4 mm to the ridge
)

const wallX = plateW/2 - wallT/2

// tent is the pass-through cutter: a house-shaped prism along X. The 2D
// profile is built with its +Y axis along model +Y (up in the print frame)
// and extruded along Z, then RotateY(90) lays the extrusion axis along X
// and leaves the apex pointing at model +Y.
func tent() *solid.Solid {
	h := holeSq / 2
	profile := shape.Polygon([]v2.Vec{
		v2.XY(-h, -h),
		v2.XY(h, -h),
		v2.XY(h, h),
		v2.XY(0, h+rise), // ridge, over the middle of the opening
		v2.XY(-h, h),
	})
	return profile.Extrude(wallT + 2*over).
		RotateY(90).
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
