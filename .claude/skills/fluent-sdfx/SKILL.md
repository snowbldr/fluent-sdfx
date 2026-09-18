---
name: fluent-sdfx
description: Build, modify or verify 3D-printable parts with the fluent-sdfx Go CAD library (github.com/snowbldr/fluent-sdfx, a fluent layer over sdfx). Use whenever a task involves writing or editing fluent-sdfx or sdfx code, turning a person's description of a physical object into an STL, checking a part fits or prints, or debugging why a modelled part came out wrong. Covers intake (what to ask before building), the reference to load, and the verification loop that catches errors before a print.
---

# fluent-sdfx

A part is a signed distance field: a function from any point to how far it is from the surface, negative inside. Everything below follows from that. You can ask the field anything, so you never have to guess whether what you built is what was asked for.

The record this skill comes from is a two-month project modelled entirely this way. Its failures were not API failures. They were facts the person left out and the model assumed, size words on curved features read on the wrong axis, a single word with two readings, and claims made from a render instead of a measurement. Each step here exists because skipping it cost a print.

## 1. Load the reference

Read the library's own agent reference before writing code, from the version the project actually uses:

```bash
D=$(go list -m -f '{{.Dir}}' github.com/snowbldr/fluent-sdfx 2>/dev/null || pwd)
cat "$D/docs/llms.txt"            # always: conventions, traps, verification
grep -n 'Name' "$D/docs/llms-api.txt"   # when you need a signature or want to know if something exists
```

`llms.txt` is short and every claim in it was probe-verified. Trust it over your memory of sdfx: Revolve is about Z, `FrontAt` pins the minimum-Y face, `ExtrudeRounded` grows the footprint, `Shell` straddles the surface, `Intersect(a, b)` unions its arguments first, `Scale` is about the origin.

## 2. Intake: write the spec before the geometry

Follow `architect.md` in this skill's directory. In short: split the asks and carry every one; name the standard form for anything described by analogy (buttress thread, teardrop hole, corner gusset, labyrinth seal); translate every frame word the person used ("bottom", "front", "long", "clockwise") into a model axis and a sign, and show the translation; state what is frozen; list what you are deliberately not adding.

Then list the unknowns. Two kinds count. A fact that is absent: manufacturing orientation, nozzle or wall minimum, which surfaces are visible, a bought part's package, a tolerance. And a word with two readings that produce different geometry, however strongly you prefer one. If there are any, ask them all in one message and wait. Do not build on a guess when the guess changes the geometry.

Write the spec down, in the reply or in a comment at the top of the file. A spec you can point at is one the person can correct before the print.

## 3. Build

One expression per part, ingredients at the top as bare primitives, the fluent chain at the bottom. Cutting tools 1 to 2 mm longer than the body. Nothing added that was not asked for: no taper, fillet, clearance, lead-in or vent "to be helpful". Those are the features that fail.

## 4. Verify before you claim

A render is not evidence; a measurement is. Do these in a `_test.go` or a throwaway `main` using the `validate` package, and read the report:

- **Probe the intent.** For each thing the person asked for, a named `validate.Probe` at a point that is only inside (or only outside) if the ask was met, plus probes on what must not have changed. `validate.RequireProbes` names the failure.
- **Measure, do not derive.** A radius, a wall, a gap: `SurfaceR`, `WallThickness`, `Opening`, `MinWall`. Handedness: `RequireHelixHand`. A number you computed is a hypothesis; a number you measured is a fact.
- **Anything that mates:** `RequireClearance` between the two bodies at the design gap. Anything keyed: `RequireSweepBlocked` on the flipped assembly.
- **Anything that moves:** `RequireSweepClear` through the real motion (`ScrewPoses`, `LinearPoses`). A static fit check passes and the part still jams.
- **Anything printed:** rotate into the print frame, then `RequireOverhangArea` at the angle the person's slicer paints, and `RequireOneBody`.
- **A change to a working part:** `RequireUnchangedOutside` the edited region, and `RequireIdentical` for a refactor.
- **Look, too.** `go run github.com/snowbldr/fluent-sdfx/tools/inspect part.stl` gives six views; `CutPlane` gives a section. Probes verify what you thought to check; renders catch what you did not.

Turn every bug you find into a probe that stays. That is what makes the next change safe.

## 5. Deliver

Report what you measured, not what you believe: the probe table's pass count, the clearance minimum, the sweep minimum, the overhang area. State the assumptions you made and the frame every direction word is in. If a render is involved, send the six views, and say which artifact each shows.
