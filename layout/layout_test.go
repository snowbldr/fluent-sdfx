package layout

import (
	"math"
	"testing"

	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const eps = 1e-9

func vecClose(a, b v3.Vec) bool {
	return math.Abs(a.X-b.X) < eps && math.Abs(a.Y-b.Y) < eps && math.Abs(a.Z-b.Z) < eps
}

func TestPolarFirstAtPlusX(t *testing.T) {
	pts := Polar(20, 6)
	if len(pts) != 6 {
		t.Fatalf("Polar(20,6) length = %d, want 6", len(pts))
	}
	if !vecClose(pts[0], v3.XYZ(20, 0, 0)) {
		t.Fatalf("Polar(20,6)[0] = %+v, want (20,0,0)", pts[0])
	}
	for _, p := range pts {
		r := math.Hypot(p.X, p.Y)
		if math.Abs(r-20) > 1e-9 {
			t.Fatalf("Polar produced point off radius: %+v r=%v", p, r)
		}
	}
}

func TestRectCornersCentered(t *testing.T) {
	pts := RectCorners(80, 50)
	if len(pts) != 4 {
		t.Fatalf("RectCorners length = %d, want 4", len(pts))
	}
	var sum v3.Vec
	for _, p := range pts {
		sum = sum.Add(p)
	}
	if math.Abs(sum.X) > eps || math.Abs(sum.Y) > eps || math.Abs(sum.Z) > eps {
		t.Fatalf("RectCorners centroid = %+v, want origin", sum)
	}
}

func TestGridLengthAndCentroid(t *testing.T) {
	pts := Grid(10, 10, 4, 4)
	if len(pts) != 16 {
		t.Fatalf("Grid(10,10,4,4) length = %d, want 16", len(pts))
	}
	var sum v3.Vec
	for _, p := range pts {
		sum = sum.Add(p)
	}
	if math.Abs(sum.X) > 1e-9 || math.Abs(sum.Y) > 1e-9 {
		t.Fatalf("Grid centroid = %+v, want origin", sum)
	}
}
