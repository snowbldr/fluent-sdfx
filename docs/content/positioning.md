# Positioning

Place parts by name, not by hand-rolled bbox math. Every solid has 27 named anchors on its bounding box; you align them to other anchors with verbs like `On`, `Above`, `RightOf`, or to absolute coordinates with `BottomAt`, `CenterAt`, etc. Pair with the `layout` package and a single `.Multi(...)` call for grids and rings.

This is the canonical way to position things in fluent-sdfx. `Translate` still works — but reach for an anchor first; the code is shorter, the intent is louder, and the bug surface (off-by-one bbox math) shrinks to zero.

## The pattern in 30 seconds

```go
// Stack a sphere on top of a box.
result := sphere.OnTopOf(body.Top()).Union()

// Drill a hole down through the body's top, 2mm in from the right edge.
hole := solid.Cylinder(20, 2, 0).
    Bottom().Above(body.Top()).Solid().
    TranslateX(body.Bounds().Anchor(1, 0, 1).X - 2)
result = body.Cut(hole)

// Sit something flat on the build plate.
part := part.BottomAt(0)

// Stamp 6 pegs on a 20 mm circle.
pegs := peg.Multi(layout.Polar(20, 6)...)
```

Three concepts:

1. **Anchors** — `s.Top()`, `s.BottomFrontRight()`, `s.AnchorAt(x,y,z)` — return an `AnchoredSolid` (the solid + the anchor point in world space).
2. **Placement verbs** on `AnchoredSolid` — `On`, `Above`, `Below`, `RightOf`, `LeftOf`, `Behind`, `InFrontOf`, `At`, `AtX`, `AtY`, `AtZ`, `ShiftX/Y/Z` — move the solid so its anchor lands on a target. Relative verbs return a `Placement`; absolute verbs return a `*Solid`.
3. **`Placement` finalizers** — `.Union()`, `.Cut()`, `.Intersect()`, `.Solid()` — finish the chain. `.Solid()` is the escape hatch when you want the moved part on its own.

## Anchors

A solid's bounding box has 27 named anchors: 6 face centers, 12 edge midpoints, 8 corners, plus the box center. Each tab below shows where the anchor lives on the box (the indicator pin sticks outward).

### Faces

<div class="tab-gallery">
  <div class="tab-gallery-tabs">
    <button data-tab="top">Top</button>
    <button data-tab="bottom">Bottom</button>
    <button data-tab="right">Right</button>
    <button data-tab="left">Left</button>
    <button data-tab="back">Back</button>
    <button data-tab="front">Front</button>
  </div>
  <div data-pane="top">
    <figure>
      <img src="../images/positioning-anchor-top.png" alt="Top face anchor">
    </figure>
    <pre><code class="language-go">box.Top()    // (0, 0, +Z)</code></pre>
  </div>
  <div data-pane="bottom">
    <figure>
      <img src="../images/positioning-anchor-bottom.png" alt="Bottom face anchor">
    </figure>
    <pre><code class="language-go">box.Bottom() // (0, 0, -Z)</code></pre>
  </div>
  <div data-pane="right">
    <figure>
      <img src="../images/positioning-anchor-right.png" alt="Right face anchor">
    </figure>
    <pre><code class="language-go">box.Right()  // (+X, 0, 0)</code></pre>
  </div>
  <div data-pane="left">
    <figure>
      <img src="../images/positioning-anchor-left.png" alt="Left face anchor">
    </figure>
    <pre><code class="language-go">box.Left()   // (-X, 0, 0)</code></pre>
  </div>
  <div data-pane="back">
    <figure>
      <img src="../images/positioning-anchor-back.png" alt="Back face anchor">
    </figure>
    <pre><code class="language-go">box.Back()   // (0, +Y, 0)</code></pre>
  </div>
  <div data-pane="front">
    <figure>
      <img src="../images/positioning-anchor-front.png" alt="Front face anchor">
    </figure>
    <pre><code class="language-go">box.Front()  // (0, -Y, 0)</code></pre>
  </div>
</div>

### Edges

