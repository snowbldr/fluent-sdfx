# API reference

A condensed table of every type and method exposed by fluent-sdfx, package by package.

This page is a flat reference — every constructor and method, grouped by package. For tutorials and worked examples, see the rest of the docs.

## `shape` — 2D primitives

`*shape.Shape` wraps `sdf.SDF2`. All methods return a new `*Shape`.

### Constructors

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
| `Flange1(distance, centerR, sideR)` | Two tangent-joined circles |
| `CubicSpline(knots)` | Closed cubic spline |
| `Nagon(n, radius)` | Vertices of a regular N-gon |
| `AcmeThread`, `ISOThread`, `ANSIButtressThread`, `PlasticButtressThread` | Screw thread profiles |
| `ThreadLookup(name)` | Standard thread by name |
| `FlatFlankCam` / `MakeFlatFlankCam` | Flat-flank cam profile |
| `ThreeArcCam` / `MakeThreeArcCam` | Three-arc cam profile |
| `GearRack(params)` | Linear gear rack |
| `Text(font, str, height)` + `LoadFont(path)` | Truetype-rendered text |
| `NewPoly()` | Fluent polygon builder with `.Smooth/.Chamfer/.Arc/.Rel/.Polar` |
| `NewBezier()` | Fluent bezier builder with slope handles |
| `Wrap2D(sdf2)` | Wrap a raw `sdf.SDF2` |

### Transforms

`Translate`, `TranslateX`, `TranslateY`, `TranslateXY`, `Rotate`, `Scale`, `MirrorX`, `MirrorY`, `ScaleUniform`, `Center`, `CenterAndScale`, `Transform`

### 2D → 3D

Each returns a `*solid.Solid`. Package-level constructors with the same names also exist in `solid` for callers holding raw `sdf.SDF2`.

| Method | Description |
|---|---|
| `Extrude(height)` | Linear extrusion |
| `ExtrudeRounded(height, round)` | Extrusion with rounded edges |
| `TwistExtrude(height, twist)` | Twist around Z (radians) |
| `ScaleExtrude(height, scale)` | Extrude with linear scaling |
| `ScaleTwistExtrude(height, twist, scale)` | Both |
| `Revolve()` | Full revolution around Y |
| `RevolveAngle(angleDeg)` | Partial revolution |
| `Screw(height, start, pitch, n)` | Helical screw thread |
| `SweepHelix(radius, turns, height, flatEnds)` | Sweep along a helix |
| `LoftTo(top, height, round)` | Transition to another profile |

### Booleans, blends, patterns

**Booleans:** `Union` / `Add`, `Cut` / `Difference`, `Intersect`

**Smooth blends:** `SmoothUnion` / `SmoothAdd`, `SmoothCut` / `SmoothDifference`, `SmoothIntersect` (paired with min/max funcs from `solid`)

**Modifiers:** `Offset`, `CutLine`, `Split`, `Elongate`, `Cache`

**Patterns:** `Array`, `SmoothArray`, `RotateCopy`, `RotateUnion`, `SmoothRotateUnion`, `Multi`, `LineOf`

### Inspection & sampling

- `.Benchmark(description)` reports SDF2 evaluation speed.
- `.Normal(p, eps)` returns the surface normal at p (finite differences).
- `.Raycast(from, dir, ...)` sphere-traces a ray; returns hit point, distance, step count.
- `.GenerateMesh(grid)` samples interior points on an integer grid.
- `.MeshBoxes()` returns the acceleration-structure boxes of a mesh-backed shape.

### Bounding boxes

- `Box2` type alias + `NewBox2(center, size)` for 2D AABBs.
- `s.Bounds()` returns the AABB with `Center`, `Size`, `Vertices`, `ScaleAboutCenter`, `RandomSet`.

### Mesh-from-lines

`Mesh2D(segments)` / `Mesh2DSlow(segments)` build a `*Shape` from `[]*Line2`.

---

## `solid` — 3D primitives

`*solid.Solid` wraps `sdf.SDF3`. All methods return a new `*Solid`.

### Constructors

| Function | Description |
|---|---|
| `Cylinder(height, radius, round)` | Cylinder along Z |
| `Box(size, round)` | Rounded box |
| `Sphere(radius)` | Sphere |
| `Cone(height, r0, r1, round)` | Truncated cone |
| `Capsule(height, radius)` | Cylinder with hemispherical caps |
| `Torus(majorR, minorR)` | Torus |
| `Gyroid(scale)` | Infinite gyroid surface |

