import { button, div, span } from '@srfnstack/fntags'
import { fnstate } from '@srfnstack/fntags'
import { pathState } from '@srfnstack/fntags/src/fnroute.mjs'

// On narrow screens the sidebar is a drawer. navOpen drives a class on
// .app; the CSS only honours it below the mobile breakpoint, so on a
// desktop the state is inert.
export const navOpen = fnstate(false)

// Following a link closes the drawer, so the page it opened is visible.
pathState.subscribe(() => navOpen(false))

export const NavToggle = () =>
  button({
    class: 'nav-toggle',
    type: 'button',
    'aria-label': 'Open navigation',
    'aria-expanded': navOpen.bindAttr(() => String(navOpen())),
    onclick: () => navOpen(!navOpen())
  },
    span({ class: 'nav-toggle-bar' }),
    span({ class: 'nav-toggle-bar' }),
    span({ class: 'nav-toggle-bar' })
  )

// The dimmed page behind the open drawer; tapping it closes the drawer.
export const NavBackdrop = () =>
  div({ class: 'nav-backdrop', onclick: () => navOpen(false) })