<div class="tab-gallery">
  <div class="tab-gallery-tabs">
    <button data-tab="top-right">TopRight</button>
    <button data-tab="top-left">TopLeft</button>
    <button data-tab="top-front">TopFront</button>
    <button data-tab="top-back">TopBack</button>
    <button data-tab="bottom-right">BottomRight</button>
    <button data-tab="bottom-left">BottomLeft</button>
    <button data-tab="bottom-front">BottomFront</button>
    <button data-tab="bottom-back">BottomBack</button>
    <button data-tab="front-right">FrontRight</button>
    <button data-tab="front-left">FrontLeft</button>
    <button data-tab="back-right">BackRight</button>
    <button data-tab="back-left">BackLeft</button>
  </div>
  <div data-pane="top-right">
    <figure>
      <img src="../images/positioning-anchor-top-right.png" alt="Top-right edge anchor">
    </figure>
    <pre><code class="language-go">box.TopRight()</code></pre>
  </div>
  <div data-pane="top-left">
    <figure>
      <img src="../images/positioning-anchor-top-left.png" alt="Top-left edge anchor">
    </figure>
    <pre><code class="language-go">box.TopLeft()</code></pre>
  </div>
  <div data-pane="top-front">
    <figure>
      <img src="../images/positioning-anchor-top-front.png" alt="Top-front edge anchor">
    </figure>
    <pre><code class="language-go">box.TopFront()</code></pre>
  </div>
  <div data-pane="top-back">
    <figure>
      <img src="../images/positioning-anchor-top-back.png" alt="Top-back edge anchor">
    </figure>
    <pre><code class="language-go">box.TopBack()</code></pre>
  </div>
  <div data-pane="bottom-right">
    <figure>
      <img src="../images/positioning-anchor-bottom-right.png" alt="Bottom-right edge anchor">
    </figure>
    <pre><code class="language-go">box.BottomRight()</code></pre>
  </div>
  <div data-pane="bottom-left">
    <figure>
      <img src="../images/positioning-anchor-bottom-left.png" alt="Bottom-left edge anchor">
    </figure>
    <pre><code class="language-go">box.BottomLeft()</code></pre>
  </div>
  <div data-pane="bottom-front">
    <figure>
      <img src="../images/positioning-anchor-bottom-front.png" alt="Bottom-front edge anchor">
    </figure>
    <pre><code class="language-go">box.BottomFront()</code></pre>
  </div>
  <div data-pane="bottom-back">
    <figure>
      <img src="../images/positioning-anchor-bottom-back.png" alt="Bottom-back edge anchor">
    </figure>
    <pre><code class="language-go">box.BottomBack()</code></pre>
  </div>
  <div data-pane="front-right">
    <figure>
      <img src="../images/positioning-anchor-front-right.png" alt="Front-right edge anchor">
    </figure>
    <pre><code class="language-go">box.FrontRight()</code></pre>
  </div>
  <div data-pane="front-left">
    <figure>
      <img src="../images/positioning-anchor-front-left.png" alt="Front-left edge anchor">
    </figure>
    <pre><code class="language-go">box.FrontLeft()</code></pre>
  </div>
  <div data-pane="back-right">
    <figure>
      <img src="../images/positioning-anchor-back-right.png" alt="Back-right edge anchor">
    </figure>
    <pre><code class="language-go">box.BackRight()</code></pre>
  </div>
  <div data-pane="back-left">
    <figure>
      <img src="../images/positioning-anchor-back-left.png" alt="Back-left edge anchor">
    </figure>
    <pre><code class="language-go">box.BackLeft()</code></pre>
  </div>
</div>

### Corners

<div class="tab-gallery">
  <div class="tab-gallery-tabs">
    <button data-tab="tfr">TopFrontRight</button>
    <button data-tab="tfl">TopFrontLeft</button>
    <button data-tab="tbr">TopBackRight</button>
    <button data-tab="tbl">TopBackLeft</button>
    <button data-tab="bfr">BottomFrontRight</button>
    <button data-tab="bfl">BottomFrontLeft</button>
    <button data-tab="bbr">BottomBackRight</button>
    <button data-tab="bbl">BottomBackLeft</button>
  </div>
  <div data-pane="tfr">
    <figure>
      <img src="../images/positioning-anchor-top-front-right.png" alt="Top-front-right corner">
    </figure>
    <pre><code class="language-go">box.TopFrontRight()</code></pre>
  </div>
  <div data-pane="tfl">
    <figure>
      <img src="../images/positioning-anchor-top-front-left.png" alt="Top-front-left corner">
    </figure>
    <pre><code class="language-go">box.TopFrontLeft()</code></pre>
  </div>
  <div data-pane="tbr">
    <figure>
      <img src="../images/positioning-anchor-top-back-right.png" alt="Top-back-right corner">
    </figure>
    <pre><code class="language-go">box.TopBackRight()</code></pre>
  </div>
  <div data-pane="tbl">
    <figure>
      <img src="../images/positioning-anchor-top-back-left.png" alt="Top-back-left corner">
    </figure>
    <pre><code class="language-go">box.TopBackLeft()</code></pre>
  </div>
  <div data-pane="bfr">
    <figure>
      <img src="../images/positioning-anchor-bottom-front-right.png" alt="Bottom-front-right corner">
    </figure>
    <pre><code class="language-go">box.BottomFrontRight()</code></pre>
  </div>
  <div data-pane="bfl">
    <figure>
      <img src="../images/positioning-anchor-bottom-front-left.png" alt="Bottom-front-left corner">
    </figure>
    <pre><code class="language-go">box.BottomFrontLeft()</code></pre>
  </div>
  <div data-pane="bbr">
    <figure>
      <img src="../images/positioning-anchor-bottom-back-right.png" alt="Bottom-back-right corner">
    </figure>
    <pre><code class="language-go">box.BottomBackRight()</code></pre>
  </div>
  <div data-pane="bbl">
    <figure>
      <img src="../images/positioning-anchor-bottom-back-left.png" alt="Bottom-back-left corner">
    </figure>
    <pre><code class="language-go">box.BottomBackLeft()</code></pre>
  </div>
