package solid

import (
	"math"

	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
	"github.com/snowbldr/fluent-sdfx/shape"
)

// SweepHelix sweeps a 2D profile along a helix at the given radius.
func SweepHelix(profile *shape.Shape, radius, turns, height float64, flatEnds bool) *Solid {
	pitch := height / turns

	// Prepare profile for Screw3D: mirror+rotate to orient the cross-section,
	// then translate to the helix radius.
	screwProfile := profile.
		Rotate(-90).MirrorX().
		Translate(v2.Vec{Y: radius})

	if !flatEnds {
		return Screw(screwProfile, height, 0, pitch, 1)
	}

	// Extend helix so tapers fall outside desired height range.
	bb := profile.BoundingBox()
	taperHeight := (bb.Max.Y + math.Abs(bb.Min.Y)) / 2
	extHeight := height + 2*taperHeight
	extTurns := extHeight / pitch

	helix := Screw(screwProfile, extHeight, 0, pitch, 1)

	// Taper geometry
	tailExt := bb.Max.Y / pitch * 2 * math.Pi
	circleExt := math.Abs(bb.Min.Y) / pitch * 2 * math.Pi
	wedgeR := radius * 1.5
	big := extHeight * 4
	helixR := radius + (bb.Max.X - bb.Min.X)
	pitchAngle := math.Atan2(pitch, 2*math.Pi*helixR)

	centerAngle := func(halfTurns float64) float64 {
		frac := halfTurns - math.Floor(halfTurns)
		return frac * 2 * math.Pi
	}

	topCenter := centerAngle(extTurns / 2)
	bottomCenter := centerAngle(-extTurns / 2)

	// Slightly grown helix for taper isolation — ensures the cut fully clears
	grownHelix := helix.Grow(0.1)

	isolateTaper := func(center, zOffset float64) *Solid {
		// Wedge spanning the taper's angular range
		start := center - tailExt
		span := tailExt + circleExt
		pts := make([]v2.Vec, 0, 66)
		pts = append(pts, v2.Vec{X: 0, Y: 0})
		for i := 0; i <= 64; i++ {
			a := start + span*float64(i)/64
			pts = append(pts, v2.Vec{X: wedgeR * math.Cos(a), Y: wedgeR * math.Sin(a)})
		}
		wedge2d := shape.Polygon(pts)
		cutterHeight := extHeight * 2
		angularCutter := Cylinder(cutterHeight, wedgeR, 0).
			Cut(Extrude(wedge2d, cutterHeight))

		// Pitch-tilted box to isolate taper from adjacent coils
		pitchAngleDeg := pitchAngle * 180 / math.Pi
		axis := v3.Vec{X: math.Cos(center), Y: math.Sin(center)}
		zCutter := Box(v3.Vec{X: big, Y: big, Z: big}, 0).
			RotateAxis(axis, pitchAngleDeg).
			Translate(v3.Vec{Z: zOffset})

		return grownHelix.Cut(angularCutter, zCutter)
	}

	topTaper := isolateTaper(topCenter, extHeight/2-big/2-pitch*0.6)
	bottomTaper := isolateTaper(bottomCenter, -extHeight/2+big/2+pitch*0.4)

	return helix.Cut(topTaper, bottomTaper)
}
