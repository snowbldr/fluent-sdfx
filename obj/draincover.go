package obj

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/sdfx/obj"
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
