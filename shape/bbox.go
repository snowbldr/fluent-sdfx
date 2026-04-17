package shape

import (
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

// Box2 is a 2D axis-aligned bounding box.
type Box2 = v2.Box

// NewBox2 returns a Box2 with the given center and size.
func NewBox2(center, size v2.Vec) Box2 { return v2.NewBox(center, size) }
