import { button } from '@srfnstack/fntags'

const KEY = 'fluent-sdfx-theme'

function current() {
  return document.documentElement.getAttribute('data-theme') || 'dark'
}

function set(theme) {
  document.documentElement.setAttribute('data-theme', theme)
  try {
    localStorage.setItem(KEY, theme)
  } catch {}
}

// Restore on load.
try {
  const saved = localStorage.getItem(KEY)
  if (saved) set(saved)
} catch {}

export const ThemeToggle = () =>
  button({
    class: 'theme-toggle',
    type: 'button',
    'aria-label': 'Toggle theme',
    onclick: () => set(current() === 'dark' ? 'light' : 'dark')
  }, '◐')
