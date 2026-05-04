#!/bin/sh
# Sync code blocks in docs/content/*.md with their canonical Go source.
#
# For each `<!-- src: PATH -->` HTML comment in any markdown file,
# replace the body of the immediately-following ```go fenced block with
# the contents of <repo-root>/PATH. This keeps the docs in lock-step with
# tutorial code without a build step.
#
# Run it directly, or via `make build`, or via the pre-commit hook
# installed by tools/install-hooks.sh.

set -e
exec python3 "$(cd "$(dirname "$0")" && pwd)/update-tutorial-code.py" "$@"
