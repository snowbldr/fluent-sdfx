# fluent-sdfx

A fluent, chainable API for [sdfx](https://github.com/snowbldr/sdfx) — Go's signed distance function CAD library.

fluent-sdfx wraps sdfx's SDF2 and SDF3 types with `shape.Shape` and `solid.Solid`, giving you a chainable API that reads like a description of the part you're building. All angles are in degrees. All constructors handle errors internally so you can chain without interruption.

📚 **[Read the docs](https://snowbldr.github.io/fluent-sdfx/)** — install, project setup, the dev loop, foundations, operations, cookbook recipes, and the full API reference.

## Install

```bash
go get github.com/snowbldr/fluent-sdfx
```

## Quick example

![Example output](example.png)

```go
package main

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func main() {
    // A cylinder with 4 holes drilled through it
	solid.Cylinder(20, 10, 1).Cut(
		solid.Cylinder(25, 2, 0).Multi(v3.X(5), v3.X(-5), v3.Y(5), v3.Y(-5)),
    ).STL("part.stl", 3.0)
}
```

The [quickstart](https://snowbldr.github.io/fluent-sdfx/quickstart/) walks through this exact part in five incremental steps.

## What's in the docs

- [**Install**](https://snowbldr.github.io/fluent-sdfx/install/) — Go, fluent-sdfx, f3d.
- [**Project setup**](https://snowbldr.github.io/fluent-sdfx/project-setup/) — scaffold a Go module and produce your first STL.
- [**Dev loop with stldev**](https://snowbldr.github.io/fluent-sdfx/dev-loop/) — watch-rebuild-preview iteration cycle.
- [**Quickstart**](https://snowbldr.github.io/fluent-sdfx/quickstart/) — five steps to a non-trivial part.
- **Foundations** — vectors, 2D shapes, 3D solids, booleans, transforms.
- **Operations** — extrude/revolve/loft/sweep, smooth blends, modifiers, patterns, cross-sections, text, output resolution, parametric helpers.
- **Cookbook** — bolt assembly, enclosure, gear.
- [**API reference**](https://snowbldr.github.io/fluent-sdfx/api-reference/) — every type and method, package by package.

## Repo layout

```
fluent-sdfx/
├── shape/      2D primitives, builders, transforms, booleans
├── solid/      3D primitives, transforms, booleans, modifiers
├── obj/        Parametric helpers (bolts, panels, gears, …)
├── render/     Output formats (STL, 3MF, DXF, SVG, PNG)
├── mesh/       Triangle-mesh utilities
├── plane/      Plane helpers for cross-sections
├── units/      Constants and unit conversions
├── vec/        Vector types (v2, v3, v2i, v3i, p2)
├── examples/   77 example projects
├── tutorial/   Runnable code paired with the docs site
└── docs/       The fntags-based documentation site
```

## License

MIT