### From 2D profiles

| Function | Description |
|---|---|
| `Extrude(profile, height)` | Linear extrusion |
| `ExtrudeRounded(profile, height, round)` | Rounded edges |
| `TwistExtrude(profile, height, twist)` | Twisted (radians) |
| `ScaleExtrude(profile, height, scale)` | Linearly scaled |
| `ScaleTwistExtrude(profile, height, twist, scale)` | Both |
| `Revolve(profile)` | Full Y-axis revolution |
| `RevolveAngle(profile, angle)` | Partial |
| `Screw(profile, height, start, pitch, n)` | Helical thread |
| `Loft(bottom, top, height, round)` | Profile transition |
| `SweepHelix(profile, radius, turns, height, flatEnds)` | Sweep along helix |

### Cross-section

| API | Description |
|---|---|
| `s.Slice2D(origin, normal)` | Cross-section through a solid |
| `s.SliceAt(plane.Plane)` | Cross-section at a plane |
| `solid.Slice(s, origin, normal)` | Package-level form |
| `shape.SliceOf(s, origin, normal)` | Slice and wrap as `*Shape` |
| `shape.SliceAt(s, plane.Plane)` | Slice at a plane and wrap as `*Shape` |

### Transforms

`Translate`, `TranslateX/Y/Z`, `TranslateXY/XZ/YZ`, `TranslateXYZ`, `RotateX/Y/Z`, `RotateAxis`, `Scale`, `ScaleUniform`, `Transform`, `ZeroZ`, `Center`, `RotateToVector`

**Mirrors:** `MirrorXY`, `MirrorXZ`, `MirrorYZ`, `MirrorXeqY`

### Booleans, blends, patterns

**Booleans:** `Union` / `Add`, `UnionAll`, `Cut` / `Difference`, `Intersect`

**Smooth blends:** `SmoothUnion` / `SmoothAdd`, `SmoothCut` / `SmoothDifference`, `SmoothIntersect`, paired with `RoundMin`, `ChamferMin`, `ExpMin`, `PowMin`, `PolyMin`, `PolyMax`

**Modifiers:** `Shrink`, `Grow`, `Correct`, `Shell`, `CutPlane`, `Split`, `Elongate`, `Offset`

**Patterns:** `Array`, `SmoothArray`, `RotateCopyZ`, `RotateUnionZ`, `SmoothRotateUnionZ`, `Multi`, `LineOf`, `Orient`

### Mesh / voxel

`Mesh(triangles)`, `MeshSlow(triangles)`, `.Voxel(cells, progress)`

### Inspection & sampling

`Benchmark`, `Normal`, `Raycast`

### Bounding boxes

- `Box3` type alias + `NewBox3(center, size)` for 3D AABBs.
- `s.Bounds()` returns the AABB with `Center`, `Size`, `ScaleAboutCenter`, `Anchor(x,y,z int)` (returns the world-space point on the box for the given unit-cube position), `.Solid()` (turn the box back into a `*Solid`).

### Positioning

The canonical way to place parts. See [/positioning](/positioning/) for the full visual reference.

**Anchor selectors on `*Solid`** (return `AnchoredSolid`):

- 6 faces: `Top`, `Bottom`, `Right`, `Left`, `Back`, `Front`
- 12 edges: `TopRight`, `TopLeft`, `TopFront`, `TopBack`, `BottomRight`, `BottomLeft`, `BottomFront`, `BottomBack`, `FrontRight`, `FrontLeft`, `BackRight`, `BackLeft`
- 8 corners: `TopFrontRight`, `TopFrontLeft`, `TopBackRight`, `TopBackLeft`, `BottomFrontRight`, `BottomFrontLeft`, `BottomBackRight`, `BottomBackLeft`
- Arbitrary: `AnchorAt(x, y, z int)` — each component in `{-1, 0, +1}`

**Placement verbs on `AnchoredSolid`:**

- Relative (return `Placement`): `On(target)`, `Above(target, gap...)`, `Below`, `RightOf`, `LeftOf`, `Behind`, `InFrontOf`
- Absolute (return `*Solid`): `At(v3.Vec)`, `AtX`, `AtY`, `AtZ`
- Anchor tweaks (return `AnchoredSolid`): `ShiftX`, `ShiftY`, `ShiftZ`

