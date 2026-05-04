// Fetch and render a markdown file at runtime.
//
//   - fetches `${BASE}content/<slug>.md`
//   - parses with `marked` (loaded via importmap)
//   - extends marked with a tokenizer for GitHub alerts
//     (`> [!TIP]\n> body` → <aside class="callout callout-tip">…</aside>)
//   - rewrites image src from `../images/foo` to `${BASE}images/foo`
//     so paths work both locally and under a GH Pages base
//   - hands the result to Prism for syntax highlighting

import { div } from '@srfnstack/fntags'
import { goTo } from '@srfnstack/fntags/src/fnroute.mjs'
import { marked } from 'marked'
import { BASE } from '../constants.mjs'
import { attachViewer } from './StlViewer.mjs'

// Manifest of figures with a viewer-grade STL on disk. Loaded once, cached.
let stlSetPromise = null
function loadStlManifest() {
  if (!stlSetPromise) {
    stlSetPromise = fetch(BASE + 'stl/manifest.json')
      .then(r => (r.ok ? r.json() : []))
      .then(list => new Set(list))
      .catch(() => new Set())
  }
  return stlSetPromise
}

const ALERT_VARIANTS = {
  tip: 'tip',
  warning: 'warning',
  note: 'note',
  important: 'note',
  caution: 'warning'
}

marked.use({
  extensions: [
    {
      name: 'githubAlert',
      level: 'block',
      start(src) {
        const m = src.match(/^>\s*\[!(TIP|WARNING|NOTE|IMPORTANT|CAUTION)\]/m)
        return m ? src.indexOf(m[0]) : undefined
      },
      tokenizer(src) {
        const open = src.match(
          /^>\s*\[!(TIP|WARNING|NOTE|IMPORTANT|CAUTION)\]\s*\n/
        )
        if (!open) return
        // Walk subsequent lines: each `> ...` line that ISN'T another
        // alert opener belongs to this alert's body. Stop at the first
        // non-quoted line OR the next `> [!ALERT]` line so adjacent
        // alerts don't bleed into each other.
        const after = src.slice(open[0].length)
        const lines = after.split('\n')
        const bodyLines = []
        let consumed = open[0].length
        for (const line of lines) {
          if (!line.startsWith('>')) break
          if (/^>\s*\[!(TIP|WARNING|NOTE|IMPORTANT|CAUTION)\]/.test(line)) break
          bodyLines.push(line)
          consumed += line.length + 1 // +1 for the \n
        }
        const variant = ALERT_VARIANTS[open[1].toLowerCase()] || 'note'
        const body = bodyLines.map(l => l.replace(/^>\s?/, '')).join('\n')
        return {
          type: 'githubAlert',
          raw: src.slice(0, consumed),
          variant,
          tokens: this.lexer.blockTokens(body)
        }
      },
      renderer(token) {
        const inner = this.parser.parse(token.tokens)
        return (
          `<aside class="callout callout-${token.variant}">` +
          `<div class="callout-icon" aria-hidden="true">${calloutIcon(token.variant)}</div>` +
          `<div class="callout-body">${inner}</div>` +
          `</aside>`
        )
      }
    }
  ]
})

function calloutIcon(variant) {
  switch (variant) {
    case 'tip': return '✦'
    case 'warning': return '⚠'
    case 'note': return '◆'
    default: return '•'
  }
}

function rewriteImages(root) {
  for (const img of root.querySelectorAll('img')) {
    const src = img.getAttribute('src') || ''
    // Markdown is authored with `../images/foo` (works when GitHub renders
    // the .md from /docs/content/); rewrite to base-absolute for the app.
    if (src.startsWith('../images/')) {
      img.setAttribute('src', BASE + 'images/' + src.slice('../images/'.length))
    }
  }
}

// For each <figure> whose <img> stem matches an entry in the STL manifest,
// attach a "View in 3D" toggle that swaps the static image for a Three.js
// canvas of the same part. Three.js itself is loaded lazily on first click.
async function attachStlViewers(root) {
  const stems = await loadStlManifest()
  if (!stems.size) return
  for (const figure of root.querySelectorAll('figure')) {
    const img = figure.querySelector('img')
    if (!img) continue
    const m = (img.getAttribute('src') || '').match(/\/images\/([^/]+)\.png$/)
    if (!m || !stems.has(m[1])) continue
    attachViewer(figure, BASE + 'stl/' + m[1] + '.stl')
  }
}

// Rewrite root-relative in-app links (e.g. `/quickstart/`) to base-prefixed
// hrefs and intercept clicks so navigation stays in-SPA — without this,
// clicking [Smooth blends](/smooth-blends/) does a full reload to the
// host root, which 404's under any non-root deploy.
function rewriteLinks(root) {
  for (const a of root.querySelectorAll('a[href]')) {
    const href = a.getAttribute('href')
    if (!href || href.startsWith('#')) continue
    // External links: open in a new tab.
    if (/^[a-z]+:\/\//i.test(href) || href.startsWith('mailto:')) {
      a.setAttribute('target', '_blank')
      a.setAttribute('rel', 'noopener noreferrer')
      continue
    }
    // Authored as absolute in-app links — rewrite to base-prefixed.
    if (href.startsWith('/')) {
      a.setAttribute('href', BASE.replace(/\/$/, '') + href)
    }
    // Intercept the click so fnroute handles navigation in-SPA.
    a.addEventListener('click', (e) => {
      // Honour modifier-clicks (open in new tab, etc).
      if (e.metaKey || e.ctrlKey || e.shiftKey || e.altKey || e.button !== 0) return
      e.preventDefault()
      // goTo expects an in-app path (no base prefix); restore the original.
      goTo(href.startsWith('/') ? href : new URL(a.href, location.href).pathname.replace(BASE.replace(/\/$/, ''), '') || '/')
    })
  }
}

function highlight(root) {
  if (typeof window !== 'undefined' && window.Prism) {
    window.Prism.highlightAllUnder(root)
  }
}

// Tab galleries (.tab-gallery) authored as raw HTML in markdown:
// a header bar of <button data-tab="key"> elements, then a stack of
// <div data-pane="key"> panes. Clicking a tab activates its pane.
function attachTabGalleries(root) {
  for (const gallery of root.querySelectorAll('.tab-gallery')) {
    const tabs = gallery.querySelectorAll('[data-tab]')
    const panes = gallery.querySelectorAll('[data-pane]')
    if (!tabs.length || !panes.length) continue
    const setActive = (key) => {
      tabs.forEach(t => t.classList.toggle('active', t.dataset.tab === key))
      panes.forEach(p => p.classList.toggle('active', p.dataset.pane === key))
    }
    tabs.forEach(t => t.addEventListener('click', () => setActive(t.dataset.tab)))
    setActive(tabs[0].dataset.tab)
  }
}

export const Markdown = (slug) => {
  const root = div({ class: 'markdown' })
  root.innerHTML = '<p class="md-loading">Loading…</p>'

  fetch(BASE + 'content/' + slug + '.md')
    .then(r => {
      if (!r.ok) throw new Error(`fetch ${slug}: ${r.status}`)
      return r.text()
    })
    .then(text => {
      // Skip any leftover front-matter (defensive — migration strips it).
      const body = text.replace(/^---\n[\s\S]*?\n---\n/, '')
      root.innerHTML = marked.parse(body)
      rewriteImages(root)
      rewriteLinks(root)
      attachStlViewers(root)
      attachTabGalleries(root)
      highlight(root)
    })
    .catch(err => {
      root.innerHTML = `<p class="md-error">Failed to load <code>${slug}</code>: ${err.message}</p>`
    })

  return root
}
