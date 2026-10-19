// Package sculpt is an experiment: the operations a person reaches for in
// a 3D sculpting app, as functions on fluent-sdfx solids. Pull a point of
// the surface out, smooth every edge at a radius, bucket-fill an enclosed
// void, drop a piece that is not attached. These are the verbs for organic
// and decorative parts, where there is no bore or flange to name and the
// description is "fatten the belly, soften that ridge, fill that gap".
//
// Two kinds of operation live here, and the split matters:
//
//   - Closed-form warps on the field (Pull). Cheap, exact in position,
//     but they leave the field a distance bound rather than a distance:
//     the sign and the zero surface are right, the magnitudes are not.
//     Anything that reads magnitudes afterwards (Shell, Offset, a smooth
//     blend, a clearance probe) will be wrong until the field is
//     redistanced.
//
//   - Grid operations (Grid.Redistance, Fill, Keep, Smooth). The field is
//     sampled once onto a lattice; the grid verbs work on the samples and
//     the result is a true signed distance to within the cell size, so
//     everything built on it means what it says. Flood fill and
//     connectivity are only possible here, because a distance function
//     knows how far the wall is, not what the wall encloses.
//
// The usual chain is: build with primitives and Pull, Sample onto a Grid,
// apply the grid verbs, take Solid() back and finish with Shell, booleans
// and export. Resolution is the one thing to choose: features smaller than
// a cell are lost, and every surface is placed to within about half a cell.
//
// Status: experimental. Not part of docs/llms.txt, no API promises. The
// tests in this directory are the realistic cases the verbs are judged on.
package sculpt
