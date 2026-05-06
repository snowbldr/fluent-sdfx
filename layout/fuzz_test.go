package layout

import (
	"math"
	"testing"

	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// nFuzzCap bounds how many positions a fuzzed call may request, to keep
// memory bounded regardless of what the fuzzer chooses.
const nFuzzCap = 4096

// finite reports whether x is neither NaN nor infinite.
func finite(x float64) bool {
	return !math.IsNaN(x) && !math.IsInf(x, 0)
}

// FuzzPolar checks that Polar never panics for sane inputs and always
// returns exactly n positions in the XY plane.
func FuzzPolar(f *testing.F) {
	f.Add(10.0, 6)
	f.Add(0.0, 1)
	f.Add(1.0, 0)
	f.Add(1e9, 100)
	f.Add(-5.0, 8)
	f.Fuzz(func(t *testing.T, radius float64, n int) {
		if n < 0 || n > nFuzzCap {
			t.Skip()
		}
		if !finite(radius) {
			t.Skip()
		}
		defer func() {
			if rec := recover(); rec != nil {
				t.Fatalf("layout.Polar(%v, %v) panicked: %v", radius, n, rec)
			}
		}()
		out := Polar(radius, n)
		if len(out) != n {
			t.Fatalf("layout.Polar(%v, %v) len = %d, want %d", radius, n, len(out), n)
		}
		for i, p := range out {
			if p.Z != 0 {
				t.Fatalf("layout.Polar[%d].Z = %v, want 0", i, p.Z)
			}
		}
	})
}

// FuzzPolarArc checks that PolarArc returns exactly n positions in the XY
// plane and never panics on finite inputs (including n == 1, which the
// implementation special-cases).
func FuzzPolarArc(f *testing.F) {
	f.Add(10.0, 6, 0.0, 90.0)
	f.Add(10.0, 1, 45.0, 0.0)
	f.Add(0.0, 4, 0.0, 360.0)
	f.Add(1e6, 100, -720.0, 1440.0)
	f.Fuzz(func(t *testing.T, radius float64, n int, startDeg, sweepDeg float64) {
		if n < 1 || n > nFuzzCap {
			t.Skip()
		}
		if !finite(radius) || !finite(startDeg) || !finite(sweepDeg) {
			t.Skip()
		}
		defer func() {
			if rec := recover(); rec != nil {
				t.Fatalf("layout.PolarArc(%v, %v, %v, %v) panicked: %v",
					radius, n, startDeg, sweepDeg, rec)
			}
		}()
		out := PolarArc(radius, n, startDeg, sweepDeg)
		if len(out) != n {
			t.Fatalf("layout.PolarArc len = %d, want %d", len(out), n)
		}
		for i, p := range out {
			if p.Z != 0 {
				t.Fatalf("layout.PolarArc[%d].Z = %v, want 0", i, p.Z)
			}
		}
	})
}

// FuzzGrid checks Grid output length and centroid (must be at the origin
// for finite step values).
func FuzzGrid(f *testing.F) {
	f.Add(10.0, 10.0, 4, 4)
	f.Add(0.0, 0.0, 1, 1)
	f.Add(-3.0, 7.0, 5, 2)
	f.Add(1.0, 1.0, 0, 0)
	f.Fuzz(func(t *testing.T, stepX, stepY float64, nx, ny int) {
		if nx < 0 || ny < 0 || nx > 256 || ny > 256 {
			t.Skip()
		}
		if !finite(stepX) || !finite(stepY) {
			t.Skip()
		}
		defer func() {
			if rec := recover(); rec != nil {
				t.Fatalf("layout.Grid(%v, %v, %v, %v) panicked: %v",
					stepX, stepY, nx, ny, rec)
			}
		}()
		out := Grid(stepX, stepY, nx, ny)
		if len(out) != nx*ny {
			t.Fatalf("layout.Grid len = %d, want %d", len(out), nx*ny)
		}
		// Centroid must be at the origin (z always 0).
		if nx > 0 && ny > 0 {
			var sx, sy float64
			for _, p := range out {
				sx += p.X
				sy += p.Y
				if p.Z != 0 {
					t.Fatalf("Grid produced point with Z != 0: %+v", p)
				}
			}
			// Allow generous tolerance scaled to magnitudes.
			tol := 1e-6 * (1 + math.Abs(stepX)*float64(nx) + math.Abs(stepY)*float64(ny))
			if math.Abs(sx) > tol || math.Abs(sy) > tol {
				t.Fatalf("Grid centroid (%v, %v) not at origin (tol=%v)", sx, sy, tol)
			}
		}
	})
}

// FuzzLine checks that Line returns n positions and that endpoints land on
// p0 and p1 (within fp tolerance) when n > 1.
func FuzzLine(f *testing.F) {
	f.Add(0.0, 0.0, 0.0, 10.0, 0.0, 0.0, 5)
	f.Add(1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 1)
	f.Add(-1e6, 0.0, 0.0, 1e6, 0.0, 0.0, 100)
	f.Fuzz(func(t *testing.T, ax, ay, az, bx, by, bz float64, n int) {
		if n < 1 || n > nFuzzCap {
			t.Skip()
		}
		for _, v := range []float64{ax, ay, az, bx, by, bz} {
			if !finite(v) {
				t.Skip()
			}
		}
		p0 := v3.XYZ(ax, ay, az)
		p1 := v3.XYZ(bx, by, bz)
		defer func() {
			if rec := recover(); rec != nil {
				t.Fatalf("layout.Line(%+v, %+v, %v) panicked: %v", p0, p1, n, rec)
			}
		}()
		out := Line(p0, p1, n)
		if len(out) != n {
			t.Fatalf("layout.Line len = %d, want %d", len(out), n)
		}
		if n == 1 {
			if out[0] != p0 {
				t.Fatalf("Line(_, _, 1)[0] = %+v, want p0 %+v", out[0], p0)
			}
			return
		}
		// Endpoints should be exact (multiplied by 0 and 1 respectively).
		if out[0] != p0 {
			t.Fatalf("Line endpoint[0] = %+v, want %+v", out[0], p0)
		}
		// p1 may differ in last bit due to fp; allow tiny tolerance.
		span := math.Max(math.Abs(bx-ax), math.Max(math.Abs(by-ay), math.Abs(bz-az)))
		tol := 1e-9 * (1 + span)
		end := out[n-1]
		if math.Abs(end.X-bx) > tol || math.Abs(end.Y-by) > tol || math.Abs(end.Z-bz) > tol {
			t.Fatalf("Line endpoint[n-1] = %+v, want %+v (tol=%v)", end, p1, tol)
		}
	})
}

// FuzzRectCorners checks that RectCorners always returns 4 corners summing
// to the origin.
func FuzzRectCorners(f *testing.F) {
	f.Add(80.0, 50.0)
	f.Add(0.0, 0.0)
	f.Add(-10.0, -20.0)
	f.Add(1e15, 1e15)
	f.Fuzz(func(t *testing.T, w, d float64) {
		if !finite(w) || !finite(d) {
			t.Skip()
		}
		defer func() {
			if rec := recover(); rec != nil {
				t.Fatalf("layout.RectCorners(%v, %v) panicked: %v", w, d, rec)
			}
		}()
		out := RectCorners(w, d)
		if len(out) != 4 {
			t.Fatalf("RectCorners len = %d, want 4", len(out))
		}
		var sum v3.Vec
		for _, p := range out {
			sum = sum.Add(p)
		}
		tol := 1e-9 * (1 + math.Abs(w) + math.Abs(d))
		if math.Abs(sum.X) > tol || math.Abs(sum.Y) > tol || math.Abs(sum.Z) > tol {
			t.Fatalf("RectCorners centroid = %+v, want origin (tol=%v)", sum, tol)
		}
	})
}

// FuzzBoxCorners checks BoxCorners always returns 8 corners summing to the
// origin.
func FuzzBoxCorners(f *testing.F) {
	f.Add(10.0, 10.0, 10.0)
	f.Add(0.0, 0.0, 0.0)
	f.Add(-1.0, 2.0, -3.0)
	f.Add(1e12, 1e12, 1e12)
	f.Fuzz(func(t *testing.T, sx, sy, sz float64) {
		if !finite(sx) || !finite(sy) || !finite(sz) {
			t.Skip()
		}
		size := v3.XYZ(sx, sy, sz)
		defer func() {
			if rec := recover(); rec != nil {
				t.Fatalf("layout.BoxCorners(%+v) panicked: %v", size, rec)
			}
		}()
		out := BoxCorners(size)
		if len(out) != 8 {
			t.Fatalf("BoxCorners len = %d, want 8", len(out))
		}
		var sum v3.Vec
		for _, p := range out {
			sum = sum.Add(p)
		}
		tol := 1e-9 * (1 + math.Abs(sx) + math.Abs(sy) + math.Abs(sz))
		if math.Abs(sum.X) > tol || math.Abs(sum.Y) > tol || math.Abs(sum.Z) > tol {
			t.Fatalf("BoxCorners centroid = %+v, want origin (tol=%v)", sum, tol)
		}
	})
}

// FuzzPolar2 mirrors FuzzPolar for the 2D variant.
func FuzzPolar2(f *testing.F) {
	f.Add(10.0, 6)
	f.Add(0.0, 1)
	f.Add(1.0, 0)
	f.Fuzz(func(t *testing.T, radius float64, n int) {
		if n < 0 || n > nFuzzCap {
			t.Skip()
		}
		if !finite(radius) {
			t.Skip()
		}
		defer func() {
			if rec := recover(); rec != nil {
				t.Fatalf("layout.Polar2(%v, %v) panicked: %v", radius, n, rec)
			}
		}()
		out := Polar2(radius, n)
		if len(out) != n {
			t.Fatalf("Polar2 len = %d, want %d", len(out), n)
		}
	})
}

// FuzzGrid2 mirrors FuzzGrid for the 2D variant.
func FuzzGrid2(f *testing.F) {
	f.Add(10.0, 10.0, 4, 4)
	f.Add(1.0, 1.0, 0, 0)
	f.Fuzz(func(t *testing.T, stepX, stepY float64, nx, ny int) {
		if nx < 0 || ny < 0 || nx > 256 || ny > 256 {
			t.Skip()
		}
		if !finite(stepX) || !finite(stepY) {
			t.Skip()
		}
		defer func() {
			if rec := recover(); rec != nil {
				t.Fatalf("layout.Grid2(%v, %v, %v, %v) panicked: %v",
					stepX, stepY, nx, ny, rec)
			}
		}()
		out := Grid2(stepX, stepY, nx, ny)
		if len(out) != nx*ny {
			t.Fatalf("Grid2 len = %d, want %d", len(out), nx*ny)
		}
	})
}
