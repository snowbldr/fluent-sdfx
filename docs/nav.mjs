// Single source of truth for navigation. Each page slug also names its
// content file: src/content/<slug>.md.
//
// `path` is the URL the page is served at. It defaults to `/<slug>/`,
// except the index page which is served at `/`.

export const sections = [
  {
    title: 'Getting started',
    pages: [
      { slug: 'index', title: 'Introduction', path: '/' },
      { slug: 'install', title: 'Install' },
      { slug: 'project-setup', title: 'Project setup' },
      { slug: 'dev-loop', title: 'Dev loop with stldev' },
      { slug: 'quickstart', title: 'Quickstart' }
    ]
  },
  {
    title: 'Foundations',
    pages: [
      { slug: 'vectors-types', title: 'Vectors & types' },
      { slug: 'shapes-2d', title: '2D shapes' },
      { slug: 'solids-3d', title: '3D solids' },
      { slug: 'booleans', title: 'Booleans' },
      { slug: 'transforms', title: 'Transforms' },
      { slug: 'positioning', title: 'Positioning' }
    ]
  },
  {
    title: 'Operations',
    pages: [
      { slug: '2d-to-3d', title: '2D → 3D' },
      { slug: 'smooth-blends', title: 'Smooth blends' },
      { slug: 'modifiers', title: 'Modifiers' },
      { slug: 'patterns', title: 'Patterns' },
      { slug: 'cross-sections', title: 'Cross-sections' },
      { slug: 'text-2d-output', title: 'Text & 2D output' },
      { slug: 'obj-overview', title: 'Parametric helpers' },
      { slug: 'output-resolution', title: 'Output & resolution' }
    ]
  },
  {
    title: 'Cookbook',
    pages: [
      { slug: 'cookbook-bolt', title: 'Bolt assembly' },
      { slug: 'cookbook-enclosure', title: 'Enclosure' },
      { slug: 'cookbook-gear', title: 'Gear' },
      { slug: 'cookbook-lantern', title: 'Lantern' }
    ]
  },
  {
    title: 'Reference',
    pages: [{ slug: 'api-reference', title: 'API reference' }]
  }
]

export const allPages = sections.flatMap(s =>
  s.pages.map(p => ({ ...p, section: s.title, path: p.path ?? `/${p.slug}/` }))
)

export const pageBySlug = Object.fromEntries(allPages.map(p => [p.slug, p]))
export const pageByPath = Object.fromEntries(allPages.map(p => [p.path, p]))
