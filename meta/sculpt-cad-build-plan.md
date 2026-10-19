# SDF CAD — framework build plan

The framework changes that have to exist before a UI can be built on
fluent-sdfx: a sculpting kernel, anchors, a data history that runs by
reflection, materials, and render backends. The UI is a follow-up plan and
is only mentioned here where the framework has to expose something for it.

Written 2026-10-18 from the conversation that started with "bucket fill
and smooth corners as Go functions". The sculpt package in `sculpt/` is
the first piece and already exists (uncommitted at time of writing).

All names in this document are placeholders. Package names, verb names and
YAML keys get decided by the owner before anything is public.

## What we are building

A history of operations, stored as data, evaluated by the library, shown
and edited by a UI, exported as Go. Every operation is a fluent-sdfx
function. Every click will be one step. Deleting a step and re-running
makes it never have happened. The form is always a signed distance field,
so booleans never fail, the mesh is always watertight, and the validate
package can report walls and overhangs at any point.

Same architecture as Dreams (edit list evaluated into volume blocks) and
Houdini (node graph that recooks), with three things neither has: a
history that exports to a real program, print checks in the loop, and an
open library underneath that an AI can also drive.

## Principles that hold through every phase

- **The kernel never changes to serve the UI.** The UI consumes the
  history runtime; the runtime consumes the library. Anything the UI
  needs that the library lacks is a library feature with its own tests
  and its own reason to exist.
- **Public API changes go through the bench.** The AI-friendly findings
  rule stands: measure before changing `solid` or `shape`. New packages
  are free until they are documented in `docs/llms.txt`.
- **No control flow and no arithmetic in the history.** Named params,
  references to steps and anchors, literal numbers. Anything procedural
  is a registered Go function. The moment the YAML grows an `if`, it is a
  worse Go.
- **Every phase ends with a scenario test** in the style of `sculpt/`: a
  realistic piece, numbers asserted against geometry computed another
  way, log lines carrying the measurements.
- **Bake is a cache, never a flatten.** The history is always the truth.
  A grid is recomputable from a prefix of the history and is discarded
  when that prefix changes.

## Phase 1 — Sculpt kernel (`sculpt/`, exists; extend)

What exists: `Pull`, `Grid` with `Sample`, `Redistance`, `Fill`,
`Enclosed`, `Components`, `Keep`, `Smooth`, `Volume`, `Solid`. Redistance
is seeds from axis crossings plus gradient projection (including diagonal
contacts, which is what seeds a sharp edge), a separable Euclidean
distance transform with argmin, and two neighbour-refinement sweeps.
Sphere error 0.19 cells max, 0.017 mean. Three scenario tests pass.

### 1.1 Performance floor

- Parallelise the refinement sweeps in `Redistance` (wavefront over
  z-slabs or a Jacobi variant; must stay deterministic). Target: a
  128-cell `Smooth` under 0.5 s on a laptop. Today it is about 2 s.
- Add a preview-resolution path: `Sample` at 48 to 64 cells runs under
  50 ms for histories of tens of ops. Measure and record in a table in
  the package README.
- Benchmark file in `sculpt/` with the three scenario forms at 64, 96,
  128, 192 cells so regressions are visible.

### 1.2 The brush family

Add verbs only when a test piece asks for one. The candidates, in the
order the comparison with Dreams, Modeler and MagicaCSG suggests:

| Verb | Kind | Note |
|---|---|---|
| Push | closed form | `Pull` with a negative `by`; probably just documentation |
| Inflate / deflate in a region | grid | signed offset limited to a falloff region: `shift` masked by a weight |
| Local smooth | grid | `Smooth` blended with the original by a falloff weight; the mask machinery that masks and layers in Modeler use |
| Flatten | grid | Dreams' flatten: set the field toward a plane sampled at brush centre, blended by falloff |
| Pinch / spread | closed form | radial warp toward or away from an axis through the brush centre |
| Symmetry | history-level | not a verb: a step flag that runs the step twice under a mirror. Belongs to Phase 3 |