**`Placement` finalizers** — the chain's subject (the moved/active solid) is what's kept, matching the `s.Cut(other)` convention:

- Plain: `.Union()` / `.Add()` (commutative), `.Cut()` / `.Difference()` (subtracts base from moved), `.Intersect()` (commutative)
- Smooth: `.SmoothUnion(min)` / `.SmoothAdd(min)`, `.SmoothCut(max)` / `.SmoothDifference(max)`, `.SmoothIntersect(max)`
- Escape: `.Solid()` returns the moved solid alone, no boolean — useful for drilling: `body.Cut(tool.Top().On(body.Top()).Solid())`

**Sugar verbs on `*Solid`** (relative; return `Placement`):

- `OnTopOf`, `UnderneathOf`, `LeftOf`, `RightOf`, `InFrontOf`, `BehindOf` (each takes `target AnchoredSolid, gap ...float64`)
- `Inside(other *Solid)`

**Absolute scalar setters on `*Solid`** (return `*Solid`):

- `BottomAt(z)`, `TopAt(z)`, `LeftAt(x)`, `RightAt(x)`, `FrontAt(y)`, `BackAt(y)`
- `CenterAt(p v3.Vec)`

**2D mirror on `*Shape`** — same names, halved set:

- 4 faces: `Top`, `Bottom`, `Left`, `Right`
- 4 corners: `TopRight`, `TopLeft`, `BottomRight`, `BottomLeft`
- `AnchorAt(x, y int)` (2D follows screen convention: `+Y` is up)
- All placement verbs and sugar — `On`, `Above`, `Below`, `RightOf`, `LeftOf`, `At`, `AtX`, `AtY`, `ShiftX`, `ShiftY`, `OnTopOf`, `UnderneathOf`, `LeftOf`, `RightOf`, `Inside`, `BottomAt`, `TopAt`, `LeftAt`, `RightAt`, `CenterAt`
- `Placement2D` finalizers: `.Union()`, `.Cut()`, `.Intersect()`, `.Shape()`

---

## `layout` — Position arrays

Pure functions returning `[]v3.Vec` (or `[]v2.Vec`), designed to flow into the variadic `.Multi(positions...)`.

| Function | Returns | Use |
|---|---|---|
| `Polar(radius, n)` | `[]v3.Vec` | n positions on a circle (z=0); first at +X |
| `PolarArc(radius, n, startDeg, sweepDeg)` | `[]v3.Vec` | n along an arc |
| `Grid(stepX, stepY, nx, ny)` | `[]v3.Vec` | XY grid centered on origin |
| `Line(p0, p1, n)` | `[]v3.Vec` | n equally spaced from p0 to p1 |
| `RectCorners(width, depth)` | `[]v3.Vec` | 4 XY corners centered on origin |
| `BoxCorners(size)` | `[]v3.Vec` | 8 box corners centered on origin |
| `Polar2(radius, n)` | `[]v2.Vec` | 2D circle |
| `PolarArc2(radius, n, startDeg, sweepDeg)` | `[]v2.Vec` | 2D arc |
| `Grid2(stepX, stepY, nx, ny)` | `[]v2.Vec` | 2D grid |
| `Line2(p0, p1, n)` | `[]v2.Vec` | 2D evenly spaced from p0 to p1 |
| `RectCorners2(width, depth)` | `[]v2.Vec` | 2D 4 corners centered on origin |

---

## `obj` — Parametric helpers

Wraps `sdfx/obj`'s parametric parts so they plug into the fluent API directly. All helpers take parameter structs by value and panic on invalid input. 2D helpers return `*shape.Shape`; 3D helpers return `*solid.Solid`.

