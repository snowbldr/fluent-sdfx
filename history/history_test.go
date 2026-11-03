package history

import (
	"math"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/snowbldr/fluent-sdfx/anchor"
	"github.com/snowbldr/fluent-sdfx/sculpt"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func load(t *testing.T, path string) *History {
	t.Helper()
	h, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func eval(t *testing.T, h *History) *Result {
	t.Helper()
	r := NewRuntime(Options{}).Eval(h)
	if err := r.Err(); err != nil {
		t.Fatal(err)
	}
	return r
}

// congruent asserts two solids have the same field at random points in
// and around them, to within tol.
func congruent(t *testing.T, what string, got, want *solid.Solid, tol float64) {
	t.Helper()
	bb := want.Bounds().Box.Enlarge(v3.XYZ(3, 3, 3))
	rng := rand.New(rand.NewSource(1))
	worst := 0.0
	for i := 0; i < 3000; i++ {
		s := bb.Size()
		p := bb.Min.Add(v3.XYZ(rng.Float64()*s.X, rng.Float64()*s.Y, rng.Float64()*s.Z))
		worst = math.Max(worst, math.Abs(validate.At(got, p)-validate.At(want, p)))
	}
	if worst > tol {
		t.Errorf("%s: fields differ by up to %.3g (tolerance %.3g)", what, worst, tol)
	} else {
		t.Logf("%s: fields agree to %.2g at 3000 points", what, worst)
	}
}

func mirrorX(v v3.Vec) v3.Vec { return v3.XYZ(-v.X, v.Y, v.Z) }

// figureInGo is testdata/figure.yaml written as Go, by hand.
func figureInGo(headR float64) (hollow *solid.Solid, noseTip anchor.Point) {
	body := solid.Sphere(12)
	head := solid.Sphere(headR).Translate(v3.XYZ(0, 0, 14))
	figure := solid.SmoothUnion(solid.RoundMin(3), body, head)
	noseTip = anchor.Ray(figure, v3.XYZ(0, -40, 14), v3.XYZ(0, 1, 0))
	nose := sculpt.Pull(figure, noseTip.Pos, v3.XYZ(0, -4, 0), 7)
	cheek := anchor.Toward(nose, v3.XYZ(-30, -14, 12), v3.XYZ(0, 0, 12))
	cheeks := sculpt.Pull(nose, cheek.Pos, v3.XYZ(-1.5, -1, 0), 4)
	cheeks = sculpt.Pull(cheeks, mirrorX(cheek.Pos), v3.XYZ(1.5, -1, 0), 4)
	return sculpt.Sample(cheeks, 96).Redistance().Solid().Shell(1.6), noseTip
}

func TestFigureMatchesGo(t *testing.T) {
	r := eval(t, load(t, "testdata/figure.yaml"))
	if len(r.Bodies) != 1 || r.Bodies[0].Name != "hollow" {
		t.Fatalf("bodies %+v, want just hollow", r.Bodies)
	}
	want, tip := figureInGo(8)
	congruent(t, "figure", r.Solid(), want, 1e-9)
	got := r.Value("NoseTip").(anchor.Point)
	if !got.Pos.Equals(tip.Pos, 1e-12) {
		t.Errorf("NoseTip %v, want %v", got.Pos, tip.Pos)
	}
	for _, s := range r.Steps {
		t.Logf("%-8s %-16s %-6s %6.1f ms", s.label(), s.Op, s.Status, s.Millis)
	}
}

// Phase 2.5 in the history: change the head radius and the nose still
// comes out of the face, because it is pulled at an anchor and not at the
// coordinates the face had.
func TestFigureReflowsWithTheHead(t *testing.T) {
	h := load(t, "testdata/figure.yaml")
	for _, r := range []float64{8, 10} {
		h.Params["headR"] = r
		res := eval(t, h)
		tip := res.Value("NoseTip").(anchor.Point)
		validate.RequireNear(t, "face found on the head", tip.Pos.Y, -r, 1e-4)
		// The nose's own surface, before the shell: the landing point is on it.
		nose := res.SolidOf("nose")
		validate.RequireNear(t, "nose lands 4 mm out", validate.At(nose, tip.Pos.Add(v3.Y(-4))), 0, 1e-6)
		want, _ := figureInGo(r)
		congruent(t, "reflowed figure", res.Solid(), want, 1e-9)
	}
}

func TestDeleteTheNoseReflows(t *testing.T) {
	h := load(t, "testdata/figure.yaml")
	if err := h.Remove(h.Index("nose")); err != nil {
		t.Fatal(err)
	}
	r := eval(t, h)
	// The cheeks now pull the figure directly; the face is where it was.
	cheeks := r.SolidOf("cheeks")
	if validate.Inside(cheeks, v3.XYZ(0, -11, 14)) {
		t.Error("the nose is still there")
	}
	validate.RequireNear(t, "face back at the head", validate.At(cheeks, v3.XYZ(0, -8, 14)), 0, 1e-6)
}

func TestDeleteTheBodyBlocksAndTheHeadStillShows(t *testing.T) {
	h := load(t, "testdata/figure.yaml")
	if err := h.Remove(h.Index("body")); err != nil {
		t.Fatal(err)
	}
	r := NewRuntime(Options{}).Eval(h)
	for _, s := range r.Steps {
		t.Logf("%-8s %-8s %s", s.label(), s.Status, s.Error)
	}
	fig := r.Steps[h.Index("figure")]
	if fig.Status != Blocked || !strings.Contains(fig.Error, "body") || !strings.Contains(fig.Error, "does not exist") {
		t.Errorf("figure: %s %q, want blocked naming the missing body", fig.Status, fig.Error)
	}
	for _, s := range r.Steps[h.Index("figure"):] {
		if s.Status != Blocked {
			t.Errorf("step %s is %s, want blocked", s.label(), s.Status)
		}
	}
	if len(r.Bodies) != 1 || r.Bodies[0].Name != "head" {
		t.Fatalf("bodies %+v, want the head alone", r.Bodies)
	}
	validate.RequireNear(t, "the head renders", validate.At(r.Solid(), v3.Z(22)), 0, 1e-9)
}

func TestTurningAStepOffPassesItsInput(t *testing.T) {
	h := load(t, "testdata/figure.yaml")
	h.SetOff(h.Index("nose"), true)
	off := eval(t, h)
	h.Remove(h.Index("nose"))
	gone := eval(t, h)
	congruent(t, "off is as if deleted", off.Solid(), gone.Solid(), 1e-12)
	if s := off.Steps[4]; s.Status != Off {
		t.Errorf("nose status %s", s.Status)
	}
}

func TestParamChangesEveryUse(t *testing.T) {
	h := load(t, "testdata/figure.yaml")
	for _, wall := range []float64{1.6, 3} {
		h.Params["wall"] = wall
		r := eval(t, h)
		// Through the back of the body, where nothing was pulled: Shell
		// straddles the surface, so the inner wall is half the wall in.
		d, _ := validate.SurfaceAlong(r.Solid(), v3.Zero, v3.Y(1), 20, 0.01)
		validate.RequireNear(t, "inner wall", d, 12-wall/2, 0.1)
	}
}

func TestProvenanceRoundTripsAndIsNotHashed(t *testing.T) {
	src, err := os.ReadFile("testdata/figure.yaml")
	if err != nil {
		t.Fatal(err)
	}
	h, err := Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	// Edit another step, save, load: provenance comes back exactly.
	cheeks := h.Steps[h.Index("cheeks")]
	cheeks.Args.(map[string]any)["radius"] = 5
	h.Set(h.Index("cheeks"), cheeks)
	b, err := h.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	back, err := Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	nose := back.Steps[back.Index("nose")]
	if nose.By != "agent" || nose.From != "chat/123" {
		t.Errorf("provenance after a round trip: by %q from %q", nose.By, nose.From)
	}
	if !strings.HasPrefix(string(b), "# The figure from sculpt's scenario tests") {
		t.Errorf("header lost:\n%s", b)
	}

	// Relabel a step's provenance: nothing recomputes.
	rt := NewRuntime(Options{})
	if err := rt.Eval(back).Err(); err != nil {
		t.Fatal(err)
	}
	nose.From, nose.By = "chat/456", "user"
	back.Set(back.Index("nose"), nose)
	r := rt.Eval(back)
	for _, s := range r.Steps {
		if !s.Cached {
			t.Errorf("step %s recomputed after a provenance change", s.label())
		}
	}
	if s := r.Steps[back.Index("nose")]; s.From != "chat/456" || s.By != "user" {
		t.Errorf("result carries by %q from %q", s.By, s.From)
	}
}

func TestEditRecomputesOnlyWhatFollows(t *testing.T) {
	h := load(t, "testdata/figure.yaml")
	rt := NewRuntime(Options{})
	rt.Eval(h)
	cheeks := h.Steps[h.Index("cheeks")]
	cheeks.Args.(map[string]any)["radius"] = 5
	h.Set(h.Index("cheeks"), cheeks)
	r := rt.Eval(h)
	for _, s := range r.Steps {
		want := s.Index < h.Index("cheeks")
		if s.Cached != want {
			t.Errorf("step %s cached=%v, want %v", s.label(), s.Cached, want)
		}
	}
}

func TestMissedAnchorIsAStepError(t *testing.T) {
	h := load(t, "testdata/figure.yaml")
	h.Steps[h.Index("NoseTip")].Anchor.Ray.Dir = []any{0, -1, 0} // away from the face
	r := NewRuntime(Options{}).Eval(h)
	s := r.Steps[h.Index("NoseTip")]
	if s.Status != Failed || !strings.Contains(s.Error, "misses") {
		t.Errorf("NoseTip: %s %q", s.Status, s.Error)
	}
	if n := r.Steps[h.Index("nose")]; n.Status != Blocked || !strings.Contains(n.Error, "NoseTip") {
		t.Errorf("nose: %s %q, want blocked by NoseTip", n.Status, n.Error)
	}
}

func TestPullTooFarSaysHowToFix(t *testing.T) {
	h := load(t, "testdata/figure.yaml")
	h.Steps[h.Index("nose")].Args.(map[string]any)["by"] = []any{0, -9, 0}
	r := NewRuntime(Options{}).Eval(h)
	s := r.Steps[h.Index("nose")]
	t.Logf("nose: %s: %s", s.Status, s.Error)
	if s.Status != Failed || !strings.Contains(s.Error, "keep the pull under") {
		t.Errorf("nose: %s %q", s.Status, s.Error)
	}
	// The steps before still evaluated, and the figure is what shows.
	if len(r.Bodies) != 1 || r.Bodies[0].Name != "figure" {
		t.Errorf("bodies %+v, want figure", r.Bodies)
	}
}
