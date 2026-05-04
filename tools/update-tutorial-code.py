#!/usr/bin/env python3
"""Replace `<!-- src: PATH -->` code blocks in docs/content/*.md with the
on-disk contents of <repo-root>/PATH.

Idempotent: running it twice on a clean tree changes nothing.
"""

import re
import sys
from pathlib import Path

REPO = Path(__file__).resolve().parent.parent
CONTENT = REPO / "docs" / "content"

# Match: <!-- src: PATH -->\n```go\n...body...\n```
PATTERN = re.compile(
    r"(<!-- src: ([^\s]+) -->\n```go\n)(?:.*?)(\n```)",
    re.DOTALL,
)


def update(md_path):
    text = md_path.read_text()
    changed = False
    missing = []

    def repl(match):
        nonlocal changed
        prefix, src_rel, suffix = match.group(1), match.group(2), match.group(3)
        src_abs = REPO / src_rel
        if not src_abs.is_file():
            missing.append(src_rel)
            return match.group(0)
        body = src_abs.read_text().rstrip()
        new = prefix + body + suffix
        if new != match.group(0):
            changed = True
        return new

    new_text = PATTERN.sub(repl, text)
    if changed:
        md_path.write_text(new_text)
        print(f"updated: {md_path.relative_to(REPO)}")
    if missing:
        for m in missing:
            print(f"  missing source: {m}", file=sys.stderr)
    return bool(missing)


def main():
    any_missing = False
    for md in sorted(CONTENT.glob("*.md")):
        any_missing |= update(md)
    return 1 if any_missing else 0


if __name__ == "__main__":
    sys.exit(main())
