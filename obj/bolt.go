package obj

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/sdfx/obj"
)

// BoltParms configures a bolt.
type BoltParms = obj.BoltParms

// Bolt returns a 3D bolt with a hex or socket cap head.
func Bolt(p BoltParms) *solid.Solid {
	s, err := obj.Bolt(&p)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}
