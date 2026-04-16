package meshopt

// #cgo CXXFLAGS: -O2 -std=c++11
// #include "meshopt_wrapper.h"
import "C"
import "unsafe"

// Simplify reduces an unindexed triangle mesh using quadric edge collapse.
// vertices is a flat []float32 with 9 floats per triangle (3 vertices * 3 coords).
// targetTriangles is the desired output triangle count.
// targetError controls maximum allowed error (e.g. 0.01 = 1% of mesh extent).
// Returns the simplified vertices (flat float32 slice, 9 per triangle) and triangle count.
func Simplify(vertices []float32, numTriangles int, targetTriangles int, targetError float32) ([]float32, int) {
	out := make([]float32, numTriangles*9) // worst case: no reduction
	n := C.meshopt_simplify_unindexed(
		(*C.float)(unsafe.Pointer(&out[0])),
		(*C.float)(unsafe.Pointer(&vertices[0])),
		C.size_t(numTriangles),
		C.size_t(targetTriangles),
		C.float(targetError),
	)
	count := int(n)
	return out[:count*9], count
}
