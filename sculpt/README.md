# sculpt — sculpting-app verbs on fluent-sdfx solids (experimental)

The operations a person reaches for in a 3D sculpting app, as Go functions:
pull a point of the surface out, smooth every edge at a radius, bucket-fill
an enclosed void, drop a piece that is not attached. These are the verbs for
organic and decorative parts, where there is no bore or flange to name.

This folder is self-contained. Nothing in `solid`, `shape` or `validate`
changed, it is not in `docs/llms.txt`, and there are no API promises yet.
The tests here are the realistic cases the verbs are judged on.

## What is here

| Verb | Kind | What it does |
|---|---|---|
| `Pull(s, at, by, radius)` | closed-form warp | Grabs the surface at `at` and drags it by `by`; the point lands exactly, the surface within `radius` follows with a smooth falloff, nothing farther moves. Panics past 0.65 × radius (the surface would fold). |
| `Sample(s, cells)` | grid | Samples a solid onto a cubic lattice, `cells` along its longest side, three cells of air around it. |
| `Grid.Redistance()` | grid | Recomputes every sample as the true signed distance to the surface the samples describe. Error is a fraction of a cell, deep interior included. |
| `Grid.Fill(seed)` | grid | Bucket fill: the air connected to `seed` becomes material. Panics if the region is open to the outside; `Enclosed(seed)` asks without panicking. |
| `Grid.Components()` | grid | The connected bodies of material, largest first, with cell count, volume, bounds and a point surely inside each. |
| `Grid.Keep(seed)` | grid | Keeps only the body containing `seed`. |
| `Grid.Smooth(r)` | grid | Rounds every convex edge and fills every concave crease at radius r: shrink/grow then grow/shrink, redistanced between steps. Removes features thinner than 2r, by definition. |
| `Grid.Solid()` | grid | Back to a `*solid.Solid` by trilinear interpolation, for Shell, booleans, export. |

The usual chain: build with primitives and `Pull`, `Sample`, grid verbs,
`Solid()`, then finish with `Shell`, booleans and `STL`.

## The one thing to understand

A `Pull` (and a twist, a smooth blend, a difference) leaves the field a
*bound* on the distance, not the distance. The surface is right; the
magnitudes are not. `Shell`, `Offset`, clearance probes and smooth blends
read magnitudes, so on a pulled solid they are wrong, and not subtly:
`TestFigureShellNeedsRedistance` measures a 1.6 mm shell coming out at 5.7
to 9.0 mm on the raw pulled figure and at 1.60 after `Redistance`.
`Redistance` is what makes chaining sculpt operations honest.

## The test cases

- **`ornament_test.go`** — a hollow bauble with a hanging loop, going into a
  plaster mould where the pattern must be solid. `Fill` from the centre:
  the void fills, the outer surface does not move, the loop and its hole
  survive, the volume matches a solid ball. A vented bauble is not enclosed
  and `Fill` refuses.
- **`figure_test.go`** — a body-and-head figure with a nose pulled from the
  face and a stem from the top. The grabbed points land exactly, nothing
  outside the brushes moves, a pull that would fold is refused. Then the
  shell test above.
- **`cluster_test.go`** — a form built from blobs with a sharp crease, a
  block with sharp edges and a bead that floats. `Components` finds the
  bead, `Keep` drops it, `Smooth(2)` fills the crease to the depth a 2 mm
  ball leaves (checked against the exact rolling-ball geometry) and rounds
  the block edge to a 2 mm fillet (checked against sdfx's rounded box).
  A stem thinner than 4 mm does not survive `Smooth(2)`, and the test says so.
- **`redistance_test.go`** — the primitive on its own: a sphere's field is
  recovered everywhere from scaled samples to within 0.2 cells; a twisted
  column's shell, 2.28 mm on sdfx's raw field, measures 1.99 after
  redistance; smoothing a sphere leaves it alone.

Run them with `go test ./sculpt/ -v`; the log lines carry the measured numbers.

## How Redistance works

Seeds are where the surface is: along every lattice edge whose ends differ
in sign, the zero crossing by linear interpolation; and for every corner
within a cell of the other side, the point found by projecting along the
field's gradient (which is what seeds a sharp edge, where no lattice edge
crosses). A separable Euclidean distance transform (Felzenszwalb and
Huttenlocher) finds each corner's nearest seed, two neighbour-comparison
sweeps correct its choice, and the distance to that seed's surface point is
the new sample. Signs are kept from the input.

## Known limits

- Everything grid-side is accurate to a fraction of a cell and loses
  features smaller than a cell. Choose `cells` for the smallest feature.
- The refinement sweeps in `Redistance` are single-threaded; `Smooth` runs
  five redistances. A 128-cell grid smooths in about two seconds.
- `Fill` fills enclosed regions only. A pocket open on one side is filled
  by capping it (union a plate), filling, then cutting the cap away.
- `Pull` is one brush. Symmetry, pinch, inflate and a region-limited smooth
  are the next candidates, to be added when a real piece asks for them.
