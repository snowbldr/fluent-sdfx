// Package asfailed is the incident construction: the groove cutter was made
// generously tall (twice the disc thickness) and CENTRED on the top face, on
// the usual "make cutting tools longer than the body" reflex. Centred on
// z = 4 a cutter 8 mm tall spans z = 0 .. 8, so every groove is punched
// clean through the disc and the "floor" under it does not exist.
//
// Incident: mcweed heating-core tube — "What are these huge holes on the
// bottom of the heating core tube?" The grooves had been punched through the
// 1.25 mm floor. Overshoot is only safe in a direction where there is
// nothing to protect; a cutter that has to leave a floor must be pinned to
// that floor (BottomAt), not centred on the surface it enters.
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	discR = 15.0
	discH = 4.0

	grooveCount  = 6
	grooveW      = 2.0
	grooveDeep   = 1.0
	grooveStartR = 5.0
	grooveOverR  = 2.0

	grooveLen = discR + grooveOverR - grooveStartR
)

func Build() *solid.Solid {
	disc := solid.Cylinder(discH, discR, 0).ZeroZ()

	// WRONG: 8 mm tall and centred on the top face => z = 0 .. 8. The stated
	// 1 mm depth never enters the geometry at all.
	cutter := solid.Box(v3.XYZ(grooveLen, grooveW, discH*2), 0).
		TranslateX(grooveStartR + grooveLen/2).
		TranslateZ(discH)

	return disc.Cut(cutter.RotateCopyZ(grooveCount))
}
