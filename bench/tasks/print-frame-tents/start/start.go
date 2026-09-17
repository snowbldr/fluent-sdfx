// Package start is the bracket before the change: a wall bracket used with
// Z up, modelled in the use frame.
//
//   - back plate: 30 (X) x 6 (Y) x 40 (Z), y -6..0, z 0..40. The -Y face is
//     the one that sits against the wall in use, and flat on the build plate
//     when it prints.
//   - shelf: 30 (X) x 20 (Y) x 5 (Z) projecting from the plate's +Y face,
//     z 10..15.
//   - two side walls standing on the shelf, 5 thick in X, flush with the
//     plate's sides, y 0..20, z 15..35. Plate + shelf + walls make a cable
//     tray open at the top (+Z) and the front (+Y).
//   - two 4 mm mounting holes straight down (Z) through the shelf.
//
// Nothing has been cut through the side walls yet.
package start

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plateW = 30.0 // X
	plateT = 6.0  // Y, plate spans y -6..0
	plateH = 40.0 // Z, plate spans z 0..40

	shelfD  = 20.0 // Y, out from the plate's +Y face
	shelfT  = 5.0  // Z
	shelfZ0 = 10.0 // shelf spans z 10..15

	wallT = 5.0  // X, thickness of each side wall
	wallD = 20.0 // Y, same footprint as the shelf
	wallH = 20.0 // Z, walls stand on the shelf: z 15..35

	mountD = 4.0  // mounting hole diameter, through the shelf
	mountX = 5.0  // +-5
	mountY = 12.0 // out from the plate face

	over = 2.0 // cutter overshoot past each face
)

// wallX is the X centre of the +X side wall; its outer face is flush with
// the plate's +X side.
const wallX = plateW/2 - wallT/2

func Build() *solid.Solid {
	plate := solid.Box(v3.XYZ(plateW, plateT, plateH), 0).
		Translate(v3.XYZ(0, -plateT/2, plateH/2))

	shelf := solid.Box(v3.XYZ(plateW, shelfD, shelfT), 0).
		Translate(v3.XYZ(0, shelfD/2, shelfZ0+shelfT/2))

	wall := solid.Box(v3.XYZ(wallT, wallD, wallH), 0).
		Translate(v3.XYZ(wallX, wallD/2, shelfZ0+shelfT+wallH/2))

	mount := solid.Cylinder(shelfT+2*over, mountD/2, 0).
		Translate(v3.XYZ(mountX, mountY, shelfZ0+shelfT/2))

	body := plate.Union(shelf, wall, wall.MirrorYZ())
	return body.Cut(mount, mount.MirrorYZ())
}
