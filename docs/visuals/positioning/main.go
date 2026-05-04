// Generates the visual reference for docs/content/positioning.md.
//
// For each anchor selector and placement verb, builds a small scene and
// writes it as an STL to the output directory passed as the first
// argument (or ./build by default). Per-scene camera flags are written
// alongside each STL as <name>.flags so the render script can point f3d
// at the anchor from a direction that always shows it.
//
//	go run ./docs/visuals/positioning <out-dir>
//
// Conventions:
//   - Host box: 24×16×12, asymmetric so X/Y/Z faces are visually
//     distinct, with TOP/BOTTOM/RIGHT/LEFT/BACK/FRONT embossed (cut)
//     into each face. Each label is oriented so it reads correctly when
//     the camera looks at that face from outside.
//   - Anchor indicator: a ball at the anchor + a thin pin extending
//     outward along the anchor direction, so the anchor location is
//     unambiguous.
//   - Per-anchor camera: each anchor's scene is rendered with the
//     camera placed roughly along the anchor direction, so the anchor
//     and its visible face labels are always in view.
//   - Placement verbs: a host box plus a sphere mover. Sphere shape
//     contrasts with the box host so the eye separates them
//     immediately.
package main

import (
	_ "embed"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/snowbldr/fluent-sdfx/render"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	sdfrender "github.com/snowbldr/sdfx/render"
)

const (
	hostX, hostY, hostZ = 24.0, 16.0, 12.0
	indicatorR          = 1.1
	pinR                = 0.5
	pinLen              = 5.5
	labelHeight         = 3.5
	labelDepth          = 0.55
	cellsPerMM          = 14.0
)

//go:embed cmr10.ttf
var cmr10 []byte

var font *shape.Font

func init() {
	tmp, err := os.CreateTemp("", "positioning-vis-*.ttf")
	if err != nil {
		panic(err)
	}
	defer tmp.Close()
	if _, err := tmp.Write(cmr10); err != nil {
		panic(err)
	}
	font = shape.LoadFont(tmp.Name())
}

// faceLabel builds a 3D text solid sized to be CUT from the host. The
// shape of the text (as a 2D shape) is pre-rotated so the resulting
// embossed text reads correctly when the face is viewed from outside.
//
// The text is extruded with extra depth and positioned with its outer
// face beyond the host surface, so subtracting it leaves a clean
// embossed impression no matter the wall thickness.
func faceLabel(text, face string) *solid.Solid {
	// Build the 2D text centered on origin.
	t := shape.Text(font, text, labelHeight).Center()
	tag := t.Extrude(labelDepth * 4) // thick extrude → safe cut depth

	// Rotation map per face: bring the readable face of the extruded
	// text to align with the box face's outward normal, with the text
	// reading correctly from outside.
	// Each transform brings the text's frame (reading +X, top +Y,
	// readable face +Z) onto the target face such that, when the camera
	// is OUTSIDE that face looking at it, the letters read normally.
	// Derivation in comments alongside.
	switch face {
	case "top":
		// Camera up=+Y, image right=+X. Identity.
		return tag.TopAt(hostZ/2 + labelDepth)
	case "bottom":
		// The camera nudge for the BOTTOM anchor puts image-right at
		// roughly +X+Y, so we want the text to keep its +X reading
		// direction (rather than mirror it). MirrorXY flips just the
		// Z axis — readable face moves to -Z without touching X/Y.
		return tag.MirrorXY().BottomAt(-hostZ/2 - labelDepth)
	case "right":
		// Camera up=+Z, image right=+Y. Need reading→+Y, top→+Z, readable→+X.
		// Z(+90) sends reading to +Y; Y(+90) lifts the rotated top onto +Z
		// and cycles the readable face onto +X.
		return tag.RotateZ(90).RotateY(90).RightAt(hostX/2 + labelDepth)
	case "left":
		// Camera up=+Z, image right=-Y. Need reading→-Y, top→+Z, readable→-X.
		return tag.RotateZ(-90).RotateY(-90).LeftAt(-hostX/2 - labelDepth)
	case "back":
		// Camera up=+Z, image right=-X. Need reading→-X, top→+Z, readable→+Y.
		// Z(180) flips reading and top; X(-90) lifts top to +Z.
		return tag.RotateZ(180).RotateX(-90).BackAt(hostY/2 + labelDepth)
	case "front":
		// Camera up=+Z, image right=+X. Need reading→+X, top→+Z, readable→-Y.
		// X(90) maps (x,y,z)→(x,-z,y) — single rotation suffices.
		return tag.RotateX(90).FrontAt(-hostY/2 - labelDepth)
	}
	panic("unknown face: " + face)
}

