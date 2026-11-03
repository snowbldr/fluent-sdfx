package history

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/snowbldr/fluent-sdfx/history/internal/testparts"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func init() {
	Register("test.Ring", testparts.Ring, "A washer with notches cut in its rim.")
	Register("test.Peg", testparts.Peg, "A peg with a rounded end.", "height", "radius")
}

func TestRegisteredFunctionsRun(t *testing.T) {
	h, err := Parse([]byte(`
steps:
  - name: ring
    op: test.Ring
    args: {outer: 10, inner: 6, height: 3, notches: 4}
  - name: peg
    op: test.Peg
    args: {height: 12, radius: 1}
    then: [{op: TranslateX, args: [8]}]
  - op: Union
    on: ring
    args: [peg]
`))
	if err != nil {
		t.Fatal(err)
	}
	r := eval(t, h)
	want := testparts.Ring(testparts.RingParams{Outer: 10, Inner: 6, Height: 3, Notches: 4}).Union(testparts.Peg(12, 1).TranslateX(8))
	congruent(t, "registered parts", r.Solid(), want, 1e-12)
	validate.RequireNear(t, "notch cut", validate.At(r.Solid(), v3.XYZ(-9.5, 0, 0)), 0.5, 0.01)
}

func TestPaletteDescribesEverything(t *testing.T) {
	p := Palette()
	byOp := map[string]OpInfo{}
	for _, o := range p {
		byOp[o.Kind+" "+o.Receiver+" "+o.Op] = o
	}
	t.Logf("%d ops", len(p))

	ring, ok := byOp["function  test.Ring"]
	if !ok || !ring.Registered || ring.Doc == "" || len(ring.Params) != 1 || ring.Params[0].Widget != "struct" || len(ring.Params[0].Fields) != 4 {
		t.Errorf("test.Ring: %+v", ring)
	}
	pull := byOp["function  sculpt.Pull"]
	want := map[string]string{"s": "step", "at": "point", "by": "offset", "radius": "length"}
	for _, q := range pull.Params {
		if want[q.Name] != q.Widget {
			t.Errorf("sculpt.Pull %s: widget %s, want %s", q.Name, q.Widget, want[q.Name])
		}
	}
	if pull.Input != "s" {
		t.Errorf("sculpt.Pull input %q, want s", pull.Input)
	}
	tr := byOp["method solid.Solid Translate"]
	if tr.Returns != "solid.Solid" || len(tr.Params) != 1 || tr.Params[0].Name != "v" {
		t.Errorf("Translate: %+v", tr)
	}
	if _, ok := byOp["method solid.Solid STL"]; ok {
		t.Error("STL writes a file and must not be an op")
	}
	if u := byOp["function  solid.SmoothUnion"]; len(u.Params) != 2 || !u.Params[1].Variadic || u.Params[0].Widget != "blend" {
		t.Errorf("SmoothUnion: %+v", u)
	}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"op":"solid.Sphere"`) {
		t.Error("palette JSON lacks solid.Sphere")
	}
	t.Logf("palette JSON: %d KB", len(b)/1024)
}

func TestRegisterRefusesBadOps(t *testing.T) {
	for name, fn := range map[string]func(){
		"no dot":       func() { Register("Bracket", testparts.Peg, "") },
		"library name": func() { Register("solid.Sphere", testparts.Peg, "") },
		"not a func":   func() { Register("my.X", 3, "") },
		"no result":    func() { Register("my.Y", func() {}, "") },
		"names":        func() { Register("my.Z", testparts.Peg, "", "height") },
		"twice":        func() { Register("test.Peg", testparts.Peg, "", "height", "radius") },
	} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("%s: Register did not panic", name)
				}
			}()
			fn()
		}()
	}
}
