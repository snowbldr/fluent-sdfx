# fluent-sdfx

A fluent, chainable API for [sdfx](https://github.com/snowbldr/sdfx) — Go's signed distance function CAD library.

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
    v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
    "github.com/snowbldr/fluent-sdfx/solid"
)

func main() {
    // A cylinder with 4 holes drilled through it
    body := solid.Cylinder(20, 10, 1)
    hole := solid.Cylinder(25, 2, 0)

    part := body.Cut(
        hole.Translate(v3.X(5)),
        hole.Translate(v3.X(-5)),
        hole.Translate(v3.Y(5)),
        hole.Translate(v3.Y(-5)),
    )

    part.STL("part.stl", 3.0)
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
| `Flange1(distance, centerR, sideR)` | Flange profile (two tangent-joined circles) |
| `CubicSpline(knots)` | Closed cubic spline |
| `Nagon(n, radius)` | Vertices of a regular N-gon |
| `AcmeThread/ISOThread/ANSIButtressThread/PlasticButtressThread` | Screw thread profiles |
| `ThreadLookup(name)` | Look up a standard thread by name |
| `FlatFlankCam` / `MakeFlatFlankCam` | Flat-flank cam profile |
| `ThreeArcCam` / `MakeThreeArcCam` | Three-arc cam profile |
| `GearRack(params)` | Linear gear rack |
| `Text(font, str, height)` + `LoadFont(path)` | Truetype-rendered text |
| `NewPoly()` | Fluent polygon builder with `.Smooth/.Chamfer/.Arc/.Rel/.Polar` vertices |
| `NewBezier()` | Fluent bezier builder with slope handles |
| `Wrap2D(sdf2)` | Wrap a raw `sdf.SDF2` |

**Transforms:** `Translate`, `Rotate`, `Scale`, `MirrorX`, `MirrorY`, `ScaleUniform`, `Center`, `CenterAndScale`, `Transform`

**2D → 3D:** each returns a `*solid.Solid`. Package-level constructors with the same names remain available in `solid` for callers that already hold a raw `sdf.SDF2`.

| Method | Description |
|---|---|
| `Extrude(height)` | Linear extrusion |
| `ExtrudeRounded(height, round)` | Extrusion with rounded top/bottom edges |
| `TwistExtrude(height, twist)` | Extrude while rotating about Z (radians) |
| `ScaleExtrude(height, scale)` | Extrude while scaling the profile |
| `ScaleTwistExtrude(height, twist, scale)` | Scaled + twisted extrusion |
| `Revolve()` | Full revolution around the Y axis |
| `RevolveAngle(angleDeg)` | Partial revolution |
| `Screw(height, start, pitch, n)` | Helical screw thread |
| `SweepHelix(radius, turns, height, flatEnds)` | Sweep along a helix path |
| `LoftTo(top, height, round)` | Transition from this profile to another |

**Booleans:** `Union`, `Cut`, `Intersect`

**Smooth blends:** `SmoothUnion`, `SmoothCut`, `SmoothIntersect` (pair with min/max funcs from `solid`)

**Modifiers:** `Offset`, `CutLine`, `Split`, `Elongate`, `Cache`

**Patterns:** `Array`, `SmoothArray`, `RotateCopy`, `RotateUnion`, `SmoothRotateUnion`, `Multi`, `LineOf`

**Bench / inspect:** `.Benchmark(description)` reports SDF2 evaluation speed; `.MeshBoxes()` returns the acceleration-structure boxes of a mesh-backed shape.

**Bounding boxes:** `Box2` type alias + `NewBox2(center, size)` for creating 2D AABBs (e.g., for random point sampling).

**Mesh-from-lines:** `Mesh2D(segments)` / `Mesh2DSlow(segments)` build a `*Shape` from `[]*Line2`.

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

Methods on `*Solid` return raw `sdf.SDF2`; helpers in the `shape` package wrap the result as `*Shape` so you can chain `.ToDXF / .ToSVG / .ToPNG` and further 2D ops.

| API | Description |
|---|---|
| `s.Slice2D(origin, normal)` | Cross-section through a solid (method on `*Solid`) |
| `s.SliceAt(plane.Plane)` | Cross-section at a `plane.Plane` (method on `*Solid`) |
| `solid.Slice(s, origin, normal)` | Package-level form, same as `s.Slice2D` |
| `shape.SliceOf(s, origin, normal)` | Slice and wrap as `*Shape` |
| `shape.SliceAt(s, plane.Plane)` | Slice at a plane and wrap as `*Shape` |

**Transforms:** `Translate`, `RotateX`, `RotateY`, `RotateZ`, `RotateAxis`, `Scale`, `ScaleUniform`, `Transform`, `ZeroZ`, `Center`, `RotateToVector`

**Mirrors:** `MirrorXY`, `MirrorXZ`, `MirrorYZ`, `MirrorXeqY`

**Booleans:** `Union`, `UnionAll`, `Cut`, `Intersect`

**Smooth blends:** `SmoothUnion`, `SmoothDifference`, `SmoothIntersection` with `RoundMin`, `ChamferMin`, `ExpMin`, `PowMin`, `PolyMin`, `PolyMax`

**Mesh / voxel:** `Mesh(triangles)`, `MeshSlow(triangles)`, `.Voxel(cells, progress)`

**Modifiers:** `Shrink`, `Grow`, `Correct`, `Shell`, `CutPlane`, `Split`, `Elongate`, `Offset`

**Patterns:** `Array`, `SmoothArray`, `RotateCopyZ`, `RotateUnionZ`, `SmoothRotateUnionZ`, `Multi`, `LineOf`, `Orient`

**Bench:** `.Benchmark(description)` reports SDF3 evaluation speed.

**Bounding boxes:** `Box3` type alias + `NewBox3(center, size)` for creating 3D AABBs.

### `obj` — Parametric Helpers

Wraps `sdfx/obj` higher-level parametric parts so they plug into the fluent API directly. All helpers take parameter structs by value and panic on invalid input. 2D helpers return `*shape.Shape`; 3D helpers return `*solid.Solid`.

| Function | Description |
|---|---|
| `Angle2D/3D(AngleParams)` | L-profile angle bracket |
| `Arrow3D(ArrowParms)` / `DirectedArrow3D` / `Axes3D(p0, p1)` | Arrows and axis indicators |
| `Bolt(BoltParms)` | Bolt with hex or socket cap head |
| `Nut(NutParms)` / `ThreadedCylinder(ThreadedCylinderParms)` | Threaded fasteners |
| `Washer2D/3D(WasherParms)` | Washers |
| `Hex2D(r, round)` / `Hex3D(r, h, round)` / `HexHead3D(r, h, "tb")` | Hex shapes/prisms |
| `ChamferedCylinder(s, kb, kt)` | Chamfer a cylinder top/bottom |
| `CounterBoredHole3D` / `ChamferedHole3D` / `CounterSunkHole3D` | Hole styles |
| `BoltCircle2D/3D` / `CircleGrille2D/3D(CircleGrilleParms)` | Hole patterns |
| `KeyedHole2D/3D(KeyedHoleParms)` | Circle with N key slots |
| `Keyway2D/3D(KeywayParameters)` | Shaft with keyway |
| `Knurl3D(KnurlParms)` / `KnurledHead3D(r, h, pitch)` | Knurled surfaces |
| `Panel2D/3D(PanelParms)` / `PanelHole3D` | Panels with mount holes |
| `EuroRackPanel2D/3D(EuroRackParms)` | Eurorack module panels |
| `PanelBox3D(PanelBoxParms)` | Panel-box enclosure — returns `[]*solid.Solid` |
| `Standoff3D(StandoffParms)` | PCB standoff |
| `Pipe3D(oR, iR, L)` / `StdPipe3D(name, units, L)` / `PipeLookup(name, units)` | Standard pipes |
| `PipeConnector3D` / `StdPipeConnector3D` | Multi-port pipe connectors |
| `Servo3D` / `Servo2D` / `ServoLookup(name)` / `ServoHorn(ServoHornParms)` | Servo bodies and horns |
| `Spring2D/3D(SpringParms)` / `SpringLength` | Flat planar springs |
| `GfBase(GfBaseParms)` / `GfBody(GfBodyParms)` | Gridfinity base/body |
| `DrainCover(DrainCoverParms)` | Drain covers |
| `Display(DisplayParms, negative)` | Display bezels |
| `DroneMotorArm` / `DroneMotorArmSocket` | Drone arm parts |
| `FingerButton2D(FingerButtonParms)` | Finger button profile |
| `InvoluteGear(InvoluteGearParms)` | Involute gear profile |
| `Geneva2D(GenevaParms)` | Geneva drive — returns driver and driven |
| `IsocelesTrapezoid2D` / `IsocelesTriangle2D` | Simple polygon primitives |
| `TruncRectPyramid3D(TruncRectPyramidParms)` | Truncated rectangular pyramid |
| `ImportSTL(path, ...)` / `ImportTriMesh(tris, ...)` | Load a mesh as SDF3 |
| `NewStraightTab` / `NewAngleTab` / `NewScrewTab` / `AddTabs(s, tab, upper, mset)` | Splitting tabs |

### `render` — Output Formats

```go
part.STL("output.stl", 6.0)      // 6 cells per mm
part.STL("output.stl", 6.0, 0.9) // optional decimation — remove 90% of triangles (keep 10%)
part.ThreeMF("output.3mf", 6.0)  // alias: .MF3(...)

profile.ToDXF("profile.dxf", 400)
profile.ToSVG("profile.svg", 400)
profile.ToPNG("preview.png", 800, 600)
```

**Picking `cellsPerMM`.** 3D output resolution is a mesh *density* in cells per millimeter along the longest bounding-box axis. Pick the number based on how much detail you want per unit of real-world size, independent of how big the part is:

| Scenario | `cellsPerMM` | Resulting cells on longest axis |
| --- | --- | --- |
| 500 mm enclosure, rough preview | `0.2` | 100 |
| 500 mm enclosure, final | `2.0` | 1000 |
| 50 mm bracket, typical | `5.0` | 250 |
| 10 mm gear, detailed | `20.0` | 200 |
| 1 mm sphere | `12.0` | 32 (floored by `solid.MinCells`) |

Render time scales roughly with cells³ — halving `cellsPerMM` is an 8× speedup. Drop it low for iteration, crank it for final output.

Tiny parts are floored at `solid.MinCells` (default 32) so a 1 mm sphere at `cellsPerMM=3` still produces a recognizable mesh instead of an empty STL. Raise `solid.MinCells` for more sub-mm detail, lower it (or set to `1`) for raw behavior.

3D uses the parallel marching cubes octree renderer. 2D uses the quadtree marching-squares renderer. Optional STL mesh decimation via [meshoptimizer](https://github.com/zeux/meshoptimizer) (requires CGo).

For lower-level access the `render` package exposes `ToSTL`, `To3MF`, `ToDXF`, `ToSVG`, `ToDXFWith`, `ToSVGWith`, `ToPNG`, `SaveDXF`, `SaveSVG`, and renderer constructors `NewMarchingCubesOctreeParallel`, `NewMarchingSquaresQuadtree`, `NewDualContouring2D`. For interactive 2D output, `NewPNG(path, bb, pixels)` and `NewDXF(path)` return `*PNG` / `*DXF` drawing targets with `RenderSDF2`, `Triangle`, `Line`, `Lines`, `Box`, `Points`, `Save` methods.

### `mesh` — Triangle-Mesh Utilities

Type aliases (`Triangle3`, `Triangle2`, `Line2`, `TriangleISet`) plus helpers that bypass the SDF pipeline: `ToTriangles(solid, r)`, `CollectTriangles(solid, r)`, `CountBoundaryEdges`, `SaveSTL(path, tris)`, `Delaunay2d(vs)`, `Delaunay2dSlow(vs)`, `VertexToLine(vs, closed)`.

### `units` — Constants and Conversions

Re-exports the constants and helpers you'll want from `sdf`: `Pi`, `Tau`, `Mil`, `MillimetresPerInch`, `InchesPerMillimetre`, `DtoR(deg)`, `RtoD(rad)`, `EqualFloat64`, `ErrMsg`.

### `vec/{v2,v3,v2i,v3i,p2,conv}` — Vector Types

Re-exports of sdfx's vector packages plus named constructors so you can skip the `Vec{X: ..., Y: ...}` boilerplate:

```go
v3.X(5)          // Vec{X: 5}
v3.YZ(2, 3)      // Vec{Y: 2, Z: 3}
v3.XYZ(1, 2, 3)  // Vec{X: 1, Y: 2, Z: 3}
v3.Zero          // Vec{}
```

`v2`/`v3` provide `X`, `Y`, `Z`, `XY`, `XZ`, `YZ`, `XYZ`, `Zero` (float64). `v2i`/`v3i` provide integer-component variants. `p2.R(r)`, `p2.T(theta)`, `p2.RT(r, theta)` build polar vectors. `conv` provides `V2ToV3`, `P2ToV2`, `V2ToV2i`, and the other cross-type conversions.

`*Shape` and `*Solid` satisfy `sdf.SDF2` / `sdf.SDF3` directly — their `BoundingBox()` and `Evaluate()` return the underlying sdfx types. Use `s.Bounds()` when you want the fluent `v2.Box` / `v3.Box` wrappers with chainable methods like `ScaleAboutCenter`, `RandomSet`, `Include`.

### `plane` — Slice Plane Helpers

Axis normals and plane constructors for cross-sections:

```go
plane.X, plane.Y, plane.Z       // unit normals
plane.XY, plane.XZ, plane.YZ    // planes through the origin
plane.AtX(5), plane.AtY(0), plane.AtZ(10)  // axis-aligned at offset
plane.At(origin, normal)        // arbitrary
```

Pair with `s.SliceAt(plane.AtZ(10))` or `shape.SliceAt(s, plane.XY)`.

## Dev loop: `stldev`

For an iterative dev loop, pair fluent-sdfx with [`stldev`](https://github.com/snowbldr/stldev) — a small CLI that watches your Go source, re-runs your build command, and previews the generated STLs in [f3d](https://f3d.app) with auto-reload and tiled windows.

```bash
go install github.com/snowbldr/stldev@latest

stldev -cmd "go run ." part.stl
```

See the [`stldev` README](https://github.com/snowbldr/stldev) for flags, the `--` passthrough for f3d args, and a Makefile-driven workflow that scales cleanly to multi-part assemblies.

## Design

- **Chainable**: every transform/boolean returns a new object, so you can chain: `solid.Cylinder(10, 5, 0).RotateX(90).Translate(v3.Z(10))`
- **Degrees everywhere**: all angle parameters are in degrees, converted to radians internally
- **No error returns**: constructors panic on invalid input rather than returning errors — CAD geometry errors are programming bugs, not runtime conditions

## License

MIT
