// Parametric helpers: a PCB standoff via obj.Standoff3D.
//
// PillarHeight × PillarDiameter is the cylinder body. HoleDepth > 0 makes
// a tapped/screw-receiving hole; HoleDepth < 0 produces a support stub.
// NumberWebs adds triangular gussets around the base for stiffness.
package main

import "github.com/snowbldr/fluent-sdfx/obj"

func main() {
	obj.Standoff3D(obj.StandoffParms{
		PillarHeight:   12,
		PillarDiameter: 6,
		HoleDepth:      8,
		HoleDiameter:   2.5,
		NumberWebs:     4,
		WebHeight:      4,
		WebDiameter:    10,
		WebWidth:       1,
	}).STL("out.stl", 8.0)
}
