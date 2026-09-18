// Package asfailed: the model made the requested 0.4 mm change and then
// "improved" the part on its own -- a 1 mm 45 degree taper at the top of
// the boss and a fillet where the boss meets the plate. Neither was
// asked for; the prompt explicitly said "just that, nothing else".
//
// Incident: mcweed U222-U223, "no taper ffs".
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plateW = 30.0
	plateD = 30.0
	plateH = 4.0

	bossR   = 4.0
	bossH   = 5.6
	taperH  = 1.0 // UNREQUESTED: 45 degree lead-in at the boss top
	filletR = 1.2 // UNREQUESTED: root fillet (RoundMin blend, ~0.8 mm)
)

func Build() *solid.Solid {
	plate := solid.Box(v3.XYZ(plateW, plateD, plateH), 0).TranslateZ(plateH / 2)

	shank := solid.Cylinder(bossH-taperH, bossR, 0).TranslateZ(plateH + (bossH-taperH)/2)
	taper := solid.Cone(taperH, bossR, bossR-taperH, 0).TranslateZ(plateH + bossH - taperH/2)

	return plate.SmoothUnion(solid.RoundMin(filletR), shank.Union(taper))
}
