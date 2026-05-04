// Detect the docs site base URL from the module's own location.
//
// Works for any hosting setup without configuration:
//   - GitHub Pages    https://user.github.io/fluent-sdfx/  → '/fluent-sdfx/'
//   - IntelliJ server localhost:63342/fluent-sdfx/docs/    → '/fluent-sdfx/docs/'
//   - Root deploys    localhost:7174/                       → '/'
//   - Custom domains  https://docs.example.com/            → '/'
//
// We anchor on import.meta.url — the absolute URL this module was loaded
// from. The directory part of the pathname is the base.

const moduleUrl = new URL(import.meta.url)
export const BASE = moduleUrl.pathname.replace(/\/[^/]*$/, '/')
