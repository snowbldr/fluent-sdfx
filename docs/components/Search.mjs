// Client-side fuzzy search over the docs.
//
// On first focus, fetches every page's markdown in parallel, strips it to
// plain text + headings, and builds a flat record list. Each keystroke
// re-scores against the records: title hits weight more than heading hits,
// which weight more than body hits; an unbroken substring weights more
// than letter-soup. Top 8 results are rendered into a popover; clicking
// (or Enter) navigates via fnroute.
//
// Built on fntags' reactive state (fnstate) — same idiom as the rest of
// the components.

import { div, input, ul, li, button, mark, span } from '@srfnstack/fntags'
import { fnstate } from '@srfnstack/fntags'
import { goTo } from '@srfnstack/fntags/src/fnroute.mjs'
import { allPages } from '../nav.mjs'
import { BASE } from '../constants.mjs'

const RESULT_LIMIT = 8

let recordsPromise = null

const stripMarkdown = (md) =>
  md
    .replace(/^---\n[\s\S]*?\n---\n/, '')
    .replace(/```[\s\S]*?```/g, ' ')
    .replace(/`[^`]+`/g, m => m.slice(1, -1))
    .replace(/<[^>]+>/g, ' ')
    .replace(/[#*_>~]/g, ' ')
    .replace(/\s+/g, ' ')
    .toLowerCase()

const extractHeadings = (md) =>
  Array.from(md.matchAll(/^(#{1,3})\s+(.+)$/gm)).map(m => ({
    level: m[1].length,
    text: m[2].trim()
  }))

// Slug-style anchor — matches marked's GitHub-style heading IDs.
const headingAnchor = (text) =>
  text.toLowerCase().replace(/[^\w\s-]/g, '').trim().replace(/\s+/g, '-')

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
            text: stripMarkdown(md),
            headings: extractHeadings(md)
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

function scoreRecord(rec, q) {
  const inTitle = rec.title.toLowerCase().includes(q)
  const inSection = (rec.section || '').toLowerCase().includes(q)
  const headingHit = rec.headings.find(h => h.text.toLowerCase().includes(q))
  const bodyIdx = rec.text.indexOf(q)
  if (!inTitle && !headingHit && bodyIdx < 0 && !inSection) return null

  let score = 0
  if (inTitle) score += 100
  if (inSection) score += 20
  if (headingHit) score += 60 - 5 * Math.min(headingHit.level, 3)
  if (bodyIdx >= 0) score += 30 + Math.max(0, 20 - Math.floor(bodyIdx / 100))

  const idx = rec.text.indexOf(q)
  let excerpt = ''
  if (idx >= 0) {
    const start = Math.max(0, idx - 40)
    const end = Math.min(rec.text.length, idx + q.length + 80)
    excerpt = rec.text.slice(start, end)
    if (start > 0) excerpt = '…' + excerpt
    if (end < rec.text.length) excerpt += '…'
  }

  return {
    score,
    rec,
    heading: headingHit?.text,
    anchor: headingHit ? headingAnchor(headingHit.text) : null,
    excerpt,
    query: q,
    // bindChildren needs a stable key per item.
    key: rec.slug + (headingHit ? '#' + headingHit.text : '')
  }
}

function search(records, q) {
  const norm = q.trim().toLowerCase()
  if (!norm) return []
  const hits = []
  for (const r of records) {
    const h = scoreRecord(r, norm)
    if (h) hits.push(h)
  }
  hits.sort((a, b) => b.score - a.score)
  return hits.slice(0, RESULT_LIMIT)
}

// Render a string with the matched substring wrapped in <mark>.
function highlight(text, q) {
  if (!q) return text
  const idx = text.toLowerCase().indexOf(q)
  if (idx < 0) return text
  return span({},
    text.slice(0, idx),
    mark({}, text.slice(idx, idx + q.length)),
    text.slice(idx + q.length)
  )
}

export const Search = () => {
  const query = fnstate('')
  const results = fnstate([], r => r.key)
  const cursor = fnstate(-1)
  const focused = fnstate(false)
  const showResults = fnstate(false)
  let recordsLoaded = false

  // Show results only when there are some AND the input has focus (so clicks
  // outside dismiss the popover). bindStyle on the wrapper drives display.
  const recompute = () => showResults(results().length > 0 && focused())
  results.subscribe(recompute)
  focused.subscribe(recompute)

  const ensureLoaded = async () => {
    if (recordsLoaded) return
    recordsLoaded = true
    const records = await loadRecords()
    if (query()) results(search(records, query()))
    return records
  }

  query.subscribe(async (q) => {
    if (!q) { results([]); cursor(-1); return }
    const records = await loadRecords()
    results(search(records, q))
    cursor(-1)
  })

  const inputEl = input({
    type: 'search',
    class: 'search-input',
    placeholder: 'Search docs',
    'aria-label': 'Search docs',
    autocomplete: 'off',
    spellcheck: 'false',
    value: query.bindAttr(),
    onfocus: () => { focused(true); ensureLoaded() },
    // Delay so a click on a result fires before we close the dropdown.
    onblur: () => setTimeout(() => focused(false), 200),
    oninput: (e) => query(e.target.value)
  })

  const navigate = (hit) => {
    const path = hit.anchor ? hit.rec.path + '#' + hit.anchor : hit.rec.path
    query('')
    inputEl.blur()
    goTo(path)
    if (hit.anchor) {
      setTimeout(() => {
        document.getElementById(hit.anchor)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
      }, 100)
    }
  }

  const ResultItem = (hitState) =>
    li({
      class: cursor.bindAttr(() => {
        const i = results().findIndex(h => h.key === hitState().key)
        return i === cursor() ? 'search-result active' : 'search-result'
      })
    },
      button({
        type: 'button',
        class: 'search-result-btn',
        // mousedown fires before the input's blur, preventing the dropdown
        // from disappearing before the click registers.
        onmousedown: (e) => { e.preventDefault(); navigate(hitState()) }
      },
        div({ class: 'search-result-title' },
          hitState.bindAs(() => highlight(hitState().rec.title, hitState().query))
        ),
        hitState.bindAs(() =>
          hitState().heading
            ? div({ class: 'search-result-heading' },
                highlight(hitState().heading, hitState().query))
            : div({ hidden: true })
        ),
        hitState.bindAs(() =>
          hitState().excerpt
            ? div({ class: 'search-result-excerpt' },
                highlight(hitState().excerpt, hitState().query))
            : div({ hidden: true })
        ),
        div({ class: 'search-result-section' },
          hitState.bindAs(() => hitState().rec.section || '')
        )
      )
    )

  const resultsEl = results.bindChildren(
    ul({
      class: 'search-results',
      hidden: showResults.bindAttr(() => showResults() ? null : 'hidden')
    }),
    ResultItem
  )

  inputEl.addEventListener('keydown', (e) => {
    const list = results()
    if (e.key === 'Escape') {
      query('')
      inputEl.blur()
      return
    }
    if (!list.length) return
    if (e.key === 'ArrowDown') {
      e.preventDefault()
      cursor(Math.min(list.length - 1, cursor() + 1))
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()
      cursor(Math.max(0, cursor() - 1))
    } else if (e.key === 'Enter') {
      e.preventDefault()
      navigate(list[cursor() >= 0 ? cursor() : 0])
    }
  })

  // Cmd/Ctrl-K from anywhere focuses the input.
  document.addEventListener('keydown', (e) => {
    if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
      e.preventDefault()
      inputEl.focus()
      inputEl.select()
    }
  })

  return div({ class: 'search' }, inputEl, resultsEl)
}
