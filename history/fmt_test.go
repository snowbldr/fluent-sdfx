package history

import (
	"strings"
	"testing"
)

func TestMarshalReadsWell(t *testing.T) {
	h := load(t, "testdata/figure.yaml")
	b, err := h.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%s", b)
	for _, want := range []string{
		"args: {radius: 12}",
		"min: {op: solid.RoundMin, args: {k: 3}}",
		"ray: {from: [0, -40, 14], dir: [0, 1, 0]}",
		"  - name: body\n    op: solid.Sphere",
	} {
		if !strings.Contains(string(b), want) {
			t.Errorf("marshalled history lacks %q", want)
		}
	}
	back, err := Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	b2, _ := back.Marshal()
	if string(b2) != string(b) {
		t.Error("marshal is not stable across a round trip")
	}
}
