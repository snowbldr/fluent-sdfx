// Package start is the part before the change: a chamfered plate carrying two
// bosses, each with a through hole, and a pocket in its top face. Five
// features, all of them already validated on a printed part -- only one of
// them is being changed.
package start

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plateW  = 44.0 // X
	plateD  = 30.0 // Y
	plateH  = 5.0  // Z, the plate spans z 0..5
	chamfer = 2.0  // 45 degree chamfer on all four top edges of the plate

	bossR = 5.0  // boss radius
	bossH = 8.0  // boss height, z 5..13
	bossX = 10.0 // BOTH bosses: one at -bossX, one at +bossX
	holeR = 1.5  // 3 mm through hole on each boss axis

	pocketW = 6.0  // X
	pocketD = 20.0 // Y, the pocket's long axis
	pocketZ = 2.0  // depth below the top face, so its floor is at z 3
)

// chamferedPlate is the plate with a 45 degree chamfer round all four edges
// of its top face: a straight box up to plateH-chamfer, then a band that
// insets by chamfer on every side over the last chamfer of height.
func chamferedPlate() *solid.Solid {
	base := solid.Box(v3.XYZ(plateW, plateD, plateH-chamfer), 0).
		TranslateZ((plateH - chamfer) / 2)
	band := shape.Rect(v2.XY(plateW, plateD), 0).
		ScaleExtrude(chamfer, v2.XY((plateW-2*chamfer)/plateW, (plateD-2*chamfer)/plateD)).
		TranslateZ(plateH - chamfer/2)
	return base.Union(band)
}

func Build() *solid.Solid {
	plate := chamferedPlate()

	boss := solid.Cylinder(bossH, bossR, 0).TranslateZ(plateH + bossH/2)
	hole := solid.Cylinder(plateH+bossH+2, holeR, 0).TranslateZ((plateH + bossH) / 2)

	// Pocket in the top face, long axis along Y, centred on the origin. The
	// tool is twice the depth so it breaks cleanly through the top face.
	pocket := solid.Box(v3.XYZ(pocketW, pocketD, pocketZ*2), 0).TranslateZ(plateH)

	return plate.
		Union(boss.TranslateX(-bossX), boss.TranslateX(bossX)).
		Cut(pocket, hole.TranslateX(-bossX), hole.TranslateX(bossX))
}
