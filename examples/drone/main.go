package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const shrink = 1.0 / 0.999 // PLA ~0.1%

var kArm = obj.DroneArmParms{
	MotorSize:     v2.XY(28, 30),
	MotorMount:    v3.XYZ(16, 19, 3.4),
	RotorCavity:   v2.XY(9, 1.5),
	WallThickness: 3.0,
	SideClearance: 1.5,
	MountHeight:   0.7,
	ArmHeight:     0.9,
	ArmLength:     70.0,
}

var kSocket = obj.DroneArmSocketParms{
	Arm:       &kArm,
	Size:      v3.XYZ(40, 30, 30),
	Clearance: 0.5,
	Stop:      35,
}

func main() {
	obj.DroneMotorArm(kArm).ScaleUniform(shrink).ToSTL("arm.stl", 300)
	obj.DroneMotorArmSocket(kSocket).ScaleUniform(shrink).ToSTL("socket.stl", 300)
}
