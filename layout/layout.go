// Package layout returns position arrays meant to be spread into the
// variadic .Multi(positions...) on Solid and Shape. The functions are
// pure — they do no SDF evaluation and produce no geometry.
//
//	peg.Multi(layout.Polar(20, 6)...)
//	hole.Multi(layout.Grid(10, 10, 4, 4)...)
//	standoff.Multi(layout.RectCorners(panelW-16, panelD-16)...)
package layout

import (
	"fmt"
	"math"
	"sync/atomic"

	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// gridLimit caps Grid / Grid2 nx*ny to prevent accidental OOMs (e.g. a
// `Grid(1, 1, 100000, 100000)` typo allocates ~240 GB). The default of 1M
// is comfortably above any sensible CAD grid; raise via SetGridLimit if
// you genuinely need denser arrays.
var gridLimitAtomic int64 = 1_000_000

var gridLimit = 1_000_000 // exposed read path; the atomic copy is the source of truth

// SetGridLimit overrides the maximum nx*ny accepted by Grid / Grid2.
// Atomic and safe from any goroutine.
func SetGridLimit(n int) {
	atomic.StoreInt64(&gridLimitAtomic, int64(n))
	gridLimit = n
}

// Polar returns n positions evenly spaced on a circle of the given radius
// in the XY plane (z = 0). The first position is at angle 0 (+X axis).
// Panics if n is negative; n == 0 returns an empty slice.
func Polar(radius float64, n int) []v3.Vec {
	if n < 0 {
		panic(fmt.Errorf("layout.Polar: n must be >= 0, got %d", n))
	}
	out := make([]v3.Vec, n)
	for i := 0; i < n; i++ {
		theta := 2 * math.Pi * float64(i) / float64(n)
		out[i] = v3.XY(radius*math.Cos(theta), radius*math.Sin(theta))
	}
	return out
}

// PolarArc returns n positions evenly spaced along an arc of sweepDeg
// degrees beginning at startDeg, on a circle of the given radius in the
// XY plane (z = 0). For n == 1 the single point lands at startDeg.
func PolarArc(radius float64, n int, startDeg, sweepDeg float64) []v3.Vec {
	out := make([]v3.Vec, n)
	if n == 1 {
		theta := startDeg * math.Pi / 180
		out[0] = v3.XY(radius*math.Cos(theta), radius*math.Sin(theta))
		return out
	}
	for i := 0; i < n; i++ {
		deg := startDeg + sweepDeg*float64(i)/float64(n-1)
		theta := deg * math.Pi / 180
		out[i] = v3.XY(radius*math.Cos(theta), radius*math.Sin(theta))
	}
	return out
}

// Grid returns nx*ny positions on an XY grid centered on the origin
// (z = 0), with stepX between columns and stepY between rows.
// Panics if nx or ny is negative or if nx*ny exceeds gridLimit (a sanity
// cap that prevents accidental OOM from a parameter mistake).
func Grid(stepX, stepY float64, nx, ny int) []v3.Vec {
	if nx < 0 || ny < 0 {
		panic(fmt.Errorf("layout.Grid: nx and ny must be >= 0, got nx=%d ny=%d", nx, ny))
	}
	if nx*ny > gridLimit {
		panic(fmt.Errorf("layout.Grid: nx*ny = %d exceeds limit %d (likely a parameter mistake; raise layout.SetGridLimit if intentional)", nx*ny, gridLimit))
	}
	out := make([]v3.Vec, 0, nx*ny)
	x0 := -float64(nx-1) * stepX / 2
	y0 := -float64(ny-1) * stepY / 2
	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			out = append(out, v3.XY(x0+float64(i)*stepX, y0+float64(j)*stepY))
		}
	}
	return out
}

// Line returns n equally spaced positions from p0 to p1 inclusive.
// For n == 1 the single point lands at p0.
func Line(p0, p1 v3.Vec, n int) []v3.Vec {
	out := make([]v3.Vec, n)
	if n == 1 {
		out[0] = p0
		return out
	}
	d := p1.Sub(p0)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		out[i] = v3.XYZ(p0.X+d.X*t, p0.Y+d.Y*t, p0.Z+d.Z*t)
	}
	return out
}

