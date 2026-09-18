// Clearance, interference and congruence between two solids.
//
// Every measurement in this file works the same way: sample one solid's
// SURFACE (marching-cubes vertices at cellsPerMM) and read the OTHER solid's
// signed field there. Negative means that surface point is inside the other
// solid's material, positive means it is in air that distance away.
//
// Compare a surface to a field, never two fields. An SDF value is an exact
// distance only for primitives and rigid transforms of them; after booleans,
// twists, scales or Correct() it is a Lipschitz BOUND, so a reading can
// understate the true distance but never get the sign wrong. That is why a
// gap of +0.30 mm should be read as "at least 0.30 mm of air", and why the
// Require* helpers here assert lower bounds on clearance and upper bounds on
// penetration rather than exact distances.

package validate

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/snowbldr/fluent-sdfx/mesh"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Gap describes how one solid's surface sits relative to another solid's
// field: every sampled value is the second solid's signed field (mm) at a
// surface point of the first. Negative means the first solid penetrates the
// second; positive is clear air. Min is the worst case.
type Gap struct {
	Min, Median, P95, Max float64 // mm, signed: negative = penetration
	MinAt, MaxAt          v3.Vec  // the sample points where Min and Max occurred
	N                     int     // number of surface points sampled
}

// String renders the gap as one test-log line: worst and best readings with
// the points they occurred at, plus the median and 95th percentile.
func (g Gap) String() string {
	if g.N == 0 {
		return "gap: no surface points sampled"
	}
	return fmt.Sprintf("gap min=%+.3f at (%.3f, %.3f, %.3f) median=%+.3f p95=%+.3f max=%+.3f at (%.3f, %.3f, %.3f) n=%d",
		g.Min, g.MinAt.X, g.MinAt.Y, g.MinAt.Z,
		g.Median, g.P95,
		g.Max, g.MaxAt.X, g.MaxAt.Y, g.MaxAt.Z, g.N)
}

// Clearance samples a's surface at cellsPerMM and reads b's signed field at
// every one of those points: Min is the worst clearance in mm (negative means
// a penetrates b). Nothing is meshed for b, so b may be any solid.
func Clearance(a, b *solid.Solid, cellsPerMM float64) Gap {
	pts := Surface(a, cellsPerMM)
	return gapOf(evalAll(b, pts), pts)
}

// ClearanceIn is Clearance restricted to the surface points of a that lie
// inside region (an axis-aligned box in world coordinates) — the way to ask
// about one feature of a big part without meshing it differently.
func ClearanceIn(a, b *solid.Solid, region v3.Box, cellsPerMM float64) Gap {
	pts := gapWithin(Surface(a, cellsPerMM), region)
	return gapOf(evalAll(b, pts), pts)
}

// RequireClearance fails the test unless every surface point of a is at least
// minGap mm clear of b's material (g.Min >= minGap). Pass a negative minGap to
// allow a known interference fit.
func RequireClearance(t testing.TB, a, b *solid.Solid, cellsPerMM, minGap float64) {
	t.Helper()
	g := Clearance(a, b, cellsPerMM)
	t.Logf("clearance: %s (want min >= %+.3f)", g, minGap)
	if g.N == 0 {
		t.Errorf("clearance: a has no surface points at cellsPerMM=%g", cellsPerMM)
		return
	}
	if g.Min < minGap {
		t.Errorf("clearance min = %+.3f mm at (%.3f, %.3f, %.3f), want >= %+.3f mm",
			g.Min, g.MinAt.X, g.MinAt.Y, g.MinAt.Z, minGap)
	}
}

// RequireNoContact fails the test unless a's surface is strictly outside b
// everywhere (g.Min > 0): the two bodies do not touch at the sampled density.
func RequireNoContact(t testing.TB, a, b *solid.Solid, cellsPerMM float64) {
	t.Helper()
	g := Clearance(a, b, cellsPerMM)
	t.Logf("contact: %s (want min > 0)", g)
	if g.N == 0 {
		t.Errorf("contact: a has no surface points at cellsPerMM=%g", cellsPerMM)
		return
	}
	if g.Min <= 0 {
		t.Errorf("bodies touch or overlap: min = %+.3f mm at (%.3f, %.3f, %.3f), want > 0",
			g.Min, g.MinAt.X, g.MinAt.Y, g.MinAt.Z)
	}
}

