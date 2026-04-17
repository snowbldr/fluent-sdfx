package shape

import "github.com/deadsy/sdfx/sdf"

// Flange1 returns a 2D flange profile built from a center circle and two side circles
// connected by tangent lines.
//
// distance is from the center circle to each side circle's center; centerRadius and
// sideRadius are the respective circle radii.
func Flange1(distance, centerRadius, sideRadius float64) *Shape {
	return &Shape{sdf.NewFlange1(distance, centerRadius, sideRadius)}
}
