#!/bin/sh
# Render the positioning visual reference.
#
# 1. Build STLs by running ./main.go. Each scene also writes a sibling
#    <name>.flags file with the per-scene camera flags so f3d can put
#    the camera on the side that shows the indicator pin.
# 2. f3d each STL → docs/images/positioning-<name>.png.
#
# env -i / HOME=/tmp scrubs any local f3d config that survives --no-config.

set -eu

ROOT=$(cd "$(dirname "$0")/../../.." && pwd)
IMG_DIR="$ROOT/docs/images"
STL_DIR="$ROOT/docs/stl"
WORK=$(mktemp -d)
trap 'rm -rf "$WORK"' EXIT

mkdir -p "$IMG_DIR" "$STL_DIR"

# Triangle target for the in-browser viewer-grade STL.
VIEWER_TRIS=12000

if ! command -v f3d >/dev/null 2>&1; then
    echo "f3d not found on PATH. Install from https://f3d.app." >&2
    exit 1
fi

echo "Building STLs in $WORK"
(cd "$ROOT" && go run ./docs/visuals/positioning "$WORK") >/dev/null

BASE_FLAGS="--no-config --up=+Z --no-background=1 --anti-aliasing ssaa --tone-mapping --resolution 1000,750 --color #b6bdc9 --ambient-occlusion"

count=0

for stl in "$WORK"/anchor-*.stl "$WORK"/verb-*.stl; do
    name=$(basename "${stl%.stl}")
    png="$IMG_DIR/positioning-$name.png"
    cam=$(cat "${stl%.stl}.flags")

    env -i HOME=/tmp PATH="$PATH" f3d "$stl" \
        --output "$png" \
        $BASE_FLAGS $cam \
        >/dev/null 2>&1

    # Decimate for the in-browser STL viewer.
    decimated="$STL_DIR/positioning-$name.stl"
    (cd "$ROOT" && go run ./tools/decimate-stl "$stl" "$decimated" "$VIEWER_TRIS") >/dev/null 2>&1 || cp "$stl" "$decimated"

    echo "  positioning-$name.png + .stl"
    count=$((count + 1))
done

# Refresh the manifest from whatever STLs are on disk — same logic as
# tutorial/render-screenshots.sh so a partial re-render doesn't drop the
# rest of the entries.
{
    echo "["
    (cd "$STL_DIR" && ls *.stl 2>/dev/null) | sed 's/\.stl$//' | sort -u | sed 's/.*/  "&"/' | paste -sd, -
    echo "]"
} > "$STL_DIR/manifest.json"

echo "Rendered $count positioning visuals to $IMG_DIR (+ decimated STLs in $STL_DIR)"
