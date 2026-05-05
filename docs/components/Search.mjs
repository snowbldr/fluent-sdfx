// Client-side fuzzy search over the docs.
//
// On first focus, fetches every page's markdown in parallel, strips it to
// plain text + headings, and builds a flat record list. Each keystroke
// re-scores against the records: title hits weight more than heading hits,
// which weight more than body hits; an unbroken substring weights more
// than letter-soup. Top 8 results are rendered into a popover; clicking a
// result navigates via fnroute (so the SPA stays in-app).
//
// No external library — keeps payload small and lets us tune ranking
// for fluent-sdfx's vocabulary (e.g., method names beat surrounding prose).

import { div, input, ul, li, button } from '@srfnstack/fntags'
import { goTo } from '@srfnstack/fntags/src/fnroute.mjs'
import { allPages } from '../nav.mjs'
import { BASE } from '../constants.mjs'

let recordsPromise = null

// Strip HTML tags, code-fence markers, and front-matter from raw markdown so
// queries match the actual prose, not the syntax.
function plainText(md) {
  return md
    .replace(/^---\n[\s\S]*?\n---\n/, '')
    .replace(/```[\s\S]*?```/g, ' ')
    .replace(/`[^`]+`/g, m => m.slice(1, -1))
    .replace(/<[^>]+>/g, ' ')
    .replace(/[#*_>~]/g, ' ')
    .replace(/\s+/g, ' ')
    .toLowerCase()
}

// Pull headings (## H2, ### H3) so we can show a meaningful subtitle on hits.
function headings(md) {
  const out = []
  for (const m of md.matchAll(/^(#{1,3})\s+(.+)$/gm)) {
    out.push({ level: m[1].length, text: m[2].trim() })
  }
  return out
}

// Find the slug-style anchor for a heading — marked uses lowercased,
// hyphen-joined IDs, matching GitHub's algorithm closely enough for our
// uses. Strips punctuation; collapses spaces to hyphens.
function headingAnchor(text) {
  return text
    .toLowerCase()
    .replace(/[^\w\s-]/g, '')
    .trim()
    .replace(/\s+/g, '-')
}

async function loadRecords() {
  if (recordsPromise) return recordsPromise
  recordsPromise = (async () => {
    const out = []
    await Promise.all(
      allPages.map(async (p) => {
        try {
          const r = await fetch(BASE + 'content/' + p.slug + '.md')
          if (!r.ok) return
          const md = await r.text()
          out.push({
            slug: p.slug,
            path: p.path,
            title: p.title,
            section: p.section,
            text: plainText(md),
            headings: headings(md)
          })
        } catch {
          // Skip failed pages silently — search degrades gracefully.
        }
      })
    )
    return out
  })()
  return recordsPromise
}

// Score a record against a normalized query. Higher = better.
function scoreRecord(rec, q) {
  if (!q) return 0
  const inTitle = rec.title.toLowerCase().includes(q)
  const inSection = (rec.section || '').toLowerCase().includes(q)
  const headingHit = rec.headings.find(h => h.text.toLowerCase().includes(q))
  const bodyIdx = rec.text.indexOf(q)

  if (!inTitle && !headingHit && bodyIdx < 0 && !inSection) return 0

  let score = 0
  if (inTitle) score += 100
  if (inSection) score += 20
  if (headingHit) score += 60 - 5 * Math.min(headingHit.level, 3)
  if (bodyIdx >= 0) score += 30 + Math.max(0, 20 - Math.floor(bodyIdx / 100))

  return score
}

// Build a 120-char excerpt around the first body match.
function excerpt(rec, q) {
  if (!q) return ''
  const idx = rec.text.indexOf(q)
  if (idx < 0) return ''
  const start = Math.max(0, idx - 40)
  const end = Math.min(rec.text.length, idx + q.length + 80)
  let s = rec.text.slice(start, end)
  if (start > 0) s = '…' + s
  if (end < rec.text.length) s += '…'
  return s
}

function search(records, query) {
  const q = query.trim().toLowerCase()
  if (!q) return []
  const hits = []
  for (const r of records) {
    const s = scoreRecord(r, q)
    if (s > 0) {
      const heading = r.headings.find(h => h.text.toLowerCase().includes(q))
      hits.push({
        rec: r,
        score: s,
        heading: heading?.text,
        anchor: heading ? headingAnchor(heading.text) : null,
        excerpt: excerpt(r, q)
      })
    }
  }
  hits.sort((a, b) => b.score - a.score)
  return hits.slice(0, 8)
}

function highlight(text, q) {
  if (!q) return text
  const idx = text.toLowerCase().indexOf(q)
  if (idx < 0) return text
  const before = text.slice(0, idx)
  const match = text.slice(idx, idx + q.length)
  const after = text.slice(idx + q.length)
  const span = document.createElement('span')
  span.appendChild(document.createTextNode(before))
  const mark = document.createElement('mark')
  mark.textContent = match
  span.appendChild(mark)
  span.appendChild(document.createTextNode(after))
  return span
}

function ResultLi({ hit, onPick }) {
  const q = (hit._q || '').toLowerCase()
  const children = [div({ class: 'search-result-title' })]
  if (hit.heading) children.push(div({ class: 'search-result-heading' }))
  if (hit.excerpt) children.push(div({ class: 'search-result-excerpt' }))
  children.push(div({ class: 'search-result-section' }, hit.rec.section || ''))

  const node = li({ class: 'search-result' },
    button({
      type: 'button',
      class: 'search-result-btn',
      onclick: () => onPick(hit)
    }, ...children)
  )
  // Fill highlighted text in after construction so we can use real DOM nodes
  // (mark elements) rather than escaped strings.
  node.querySelector('.search-result-title').appendChild(highlight(hit.rec.title, q))
  if (hit.heading) {
    node.querySelector('.search-result-heading').appendChild(highlight(hit.heading, q))
  }
  if (hit.excerpt) {
    node.querySelector('.search-result-excerpt').appendChild(highlight(hit.excerpt, q))
  }
  return node
}

export const Search = () => {
  const inputEl = input({
    type: 'search',
    class: 'search-input',
    placeholder: 'Search docs',
    'aria-label': 'Search docs',
    autocomplete: 'off',
    spellcheck: 'false'
  })

  const resultsEl = ul({ class: 'search-results', hidden: true })
  let activeRecords = null
  let lastHits = []
  let cursor = -1

  const close = () => {
    resultsEl.hidden = true
    cursor = -1
  }

  const navigate = (hit) => {
    const path = hit.anchor
      ? hit.rec.path + '#' + hit.anchor
      : hit.rec.path
    close()
    inputEl.value = ''
    inputEl.blur()
    goTo(path)
    // fnroute won't scroll to the hash on its own — do it after a tick.
    if (hit.anchor) {
      setTimeout(() => {
        const target = document.getElementById(hit.anchor)
        if (target) target.scrollIntoView({ behavior: 'smooth', block: 'start' })
      }, 100)
    }
  }

  const render = (q) => {
    const hits = search(activeRecords, q)
    hits.forEach(h => { h._q = q })
    resultsEl.innerHTML = ''
    if (!hits.length) {
      resultsEl.hidden = true
      lastHits = []
      cursor = -1
      return
    }
    for (const h of hits) {
      resultsEl.appendChild(ResultLi({ hit: h, onPick: navigate }))
    }
    resultsEl.hidden = false
    lastHits = hits
    cursor = -1
  }

  const setCursor = (i) => {
    cursor = i
    for (const [idx, el] of resultsEl.children.entries()) {
      el.classList.toggle('active', idx === i)
    }
    if (i >= 0) {
      resultsEl.children[i]?.scrollIntoView({ block: 'nearest' })
    }
  }

  const ensureLoaded = async () => {
    if (activeRecords) return
    activeRecords = await loadRecords()
  }

  inputEl.addEventListener('focus', ensureLoaded)

  inputEl.addEventListener('input', async () => {
    await ensureLoaded()
    render(inputEl.value)
  })

  inputEl.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
      close()
      inputEl.blur()
      return
    }
    if (resultsEl.hidden || !lastHits.length) return
    if (e.key === 'ArrowDown') {
      e.preventDefault()
      setCursor(Math.min(lastHits.length - 1, cursor + 1))
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()
      setCursor(Math.max(0, cursor - 1))
    } else if (e.key === 'Enter') {
      e.preventDefault()
      const idx = cursor >= 0 ? cursor : 0
      navigate(lastHits[idx])
    }
  })

  // Click outside closes the dropdown.
  const wrapper = div({ class: 'search' }, inputEl, resultsEl)
  document.addEventListener('click', (e) => {
    if (!wrapper.contains(e.target)) close()
  })

  // Cmd/Ctrl-K to focus.
  document.addEventListener('keydown', (e) => {
    if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
      e.preventDefault()
      inputEl.focus()
      inputEl.select()
    }
  })

  return wrapper
}
