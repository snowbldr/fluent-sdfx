// Package shape is the 2D side of fluent-sdfx: a chainable wrapper around
// sdfx's SDF2 type with primitives, transforms, booleans, polygon and
// bezier builders, text, threads, cams, and gears.
//
// 2D follows screen convention: +Y is up. The 2D anchor set on *Shape is a
// halved version of the 3D set on *Solid: 4 face midpoints (Top, Bottom,
// Left, Right), 4 corners (TopRight, TopLeft, BottomRight, BottomLeft),
// plus AnchorAt(x, y int) for arbitrary -1/0/+1 coordinates.
//
// Conventions match the solid package:
//
//   - All angles are in degrees.
//   - Constructors panic on invalid input.
//   - Every method returns a new *Shape; the receiver is never mutated.
//
// 2D shapes extrude/revolve/sweep/loft to *solid.Solid. Cross-sections of a
// *Solid come back as a *Shape via shape.SliceOf / shape.SliceAt.
package shape
