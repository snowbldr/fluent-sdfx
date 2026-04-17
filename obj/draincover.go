package obj

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
)

// DrainCoverParms configures a drain cover.
type DrainCoverParms = obj.DrainCoverParms

// DrainCover returns a 3D drain cover.
func DrainCover(p DrainCoverParms) *solid.Solid {
	s, err := obj.DrainCover(&p)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}
