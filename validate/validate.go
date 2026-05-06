// Package validate inspects a solid's rendered mesh for printability and
// regression-test signals: triangle count, surface area, volume, watertightness
// (boundary-edge count), and overhang area for FDM-style 3D printing.
//
// All metrics work on the marching-cubes mesh, so they reflect what would
// actually be exported. cellsPerMM controls render density (same units as the
// STL methods on *solid.Solid). Higher density → more accurate metrics, slower.
//
// The Require* helpers integrate with Go's testing package so you can guard
// real geometric properties from a unit test:
//
//	func TestBracketPrintable(t *testing.T) {
//		s := buildBracket()
//		validate.RequireWatertight(t, s, 5.0)
//		validate.RequireMaxOverhang(t, s, 5.0, 45.0)  // FDM rule of thumb
//		validate.RequireVolumeNear(t, s, 5.0, 12500, 0.02)  // mm³ (12.5 cm³)
//	}
package validate

import (
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/mesh"
	"github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Stats reports validation metrics on a solid's rendered mesh.
type Stats struct {
	Triangles     int          // triangle count after marching cubes
	SurfaceArea   float64      // mm² — sum of triangle areas
	Volume        float64      // mm³ — signed-tetrahedron volume sum
	BoundaryEdges int          // edges shared by exactly 1 triangle (0 = closed mesh)
	Watertight    bool         // BoundaryEdges == 0
	Bounds        solid.Box3   // axis-aligned bounding box
	OverhangArea  float64      // mm² overhanging > 45° from vertical (FDM threshold)
}

// Of computes validation stats by rendering s at cellsPerMM density.
// The bounds reflect the SDF's bounding box (not the rendered mesh).
func Of(s *solid.Solid, cellsPerMM float64) Stats {
	tris := mesh.CollectTriangles(s, render.NewMarchingCubesOctreeParallel(solid.CellsFor(s, cellsPerMM)))
	st := OfMesh(tris)
	st.Bounds = s.Bounds()
	return st
}

// OfMesh computes validation stats from a precomputed triangle mesh. Bounds
// is left zero-valued — Of(*solid.Solid, ...) fills it in.
func OfMesh(tris []mesh.Triangle3) Stats {
	bnd := mesh.CountBoundaryEdges(tris)
	return Stats{
		Triangles:     len(tris),
		SurfaceArea:   SurfaceArea(tris),
		Volume:        Volume(tris),
		BoundaryEdges: bnd,
		Watertight:    bnd == 0,
		OverhangArea:  OverhangArea(tris, 45.0),
	}
}

// SurfaceArea returns the total area of all triangles in mm².
func SurfaceArea(tris []mesh.Triangle3) float64 {
	var total float64
	for _, t := range tris {
		ab := t[1].Sub(t[0])
		ac := t[2].Sub(t[0])
		total += 0.5 * ab.Cross(ac).Length()
	}
	return total
}

// Volume returns the enclosed volume in mm³, computed as the sum of signed
// tetrahedron volumes from the origin. Accurate for any closed orientable
// mesh regardless of where the origin sits.
//
// Returns garbage for non-watertight meshes — call IsWatertight first if
// you don't trust the input.
func Volume(tris []mesh.Triangle3) float64 {
	var total float64
	for _, t := range tris {
		a, b, c := t[0], t[1], t[2]
		// a · (b × c) / 6
		cross := b.Cross(c)
		total += a.Dot(cross) / 6.0
	}
	return math.Abs(total)
}

// IsWatertight reports whether the mesh is a closed surface and the
// boundary-edge count (0 = sealed). Mirror of mesh.IsWatertight kept here
// so validation users only need one import.
func IsWatertight(tris []mesh.Triangle3) (bool, int) {
	return mesh.IsWatertight(tris)
}

// OverhangArea returns the total area of triangles whose downward-facing
// normal makes an angle greater than maxAngleDeg from vertical (Z+).
//
// FDM rule of thumb: faces with overhang angle > 45° typically need supports.
// Slicers report this same angle; passing 45.0 here matches Cura/PrusaSlicer
// defaults. Pass 0 to count any downward-facing area; 90 to count only
// horizontal ceilings.
//
// Z+ is assumed to be the print direction. For non-default orientations, run
// the solid through Rotate*() before validating.
func OverhangArea(tris []mesh.Triangle3, maxAngleDeg float64) float64 {
	// A triangle's outward normal n; the face overhangs if n.z < 0 (faces
	// downward). The overhang angle from vertical is asin(-n.z / |n|).
	// face overhangs > maxAngleDeg iff -n.z/|n| > sin(maxAngleDeg).
	threshold := math.Sin(maxAngleDeg * math.Pi / 180)
	var total float64
	for _, t := range tris {
		ab := t[1].Sub(t[0])
		ac := t[2].Sub(t[0])
		n := ab.Cross(ac) // length = 2 * area, direction = outward normal
		nlen := n.Length()
		if nlen == 0 {
			continue
		}
		if -n.Z/nlen > threshold {
			total += 0.5 * nlen
		}
	}
	return total
}

// OverhangFaces returns the triangles overhanging more than maxAngleDeg
// — useful for visualizing problem areas (write them to their own STL and
// load alongside the part).
func OverhangFaces(tris []mesh.Triangle3, maxAngleDeg float64) []mesh.Triangle3 {
	threshold := math.Sin(maxAngleDeg * math.Pi / 180)
	out := make([]mesh.Triangle3, 0)
	for _, t := range tris {
		n := triNormal(t)
		nlen := n.Length()
		if nlen == 0 {
			continue
		}
		if -n.Z/nlen > threshold {
			out = append(out, t)
		}
	}
	return out
}

func triNormal(t mesh.Triangle3) v3.Vec {
	return t[1].Sub(t[0]).Cross(t[2].Sub(t[0]))
}

// --- testing helpers ---

// RequireWatertight fails the test if the rendered mesh has any boundary edges.
func RequireWatertight(t *testing.T, s *solid.Solid, cellsPerMM float64) {
	t.Helper()
	st := Of(s, cellsPerMM)
	if !st.Watertight {
		t.Fatalf("solid is not watertight: %d boundary edges (%d triangles)", st.BoundaryEdges, st.Triangles)
	}
}

// RequireVolumeNear fails the test if the rendered solid's volume differs from
// expected by more than relTol (e.g. 0.02 = 2%). Useful as a regression guard:
// a refactor that accidentally inverts an SDF or breaks a boolean shows up as
// a wildly different volume.
//
// expectedMM3 must be positive — passing 0 makes the relative tolerance
// undefined and would silently pass for any computed volume.
func RequireVolumeNear(t *testing.T, s *solid.Solid, cellsPerMM, expectedMM3, relTol float64) {
	t.Helper()
	if expectedMM3 <= 0 {
		t.Fatalf("RequireVolumeNear: expectedMM3 must be > 0, got %v", expectedMM3)
	}
	st := Of(s, cellsPerMM)
	rel := math.Abs(st.Volume-expectedMM3) / expectedMM3
	if rel > relTol {
		t.Fatalf("volume = %.4f mm³, want %.4f ± %.2f%% (got %.2f%% off)",
			st.Volume, expectedMM3, relTol*100, rel*100)
	}
}

// RequireMaxOverhang fails the test if any face overhangs more than maxAngleDeg
// — i.e. the part would need supports when FDM-printed in default Z-up
// orientation. Use 45.0 for typical PLA/PETG limits.
//
// Accepts a tinyArea tolerance (mm²) to ignore numerical-noise triangles from
// marching-cubes discretisation. 0 = strict.
func RequireMaxOverhang(t *testing.T, s *solid.Solid, cellsPerMM, maxAngleDeg float64, tinyArea ...float64) {
	t.Helper()
	tol := 0.0
	if len(tinyArea) > 0 {
		tol = tinyArea[0]
	}
	tris := mesh.CollectTriangles(s, render.NewMarchingCubesOctreeParallel(solid.CellsFor(s, cellsPerMM)))
	area := OverhangArea(tris, maxAngleDeg)
	if area > tol {
		t.Fatalf("overhanging area = %.4f mm² > %.4f mm² (faces > %.1f° from vertical)", area, tol, maxAngleDeg)
	}
}
