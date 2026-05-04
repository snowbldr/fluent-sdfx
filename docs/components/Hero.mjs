import { div, h1, p, span, a, section, code, pre } from '@srfnstack/fntags'
import { fnlink } from '@srfnstack/fntags/src/fnroute.mjs'

const HERO_CODE = `solid.Cylinder(20, 10, 1).
    Cut(
        solid.Cylinder(25, 2, 0).TranslateX(5),
        solid.Cylinder(25, 2, 0).TranslateX(-5),
        solid.Cylinder(25, 2, 0).TranslateY(5),
        solid.Cylinder(25, 2, 0).TranslateY(-5),
    ).
    STL("part.stl", 3.0)`

const FEATURES = [
  {
    title: 'Chainable',
    body: 'Every transform and boolean returns a new value. The whole part is one expression.'
  },
  {
    title: 'Degrees by default',
    body: 'No more math.Pi/180. fluent-sdfx converts internally.'
  },
  {
    title: 'No error returns',
    body: 'CAD-geometry errors are programming bugs. Constructors panic; you chain.'
  },
  {
    title: 'Real renderer',
    body: 'STL, 3MF, DXF, SVG, PNG. Parallel marching cubes, optional decimation.'
  }
]

const Cta = (slug, text, kind = 'primary') =>
  fnlink({ to: `/${slug}/`, class: `cta cta-${kind}` }, text)

const FeatureCard = (f) =>
  div({ class: 'feature' },
    div({ class: 'feature-title' }, f.title),
    p({ class: 'feature-body' }, f.body)
  )

export const Hero = () =>
  section({ class: 'hero' },
    div({ class: 'hero-grid' },
      div({ class: 'hero-text' },
        div({ class: 'hero-eyebrow' }, 'Go · CAD · SDF'),
        h1({ class: 'hero-title' },
          'Build 3D parts the way you ',
          span({ class: 'hero-accent' }, 'describe them'),
          '.'
        ),
        p({ class: 'hero-lede' },
          'fluent-sdfx wraps the sdfx CAD kernel with two chainable types — ',
          code({ class: 'hero-inline' }, 'shape.Shape'),
          ' and ',
          code({ class: 'hero-inline' }, 'solid.Solid'),
          ' — and a fluent API that reads like a description of the part.'
        ),
        div({ class: 'hero-actions' },
          Cta('quickstart', 'Quickstart →', 'primary'),
          Cta('install', 'Install', 'ghost')
        )
      ),
      div({ class: 'hero-code' },
        div({ class: 'hero-code-header' },
          span({ class: 'hero-code-dot' }),
          span({ class: 'hero-code-dot' }),
          span({ class: 'hero-code-dot' }),
          span({ class: 'hero-code-name' }, 'part.go')
        ),
        pre({ class: 'hero-code-body' },
          code({}, HERO_CODE)
        ),
        div({ class: 'hero-code-result' },
          span({ class: 'hero-code-prompt' }, '$'),
          ' go run .',
          span({ class: 'hero-code-out' }, ' → part.stl')
        )
      )
    ),
    div({ class: 'features' },
      ...FEATURES.map(FeatureCard)
    )
  )
