package validate

import (
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const mdCells = 5.0

// A cone releases along its own taper and jams against it. The second half
// is the real assertion: the normal test has to change its answer when the
// draw flips, or it is not measuring anything.
func TestMoldabilityConeDrawsOneWay(t *testing.T) {
	cone := solid.Cone(20, 10, 2, 0) // wide at the bottom, narrow at the top

	up := Moldability(cone, v3.Z(1), mdCells, MoldOptions{})
	t.Log(up.String())
	if !up.Moldable() {
		t.Fatalf("cone should pull along +Z: %s", up)
	}

	down := Moldability(cone, v3.Z(-1), mdCells, MoldOptions{})
	t.Log(down.String())
	if down.Moldable() {
		t.Fatalf("a cone that widens downward cannot pull along -Z: %s", down)
	}
	if down.WorstDraftDeg >= 0 {
		t.Fatalf("worst draft pulling -Z is %.1f°, want negative", down.WorstDraftDeg)
	}
}

// The flat-back ornament: a slab on z=0 with a raised boss. This is the
// M1.5 fast path, and it only passes because the back face is read as the
// mold's opening rather than a 90° undercut.
func TestMoldabilityFlatBackIsTheOpening(t *testing.T) {
	slab := solid.Box(v3.XYZ(30, 30, 3), 0).ZeroZ()
	boss := solid.Cylinder(2, 6, 0).ZeroZ().TranslateZ(3)
	part := slab.Union(boss)

	open := Moldability(part, v3.Z(1), mdCells, MoldOptions{})
	t.Log(open.String())
	if open.UndercutFaces != 0 {
		t.Fatalf("flat back should be the opening, not an undercut: %s", open)
	}

	// Switch the exemption off and the same part is dominated by its back
	// face: that is what a closed two-part mold sees.
	closed := Moldability(part, v3.Z(1), mdCells, MoldOptions{OpenFaceTol: -1})
	t.Log(closed.String())
	if closed.UndercutArea < 30*30*0.9 {
		t.Fatalf("with the opening disabled the 900 mm² back should be undercut, got %.1f mm²", closed.UndercutArea)
	}
}

// A mushroom is the textbook undercut: the cap's underside faces back
// against every upward draw, and it is nowhere near the parting plane.
func TestMoldabilityFindsUndercutAndSaysWhere(t *testing.T) {
	stem := solid.Cylinder(10, 3, 0).ZeroZ()
	cap := solid.Cylinder(3, 8, 0).ZeroZ().TranslateZ(10)
	part := stem.Union(cap)

	rep := Moldability(part, v3.Z(1), mdCells, MoldOptions{})
	t.Log(rep.String())
	if rep.UndercutFaces == 0 {
		t.Fatal("the underside of the cap is an undercut and was not found")
	}

	// Two surfaces cannot release, not one. The annulus under the cap is
	// back-facing, π(8² - 3²) ≈ 172.8 mm²; the stem wall passes the normal
	// test and is still trapped, because a ray leaving it runs straight
	// into the overhanging cap — 2π·3·10 ≈ 188.5 mm² the shadow march
	// finds and the normal test alone would miss.
	want := math.Pi*(8*8-3*3) + 2*math.Pi*3*10
	if rep.UndercutArea < want*0.8 || rep.UndercutArea > want*1.25 {
		t.Fatalf("undercut area %.1f mm², want about %.1f (cap underside plus shadowed stem)", rep.UndercutArea, want)
	}
	if math.Abs(rep.WorstAt.Z-10) > 1.0 {
		t.Fatalf("worst face at z=%.2f, want the cap underside at z=10", rep.WorstAt.Z)
	}
}

// A face can point straight down the draw and still not release, because
// something else is in the way. Only the ray march sees it, so run the same
// part both ways: this is the test that fails if the shadow test is a no-op.
func TestMoldabilityShadowCatchesWhatNormalsMiss(t *testing.T) {
	base := solid.Box(v3.XYZ(20, 20, 3), 0).ZeroZ()               // z 0..3, top face points +Z
	roof := solid.Box(v3.XYZ(20, 20, 3), 0).ZeroZ().TranslateZ(9) // z 9..12, directly above
	wall := solid.Box(v3.XYZ(3, 20, 6), 0).ZeroZ().TranslateZ(3).TranslateX(8.5)
	part := base.Union(roof, wall)

	blind := Moldability(part, v3.Z(1), mdCells, MoldOptions{NoShadow: true})
	t.Log("normals only: " + blind.String())
	seeing := Moldability(part, v3.Z(1), mdCells, MoldOptions{})
	t.Log("with shadow:  " + seeing.String())

	if seeing.UndercutArea <= blind.UndercutArea {
		t.Fatalf("the shadow test found nothing extra: %.1f mm² with, %.1f mm² without",
			seeing.UndercutArea, blind.UndercutArea)
	}
	// The roofed-over part of the base top is about 17 x 20 mm.
	if seeing.UndercutArea-blind.UndercutArea < 200 {
		t.Fatalf("shadowed area %.1f mm², want roughly the 340 mm² under the roof",
			seeing.UndercutArea-blind.UndercutArea)
	}
}

// Draft is a separate verdict from undercut: a vertical wall releases in
// principle and drags in practice, and the repair is to taper it, not to
// split the mold.
func TestMoldabilityDraftIsItsOwnVerdict(t *testing.T) {
	straight := solid.Cylinder(10, 8, 0)
	rep := Moldability(straight, v3.Z(1), mdCells, MoldOptions{MinDraftDeg: 1})
	t.Log(rep.String())
	if rep.UndercutFaces != 0 {
		t.Fatalf("a plain cylinder has no undercut along its axis: %s", rep)
	}
	if rep.LowDraftFaces == 0 {
		t.Fatal("a vertical wall is 0° draft and should be reported as dragging")
	}
	if math.Abs(rep.WorstDraftDeg) > 2 {
		t.Fatalf("worst draft %.2f°, want about 0° for a vertical wall", rep.WorstDraftDeg)
	}

	// Give the same part 6° of taper and the drag goes away.
	tapered := solid.Cone(10, 8, 8-10*math.Tan(6*math.Pi/180), 0)
	ok := Moldability(tapered, v3.Z(1), mdCells, MoldOptions{MinDraftDeg: 1})
	t.Log(ok.String())
	if ok.LowDraftArea > 5 {
		t.Fatalf("a 6° taper should clear a 1° minimum, got %.1f mm² dragging: %s", ok.LowDraftArea, ok)
	}
}

// Disabling the draft check leaves undercuts still reported.
func TestMoldabilityDraftCheckCanBeDisabled(t *testing.T) {
	cyl := solid.Cylinder(10, 8, 0)
	rep := Moldability(cyl, v3.Z(1), mdCells, MoldOptions{MinDraftDeg: -1})
	t.Log(rep.String())
	if rep.LowDraftFaces != 0 {
		t.Fatalf("MinDraftDeg -1 should switch the drag check off: %s", rep)
	}
}
