#!/usr/bin/env python3
"""Local dev server for the fluent-sdfx docs.

Serves the docs/ directory and falls back to docs/404.html with a 404
status when the requested path isn't a real file — matching how GitHub
Pages behaves. Without this, deep-link reloads (e.g. /quickstart/) would
show Python's plain 404 instead of triggering our SPA recovery.
"""

import http.server
import os
import socketserver
import sys
from pathlib import Path
from urllib.parse import urlparse

DOCS = Path(__file__).resolve().parent.parent / "docs"
PORT = int(sys.argv[1]) if len(sys.argv) > 1 else 7174


class Handler(http.server.SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory=str(DOCS), **kwargs)

    def do_GET(self):
        url = urlparse(self.path)
        path = (DOCS / url.path.lstrip("/")).resolve()
        if path.is_dir():
            path = path / "index.html"
        if not (path.is_file() and DOCS in path.parents or path == DOCS):
            # Not a real file — serve 404.html with 404 status (mimics GH Pages).
            fb = DOCS / "404.html"
            if fb.is_file():
                self.send_response(404)
                self.send_header("Content-Type", "text/html")
                self.end_headers()
                self.wfile.write(fb.read_bytes())
                return
        super().do_GET()


def main():
    os.chdir(DOCS)
    print(f"Serving {DOCS} on http://localhost:{PORT}/", flush=True)
    with socketserver.TCPServer(("", PORT), Handler) as httpd:
        try:
            httpd.serve_forever()
        except KeyboardInterrupt:
            pass


if __name__ == "__main__":
    main()