// RequireInterference fails the test unless a penetrates b by at least
// atLeast mm somewhere (g.Min <= -atLeast). This is the anti-flip guard for a
// rejection test: it proves the clearance check above would actually have
// caught a collision, rather than passing because the sample missed.
func RequireInterference(t testing.TB, a, b *solid.Solid, cellsPerMM, atLeast float64) {
	t.Helper()
	g := Clearance(a, b, cellsPerMM)
	t.Logf("interference: %s (want min <= %+.3f)", g, -atLeast)
	if g.N == 0 {
		t.Errorf("interference: a has no surface points at cellsPerMM=%g", cellsPerMM)
		return
	}
	if g.Min > -atLeast {
		t.Errorf("expected interference of at least %.3f mm, got min = %+.3f mm at (%.3f, %.3f, %.3f)",
			atLeast, g.Min, g.MinAt.X, g.MinAt.Y, g.MinAt.Z)
	}
}

// Congruent reports whether a and b occupy the same surface: ab reads b's
// field on a's surface, ba reads a's field on b's surface, and ok is true
// when the 95th percentile of |field| is within tol mm in BOTH directions.
// Both directions are needed — one surface can lie entirely on another that
// extends far beyond it.
func Congruent(a, b *solid.Solid, cellsPerMM, tol float64) (ab, ba Gap, ok bool) {
	aPts := Surface(a, cellsPerMM)
	bPts := Surface(b, cellsPerMM)
	aVals := evalAll(b, aPts)
	bVals := evalAll(a, bPts)
	ab, ba = gapOf(aVals, aPts), gapOf(bVals, bPts)
	if len(aPts) == 0 || len(bPts) == 0 {
		return ab, ba, false
	}
	return ab, ba, gapAbsP95(aVals) <= tol && gapAbsP95(bVals) <= tol
}

// RequireCongruent fails the test unless a and b describe the same surface
// within tol mm (see Congruent). Use it as a refactor guard: the rewritten
// builder must produce the geometry the old one did.
func RequireCongruent(t testing.TB, a, b *solid.Solid, cellsPerMM, tol float64) {
	t.Helper()
	ab, ba, ok := Congruent(a, b, cellsPerMM, tol)
	t.Logf("congruent a->b: %s", ab)
	t.Logf("congruent b->a: %s", ba)
	if !ok {
		t.Errorf("solids are not congruent within %.3f mm: a->b p95|f|=%.3f, b->a p95|f|=%.3f",
			tol, gapAbsP95Of(ab), gapAbsP95Of(ba))
	}
}

// DeviationSTL measures an already-written mesh file against the solid it was
// rendered from: it reads the STL at path (binary or ASCII), deduplicates the
// vertices, and reports |field of s| at each one. Every value is an absolute
// distance in mm, so Min/Median/P95/Max are all >= 0 and Max is the worst
// deviation of the exported triangle soup from the true surface — the check
// for "did decimation ruin this part".
func DeviationSTL(path string, s *solid.Solid) (Gap, error) {
	tris, err := gapReadSTL(path)
	if err != nil {
		return Gap{}, err
	}
	pts := Vertices(tris)
	vals := evalAll(s, pts)
	for i := range vals {
		vals[i] = math.Abs(vals[i])
	}
	return gapOf(vals, pts), nil
}

// --- unexported helpers ---

// gapOf summarises signed field readings vals taken at points pts.
func gapOf(vals []float64, pts []v3.Vec) Gap {
	if len(vals) == 0 {
		return Gap{}
	}
	g := Gap{Min: math.Inf(1), Max: math.Inf(-1), N: len(vals)}
	for i, v := range vals {
		if v < g.Min {
			g.Min, g.MinAt = v, pts[i]
		}
		if v > g.Max {
			g.Max, g.MaxAt = v, pts[i]
		}
	}
	g.Median = Percentile(vals, 50)
	g.P95 = Percentile(vals, 95)
	return g
}

