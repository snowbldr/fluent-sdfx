// Package asfailed: the model added registration keys but made both of them
// the SAME round post and both pockets the same round hole. The key pair is
// then symmetric under a 180 degree rotation about Z: the halves mate in
// either orientation, which is exactly the failure the user was trying to
// design out, and nothing about the part tells you which half is which.
//
// Incident: mcweed U482, "on wrong end + asymmetric" - the model delivered
// keys that registered the halves but did not POLARISE them.
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plateT = 4.0
	plateA = -20.0
	plateB = 20.0

	keyOff  = 8.0
	nubD    = 3.0
	nubH    = 2.0
	pockD   = 3.3
	pockDep = 2.2
	pad     = 0.5
)

func plate(cx float64) *solid.Solid {
	return solid.Box(v3.XYZ(24, 24, plateT), 0).Translate(v3.XYZ(cx, 0, plateT/2))
}

func Build() *solid.Solid {
	// WRONG: two identical round nubs, two identical round pockets.
	nub := solid.Cylinder(nubH, nubD/2, 0).TranslateZ(plateT + nubH/2)
	nubs := nub.Translate(v3.X(plateA + keyOff)).
		Union(nub.Translate(v3.X(plateA - keyOff)))

	cutH := pockDep + pad
	pock := solid.Cylinder(cutH, pockD/2, 0).TranslateZ(plateT - pockDep + cutH/2)
	pocks := pock.Translate(v3.X(plateB - keyOff)).
		Union(pock.Translate(v3.X(plateB + keyOff)))

	return plate(plateA).Union(plate(plateB), nubs).Cut(pocks)
}
