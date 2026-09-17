package solid

import (
	"github.com/snowbldr/sdfx/sdf"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
)

// Field manipulators: wrappers that change what a solid's SDF reports
// without changing the shape through any ordinary boolean or transform.
// They exist because some things a model needs are properties of the
// FIELD, not of the geometry.

// zClampBelowSDF replaces every cross-section below z0 with the z0
// cross-section — a straight prismatic extension downward, floored at
// zBot.
type zClampBelowSDF struct {
	s        sdf.SDF3
	z0, zBot float64
}

func (c *zClampBelowSDF) Evaluate(p v3sdf.Vec) float64 {
	if p.Z >= c.z0 {
		return c.s.Evaluate(p)
	}
	d := c.s.Evaluate(v3sdf.Vec{X: p.X, Y: p.Y, Z: c.z0})
	if floor := c.zBot - p.Z; floor > d {
		d = floor
	}
	return d
}

func (c *zClampBelowSDF) BoundingBox() sdf.Box3 { return c.s.BoundingBox() }

// StraightenBelow returns s with everything below z0 replaced by the
// straight prismatic extension of its z0 cross-section, ending at the
// solid's own bottom.
//
// The use is mating: a twisted or swept feature (a helix, a draft-angled
// boss) that has to enter a straight socket needs a straight lead-in, so
// the two parts push together with no rotation. Clamping the field is
// exact and costs one extra evaluation below z0, where lofting a separate
// straight section onto the bottom would leave a seam in the field.
func (s *Solid) StraightenBelow(z0 float64) *Solid {
	return Wrap(&zClampBelowSDF{s: s.SDF3, z0: z0, zBot: s.Bounds().Min.Z})
}

// boundsOverrideSDF evaluates s but reports a caller-supplied bounding box.
type boundsOverrideSDF struct {
	s  sdf.SDF3
	bb sdf.Box3
}

func (b *boundsOverrideSDF) Evaluate(p v3sdf.Vec) float64 { return b.s.Evaluate(p) }
func (b *boundsOverrideSDF) BoundingBox() sdf.Box3        { return b.bb }

// WithBounds returns s evaluating exactly as before but reporting bb as its
// bounding box. bb must contain s's own bounds.
//
// This is the escape hatch for bbox coupling. A solid's bounds drive the
// render domain and, through helpers that size themselves off Bounds(),
// the dimensions of neighbouring parts — so a boolean that legitimately
// shrinks the geometry can move something downstream that only ever wanted
// the original envelope. Overriding the reported box decouples the two
// without perturbing the field.
//
// Shrinking the reported box below the real surface will silently clip the
// render, so WithBounds enlarges rather than replaces: the result is the
// union of bb and s's own bounds.
func (s *Solid) WithBounds(bb Box3) *Solid {
	own := s.Bounds()
	merged := sdf.Box3{
		Min: v3sdf.Vec{
			X: min(bb.Min.X, own.Min.X),
			Y: min(bb.Min.Y, own.Min.Y),
			Z: min(bb.Min.Z, own.Min.Z),
		},
		Max: v3sdf.Vec{
			X: max(bb.Max.X, own.Max.X),
			Y: max(bb.Max.Y, own.Max.Y),
			Z: max(bb.Max.Z, own.Max.Z),
		},
	}
	return Wrap(&boundsOverrideSDF{s: s.SDF3, bb: merged})
}
