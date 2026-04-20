package main

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// material shrinkage
var shrink = 1.0 / 0.999 // PLA ~0.1%

func colSpace(radius float64) float64 {
	return (4.0 * radius) / math.Sqrt(3.0)
}

func rowSpace(radius float64) float64 {
	return 2.0 * radius
}

func xOffset(radius float64) float64 {
	return (2.0 * radius) / math.Sqrt(3.0)
}

func yOffset(radius float64) float64 {
	return (2.0 * radius) / 3.0
}

func zOffset(radius float64) float64 {
	return (4.0 * radius) / 3.0
}

func ballRow(ncol int, radius float64) *solid.Solid {
	space := colSpace(radius)
	x := v3.X(-0.5 * ((float64(ncol) - 1) * space))
	dx := v3.X(space)

	sphere := solid.Sphere(radius)
	balls := make([]*solid.Solid, ncol)
	for i := 0; i < ncol; i++ {
		balls[i] = sphere.Translate(x)
		x = x.Add(dx)
	}
	return solid.UnionAll(balls...)
}

func ballGrid(ncol, nrow int, radius float64) *solid.Solid {
	space := rowSpace(radius)
	x := v3.Y(-0.5 * ((float64(nrow) - 1) * space))
	dy0 := v3.XY(-xOffset(radius), space)
	dy1 := v3.XY(xOffset(radius), space)

	row := ballRow(ncol, radius)
	rows := make([]*solid.Solid, nrow)
	for i := 0; i < nrow; i++ {
		rows[i] = row.Translate(x)
		if i%2 == 0 {
			x = x.Add(dy0)
		} else {
			x = x.Add(dy1)
		}
	}
	return solid.UnionAll(rows...)
}

func macCheeseGrater(ncol, nrow int, radius float64) *solid.Solid {
	dx := v3.XYZ(xOffset(radius), yOffset(radius), zOffset(radius)).MulScalar(0.5)
	g := ballGrid(ncol, nrow, radius)
	g0 := g.Translate(dx.Neg())
	g1 := g.Translate(dx)
	balls := g0.Union(g1)

	pX := colSpace(radius) * (float64(ncol) - 1)
	pY := rowSpace(radius) * (float64(nrow) - 1)
	pZ := 0.5 * colSpace(radius)
	plate := solid.Box(v3.XYZ(pX, pY, pZ), 0)
	return plate.Cut(balls)
}

func main() {
	macCheeseGrater(15, 6, 10.0).ScaleUniform(shrink).STL("mcg.stl", 5.0)
}
