package layout

import (
	"math"
	"testing"

	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// vecClose2 reports whether two 2D vecs match within eps.
func vecClose2(a, b v2.Vec) bool {
	return math.Abs(a.X-b.X) < eps && math.Abs(a.Y-b.Y) < eps
}

// distinctXY reports whether all 2D points are pairwise distinct in XY.
func distinctXY(pts []v3.Vec) bool {
	for i := 0; i < len(pts); i++ {
		for j := i + 1; j < len(pts); j++ {
			if math.Abs(pts[i].X-pts[j].X) < 1e-9 && math.Abs(pts[i].Y-pts[j].Y) < 1e-9 {
				return false
			}
		}
	}
	return true
}

// --- Polar evenly-spaced. ---

func TestPolarEvenlySpaced(t *testing.T) {
	pts := Polar(20, 8)
	if len(pts) != 8 {
		t.Fatalf("Polar(20,8) length = %d, want 8", len(pts))
	}
	if !distinctXY(pts) {
		t.Errorf("Polar produced duplicate points")
	}
	// Adjacent point separations must all be equal — chord length = 2r sin(π/n).
	want := 2 * 20 * math.Sin(math.Pi/8)
	for i := 0; i < 8; i++ {
		j := (i + 1) % 8
		d := math.Hypot(pts[i].X-pts[j].X, pts[i].Y-pts[j].Y)
		if math.Abs(d-want) > 1e-9 {
			t.Errorf("Polar[%d→%d] chord = %.6f, want %.6f", i, j, d, want)
		}
	}
}

// Polar with n=1 → single point at +X.
func TestPolarOnePoint(t *testing.T) {
	pts := Polar(15, 1)
	if len(pts) != 1 {
		t.Fatalf("Polar(15,1) length = %d, want 1", len(pts))
	}
	if math.Abs(pts[0].X-15) > 1e-9 || math.Abs(pts[0].Y) > 1e-9 {
		t.Errorf("Polar(15,1)[0] = %+v, want (15, 0, 0)", pts[0])
	}
}

// --- PolarArc. ---

func TestPolarArcSweep(t *testing.T) {
	// 90° sweep starting at 0°, 3 points → 0°, 45°, 90°.
	pts := PolarArc(10, 3, 0, 90)
	if len(pts) != 3 {
		t.Fatalf("PolarArc(10,3,0,90) length = %d, want 3", len(pts))
	}
	if math.Abs(pts[0].X-10) > 1e-9 || math.Abs(pts[0].Y) > 1e-9 {
		t.Errorf("PolarArc[0] = %+v, want (10, 0, 0) at 0°", pts[0])
	}
	if math.Abs(pts[2].X) > 1e-9 || math.Abs(pts[2].Y-10) > 1e-9 {
		t.Errorf("PolarArc[2] = %+v, want (0, 10, 0) at 90°", pts[2])
	}
}

func TestPolarArcSinglePoint(t *testing.T) {
	pts := PolarArc(10, 1, 30, 60)
	if len(pts) != 1 {
		t.Fatalf("PolarArc n=1 length = %d, want 1", len(pts))
	}
	// Single point should land at startDeg, NOT at startDeg+sweep/2.
	wantX := 10 * math.Cos(30*math.Pi/180)
	wantY := 10 * math.Sin(30*math.Pi/180)
	if math.Abs(pts[0].X-wantX) > 1e-9 || math.Abs(pts[0].Y-wantY) > 1e-9 {
		t.Errorf("PolarArc n=1 = %+v, want (%v, %v) at 30°", pts[0], wantX, wantY)
	}
}

// --- Grid: count and stepping. ---

func TestGridStepping(t *testing.T) {
	pts := Grid(10, 5, 3, 2) // 6 points
	if len(pts) != 6 {
		t.Fatalf("Grid(10,5,3,2) length = %d, want 6", len(pts))
	}
	if !distinctXY(pts) {
		t.Errorf("Grid produced duplicate points")
	}
	// X step should be exactly 10 between adjacent X-direction neighbors.
	dx := pts[1].X - pts[0].X
	if math.Abs(dx-10) > 1e-9 {
		t.Errorf("Grid X step = %.4f, want 10", dx)
	}
	// Y step between rows should be 5.
	dy := pts[3].Y - pts[0].Y
	if math.Abs(dy-5) > 1e-9 {
		t.Errorf("Grid Y step = %.4f, want 5", dy)
	}
}

func TestGridSingleCell(t *testing.T) {
	pts := Grid(10, 10, 1, 1)
	if len(pts) != 1 {
		t.Fatalf("Grid 1x1 length = %d, want 1", len(pts))
	}
	if !vecClose(pts[0], v3.XYZ(0, 0, 0)) {
		t.Errorf("Grid(10,10,1,1)[0] = %+v, want origin", pts[0])
	}
}

// --- Line: count, distinctness, endpoints. ---

func TestLineEndpoints(t *testing.T) {
	pts := Line(v3.XYZ(0, 0, 0), v3.XYZ(10, 0, 0), 5)
	if len(pts) != 5 {
		t.Fatalf("Line length = %d, want 5", len(pts))
	}
	if !vecClose(pts[0], v3.XYZ(0, 0, 0)) {
		t.Errorf("Line first = %+v, want origin", pts[0])
	}
	if !vecClose(pts[4], v3.XYZ(10, 0, 0)) {
		t.Errorf("Line last = %+v, want (10,0,0)", pts[4])
	}
}

func TestLineSinglePoint(t *testing.T) {
	pts := Line(v3.XYZ(2, 3, 4), v3.XYZ(99, 99, 99), 1)
	if len(pts) != 1 {
		t.Fatalf("Line n=1 length = %d, want 1", len(pts))
	}
	if !vecClose(pts[0], v3.XYZ(2, 3, 4)) {
		t.Errorf("Line n=1 should land at p0; got %+v, want (2,3,4)", pts[0])
	}
}

// --- RectCorners produces 4 distinct points. ---

func TestRectCornersDistinct(t *testing.T) {
	pts := RectCorners(80, 50)
	if len(pts) != 4 {
		t.Fatalf("RectCorners length = %d, want 4", len(pts))
	}
	if !distinctXY(pts) {
		t.Errorf("RectCorners produced duplicate corners: %+v", pts)
	}
}

// BoxCorners → 8 distinct points.
func TestBoxCornersDistinct(t *testing.T) {
	pts := BoxCorners(v3.XYZ(20, 10, 6))
	if len(pts) != 8 {
		t.Fatalf("BoxCorners length = %d, want 8", len(pts))
	}
	for i := 0; i < 8; i++ {
		for j := i + 1; j < 8; j++ {
			if vecClose(pts[i], pts[j]) {
				t.Errorf("BoxCorners[%d] == BoxCorners[%d] = %+v", i, j, pts[i])
			}
		}
	}
}

// --- 2D variants. ---

func TestPolar2(t *testing.T) {
	pts := Polar2(20, 6)
	if len(pts) != 6 {
		t.Fatalf("Polar2 length = %d, want 6", len(pts))
	}
	if !vecClose2(pts[0], v2.XY(20, 0)) {
		t.Errorf("Polar2[0] = %+v, want (20, 0)", pts[0])
	}
}

func TestGrid2(t *testing.T) {
	pts := Grid2(10, 10, 3, 3)
	if len(pts) != 9 {
		t.Fatalf("Grid2 3x3 length = %d, want 9", len(pts))
	}
}

func TestRectCorners2(t *testing.T) {
	pts := RectCorners2(80, 50)
	if len(pts) != 4 {
		t.Fatalf("RectCorners2 length = %d, want 4", len(pts))
	}
	// Centroid should be origin.
	var sx, sy float64
	for _, p := range pts {
		sx += p.X
		sy += p.Y
	}
	if math.Abs(sx) > eps || math.Abs(sy) > eps {
		t.Errorf("RectCorners2 centroid = (%v, %v), want origin", sx, sy)
	}
}

func TestLine2Endpoints(t *testing.T) {
	pts := Line2(v2.XY(0, 0), v2.XY(10, 0), 3)
	if len(pts) != 3 {
		t.Fatalf("Line2 length = %d, want 3", len(pts))
	}
	if !vecClose2(pts[0], v2.XY(0, 0)) {
		t.Errorf("Line2 first = %+v, want origin", pts[0])
	}
	if !vecClose2(pts[2], v2.XY(10, 0)) {
		t.Errorf("Line2 last = %+v, want (10,0)", pts[2])
	}
}
