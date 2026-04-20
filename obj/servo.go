package obj

import (
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/sdfx/obj"
	v2sdf "github.com/snowbldr/sdfx/vec/v2"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
)

// ServoParms describes a servo body and mounting.
type ServoParms struct {
	Body        v3.Vec
	Mount       v3.Vec
	Hole        v2.Vec
	MountOffset float64
	ShaftOffset float64
	ShaftLength float64
	ShaftRadius float64
	HoleRadius  float64
}

func (p *ServoParms) toSDF() *obj.ServoParms {
	return &obj.ServoParms{
		Body:        v3sdf.Vec(p.Body),
		Mount:       v3sdf.Vec(p.Mount),
		Hole:        v2sdf.Vec(p.Hole),
		MountOffset: p.MountOffset,
		ShaftOffset: p.ShaftOffset,
		ShaftLength: p.ShaftLength,
		ShaftRadius: p.ShaftRadius,
		HoleRadius:  p.HoleRadius,
	}
}

// ServoHornParms configures a servo horn.
type ServoHornParms = obj.ServoHornParms

// ServoLookup returns parameters for a named standard servo.
func ServoLookup(name string) *ServoParms {
	p, err := obj.ServoLookup(name)
	if err != nil {
		panic(err)
	}
	return &ServoParms{
		Body:        v3.Vec(p.Body),
		Mount:       v3.Vec(p.Mount),
		Hole:        v2.Vec(p.Hole),
		MountOffset: p.MountOffset,
		ShaftOffset: p.ShaftOffset,
		ShaftLength: p.ShaftLength,
		ShaftRadius: p.ShaftRadius,
		HoleRadius:  p.HoleRadius,
	}
}

// Servo3D returns a 3D servo body.
func Servo3D(p ServoParms) *solid.Solid {
	s, err := obj.Servo3D(p.toSDF())
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}

// Servo2D returns a 2D mounting profile for a servo.
func Servo2D(p ServoParms, holeRadius float64) *shape.Shape {
	s, err := obj.Servo2D(p.toSDF(), holeRadius)
	if err != nil {
		panic(err)
	}
	return shape.Wrap2D(s)
}

// ServoHorn returns a 2D servo horn profile.
func ServoHorn(p ServoHornParms) *shape.Shape {
	s, err := obj.ServoHorn(&p)
	if err != nil {
		panic(err)
	}
	return shape.Wrap2D(s)
}