</div>

### Center & arbitrary

`AnchorAt(x, y, z int)` is the escape hatch — each component is min at `-1`, center at `0`, max at `+1`. The center is `AnchorAt(0, 0, 0)`. There's no `Center()` selector because `Center()` is already a transform that recentres a solid on the origin.

```go
box.AnchorAt(0, 0, 0)   // bbox center (no protrusion to draw)
box.AnchorAt(1, 1, 0)   // back-right edge midpoint, same as box.BackRight()
```

## Placement verbs

Once you have an anchor, you place it. **Relative** verbs take another anchor and return a `Placement` so you can finish with `.Union()`, `.Cut()`, `.Intersect()`, or `.Solid()`. **Absolute** verbs take a literal coordinate and return a `*Solid` directly.

<div class="tab-gallery">
  <div class="tab-gallery-tabs">
    <button data-tab="on">On</button>
    <button data-tab="above">Above</button>
    <button data-tab="below">Below</button>
    <button data-tab="right-of">RightOf</button>
    <button data-tab="left-of">LeftOf</button>
    <button data-tab="behind">Behind</button>
    <button data-tab="in-front-of">InFrontOf</button>
    <button data-tab="inside">Inside</button>
    <button data-tab="bottom-at">BottomAt</button>
  </div>
  <div data-pane="on">
    <figure>
      <img src="../images/positioning-verb-on.png" alt="cap.Bottom().On(host.Top())">
    </figure>
    <pre><code class="language-go">cap.Bottom().On(host.Top()).Union()</code></pre>
    <div class="tab-caption">Cap's bottom anchor lands flush on the host's top anchor.</div>
  </div>
  <div data-pane="above">
    <figure>
      <img src="../images/positioning-verb-above.png" alt="m.Bottom().Above(host.Top(), 2)">
    </figure>
    <pre><code class="language-go">m.Bottom().Above(host.Top(), 2).Union()</code></pre>
    <div class="tab-caption">Same as <code>On</code>, but with a 2mm gap along +Z.</div>
  </div>
  <div data-pane="below">
    <figure>
      <img src="../images/positioning-verb-below.png" alt="m.Top().Below(host.Bottom(), 2)">
    </figure>
    <pre><code class="language-go">m.Top().Below(host.Bottom(), 2).Union()</code></pre>
  </div>
  <div data-pane="right-of">
    <figure>
      <img src="../images/positioning-verb-right-of.png" alt="m.Left().RightOf(host.Right(), 2)">
    </figure>
    <pre><code class="language-go">m.Left().RightOf(host.Right(), 2).Union()</code></pre>
  </div>
  <div data-pane="left-of">
    <figure>
      <img src="../images/positioning-verb-left-of.png" alt="m.Right().LeftOf(host.Left(), 2)">
    </figure>
    <pre><code class="language-go">m.Right().LeftOf(host.Left(), 2).Union()</code></pre>
  </div>
  <div data-pane="behind">
    <figure>
      <img src="../images/positioning-verb-behind.png" alt="m.Front().Behind(host.Back(), 2)">
    </figure>
    <pre><code class="language-go">m.Front().Behind(host.Back(), 2).Union()</code></pre>
  </div>
  <div data-pane="in-front-of">
    <figure>
      <img src="../images/positioning-verb-in-front-of.png" alt="m.Back().InFrontOf(host.Front(), 2)">
    </figure>
    <pre><code class="language-go">m.Back().InFrontOf(host.Front(), 2).Union()</code></pre>
  </div>
  <div data-pane="inside">
    <figure>
      <img src="../images/positioning-verb-inside.png" alt="m.Inside(host)">
    </figure>
    <pre><code class="language-go">m.Inside(host).Union()</code></pre>
    <div class="tab-caption">Centers the mover's bbox on the host's bbox center.</div>
  </div>
  <div data-pane="bottom-at">
    <figure>
      <img src="../images/positioning-verb-bottom-at.png" alt="part.BottomAt(0)">
    </figure>
    <pre><code class="language-go">part.BottomAt(0)</code></pre>
    <div class="tab-caption">Absolute setter — drops the part so its bottom face sits on z=0. Returns <code>*Solid</code> directly.</div>
  </div>
