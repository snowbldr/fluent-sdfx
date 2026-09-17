// Package start: the two halves of a housing, laid out side by side on the
// build plate. Plate A is on the left (centre x = -20), plate B on the right
// (centre x = +20). In service B is flipped over and mated to A; for scoring
// they simply sit side by side in the same frame. No registration features
// exist yet.
package start

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	plateX = 24.0
	plateY = 24.0
	plateT = 4.0   // z 0..4
	plateA = -20.0 // centre X of half A
	plateB = 20.0  // centre X of half B
)

func plate(cx float64) *solid.Solid {
	return solid.Box(v3.XYZ(plateX, plateY, plateT), 0).Translate(v3.XYZ(cx, 0, plateT/2))
}

func Build() *solid.Solid {
	return plate(plateA).Union(plate(plateB))
}
