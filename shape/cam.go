package shape

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
)

// FlatFlankCam returns a 2D cam profile made from a base circle, a nose circle, and tangential flats.
// distance is the base-to-nose center distance; baseRadius > noseRadius.
// The base circle is centered at the origin, nose on the +Y axis.
func FlatFlankCam(distance, baseRadius, noseRadius float64) *Shape {
	s, err := sdf.FlatFlankCam2D(distance, baseRadius, noseRadius)
	if err != nil {
		panic(err)
	}
	return &Shape{s}
}

// MakeFlatFlankCam builds a flat-flank cam profile from follower-design parameters.
// lift is the follower lift, durationDeg is the angular duration of the lift in degrees,
// and maxDiameter is the maximum rotation diameter of the cam.
func MakeFlatFlankCam(lift, durationDeg, maxDiameter float64) *Shape {
	s, err := sdf.MakeFlatFlankCam(lift, durationDeg*math.Pi/180, maxDiameter)
	if err != nil {
		panic(err)
	}
	return &Shape{s}
}

// ThreeArcCam returns a 2D cam profile made from a base circle, nose circle, and circular flank arcs.
func ThreeArcCam(distance, baseRadius, noseRadius, flankRadius float64) *Shape {
	s, err := sdf.ThreeArcCam2D(distance, baseRadius, noseRadius, flankRadius)
	if err != nil {
		panic(err)
	}
	return &Shape{s}
}

// MakeThreeArcCam builds a three-arc cam from follower-design parameters.
func MakeThreeArcCam(lift, durationDeg, maxDiameter, k float64) *Shape {
	s, err := sdf.MakeThreeArcCam(lift, durationDeg*math.Pi/180, maxDiameter, k)
	if err != nil {
		panic(err)
	}
	return &Shape{s}
}
