package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

func servos() {

	names := []string{
		"nano",
		"submicro",
		"micro",
		"mini",
		"standard",
		"large",
		"giant",
	}

	var s *solid.Solid
	yOfs := 0.0

	for _, n := range names {
		k := obj.ServoLookup(n)

		yOfs += 0.5*k.Body.Y + 10.0

		servo := obj.Servo3D(*k).Translate(v3.XYZ(0, yOfs, 20))

		outline := obj.Servo2D(*k, -1).Extrude(5).Translate(v3.XYZ(0, yOfs, 0))

		if s == nil {
			s = servo.Union(outline)
		} else {
			s = s.Union(servo, outline)
		}
		yOfs += 0.5 * k.Body.Y
	}

	s.STL("servos.stl", 3.0)
}

func main() {
	servos()
}
