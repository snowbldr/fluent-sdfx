package shape_test

import (
	"math"
	"testing"

	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Hole is a diameter, Circle is a radius. The whole point of the helper is
// that the two disagree, so prove it at the surface.
func TestHoleIsADiameter(t *testing.T) {
	h := shape.Hole(6).Extrude(2)
	if got := validate.At(h, v3.XYZ(2.9, 0, 0)); got >= 0 {
		t.Fatalf("Hole(6) at r=2.9: want inside, got sdf %.4f", got)
	}
	if got := validate.At(h, v3.XYZ(3.1, 0, 0)); got <= 0 {
		t.Fatalf("Hole(6) at r=3.1: want outside, got sdf %.4f", got)
	}
}

// A through-cut must leave the bore open end to end. Opening measures the
// clear span along the axis: it is the thing an exact-height tool gets
// wrong, and no render shows.
func TestThroughLeavesNoSkin(t *testing.T) {
	plate := solid.Box(v3.XYZ(20, 20, 5), 0) // z = -2.5 .. 2.5
	part := plate.Cut(shape.Hole(4).Through(plate, v3.Z(1), 1))

	validate.RequireProbes(t, part, []validate.Probe{
		validate.Out("bore-open-at-top-face", v3.XYZ(0, 0, 2.49)),
		validate.Out("bore-open-at-mid", v3.XYZ(0, 0, 0)),
		validate.Out("bore-open-at-bottom-face", v3.XYZ(0, 0, -2.49)),
		validate.In("plate-kept-beside-bore", v3.XYZ(6, 0, 0)),
	})

	// A sealed bore and an open one differ by exactly the bore, so measure
	// the volume rather than trusting the probes to have landed on a skin.
	validate.RequireVolumeNear(t, part, 10, 20*20*5-math.Pi*2*2*5, 0.02)
}

// The tool is sized from the target, so a thickness change cannot leave the
// tool behind. Build the same cut against two plates and check both open.
func TestThroughTracksTargetThickness(t *testing.T) {
	for _, thick := range []float64{2, 5, 30} {
		plate := solid.Box(v3.XYZ(20, 20, thick), 0)
		part := plate.Cut(shape.Hole(4).Through(plate, v3.Z(1), 1))
		validate.RequireProbes(t, part, []validate.Probe{
			validate.Out("bore-open-at-mid", v3.XYZ(0, 0, 0)),
			validate.In("plate-kept-beside-bore", v3.XYZ(6, 0, 0)),
		})
		validate.RequireVolumeNear(t, part, 10, 20*20*thick-math.Pi*2*2*thick, 0.02)
	}
}

// A profile positioned before the call keeps its position for a ±Z draw.
func TestThroughKeepsProfilePositionForZ(t *testing.T) {
	plate := solid.Box(v3.XYZ(40, 20, 5), 0)
	part := plate.Cut(shape.Hole(4).TranslateX(12).Through(plate, v3.Z(1), 1))

	validate.RequireProbes(t, part, []validate.Probe{
		validate.Out("bore-at-x12", v3.XYZ(12, 0, 0)),
		validate.In("material-at-origin", v3.XYZ(0, 0, 0)),
	})
}

// -Z must give the same tool as +Z: the extrusion is symmetric, and the
// centring must not land it on the wrong side.
func TestThroughNegativeZMatchesPositive(t *testing.T) {
	plate := solid.Box(v3.XYZ(20, 20, 5), 0).TranslateZ(17) // well off the origin
	up := plate.Cut(shape.Hole(4).Through(plate, v3.Z(1), 1))
	down := plate.Cut(shape.Hole(4).Through(plate, v3.Z(-1), 1))

	pts := validate.GridPoints(plate.Bounds().Box, 1)
	validate.RequireIdentical(t, up, down, pts, 1e-9)
}

// An off-axis draw bores along that axis instead.
func TestThroughAlongX(t *testing.T) {
	block := solid.Box(v3.XYZ(20, 10, 10), 0)
	part := block.Cut(shape.Hole(4).Through(block, v3.X(1), 1))

	validate.RequireProbes(t, part, []validate.Probe{
		validate.Out("bore-at-centre", v3.XYZ(0, 0, 0)),
		validate.Out("bore-near-face", v3.XYZ(9.9, 0, 0)),
		validate.In("material-above-bore", v3.XYZ(0, 0, 4)),
	})

	validate.RequireVolumeNear(t, part, 10, 20*10*10-math.Pi*2*2*20, 0.02)
}

// The pad is real clearance, not a rounding allowance: the tool must stand
// proud of the target at both ends.
func TestThroughPadStandsProud(t *testing.T) {
	plate := solid.Box(v3.XYZ(20, 20, 5), 0)
	tool := shape.Hole(4).Through(plate, v3.Z(1), 1)
	b := tool.Bounds().Box
	if got := b.Max.Z; math.Abs(got-3.5) > 1e-9 {
		t.Fatalf("tool top at %.4f, want 3.5 (2.5 plate + 1 pad)", got)
	}
	if got := b.Min.Z; math.Abs(got+3.5) > 1e-9 {
		t.Fatalf("tool bottom at %.4f, want -3.5", got)
	}
}
