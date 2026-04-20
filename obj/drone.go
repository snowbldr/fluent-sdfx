package obj

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/sdfx/obj"
	v2sdf "github.com/snowbldr/sdfx/vec/v2"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
)

// DroneArmParms configures a drone motor arm.
type DroneArmParms struct {
	MotorSize     v2.Vec  // motor diameter/height
	MotorMount    v3.Vec  // motor mount l0, l1, diameter
	RotorCavity   v2.Vec  // cavity for bottom of rotor
	WallThickness float64 // wall thickness
	SideClearance float64 // wall to motor clearance
	MountHeight   float64 // height of motor mount wrt motor height
	ArmHeight     float64 // height of arm wrt motor mount height
	ArmLength     float64 // length of rotor arm
}

func (p *DroneArmParms) toSDF() *obj.DroneArmParms {
	return &obj.DroneArmParms{
		MotorSize:     v2sdf.Vec(p.MotorSize),
		MotorMount:    v3sdf.Vec(p.MotorMount),
		RotorCavity:   v2sdf.Vec(p.RotorCavity),
		WallThickness: p.WallThickness,
		SideClearance: p.SideClearance,
		MountHeight:   p.MountHeight,
		ArmHeight:     p.ArmHeight,
		ArmLength:     p.ArmLength,
	}
}

// DroneArmSocketParms configures a drone motor arm socket.
type DroneArmSocketParms struct {
	Arm       *DroneArmParms // drone arm parameters
	Size      v3.Vec         // body size for socket
	Clearance float64        // clearance between arm and socket
	Stop      float64        // depth of arm stop
}

func (p *DroneArmSocketParms) toSDF() *obj.DroneArmSocketParms {
	return &obj.DroneArmSocketParms{
		Arm:       p.Arm.toSDF(),
		Size:      v3sdf.Vec(p.Size),
		Clearance: p.Clearance,
		Stop:      p.Stop,
	}
}

// DroneMotorArm returns a 3D drone motor arm.
func DroneMotorArm(p DroneArmParms) *solid.Solid {
	s, err := obj.DroneMotorArm(p.toSDF())
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}

// DroneMotorArmSocket returns a 3D drone motor arm socket.
func DroneMotorArmSocket(p DroneArmSocketParms) *solid.Solid {
	s, err := obj.DroneMotorArmSocket(p.toSDF())
	if err != nil {
		panic(err)
	}
	return solid.Wrap(s)
}