</div>

### Sugar on `*Solid`

For the most common case — placing a solid relative to another solid's matching face — there's a shorter form on the moved solid:

| Sugar | Equivalent |
|---|---|
| `s.OnTopOf(t.Top())` | `s.Bottom().Above(t.Top())` |
| `s.UnderneathOf(t.Bottom())` | `s.Top().Below(t.Bottom())` |
| `s.RightOf(t.Right())` | `s.Left().RightOf(t.Right())` |
| `s.LeftOf(t.Left())` | `s.Right().LeftOf(t.Left())` |
| `s.InFrontOf(t.Front())` | `s.Back().InFrontOf(t.Front())` |
| `s.BehindOf(t.Back())` | `s.Front().Behind(t.Back())` |
| `s.Inside(t)` | `s.AnchorAt(0,0,0).On(t.AnchorAt(0,0,0))` |

All return a `Placement`. The absolute setters — `BottomAt`, `TopAt`, `LeftAt`, `RightAt`, `FrontAt`, `BackAt`, `CenterAt` — each leave the other axes alone and return `*Solid`.

> [!TIP]
> The sugar verbs collapse the two-solid case into one chain, but the underlying anchor form is more flexible — use it when you need cross-axis pairings like "this part's right face on that part's top face" (`s.Right().On(t.Top())`).

## `ShiftX/Y/Z` — anchor tweaks

Sometimes the target you want isn't an anchor exactly — it's an anchor *plus* a small offset. `ShiftX/Y/Z` move the anchor point without moving the solid yet, so you can chain it before the placement verb:

```go
// boss centered on plate, but shifted 4mm forward of plate's top center
boss.OnTopOf(plate.Top().ShiftY(-4)).Union()
```

## Layouts

The `layout` package returns position arrays designed to flow straight into the variadic `.Multi(...)`:

```go
import "github.com/snowbldr/fluent-sdfx/layout"

peg.Multi(layout.Polar(20, 6)...)            // 6 pegs on a r=20 circle
hole.Multi(layout.Grid(10, 10, 4, 4)...)     // 16 holes on a 10mm grid
standoff.Multi(layout.RectCorners(76, 46)...) // 4 standoffs at panel corners
post.Multi(layout.PolarArc(15, 5, 0, 90)...)  // 5 posts spanning a quarter arc
```

Available functions:

| 3D | 2D | Use |
|---|---|---|
| `layout.Polar(r, n)` | `layout.Polar2(r, n)` | n positions evenly spaced on a circle |
| `layout.PolarArc(r, n, startDeg, sweepDeg)` | — | arc instead of full circle |
| `layout.Grid(stepX, stepY, nx, ny)` | `layout.Grid2(...)` | XY grid centered on origin |
| `layout.Line(p0, p1, n)` | — | n equally spaced from p0 to p1 |
| `layout.RectCorners(w, d)` | — | the 4 XY corners of a rectangle |
| `layout.BoxCorners(size)` | — | the 8 corners of a box |

## 2D shapes

Everything above mirrors to `*Shape` with the 2D anchor set: `Top`, `Bottom`, `Left`, `Right`, the 4 corners, and `AnchorAt(x, y int)`. 2D follows screen convention — `+Y` is up.

```go
label.OnTopOf(panel.Top()).Union()
hole.Inside(plate).Cut()
```

## When to reach for what

- **`.Translate(v)`** — when the move is a literal vector unrelated to any other geometry.
- **`anchor.At(p)`** — when the move IS to a literal coordinate, but you want to land a specific edge / corner there.
- **`anchor.On(target)` / sugar verbs** — anything relative to another part.
- **`BottomAt(z)` / `LeftAt(x)` etc.** — flush a specific face against an axis plane.
- **`layout.X` + `.Multi`** — repeated copies in a regular pattern.
