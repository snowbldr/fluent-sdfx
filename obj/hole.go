package obj

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/sdfx/obj"
)

// CircleGrilleParms configures a circular grille pattern.
type CircleGrilleParms = obj.CircleGrilleParms

// KeyedHoleParms configures a keyed hole (circle with N key slots).
type KeyedHoleParms = obj.KeyedHoleParms

// CircleGrille2D returns the 2D profile of a circle grille.
func CircleGrille2D(p CircleGrilleParms) *shape.Shape {
	s, err := obj.CircleGrille2D(&p)
	if err != nil {
		panic(err)
	}
	return shape.Wrap2D(s)
}

// CircleGrille3D returns a 3D extruded circle grille.
func CircleGrille3D(p CircleGrilleParms) *solid.Solid {
	s, err := obj.CircleGrille3D(&p)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}

// CounterBoredHole3D returns a counterbored hole through a given length.
func CounterBoredHole3D(l, r, cbRadius, cbDepth float64) *solid.Solid {
	s, err := obj.CounterBoredHole3D(l, r, cbRadius, cbDepth)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}

// ChamferedHole3D returns a chamfered hole.
func ChamferedHole3D(l, r, chRadius float64) *solid.Solid {
	s, err := obj.ChamferedHole3D(l, r, chRadius)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}

// CounterSunkHole3D returns a countersunk hole.
func CounterSunkHole3D(l, r float64) *solid.Solid {
	s, err := obj.CounterSunkHole3D(l, r)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}

// BoltCircle2D returns a 2D bolt circle pattern.
func BoltCircle2D(holeRadius, circleRadius float64, numHoles int) *shape.Shape {
	s, err := obj.BoltCircle2D(holeRadius, circleRadius, numHoles)
	if err != nil {
		panic(err)
	}
	return shape.Wrap2D(s)
}

// BoltCircle3D returns a 3D bolt circle hole pattern.
func BoltCircle3D(holeDepth, holeRadius, circleRadius float64, numHoles int) *solid.Solid {
	s, err := obj.BoltCircle3D(holeDepth, holeRadius, circleRadius, numHoles)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}

// KeyedHole2D returns a 2D keyed hole.
func KeyedHole2D(p KeyedHoleParms) *shape.Shape {
	s, err := obj.KeyedHole2D(&p)
	if err != nil {
		panic(err)
	}
	return shape.Wrap2D(s)
}

// KeyedHole3D returns a 3D keyed hole.
func KeyedHole3D(p KeyedHoleParms) *solid.Solid {
	s, err := obj.KeyedHole3D(&p)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}
