#ifndef MESHOPT_WRAPPER_H
#define MESHOPT_WRAPPER_H

#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

// Simplify an unindexed triangle mesh (3 floats per vertex, 3 vertices per triangle).
// Returns the number of triangles in the output.
// out_vertices must be pre-allocated to at least num_triangles * 9 floats.
size_t meshopt_simplify_unindexed(
    float* out_vertices,
    const float* vertices,
    size_t num_triangles,
    size_t target_triangles,
    float target_error
);

#ifdef __cplusplus
}
#endif

#endif
