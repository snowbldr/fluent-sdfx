package solid

import (
	"math"

	"github.com/snowbldr/sdfx/sdf"
	v2sdf "github.com/snowbldr/sdfx/vec/v2"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
)

// flatHelixSDF3 evaluates a helical sweep with flat end faces perpendicular
// to the sweep direction. Each end shows a full profile cross-section rather
// than Screw3D's horizontal z-plane slice through a partial thread.
//
// The body is the continuous sweep of a 2D profile (X=axial, Y=radial) along
// the helical centerline for turn number t ∈ [-halfTurns, halfTurns]. Its
// boundary has two parts: the lateral helicoidal surface and two flat end
// caps, each a copy of the profile embedded in the meridional plane at
// t=±halfTurns.
//
// SDF algorithm: the lateral surface contribution is the min profile distance
// over nearby in-range integer turns k (t=phi+k lies in the valid range where
// phi=θ/2π). The end-cap contribution is the 3D distance to the profile
// cross-section in its meridional plane: perpendicular distance to the plane
// combined with in-plane profile distance. Combining gives interior points
// the closer of lateral/cap distance, and exterior points likewise.
type flatHelixSDF3 struct {
	profile   sdf.SDF2
	pitch     float64
	halfTurns float64
	jScale    float64
	bb        sdf.Box3
}

func newFlatHelixSDF3(profile sdf.SDF2, radius, pitch, turns float64) *flatHelixSDF3 {
	bb2 := profile.BoundingBox()
	axialMin := bb2.Min.X
	axialMax := bb2.Max.X
	rMax := bb2.Max.Y

	halfTurns := turns / 2
	slope := pitch / (2 * math.Pi) / radius
	jScale := 1.0 / math.Sqrt(1+slope*slope)

	halfLen := pitch * halfTurns
	halfAxial := math.Max(math.Abs(axialMin), math.Abs(axialMax))
	zExt := halfLen + halfAxial

	return &flatHelixSDF3{
		profile:   profile,
		pitch:     pitch,
		halfTurns: halfTurns,
		jScale:    jScale,
		bb: sdf.Box3{
			Min: v3sdf.Vec{X: -rMax, Y: -rMax, Z: -zExt},
			Max: v3sdf.Vec{X: rMax, Y: rMax, Z: zExt},
		},
	}
}

func (h *flatHelixSDF3) Evaluate(p v3sdf.Vec) float64 {
	r := math.Sqrt(p.X*p.X + p.Y*p.Y)
	theta := math.Atan2(p.Y, p.X)
	phi := theta / (2 * math.Pi)

	// Valid integer k's: t = phi + k ∈ [-halfTurns, halfTurns].
	kMin := math.Ceil(-h.halfTurns - phi - 1e-12)
	kMax := math.Floor(h.halfTurns - phi + 1e-12)

	// Lateral distance: min 2D profile distance over all in-range k's. The
	// profile may be wider than one pitch and offset from axial=0, so multiple
	// k's can contribute — just iterate the whole range.
	dLat := math.Inf(1)
	for k := kMin; k <= kMax; k++ {
		t := phi + k
		axial := p.Z - h.pitch*t
		pd := h.profile.Evaluate(v2sdf.Vec{X: axial, Y: r})
		if pd < dLat {
			dLat = pd
		}
	}
	dLat *= h.jScale

	// End caps: 3D distance to the profile cross-section at each end.
	dCap := math.Min(h.endCapDist(p, h.halfTurns), h.endCapDist(p, -h.halfTurns))

	if dLat < 0 {
		// Interior: nearest boundary is whichever of lateral or cap is closer.
		if -dCap > dLat {
			return -dCap
		}
		return dLat
	}
	if dCap < dLat {
		return dCap
	}
	return dLat
}

// endCapDist returns the unsigned 3D distance from p to the flat end-cap
// cross-section at turn number tEnd. The cap is the 2D profile region
// embedded in the meridional plane at θ = 2π·tEnd, z-offset pitch·tEnd.
func (h *flatHelixSDF3) endCapDist(p v3sdf.Vec, tEnd float64) float64 {
	thetaEnd := 2 * math.Pi * tEnd
	cosT := math.Cos(thetaEnd)
	sinT := math.Sin(thetaEnd)
	// Perpendicular distance from p to the meridional plane.
	perp := math.Abs(-p.X*sinT + p.Y*cosT)
	// Signed radial coordinate of p's projection within the plane.
	rad := p.X*cosT + p.Y*sinT
	ax := p.Z - h.pitch*tEnd
	profD := h.profile.Evaluate(v2sdf.Vec{X: ax, Y: rad})
	if profD <= 0 {
		return perp
	}
	return math.Sqrt(perp*perp + profD*profD)
}

func (h *flatHelixSDF3) BoundingBox() sdf.Box3 {
	return h.bb
}
