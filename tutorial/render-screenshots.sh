#!/bin/sh
# Render docs assets for every tutorial step.
#
# For each tutorial/<page>/<step>/ directory:
#   1. go run produces out.stl (3D) or out.png (2D)
#   2. f3d renders a hero PNG → docs/images/<page>-<step>.png
#   3. The STL is decimated for the in-browser 3D viewer →
#      docs/stl/<page>-<step>.stl (a few hundred KB).
#
# Per-step f3d.flags or per-page _camera.flags overrides the default f3d
# camera string.

set -eu

ROOT=$(cd "$(dirname "$0")/.." && pwd)
IMG_DIR="$ROOT/docs/images"
STL_DIR="$ROOT/docs/stl"
mkdir -p "$IMG_DIR" "$STL_DIR"

if ! command -v f3d >/dev/null 2>&1; then
	echo "f3d not found on PATH. Install from https://f3d.app." >&2
	exit 1
fi

F3D_DEFAULTS="--no-config --up=+Z --no-background=1 --camera-elevation-angle 25 --camera-azimuth-angle 40 --anti-aliasing ssaa --ambient-occlusion --tone-mapping --resolution 1200,900"

# Triangle target for the viewer-grade STL. Big enough that even an
# array of 20 spheres reads smoothly. Each sphere ends up with several
# hundred triangles; total docs/stl/ footprint lands around 25 MB.
VIEWER_TRIS=12000

# Filter: optional first argument is a glob like "04-quickstart" to limit scope.
FILTER=${1:-}

for stepdir in "$ROOT"/tutorial/*/*/; do
	page=$(basename "$(dirname "$stepdir")")
	step=$(basename "$stepdir")

	# Skip non-step subdirectories (e.g. tutorial/internal/...).
	if [ ! -f "$stepdir/main.go" ]; then
		continue
	fi
	# Skip step dirs whose main.go isn't a `package main` (helpers).
	if ! grep -q '^package main$' "$stepdir/main.go"; then
		continue
	fi

	if [ -n "$FILTER" ] && [ "$page" != "$FILTER" ]; then
		continue
	fi

	stem="${page}-${step}"
	img_out="$IMG_DIR/${stem}.png"
	stl_out="$STL_DIR/${stem}.stl"
	echo "→ $page/$step"

	if ! (cd "$stepdir" && go run . >/dev/null); then
		echo "  build failed, skipping" >&2
		continue
	fi

	if [ -f "$stepdir/out.stl" ]; then
		flags="$F3D_DEFAULTS"
		if [ -f "$ROOT/tutorial/$page/_camera.flags" ]; then
			flags=$(cat "$ROOT/tutorial/$page/_camera.flags")
		fi
		if [ -f "$stepdir/f3d.flags" ]; then
			flags=$(cat "$stepdir/f3d.flags")
		fi
		# shellcheck disable=SC2086
		f3d "$stepdir/out.stl" --output "$img_out" $flags >/dev/null

		# Decimate for the in-browser viewer.
		(cd "$ROOT" && go run ./tools/decimate-stl "$stepdir/out.stl" "$stl_out" "$VIEWER_TRIS") >/dev/null 2>&1 || {
			echo "  decimation failed; copying full STL" >&2
			cp "$stepdir/out.stl" "$stl_out"
		}

		rm -f "$stepdir/out.stl"
	elif [ -f "$stepdir/out.png" ]; then
		cp "$stepdir/out.png" "$img_out"
		rm -f "$stepdir/out.png"
	else
		echo "  produced no out.stl or out.png" >&2
	fi
done

# Emit a manifest of viewer-available STLs from whatever's on disk. Building
# this from `ls` (rather than tracking during the loop) means a filtered
# rerun like `./render-screenshots.sh 11-smooth-blends` only refreshes that
# page's STLs without dropping the rest from the manifest.
{
	echo "["
	(cd "$STL_DIR" && ls *.stl 2>/dev/null) | sed 's/\.stl$//' | sort -u | sed 's/.*/  "&"/' | paste -sd, -
	echo "]"
} > "$STL_DIR/manifest.json"

echo "Wrote PNGs to $IMG_DIR and viewer STLs to $STL_DIR"
