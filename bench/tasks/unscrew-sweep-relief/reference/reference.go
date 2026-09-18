// Package reference: the ribs relieved along the path the lug actually
// takes on the way out.
//
// The insert unscrews: it turns 360 degrees counter-clockwise and rises
// 6 mm, and the two are locked together, so the lug is 1 mm higher by the
// time it reaches the rib at 120 degrees, 3 mm higher at 240, and 5 mm
// higher when it comes back round to the rib at 0. The relief is therefore
// a HELICAL band, not a horizontal one: each rib is notched at a different
// height, and each notch is tilted by the rise across the rib's own width.
//
// The lug starts at 60 degrees, in the gap between two ribs, so at the
// start pose there is no interference at all. Nothing short of sweeping the
// coupled rotation and rise finds the collision.
//
// The swept volume is built as a union of the lug's poses, one per degree
// of rotation, grown by the 0.3 clearance (a rounded box whose faces are
// the lug's faces offset 0.3). One degree of rotation is 6/360 = 0.017 mm
// of rise and 0.0003 mm of chord sag at the lug's outer corner, both far
// below the scoring grid cell, so the scallops between poses are invisible.
// Only the poses within 45 degrees of a rib can touch it (the lug's and the
// rib's angular half-widths add up to at most 33 degrees), so each rib is
// cut by its own 91-pose window and the union's bounding boxes prune the
// rest.
//
// Only the ribs are cut: the bore wall and the outside of the barrel are
// left exactly as they were.
//
// Incident: mcweed I-7 (U533-U541) and mined2 I-6 -- interference that a
// static fit check passes and only a sweep of the real motion finds. "Swept
// the fingers two full turns through it"; "pin surface swept through two
// unscrew turns". The fix pattern is to encode the motion, not the pose.
package reference

import (
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const (
	bodyOuterR = 14.0
	boreR      = 8.0
	bodyH      = 20.0

	ribDepth = 2.0
	ribW     = 3.0
	ribCount = 3

	// The insert (not part of this solid): a post of radius 6 carrying one
	// lug, 4 wide tangentially, 1.5 out radially (r = 6..7.5), 5 tall,
	// centred on z = 0 at the start pose, starting at 60 degrees.
	insertR     = 6.0
	lugW        = 4.0
	lugOut      = 1.5
	lugH        = 5.0
	lugStartDeg = 60.0

	// Unscrewing: one full turn counter-clockwise, 6 mm of rise, coupled.
	unscrewDeg  = 360.0
	unscrewRise = 6.0

	clearance = 0.3 // lug to rib, all round

	poseStepDeg     = 1.0  // one pose per degree of rotation
	reliefWindowDeg = 45.0 // poses further than this cannot touch a rib
)

var ribAngles = [ribCount]float64{0, 120, 240}

// lugMidR is the radius of the lug's mid-thickness.
const lugMidR = insertR + lugOut/2

// sweptNear returns the volume swept by the lug (grown by the clearance)
// over the poses that can reach the rib at ribDeg.
func sweptNear(ribDeg float64) *solid.Solid {
	// The lug's faces offset outward by the clearance: a box enlarged by
	// 2*clearance with a clearance-radius round on its edges.
	tool := solid.Box(v3.XYZ(lugOut+2*clearance, lugW+2*clearance, lugH+2*clearance), clearance).
		TranslateX(lugMidR)

	// Turn fraction at which the lug's centre line reaches this rib.
	cross := mod360(ribDeg-lugStartDeg) / unscrewDeg

	var poses []*solid.Solid
	for d := -reliefWindowDeg; d <= reliefWindowDeg+1e-9; d += poseStepDeg {
		t := cross + d/unscrewDeg
		poses = append(poses, tool.
			TranslateZ(unscrewRise*t).
			RotateZ(lugStartDeg+unscrewDeg*t))
	}
	return solid.UnionAll(poses...)
}

func mod360(deg float64) float64 {
	for deg < 0 {
		deg += 360
	}
	for deg >= 360 {
		deg -= 360
	}
	return deg
}

func Build() *solid.Solid {
	body := solid.Cylinder(bodyH, bodyOuterR, 0).
		Cut(solid.Cylinder(bodyH+2, boreR, 0))

	ribs := make([]*solid.Solid, 0, ribCount)
	for _, a := range ribAngles {
		rib := solid.Box(v3.XYZ(ribDepth, ribW, bodyH), 0).
			TranslateX(boreR - ribDepth/2).
			RotateZ(a)
		ribs = append(ribs, rib.Cut(sweptNear(a)))
	}

	return body.Union(ribs...)
}
