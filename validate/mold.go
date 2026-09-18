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

// Moldability audit for one draw direction: the direction a mold half pulls
// away from the part. Nothing here depends on the solid's frame — the draw
// direction is an argument, so one part can be scored along many candidate
// directions without being rotated between calls.
//
// Every surface face gets a draft angle: the angle between the face and the
// draw, asin(n·d) for an outward unit normal n and a unit draw d.
//
//	+90°  the face looks straight down the draw and releases freely
//	  0°  the face is parallel to the draw: a wall that drags
//	 <0°  the face looks back against the draw: an undercut
//
// The normal test alone misses a face that points the right way but sits
// under a lip, so every forward-facing face is also ray-marched along the
// draw. If the ray re-enters material before leaving the bounding box the
// face is shadowed: it cannot reach the parting plane either. Back-facing
// and shadowed are both reported as undercuts, because both need a second
// mold part to release.
//
// The opening is exempt. A one-part open-face mold parts on the plane at
// the back of the draw, and the faces lying in that plane are the opening,
// not undercuts — without the exemption a flat-backed ornament reports its
// entire back face as the worst undercut in the part.

const (
	mdDefaultMinDraftDeg   = 1.0    // degrees from the draw direction
	mdDefaultStep          = 0.2    // mm, minimum ray-march step
	mdDefaultOpenFaceCells = 1.0    // mesh cells either side of the parting plane
	mdMaterial             = -0.005 // a hair inside the surface counts as material

	// mdUndercutTolDeg is the dead band around 0° draft. A wall parallel
	// to the draw is the most common surface in a moldable part, and
	// marching cubes renders it with normals that wobble a fraction of a
	// degree either side of zero. Without a dead band that noise splits
	// one wall arbitrarily between "drags" and "undercut", which are
	// different repairs. Inside the band a face is a wall: the draft
	// check judges it, and the undercut verdict is kept for faces that
	// genuinely look back against the draw.
	mdUndercutTolDeg = 0.5
)

// MoldOptions controls the moldability audit. Every zero field takes its
// documented default; pass a negative value to switch a check off.
type MoldOptions struct {
	// MinDraftDeg is the smallest draft angle that still releases cleanly
	// (default 1.0°; negative reports undercuts only). Faces at or above
	// it pass, faces below it but still forward-facing are reported as
	// under-drafted, which is a different repair from an undercut: add
	// draft, do not add a mold part.
	MinDraftDeg float64
	// Step is the minimum ray-march step in mm for the shadow test
	// (default 0.2). The march is sphere-traced off the field, so this
	// only bounds how small a step it takes through a tight gap.
	Step float64
	// OpenFaceTol is how close to the parting plane at the back of the
	// draw a back-facing face may sit and still count as the mold's
	// opening rather than an undercut. The default is one mesh cell
	// (1/cellsPerMM), because marching cubes renders a sharp rim as a
	// band of back-facing triangles a cell thick and can only locate the
	// parting plane that precisely. The cost of the exemption is its
	// limit: an undercut shallower than a cell at the parting plane is
	// invisible, so audit at a density that resolves the features you
	// care about. Pass a negative value to treat the opening as an
	// undercut too, which is what a closed two-part mold wants.
	OpenFaceTol float64
	// NoShadow classifies on face normals alone, skipping the ray march.
	// Much faster, and exact for a convex part.
	NoShadow bool
}

// MoldReport is the result of a moldability audit along one draw direction.
type MoldReport struct {
	Draw          v3.Vec           // the draw direction audited, normalised
	UndercutArea  float64          // mm² that cannot release along Draw
	UndercutFaces int              // number of such faces
	LowDraftArea  float64          // mm² that releases, but below MinDraftDeg
	LowDraftFaces int              // number of such faces
	WorstDraftDeg float64          // smallest draft angle on any non-exempt face (90 = nothing seen)
	WorstAt       v3.Vec           // centroid of that face
	Undercut      []mesh.Triangle3 // the undercut faces, for writing to their own STL
	LowDraft      []mesh.Triangle3 // the under-drafted faces, likewise
}

// Moldable reports whether the part pulls out of a one-part mold along
// Draw: no undercuts and nothing below the minimum draft.
func (r MoldReport) Moldable() bool {
	return r.UndercutFaces == 0 && r.LowDraftFaces == 0
}

// String renders the audit as one log line: the draw direction, the two
// failing areas, and the worst draft angle with where it is.
func (r MoldReport) String() string {
	head := fmt.Sprintf("moldability: draw (%.3f, %.3f, %.3f)", r.Draw.X, r.Draw.Y, r.Draw.Z)
	worst := fmt.Sprintf("minimum draft %.1f° at (%.3f, %.3f, %.3f)",
		r.WorstDraftDeg, r.WorstAt.X, r.WorstAt.Y, r.WorstAt.Z)
	if r.Moldable() {
		return fmt.Sprintf("%s; pulls in one part, no undercuts; %s", head, worst)
	}
	return fmt.Sprintf("%s; %d undercut faces %.3f mm², %d under-drafted faces %.3f mm²; %s",
		head, r.UndercutFaces, r.UndercutArea, r.LowDraftFaces, r.LowDraftArea, worst)
}

// Moldability audits s along draw at mesh density cellsPerMM. See
// MoldOptions for the draft threshold, the shadow test and the open-face
// exemption.
func Moldability(s *solid.Solid, draw v3.Vec, cellsPerMM float64, o MoldOptions) MoldReport {
	tris := mesh.CollectTriangles(s, render.NewMarchingCubesOctreeParallel(solid.CellsFor(s, cellsPerMM)))
	return mdAudit(s, tris, draw, cellsPerMM, o)
}

