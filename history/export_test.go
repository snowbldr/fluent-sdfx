package history

import (
	"encoding/json"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// probeMain evaluates the exported Build at the points in points.json and
// prints the values per body.
const probeMain = `package main

import (
	"encoding/json"
	"os"

	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
	b, _ := os.ReadFile(os.Args[1])
	var pts []v3.Vec
	json.Unmarshal(b, &pts)
	out := map[string][]float64{}
	for name, s := range Build() {
		for _, p := range pts {
			out[name] = append(out[name], validate.At(s, p))
		}
	}
	json.NewEncoder(os.Stdout).Encode(out)
}
`

// The exported Go compiles and builds the same part: every body's field
// matches the runtime's at random points around it.
func TestExportCompilesAndIsCongruent(t *testing.T) {
	if testing.Short() {
		t.Skip("compiles a program per history")
	}
	for _, name := range []string{"figure", "ornament", "cluster", "kitchen", "registered"} {
		t.Run(name, func(t *testing.T) {
			h := load(t, "testdata/"+name+".yaml")
			src, err := ExportGo(h, ExportOptions{Source: name + ".yaml"})
			if err != nil {
				t.Fatal(err)
			}
			if name == "figure" {
				t.Logf("exported:\n%s", src)
				if !strings.Contains(string(src), "// agent, chat/123") {
					t.Error("the nose step's provenance is not in the export")
				}
			}
			// A scratch package inside this one, so its imports resolve
			// (internal/testparts included); go ignores dot directories.
			dir, err := os.MkdirTemp(".", ".exportcheck-")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)
			os.WriteFile(filepath.Join(dir, "part.go"), src, 0o644)
			os.WriteFile(filepath.Join(dir, "probe.go"), []byte(probeMain), 0o644)

			r := NewRuntime(Options{}).Eval(h)
			rng := rand.New(rand.NewSource(2))
			bb := r.Solid().Bounds().Box.Enlarge(v3.XYZ(2, 2, 2))
			var pts []v3.Vec
			for i := 0; i < 500; i++ {
				s := bb.Size()
				pts = append(pts, bb.Min.Add(v3.XYZ(rng.Float64()*s.X, rng.Float64()*s.Y, rng.Float64()*s.Z)))
			}
			pj, _ := json.Marshal(pts)
			os.WriteFile(filepath.Join(dir, "points.json"), pj, 0o644)

			cmd := exec.Command("go", "run", ".", "points.json")
			cmd.Dir = dir
			out, err := cmd.Output()
			if err != nil {
				msg := err.Error()
				if ee, ok := err.(*exec.ExitError); ok {
					msg = string(ee.Stderr)
				}
				t.Fatalf("the export does not run: %s\n%s", msg, src)
			}
			var got map[string][]float64
			if err := json.Unmarshal(out, &got); err != nil {
				t.Fatal(err)
			}
			if len(got) != len(r.Bodies) {
				t.Errorf("export has %d bodies, runtime %d", len(got), len(r.Bodies))
			}
			for _, b := range r.Bodies {
				vals := got[b.Name]
				if len(vals) != len(pts) {
					t.Fatalf("body %s missing from the export", b.Name)
				}
				for k, p := range pts {
					if want := validate.At(b.Solid, p); vals[k] != want {
						t.Fatalf("body %s at %v: export %v, runtime %v", b.Name, p, vals[k], want)
					}
				}
			}
		})
	}
}
