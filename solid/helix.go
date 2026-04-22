package solid

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

// SweepHelix sweeps a 2D profile along a helix at the given radius.
//
// When flatEnds is false the ends are horizontally flat (partial-thread
// cross-sections sliced by z=±height/2). When flatEnds is true each end
// shows a full profile cross-section — the thread is cut perpendicular to
// its own sweep direction, so material may extend slightly past ±height/2
// along the helix tangent.
func SweepHelix(profile *shape.Shape, radius, turns, height float64, flatEnds bool) *Solid {
	pitch := height / turns

	// Screw3D convention: profile X is axial, Y is radial from the helix axis.
	screwProfile := profile.
		Rotate(-90).MirrorX().
		Translate(v2.Vec{Y: radius})

	if !flatEnds {
		return Screw(screwProfile, height, 0, pitch, 1)
	}
	return &Solid{newFlatHelixSDF3(screwProfile.SDF2, radius, pitch, turns)}
}
