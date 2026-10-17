# bench results

The run directories under `_runs/` are gitignored and live on one machine.
This file is the record that survives. Each entry names the runs it came
from and the model ids `run.json` recorded, where the harness recorded them
(from 2026-10-16; earlier runs are attributed by name and commit message).

## 2026-10-17: does a mixed-model pair beat a single model?

Hard suite, 11 tasks (8 hard + 3 clarify), tools=full, architect stage
on. The 2×2 of architect model × builder model:

| architect → builder | passed | asked first | turns | cost | runs |
|---|---|---|---|---|---|
| fable → fable | 11/11 | 4/4 | 2.36 | $19.48 | hard-fable-full-arch + fix-hard-arch + fix2 + fix3 + fix3b (2026-09-17) |
| opus → fable  | 11/11 | 4/4 | 2.36 | $11.95 | mix-opus2fable-hard-arch |
| fable → opus  | 11/11 | 4/4 | 2.36 | $11.09 | mix-fable2opus-hard-arch |
| opus → opus   | 11/11 | 4/4 | 2.36 | $4.76  | mono-opus-hard-arch |

Resolved ids: `claude-opus-5-5`, `claude-fable-5-1`, CLI 2.1.302. Every
cell passed every task, asked on every clarify task, and used the same
number of turns on every task. Nothing a mixed pair did, a single model
did not also do.

The per-stage cost is the finding. Same spec stage, same 15 architect
turns; same build stage, same 11 build turns:

| stage | on opus | on fable |
|---|---|---|
| architect | $1.48–$2.13 | $7.74 |
| builder   | $3.28–$3.35 | $9.82 |

Fable costs 2.5–3× Opus for the same work at the same result. With the
architect stage in front, the hypothesis that Opus gets most of the way
and Fable finishes is not supported here: Opus finishes.

Two caveats. The fable → fable cost is a month older than the others, so
a same-day fable → fable control would make the cost column clean (the
other three were run within two hours of each other). And the suite is
saturated: four cells at 11/11 cannot rank anything, so the next question
needs harder tasks before it needs more runs.

Where Fable did earn its place is the record below: without the
architect stage, Fable with a workspace passed 23/23 where Opus passed
21/23 (2026-09-15), and 9/11 on the hard suite where no Opus plain run
exists. The architect stage is what closes that gap for Opus.

## 2026-09-17: the architect stage (hard suite, Fable)

| condition | passed | asked first | turns | cost |
|---|---|---|---|---|
| plain     | 9/11  | 3/4 | 1.27 | $13.01 |
| architect | 11/11 | 4/4 | 2.36 | $19.48 |

Gained `gusset-triangular` and `three-unknowns-bracket`, both cases where
one good question was not enough. `_runs/compare-hard.md`.

## 2026-09-15: model × tools (23 tasks)

| condition | passed | asked first | turns | cost |
|---|---|---|---|---|
| opus, none  | 18/23 | 1/3 | 1.04 | $3.56 |
| fable, none | 18/23 | 1/3 | 1.04 | $6.56 |
| opus, full  | 21/23 | 2/3 | 1.13 | $21.49 |
| fable, full | 23/23 | 3/3 | 1.17 | $17.85 |

Opus lost `chamfer-under-tab` and `chamfer-print-orientation` with
tools. This suite stopped discriminating after this run, which is why
the hard suite exists. `_runs/compare-2026-09-15.md`.
