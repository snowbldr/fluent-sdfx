# Documentation site

Static site for fluent-sdfx, built with [fntags](https://github.com/srfnstack/fntags). **No build step:** the `docs/` directory is the deployable site. Open `index.html` in any modern browser via a static server and it works.

## Local dev

```bash
make serve
# or, from this directory:
../tools/serve.py
```

Then open `http://localhost:7174/`. The dev server (a small Python script) mimics GitHub Pages 404 behavior so deep-link reloads (`/quickstart/`) trigger the same SPA-recovery path used in production.

## Deploy

GitHub Pages: configure the repo to publish from the `docs/` directory on the main branch. Nothing to build, nothing to ship — the directory *is* the site.

If you fork to a different repo name, edit `KNOWN_PREFIXES` in `404.html` and `constants.mjs` to match your project URL prefix.

## Structure

```
docs/
├── index.html             entry HTML — importmap + Prism CDN scripts + bootstrap module
├── 404.html               GitHub-Pages SPA deep-link recovery
├── llms.txt               agent-oriented project reference
├── main.mjs               app boot (mounts App, handles 404 redirect)
├── App.mjs                routeSwitch over allPages
├── constants.mjs          BASE detection (root vs /fluent-sdfx/)
├── nav.mjs                navigation config (single source of truth)
├── lib/                   vendored fntags (fntags.mjs, fnroute.mjs, fnelements.mjs, svgelements.mjs)
├── components/
│   ├── Sidebar.mjs
│   ├── ThemeToggle.mjs
│   ├── Hero.mjs           landing-page hero
│   ├── Markdown.mjs       fetch + parse + Prism highlight
│   └── PageView.mjs
├── content/*.md           one markdown file per page
├── images/                f3d-rendered hero PNGs
├── stl/                   viewer-grade STLs + manifest.json
└── styles/global.css      design tokens, layout, markdown styling
```

## Authoring conventions

Each markdown file in `content/` is a single page rendered as a route. The first `# H1` is the page title; the paragraph after it is the lede. Section name (e.g. "Getting started") is rendered as an eyebrow above the title, sourced from `nav.mjs`.

### Code blocks tied to tutorial files

```` markdown
<!-- src: tutorial/04-quickstart/01-cylinder/main.go -->
```go
// inlined source from the file above
package main
…
```
````

Run `make build` (or the pre-commit hook) to re-inline the code from the canonical Go file. The HTML comment is the source-of-truth pointer; the inlined fence is what GitHub renders directly when viewing the .md.

### Figures

```` markdown
<figure>
  <img src="../images/04-quickstart-01-cylinder.png" alt="caption">
  <figcaption>caption</figcaption>
</figure>
````

Image paths are relative to the .md file's location so GitHub renders them correctly. The Markdown component rewrites them to base-absolute paths at runtime.

### Callouts

```` markdown
> [!TIP]
> Make tools longer than the body so cuts go all the way through.
````

Standard [GitHub alert syntax](https://github.com/orgs/community/discussions/16925). Variants: `[!TIP]`, `[!WARNING]`, `[!NOTE]`, `[!IMPORTANT]`, `[!CAUTION]`. GitHub renders them as styled blockquotes natively; our Markdown component upgrades them to `<aside class="callout callout-tip">…</aside>` for matching docs styling.

## Regenerating screenshots & viewer STLs

```bash
make screenshots
# or directly:
../tutorial/render-screenshots.sh
```

For each tutorial step under `../tutorial/` the pipeline:

1. Runs `go run` to produce a full-resolution `out.stl`.
2. Hands the STL to **f3d** for an isometric hero PNG → `images/<page>-<step>.png`.
3. Decimates the STL to ~1500 triangles via `tools/decimate-stl` (uses fluent-sdfx's bundled meshoptimizer) → `stl/<page>-<step>.stl` (~75 KB each, ~6 MB for all 80 steps).
4. Writes `stl/manifest.json` listing every available STL. The Markdown component reads it on first render and attaches a "View in 3D" toggle to any figure whose stem matches.

The toggle button replaces the static `<img>` with a Three.js + OrbitControls canvas the reader can spin. Three.js itself loads lazily — only on the first click, anywhere on the site — so the docs stay light by default.

Slow (~80 `go run`s × `f3d` × decimate); not part of `make build`.
