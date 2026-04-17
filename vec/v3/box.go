package v3

import (
	"github.com/deadsy/sdfx/sdf"
	v3sdf "github.com/deadsy/sdfx/vec/v3"
)

// Box is a 3D axis-aligned bounding box.
type Box struct {
	Min, Max Vec
}

func fromSDFBox(b sdf.Box3) Box {
	return Box{Min: Vec(b.Min), Max: Vec(b.Max)}
}

// SDF returns the sdfx-compatible representation of the box.
// Internal use only.
func (a Box) SDF() sdf.Box3 {
	return sdf.Box3{Min: v3sdf.Vec(a.Min), Max: v3sdf.Vec(a.Max)}
}

// FromSDF promotes an sdfx Box3 to our Box.
// Internal use only.
func FromSDF(b sdf.Box3) Box { return fromSDFBox(b) }

// NewBox returns a Box with the given center and size.
func NewBox(center, size Vec) Box {
	return fromSDFBox(sdf.NewBox3(v3sdf.Vec(center), v3sdf.Vec(size)))
}

// Center returns the center of the box.
func (a Box) Center() Vec { return Vec(a.SDF().Center()) }

// Contains reports whether v is inside the box.
func (a Box) Contains(v Vec) bool { return a.SDF().Contains(v3sdf.Vec(v)) }

// Cube returns the smallest cubic box containing a.
func (a Box) Cube() Box { return fromSDFBox(a.SDF().Cube()) }

// Enlarge grows the box by v on each side.
func (a Box) Enlarge(v Vec) Box { return fromSDFBox(a.SDF().Enlarge(v3sdf.Vec(v))) }

// Equals reports whether a and b are within delta of each other.
func (a Box) Equals(b Box, delta float64) bool { return a.SDF().Equals(b.SDF(), delta) }

// Extend returns the bounding box containing both a and b.
func (a Box) Extend(b Box) Box { return fromSDFBox(a.SDF().Extend(b.SDF())) }

// Include returns a box extended to include v.
func (a Box) Include(v Vec) Box { return fromSDFBox(a.SDF().Include(v3sdf.Vec(v))) }

// ScaleAboutCenter scales the box about its center by k.
func (a Box) ScaleAboutCenter(k float64) Box { return fromSDFBox(a.SDF().ScaleAboutCenter(k)) }

// Size returns the size of the box.
func (a Box) Size() Vec { return Vec(a.SDF().Size()) }

// Translate returns a translated box.
func (a Box) Translate(v Vec) Box { return fromSDFBox(a.SDF().Translate(v3sdf.Vec(v))) }

// Vertices returns the eight corners of the box.
func (a Box) Vertices() VecSet {
	vs := a.SDF().Vertices()
	out := make(VecSet, len(vs))
	for i, v := range vs {
		out[i] = Vec(v)
	}
	return out
}

// Random returns a random point inside the box.
func (a *Box) Random() Vec {
	sdfBox := a.SDF()
	return Vec(sdfBox.Random())
}

// RandomSet returns n random points inside the box.
func (a *Box) RandomSet(n int) VecSet {
	sdfBox := a.SDF()
	vs := sdfBox.RandomSet(n)
	out := make(VecSet, len(vs))
	for i, v := range vs {
		out[i] = Vec(v)
	}
	return out
}
