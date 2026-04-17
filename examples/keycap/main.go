package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

var shrink = 1.0 / 0.999 // PLA ~0.1%

const stemX = 6.0
const stemY = 5.0

const crossDepth = 4.0
const crossWidth = 1.0
const crossX = 4.0
const stemRound = 0.05

// keyStem returns a keycap stem of a given length.
func keyStem(length float64) *solid.Solid {
	ofs := length - crossDepth
	s0 := solid.Box(v3.XYZ(crossX, crossWidth, length), crossX*stemRound)
	s1 := solid.Box(v3.XYZ(crossWidth, stemY*(1.0+2.0*stemRound), length), crossX*stemRound)
	cavity := s0.Union(s1).Translate(v3.XYZ(0, 0, ofs))
	stem := solid.Box(v3.XYZ(stemX, stemY, length), stemX*stemRound)
	return stem.Cut(cavity)
}

const stemLength = 15.0

// roundCap returns a round keycap.
func roundCap(diameter, height, wall float64) *solid.Solid {
	rOuter := 0.5 * diameter
	rInner := 0.5 * (diameter - (2.0 * wall))

	outer := solid.Cylinder(height, rOuter, 0)
	inner := solid.Cylinder(height, rInner, 0).Translate(v3.XYZ(0, 0, wall))
	keycap := outer.Cut(inner)

	stem := keyStem(stemLength)
	ofs := (stemLength - height) * 0.5
	stem = stem.Translate(v3.XYZ(0, 0, ofs))

	return keycap.Union(stem)
}

func main() {
	roundCap(18, 6, 1.5).ScaleUniform(shrink).ToSTL("round_cap.stl", 150)
}
