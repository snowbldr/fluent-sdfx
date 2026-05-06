package shape

import (
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	"github.com/snowbldr/sdfx/sdf"
)

// AnchoredShape is a *Shape paired with a marker point on (or inside) its
// bounding box. Returned by every anchor selector on Shape (Top, BottomRight,
// ...); receives placement verbs (On, At, Above, ...) that return a moved
// *Shape (or a Placement2D when the boolean partner is implicit).
//
// 2D follows screen convention: +Y is "up" / Top.
type AnchoredShape struct {
	Shape *Shape
	Point v2.Vec
}

// Placement2D is the intermediate produced by relative placement verbs in 2D.
// It carries the moved active shape and the implicit boolean partner (the
// target's owner), and is finalized with Union/Cut/Intersect/Shape — the
// chain's subject (the moved shape) is what's kept, matching s.Cut(other).
type Placement2D struct {
	Moved *Shape // the active shape, after positioning
	Base  *Shape // the target's owner — implicit boolean partner
}

// Union unions the moved shape with the base (commutative).
func (p Placement2D) Union() *Shape { return p.Base.Union(p.Moved) }

// Add is an alias for Union.
func (p Placement2D) Add() *Shape { return p.Union() }

// SmoothUnion smoothly unions the moved shape with the base.
func (p Placement2D) SmoothUnion(min sdf.MinFunc) *Shape {
	return SmoothUnion(min, p.Base, p.Moved)
}

// SmoothAdd is an alias for SmoothUnion.
func (p Placement2D) SmoothAdd(min sdf.MinFunc) *Shape { return p.SmoothUnion(min) }

// Cut subtracts the base from the moved shape; the moved (subject) is kept.
func (p Placement2D) Cut() *Shape { return p.Moved.Cut(p.Base) }

// Difference is an alias for Cut.
func (p Placement2D) Difference() *Shape { return p.Cut() }

// SmoothCut smoothly subtracts the base from the moved shape.
func (p Placement2D) SmoothCut(max sdf.MaxFunc) *Shape {
	return SmoothCut(max, p.Moved, p.Base)
}

// SmoothDifference is an alias for SmoothCut.
func (p Placement2D) SmoothDifference(max sdf.MaxFunc) *Shape { return p.SmoothCut(max) }

// Intersect intersects the moved with the base (commutative).
func (p Placement2D) Intersect() *Shape { return p.Base.Intersect(p.Moved) }

// SmoothIntersect smoothly intersects the moved with the base.
func (p Placement2D) SmoothIntersect(max sdf.MaxFunc) *Shape {
	return SmoothIntersect(max, p.Base, p.Moved)
}

// Shape returns the moved shape alone, discarding the boolean partner.
func (p Placement2D) Shape() *Shape { return p.Moved }

// --- Anchor selectors on Shape (the 8 + center) ---
// 2D follows screen convention: +Y is up / Top.

// Top returns the anchor at the middle of the shape's top edge (max Y).
func (s *Shape) Top() AnchoredShape { return s.anchor(0, 1) }

// Bottom returns the anchor at the middle of the shape's bottom edge (min Y).
func (s *Shape) Bottom() AnchoredShape { return s.anchor(0, -1) }

// Right returns the anchor at the middle of the shape's right edge (max X).
func (s *Shape) Right() AnchoredShape { return s.anchor(1, 0) }

// Left returns the anchor at the middle of the shape's left edge (min X).
func (s *Shape) Left() AnchoredShape { return s.anchor(-1, 0) }

// TopRight returns the anchor at the shape's top-right corner.
func (s *Shape) TopRight() AnchoredShape { return s.anchor(1, 1) }

// TopLeft returns the anchor at the shape's top-left corner.
func (s *Shape) TopLeft() AnchoredShape { return s.anchor(-1, 1) }

// BottomRight returns the anchor at the shape's bottom-right corner.
func (s *Shape) BottomRight() AnchoredShape { return s.anchor(1, -1) }

// BottomLeft returns the anchor at the shape's bottom-left corner.
func (s *Shape) BottomLeft() AnchoredShape { return s.anchor(-1, -1) }

// AnchorAt returns the anchor for an arbitrary unit-square coordinate; each
// component is min at -1, center at 0, max at +1.
func (s *Shape) AnchorAt(x, y int) AnchoredShape { return s.anchor(x, y) }

func (s *Shape) anchor(x, y int) AnchoredShape {
	return AnchoredShape{Shape: s, Point: s.Bounds().Anchor(x, y)}
}

