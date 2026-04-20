package main

import (
	"math"

	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// phone body
var phone_w = 78.0  // width
var phone_h = 146.5 // height
var phone_t = 13.0  // thickness
var phone_r = 12.0  // corner radius

// camera hole
var camera_w = 23.5
var camera_h = 33.0
var camera_r = 3.0
var camera_xofs = 0.0
var camera_yofs = 48.0

// speaker hole
var speaker_w = 12.5
var speaker_h = 10.0
var speaker_r = 3.0
var speaker_xofs = 23.0
var speaker_yofs = -46.0

// wall thickness
var wall_t = 3.0

func phone_body() *solid.Solid {
	return solid.Extrude(shape.Rect(v2.XY(phone_w, phone_h), phone_r), phone_t).
		Translate(v3.Z(wall_t / 2.0))
}

func camera_hole() *solid.Solid {
	return solid.Extrude(shape.Rect(v2.XY(camera_w, camera_h), camera_r), wall_t+phone_t).
		Translate(v3.XY(camera_xofs, camera_yofs))
}

func speaker_hole() *solid.Solid {
	return solid.Extrude(shape.Rect(v2.XY(speaker_w, speaker_h), speaker_r), wall_t+phone_t).
		Translate(v3.XY(speaker_xofs, speaker_yofs))
}

// holes for buttons, jacks, etc.
var hole_r = 2.0

func hole_left(length, yofs, zofs float64) *solid.Solid {
	w := phone_t * 2.0
	xofs := -(phone_w + wall_t) / 2.0
	yofs = (phone_h-length)/2.0 - yofs
	zofs = phone_t + ((phone_t + wall_t) / 2.0) - zofs
	return solid.Extrude(shape.Rect(v2.XY(w, length), hole_r), wall_t).
		RotateY(90).
		Translate(v3.XYZ(xofs, yofs, zofs))
}

func hole_right(length, yofs, zofs float64) *solid.Solid {
	w := phone_t * 2.0
	xofs := (phone_w + wall_t) / 2.0
	yofs = (phone_h-length)/2.0 - yofs
	zofs = phone_t + ((phone_t + wall_t) / 2.0) - zofs
	return solid.Extrude(shape.Rect(v2.XY(w, length), hole_r), wall_t).
		RotateY(90).
		Translate(v3.XYZ(xofs, yofs, zofs))
}

func hole_top(length, xofs, zofs float64) *solid.Solid {
	w := phone_t * 2.0
	xofs = -(phone_w-length)/2.0 + xofs
	yofs := (phone_h + wall_t) / 2.0
	zofs = phone_t + ((phone_t + wall_t) / 2.0) - zofs
	return solid.Extrude(shape.Rect(v2.XY(length, w), hole_r), wall_t).
		RotateX(90).
		Translate(v3.XYZ(xofs, yofs, zofs))
}

func hole_bottom(length, xofs, zofs float64) *solid.Solid {
	w := phone_t * 2.0
	xofs = -(phone_w-length)/2.0 + xofs
	yofs := -(phone_h + wall_t) / 2.0
	zofs = phone_t + ((phone_t + wall_t) / 2.0) - zofs
	return solid.Extrude(shape.Rect(v2.XY(length, w), hole_r), wall_t).
		RotateX(90).
		Translate(v3.XYZ(xofs, yofs, zofs))
}

func outer_shell() *solid.Solid {
	w := phone_w + (2.0 * wall_t)
	h := phone_h + (2.0 * wall_t)
	r := phone_r + wall_t
	t := phone_t + wall_t
	return solid.Extrude(shape.Rect(v2.XY(w, h), r), t)
}

func clip() *solid.Solid {
	theta := 35.0
	p := shape.NewPoly()
	p.Add(0, 0)
	p.Add(12.0, 0).Rel()
	p.Add(0, 2.0).Rel()
	p.Add(-10.0, 0).Rel()
	p.Add(0, 4.5).Rel()
	p.Add(-19.5411, 0).Rel()
	p.Add(14.8717, (270.0-theta)*math.Pi/180).Polar().Rel()
	p.Add(0, -7.8612).Rel()
	p.Add(4.3306, (270.0+theta)*math.Pi/180).Polar().Rel()
	p.Add(2.0, theta*math.Pi/180).Polar().Rel()
	p.Add(3.7, (90.0+theta)*math.Pi/180).Polar().Rel()
	p.Add(0, 6.6).Rel()
	p.Add(13.2, (90.0-theta)*math.Pi/180).Polar().Rel()
	p.Add(16.5, 0).Rel()
	p.Close()
	return solid.Extrude(p.Build(), 8.0)
}

func additive() *solid.Solid {
	return outer_shell()
}

func subtractive() *solid.Solid {
	return solid.UnionAll(
		phone_body(),
		camera_hole(),
		speaker_hole(),
		hole_left(31.0, 19.5, 8.0),
		hole_right(20.0, 34.0, 8.0),
		hole_top(13.0, 16.0, 8.0),
		hole_top(13.0, 49.5, 9.0),
		hole_bottom(35.0, 20.5, 9.0),
	)
}

func main() {
	clip().STL("clip.stl", 3.0)
	additive().Cut(subtractive()).STL("holder.stl", 3.0)
}
