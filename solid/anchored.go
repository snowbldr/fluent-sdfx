package solid

import (
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// AnchoredSolid is a *Solid paired with a marker point on (or inside) its
// bounding box. Returned by every anchor selector on Solid (Top, BottomRight,
// ...); receives placement verbs (On, At, Above, ...) that return a moved
// *Solid (or a Placement when the boolean partner is implicit).
type AnchoredSolid struct {
	Solid *Solid
	Point v3.Vec // current world-space anchor point on Solid
}

// Placement is the intermediate produced by relative placement verbs. It
// carries the moved active solid and the implicit boolean partner (the
// target's owner), and is finalized with Union/Cut/Intersect/Solid.
type Placement struct {
	Moved *Solid // the active solid, after positioning
	Base  *Solid // the target's owner — implicit boolean partner
}

// Union unions the moved solid into the base.
func (p Placement) Union() *Solid { return p.Base.Union(p.Moved) }

// Add alias for Union
func (p Placement) Add() *Solid { return p.Union() }

// SmoothUnion smoothly unions the moved solid into the base.
func (p Placement) SmoothUnion(min MinFunc) *Solid { return p.Base.SmoothUnion(min, p.Moved) }

// SmoothAdd alias for SmoothUnion
func (p Placement) SmoothAdd(min MinFunc) *Solid { return p.SmoothUnion(min) }

// Difference subtracts the base from the moved (active) solid and returns
// the moved. Mirrors s.Cut(other): the subject of the chain is what's kept.
func (p Placement) Difference() *Solid { return p.Moved.Difference(p.Base) }

// Cut alias for Difference.
func (p Placement) Cut() *Solid { return p.Difference() }

// SmoothDifference smoothly subtracts the base from the moved solid and
// returns the moved.
func (p Placement) SmoothDifference(max MaxFunc) *Solid {
	return p.Moved.SmoothDifference(max, p.Base)
}

// SmoothCut alias for SmoothDifference.
func (p Placement) SmoothCut(max MaxFunc) *Solid { return p.SmoothDifference(max) }

// SmoothIntersect smoothly intersects the base with the moved solid (commutative).
func (p Placement) SmoothIntersect(max MaxFunc) *Solid { return p.Base.SmoothIntersect(max, p.Moved) }

// Intersect intersects the base with the moved solid (commutative).
func (p Placement) Intersect() *Solid { return p.Base.Intersect(p.Moved) }

// Solid returns the moved (active) solid alone, discarding the boolean partner.
func (p Placement) Solid() *Solid { return p.Moved }

// --- Anchor selectors on Solid (the full 27) ---

// 6 face centers.
func (s *Solid) Top() AnchoredSolid    { return s.anchor(0, 0, 1) }
func (s *Solid) Bottom() AnchoredSolid { return s.anchor(0, 0, -1) }
func (s *Solid) Right() AnchoredSolid  { return s.anchor(1, 0, 0) }
func (s *Solid) Left() AnchoredSolid   { return s.anchor(-1, 0, 0) }
func (s *Solid) Back() AnchoredSolid   { return s.anchor(0, 1, 0) }
func (s *Solid) Front() AnchoredSolid  { return s.anchor(0, -1, 0) }

// 12 edge midpoints.
func (s *Solid) TopRight() AnchoredSolid    { return s.anchor(1, 0, 1) }
func (s *Solid) TopLeft() AnchoredSolid     { return s.anchor(-1, 0, 1) }
func (s *Solid) TopFront() AnchoredSolid    { return s.anchor(0, -1, 1) }
func (s *Solid) TopBack() AnchoredSolid     { return s.anchor(0, 1, 1) }
func (s *Solid) BottomRight() AnchoredSolid { return s.anchor(1, 0, -1) }
func (s *Solid) BottomLeft() AnchoredSolid  { return s.anchor(-1, 0, -1) }
func (s *Solid) BottomFront() AnchoredSolid { return s.anchor(0, -1, -1) }
func (s *Solid) BottomBack() AnchoredSolid  { return s.anchor(0, 1, -1) }
func (s *Solid) FrontRight() AnchoredSolid  { return s.anchor(1, -1, 0) }
func (s *Solid) FrontLeft() AnchoredSolid   { return s.anchor(-1, -1, 0) }
func (s *Solid) BackRight() AnchoredSolid   { return s.anchor(1, 1, 0) }
func (s *Solid) BackLeft() AnchoredSolid    { return s.anchor(-1, 1, 0) }

// 8 corners.
func (s *Solid) TopFrontRight() AnchoredSolid    { return s.anchor(1, -1, 1) }
func (s *Solid) TopFrontLeft() AnchoredSolid     { return s.anchor(-1, -1, 1) }
func (s *Solid) TopBackRight() AnchoredSolid     { return s.anchor(1, 1, 1) }
func (s *Solid) TopBackLeft() AnchoredSolid      { return s.anchor(-1, 1, 1) }
func (s *Solid) BottomFrontRight() AnchoredSolid { return s.anchor(1, -1, -1) }
func (s *Solid) BottomFrontLeft() AnchoredSolid  { return s.anchor(-1, -1, -1) }
func (s *Solid) BottomBackRight() AnchoredSolid  { return s.anchor(1, 1, -1) }
func (s *Solid) BottomBackLeft() AnchoredSolid   { return s.anchor(-1, 1, -1) }

// AnchorAt returns the anchor for an arbitrary unit-cube coordinate; each
// component is min at -1, center at 0, max at +1.
func (s *Solid) AnchorAt(x, y, z int) AnchoredSolid { return s.anchor(x, y, z) }

func (s *Solid) anchor(x, y, z int) AnchoredSolid {
	return AnchoredSolid{Solid: s, Point: s.Bounds().Anchor(x, y, z)}
}

// --- Placement verbs on AnchoredSolid ---

// On aligns this anchor's point with the target's point and returns a
// Placement carrying the moved solid and the target's owner.
func (a AnchoredSolid) On(target AnchoredSolid) Placement {
	moved := a.Solid.Translate(target.Point.Sub(a.Point))
	return Placement{Moved: moved, Base: target.Solid}
}

// Above places this anchor at target.Point + (0,0,gap). Default gap is 0.
func (a AnchoredSolid) Above(target AnchoredSolid, gap ...float64) Placement {
	return a.On(target.shift(v3.Z(firstOr0(gap))))
}

// Below places this anchor at target.Point + (0,0,-gap).
func (a AnchoredSolid) Below(target AnchoredSolid, gap ...float64) Placement {
	return a.On(target.shift(v3.Z(-firstOr0(gap))))
}

// RightOf places this anchor at target.Point + (gap,0,0).
func (a AnchoredSolid) RightOf(target AnchoredSolid, gap ...float64) Placement {
	return a.On(target.shift(v3.X(firstOr0(gap))))
}

// LeftOf places this anchor at target.Point + (-gap,0,0).
func (a AnchoredSolid) LeftOf(target AnchoredSolid, gap ...float64) Placement {
	return a.On(target.shift(v3.X(-firstOr0(gap))))
}

// Behind places this anchor at target.Point + (0,gap,0).
func (a AnchoredSolid) Behind(target AnchoredSolid, gap ...float64) Placement {
	return a.On(target.shift(v3.Y(firstOr0(gap))))
}

// InFrontOf places this anchor at target.Point + (0,-gap,0).
func (a AnchoredSolid) InFrontOf(target AnchoredSolid, gap ...float64) Placement {
	return a.On(target.shift(v3.Y(-firstOr0(gap))))
}

// At aligns this anchor with a literal world-space point and returns the moved solid.
func (a AnchoredSolid) At(target v3.Vec) *Solid {
	return a.Solid.Translate(target.Sub(a.Point))
}

// AtX moves only along X so this anchor lands at x.
func (a AnchoredSolid) AtX(x float64) *Solid { return a.Solid.TranslateX(x - a.Point.X) }

// AtY moves only along Y so this anchor lands at y.
func (a AnchoredSolid) AtY(y float64) *Solid { return a.Solid.TranslateY(y - a.Point.Y) }

// AtZ moves only along Z so this anchor lands at z.
func (a AnchoredSolid) AtZ(z float64) *Solid { return a.Solid.TranslateZ(z - a.Point.Z) }

// ShiftX moves the anchor point d along X without moving the solid;
// useful when chaining a target like "body's top, but 2mm up".
func (a AnchoredSolid) ShiftX(d float64) AnchoredSolid { return a.shift(v3.X(d)) }
func (a AnchoredSolid) ShiftY(d float64) AnchoredSolid { return a.shift(v3.Y(d)) }
func (a AnchoredSolid) ShiftZ(d float64) AnchoredSolid { return a.shift(v3.Z(d)) }

func (a AnchoredSolid) shift(d v3.Vec) AnchoredSolid {
	return AnchoredSolid{Solid: a.Solid, Point: a.Point.Add(d)}
}

func firstOr0(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	return xs[0]
}

// --- Solid sugar layer (Style C) ---

// OnTopOf is sugar for s.Bottom().Above(target, gap...).
func (s *Solid) OnTopOf(target AnchoredSolid, gap ...float64) Placement {
	return s.Bottom().Above(target, gap...)
}

// UnderneathOf is sugar for s.Top().Below(target, gap...).
func (s *Solid) UnderneathOf(target AnchoredSolid, gap ...float64) Placement {
	return s.Top().Below(target, gap...)
}

// LeftOf is sugar for s.Right().LeftOf(target, gap...).
func (s *Solid) LeftOf(target AnchoredSolid, gap ...float64) Placement {
	return s.Right().LeftOf(target, gap...)
}

// RightOf is sugar for s.Left().RightOf(target, gap...).
func (s *Solid) RightOf(target AnchoredSolid, gap ...float64) Placement {
	return s.Left().RightOf(target, gap...)
}

// InFrontOf is sugar for s.Back().InFrontOf(target, gap...).
func (s *Solid) InFrontOf(target AnchoredSolid, gap ...float64) Placement {
	return s.Back().InFrontOf(target, gap...)
}

// BehindOf is sugar for s.Front().Behind(target, gap...).
func (s *Solid) BehindOf(target AnchoredSolid, gap ...float64) Placement {
	return s.Front().Behind(target, gap...)
}

// Inside places s so its bbox center matches other's bbox center.
func (s *Solid) Inside(other *Solid) Placement {
	return s.AnchorAt(0, 0, 0).On(other.AnchorAt(0, 0, 0))
}

// Absolute scalar setters — leave other axes alone, return *Solid.
func (s *Solid) BottomAt(z float64) *Solid { return s.Bottom().AtZ(z) }
func (s *Solid) TopAt(z float64) *Solid    { return s.Top().AtZ(z) }
func (s *Solid) LeftAt(x float64) *Solid   { return s.Left().AtX(x) }
func (s *Solid) RightAt(x float64) *Solid  { return s.Right().AtX(x) }
func (s *Solid) FrontAt(y float64) *Solid  { return s.Front().AtY(y) }
func (s *Solid) BackAt(y float64) *Solid   { return s.Back().AtY(y) }

// CenterAt translates s so its bounding box center lands at p.
func (s *Solid) CenterAt(p v3.Vec) *Solid { return s.AnchorAt(0, 0, 0).At(p) }