| Function | Description |
|---|---|
| `Angle2D/3D(AngleParams)` | L-profile angle bracket |
| `Arrow3D(ArrowParms)` / `DirectedArrow3D` / `Axes3D(p0, p1)` | Arrows and axes |
| `Bolt(BoltParms)` | Bolt with hex or socket cap head |
| `Nut(NutParms)` / `ThreadedCylinder(ThreadedCylinderParms)` | Threaded fasteners |
| `Washer2D/3D(WasherParms)` | Washers |
| `Hex2D(r, round)` / `Hex3D(r, h, round)` / `HexHead3D(r, h, "tb")` | Hex shapes & prisms |
| `ChamferedCylinder(s, kb, kt)` | Chamfer top/bottom of a cylinder |
| `CounterBoredHole3D` / `ChamferedHole3D` / `CounterSunkHole3D` | Hole styles |
| `BoltCircle2D/3D` / `CircleGrille2D/3D(CircleGrilleParms)` | Hole patterns |
| `KeyedHole2D/3D(KeyedHoleParms)` | Circle with key slots |
| `Keyway2D/3D(KeywayParameters)` | Shaft with keyway |
| `Knurl3D(KnurlParms)` / `KnurledHead3D(r, h, pitch)` | Knurled surfaces |
| `Panel2D/3D(PanelParms)` / `PanelHole3D` | Panels with mount holes |
| `EuroRackPanel2D/3D(EuroRackParms)` | Eurorack module panels |
| `PanelBox3D(PanelBoxParms)` | Panel-box enclosure (returns `[]*solid.Solid`) |
| `Standoff3D(StandoffParms)` | PCB standoff |
| `Pipe3D(oR, iR, L)` / `StdPipe3D(name, units, L)` / `PipeLookup(name, units)` | Standard pipes |
| `PipeConnector3D` / `StdPipeConnector3D` | Multi-port pipe connectors |
| `Servo3D` / `Servo2D` / `ServoLookup(name)` / `ServoHorn(ServoHornParms)` | Servo bodies & horns |
| `Spring2D/3D(SpringParms)` / `SpringLength` | Flat planar springs |
| `GfBase(GfBaseParms)` / `GfBody(GfBodyParms)` | Gridfinity base/body |
| `DrainCover(DrainCoverParms)` | Drain covers |
| `Display(DisplayParms, negative)` | Display bezels |
| `DroneMotorArm` / `DroneMotorArmSocket` | Drone arm parts |
| `FingerButton2D(FingerButtonParms)` | Finger button profile |
| `InvoluteGear(InvoluteGearParms)` | Involute gear profile |
| `Geneva2D(GenevaParms)` | Geneva drive |
| `IsocelesTrapezoid2D` / `IsocelesTriangle2D` | Simple polygon primitives |
| `TruncRectPyramid3D(TruncRectPyramidParms)` | Truncated rect pyramid |
| `ImportSTL(path, ...)` / `ImportTriMesh(tris, ...)` | Load mesh as SDF3 |
| `NewStraightTab` / `NewAngleTab` / `NewScrewTab` / `AddTabs(s, tab, upper, mset)` | Splitting tabs |

---

## `render` — Output formats

```go
part.STL("output.stl", 6.0)              // 6 cells per mm
part.STL("output.stl", 6.0, 0.9)         // optional decimation (90% removed = keep 10%)
part.ThreeMF("output.3mf", 6.0)          // alias: .MF3(...)

profile.ToDXF("profile.dxf", 400)
profile.ToSVG("profile.svg", 400)
profile.ToPNG("preview.png", 800, 600)
```

### Picking `cellsPerMM`

Mesh density in cells per millimetre along the longest bounding-box axis. Pick by detail-per-real-millimetre, independent of part size.

| Scenario | `cellsPerMM` | Resulting cells |
|---|---|---|
| 500mm enclosure, rough preview | `0.2` | 100 |
| 500mm enclosure, final | `2.0` | 1000 |
| 50mm bracket, typical | `5.0` | 250 |
| 10mm gear, detailed | `20.0` | 200 |
| 1mm sphere | `12.0` | 32 (floored) |

Render time scales roughly with cells³ — halving `cellsPerMM` is an 8× speedup.

Tiny parts are floored at `solid.MinCells` (default 32) so a 1mm sphere at `cellsPerMM=3` still produces a recognisable mesh. Raise `solid.MinCells` for more sub-mm detail; lower (or set to `1`) for raw behavior.

