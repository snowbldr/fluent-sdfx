# Coming from another CAD library

If you've used OpenSCAD, CadQuery, Build123d, or sdfx directly, here's how the idioms map. fluent-sdfx is closest to OpenSCAD in spirit — implicit geometry, no parametric history — and closest to Build123d in API style — chainable, expression-based.

## At a glance

| Need to do | OpenSCAD | CadQuery | Build123d | sdfx (raw) | fluent-sdfx |
|---|---|---|---|---|---|
| Make a box | `cube([20,20,20])` | `Workplane().box(20,20,20)` | `Box(20,20,20)` | `sdf.Box3D(...)` | `solid.Box(v3.XYZ(20,20,20), 0)` |
| Make a cylinder | `cylinder(h=20, r=5)` | `.circle(5).extrude(20)` | `Cylinder(5, 20)` | `sdf.Cylinder3D(20, 5, 0)` | `solid.Cylinder(20, 5, 0)` |
| Subtract | `difference()` | `.cut()` | `-` | `sdf.Difference3D(a, b)` | `a.Cut(b)` |
| Union | `union()` | `.union()` | `+` | `sdf.Union3D(a, b)` | `a.Union(b)` |
| Translate | `translate([0,0,5])` | `.translate((0,0,5))` | `Pos(z=5) * obj` | `sdf.Transform3D(s, m)` | `s.TranslateZ(5)` |
| Rotate | `rotate([0,0,90])` | `.rotate(0,0,90)` | `Rot(z=90) * obj` | `sdf.Transform3D(s, m)` | `s.RotateZ(90)` |
| Sit on build plate | `translate([0,0,h/2])` | manual | `Plane.XY * obj` | manual | `s.BottomAt(0)` |
| Stack on top of body | `translate([0,0,h])` | manual | manual | manual | `cap.OnTopOf(body.Top()).Union()` |
| Pattern around a circle | `for(i=...)` | `.polarArray(...)` | `PolarLocations(...)` | manual | `s.Multi(layout.Polar(r, n)...)` |
| Hex bolt | manual or library | `.hex(...)`-ish | manual | `sdf.HexHead3D(...)` | `obj.Bolt(obj.BoltParms{...})` |
| Output an STL | `--export-format=stl` | `cq.exporters.export(...)` | `export_stl(...)` | render API | `s.STL("part.stl", 5.0)` |

## Coming from OpenSCAD