Each verb: doc comment first in the owner's words, then code, then a
scenario test with analytic expectations (the rolling-ball depth for
Smooth is the model).

### 1.3 Grid persistence

Grids need to survive a process restart so caches do (Phase 3.5) and so
a baked form can be a step input. Format: a small header plus raw
float32 corner samples, gzip optional. Load and save round-trip test.
Not a public interchange format.

### 1.4 Acceptance

- All current sculpt tests pass; the benchmark table exists.
- One new verb (local smooth or flatten) landed through the full loop:
  words, code, scenario test.

## Phase 2 — Anchors

Anchors are what make a reflow land in the right place rather than at the
same coordinates. Three kinds.

### 2.1 Body anchors (exist)

The 27 bounding-box anchors and the placement verbs. No change, but the
history runtime (Phase 3) must be able to name them: `head.Top`,
`body.BottomFrontLeft`. Also the rename the findings asked for:
unambiguous aliases for `FrontAt`/`BackAt`. Goes through the bench.

### 2.2 Surface anchors (new)

A point on a surface, remembered by how it was found, resolved at run
time against the current surface:

- **Ray anchor**: from a point (a body anchor, a literal, or the camera)
  toward a point or along a direction; `Raycast` finds the hit. Misses
  are a step error, never a silent fallback.
- **Nearest anchor**: the surface point nearest a given point.

Both return position and normal (`Normal` exists). This is how a pull
remembers "the nose", and how a fill remembers "inside the body" when
the seed must be air: nearest-surface plus an inward offset along the
normal.

### 2.3 Feature anchors (new, the one the findings asked for)

Constructors know their geometry. A `Cylinder` has an axis and two end
centres; a `Box` has six face centres; a `Hole` has a centre and an axis;
a `Torus` has its ring. Expose these as named anchors on the step that
made them, so later steps say `boss.Axis` or `plate.TopFace`, not a
bounding-box corner of a composite that happens to coincide today.

Implementation: a small table keyed by constructor, giving the anchors
as functions of the constructor's parameters and the step's transform.
The transform is the composition of the method chain that followed
(Translate, Rotate), which the runtime already has. Derivable in the same
parse that builds the palette.

Boolean results do not get feature anchors; they get body anchors plus
whatever surface anchors the user drops. That is the honest answer and
also what parametric CAD does.

### 2.4 User anchors

A named anchor defined by any of the above and stored in the history.
"NoseTip" is then a word later steps, the UI and an AI can all use.

### 2.5 Acceptance

- A scenario where an early dimension changes (head radius 8 to 10) and a
  pull anchored by ray still lands on the face, a stem anchored to
  `head.Top` still sits on top, and a boss anchored to `plate.TopFace`
  still sits flush. Asserted by probes after the reflow.
- A scenario where a ray anchor misses and the step reports it.

## Phase 3 — History runtime

The DSL. A package that parses a history, evaluates it against the
library by reflection, caches, reports per-step errors, and exports Go.

### 3.1 Schema

```yaml
params:
  wall: 1.6
steps:
  - name: body
    op: solid.Sphere
    args: {radius: 12}
  - name: head
    op: solid.Sphere
    args: {radius: 8}
    then:
      - {op: Translate, args: {v: [0, 0, 14]}}
  - name: figure
    op: solid.SmoothUnion
    args: {min: {op: solid.RoundMin, args: {k: 3}}, solids: [body, head]}
  - name: nose
    op: sculpt.Pull
    args:
      s: figure
      at: {anchor: {ray: {from: head.Front, toward: head.Center}}}
      by: [0, -4, 0]
      radius: 7
    mirror: x          # symmetry: the step runs twice, mirrored
  - name: hollow
    op: Shell
    on: nose
    args: {thickness: $wall}
anchors:
  NoseTip: {nearest: [0, -12, 14], on: nose}
```