// --- Placement verbs on AnchoredShape ---

// On aligns this anchor's point with the target's point and returns a
// Placement2D carrying the moved shape and the target's owner.
func (a AnchoredShape) On(target AnchoredShape) Placement2D {
	moved := a.Shape.Translate(target.Point.Sub(a.Point))
	return Placement2D{Moved: moved, Base: target.Shape}
}

// Above places this anchor at target.Point + (0, gap). Default gap is 0.
func (a AnchoredShape) Above(target AnchoredShape, gap ...float64) Placement2D {
	return a.On(target.shift(v2.Y(firstOr0(gap))))
}

// Below places this anchor at target.Point + (0, -gap).
func (a AnchoredShape) Below(target AnchoredShape, gap ...float64) Placement2D {
	return a.On(target.shift(v2.Y(-firstOr0(gap))))
}

// RightOf places this anchor at target.Point + (gap, 0).
func (a AnchoredShape) RightOf(target AnchoredShape, gap ...float64) Placement2D {
	return a.On(target.shift(v2.X(firstOr0(gap))))
}

// LeftOf places this anchor at target.Point + (-gap, 0).
func (a AnchoredShape) LeftOf(target AnchoredShape, gap ...float64) Placement2D {
	return a.On(target.shift(v2.X(-firstOr0(gap))))
}

// At aligns this anchor with a literal world-space point and returns the moved shape.
func (a AnchoredShape) At(target v2.Vec) *Shape {
	return a.Shape.Translate(target.Sub(a.Point))
}

// AtX moves only along X so this anchor lands at x.
func (a AnchoredShape) AtX(x float64) *Shape { return a.Shape.TranslateX(x - a.Point.X) }

// AtY moves only along Y so this anchor lands at y.
func (a AnchoredShape) AtY(y float64) *Shape { return a.Shape.TranslateY(y - a.Point.Y) }

// ShiftX moves the anchor point d along X without moving the shape;
// useful when chaining a target like "body's right edge, but 2mm out".
func (a AnchoredShape) ShiftX(d float64) AnchoredShape { return a.shift(v2.X(d)) }

// ShiftY moves the anchor point d along Y without moving the shape.
func (a AnchoredShape) ShiftY(d float64) AnchoredShape { return a.shift(v2.Y(d)) }

func (a AnchoredShape) shift(d v2.Vec) AnchoredShape {
	return AnchoredShape{Shape: a.Shape, Point: a.Point.Add(d)}
}

func firstOr0(xs []float64) float64 {
	if len(xs) == 0 {
		return 0
	}
	return xs[0]
}

// --- Shape sugar layer ---

// OnTopOf is sugar for s.Bottom().Above(target, gap...).
func (s *Shape) OnTopOf(target AnchoredShape, gap ...float64) Placement2D {
	return s.Bottom().Above(target, gap...)
}

// UnderneathOf is sugar for s.Top().Below(target, gap...).
func (s *Shape) UnderneathOf(target AnchoredShape, gap ...float64) Placement2D {
	return s.Top().Below(target, gap...)
}

// LeftOf is sugar for s.Right().LeftOf(target, gap...).
func (s *Shape) LeftOf(target AnchoredShape, gap ...float64) Placement2D {
	return s.Right().LeftOf(target, gap...)
}

// RightOf is sugar for s.Left().RightOf(target, gap...).
func (s *Shape) RightOf(target AnchoredShape, gap ...float64) Placement2D {
	return s.Left().RightOf(target, gap...)
}

// Inside places s so its bbox center matches other's bbox center.
func (s *Shape) Inside(other *Shape) Placement2D {
	return s.AnchorAt(0, 0).On(other.AnchorAt(0, 0))
}

// Absolute scalar setters — leave the other axis alone, return *Shape.

// BottomAt translates s so the bottom edge lands at y.
func (s *Shape) BottomAt(y float64) *Shape { return s.Bottom().AtY(y) }

// TopAt translates s so the top edge lands at y.
func (s *Shape) TopAt(y float64) *Shape { return s.Top().AtY(y) }

// LeftAt translates s so the left edge lands at x.
func (s *Shape) LeftAt(x float64) *Shape { return s.Left().AtX(x) }

// RightAt translates s so the right edge lands at x.
func (s *Shape) RightAt(x float64) *Shape { return s.Right().AtX(x) }

// CenterAt translates s so its bounding box center lands at p.
func (s *Shape) CenterAt(p v2.Vec) *Shape { return s.AnchorAt(0, 0).At(p) }
