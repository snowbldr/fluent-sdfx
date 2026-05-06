// Package solid is the 3D side of fluent-sdfx: a chainable wrapper around
// sdfx's SDF3 type with primitives, transforms, booleans, smooth blends, and
// anchor-based positioning.
//
// The two top-level types in fluent-sdfx are *solid.Solid (3D, here) and
// *shape.Shape (2D, in the shape package). They interconvert: a *Shape
// becomes a *Solid via Extrude/Revolve/Loft/SweepHelix/Screw; a *Solid
// becomes a *Shape via Slice (2D cross-section).
//
// Conventions:
//
//   - All angles are in degrees.
//   - Constructors panic on invalid input; CAD-geometry errors are
//     programming bugs, not runtime conditions. Chain without if err != nil.
//   - Every transform and boolean returns a new *Solid; the receiver is
//     never mutated.
//
// See the positioning page (https://snowbldr.github.io/fluent-sdfx/positioning/)
// for the anchor verbs (Top, OnTopOf, Inside, BottomAt, …) and the layout
// package for repeated-copy patterns (Polar, Grid, RectCorners, …).
package solid
