#include "meshopt_wrapper.h"
#include "meshoptimizer.h"
#include <vector>
#include <cstring>

extern "C" size_t meshopt_simplify_unindexed(
    float* out_vertices,
    const float* vertices,
    size_t num_triangles,
    size_t target_triangles,
    float target_error
) {
    size_t vertex_count = num_triangles * 3;
    size_t index_count = vertex_count;

    // Generate index buffer from unindexed vertices
    std::vector<unsigned int> remap(vertex_count);
    size_t unique_vertices = meshopt_generateVertexRemap(
        remap.data(), NULL, index_count,
        vertices, vertex_count, sizeof(float) * 3
    );

    // Build indexed mesh
    std::vector<unsigned int> indices(index_count);
    std::vector<float> indexed_verts(unique_vertices * 3);
    for (size_t i = 0; i < index_count; i++) {
        indices[i] = remap[i];
    }
    for (size_t i = 0; i < vertex_count; i++) {
        if (remap[i] < unique_vertices) {
            memcpy(&indexed_verts[remap[i] * 3], &vertices[i * 3], sizeof(float) * 3);
        }
    }

    // Simplify
    std::vector<unsigned int> simplified(index_count);
    size_t simplified_count = meshopt_simplify(
        simplified.data(), indices.data(), index_count,
        indexed_verts.data(), unique_vertices, sizeof(float) * 3,
        target_triangles * 3, target_error, 0, NULL
    );

    // Convert back to unindexed triangle soup
    size_t out_triangles = simplified_count / 3;
    for (size_t i = 0; i < simplified_count; i++) {
        memcpy(&out_vertices[i * 3], &indexed_verts[simplified[i] * 3], sizeof(float) * 3);
    }

    return out_triangles;
}
