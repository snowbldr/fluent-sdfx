import { div, article, p, footer } from '@srfnstack/fntags'
import { fnlink } from '@srfnstack/fntags/src/fnroute.mjs'
import { allPages, pageBySlug } from '../nav.mjs'
import { Markdown } from './Markdown.mjs'
import { Hero } from './Hero.mjs'

function neighbors(slug) {
  const i = allPages.findIndex(p => p.slug === slug)
  if (i < 0) return { prev: null, next: null }
  return {
    prev: i > 0 ? allPages[i - 1] : null,
    next: i + 1 < allPages.length ? allPages[i + 1] : null
  }
}

const PageFooter = (slug) => {
  const { prev, next } = neighbors(slug)
  return footer({ class: 'page-footer' },
    prev
      ? fnlink({ to: prev.path, class: 'page-nav prev' },
          div({ class: 'page-nav-label' }, '← Previous'),
          div({ class: 'page-nav-title' }, prev.title)
        )
      : div({ class: 'page-nav-spacer' }),
    next
      ? fnlink({ to: next.path, class: 'page-nav next' },
          div({ class: 'page-nav-label' }, 'Next →'),
          div({ class: 'page-nav-title' }, next.title)
        )
      : div({ class: 'page-nav-spacer' })
  )
}

export const PageView = (slug) => {
  const page = pageBySlug[slug]
  if (!page) {
    return article({ class: 'page page-missing' },
      p({}, `Page “${slug}” not found.`)
    )
  }

  // Index gets a hero in place of the markdown's H1 + lede.
  if (slug === 'index') {
    return article({ class: 'page page-index' },
      Hero(),
      Markdown(slug),
      PageFooter(slug)
    )
  }

  return article({ class: 'page' },
    div({ class: 'page-eyebrow' }, page.section),
    Markdown(slug),
    PageFooter(slug)
  )
}
