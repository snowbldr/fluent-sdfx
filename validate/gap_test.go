package validate

import (
	"encoding/binary"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// gapCube is a cube of the given side centred at (cx, 0, 0): its faces sit at
// x = cx ± side/2, so two of them have an analytically known gap.
func gapCube(side, cx float64) *solid.Solid {
	return solid.Box(v3.XYZ(side, side, side), 0).TranslateX(cx)
}

const gapCPM = 4.0 // cells/mm for the gap tests: cheap and exact on flat faces

func TestClearanceKnownGaps(t *testing.T) {
	// a spans x[-5,5]; the partner cube is placed so the gap is exact.
	a := gapCube(10, 0)
	tests := []struct {
		name    string
		b       *solid.Solid
		wantMin float64 // b's field at a's nearest surface point
		wantMax float64 // b's field at a's farthest surface point
	}{
		// b spans x[7,17]: a's +X face is 2 mm clear, its -X face 12 mm.
		{"clear by 2mm", gapCube(10, 12), 2, 12},
		// b spans x[10.5,20.5]: gap 5.5.
		{"clear by 5.5mm", gapCube(10, 15.5), 5.5, 15.5},
		// b spans x[4,14]: a's +X face is 1 mm INSIDE b.
		{"overlapping by 1mm", gapCube(10, 9), -1, 9},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := Clearance(a, tc.b, gapCPM)
			if g.N == 0 {
				t.Fatal("no surface points sampled")
			}
			if math.Abs(g.Min-tc.wantMin) > 0.1 {
				t.Errorf("Min = %+.3f, want %+.3f (%s)", g.Min, tc.wantMin, g)
			}
			if math.Abs(g.Max-tc.wantMax) > 0.1 {
				t.Errorf("Max = %+.3f, want %+.3f (%s)", g.Max, tc.wantMax, g)
			}
			if g.Median < g.Min || g.Median > g.Max || g.P95 < g.Median || g.P95 > g.Max {
				t.Errorf("percentiles out of order: %s", g)
			}
		})
	}
}

func TestClearanceMinAtIsOnTheNearFace(t *testing.T) {
	a, b := gapCube(10, 0), gapCube(10, 12)
	g := Clearance(a, b, gapCPM)
	if math.Abs(g.MinAt.X-5) > 0.1 {
		t.Errorf("MinAt = (%.3f, %.3f, %.3f), want the +X face at x=5", g.MinAt.X, g.MinAt.Y, g.MinAt.Z)
	}
	if math.Abs(g.MaxAt.X+5) > 0.1 {
		t.Errorf("MaxAt = (%.3f, %.3f, %.3f), want the -X face at x=-5", g.MaxAt.X, g.MaxAt.Y, g.MaxAt.Z)
	}
}

func TestClearanceIn(t *testing.T) {
	a, b := gapCube(10, 0), gapCube(10, 12)
	full := Clearance(a, b, gapCPM)

	// Restrict to a's -X half: the nearest surface there is the x=-5 face.
	farHalf := v3.Box{Min: v3.XYZ(-6, -6, -6), Max: v3.XYZ(-4.9, 6, 6)}
	g := ClearanceIn(a, b, farHalf, gapCPM)
	if g.N == 0 || g.N >= full.N {
		t.Fatalf("region sample = %d points, want 0 < n < %d", g.N, full.N)
	}
	if math.Abs(g.Min-12) > 0.1 {
		t.Errorf("Min in -X half = %+.3f, want +12.000 (%s)", g.Min, g)
	}

	// A region containing nothing yields an empty Gap that says so.
	empty := ClearanceIn(a, b, v3.Box{Min: v3.XYZ(100, 100, 100), Max: v3.XYZ(101, 101, 101)}, gapCPM)
	if empty.N != 0 {
		t.Errorf("N = %d for an empty region, want 0", empty.N)
	}
	if !strings.Contains(empty.String(), "no surface points") {
		t.Errorf("String() = %q, want the empty-sample message", empty.String())
	}
}

