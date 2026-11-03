# history: a part as data, run by reflection (experimental)

A history is a list of steps in YAML or JSON. Each step is one fluent-sdfx
function or method. The runtime evaluates it by reflection, caches every
step, reports a status per step, and exports the same part as Go. This is
the runtime from Phase 3 of `meta/sculpt-cad-build-plan.md`; a UI or an
agent edits the file, and the library does the rest.

```yaml
params:
  headR: 8
  wall: 1.6
steps:
  - name: body
    op: solid.Sphere
    args: {radius: 12}
  - name: head
    op: solid.Sphere
    args: {radius: $headR}
    then:
      - {op: Translate, args: {v: [0, 0, 14]}}
  - name: figure
    op: solid.SmoothUnion
    args:
      min: {op: solid.RoundMin, args: {k: 3}}
      solids: [body, head]
  - name: NoseTip
    anchor:
      on: figure
      ray: {from: [0, -40, 14], dir: [0, 1, 0]}
  - name: nose
    op: sculpt.Pull
    args: {at: NoseTip, by: [0, -4, 0], radius: 7}
    by: agent
    from: chat/123
  - name: hollow
    op: Shell
    args: {thickness: $wall}
```

`testdata/` has more: `figure.yaml` (anchors, a mirrored step, a grid),
`ornament.yaml` (a fill seeded by an anchor), `cluster.yaml` (keep and
smooth), `kitchen.yaml` (every way of writing args).

## Steps

| Key | Meaning |
|---|---|
| `name` | What later steps call it. Letters, digits and `_`, unique. Optional. |
| `op` | A function (`solid.Sphere`, `sculpt.Pull`, a registered `my.Bracket`) or a method on the input (`Translate`, `Shell`, `Union`). |
| `on` | The input: the step a method is called on, or the first argument of a function. Omitted, it is the previous step, for a method, and for a function whose first parameter is a solid that `args` does not give. |
| `args` | By name, a map (`{radius: 12}`), or by position, a list (`[12]`). |
| `then` | Methods called on the result in order: `[{op: TranslateZ, args: [14]}]`. |
| `mirror` | `x`, `y` or `z`: run the step again on its own result with every vector argument mirrored. One step pulls both cheeks. |
| `off` | `true` skips the step. Steps that used it get its input, as if it had never been there. |
| `by` | Who made it: `user` (the default) or `agent`. |
| `from` | An opaque reference to what made it (`chat/123`). Stored and returned exactly, never read. |
| `anchor` | Instead of `op`: this step is a named point on a surface (below). |

## Values in args

- A number, a bool, a string (where the parameter is a string).
- `[x, y, z]` for a `v3.Vec`, `[x, y]` for a `v2.Vec`; a list for a slice
  or a variadic parameter.
- `$name`: a param.
- `body`: another step's value. `head.Top`: a member of it, a method with
  no arguments or a field (`head.Top`, `NoseTip.Normal`). A body anchor
  goes where a point is wanted as its point.
- `{op: solid.RoundMin, args: {k: 3}}`: an inline call, for blends and
  other values that are not steps.
- `{anchor: {...}}`: an inline anchor; its surface defaults to the step's
  input.
- A map for a struct parameter. A function whose only parameter is a struct
  (`obj.Bolt`, or a registered `func(p Params)`) takes the fields straight
  in `args`.

Conversions: an anchor or a body anchor where a point is wanted gives its
point; a grid where a solid is wanted gives its `Solid()`, and a method a
grid does not have (Shell, Translate) is called on its solid.

## Anchors

A point found on a surface every time the history runs, so a later step
lands on the face even after an earlier dimension changes.

```yaml
- name: NoseTip
  anchor:
    on: figure                                   # default: the previous step
    ray: {from: [0, -40, 14], dir: [0, 1, 0]}    # or toward: [x, y, z]
- name: VoidSeed
  anchor: {on: bauble, nearest: [0, 0, 0], offset: 2}
```

`ray` is the first crossing of the surface along the ray. `nearest` is the
nearest surface point. `offset` moves the point along the normal: positive
out of the material, negative into it, so a fill seed 2 mm into a void is
the nearest inner wall point offset by 2. The value has `Pos` and `Normal`;
using the step where a point is wanted gives `Pos`. A ray that misses is
the anchor step's error, and the steps that use it are blocked. The same
finding is available in Go as package `anchor`.

## Statuses

`Runtime.Eval` never fails as a whole. Each `StepResult` has `status`:

- `ok`;
- `error`, with a message that says what went wrong and usually how to fix
  it ("keep the pull under 4.55 mm (0.65 × radius) or widen the radius",
  "has no parameter radious; did you mean radius?");
- `blocked`: something it needs failed, is off with nothing to pass on, or
  does not exist ("needs step body, which does not exist");
- `off`.

The steps that could be evaluated were, and `Result.Bodies` is what shows:
every solid that no evaluated step went on to use. Deleting the body of a
figure blocks the blend and everything after it, and the head still shows.

## Go API

| Call | What it does |
|---|---|
| `Load`, `Parse`, `(*History).Save`, `Marshal` | Read and write; unknown keys are an error; the header comment survives. |
| `Insert`, `Set`, `Remove`, `Move`, `SetOff`, `Index`, `Clone` | Edits; one that would leave the history invalid changes nothing. |
| `NewRuntime(Options)`, `(*Runtime).Eval` | Evaluate with a cache. `Options.BakeEvery` samples the form onto a grid after that many pulls in a chain. Two runtimes are two caches (preview and final). |
| `Result.Steps`, `Bodies`, `Solid`, `SolidOf`, `Value`, `Err` | Outcomes and geometry. A `StepResult` marshals to JSON for an agent. |
| `Result.Pick(from, dir)` | The nearest body hit by a ray, against the field, not a mesh. |
| `Result.Anchors()` | Every anchor a new step could use: anchor steps with normals, and the box anchors of every named solid. |
| `Palette()` | Every op with its doc, parameters, types and a widget hint (point, offset, direction, size, length, angle, count, blend, step, struct). |
| `Register(name, fn, doc, params...)` | Make a Go function an op. |
| `ExportGo(h, ExportOptions)` | The Go program that builds the same part; provenance is a comment above its step. |

The command line is `go run ./tools/history check|build|export|palette|anchors`.
A program that registers its own functions gets the same commands from
`history/cli.Run`.

## What is cached, and what a bake is

A step's cache key is its definition (without `name`, `by`, `from`, `off`),
the keys of every step it uses and the values of the params it reads.
Changing a step recomputes it and what follows; changing only its
provenance recomputes nothing. With `BakeEvery`, a chain of that many
pulls is sampled onto a grid (not redistanced, so the field keeps its
values) and later steps read the grid. The surface moves by a fraction of
a cell; the Go export has no bakes.

## Not yet

- Feature anchors (`boss.Axis`, `plate.TopFace` from the constructor that
  made them, plan 2.3). Box anchors and surface anchors cover the cases so
  far.
- Grid caches on disk (plan 1.3); bakes live in memory.
- Comments inside the file other than the header are not kept on save.
