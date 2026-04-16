# fluent-sdfx

A fluent, chainable API for [sdfx](https://github.com/deadsy/sdfx) — Go's signed distance function CAD library.

fluent-sdfx wraps sdfx's SDF2 and SDF3 types with `shape.Shape` and `solid.Solid`, giving you a chainable API that reads like a description of the part you're building. All angles are in degrees. All constructors handle errors internally so you can chain without interruption.

## Install

```bash
go get github.com/snowbldr/fluent-sdfx
```

## Quick Example

![Example output](example.png)

```go
package main

import (
    v3 "github.com/deadsy/sdfx/vec/v3"
    "github.com/snowbldr/fluent-sdfx/solid"
)

func main() {
    // A cylinder with 4 holes drilled through it
    body := solid.Cylinder(20, 10, 1)
    hole := solid.Cylinder(25, 2, 0)

    part := body.Cut(
        hole.Translate(v3.Vec{X: 5}),
        hole.Translate(v3.Vec{X: -5}),
        hole.Translate(v3.Vec{Y: 5}),
        hole.Translate(v3.Vec{Y: -5}),
    )

    part.ToSTL("part.stl", 300)
}
```

## Packages

### `shape` — 2D Primitives

`shape.Shape` wraps `sdf.SDF2`. All methods return a new `*Shape`.

**Constructors:**

| Function | Description |
|---|---|
| `Rect(size, round)` | Rounded rectangle |
| `Circle(radius)` | Circle |
| `Polygon(pts)` | Arbitrary polygon (cache-friendly SDF) |
| `Line(length, round)` | Line segment |
| `ArcSpiral(a, k, start, end, d)` | Archimedean spiral |
| `Star(outer, inner, points)` | Star polygon |
| `Hexagon(radius)` | Regular hexagon |
| `Triangle(radius)` | Equilateral triangle |
| `Cross(width, thickness)` | Cross / plus sign |
| `WireGroove(radius, depth, tailAngle)` | Wire groove profile |
| `Wrap2D(sdf2)` | Wrap a raw `sdf.SDF2` |

**Transforms:** `Translate`, `Rotate`, `Scale`, `MirrorX`, `MirrorY`, `ScaleUniform`, `Center`, `CenterAndScale`, `Transform`

**Booleans:** `Union`, `Cut`, `Intersect`

**Modifiers:** `Offset`, `CutLine`, `Elongate`

**Patterns:** `Array`, `RotateCopy`, `RotateUnion`, `Multi`, `LineOf`

### `solid` — 3D Primitives

`solid.Solid` wraps `sdf.SDF3`. All methods return a new `*Solid`.

**Constructors:**

| Function | Description |
|---|---|
| `Cylinder(height, radius, round)` | Cylinder |
| `Box(size, round)` | Box |
| `Sphere(radius)` | Sphere |
| `Cone(height, r0, r1, round)` | Truncated cone |
| `Capsule(height, radius)` | Cylinder with hemispherical caps |
| `Torus(majorR, minorR)` | Torus |
| `Gyroid(scale)` | Infinite gyroid surface |

**From 2D profiles:**

| Function | Description |
|---|---|
| `Extrude(profile, height)` | Linear extrusion |
| `ExtrudeRounded(profile, height, round)` | Extrusion with rounded edges |
| `TwistExtrude(profile, height, twist)` | Twisted extrusion |
| `ScaleExtrude(profile, height, scale)` | Scaled extrusion |
| `ScaleTwistExtrude(profile, height, twist, scale)` | Scaled + twisted extrusion |
| `Revolve(profile)` | Full revolution around Y |
| `RevolveAngle(profile, angle)` | Partial revolution |
| `Screw(profile, height, start, pitch, n)` | Helical screw thread |
| `Loft(bottom, top, height, round)` | Transition between two profiles |
| `SweepHelix(profile, radius, turns, height, flatEnds)` | Sweep along helix path |

**Cross-section:**

| Function | Description |
|---|---|
| `Slice(solid, origin, dir)` | Cut a 2D cross-section, returns `*shape.Shape` |

**Transforms:** `Translate`, `RotateX`, `RotateY`, `RotateZ`, `RotateAxis`, `Scale`, `ScaleUniform`, `Transform`, `ZeroZ`, `Center`, `RotateToVector`

**Mirrors:** `MirrorXY`, `MirrorXZ`, `MirrorYZ`, `MirrorXeqY`

**Booleans:** `Union`, `UnionAll`, `Cut`, `Intersect`

**Modifiers:** `Shrink`, `Grow`, `Correct`, `Shell`, `CutPlane`, `Elongate`

**Patterns:** `Array`, `RotateCopyZ`, `RotateUnionZ`, `Multi`, `LineOf`, `Orient`

### `render` — STL Output

```go
part.ToSTL("output.stl", 600)      // 600 cells along longest axis
part.ToSTL("output.stl", 600, 0.5) // optional decimation (keep 50%)
```

Uses sdfx's parallel marching cubes octree renderer. Optional mesh decimation via [meshoptimizer](https://github.com/zeux/meshoptimizer) (requires CGo).

For lower-level access: `render.ToSTL(sdf, path, renderer, factor...)` accepts any `render.Render3`.

## Design

- **Chainable**: every transform/boolean returns a new object, so you can chain: `solid.Cylinder(10, 5, 0).RotateX(90).Translate(v3.Vec{Z: 10})`
- **Degrees everywhere**: all angle parameters are in degrees, converted to radians internally
- **No error returns**: constructors panic on invalid input rather than returning errors — CAD geometry errors are programming bugs, not runtime conditions
- **Zero raw sdfx needed**: the API covers every operation in sdfx's SDF2, SDF3, and matrix packages

## License

MIT
