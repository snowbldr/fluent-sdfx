# Install

Install fluent-sdfx, the f3d viewer, and the optional stldev dev-loop tool.

fluent-sdfx is a Go library. You'll want three things on your machine before you build anything:

1. **Go 1.21 or newer** — to compile the library and your designs.
2. **fluent-sdfx itself** — added with `go get`.
3. **f3d** *(optional but recommended)* — to view STL files. Free, fast, headless-capable.
4. **stldev** *(optional)* — a watch-rebuild-preview loop that pairs nicely with fluent-sdfx; covered on the [dev-loop page](/dev-loop/).

## 1. Install Go

If you don't already have it, grab Go from [go.dev/dl](https://go.dev/dl/) or your package manager of choice.

```bash
# macOS
brew install go

# Debian / Ubuntu
sudo apt install golang-go

# Fedora
sudo dnf install golang
```

Verify:

```bash
go version
# go version go1.25.0 darwin/arm64
```

## 2. Install fluent-sdfx

Inside any Go module:

```bash
go get github.com/snowbldr/fluent-sdfx
```

The library is pure Go and has no system dependencies (one transitively-pulled package, [meshoptimizer](https://github.com/zeux/meshoptimizer), uses CGo for STL decimation — only relevant if you pass a decimation factor to `STL`).

## 3. Install f3d

f3d is a fast, modern STL/OBJ/3MF viewer. We use it both interactively (open a part to look at it) and headlessly (the docs site renders every figure on every page through f3d).

```bash
# macOS
brew install f3d

# Debian / Ubuntu
sudo apt install f3d

# Other platforms
# https://f3d.app — releases for Windows, Linux, macOS
```

Verify:

```bash
f3d --version
# F3D 3.4.1
```

> [!TIP]
> You can re-render every figure on this site locally with `./tutorial/render-screenshots.sh` — it walks every `tutorial/<page>/<step>/main.go`, builds the part, and runs f3d to produce a PNG. The same pipeline keeps the docs in sync with the code.

## 4. Verify your install

The fastest way to confirm everything works end-to-end is to scaffold a tiny module and run a one-liner. The next page, [Project setup](/project-setup/), walks through that in detail.