// gapWithin keeps the points of pts inside region (inclusive).
func gapWithin(pts []v3.Vec, region v3.Box) []v3.Vec {
	out := make([]v3.Vec, 0, len(pts))
	for _, p := range pts {
		if p.X >= region.Min.X && p.X <= region.Max.X &&
			p.Y >= region.Min.Y && p.Y <= region.Max.Y &&
			p.Z >= region.Min.Z && p.Z <= region.Max.Z {
			out = append(out, p)
		}
	}
	return out
}

// gapAbsP95 is the 95th percentile of |vals|.
func gapAbsP95(vals []float64) float64 {
	abs := make([]float64, len(vals))
	for i, v := range vals {
		abs[i] = math.Abs(v)
	}
	return Percentile(abs, 95)
}

// gapAbsP95Of estimates the p95 of |field| from a summarised Gap, for the
// failure message only: for a near-congruent pair the extremes bound it.
func gapAbsP95Of(g Gap) float64 {
	return math.Max(math.Abs(g.Min), math.Abs(g.P95))
}

// gapReadSTL reads a binary or ASCII STL file into triangles.
func gapReadSTL(path string) ([]mesh.Triangle3, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if gapIsASCIISTL(data) {
		return gapParseASCIISTL(data)
	}
	return gapParseBinarySTL(data)
}

// gapIsASCIISTL reports whether the bytes look like an ASCII STL: the file
// starts with "solid" AND its length does not match the binary layout.
func gapIsASCIISTL(data []byte) bool {
	if len(data) < 84 {
		return len(data) >= 5 && strings.HasPrefix(strings.TrimSpace(string(data[:min(5, len(data))])), "solid")
	}
	n := binary.LittleEndian.Uint32(data[80:84])
	if int64(len(data)) == 84+50*int64(n) {
		return false
	}
	return strings.HasPrefix(strings.TrimSpace(string(data[:5])), "solid")
}

// gapParseBinarySTL parses the binary layout: an 80-byte header, a uint32
// little-endian triangle count at offset 80, then 50 bytes per triangle
// (12 float32s: normal then three vertices, plus a uint16 attribute count).
func gapParseBinarySTL(data []byte) ([]mesh.Triangle3, error) {
	if len(data) < 84 {
		return nil, fmt.Errorf("stl: file is %d bytes, too short for a binary STL header", len(data))
	}
	n := binary.LittleEndian.Uint32(data[80:84])
	want := int64(84) + 50*int64(n)
	if int64(len(data)) < want {
		return nil, fmt.Errorf("stl: header claims %d triangles (%d bytes) but file is %d bytes", n, want, len(data))
	}
	tris := make([]mesh.Triangle3, n)
	for i := range tris {
		rec := data[84+50*i : 84+50*i+50]
		for v := 0; v < 3; v++ {
			off := 12 + 12*v // skip the normal
			tris[i][v] = v3.XYZ(
				float64(math.Float32frombits(binary.LittleEndian.Uint32(rec[off:off+4]))),
				float64(math.Float32frombits(binary.LittleEndian.Uint32(rec[off+4:off+8]))),
				float64(math.Float32frombits(binary.LittleEndian.Uint32(rec[off+8:off+12]))),
			)
		}
	}
	return tris, nil
}

// gapParseASCIISTL parses the text layout, reading every "vertex x y z" line
// in groups of three.
func gapParseASCIISTL(data []byte) ([]mesh.Triangle3, error) {
	var tris []mesh.Triangle3
	var cur mesh.Triangle3
	k := 0
	for ln, line := range strings.Split(string(data), "\n") {
		f := strings.Fields(line)
		if len(f) == 0 || f[0] != "vertex" {
			continue
		}
		if len(f) < 4 {
			return nil, fmt.Errorf("stl: line %d: malformed vertex %q", ln+1, line)
		}
		var xyz [3]float64
		for i := 0; i < 3; i++ {
			v, err := strconv.ParseFloat(f[i+1], 64)
			if err != nil {
				return nil, fmt.Errorf("stl: line %d: %w", ln+1, err)
			}
			xyz[i] = v
		}
		cur[k] = v3.XYZ(xyz[0], xyz[1], xyz[2])
		if k++; k == 3 {
			tris = append(tris, cur)
			k = 0
		}
	}
	if k != 0 {
		return nil, fmt.Errorf("stl: file ends with %d leftover vertices", k)
	}
	return tris, nil
}