Rules: step names are the variables; `on` is the receiver for a method;
`then` is a method chain on the step's result; a nested `{op: ...}` is an
inline value; `$name` is a param; `[x, y, z]` is a vector; `step.Anchor`
is an anchor reference. No expressions.

Pin the schema with a handful of files before the runtime exists: a
method chain, a blend, a param, a registered function, an anchor, a
symmetry step. These files are the first tests.

### 3.2 Evaluator

- Methods resolve by `reflect.Value.MethodByName` on the step's value.
  Zero registration; new library methods appear automatically.
- Package-level functions resolve through one map, `"solid.Sphere":
  solid.Sphere`, about eighty entries. `tools/llms-check` fails the build
  when an exported constructor is missing from the map, the same way it
  fails on a missing doc comment.
- Argument conversion by the parameter's reflected type: number to
  `float64`/`int`, 3-list to `v3.Vec`, 2-list to `v2.Vec`, string to a
  step's value or an anchor, nested op to its result, list to a variadic
  slice. Interface parameters (`sdf.SDF2`) check assignability.
- Panics inside a step are recovered and attached to the step. Later
  steps that depend on a failed step are marked blocked, not run.
  Everything else still evaluates. This is Houdini's red node.
- Registry for user functions: `history.Register("my.Bracket",
  Bracket)`. Reflection reads the signature, so the UI palette and the
  argument conversion treat it like a built-in.

### 3.3 Incremental evaluation and the cache

- Each step has a content hash of its own definition plus the hashes of
  what it references. Unchanged hash, unchanged result.
- A step whose result is a grid (anything that went through `Sample`)
  caches the grid in memory and on disk (Phase 1.3) keyed by hash.
- Baking policy: the runtime inserts an implicit bake after every N
  closed-form sculpt steps or when evaluation time per sample exceeds a
  budget; the bake is a cache entry, never a step. Editing an earlier step
  invalidates the bake and recomputes. The user never sees a flatten.
- Two resolutions: preview (fast, for the UI) and final (for export).
  Same history, different `cells`.

### 3.4 Go export

