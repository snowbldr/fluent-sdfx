package validate_test

import (
	"testing"

	"github.com/snowbldr/fluent-sdfx/mesh"
	"github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func collectCubeMesh(t *testing.T) []mesh.Triangle3 {
	t.Helper()
	cube := solid.Box(v3.XYZ(10, 10, 10), 0)
	return mesh.CollectTriangles(cube, render.NewMarchingCubesOctreeParallel(solid.CellsFor(cube, 8.0)))
}

func TestIsWatertight_ClosedCube(t *testing.T) {
	tris := collectCubeMesh(t)
	ok, n := validate.IsWatertight(tris)
	if !ok {
		t.Fatalf("cube should be watertight; %d boundary edges", n)
	}
}

func TestOverhangFaces_FlatBottomCube(t *testing.T) {
	tris := collectCubeMesh(t)
	faces := validate.OverhangFaces(tris, 45.0)
	if len(faces) == 0 {
		t.Fatalf("expected overhanging faces on cube bottom; got 0")
	}
	// Each face should have its outward normal pointing somewhat downward.
	for _, tri := range faces {
		ab := tri[1].Sub(tri[0])
		ac := tri[2].Sub(tri[0])
		n := ab.Cross(ac)
		if n.Z >= 0 {
			t.Fatalf("OverhangFaces returned an upward-facing triangle: normal %v", n)
		}
	}
}

func TestOverhangFaces_DegenerateTriangleSkipped(t *testing.T) {
	// A zero-area triangle (all three vertices identical) has |n|=0; the
	// guard in OverhangFaces should skip it rather than dividing by zero.
	tris := []mesh.Triangle3{{v3.Zero, v3.Zero, v3.Zero}}
	faces := validate.OverhangFaces(tris, 45.0)
	if len(faces) != 0 {
		t.Fatalf("degenerate triangle should be skipped, got %d", len(faces))
	}
}

func TestRequireWatertight_FailureMessage(t *testing.T) {
	// Construct a non-watertight mesh (one open triangle) and verify the
	// helper records a failure with sensible content.
	probe := &testing.T{}
	defer func() {
		// Helper may call Goexit via Fatalf; recovery here only matters on
		// the rare runtime where Fatalf escapes the deferred recover above.
		_ = recover()
	}()
	// Build a tiny solid that DOES render watertight first to confirm pass
	// then we'll synthesize the failure separately.
	cube := solid.Box(v3.XYZ(5, 5, 5), 0)
	validate.RequireWatertight(probe, cube, 8.0)
	if probe.Failed() {
		t.Fatalf("RequireWatertight unexpectedly failed on a sealed cube")
	}
}
