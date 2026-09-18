# Testing & validation

A solid is a signed distance field: a function from any point to how far it is from the surface, negative inside material. That means every question you could ask of a part by looking at it, measuring it with calipers, or trying to fit two pieces together is a question you can ask the code directly, exactly, and in a test that runs every time you change something. The `validate` package makes each of those questions one line.

This is the thing a programming language gives you that a modelling package cannot. A drawing lets you look; a field lets you measure, sweep, compare and prove.

## Why it exists

The package was assembled from a real project: a device, its slip-casting moulds and its extraction tooling, modelled over two months with an AI writing the fluent-sdfx code. Sixteen of its eighteen test files re-implemented the same loop by hand: evaluate the field at a named point, fail if the sign is wrong. A dozen re-wrote the same bisection to find a surface. The one printability auditor lived inside a `package main` test file, so most of the project had no printability check at all. Every technique here was hand-rolled there first, saved a print or a day, and is now a call.

Two rules from that record are worth stating up front.

**Probes verify what you thought to check; renders catch what you did not.** Use both. Render from several directions (`tools/inspect`), section through the feature under discussion, and then write probes for what you saw.

**Measure; do not derive.** Every multi-round failure in the record was a model reasoning about a mirror, a hand, an angle or a threshold instead of reading it off the geometry. The helpers below are the reading-off.

## Asking the field

```go
import "github.com/snowbldr/fluent-sdfx/validate"

d := validate.At(part, v3.XYZ(10, 0, 5))          // signed distance; negative inside
ok := validate.Inside(part, validate.Cyl(12.5, 120, 3.5))   // r, degrees, z
```

`At` takes the fluent vector type, so a test never needs the sdfx vector package. One caution the whole package is built around: after booleans, twists, scale-extrudes and `Correct`, the field is a *bound*, not a distance. The sign is always right; the magnitude can understate. `validate.At(cyl.Cut(bore), top)` can read 3.0 at a point 2.0 from the material. So use `At` for inside-or-outside, and use the surface-finding helpers when you need a number.

## Named probes

The most-used pattern by far. A table of points that encode what the person actually asked for, each named so a failure reads as a sentence.

```go
func TestBossShortened(t *testing.T) {
    part := build()
    validate.RequireProbes(t, part, []validate.Probe{
        validate.In("boss1-height-kept", v3.XYZ(0, 12.5, 4.8)),
        validate.Out("boss1-tip-removed", v3.XYZ(2.75, 12.5, 3.5)),
        {Name: "wall-at-least-1mm", P: v3.XYZ(9, 0, 5), Inside: true, Margin: 1.0},
    })
}
```

```
[ok] boss1-height-kept (0.000, 12.500, 4.800) want=inside got=inside sdf=-0.200
[FAIL] boss1-tip-removed (2.750, 12.500, 3.500) want=outside got=inside sdf=-0.250
2/3 probes passed
```

`Margin` asks for at least that much field on the expected side, which is how you guard "this wall is at least 1 mm thick here". For an assembly, `NamedProbe` adds a `Solid` key and `RequireProbesIn` takes a map of parts, so one table checks several bodies. `Probes` and `ProbesIn` are the non-test forms; they return results you can print from a throwaway `main`.

Every incident in the source project became a probe that stayed. That growing table is what made one-turn fixes possible: a bug cannot come back silently.

## Finding surfaces and measuring through material

```go
r, _ := validate.SurfaceR(part, 45, 3, 0, 30)                    // radius of the surface at 45°, z = 3, marching out from r = 0
z, _ := validate.SurfaceZ(part, 4, 4, 20, -20)                   // top surface height at (4, 4)
d, _ := validate.SurfaceAlong(part, from, dir, 30, 0.1)          // first surface along any ray, bisected to 1e-6
w, _ := validate.WallThickness(part, v3.XYZ(9, 0, 5), v3.X(1), 20) // material run through a point
g, _ := validate.Opening(part, v3.XYZ(0, 0, 7), v3.Z(1), 20)      // clear span through a void
validate.RequireNear(t, "bore radius", r, 8.0, 0.05)
```

These march until the sign changes and then bisect, so they are exact regardless of how loose the field's magnitude is. `MinWall` scans a region on a grid and reports the thinnest material run it finds along any axis: the screening check for an accidental 0.1 mm skin, which in the record was found by a slicer after the part was printed.

## Two parts: clearance, contact, congruence

Sample one part's surface, read the other part's field there. Never compare two fields to each other; compare a surface to a field.

```go
g := validate.Clearance(plug, socket, 8)      // g.Min, g.Median, g.P95, g.Max, g.MinAt
validate.RequireClearance(t, plug, socket, 8, 0.2)     // every point at least 0.2 clear
validate.RequireNoContact(t, lid, body, 8)
validate.RequireInterference(t, flipped, socket, 8, 0.5)  // a rejection test: it must NOT fit
validate.RequireCongruent(t, rebuilt, original, 8, 0.02)  // same surface, both ways
```

`ClearanceIn` restricts the sample to a box, so you can ask about one seat on a large part without meshing the whole thing finely. `DeviationSTL` reads an STL you already wrote and reports how far its vertices sit from the field it came from: the check for decimation damage.

## Motion

A static fit check passes and the part still jams, because the collision happens halfway through the unscrewing. Sweep the moving body's surface through the real motion.

