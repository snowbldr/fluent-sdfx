package obj

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
)

// GenevaParms configures a Geneva drive mechanism.
type GenevaParms = obj.GenevaParms

// Geneva2D returns the driver and driven 2D profiles of a Geneva drive.
func Geneva2D(p GenevaParms) (driver, driven *shape.Shape) {
	a, b, err := obj.Geneva2D(&p)
	if err != nil {
		panic(err)
	}
	return shape.Wrap2D(a), shape.Wrap2D(b)
}
