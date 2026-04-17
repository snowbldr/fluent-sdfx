package solid

import "github.com/deadsy/sdfx/sdf"

// Benchmark reports the evaluation speed of the solid's SDF3.
func (s *Solid) Benchmark(description string) {
	sdf.BenchmarkSDF3(description, s.SDF3)
}
