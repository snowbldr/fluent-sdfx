package shape

import "github.com/snowbldr/sdfx/sdf"

// Benchmark reports the evaluation speed of the shape's SDF2.
func (s *Shape) Benchmark(description string) {
	sdf.BenchmarkSDF2(description, s.SDF2)
}
