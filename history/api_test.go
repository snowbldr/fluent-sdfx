package history

import (
	"strings"
	"testing"

	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// An agent sends JSON; it is YAML, and it parses the same.
func TestJSONIsAHistory(t *testing.T) {
	h, err := Parse([]byte(`{"steps": [
		{"name": "ball", "op": "solid.Sphere", "args": {"radius": 10}, "by": "agent", "from": "chat/1"},
		{"op": "sculpt.Pull", "args": {"at": [0, -10, 0], "by": [0, -3, 0], "radius": 6}}
	]}`))
	if err != nil {
		t.Fatal(err)
	}
	r := eval(t, h)
	validate.RequireNear(t, "pulled", validate.At(r.Solid(), v3.Y(-13)), 0, 1e-6)
}

// The messages an agent gets back name the problem and the fix.
func TestErrorsSayWhatToDo(t *testing.T) {
	for _, c := range []struct{ yaml, want string }{
		{"steps: [{op: solid.Sphere, args: {radious: 3}}]", "has no parameter radious; did you mean radius?"},
		{"steps: [{op: solid.Sphere, args: {}}]", "needs radius (float64)"},
		{"steps: [{op: solid.Spher, args: [3]}]", "no op solid.Spher; did you mean solid.Sphere?"},
		{"steps: [{op: solid.Sphere, args: [3]}, {op: Transalte, args: [[1, 0, 0]]}]", "has no method Transalte; did you mean Translate?"},
		{"steps: [{op: solid.Sphere, args: [3]}, {op: Translate, args: [[1, 0]]}]", "a v3.Vec is 3 numbers, got 2"},
		{"steps: [{op: solid.Box, args: {size: 3, round: 0}}]", "want a v3.Vec, got the number 3"},
		{"steps: [{op: Translate, args: [[1, 0, 0]]}]", "needs an input and there is no step before it"},
		{"steps: [{op: solid.Sphere, args: [$r]}]", "no param r"},
	} {
		h, err := Parse([]byte(c.yaml))
		if err != nil {
			t.Fatal(err)
		}
		err = NewRuntime(Options{}).Eval(h).Err()
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s:\n  got  %v\n  want %q", c.yaml, err, c.want)
		}
	}
	for _, c := range []struct{ yaml, want string }{
		{"steps: [{op: solid.Sphere, radius: 3}]", "field radius not found"},
		{"steps: [{name: a, op: solid.Sphere, args: [1]}, {name: a, op: solid.Sphere, args: [2]}]", "step 0 already has this name"},
		{"steps: [{op: solid.Sphere, args: [1], by: robot}]", "by is user or agent"},
		{"steps: [{op: solid.Sphere, args: [1], mirror: w}]", "mirror is x, y or z"},
		{"steps: [{name: 2a, op: solid.Sphere, args: [1]}]", "a name is letters"},
	} {
		_, err := Parse([]byte(c.yaml))
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s:\n  got  %v\n  want %q", c.yaml, err, c.want)
		}
	}
}

func TestEdits(t *testing.T) {
	h := load(t, "testdata/figure.yaml")
	n := len(h.Steps)
	if err := h.Insert(n, Step{Name: "body", Op: "solid.Sphere", Args: []any{1}}); err == nil {
		t.Error("a second step named body was accepted")
	}
	if len(h.Steps) != n {
		t.Error("a refused insert changed the history")
	}
	if err := h.Move(h.Index("cheeks"), h.Index("NoseTip")); err != nil {
		t.Fatal(err)
	}
	if h.Index("cheeks") != 3 || h.Index("NoseTip") != 4 {
		t.Errorf("after the move: cheeks at %d, NoseTip at %d", h.Index("cheeks"), h.Index("NoseTip"))
	}
	// The cheeks now come before the nose: their implicit input is the
	// figure, and the nose's is the cheeks.
	r := eval(t, h)
	if len(r.Bodies) != 1 {
		t.Errorf("bodies %+v", r.Bodies)
	}
	c := h.Clone()
	c.Steps[0].Args.(map[string]any)["radius"] = 99
	if h.Steps[0].Args.(map[string]any)["radius"] == 99 {
		t.Error("Clone shares args with the original")
	}
}

func TestPickAndAnchors(t *testing.T) {
	r := eval(t, load(t, "testdata/figure.yaml"))
	hit, ok := r.Pick(v3.XYZ(0, 0, 80), v3.Z(-1))
	if !ok || hit.Body.Name != "hollow" {
		t.Fatalf("pick: %+v %v", hit, ok)
	}
	validate.RequireNear(t, "top of the head (shell outer face)", hit.Point.Pos.Z, 22.8, 0.1)
	if _, ok := r.Pick(v3.XYZ(0, 0, 80), v3.Z(1)); ok {
		t.Error("a ray pointing away hit something")
	}
	var refs []string
	for _, a := range r.Anchors() {
		refs = append(refs, a.Ref)
	}
	all := strings.Join(refs, " ")
	for _, want := range []string{"NoseTip", "head.Top", "figure.BottomFrontLeft"} {
		if !strings.Contains(all, want) {
			t.Errorf("anchors lack %s", want)
		}
	}
}
