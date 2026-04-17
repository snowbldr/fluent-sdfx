package obj

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
)

// PipeParameters describes a standard pipe.
type PipeParameters = obj.PipeParameters

// PipeConnectorParms configures a pipe connector.
type PipeConnectorParms = obj.PipeConnectorParms

// PipeLookup returns the parameters for a named standard pipe.
func PipeLookup(name, units string) *PipeParameters {
	p, err := obj.PipeLookup(name, units)
	if err != nil {
		panic(err)
	}
	return p
}

// Pipe3D returns a 3D pipe section.
func Pipe3D(oRadius, iRadius, length float64) *solid.Solid {
	s, err := obj.Pipe3D(oRadius, iRadius, length)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}

// StdPipe3D returns a 3D pipe section looked up by standard name.
func StdPipe3D(name, units string, length float64) *solid.Solid {
	s, err := obj.StdPipe3D(name, units, length)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}

// PipeConnector3D returns a 3D pipe connector with up to 6 port recesses.
func PipeConnector3D(p PipeConnectorParms) *solid.Solid {
	s, err := obj.PipeConnector3D(&p)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}

// StdPipeConnector3D returns a 3D pipe connector for a standard pipe.
func StdPipeConnector3D(name, units string, length float64, cfg [6]bool) *solid.Solid {
	s, err := obj.StdPipeConnector3D(name, units, length, cfg)
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}