// mdAudit is Moldability on an already-rendered mesh, so a caller that
// already has the triangles pays for marching cubes once.
func mdAudit(s *solid.Solid, tris []mesh.Triangle3, draw v3.Vec, cellsPerMM float64, o MoldOptions) MoldReport {
	minDraft, step, openTol := mdResolve(o, cellsPerMM)
	d := draw.Normalize()
	box := s.Bounds().Box
	openAt := mdCenterAlong(box, d) - mdExtentAlong(box, d)/2
	maxDist := mdExtentAlong(box, d) + 2*step

	type mdFace struct {
		deg      float64
		centroid v3.Vec
		area     float64
		valid    bool
		exempt   bool
		undercut bool
		lowDraft bool
	}
	faces := make([]mdFace, len(tris))
	parallel(len(tris), func(i int) {
		tri := tris[i]
		n := triNormal(tri)
		ln := n.Length()
		if ln < 1e-12 {
			return
		}
		n = n.DivScalar(ln)
		c := tri[0].Add(tri[1]).Add(tri[2]).DivScalar(3)
		f := mdFace{
			deg:      math.Asin(mdClamp1(n.Dot(d))) * 180 / math.Pi,
			centroid: c,
			area:     ln / 2,
			valid:    true,
		}
		switch {
		case f.deg < -mdUndercutTolDeg && openTol > 0 && c.Dot(d) <= openAt+openTol:
			f.exempt = true // the opening of a one-part mold
		case f.deg < -mdUndercutTolDeg:
			f.undercut = true
		case !o.NoShadow && mdShadowed(s, c, n, d, maxDist, step):
			f.undercut = true
		case minDraft > 0 && f.deg < minDraft:
			f.lowDraft = true
		}
		faces[i] = f
	})

	rep := MoldReport{Draw: d, WorstDraftDeg: 90}
	for i, f := range faces {
		if !f.valid || f.exempt {
			continue
		}
		if f.deg < rep.WorstDraftDeg {
			rep.WorstDraftDeg, rep.WorstAt = f.deg, f.centroid
		}
		switch {
		case f.undercut:
			rep.UndercutFaces++
			rep.UndercutArea += f.area
			rep.Undercut = append(rep.Undercut, tris[i])
		case f.lowDraft:
			rep.LowDraftFaces++
			rep.LowDraftArea += f.area
			rep.LowDraft = append(rep.LowDraft, tris[i])
		}
	}
	return rep
}

// mdResolve fills in the documented defaults. The open-face tolerance
// scales with the mesh: the parting-plane band marching cubes leaves at a
// sharp rim is about one cell thick, so a fixed tolerance either misses it
// on a coarse mesh or swallows real geometry on a fine one.
func mdResolve(o MoldOptions, cellsPerMM float64) (minDraft, step, openTol float64) {
	minDraft, step, openTol = o.MinDraftDeg, o.Step, o.OpenFaceTol
	if minDraft == 0 {
		minDraft = mdDefaultMinDraftDeg
	}
	if step <= 0 {
		step = mdDefaultStep
	}
	if openTol == 0 {
		openTol = mdDefaultOpenFaceCells / cellsPerMM
	}
	return minDraft, step, openTol
}

// mdShadowed reports whether a ray leaving the face at c along the draw
// re-enters material before it has travelled maxDist. The start point is
// lifted off the surface along the normal so the face itself is not a hit,
// and the march is sphere-traced: the field is a lower bound on the
// distance to the surface, so stepping by it never jumps a wall.
func mdShadowed(s *solid.Solid, c, n, d v3.Vec, maxDist, minStep float64) bool {
	from := c.Add(n.MulScalar(minStep))
	for t := 0.0; t < maxDist; {
		v := At(s, from.Add(d.MulScalar(t)))
		if v < mdMaterial {
			return true
		}
		if v < minStep {
			v = minStep
		}
		t += v
	}
	return false
}

// mdCenterAlong returns the box centre projected onto d.
func mdCenterAlong(b v3.Box, d v3.Vec) float64 { return b.Center().Dot(d) }

// mdExtentAlong returns how far the box reaches along d.
func mdExtentAlong(b v3.Box, d v3.Vec) float64 {
	h := b.Size().MulScalar(0.5)
	return 2 * (math.Abs(d.X)*h.X + math.Abs(d.Y)*h.Y + math.Abs(d.Z)*h.Z)
}

// mdClamp1 keeps a dot product of two unit vectors inside asin's domain.
func mdClamp1(x float64) float64 {
	if x > 1 {
		return 1
	}
	if x < -1 {
		return -1
	}
	return x
}

// RequireMoldable logs the audit and fails the test when the failing area
// (undercut plus under-drafted) exceeds areaBudget in mm². The budget is a
// parameter for the same reason RequireOverhangArea takes one: marching
// cubes leaves thin artefact bands along sharp edges, so a real part
// usually needs a small non-zero allowance.
func RequireMoldable(t testing.TB, s *solid.Solid, draw v3.Vec, cellsPerMM float64, o MoldOptions, areaBudget float64) {
	t.Helper()
	rep := Moldability(s, draw, cellsPerMM, o)
	t.Log(rep.String())
	if bad := rep.UndercutArea + rep.LowDraftArea; bad > areaBudget {
		t.Fatalf("not moldable along (%.3f, %.3f, %.3f): %.3f mm² fails, budget %.3f mm²: %s",
			rep.Draw.X, rep.Draw.Y, rep.Draw.Z, bad, areaBudget, rep.String())
	}
}
