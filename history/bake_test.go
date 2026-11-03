package history

import (
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	"github.com/snowbldr/fluent-sdfx/history/internal/opsgen"
	"github.com/snowbldr/fluent-sdfx/validate"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// TestOpsTableFresh fails when ops_gen.go is stale: a library function was
// added, renamed or re-documented and `go generate ./history/` was not run.
func TestOpsTableFresh(t *testing.T) {
	want, err := opsgen.Generate("..")
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile("ops_gen.go")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Error("ops_gen.go is stale; run go generate ./history/")
	}
}

// pulls is a sphere with n pulls spiralling round it, as a session of
// dragging on a form makes.
func pulls(n int) *History {
	h := &History{Steps: []Step{{Name: "ball", Op: "solid.Sphere", Args: map[string]any{"radius": 12}}}}
	for i := 0; i < n; i++ {
		a := float64(i) * 2.4
		at := []any{12 * math.Cos(a) * 0.8, 12 * math.Sin(a) * 0.8, 12 * 0.6 * math.Cos(float64(i))}
		dir := []any{math.Cos(a), math.Sin(a), 0.3}
		h.Steps = append(h.Steps, Step{Name: fmt.Sprintf("p%d", i), Op: "sculpt.Pull", Args: map[string]any{"at": at, "by": dir, "radius": 5}})
	}
	return h
}

func TestBakesKeepLongChainsFast(t *testing.T) {
	h := pulls(24)
	raw := NewRuntime(Options{}).Eval(h)
	baked := NewRuntime(Options{BakeEvery: 6, BakeCells: 128}).Eval(h)
	if err := baked.Err(); err != nil {
		t.Fatal(err)
	}
	var at []int
	for _, s := range baked.Steps {
		if s.Baked {
			at = append(at, s.Index)
		}
	}
	t.Logf("baked at steps %v", at)
	if fmt.Sprint(at) != "[6 12 18 24]" {
		t.Errorf("baked at %v, want every sixth pull", at)
	}

	cost := func(r *Result) time.Duration {
		s, bb := r.Solid(), r.Solid().Bounds().Box
		start := time.Now()
		for i := 0; i < 20000; i++ {
			f := float64(i) / 20000
			validate.At(s, bb.Min.Add(bb.Size().Mul(v3.XYZ(f, math.Mod(f*7, 1), math.Mod(f*13, 1)))))
		}
		return time.Since(start)
	}
	cr, cb := cost(raw), cost(baked)
	t.Logf("20k evaluations: %v unbaked, %v baked", cr, cb)
	if cb > cr {
		t.Errorf("the baked chain is slower to evaluate (%v > %v)", cb, cr)
	}
	// Same surface to within a cell: compare near the surface, since away
	// from it a grid's field outside its lattice is only a lower bound.
	cell := raw.Solid().Bounds().Box.Size().MaxComponent() / 128
	bb := raw.Solid().Bounds().Box
	worst, n := 0.0, 0
	for i := 0; i < 200000 && n < 3000; i++ {
		f := float64(i) / 200000
		p := bb.Min.Add(bb.Size().Mul(v3.XYZ(math.Mod(f*101, 1), math.Mod(f*7919, 1), math.Mod(f*104729, 1))))
		want := validate.At(raw.Solid(), p)
		if math.Abs(want) > 2*cell {
			continue
		}
		n++
		worst = math.Max(worst, math.Abs(validate.At(baked.Solid(), p)-want))
	}
	t.Logf("near the surface (%d points): baked and unbaked fields agree to %.3f mm; a cell is %.3f", n, worst, cell)
	if worst > cell {
		t.Errorf("the bake moved the surface by up to %.3f, more than a cell (%.3f)", worst, cell)
	}
}
