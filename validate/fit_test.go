package validate

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// fitPinAndBlock is a pin of the given radius in a block with a 5 mm-radius
// through-bore: the smallest assembly with a clearance worth measuring.
func fitPinAndBlock(pinR float64) []Part {
	block := solid.Box(v3.XYZ(20, 20, 10), 0).Cut(solid.Cylinder(12, 5, 0))
	pin := solid.Cylinder(8, pinR, 0)
	return []Part{{Name: "block", Solid: block}, {Name: "pin", Solid: pin}}
}

func TestFitCheckWritesAndReports(t *testing.T) {
	dir := t.TempDir()
	gaps, err := FitCheck(dir, "pinfit", fitPinAndBlock(4.8), v3.Vec{}, v3.Y(1), 8)
	if err != nil {
		t.Fatalf("FitCheck: %v", err)
	}
	for _, f := range []string{"pinfit_assembly.stl", "pinfit_cutaway.stl"} {
		st, err := os.Stat(filepath.Join(dir, f))
		if err != nil {
			t.Fatalf("%s: %v", f, err)
		}
		if st.Size() < 84 {
			t.Errorf("%s is %d bytes — no triangles written", f, st.Size())
		}
	}
	if len(gaps) != 2 {
		t.Fatalf("gaps = %v, want both ordered pairs", gaps)
	}
	g, ok := gaps["pin>block"]
	if !ok {
		t.Fatalf("missing key pin>block in %v", gaps)
	}
	t.Log("pin>block:", g)
	if g.N == 0 {
		t.Fatal("no surface points sampled")
	}
	// The radial clearance is 5.0 - 4.8 = 0.2 mm; the rendered surface is
	// within about half a cell of the true one.
	if g.Min <= 0 || g.Min > 0.35 {
		t.Errorf("pin>block min clearance = %+.3f mm, want ≈+0.2 (clear of the bore)", g.Min)
	}
	if _, ok := gaps["block>pin"]; !ok {
		t.Errorf("missing the reverse key in %v", gaps)
	}
}

func TestFitCheckInterferenceIsNegative(t *testing.T) {
	dir := t.TempDir()
	gaps, err := FitCheck(dir, "tight", fitPinAndBlock(5.4), v3.Vec{}, v3.Vec{}, 8)
	if err != nil {
		t.Fatalf("FitCheck: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "tight_cutaway.stl")); !os.IsNotExist(err) {
		t.Errorf("a zero normal must not write a cutaway (stat err = %v)", err)
	}
	g := gaps["pin>block"]
	t.Log("pin>block:", g)
	if g.Min > -0.2 {
		t.Errorf("a 5.4 mm pin in a 5.0 mm bore should read ≈-0.4 mm, got %+.3f", g.Min)
	}
}

func TestFitCheckErrors(t *testing.T) {
	dir := t.TempDir()
	block := solid.Box(v3.XYZ(10, 10, 10), 0)
	cases := []struct {
		name  string
		dir   string
		parts []Part
		match string
	}{
		{"no parts", dir, nil, "no parts"},
		{"unnamed part", dir, []Part{{Solid: block}}, "no name"},
		{"nil solid", dir, []Part{{Name: "a"}}, "nil solid"},
		{"duplicate names", dir, []Part{{Name: "a", Solid: block}, {Name: "a", Solid: block}}, "duplicate"},
		{"undirectory", filepath.Join(dir, "notadir", "x"), []Part{{Name: "a", Solid: block}}, ""},
	}
	if err := os.WriteFile(filepath.Join(dir, "notadir"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := FitCheck(tc.dir, "x", tc.parts, v3.Vec{}, v3.Vec{}, 4)
			if err == nil {
				t.Fatal("want an error")
			}
			if !strings.Contains(err.Error(), tc.match) {
				t.Errorf("err = %v, want it to mention %q", err, tc.match)
			}
		})
	}
}

func TestManifestTable(t *testing.T) {
	parts := []Part{
		{Name: "cube", Solid: solid.Box(v3.XYZ(10, 10, 10), 0)},
		{Name: "pair", Solid: ohBox(6, 6, 10, 0).Translate(v3.X(-10)).Union(ohBox(6, 6, 10, 0).Translate(v3.X(10)))},
		{Name: "missing"},
	}
	md := Manifest(parts, 5)
	t.Log("\n" + md)

	lines := strings.Split(strings.TrimSpace(md), "\n")
	if len(lines) != 5 {
		t.Fatalf("got %d lines, want a header, a rule and 3 rows:\n%s", len(lines), md)
	}
	if !strings.HasPrefix(lines[0], "| name | min | max | size | volume mm³ | bodies |") {
		t.Errorf("header = %q", lines[0])
	}
	if !strings.HasPrefix(lines[1], "| --- |") {
		t.Errorf("rule = %q", lines[1])
	}
	for _, want := range []string{"| cube |", "(-5.000, -5.000, -5.000)", "10.000 × 10.000 × 10.000"} {
		if !strings.Contains(lines[2], want) {
			t.Errorf("cube row %q missing %q", lines[2], want)
		}
	}
	if !strings.HasSuffix(strings.TrimSpace(lines[2]), "| 1 |") {
		t.Errorf("cube row should end with 1 body: %q", lines[2])
	}
	if !strings.HasSuffix(strings.TrimSpace(lines[3]), "| 2 |") {
		t.Errorf("pair row should end with 2 bodies: %q", lines[3])
	}
	if !strings.Contains(lines[4], "| missing | — |") {
		t.Errorf("nil-solid row = %q", lines[4])
	}
}

func TestFitVolumeEstimate(t *testing.T) {
	cases := []struct {
		name string
		s    *solid.Solid
		want float64
	}{
		{"10 mm cube", solid.Box(v3.XYZ(10, 10, 10), 0), 1000},
		{"20x10x5 box", solid.Box(v3.XYZ(20, 10, 5), 0), 1000},
		{"r=5 sphere", solid.Sphere(5), 4.0 / 3.0 * math.Pi * 125},
		{"r=4 h=10 cylinder", solid.Cylinder(10, 4, 0), math.Pi * 16 * 10},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := fitVolume(tc.s, 8)
			if rel := math.Abs(got-tc.want) / tc.want; rel > 0.02 {
				t.Errorf("volume ≈ %.1f mm³, want %.1f (%.1f%% off)", got, tc.want, rel*100)
			}
		})
	}
	// The sample cap keeps a big part at a fine density cheap; the estimate
	// stays within a few percent.
	big := solid.Box(v3.XYZ(200, 200, 200), 0)
	if got := fitVolume(big, 20); math.Abs(got-8e6)/8e6 > 0.02 {
		t.Errorf("capped-grid volume = %.0f mm³, want ≈8000000", got)
	}
	if got := fitVolume(solid.Box(v3.XYZ(10, 10, 10), 0), 0); got != 0 {
		t.Errorf("fitVolume with cellsPerMM=0 = %v, want 0", got)
	}
}
