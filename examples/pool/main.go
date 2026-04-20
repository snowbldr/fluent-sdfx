package main

import (
	"log"

	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
)

const cubicInchesPerGallon = 231.0

// pool dimensions are in inches
const poolWidth = 234.0
const poolLength = 477.0

var poolDepth = []v2.Vec{v2.XY(0.0, 43.0), v2.XY(101.0, 46.0), v2.XY(202.0, 58.0), v2.XY(298.0, 83.0), v2.XY(394.0, 96.0), v2.XY(477.0, 96.0)}

const vol = (7738.3005 * 1000.0) / cubicInchesPerGallon // gallons

func pool() *solid.Solid {
	log.Printf("pool volume %f gallons\n", vol)
	p := shape.NewPoly()
	p.Add(0, 0)
	p.AddV2Set(poolDepth)
	p.Add(poolLength, 0)
	return solid.Extrude(p.Build(), poolWidth)
}

func main() {
	pool().STL("pool1.stl", 3.0)
}