Generate a `main.go` that builds the same form with the fluent API: one
statement per step, step names as variables, params as constants, anchors
resolved to the calls that find them, registered functions imported from
their module. It must compile and produce congruent geometry
(`validate.Congruent` against the runtime's result). This is what makes
the history a real artefact and what lets a person leave the UI for Go.

### 3.5 Palette description

For the UI follow-up, the runtime emits a JSON description of every
callable op: name, doc synopsis, parameters with type, name and the
widget hint (number, point, size, direction, selection, blend, enum).
Derived by reflection plus the parameter-name table for `v3.Vec` meaning
(size, point, axis). Same source as the API index; cannot drift.

### 3.6 Command line

`go run ./tools/<name> build part.yaml -o part.stl` and `export part.yaml
-o part.go`. No UI yet, but the whole pipeline runs headless, which is
also what the bench and an AI will call.

### 3.7 Acceptance

- The three sculpt scenario forms expressed as YAML produce geometry
  congruent to the Go tests within a cell.
- Deleting the nose step reflows to the plain figure; deleting the body
  step blocks everything downstream with clear errors and still renders
  the head.
- Changing `wall` in params changes every shell.
- The Go export compiles and is congruent.
- A registered function appears in the palette description and runs.

## Phase 4 — Materials

- A material field: a wrapper that evaluates distance and a material id
  together. Rules: union takes the nearer child's material; cut and
  intersect keep the surviving body's; smooth blends mix by the blend
  weight and record the weight; transforms pass through.
- A material channel on `Grid` next to the distance samples, so baked
  forms keep paint.
- Paint ops: assign a material to a step, or to a region (a ball, a
  box, a falloff) — the Dreams colour edit.
- Export: 3MF with per-triangle material and colour (the mesh writer
  needs a colour path; today it writes geometry only). Per-material STL
  split for printers that want one file per filament.
- Acceptance: a two-colour figure exports a 3MF that a slicer opens with
  both materials assigned; a blended seam has a soft boundary of the
  expected width.

## Phase 5 — Render backends

### 5.1 Turnaround renders (first, cheap, needed by the AI loop)

`tools/inspect` already does six views with f3d. Add a turnaround (N
frames around Z) and a contact sheet. This is how artistic pieces get
judged before any UI exists, and how an AI looks at its own work.

### 5.2 GLSL backend for the history

- A shader-side twin for each op in the palette: primitives, booleans,
  transforms, blends, `Pull` and the closed-form brushes. A few lines
  each. Ops without a twin fall back to a baked grid.
- Grids upload as 3D textures with trilinear filtering (WebGL2 and
  WebGPU both support them).
- The transpiler emits one distance function and one material function
  from a history. Incremental: ten ops first, bake the rest.
- Acceptance: for a history with only twinned ops, the GPU distance at a
  thousand random points matches the CPU within 1e-4.

### 5.3 Path tracer

- Browser: progressive accumulation in a fragment shader, reset on camera
  move, sphere tracing against the transpiled distance function. The
  Shadertoy pattern; this is MagicaCSG's rendering.
- Go: a CPU path tracer for headless final frames and docs figures.
  Parallel, slow, fine.

## Phase 6 — Sparse grid

The one real engineering project, and gated on need: start it when a
test piece needs detail beyond what a 256-cell dense grid gives.

- Tiles of 8³ or 16³ corners stored only within a narrow band of the
  surface plus a coarse inside/outside tile map (the OpenVDB shape).
- Same interface as `Grid`: `Sample`, `Redistance`, the sculpt verbs,
  `Solid`. Redistance becomes narrow-band: seeds as today, propagation
  limited to the band, far field from the tile map.
- Acceptance: the sphere redistance test at a resolution a dense grid
  cannot hold in memory; the figure scenario at 512 cells in under the
  memory of today's 128.

## Order and dependencies

```
Phase 1.1 perf ─┐
Phase 1.2 verbs ├──► Phase 3 runtime ──► Phase 4 materials ──► Phase 5.2 GLSL ──► 5.3 path tracer
Phase 1.3 grids ┘         ▲
Phase 2 anchors ──────────┘
Phase 5.1 turnarounds (any time, early)
Phase 6 sparse grid (when a piece needs it)
```

Rough sizes, in working days, single person, with the AI doing the
typing: Phase 1, five. Phase 2, five, plus a bench run. Phase 3, ten.
Phase 4, four. Phase 5.1, one; 5.2, ten for the core palette; 5.3, five.
Phase 6, fifteen. The UI follow-up plan starts after Phase 3.7 passes,
with 5.1 in hand so pieces can be judged by eye.

## What the UI will need from the framework (contract, not design)

So the follow-up plan has a fixed surface to build against:

- Load, save, and edit a history (add, change, delete, reorder steps).
- Evaluate at preview or final resolution; receive a mesh per named body
  and the per-step status (ok, error with message, blocked).
- The palette description (3.5) with widget hints.
- Anchor listing per step (body, feature, user) with world positions.
- Pick: a ray in, the hit step, point and normal out.
- Export: STL, 3MF, Go.
- Later: the transpiled GLSL for the history and grid textures.

## Not in this plan

- The UI itself, the WebSocket protocol, the web page, gizmos, snapping.
- UVs and textures. Colour is per material, not per surface point.
- Thin sheets and cloth. Give them thickness or model them elsewhere.
- Any change to `docs/llms.txt` or the bench until the verbs have been
  used on a real piece and the names are settled.

## Open decisions for the owner

1. Names: the package for the history runtime, the YAML keys, the brush
   verbs. Everything above is placeholder.
2. Where the history runtime lives: inside this module or a sibling
   module that imports it. A sibling keeps the library small; inside
   keeps one `go test ./...`.
3. Whether the Go export is a feature or the save format. This plan says
   feature, with YAML as the format, for the round-trip reason.
4. When to start the sparse grid.
