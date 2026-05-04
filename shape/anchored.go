package shape

import (
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
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
type Placement2D struct {
	Moved *Shape
	Base  *Shape
}

func (p Placement2D) Union() *Shape     { return p.Base.Union(p.Moved) }
func (p Placement2D) Cut() *Shape       { return p.Base.Cut(p.Moved) }
func (p Placement2D) Intersect() *Shape { return p.Base.Intersect(p.Moved) }
func (p Placement2D) Shape() *Shape     { return p.Moved }

// --- Anchor selectors on Shape (the 8 + center) ---

// 4 face midpoints. +Y is up.
func (s *Shape) Top() AnchoredShape    { return s.anchor(0, 1) }
func (s *Shape) Bottom() AnchoredShape { return s.anchor(0, -1) }
func (s *Shape) Right() AnchoredShape  { return s.anchor(1, 0) }
func (s *Shape) Left() AnchoredShape   { return s.anchor(-1, 0) }

// 4 corners.
func (s *Shape) TopRight() AnchoredShape    { return s.anchor(1, 1) }
func (s *Shape) TopLeft() AnchoredShape     { return s.anchor(-1, 1) }
func (s *Shape) BottomRight() AnchoredShape { return s.anchor(1, -1) }
func (s *Shape) BottomLeft() AnchoredShape  { return s.anchor(-1, -1) }

// AnchorAt returns the anchor for an arbitrary unit-square coordinate.
func (s *Shape) AnchorAt(x, y int) AnchoredShape { return s.anchor(x, y) }

func (s *Shape) anchor(x, y int) AnchoredShape {
	return AnchoredShape{Shape: s, Point: s.Bounds().Anchor(x, y)}
}

// --- Placement verbs on AnchoredShape ---

func (a AnchoredShape) On(target AnchoredShape) Placement2D {
	moved := a.Shape.Translate(target.Point.Sub(a.Point))
	return Placement2D{Moved: moved, Base: target.Shape}
}

func (a AnchoredShape) Above(target AnchoredShape, gap ...float64) Placement2D {
	return a.On(target.shift(v2.Y(firstOr0(gap))))
}

func (a AnchoredShape) Below(target AnchoredShape, gap ...float64) Placement2D {
	return a.On(target.shift(v2.Y(-firstOr0(gap))))
}

func (a AnchoredShape) RightOf(target AnchoredShape, gap ...float64) Placement2D {
	return a.On(target.shift(v2.X(firstOr0(gap))))
}

func (a AnchoredShape) LeftOf(target AnchoredShape, gap ...float64) Placement2D {
	return a.On(target.shift(v2.X(-firstOr0(gap))))
}

func (a AnchoredShape) At(target v2.Vec) *Shape {
	return a.Shape.Translate(target.Sub(a.Point))
}

func (a AnchoredShape) AtX(x float64) *Shape { return a.Shape.TranslateX(x - a.Point.X) }
func (a AnchoredShape) AtY(y float64) *Shape { return a.Shape.TranslateY(y - a.Point.Y) }

func (a AnchoredShape) ShiftX(d float64) AnchoredShape { return a.shift(v2.X(d)) }
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

func (s *Shape) OnTopOf(target AnchoredShape, gap ...float64) Placement2D {
	return s.Bottom().Above(target, gap...)
}

func (s *Shape) UnderneathOf(target AnchoredShape, gap ...float64) Placement2D {
	return s.Top().Below(target, gap...)
}

func (s *Shape) LeftOf(target AnchoredShape, gap ...float64) Placement2D {
	return s.Right().LeftOf(target, gap...)
}

func (s *Shape) RightOf(target AnchoredShape, gap ...float64) Placement2D {
	return s.Left().RightOf(target, gap...)
}

func (s *Shape) Inside(other *Shape) Placement2D {
	return s.AnchorAt(0, 0).On(other.AnchorAt(0, 0))
}

func (s *Shape) BottomAt(y float64) *Shape { return s.Bottom().AtY(y) }
func (s *Shape) TopAt(y float64) *Shape    { return s.Top().AtY(y) }
func (s *Shape) LeftAt(x float64) *Shape   { return s.Left().AtX(x) }
func (s *Shape) RightAt(x float64) *Shape  { return s.Right().AtX(x) }

// CenterAt translates s so its bounding box center lands at p.
func (s *Shape) CenterAt(p v2.Vec) *Shape { return s.AnchorAt(0, 0).At(p) }
