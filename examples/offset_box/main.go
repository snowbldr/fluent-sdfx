package main

import (
	"log"

	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

const sizeX = 30.0
const sizeY = 40.0
const sizeZ = 30.0

const wallThickness = 3.0
const outerRadius = 6.0
const lidPosition = 0.75

func box() error {
	if outerRadius < wallThickness {
		return units.ErrMsg("outerRadius < wallThickness")
	}

	innerOfs := outerRadius - wallThickness
	outerOfs := innerOfs + wallThickness

	if sizeX < outerOfs {
		return units.ErrMsg("sizeX < outerOfs")
	}
	if sizeY < outerOfs {
		return units.ErrMsg("sizeY < outerOfs")
	}
	if sizeZ < outerOfs {
		return units.ErrMsg("sizeZ < outerOfs")
	}

	baseBox := solid.Box(v3.XYZ(sizeX-outerOfs, sizeY-outerOfs, sizeZ-outerOfs), 0)

	innerBox := baseBox.Offset(innerOfs)
	outerBox := baseBox.Offset(outerOfs)
	theBox := outerBox.Cut(innerBox)

	lidZ := (lidPosition - 0.5) * sizeZ
	base := theBox.CutPlane(v3.XYZ(0, 0, lidZ), v3.XYZ(0, 0, -1))
	top := theBox.CutPlane(v3.XYZ(0, 0, lidZ), v3.XYZ(0, 0, 1))

	base.STL("base.stl", 3.0)
	top.STL("top.stl", 3.0)

	return nil
}

func main() {
	err := box()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}
