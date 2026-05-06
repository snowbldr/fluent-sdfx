#!/usr/bin/env bash
# Local coverage runner for fluent-sdfx.
#
# Runs `go test` with coverage across all packages, then strips lines
# belonging to non-kernel packages (examples, tutorial, docs, meta, tools)
# from the resulting profile so the reported percentage reflects the
# kernel surface only — matching what Codecov reports in CI.
#
# Outputs:
#   coverage.out   — filtered Go coverage profile
#   coverage.html  — browsable HTML report
# Prints overall kernel coverage % to stdout.

set -euo pipefail

# Run from the repo root regardless of where the script was invoked.
cd "$(dirname "$0")/.."

PROFILE="coverage.out"
RAW_PROFILE="coverage.raw.out"
HTML="coverage.html"

MODULE="github.com/snowbldr/fluent-sdfx"

echo "==> running go test with coverage"
go test -coverprofile="$RAW_PROFILE" -covermode=atomic ./...

echo "==> filtering non-kernel packages from profile"
# Keep the `mode:` header line and any line whose first whitespace-delimited
# field (the file path) does NOT start with one of the excluded prefixes.
awk -v mod="$MODULE" '
  NR == 1 && /^mode:/ { print; next }
  {
    path = $1
    if (path ~ "^" mod "/examples/")      next
    if (path ~ "^" mod "/tutorial/")      next
    if (path ~ "^" mod "/docs/")          next
    if (path ~ "^" mod "/meta/")          next
    if (path ~ "^" mod "/tools/")         next
    print
  }
' "$RAW_PROFILE" > "$PROFILE"

rm -f "$RAW_PROFILE"

echo "==> generating HTML report"
go tool cover -html="$PROFILE" -o "$HTML"

echo "==> overall kernel coverage"
# `go tool cover -func` prints a `total:` line at the bottom with the
# aggregate percentage across the (filtered) profile.
go tool cover -func="$PROFILE" | tail -n 1
