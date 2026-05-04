import { goTo } from '@srfnstack/fntags/src/fnroute.mjs'
import { App } from './App.mjs'

document.body.appendChild(App())

// GitHub Pages SPA deep-link recovery.
//
// When a user requests /fluent-sdfx/quickstart/ on GH Pages, the server
// 404's because there's no quickstart/index.html. Our public/404.html
// captures location.pathname into sessionStorage.redirect and refreshes
// to the app root. Here we read the stash and navigate the SPA there.
//
// The replace() strips the BASE_URL prefix from the redirect, leaving
// the in-app path that fnroute will re-prepend rootPath to.
if (sessionStorage.redirect) {
  const stashed = sessionStorage.redirect
  delete sessionStorage.redirect
  setTimeout(() => {
    const base = location.pathname.split(/[#?]/)[0]
    goTo(stashed.replace(base, ''))
  }, 1)
}
