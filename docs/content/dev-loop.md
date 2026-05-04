# Dev loop with stldev

Pair fluent-sdfx with stldev for a watch-rebuild-preview cycle that auto-reloads STLs in f3d.

Designing a 3D part is a tweak-and-look loop. You change a number, save, and want to see the new geometry in three seconds — not run a build, find the file, drag it into a viewer.

[`stldev`](https://github.com/snowbldr/stldev) is a small CLI that closes that loop. It watches your Go source, runs your build command on save, and reloads the resulting STLs in tiled f3d windows.

## Install

```bash
go install github.com/snowbldr/stldev@latest
```

Requires `f3d` on your `PATH` ([install](/install/) covers it).

Verify:

```bash
stldev --help
```

## Basic usage

From the directory of a fluent-sdfx project that has a `main.go` producing an STL:

```bash
stldev -cmd "go run ." part.stl
```

This:

1. Watches the current directory for `.go` file changes.
2. Runs `go run .` on each save.
3. Opens an f3d window for `part.stl` and reloads it whenever the file changes on disk.
4. Shows a 3D-grid placeholder during compilation, so you always know where the part used to be.

## Multiple parts

For an assembly with several outputs, name each STL on the command line. stldev opens one f3d window per file and tiles them:

```bash
stldev -cmd "go run ." body.stl lid.stl mount.stl
```

A common pattern: a `main.go` that generates several STLs side-by-side so you can iterate on the whole assembly at once.

## Makefile-driven workflow

Once you have more than two parts, the build command grows. Put it in a Makefile:

```make
.PHONY: dev build

dev:
	stldev -cmd "go run ./parts/all" body.stl lid.stl mount.stl

build:
	go run ./parts/all
```

Then `make dev` is your iteration loop and `make build` produces the final STLs.

> [!TIP]
> Keep the iteration `cellsPerMM` low — `1.0` to `2.0` — for fast rebuilds. Bump it back up to `5.0`+ for the final export. Render time scales roughly with cells³, so halving the density is an 8× speedup.

## What stldev doesn't do

- It doesn't slice for printing — that's your slicer's job.
- It doesn't transform the STLs — output what you want, exactly how you want it.
- It doesn't replace `go test` or other CI — it's purely a viewer + watcher.

For all the flags, the `--` passthrough to f3d (e.g. for camera presets), and notes on multi-monitor layouts, see the [stldev README](https://github.com/snowbldr/stldev).
