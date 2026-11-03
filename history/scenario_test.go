package history

import (
	"testing"

	"github.com/snowbldr/fluent-sdfx/sculpt"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// The sculpt scenario forms as histories give the same geometry as the
// Go that sculpt's own tests build.

func TestOrnamentMatchesGo(t *testing.T) {
	r := eval(t, load(t, "testdata/ornament.yaml"))
	ball := solid.Sphere(15).Shell(1.5)
	loop := solid.Torus(3, 0.8).RotateX(90).TranslateZ(18.3)
	bauble := ball.Union(loop)
	seed := r.Value("VoidSeed")
	t.Logf("void seed %+v", seed)
	want := sculpt.Sample(bauble, 96).Fill(v3.Zero).Solid()
	congruent(t, "filled ornament", r.Solid(), want, 1e-9)
	validate.RequireNear(t, "solid at the centre", validate.At(r.Solid(), v3.Zero), -15.75, 0.7)
}

func TestClusterMatchesGo(t *testing.T) {
	r := eval(t, load(t, "testdata/cluster.yaml"))
	body := solid.Sphere(10)
	arm := solid.Sphere(4).TranslateX(10)
	block := solid.Box(v3.XYZ(10, 10, 10), 0).TranslateZ(10)
	bead := solid.Sphere(2).TranslateZ(20)
	want := sculpt.Sample(body.Union(arm, block, bead), 128).Keep(v3.Zero).Smooth(2).Solid()
	congruent(t, "smoothed cluster", r.Solid(), want, 1e-9)
}

func TestKitchenMatchesGo(t *testing.T) {
	r := eval(t, load(t, "testdata/kitchen.yaml"))
	plate := solid.Cylinder(3, 20, 0.5)
	peg := solid.Cylinder(8, 1.5, 0).TranslateX(14).RotateCopyZ(6)
	stem := solid.Capsule(10, 2)
	mast := stem.Bottom().On(plate.Top()).Union()
	assembly := mast.Union(peg)
	want := solid.SmoothUnion(solid.RoundMin(1), assembly, solid.Torus(17, 1))
	if len(r.Bodies) != 1 || r.Bodies[0].Name != "rounded" {
		t.Errorf("bodies %+v, want rounded", r.Bodies)
	}
	congruent(t, "kitchen", r.Solid(), want, 1e-12)
}