func TestGapString(t *testing.T) {
	g := Gap{Min: -0.5, Median: 1, P95: 2, Max: 3, MinAt: v3.XYZ(1, 2, 3), MaxAt: v3.XYZ(-1, -2, -3), N: 7}
	want := "gap min=-0.500 at (1.000, 2.000, 3.000) median=+1.000 p95=+2.000 max=+3.000 at (-1.000, -2.000, -3.000) n=7"
	if got := g.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestRequireClearanceAndContact(t *testing.T) {
	a, clear, touching := gapCube(10, 0), gapCube(10, 12), gapCube(10, 9)

	// Positive: 2 mm of air really is at least 1.5 mm of air, and no contact.
	RequireClearance(t, a, clear, gapCPM, 1.5)
	RequireNoContact(t, a, clear, gapCPM)

	// Negative, via the pure form (a real Require* failure would fail this test):
	if g := Clearance(a, clear, gapCPM); g.Min >= 2.5 {
		t.Errorf("min = %+.3f, expected RequireClearance(..., 2.5) to fail", g.Min)
	}
	if g := Clearance(a, touching, gapCPM); g.Min > 0 {
		t.Errorf("min = %+.3f, expected RequireNoContact to fail on overlapping cubes", g.Min)
	}
}

func TestRequireInterference(t *testing.T) {
	a, overlapping, clear := gapCube(10, 0), gapCube(10, 9), gapCube(10, 12)

	// Positive: a penetrates the overlapping cube by 1 mm.
	RequireInterference(t, a, overlapping, gapCPM, 0.5)

	// Negative: separated cubes do not interfere at all.
	if g := Clearance(a, clear, gapCPM); g.Min <= 0 {
		t.Errorf("min = %+.3f, expected RequireInterference to fail on clear cubes", g.Min)
	}
}

func TestCongruent(t *testing.T) {
	a := gapCube(10, 0)
	tests := []struct {
		name string
		b    *solid.Solid
		tol  float64
		want bool
	}{
		{"same cube", gapCube(10, 0), 0.05, true},
		{"same cube, rebuilt via union with itself", gapCube(10, 0).Union(gapCube(10, 0)), 0.05, true},
		{"0.6mm bigger", gapCube(10.6, 0), 0.05, false},
		{"0.6mm bigger, loose tol", gapCube(10.6, 0), 0.5, true},
		{"shifted 1mm", gapCube(10, 1), 0.05, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ab, ba, ok := Congruent(a, tc.b, gapCPM, tc.tol)
			if ok != tc.want {
				t.Errorf("ok = %v, want %v\n  a->b: %s\n  b->a: %s", ok, tc.want, ab, ba)
			}
			if ab.N == 0 || ba.N == 0 {
				t.Errorf("empty sample: a->b n=%d b->a n=%d", ab.N, ba.N)
			}
		})
	}
}

func TestRequireCongruent(t *testing.T) {
	// Positive: an intersection with a box that contains it is the same solid.
	a := gapCube(10, 0)
	b := a.Intersect(gapCube(30, 0))
	RequireCongruent(t, a, b, gapCPM, 0.05)

	// Negative, via the pure form.
	if _, _, ok := Congruent(a, gapCube(12, 0), gapCPM, 0.05); ok {
		t.Error("a 12mm cube reported congruent with a 10mm cube")
	}
}

func TestDeviationSTL(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cube.stl")
	s := gapCube(10, 0)
	s.STL(path, gapCPM)

	// The file really is binary STL in the layout DeviationSTL expects.
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) < 84 {
		t.Fatalf("stl is %d bytes", len(raw))
	}
	n := binary.LittleEndian.Uint32(raw[80:84])
	if got, want := int64(len(raw)), int64(84)+50*int64(n); got != want {
		t.Fatalf("stl is %d bytes, want %d for %d triangles", got, want, n)
	}

	// Positive: the mesh sits on the solid it came from.
	g, err := DeviationSTL(path, s)
	if err != nil {
		t.Fatal(err)
	}
	if g.N == 0 {
		t.Fatal("no vertices read")
	}
	if g.N >= int(n)*3 {
		t.Errorf("N = %d, want fewer than %d (vertices should be deduplicated)", g.N, int(n)*3)
	}
	if g.Min < 0 {
		t.Errorf("Min = %+.3f, want >= 0 (values are absolute deviations)", g.Min)
	}
	if g.Max > 0.1 {
		t.Errorf("worst deviation = %.3f mm, want <= 0.1 mm (%s)", g.Max, g)
	}

	// Negative: measured against a cube 2 mm bigger, the same mesh is 1 mm off.
	off, err := DeviationSTL(path, gapCube(12, 0))
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(off.Median-1) > 0.1 {
		t.Errorf("median deviation vs a 12mm cube = %.3f, want ~1.000 (%s)", off.Median, off)
	}
}

func TestDeviationSTLASCII(t *testing.T) {
	// One triangle on the +X face of the 10 mm cube, written as ASCII STL.
	const ascii = `solid one
facet normal 1 0 0
  outer loop
    vertex 5 0 0
    vertex 5 1 0
    vertex 5 0 1
  endloop
endfacet
endsolid one
`
	path := filepath.Join(t.TempDir(), "face.stl")
	if err := os.WriteFile(path, []byte(ascii), 0o600); err != nil {
		t.Fatal(err)
	}
	g, err := DeviationSTL(path, gapCube(10, 0))
	if err != nil {
		t.Fatal(err)
	}
	if g.N != 3 {
		t.Errorf("N = %d, want 3", g.N)
	}
	if g.Max > 1e-9 {
		t.Errorf("deviation = %.9f, want 0: all three vertices lie on the face", g.Max)
	}
}

func TestDeviationSTLErrors(t *testing.T) {
	dir := t.TempDir()
	short := filepath.Join(dir, "short.bin")
	if err := os.WriteFile(short, make([]byte, 10), 0o600); err != nil {
		t.Fatal(err)
	}
	truncated := filepath.Join(dir, "truncated.stl")
	buf := make([]byte, 84)
	binary.LittleEndian.PutUint32(buf[80:84], 100) // claims 100 triangles, has none
	if err := os.WriteFile(truncated, buf, 0o600); err != nil {
		t.Fatal(err)
	}
	badASCII := filepath.Join(dir, "bad.stl")
	if err := os.WriteFile(badASCII, []byte("solid x\nvertex 1 2 zzz\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct{ name, path string }{
		{"missing file", filepath.Join(dir, "nope.stl")},
		{"too short", short},
		{"truncated", truncated},
		{"unparsable ascii", badASCII},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := DeviationSTL(tc.path, gapCube(10, 0)); err == nil {
				t.Error("want an error, got nil")
			}
		})
	}
}
