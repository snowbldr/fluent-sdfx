import { aside, nav, div, h2, ul, li, img, a, span, svg, path } from '@srfnstack/fntags'
import { fnlink, pathState } from '@srfnstack/fntags/src/fnroute.mjs'
import { sections, pageBySlug } from '../nav.mjs'

const GITHUB_URL = 'https://github.com/snowbldr/fluent-sdfx'

// Inline GitHub octocat mark — keeps us off external icon hosts.
const Octocat = () =>
  svg({ viewBox: '0 0 24 24', 'aria-hidden': 'true', class: 'octocat' },
    path({
      fill: 'currentColor',
      d: 'M12 .5C5.65.5.5 5.65.5 12c0 5.08 3.29 9.39 7.86 10.91.58.11.79-.25.79-.55 0-.27-.01-.99-.02-1.94-3.2.69-3.87-1.54-3.87-1.54-.52-1.33-1.28-1.69-1.28-1.69-1.05-.72.08-.71.08-.71 1.16.08 1.77 1.19 1.77 1.19 1.03 1.76 2.69 1.25 3.35.96.1-.74.4-1.25.73-1.54-2.55-.29-5.24-1.28-5.24-5.69 0-1.26.45-2.29 1.18-3.1-.12-.29-.51-1.46.11-3.05 0 0 .97-.31 3.18 1.18a11.1 11.1 0 015.79 0c2.21-1.49 3.18-1.18 3.18-1.18.62 1.59.23 2.76.11 3.05.74.81 1.18 1.84 1.18 3.1 0 4.42-2.69 5.4-5.25 5.68.41.36.78 1.06.78 2.13 0 1.54-.01 2.78-.01 3.16 0 .31.21.67.8.55C20.21 21.39 23.5 17.08 23.5 12 23.5 5.65 18.35.5 12 .5z'
    })
  )

const Footer = () =>
  div({ class: 'sidebar-footer' },
    a({
      href: GITHUB_URL,
      target: '_blank',
      rel: 'noopener noreferrer',
      class: 'github-link',
      'aria-label': 'View fluent-sdfx on GitHub'
    },
      Octocat(),
      span({ class: 'github-link-text' }, 'GitHub')
    )
  )

const Brand = () =>
  div({ class: 'sidebar-brand' },
    fnlink({ to: '/', class: 'brand-link' },
      img({ class: 'brand-mark', src: './public/logo/logo-128.webp', alt: '', width: '72', height: '72' }),
      div({ class: 'brand-text' },
        div({ class: 'brand-title' }, 'fluent-sdfx'),
        div({ class: 'brand-tag' }, 'fluent SDF CAD for Go')
      )
    )
  )

// pathState's value is { currentPath, rootPath, context }. We compare
// currentPath against the link target. Trailing slashes are normalised
// because the browser's pathname strips them but our nav paths don't.
const normalisePath = p => {
  if (!p) return '/'
  return p.length > 1 && p.endsWith('/') ? p.slice(0, -1) : p
}

const NavLink = (page) =>
  li({ class: 'nav-item' },
    fnlink({
      to: page.path,
      class: pathState.bindAttr(() =>
        normalisePath(pathState().currentPath) === normalisePath(page.path)
          ? 'nav-link active'
          : 'nav-link'
      )
    }, page.title)
  )

const Section = (s) => {
  // section.pages reference slugs from nav.mjs; resolve to the path-aware
  // pages from allPages (via pageBySlug).
  const resolved = s.pages.map(p => pageBySlug[p.slug])
  return div({ class: 'nav-section' },
    h2({ class: 'nav-section-title' }, s.title),
    ul({ class: 'nav-list' }, ...resolved.map(NavLink))
  )
}

export const Sidebar = () =>
  aside({ class: 'sidebar' },
    Brand(),
    nav({ class: 'nav', 'aria-label': 'Documentation' },
      ...sections.map(Section)
    ),
    Footer()
  )