// RectCorners returns the 4 XY corners of a rectangle of the given width
// (X) and depth (Y), centered on the origin (z = 0).
func RectCorners(width, depth float64) []v3.Vec {
	hx, hy := width/2, depth/2
	return []v3.Vec{
		v3.XY(-hx, -hy),
		v3.XY(hx, -hy),
		v3.XY(hx, hy),
		v3.XY(-hx, hy),
	}
}

// BoxCorners returns the 8 corners of a box of the given size, centered on the origin.
func BoxCorners(size v3.Vec) []v3.Vec {
	hx, hy, hz := size.X/2, size.Y/2, size.Z/2
	return []v3.Vec{
		v3.XYZ(-hx, -hy, -hz),
		v3.XYZ(hx, -hy, -hz),
		v3.XYZ(hx, hy, -hz),
		v3.XYZ(-hx, hy, -hz),
		v3.XYZ(-hx, -hy, hz),
		v3.XYZ(hx, -hy, hz),
		v3.XYZ(hx, hy, hz),
		v3.XYZ(-hx, hy, hz),
	}
}

// --- 2D variants ---

// Polar2 returns n positions evenly spaced on a circle of the given radius.
// Panics if n is negative; n == 0 returns an empty slice.
func Polar2(radius float64, n int) []v2.Vec {
	if n < 0 {
		panic(fmt.Errorf("layout.Polar2: n must be >= 0, got %d", n))
	}
	out := make([]v2.Vec, n)
	for i := 0; i < n; i++ {
		theta := 2 * math.Pi * float64(i) / float64(n)
		out[i] = v2.XY(radius*math.Cos(theta), radius*math.Sin(theta))
	}
	return out
}

// Grid2 returns nx*ny positions on an XY grid centered on the origin.
// Panics on negative nx/ny or if nx*ny exceeds the grid limit.
func Grid2(stepX, stepY float64, nx, ny int) []v2.Vec {
	if nx < 0 || ny < 0 {
		panic(fmt.Errorf("layout.Grid2: nx and ny must be >= 0, got nx=%d ny=%d", nx, ny))
	}
	if nx*ny > gridLimit {
		panic(fmt.Errorf("layout.Grid2: nx*ny = %d exceeds limit %d (raise via layout.SetGridLimit if intentional)", nx*ny, gridLimit))
	}
	out := make([]v2.Vec, 0, nx*ny)
	x0 := -float64(nx-1) * stepX / 2
	y0 := -float64(ny-1) * stepY / 2
	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			out = append(out, v2.XY(x0+float64(i)*stepX, y0+float64(j)*stepY))
		}
	}
	return out
}

// PolarArc2 returns n positions evenly spaced along an arc of sweepDeg degrees
// beginning at startDeg, on a circle of the given radius. For n == 1 the
// single point lands at startDeg.
func PolarArc2(radius float64, n int, startDeg, sweepDeg float64) []v2.Vec {
	out := make([]v2.Vec, n)
	if n == 1 {
		theta := startDeg * math.Pi / 180
		out[0] = v2.XY(radius*math.Cos(theta), radius*math.Sin(theta))
		return out
	}
	for i := 0; i < n; i++ {
		deg := startDeg + sweepDeg*float64(i)/float64(n-1)
		theta := deg * math.Pi / 180
		out[i] = v2.XY(radius*math.Cos(theta), radius*math.Sin(theta))
	}
	return out
}

// Line2 returns n equally spaced positions from p0 to p1 inclusive.
// For n == 1 the single point lands at p0.
func Line2(p0, p1 v2.Vec, n int) []v2.Vec {
	out := make([]v2.Vec, n)
	if n == 1 {
		out[0] = p0
		return out
	}
	d := p1.Sub(p0)
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n-1)
		out[i] = v2.XY(p0.X+d.X*t, p0.Y+d.Y*t)
	}
	return out
}

// RectCorners2 returns the 4 XY corners of a rectangle of the given width
// (X) and depth (Y), centered on the origin.
func RectCorners2(width, depth float64) []v2.Vec {
	hx, hy := width/2, depth/2
	return []v2.Vec{
		v2.XY(-hx, -hy),
		v2.XY(hx, -hy),
		v2.XY(hx, hy),
		v2.XY(-hx, hy),
	}
}
