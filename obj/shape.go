package obj

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/sdfx/obj"
)

// IsocelesTrapezoid2D returns a 2D isoceles trapezoid.
func IsocelesTrapezoid2D(base0, base1, height float64) *shape.Shape {
	s, err := obj.IsocelesTrapezoid2D(base0, base1, height)
	if err != nil {
		panic(err)
	}
	return shape.Wrap2D(s)
}

// IsocelesTriangle2D returns a 2D isoceles triangle.
func IsocelesTriangle2D(base, height float64) *shape.Shape {
	s, err := obj.IsocelesTriangle2D(base, height)
	if err != nil {
		panic(err)
	}
	return shape.Wrap2D(s)
}
