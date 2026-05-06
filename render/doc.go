// Package render provides output formats and renderer constructors used
// internally by the *solid.Solid and *shape.Shape STL/3MF/DXF/SVG/PNG
// methods.
//
// Most users won't import render directly — calling s.STL(path, cellsPerMM)
// on a *Solid is enough. Reach for this package when you want to:
//
//   - Configure the renderer explicitly (e.g., NewMarchingCubesUniform vs
//     NewMarchingCubesOctreeParallel for different speed/quality trade-offs).
//   - Drive the rendering pipeline yourself for batched output.
//   - Use the lower-level ToDXF / ToSVG / ToPNG drawing-target API.
//
// All format-writer functions panic on filesystem errors (out of disk,
// permission denied, etc.) — consistent with fluent-sdfx's "errors are
// programming bugs" convention.
package render
