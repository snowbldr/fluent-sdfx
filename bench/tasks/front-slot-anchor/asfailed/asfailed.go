// Package asfailed is the incident construction: the model read "front" as
// the +Y side (the face pointing "away" in its head) and pinned the cutter
// with FrontAt at a POSITIVE y, so the slot came out on the back (high-Y)
// face and the front face was left untouched.
//
// Incident: mcweed 2026-07-11, bit twice — the stem wire flats (grooves
// that never reached the surface) and the housing key block/keyway (pocket
// floating over the bay). FrontAt(v) pins the LOW-Y face at v; BackAt(v)
// pins the HIGH-Y face. Intuition inverts this.
package asfailed

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	blockW = 30.0
	blockD = 20.0
	blockH = 10.0

	slotWide = 4.0
	slotDeep = 2.0
	slotOver = 2.0
)

func Build() *solid.Solid {
	block := solid.Box(v3.XYZ(blockW, blockD, blockH), 0)

	// WRONG: FrontAt(+8) pins the cutter's LOW-Y face at y = +8, so the
	// cutter occupies y = +8 .. +12 and the slot lands on the BACK face.
	cutter := solid.Box(v3.XYZ(blockW+2*slotOver, 2*slotDeep, slotWide), 0).
		FrontAt(blockD/2 - slotDeep)

	return block.Cut(cutter)
}
