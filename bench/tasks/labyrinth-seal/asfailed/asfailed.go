// Package asfailed: the model built the groove at exactly the tongue's
// nominal size and forgot the clearance the user asked for. 2 wide and 3
// deep instead of 2.3 wide and 3.15 deep.
//
// Geometrically it looks right in a render and a section: tongue and groove
// are a perfect match. Physically the halves cannot close — a zero-clearance
// tongue is an interference fit on a printed part, so the seam it is meant
// to seal is held open by the seal itself.
//
// Incident: mcweed I-10, where the working seal was specified as "groove
// +0.15 clearance".
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plateLen = 30.0
	plateW   = 20.0
	plateT   = 5.0
	layoutX  = 40.0

	tongueW = 2.0
	tongueH = 3.0

	// WRONG: no clearance applied to the groove.
	grooveW = tongueW
	grooveD = tongueH

	bury    = 0.5
	overrun = 1.0
)

func Build() *solid.Solid {
	plateA := solid.Box(v3.XYZ(plateLen, plateW, plateT), 0).Translate(v3.Z(-plateT / 2))
	plateB := solid.Box(v3.XYZ(plateLen, plateW, plateT), 0).Translate(v3.XZ(layoutX, plateT/2))

	tongue := solid.Box(v3.XYZ(plateLen, tongueW, tongueH+bury), 0).
		Translate(v3.Z((tongueH - bury) / 2))

	groove := solid.Box(v3.XYZ(plateLen+2*overrun, grooveW, grooveD+overrun), 0).
		Translate(v3.XZ(layoutX, (grooveD-overrun)/2))

	return plateA.Union(plateB, tongue).Cut(groove)
}
