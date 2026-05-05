# VSCode snippets for fluent-sdfx

Drop-in snippets for the most common fluent-sdfx idioms — primitives, anchor placement, layout helpers, the recipe pattern, and validation tests.

## Install

**Project-local** (recommended — share with collaborators via git):

```bash
mkdir -p .vscode
cp editors/vscode/fluent-sdfx.code-snippets .vscode/
```

**User-global** (available across all projects):

VSCode → `Cmd/Ctrl-Shift-P` → "Snippets: Configure Snippets" → "go.json" → paste the contents of `fluent-sdfx.code-snippets` inside the outer braces.

## Snippet index

| Prefix | What it expands to |
|---|---|
| `fsdfxmain` | full `main.go` scaffold with imports + STL output |
| `fsdfxcyl` | `solid.Cylinder(h, r, chamfer)` |
| `fsdfxbox` | `solid.Box(v3.XYZ(...), chamfer)` |
| `fsdfxsphere` | `solid.Sphere(r)` |
| `fsdfxdrill` | `body.Cut(tool.Top().On(body.Top()).Solid())` |
| `fsdfxstack` | `part.OnTopOf(target.Top()).Union()` |
| `fsdfxinside` | `inner.Inside(outer).Union()` |
| `fsdfxbottom` | `part.BottomAt(0)` |
| `fsdfxpolar` | `part.Multi(layout.Polar(r, n)...)` |
| `fsdfxgrid` | `part.Multi(layout.Grid(stepX, stepY, nx, ny)...)` |
| `fsdfxcorners` | `part.Multi(layout.RectCorners(w, d)...)` |
| `fsdfxrecipe` | recipe pattern scaffold (ingredients + method) |
| `fsdfxtest` | `_test.go` scaffold with `validate.Require*` calls |

All snippets restrict themselves to `.go` files. Tab cycles through editable placeholders; `Esc` finishes.
