# Contributing to fluent-sdfx

Thanks for the interest. fluent-sdfx is a small project; the contribution process is correspondingly light.

## Quick start

```bash
git clone https://github.com/snowbldr/fluent-sdfx
cd fluent-sdfx
go build ./...
go test ./...
go vet ./...
```

CI runs the same three commands on every PR — if they're green locally, they'll be green there.

## What kinds of contributions land easily

- **Bug fixes** with a regression test. Geometric bugs in particular: a failing test that produces wrong output or crashes is the gold standard.
- **New examples** in `examples/` or new tutorial steps in `tutorial/`. Examples are graded on whether they show off something genuinely useful, not on their size.
- **Doc improvements** — typos, clearer wording, missing godoc comments on exported API.
- **`obj/` parametric helpers** — bolts, nuts, panels, brackets, etc. If it has a parameter struct and produces a `*solid.Solid`, it fits.
- **`layout/` patterns** — new helpers that return `[]v3.Vec` or `[]v2.Vec` for `.Multi(...)`.
- **`validate/` checks** — new mesh-quality assertions that work as `*testing.T` helpers.

## What's likely to bounce

- **API changes that aren't backed by a clear case** — the chainable surface is the project's main asset; breaking changes need to clear a high bar.
- **New rendering backends or output formats** without a paired use case (we're not trying to become a graphics library).
- **Refactors that don't improve readability OR performance** — both should be measurably better, with numbers attached for performance claims.
- **Code that adds a dependency** — anything outside the standard library + sdfx + the existing transitive deps needs strong justification.

If you're unsure whether something will land, open an issue first — that's usually 5 minutes of your time and saves rewriting.

## House style

- **Constructors panic on invalid input.** CAD-geometry errors are programming bugs, not runtime conditions. No `_, err :=` interruptions in the chain.
- **All angles in degrees** at the public API boundary. Internal conversions go through `units.DtoR`.
- **The recipe pattern** in cookbook/example code: bare primitives at the top, single fluent assembly expression at the bottom. See [`tutorial/21-cookbook-lantern/05-finial/main.go`](tutorial/21-cookbook-lantern/05-finial/main.go) for the canonical shape.
- **Anchor positioning preferred over `Translate(v3.Z(h/2))` math.** If a translate is relative to another part's geometry, reach for `On`, `OnTopOf`, `BottomAt`, `Inside`, etc.
- **One short comment max** per non-obvious decision. Don't write multi-paragraph docstrings; let the names carry the meaning.
- **`go vet ./...` and `gofmt -l .` clean** in every PR.

## Testing

- Library tests live next to their code: `solid/*_test.go`, `validate/*_test.go`, etc.
- Example/tutorial code is verified by `go build ./...` in CI — if your example doesn't compile, CI fails. Don't write tutorial code you can't run.
- For geometric correctness, prefer `validate.RequireWatertight` / `RequireVolumeNear` / `RequireMaxOverhang` over hand-rolled assertions where it makes sense.

## Cookbook / tutorial / docs flow

The docs site (`docs/`) is generated from the markdown under `docs/content/`. Tutorial code lives in `tutorial/` and is referenced from the docs via `<!-- src: ... -->` comments — there's a syncer that keeps them in sync. Don't edit the rendered code blocks in markdown directly; edit the tutorial source.

Visual references (positioning anchors, the gallery) are produced by code under `docs/visuals/` plus the per-tutorial-step renders pinned in `docs/images/`. New visuals need both a generator and a render committed.

## Filing issues

A useful issue includes:

- What you ran (a minimal `package main` if possible)
- What you expected
- What happened (stack trace if it panicked)
- Your Go version and OS

Geometric bug reports should ideally include the parameters that produce the bug — "the boolean produced a non-watertight mesh on this configuration" is much easier to fix than "booleans sometimes break."

## License

By contributing, you agree your contributions will be licensed under the [MIT License](LICENSE) that covers the rest of the project.
