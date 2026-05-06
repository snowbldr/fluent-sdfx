package layout_test

import (
	"fmt"
	"testing"

	"github.com/snowbldr/fluent-sdfx/layout"
)

// BenchmarkPolar measures Polar(radius, n) for a range of n.
// A regression introducing N² behavior would show the n=256 case
// growing super-linearly relative to n=4 / n=32.
func BenchmarkPolar(b *testing.B) {
	for _, n := range []int{4, 32, 256} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = layout.Polar(10, n)
			}
		})
	}
}

// BenchmarkGrid measures Grid(...) for a range of nx*ny totals.
func BenchmarkGrid(b *testing.B) {
	cases := []struct {
		nx, ny int
	}{
		{2, 2},   // 4
		{8, 8},   // 64
		{32, 32}, // 1024
	}
	for _, c := range cases {
		b.Run(fmt.Sprintf("n=%d", c.nx*c.ny), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = layout.Grid(5, 5, c.nx, c.ny)
			}
		})
	}
}
