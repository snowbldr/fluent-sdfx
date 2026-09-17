# bench: description-to-geometry benchmark

Measures how well a model turns a human description into a fluent-sdfx part,
and how many turns it takes. Tasks are seeded from real incidents in a
two-month AI-driven CAD project (see `meta/ai-friendly-findings.md`), so each
one encodes a way the translation actually failed.

## Layout

```
bench/
  score/            scoring library: staged metrics, generated scorer program
  cmd/score         score one candidate package against a task
  cmd/selfcheck     validate every task (reference passes, asfailed fails)
  cmd/run           drive a model over tasks with claude -p, score, summarise
  tasks/<name>/
    task.json       metadata, thresholds, intent probes
    prompt.md       the description, in the user's own words
    answers.md      hidden spec a scripted user answers from (clarify tasks)
    start/          package start: the part before the change (modify tasks)
    reference/      package reference: the intended result
    asfailed/       package asfailed: what the model actually built in the incident
  runs/             gitignored: scorer scratch and run outputs
```

Every `start`, `reference`, `asfailed` and candidate package exports
`func Build() *solid.Solid`. Nothing else is required of them.

## Task modes

- **build**: the prompt describes a part from nothing. Candidate is aligned to
  the reference by bounding-box centre before comparison (`align: center`).
- **modify**: the prompt asks for a change to `start`. The model is given the
  start code and edits it. Frames are shared (`align: none`).
- **clarify** (`expects_question: true`, any mode): the prompt is
  underspecified on purpose, exactly as the user phrased it. The ideal first
  response is a question. `answers.md` is the hidden spec a scripted user
  answers from. The reference is what results once the answers are known.

## Scoring stages

1. **compiled**: `go vet` on the candidate package.
2. **built**: `Build()` returned a solid without panicking.
3. **bbox**: size and centre error against the reference.
4. **IoU**: both SDFs sampled on a grid over the union bounding box
   (`grid` cells along the longest axis). Also reports `missing`
   (reference volume the candidate lacks) and `extra` (candidate volume the
   reference lacks, i.e. unrequested material), both as fractions of the
   reference volume.
5. **max_dev**: approximate two-sided Hausdorff distance, from each solid's
   SDF sampled on the other's surface band. Catches small local errors that
   IoU cannot see (a 0.4 mm shortening, a missing chamfer).
6. **probes**: named points in the reference frame with an expected inside
   or outside state. These are the intent: they encode what the user was
   actually asking for, and they name the failure when they fail.
7. **mesh** (optional, `mesh_cells_per_mm > 0`): watertight, volume,
   connected-component count.

A candidate passes when IoU >= `iou_min`, max_dev <= `max_dev_mm`,
extra <= `extra_max`, every probe passes and, if enabled, the mesh is one
watertight body. Defaults: 0.97, 0.3 mm, 5 percent.

The runner adds per-task **turns to pass**, **asked a question first** (for
clarify tasks), cost and elapsed time.

## Authoring a task

1. Pick an incident. Write `prompt.md` in the user's voice. If the incident
   was a silent assumption, keep the prompt as vague as the user's and put
   the facts in `answers.md`; otherwise include enough numbers that the
   reference is the unique correct answer.
2. Write `reference/`. Keep parts small (under about 40 mm) so scoring takes
   seconds. Name the constants; comment the geometry.
3. Write `asfailed/`: the construction the model actually produced in the
   incident, or the most plausible wrong reading. This proves the metric
   discriminates.
4. Write at least eight probes. Every probe should either guard something the
   prompt asked for (and fail on `asfailed`) or guard something the prompt
   said not to change (and pass on both). Name them so a failure reads as a
   sentence: `boss1-height-kept-z4.8`.
5. Run `go run ./bench/cmd/selfcheck <name> -v`. The reference must pass with
   IoU 1.000 and every probe green; `asfailed` must fail. Tune `grid` and
   thresholds until the failure is caught by more than one stage if you can.

What the numbers mean:

- `max_dev` is a symmetric Hausdorff distance between the two voxelised
  surfaces, so its resolution is one grid cell and a reference scored
  against itself reports exactly 0. It depends only on occupancy: SDFs that
  are not true distances (twists, scale extrudes, loose CSG bounds, internal
  coincident-face membranes) cannot inflate it. A sub-cell error (a 0.15 mm
  clearance) is invisible to it; guard those with probes, which evaluate the
  exact SDF at a point, and with `extra`.
- `max_dev_sdf` is the older SDF-magnitude measure, kept for information. It
  sees sub-cell shifts but is inflated by non-Euclidean SDFs. Not a pass
  criterion.
- Volumetric stages cannot see small local features on a large part (a 1 mm
  chamfer on a 4000 mm^3 part is 0.3 percent). Probes and `max_dev` are what
  catch those; tune `iou_min` per task rather than chasing it.
- The sample lattice is offset by an irrational fraction of a cell so that
  cell centres never sit exactly on a face plane, where a coincident-face
  SDF reads exactly 0.

Categories to cover (from the findings): vocabulary, assumption, orientation,
anchor, api-trap, printability, verification, over-engineering, multi-ask.

## Running

```bash
go run ./bench/cmd/selfcheck                      # validate all tasks
go run ./bench/cmd/score -task bench/tasks/teardrop-hole -cand path/to/candidate
go run ./bench/cmd/run -model opus -tools none    # drive a model over all tasks
go run ./bench/cmd/run -model fable -tools full -tasks pill,teardrop
```

See `cmd/run` for the driver flags: model, tool condition (none: code in a
fenced block, no tools; full: the model may read, write, run Go, render and
probe inside its workspace), max turns, parallel jobs, and which reference
material (llms.txt or a skill directory) is appended to the system prompt.