func host() *solid.Solid {
	box := solid.Box(v3.XYZ(hostX, hostY, hostZ), 0.5)
	return box.Cut(
		faceLabel("TOP", "top"),
		faceLabel("BOTTOM", "bottom"),
		faceLabel("RIGHT", "right"),
		faceLabel("LEFT", "left"),
		faceLabel("BACK", "back"),
		faceLabel("FRONT", "front"),
	)
}

func indicator(x, y, z int) *solid.Solid {
	a := solid.Box(v3.XYZ(hostX, hostY, hostZ), 0).AnchorAt(x, y, z).Point
	dir := v3.XYZ(float64(x), float64(y), float64(z))
	mag := math.Sqrt(dir.X*dir.X + dir.Y*dir.Y + dir.Z*dir.Z)
	if mag == 0 {
		return solid.Sphere(indicatorR).Translate(a)
	}
	dir = v3.XYZ(dir.X/mag, dir.Y/mag, dir.Z/mag)
	pin := solid.Capsule(pinLen, pinR).
		TranslateZ(pinLen/2).
		RotateToVector(v3.Z(1), dir).
		Translate(a)
	ball := solid.Sphere(indicatorR).Translate(a)
	return ball.Union(pin)
}

// cameraFlags returns the f3d --camera-elevation/azimuth flags for a
// camera that looks at the origin from the direction (x, y, z) — so the
// indicator pin is always on the camera-facing side. For face anchors
// (only one non-zero component), we nudge slightly off-axis so the
// adjacent labels are also visible.
func cameraFlags(x, y, z int) string {
	dx, dy, dz := float64(x), float64(y), float64(z)
	// Nudge faces off-axis so adjacent faces (and their labels) show too.
	const off = 0.45
	switch {
	case x != 0 && y == 0 && z == 0: // L/R faces
		dy = -off
		dz = off
	case y != 0 && x == 0 && z == 0: // F/B faces
		dx = off
		dz = off
	case z != 0 && x == 0 && y == 0: // T/B faces
		dx = off
		dy = -off
	case x == 0 && y == 0 && z == 0: // center (unused, but safe)
		dx, dy, dz = 0.6, -0.8, 0.4
	}
	mag := math.Sqrt(dx*dx + dy*dy + dz*dz)
	dx, dy, dz = dx/mag, dy/mag, dz/mag

	// f3d: az=0 → camera at (0, -D, 0); az rotates CCW around +Z viewed from above.
	// camera position direction = (sin(az), -cos(az), sin(el))-ish; solve:
	//   az = atan2(dx, -dy), el = asin(dz)
	az := math.Atan2(dx, -dy) * 180 / math.Pi
	el := math.Asin(dz) * 180 / math.Pi
	return fmt.Sprintf("--camera-azimuth-angle %.1f --camera-elevation-angle %.1f", az, el)
}

func writeSTL(out, name string, s *solid.Solid) {
	path := filepath.Join(out, name+".stl")
	cells := solid.CellsFor(s, cellsPerMM)
	r := sdfrender.NewMarchingCubesOctree(cells)
	render.ToSTL(s, path, r)
	fmt.Printf("  %s (%d cells)\n", filepath.Base(path), cells)
}

func writeFlags(out, name, flags string) {
	path := filepath.Join(out, name+".flags")
	if err := os.WriteFile(path, []byte(flags), 0o644); err != nil {
		panic(err)
	}
}

func writeAnchor(out, name string, x, y, z int) {
	scene := host().Union(indicator(x, y, z))
	writeSTL(out, name, scene)
	writeFlags(out, name, cameraFlags(x, y, z))
}

// --- placement verb scenes ---

func mover() *solid.Solid {
	return solid.Sphere(4)
}

func onScene() *solid.Solid {
	return mover().Bottom().On(host().Top()).Union()
}

func directionalScene(builder func(host, mover *solid.Solid) solid.Placement) *solid.Solid {
	return builder(host(), mover()).Union()
}

