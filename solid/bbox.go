package solid

import (
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Box3 is a 3D axis-aligned bounding box. It embeds v3.Box so all of
// its methods (Center, Size, Vertices, Random, etc.) are available, and
// adds a Solid() method for converting the box back into a *Solid.
type Box3 struct{ v3.Box }

// NewBox3 returns a Box3 with the given center and size.
func NewBox3(center, size v3.Vec) Box3 { return Box3{v3.NewBox(center, size)} }

// wrapBox3 lifts a v3.Box into a Box3.
func wrapBox3(b v3.Box) Box3 { return Box3{b} }

// Solid returns the bounding box as a *Solid (a Box of the same size,
// translated so its center matches this box's center).
func (b Box3) Solid() *Solid {
	return Box(b.Size(), 0).Translate(b.Center())
}

// --- Chainable overrides so methods that produce a new box keep the Box3 type ---

// Cube returns the smallest cubic box containing b.
func (b Box3) Cube() Box3 { return Box3{b.Box.Cube()} }

// Enlarge grows the box by v on each side.
func (b Box3) Enlarge(v v3.Vec) Box3 { return Box3{b.Box.Enlarge(v)} }

// Extend returns the bounding box containing both b and o.
func (b Box3) Extend(o Box3) Box3 { return Box3{b.Box.Extend(o.Box)} }

// Include returns a box extended to include v.
func (b Box3) Include(v v3.Vec) Box3 { return Box3{b.Box.Include(v)} }

// ScaleAboutCenter scales the box about its center by k.
func (b Box3) ScaleAboutCenter(k float64) Box3 { return Box3{b.Box.ScaleAboutCenter(k)} }

// Translate returns a translated box.
func (b Box3) Translate(v v3.Vec) Box3 { return Box3{b.Box.Translate(v)} }

// Equals reports whether b and o are within delta of each other.
func (b Box3) Equals(o Box3, delta float64) bool { return b.Box.Equals(o.Box, delta) }
