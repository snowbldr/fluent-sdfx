// Package testparts holds the part functions the history tests register,
// in a package of their own so the Go export of a history that uses them
// can import them.
package testparts

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// RingParams is a registered function's struct of parameters.
type RingParams struct {
	Outer, Inner, Height float64
	Notches              int
}

// Ring is a part function of the kind a person writes in Go and a history
// composes: a washer with notches cut in its rim.
func Ring(p RingParams) *solid.Solid {
	r := solid.Cylinder(p.Height, p.Outer, 0).Cut(solid.Cylinder(p.Height+1, p.Inner, 0))
	if p.Notches > 0 {
		notch := solid.Box(v3.XYZ(2, 2, p.Height+1), 0).TranslateX(p.Outer).RotateCopyZ(p.Notches)
		r = r.Cut(notch)
	}
	return r
}

// Peg is registered with parameter names.
func Peg(height, radius float64) *solid.Solid { return solid.Cylinder(height, radius, 0.3) }