// insideScene visualizes Inside by performing the actual Union, then
// carving out the front-right-top octant of the host so the sphere
// unioned inside is visible in situ. A pure Union would render as just
// the host (sphere fully contained → no surface), so we have to expose
// the interior — but the geometry IS the union the docs describe.
func insideScene() *solid.Solid {
	bare := solid.Box(v3.XYZ(hostX, hostY, hostZ), 0)
	sphere := solid.Sphere(5).Inside(bare).Solid()

	// Octant-sized box positioned with one corner at the host center,
	// the opposite corner at the host's +X/-Y/+Z corner. Cutting it from
	// the host removes that one corner, opening a cube-shaped window.
	corner := solid.Box(v3.XYZ(hostX/2, hostY/2, hostZ/2), 0).
		Translate(v3.XYZ(hostX/4, -hostY/4, hostZ/4))

	return host().Cut(corner).Union(sphere)
}

// bottomAtScene shows the sphere mover landed so its bottom face sits
// on z=0 — the floor plate is unioned in for a clear ground-plane
// reference, so the viewer sees the sphere is "on the floor", not
// floating in space.
func bottomAtScene() *solid.Solid {
	floor := solid.Box(v3.XYZ(20, 16, 0.6), 0.1).TopAt(0)
	ball := mover().BottomAt(0)
	return floor.Union(ball)
}

func writeVerb(out, name string, s *solid.Solid) {
	writeSTL(out, name, s)
	// Verb camera: stand back from the front-right-top quadrant, but
	// raise the elevation a bit so stacked-on-top scenes (On, Above,
	// etc.) read clearly.
	writeFlags(out, name, "--camera-azimuth-angle 35 --camera-elevation-angle 22")
}

func main() {
	out := "build"
	if len(os.Args) > 1 {
		out = os.Args[1]
	}
	if err := os.MkdirAll(out, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir: %v\n", err)
		os.Exit(1)
	}

	type ac struct {
		name    string
		x, y, z int
	}
	anchors := []ac{
		// 6 faces
		{"anchor-top", 0, 0, 1},
		{"anchor-bottom", 0, 0, -1},
		{"anchor-right", 1, 0, 0},
		{"anchor-left", -1, 0, 0},
		{"anchor-back", 0, 1, 0},
		{"anchor-front", 0, -1, 0},
		// 12 edges
		{"anchor-top-right", 1, 0, 1},
		{"anchor-top-left", -1, 0, 1},
		{"anchor-top-front", 0, -1, 1},
		{"anchor-top-back", 0, 1, 1},
		{"anchor-bottom-right", 1, 0, -1},
		{"anchor-bottom-left", -1, 0, -1},
		{"anchor-bottom-front", 0, -1, -1},
		{"anchor-bottom-back", 0, 1, -1},
		{"anchor-front-right", 1, -1, 0},
		{"anchor-front-left", -1, -1, 0},
		{"anchor-back-right", 1, 1, 0},
		{"anchor-back-left", -1, 1, 0},
		// 8 corners
		{"anchor-top-front-right", 1, -1, 1},
		{"anchor-top-front-left", -1, -1, 1},
		{"anchor-top-back-right", 1, 1, 1},
		{"anchor-top-back-left", -1, 1, 1},
		{"anchor-bottom-front-right", 1, -1, -1},
		{"anchor-bottom-front-left", -1, -1, -1},
		{"anchor-bottom-back-right", 1, 1, -1},
		{"anchor-bottom-back-left", -1, 1, -1},
	}
	for _, a := range anchors {
		writeAnchor(out, a.name, a.x, a.y, a.z)
	}

	writeVerb(out, "verb-on", onScene())
	writeVerb(out, "verb-above", directionalScene(func(h, m *solid.Solid) solid.Placement {
		return m.Bottom().Above(h.Top(), 2)
	}))
	writeVerb(out, "verb-below", directionalScene(func(h, m *solid.Solid) solid.Placement {
		return m.Top().Below(h.Bottom(), 2)
	}))
	writeVerb(out, "verb-right-of", directionalScene(func(h, m *solid.Solid) solid.Placement {
		return m.Left().RightOf(h.Right(), 2)
	}))
	writeVerb(out, "verb-left-of", directionalScene(func(h, m *solid.Solid) solid.Placement {
		return m.Right().LeftOf(h.Left(), 2)
	}))
	writeVerb(out, "verb-behind", directionalScene(func(h, m *solid.Solid) solid.Placement {
		return m.Front().Behind(h.Back(), 2)
	}))
	writeVerb(out, "verb-in-front-of", directionalScene(func(h, m *solid.Solid) solid.Placement {
		return m.Back().InFrontOf(h.Front(), 2)
	}))
	writeVerb(out, "verb-inside", insideScene())
	writeVerb(out, "verb-bottom-at", bottomAtScene())
}
