import { aside, nav, div, h2, ul, li, img } from '@srfnstack/fntags'
import { fnlink, pathState } from '@srfnstack/fntags/src/fnroute.mjs'
import { sections, pageBySlug } from '../nav.mjs'

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
    )
  )
