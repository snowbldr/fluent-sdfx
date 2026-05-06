package solid

import (
	"math"

	"github.com/snowbldr/sdfx/sdf"
	v2sdf "github.com/snowbldr/sdfx/vec/v2"
)

// SweepHelix sweeps a 2D profile along a helix at the given radius.
//
// When flatEnds is false the ends are horizontally flat (partial-thread
// cross-sections sliced by z=±height/2). When flatEnds is true each end
// shows a full profile cross-section — the thread is cut perpendicular to
// its own sweep direction, so material may extend slightly past ±height/2
// along the helix tangent.
//
// Panics on non-positive height or turns.
func SweepHelix(profile sdf.SDF2, radius, turns, height float64, flatEnds bool) *Solid {
	if turns <= 0 {
		panic("solid.SweepHelix: turns must be > 0")
	}
	if height <= 0 {
		panic("solid.SweepHelix: height must be > 0")
	}
	pitch := height / turns

	// Screw3D convention: profile X is axial, Y is radial from the helix axis.
	// Compose the transform matrix once: translate(Y=radius) · mirrorX · rotate(-90°).
	m := sdf.Translate2d(v2sdf.Vec{Y: radius}).
		Mul(sdf.MirrorX()).
		Mul(sdf.Rotate2d(-90 * math.Pi / 180))
	screwProfile := sdf.Transform2D(profile, m)

	if !flatEnds {
		return Screw(screwProfile, height, 0, pitch, 1)
	}
	return &Solid{newFlatHelixSDF3(screwProfile, radius, pitch, turns)}
}
