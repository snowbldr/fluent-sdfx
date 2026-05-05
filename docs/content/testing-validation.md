# Testing & validation

The `validate` package inspects a rendered solid's mesh for printability and regression-test signals — triangle count, surface area, volume, watertightness, and overhang area. It works on the marching-cubes output, so the metrics reflect what would actually be exported to STL.

Use it for two things:

- **CI guardrails on geometry** — `RequireWatertight`, `RequireVolumeNear`, `RequireMaxOverhang` in a normal `_test.go` file. A refactor that flips an SDF, breaks a boolean, or accidentally adds a hidden cavity shows up as a failed test instead of a wrecked print.
- **Pre-flight printability check** — eliminate the "render, load in slicer, find supports, fix code, repeat" loop. A `RequireMaxOverhang` test catches >45° overhangs before you spend 30 seconds dragging the STL into Cura.

## Quick example

```go
import "github.com/snowbldr/fluent-sdfx/validate"

func TestBracketPrintable(t *testing.T) {
    s := buildBracket()                       // your assembly under test
    validate.RequireWatertight(t, s, 5.0)     // sealed mesh, no holes
    validate.RequireMaxOverhang(t, s, 5.0, 45.0)  // FDM rule of thumb
    validate.RequireVolumeNear(t, s, 5.0, 12500.0, 0.02)  // ±2% of expected
}
```

`5.0` is `cellsPerMM` — the same render density argument as `s.STL("part.stl", 5.0)`.

## Stats

`validate.Of(s, cellsPerMM)` returns a `Stats` struct with everything in one pass:

```go
st := validate.Of(part, 5.0)
fmt.Printf("%.1f mm³  %.0f mm²  %d tris  watertight=%v  overhang=%.1f mm²\n",
    st.Volume, st.SurfaceArea, st.Triangles, st.Watertight, st.OverhangArea)
```

| Field | Meaning |
|---|---|
| `Triangles` | triangle count after marching cubes |
| `SurfaceArea` | mm² — sum of triangle areas |
| `Volume` | mm³ — signed-tetrahedron sum (garbage if not watertight) |
| `BoundaryEdges` | edges shared by exactly 1 triangle (0 = sealed) |
| `Watertight` | `BoundaryEdges == 0` |
| `Bounds` | axis-aligned bounding box |
| `OverhangArea` | mm² overhanging > 45° from vertical (FDM threshold) |

The triangle-level helpers (`Volume`, `SurfaceArea`, `OverhangArea`, `OverhangFaces`, `IsWatertight`) take a precomputed mesh if you've already rendered.

## What each check catches

### Watertight / boundary edges

Marching cubes shares vertex positions across cells, so `==` comparison of edge endpoints is a reliable closure check — no tolerance needed. A non-zero boundary count means **holes in the mesh**: an SDF discontinuity (e.g. a `Sub3D` that strays outside a bounding region), a Lipschitz violation in a custom SDF, or a bug in a non-monotonic blend. The slicer will silently fail or behave unpredictably on a non-watertight mesh.

```go
validate.RequireWatertight(t, part, cellsPerMM)
```

### Volume drift

A regression-style guard: pin the volume of a finalized part to its known value. Any future change that alters the volume by more than the relative tolerance fails the test. Surfaces a subtle cluster of bugs:

- A `Cut` that accidentally became a `Union` (or vice versa)
- An SDF whose sign flipped (subtracting the inverse of a part)
- A boolean that lost an argument
- A wall thickness changed when it shouldn't have

```go
validate.RequireVolumeNear(t, part, cellsPerMM, expectedMM3, relTol)  // 0.02 = ±2%
```

The marching-cubes mesh approaches the true volume as `cellsPerMM` increases, so the test value should be measured at the same density it'll run at. Lower density → looser tolerance.

### Overhang detection

For FDM 3D printing, faces overhanging more than ~45° from vertical typically need supports — extra material that wastes filament, ages the printer, and leaves witness marks on the underside. `OverhangArea` reports how many mm² of mesh fall in that category.

```go
validate.RequireMaxOverhang(t, part, cellsPerMM, 45.0)               // strict
validate.RequireMaxOverhang(t, part, cellsPerMM, 45.0, 5.0)          // tolerate 5 mm² of noise
```

The triangle's overhang angle is the angle between its outward-pointing normal and horizontal. A flat ceiling (normal = -Z) is 90°. A 45° overhang sits exactly at the threshold. The check assumes Z+ is the print direction; rotate the part first if you'll print at a different orientation.

`validate.OverhangFaces(tris, 45.0)` returns the offending triangles — write them to their own STL with `mesh.SaveSTL(...)` and load alongside the original part to visualize problem areas.

## Performance notes

Every helper renders the SDF first, so density (`cellsPerMM`) drives runtime. Treat validation tests like preview renders — pick a density that catches the bugs without burning seconds. For unit-test runs, `5.0`-`8.0` is usually enough; for a "ship gate" test you might crank it to the same density as your final STL.

For a single solid that runs through multiple checks, render once and pass the triangle slice to the lower-level helpers:

```go
tris := mesh.CollectTriangles(s, render.NewMarchingCubesOctreeParallel(solid.CellsFor(s, cellsPerMM)))
ok, edges := validate.IsWatertight(tris)
vol := validate.Volume(tris)
overhang := validate.OverhangArea(tris, 45.0)
```

## Best practices

- **Add a `_test.go` to every cookbook-style finished part.** Volume + watertight + overhang together cost a fraction of a second at typical preview density.
- **Pin the volume only on parts whose dimensions are intentional.** Don't pin parts that change with parameters in CI; pin a representative configuration.
- **For overhang failures, look at `OverhangFaces` first.** Often the fix is rotating the print orientation, not redesigning — confirm with the visual before changing geometry.
- **Skip `Volume`/`OverhangArea` on non-watertight inputs.** Both metrics return nonsense if the mesh has holes; gate them behind `RequireWatertight`.