3D uses parallel marching-cubes octree. 2D uses quadtree marching squares. Optional STL mesh decimation via [meshoptimizer](https://github.com/zeux/meshoptimizer) (requires CGo).

### Lower-level access

`render` exposes:

- `ToSTL`, `To3MF`, `ToDXF`, `ToSVG`, `ToDXFWith`, `ToSVGWith`, `ToPNG`, `SaveDXF`, `SaveSVG`
- Renderer constructors: `NewMarchingCubesOctreeParallel`, `NewMarchingSquaresQuadtree`, `NewDualContouring2D`
- `NewPNG(path, bb, pixels)` and `NewDXF(path)` return `*PNG` / `*DXF` drawing targets with `RenderSDF2`, `Triangle`, `Line`, `Lines`, `Box`, `Points`, `Save`

---

## `mesh` — Triangle-mesh utilities

Type aliases (`Triangle3`, `Triangle2`, `Line2`, `TriangleISet`) plus helpers that bypass the SDF pipeline:

- `ToTriangles(solid, r)`
- `CollectTriangles(solid, r)`
- `CountBoundaryEdges`
- `IsWatertight`
- `SaveSTL(path, tris)`
- `Delaunay2d(vs)`, `Delaunay2dSlow(vs)`
- `VertexToLine(vs, closed)`

---

## `validate` — Mesh validation & test helpers

Inspect a rendered solid for printability and regression-test signals. See [Testing & validation](/testing-validation/).

| Function | Use |
|---|---|
| `Of(s, cellsPerMM) Stats` | render and compute every metric in one pass |
| `OfMesh(tris) Stats` | same, on a precomputed mesh (no Bounds) |
| `Volume(tris) float64` | mm³ — signed-tetrahedron sum |
| `SurfaceArea(tris) float64` | mm² — sum of triangle areas |
| `IsWatertight(tris) (bool, int)` | `true` and `0` for sealed mesh |
| `OverhangArea(tris, deg) float64` | mm² of faces overhanging > deg from vertical |
| `OverhangFaces(tris, deg) []Triangle3` | the offending triangles |
| `RequireWatertight(t, s, cellsPerMM)` | `*testing.T` helper, fails on holes |
| `RequireVolumeNear(t, s, cellsPerMM, expectedMM3, relTol)` | regression guard |
| `RequireMaxOverhang(t, s, cellsPerMM, maxAngleDeg, tinyArea...)` | printability gate |

Stats fields: `Triangles, SurfaceArea, Volume, BoundaryEdges, Watertight, Bounds, OverhangArea` (the last is at the FDM 45° threshold).

---

## `units` — Constants & conversions

Re-exports from sdfx: `Pi`, `Tau`, `Mil`, `MillimetresPerInch`, `InchesPerMillimetre`, `DtoR(deg)`, `RtoD(rad)`, `EqualFloat64`, `ErrMsg`.

---

## `vec/{v2,v3,v2i,v3i,p2,conv}` — Vector types

Re-exports of sdfx's vector packages with named constructors:

```go
v3.X(5)          // {X: 5}
v3.YZ(2, 3)      // {Y: 2, Z: 3}
v3.XYZ(1, 2, 3)  // all three
v3.Zero          // {0, 0, 0}
```

`v2`/`v3` provide `X`, `Y`, `Z`, `XY`, `XZ`, `YZ`, `XYZ`, `Zero` (float64). `v2i`/`v3i` are integer variants. `p2.R(r)`, `p2.T(theta)`, `p2.RT(r, theta)` build polar vectors.

`*Shape` and `*Solid` satisfy `sdf.SDF2` / `sdf.SDF3` directly — their `BoundingBox()` and `Evaluate()` return raw sdfx types. Use `s.Bounds()` for the fluent `v2.Box` / `v3.Box` wrappers (with `ScaleAboutCenter`, `RandomSet`, `Include`).

---

## `plane` — Slice plane helpers

Axis normals and plane constructors for cross-sections:

```go
plane.X, plane.Y, plane.Z       // unit normals (v3.Vec)
plane.XY, plane.XZ, plane.YZ    // planes through the origin
plane.AtX(5), plane.AtY(0), plane.AtZ(10)  // axis-aligned at offset
plane.At(origin, normal)        // arbitrary plane
```

Pair with `s.SliceAt(plane.AtZ(10))` or `shape.SliceAt(s, plane.XY)`.

---

## Design principles

- **Chainable.** Every transform/boolean returns a new object — `solid.Cylinder(10, 5, 0).RotateX(90).Translate(v3.Z(10))` works as a single expression.
- **Degrees everywhere.** All angle parameters in degrees, converted to radians internally. Exception: `TwistExtrude` and `ScaleTwistExtrude` take twist in radians (use `units.DtoR(...)` to convert).
- **No error returns.** Constructors panic on invalid input — CAD geometry errors are programming bugs, not runtime conditions.
