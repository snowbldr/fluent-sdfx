#!/bin/sh
# Install (or re-install) the project's git hooks.
#
# Run:   make install-hooks
# Or:    ./tools/install-hooks.sh
#
# The pre-commit hook runs `make build` and re-stages any files the build
# changes (e.g. docs/content/*.md picking up new tutorial Go code).

set -e
ROOT=$(cd "$(dirname "$0")/.." && pwd)
HOOK_DIR="$ROOT/.git/hooks"

if [ ! -d "$HOOK_DIR" ]; then
	echo "Not a git repo: $ROOT/.git/hooks doesn't exist." >&2
	exit 1
fi

cp "$ROOT/tools/hooks/pre-commit" "$HOOK_DIR/pre-commit"
chmod +x "$HOOK_DIR/pre-commit"

echo "Installed: $HOOK_DIR/pre-commit"
