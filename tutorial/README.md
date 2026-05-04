# Tutorial code

Runnable Go programs paired with the [docs site](../docs). Each step is its own `package main` so you can run any one independently.

## Layout

```
tutorial/
├── 04-quickstart/
│   ├── 01-cylinder/main.go
│   ├── 02-with-one-hole/main.go
│   └── ...
├── 05-vectors-types/...
└── render-screenshots.sh
```

## Run a single step

From the repo root:

```bash
go run ./tutorial/04-quickstart/03-with-four-holes
```

Each step writes its output to `out.stl` (3D) or `out.png` (2D) inside its own directory. Both are gitignored.

## Regenerate all docs screenshots

Requires [`f3d`](https://f3d.app) on your `PATH`.

```bash
./tutorial/render-screenshots.sh
```

The script walks every step, runs it, then renders a normalised hero PNG into `docs/public/images/<page>-<step>.png`.

A step can override the default f3d camera by dropping an `f3d.flags` file alongside its `main.go` — the file's contents replace the default flag string.
