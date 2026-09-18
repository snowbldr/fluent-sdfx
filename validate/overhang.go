package validate

import (
	"fmt"
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/mesh"
	"github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Overhang audit in the solid's CURRENT frame. The solid is audited as it
// would be exported: its bounding box minimum Z is the build plate, +Z is
// the layer direction. Rotate the solid into its print orientation before
// calling — nothing here rotates for you.
//
// The angle convention is the slicer's: 0° is a vertical wall, 90° is a
// flat ceiling, and a downward-facing face fails when its angle exceeds the
// threshold. Two exemptions keep the audit honest about what actually
// prints, and both come from a production auditor:
//
//   - near-bed: a face whose centroid sits within BedTol of the plate is a
//     first-layer feature, squished onto the sheet, and is never reported.
//   - supported-below: a face with material (or the plate) SupportDepth
//     beneath its centroid, or beneath any of four lateral offsets of
//     Lateral in ±X and ±Y, is layer-adjacent to the body below and bonds
//     to it. A gradually-curving underside (the inside of a helical tube)
//     passes this way; a genuinely cantilevered face has air in the whole
//     neighbourhood and still fails.
//
// Material means the field reads below ohMaterial (a hair inside the
// surface, so a probe landing exactly on a face does not count as support).

const (
	ohDefaultMaxDeg       = 45.0 // degrees from vertical
	ohDefaultBedTol       = 0.5  // mm
	ohDefaultSupportDepth = 1.0  // mm
	ohDefaultLateral      = 0.8  // mm
	ohMaterial            = -0.005
)

// OverhangOptions controls the overhang audit. The solid is audited AS IS:
// its minimum Z is the build plate, so rotate into the print frame first.
// Every zero field takes its documented default; pass a negative value to
// switch an exemption off.
type OverhangOptions struct {
	// MaxDeg is the fixed threshold in degrees from vertical (0 = vertical
	// wall, 90 = flat ceiling). 0 means the 45° default; use Threshold to
	// express a genuine 0° threshold. Ignored when Threshold is non-nil.
	MaxDeg float64
	// Threshold gives the allowed angle as a function of the face
	// centroid's Z (mm, world frame) — e.g. to exempt a helical band whose
	// flanks lean by construction. Non-nil wins over MaxDeg.
	Threshold func(z float64) float64
	// BedTol is the height above the plate within which a face counts as a
	// squished first-layer feature and is exempt (default 0.5 mm; negative
	// disables the near-bed exemption). A downward probe landing within
	// BedTol/2 of the plate also counts as resting on the plate.
	BedTol float64
	// SupportDepth is how far below the face centroid to look for material
	// or the plate (default 1.0 mm; negative disables the support probe).
	SupportDepth float64
	// Lateral is the extra ±X/±Y offset probed at the same depth (default
	// 0.8 mm; negative probes straight down only). A face is supported if
	// ANY of the five probes finds support.
	Lateral float64
}

// OverhangReport is the result of an overhang audit: how much downward-facing
// area is over threshold and unsupported, and the worst angle seen anywhere.
type OverhangReport struct {
	Area     float64          // mm² of unsupported faces over threshold
	Faces    int              // number of such faces
	WorstDeg float64          // worst downward angle anywhere (0 = no downward face)
	WorstAt  v3.Vec           // centroid of the worst downward face
	Failing  []mesh.Triangle3 // the unsupported over-threshold faces, for writing to their own STL
}

// String renders the audit as one log line: failing face count and area,
// the worst failing angle and where, and the worst downward angle anywhere.
func (r OverhangReport) String() string {
	worstFail, at := ohWorstFailing(r.Failing)
	if r.Faces == 0 {
		return fmt.Sprintf("overhangs: none unsupported over threshold (%.3f mm²); worst downward face anywhere %.1f° at (%.3f, %.3f, %.3f)",
			r.Area, r.WorstDeg, r.WorstAt.X, r.WorstAt.Y, r.WorstAt.Z)
	}
	return fmt.Sprintf("overhangs: %d unsupported faces over threshold, %.3f mm²; worst failing %.1f° at (%.3f, %.3f, %.3f); worst downward face anywhere %.1f° at (%.3f, %.3f, %.3f)",
		r.Faces, r.Area, worstFail, at.X, at.Y, at.Z,
		r.WorstDeg, r.WorstAt.X, r.WorstAt.Y, r.WorstAt.Z)
}

// Overhangs audits s for unsupported overhangs at mesh density cellsPerMM,
// in the solid's current frame (minimum Z is the build plate). A face is
// reported when its downward angle (0° = vertical wall, 90° = flat ceiling)
// exceeds the threshold AND it is neither near the bed nor supported from
// below — see OverhangOptions for both exemptions.
func Overhangs(s *solid.Solid, cellsPerMM float64, o OverhangOptions) OverhangReport {
	tris := mesh.CollectTriangles(s, render.NewMarchingCubesOctreeParallel(solid.CellsFor(s, cellsPerMM)))
	return ohAudit(s, tris, o)
}

// ohAudit is Overhangs on an already-rendered mesh, so a caller that
// already has the triangles (and the solid, for the support probe) pays for
// marching cubes once.
func ohAudit(s *solid.Solid, tris []mesh.Triangle3, o OverhangOptions) OverhangReport {
	thr, bedTol, depth, lateral := ohResolve(o)
	bedZ := s.Bounds().Min.Z

	type ohFace struct {
		deg      float64
		centroid v3.Vec
		area     float64
		downward bool
		failing  bool
	}
	faces := make([]ohFace, len(tris))
	parallel(len(tris), func(i int) {
		tri := tris[i]
		n := triNormal(tri)
		ln := n.Length()
		if ln < 1e-12 {
			return
		}
		n = n.DivScalar(ln)
		if n.Z >= 0 {
			return // not downward-facing
		}
		c := tri[0].Add(tri[1]).Add(tri[2]).DivScalar(3)
		f := ohFace{
			deg:      math.Asin(-n.Z) * 180 / math.Pi,
			centroid: c,
			area:     ln / 2,
			downward: true,
		}
		f.failing = f.deg > thr(c.Z) && !ohSupported(s, c, bedZ, bedTol, depth, lateral)
		faces[i] = f
	})

	var rep OverhangReport
	for i, f := range faces {
		if !f.downward {
			continue
		}
		if f.deg > rep.WorstDeg {
			rep.WorstDeg, rep.WorstAt = f.deg, f.centroid
		}
		if f.failing {
			rep.Faces++
			rep.Area += f.area
			rep.Failing = append(rep.Failing, tris[i])
		}
	}
	return rep
}

// ohResolve fills in the documented defaults and returns the threshold
// function along with the three exemption distances.
func ohResolve(o OverhangOptions) (thr func(z float64) float64, bedTol, depth, lateral float64) {
	thr = o.Threshold
	if thr == nil {
		maxDeg := o.MaxDeg
		if maxDeg == 0 {
			maxDeg = ohDefaultMaxDeg
		}
		thr = func(float64) float64 { return maxDeg }
	}
	bedTol, depth, lateral = o.BedTol, o.SupportDepth, o.Lateral
	if bedTol == 0 {
		bedTol = ohDefaultBedTol
	}
	if depth == 0 {
		depth = ohDefaultSupportDepth
	}
	if lateral == 0 {
		lateral = ohDefaultLateral
	}
	return thr, bedTol, depth, lateral
}

// ohSupported reports whether the face centroid c is exempt: within bedTol
// of the plate, or with the plate or material (field < ohMaterial) depth
// below it — probed at the centroid and at four lateral offsets.
func ohSupported(s *solid.Solid, c v3.Vec, bedZ, bedTol, depth, lateral float64) bool {
	if bedTol > 0 && c.Z <= bedZ+bedTol {
		return true
	}
	if depth <= 0 {
		return false
	}
	pz := c.Z - depth
	if bedTol > 0 && pz <= bedZ+bedTol/2 {
		return true // the probe reaches the plate: the face rests on it
	}
	offsets := []v3.Vec{{}}
	if lateral > 0 {
		offsets = append(offsets,
			v3.X(lateral), v3.X(-lateral),
			v3.Y(lateral), v3.Y(-lateral))
	}
	for _, d := range offsets {
		if At(s, v3.XYZ(c.X+d.X, c.Y+d.Y, pz)) < ohMaterial {
			return true
		}
	}
	return false
}

// ohWorstFailing returns the worst angle among failing faces and its
// centroid, recomputed from the triangles so OverhangReport stays small.
func ohWorstFailing(failing []mesh.Triangle3) (float64, v3.Vec) {
	var worst float64
	var at v3.Vec
	for _, tri := range failing {
		n := triNormal(tri)
		ln := n.Length()
		if ln < 1e-12 {
			continue
		}
		deg := math.Asin(-n.Z/ln) * 180 / math.Pi
		if deg > worst {
			worst = deg
			at = tri[0].Add(tri[1]).Add(tri[2]).DivScalar(3)
		}
	}
	return worst, at
}

// RequireOverhangArea logs the audit and fails the test when the
// unsupported over-threshold area exceeds areaBudget (mm²). The budget is a
// parameter, not a magic number: marching cubes leaves thin artefact bands
// along sharp V-apex edges, so a real part usually gets a small non-zero
// allowance.
func RequireOverhangArea(t testing.TB, s *solid.Solid, cellsPerMM float64, o OverhangOptions, areaBudget float64) {
	t.Helper()
	r := Overhangs(s, cellsPerMM, o)
	t.Log(r.String())
	if r.Area > areaBudget {
		worst, at := ohWorstFailing(r.Failing)
		t.Errorf("unsupported overhang area %.3f mm² exceeds budget %.3f mm² (%d faces, worst failing %.1f° at (%.3f, %.3f, %.3f))",
			r.Area, areaBudget, r.Faces, worst, at.X, at.Y, at.Z)
	}
}
