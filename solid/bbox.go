package solid

import (
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Box3 is a 3D axis-aligned bounding box.
type Box3 = v3.Box

// NewBox3 returns a Box3 with the given center and size.
func NewBox3(center, size v3.Vec) Box3 { return v3.NewBox(center, size) }
