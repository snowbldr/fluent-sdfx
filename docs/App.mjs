import { div, header, main, span } from '@srfnstack/fntags'
import { route, routeSwitch, setRootPath } from '@srfnstack/fntags/src/fnroute.mjs'
import { allPages } from './nav.mjs'
import { BASE } from './constants.mjs'
import { Sidebar } from './components/Sidebar.mjs'
import { ThemeToggle } from './components/ThemeToggle.mjs'
import { PageView } from './components/PageView.mjs'

// fnroute initialises rootPath from window.location.pathname, which on a
// deep link makes everything relative to that — collapsing all our routes
// onto the same path. Pin rootPath to the detected base (see
// constants.mjs) so absolute routes mean what we wrote.
setRootPath(BASE)

// One <route> per known path. Each renders the matching markdown page.
const routes = allPages.map(p =>
  route({ path: p.path, absolute: 'true' }, PageView(p.slug))
)

// Fallback for unknown URLs. Non-absolute "/" matches anything, but is
// last so specific routes win first.
routes.push(
  route({ path: '/' },
    div({ class: 'page-missing' },
      span({}, 'Page not found.')
    )
  )
)

export const App = () =>
  div({ class: 'app' },
    Sidebar(),
    div({ class: 'main-col' },
      header({ class: 'topbar' },
        div({ class: 'topbar-spacer' }),
        ThemeToggle()
      ),
      main({ class: 'content' },
        routeSwitch(...routes)
      )
    )
  )