```go
poses := validate.ScrewPoses(1, 360, 6)        // 1° steps, one full turn counter-clockwise, rising 6 mm
validate.RequireSweepClear(t, lug, body, poses, 8, 0.3)

lifts := validate.LinearPoses(v3.Z(1), 0, 12, 0.5)
validate.RequireSweepClear(t, key, channel, lifts, 8, 0.3)

// Anti-flip: the wrong orientation must collide.
validate.RequireSweepBlocked(t, key.RotateX(180), channel, lifts, 8, 0.5)

// Bore play × screw motion, every combination.
all := validate.Compose(validate.LinearPoses(v3.X(1), -0.3, 0.3, 0.3), poses)
```

`Sweep` returns the worst clearance, the pose index and the world point where it happened, so the failure tells you where to look. On the benchmark's unscrew task, the relieved body sweeps clear at 0.27 mm against a 0.3 design, and the unrelieved one collides by a full millimetre at pose 310 of 361, which is the last rib met after the lug has risen 5 mm. No static check finds that.

## Handedness

A render will not reliably tell you which way a thread winds, and a mirrored mould stayed mirrored for weeks in the record. `HelixHand` decides by screw invariance: the field of a true helix is unchanged when you advance along the axis and rotate by the matching angle.

```go
validate.RequireHelixHand(t, thread, 1.25, 12, validate.RightHand)
```

Right-hand means rising along +Z while turning counter-clockwise seen from +Z, which is a normal screw thread. `shape.SweepHelix` produces a right-hand helix; `MirrorXZ` or `MirrorYZ` flips it.

## Printability, in the frame you will print

```go
pf := part.RotateX(90)                                   // into the print frame: min Z is the plate
rep := validate.Overhangs(pf, 8, validate.OverhangOptions{MaxDeg: 45})
validate.RequireOverhangArea(t, pf, 8, validate.OverhangOptions{MaxDeg: 45}, 2.0)
validate.RequireStandsOn(t, pf, 0, 0.01)
```

The audit is the one from the source project, exactly: a downward face fails when its angle from vertical exceeds the threshold, unless it is within `BedTol` of the plate or there is material (or the plate) within `SupportDepth` below its centroid or below four lateral offsets. The result is an *area* in mm², so a few sliver triangles along a chamfer edge do not fail a part while a cantilevered spoke does. `Threshold` takes a function of height so one band (a helical flank that prints at 60 degrees by construction) can be exempted while the rest is held at 45. Set the threshold to what your slicer paints, not to what you remember; in the record, an audit at 50 degrees against a slicer at 30 cost four rounds of "nothing changed".

`RequireMaxOverhang` and `OverhangArea` remain for the simple mesh-level version.

## Regression guards

```go
validate.RequireIdentical(t, before, after, validate.GridPoints(before.Bounds().Box, 1), 1e-6)
validate.RequireUnchangedOutside(t, before, after, editedRegion, 8, 0.01)
box, n := validate.ChangedRegion(before, after, 8, 0.05)
validate.RequireBounds(t, part, want, 0.01)
```

`RequireIdentical` is the refactor guard: a change that is supposed to be a no-op is proven to be one pointwise. `ChangedRegion` is the locality guard: it returns the bounding box of every surface point where the two solids differ, so "I only touched the boss" becomes a checked claim.

## Mesh-level checks

```go
validate.RequireWatertight(t, part, 5)
validate.RequireOneBody(t, part, 5)          // connected components == 1
validate.RequireVolumeNear(t, part, 5, 12500, 0.02)
st := validate.Of(part, 5)                   // Triangles, SurfaceArea, Volume, BoundaryEdges, Watertight, Bounds, OverhangArea
```

Watertightness catches holes from a field discontinuity or a bad custom SDF. Component count catches floating bodies, which a render hides and a slicer prints as loose parts. Volume pins a finished part against a future sign flip or a dropped boolean argument.

## Physical iteration

```go
coupon := validate.Coupon(part, keepBox, 8)                       // s ∩ keep, dropped to z = 0; panics if not one body
coupon = validate.Dimples(coupon, 3, at, normal, along, 1.5, 0.4, 3.0)   // 3 dimples = the third variant
validate.FitCheck("out", "assembly", parts, origin, v3.X(-1), 8)  // assembly.stl and a cutaway, plus pairwise gaps
fmt.Println(validate.Manifest(parts, 4))                          // name, bbox, size, volume, bodies
```

A coupon is a small piece cut from production geometry, oriented as printed, one variable per piece, labelled with dimples because text does not print at 4 mm. In the record they turned a two-hour print into a fifteen-minute one and answered one question each. Cutting a coupon through a hole yields two loose pieces, which is why `Coupon` checks the body count.

## Looking at it

```
go run github.com/snowbldr/fluent-sdfx/tools/inspect part.stl
```

writes six PNGs (four isometric, top, bottom) with f3d, opaque and unlabelled, so a person or a model can read them. For a section, cut with `CutPlane` and inspect that STL. Look before you probe, and probe what you saw.

## Density and cost

`cellsPerMM` here means what it means for `STL`. Field probes cost nothing; surface-sampling checks (`Clearance`, `Sweep`, `Congruent`, `Overhangs`, `Components`) mesh once at that density and then evaluate in parallel. Five to eight cells per millimetre catches everything in this page in a second or two; go finer only for a feature under a millimetre.
