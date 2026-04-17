package obj

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
)

// StandoffParms configures a PCB standoff.
type StandoffParms = obj.StandoffParms

// Standoff3D returns a 3D PCB standoff.
func Standoff3D(p StandoffParms) *solid.Solid {
	s, err := obj.Standoff3D(&p)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}
