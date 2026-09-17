// Package reference: the pocket with a tented (gable) ceiling. The pocket
// is 10 wide in X, 12 deep into the -Y face and its floor is at z=4. Its
// ceiling is where the printability lives: instead of a flat roof at z=12
// (an 10 x 12 mm unsupported overhang) the ceiling rises at 45 degrees from
// both side walls, from z=12 at x=+-5 to a ridge 5 mm higher at x=0.
// Every ceiling face is then at 45 degrees, so the pocket prints
// support-free on the block's -Z face.
//
// Incident: mcweed U189-U192 and mined4 I-3 -- "tented" pockets, a 45
// degree tent/frustum ceiling for a supportless interior; the fix always
// came from the user's slicer overhang view, never from the model.
package reference

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	blockW = 24.0
	blockD = 24.0
	blockH = 20.0

	pocketW   = 10.0        // X
	pocketH   = 8.0         // Z of the straight walls: floor z=4 to eaves z=12
	pocketDep = 12.0        // Y, in from the -Y face
	pocketZ0  = 4.0         // floor
	ridgeRise = pocketW / 2 // 45 degrees from the walls: 5 mm up to the ridge
	over      = 2.0         // cutter overshoot out through the -Y face
)

func Build() *solid.Solid {
	block := solid.Box(v3.XYZ(blockW, blockD, blockH), 0).TranslateZ(blockH / 2)

	// House-shaped profile in the XZ plane: floor, two side walls to the
	// eaves at z=12, then two 45 degree slopes meeting at the ridge.
	eaves := pocketZ0 + pocketH
	profile := shape.Polygon([]v2.Vec{
		v2.XY(-pocketW/2, pocketZ0),
		v2.XY(pocketW/2, pocketZ0),
		v2.XY(pocketW/2, eaves),
		v2.XY(0, eaves+ridgeRise),
		v2.XY(-pocketW/2, eaves),
	})

	// Extrude along Y, then set it to span y = -(12+over) .. 0 so the
	// pocket is 12 deep from the -Y face at y = -12.
	depth := pocketDep + over
	pocket := profile.Extrude(depth).RotateX(90).TranslateY(-depth / 2)

	return block.Cut(pocket)
}
