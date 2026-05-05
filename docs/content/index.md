# Introduction

A fluent, chainable Go API for signed-distance-function CAD. Build 3D parts by composing primitives, transforms, and booleans.

fluent-sdfx wraps Go's [sdfx](https://github.com/snowbldr/sdfx) library — a signed-distance-function CAD kernel — with two chainable types, `*shape.Shape` (2D) and `*solid.Solid` (3D), and a method on each operation that reads like a description of the part you're building.

<!-- src: tutorial/04-quickstart/05-polished-final/main.go -->
```go
// Quickstart step 5: print-shrinkage compensation and mesh decimation for a smaller STL.
package main

import (
	"github.com/snowbldr/fluent-sdfx/layout"
	"github.com/snowbldr/fluent-sdfx/solid"
)

const shrink = 1.0 / 0.999 // PLA shrinks ~0.1% on cooling.

func main() {
	solid.Cylinder(20, 10, 1).
		Cut(solid.Cylinder(25, 2, 0).Multi(layout.Polar(5, 4)...)).
		ScaleUniform(shrink).
		// 0.5 keeps half the triangles after meshoptimizer decimation.
		STL("out.stl", 3.0, 0.5)
}
```
That's the part you'll build by the end of the [quickstart](/quickstart/): a rounded cylinder with four through-holes, scaled for print-shrinkage compensation, decimated for a smaller mesh.

## Why fluent-sdfx

- **Chainable** — every transform and boolean returns a new object, so building a part is a single expression.
- **Degrees everywhere** — angles never need a manual `* math.Pi / 180`.
- **No error returns** — constructors panic on invalid input, so you can chain without `if err != nil` interruptions. CAD geometry errors are programming bugs, not runtime conditions.
- **Anchor-based positioning** — say `lid.OnTopOf(box.Top())` instead of doing bounding-box math. 27 anchors per part (faces, edges, corners) plus layout helpers for polar, grid, and rect-corner patterns.
- **Real renderer** — outputs STL (parallel marching cubes), 3MF, DXF, SVG, PNG.

## What's in the box

- **`shape`** — 2D primitives, booleans, transforms, polygon and bezier builders, text, threads, cams, gears.
- **`solid`** — 3D primitives, booleans, transforms, smooth blends, shells, sweeps, screws, lofts.
- **Positioning** — anchor selectors (`Top`, `BottomLeft`, …) and placement verbs (`On`, `Above`, `Inside`, …) on every Solid, plus a [`layout`](/positioning/) package of helpers (`Polar`, `Grid`, `RectCorners`, …) for spreading parts via `.Multi(...)`.
- **`obj`** — parametric helpers: bolts, nuts, washers, panels, eurorack, standoffs, pipes, servos, gridfinity, gears.
- **`vec/{v2,v3,p2}`** — vector types with named constructors so you can skip `Vec{X: ..., Y: ...}` boilerplate.
- **`render`** — output helpers and lower-level renderer constructors.
- **`validate`** — mesh validation for `_test.go`: watertight check, volume regression guards, FDM overhang detection. See [Testing & validation](/testing-validation/).

## Where to next

Start with [installing](/install/) the library and tooling, set up a [project](/project-setup/), then build something in the [quickstart](/quickstart/).