OpenSCAD users will find fluent-sdfx the most familiar conceptually — both are CSG-style implicit-geometry tools where you compose primitives with booleans. The syntax differs (Go method chains vs. OpenSCAD's nesting/`children()`) but the mental model is the same.

```scad
// OpenSCAD
difference() {
  cylinder(h=20, r=10);
  for (i = [0:3]) {
    rotate([0, 0, i*90])
      translate([5, 0, 0])
        cylinder(h=25, r=2);
  }
}
```

```go
// fluent-sdfx
solid.Cylinder(20, 10, 0).
    Cut(solid.Cylinder(25, 2, 0).Multi(layout.Polar(5, 4)...)).
    STL("part.stl", 5.0)
```

What you get over OpenSCAD:
- Real types and functions — no string-template gymnastics for parametrics.
- A real debugger and language tooling.
- Anchor-based positioning so you don't write bbox math.
- Everything renders to STL, 3MF, DXF, SVG, PNG out of the box.
- Smooth blends (round-min, chamfer-min, etc.) are first-class — no `minkowski()` workarounds.

What OpenSCAD does that we don't:
- Live preview is on you (use [stldev](/dev-loop/) or your own watcher).
- No GUI; no render-as-you-type WebGL pane.

## Coming from CadQuery / Build123d

The B-rep crowd. CadQuery and Build123d work with parametric histories of NURBS surfaces; fluent-sdfx works with signed-distance fields and converts to triangles via marching cubes. That has trade-offs:

**You'll like fluent-sdfx for:**
- Smooth blends, fillets, and unions that *just work* — no failed boolean operations from degenerate surfaces.
- Implicit infinity (`solid.Cylinder` is conceptually infinite at the SDF level, clipped at render time) — composing booleans never hits "non-manifold" errors.
- A simpler mental model: every shape is a function `f(point) → distance`. No edges, no faces, no half-edges.
- Go's tooling — types, refactoring, no Python venv.

**You'll miss:**
- Exact NURBS surfaces. fluent-sdfx output is meshed, so curves are tessellated.
- True chamfers/fillets along arbitrary edges. SDF-based smooth blends apply globally to the join between two solids — they're not edge-selectors.
- Parametric assembly history. There's no "go back two steps and redo with a different radius."
- Thin-feature accuracy. Marching cubes at low density misses sub-cell features.

```python
# CadQuery
result = (cq.Workplane("XY")
    .circle(10).extrude(20)
    .faces(">Z").workplane()
    .pushPoints([(5,0),(-5,0),(0,5),(0,-5)])
    .circle(2).cutThruAll())
```

```go
// fluent-sdfx
solid.Cylinder(20, 10, 0).
    Cut(solid.Cylinder(25, 2, 0).Multi(layout.RectCorners(10, 10)...)).
    STL("part.stl", 5.0)
```

The CadQuery-style face-selector + pushPoints workflow is replaced by anchor placement + `layout` helpers. There's no `.faces(">Z")` because fluent-sdfx parts don't have explicit faces — but there's `.Top()`, which is the equivalent for placement.

## Coming from sdfx (the underlying library)

fluent-sdfx is a thin chainable wrapper on top of [sdfx](https://github.com/snowbldr/sdfx). The SDF kernel, marching cubes, mesh decimation, threads, gears — everything compute-heavy — is sdfx. fluent-sdfx adds the chainable types, anchor positioning, layout helpers, and the validate package on top.

You can drop down to raw sdfx at any time: `s.SDF3` exposes the wrapped `sdf.SDF3` value. Conversely, any `sdf.SDF3` can be lifted into a `*solid.Solid` via `solid.Wrap(...)` (or `shape.Wrap2D(...)` for 2D).

Migration:

```go
// Raw sdfx
body, _ := sdf.Cylinder3D(20, 10, 0)
hole, _ := sdf.Cylinder3D(25, 2, 0)
hole = sdf.Transform3D(hole, sdf.Translate3d(v3.Vec{X: 5}))
result := sdf.Difference3D(body, hole)
render.ToSTL(result, "part.stl", render.NewMarchingCubesOctree(200))
```

```go
// fluent-sdfx
solid.Cylinder(20, 10, 0).
    Cut(solid.Cylinder(25, 2, 0).TranslateX(5)).
    STL("part.stl", 5.0)
```

Why bother with the wrapper:
- No `_, err := ...` interruptions; constructors panic on invalid input (CAD-geometry errors are programming bugs, not runtime conditions).
- Fluent chaining; no re-binding intermediates.
- Anchor-based positioning ([positioning](/positioning/)) and layout helpers ([positioning](/positioning/)).
- Degrees, not radians — `RotateZ(90)`, not `Transform3D(s, RotateZ(math.Pi/2))`.
- A single `s.STL(path, cellsPerMM)` call replaces renderer construction.
- The `validate` package for testing-friendly mesh checks ([testing & validation](/testing-validation/)).

## When fluent-sdfx is the wrong choice

- **Mechanical engineering with strict NURBS / G-code interop.** Use CadQuery, Build123d, or commercial parametric CAD.
- **Architectural drawing or CNC tooling that needs explicit edges, faces, and surface IDs.** SDFs don't have those.
- **You hate compiled languages and just want to script `for i in range`.** OpenSCAD or CadQuery (Python) will be lighter.

For 3D-printable parts, hobby mechanical assemblies, generative geometry, parametric helpers, threads/gears, and anything where smooth implicit blends shine — that's the sweet spot.
