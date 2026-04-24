package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// material shrinkage
const shrink = 1.0 / 0.999 // PLA ~0.1%

const wallThickness = 3.0
const round = 0.5 * wallThickness
const clearance = 0.3

func box1(upper bool) *solid.Solid {
	oSize := v3.XYZ(40, 40, 20)
	iSize := oSize.SubScalar(2.0 * wallThickness)

	// build the box
	outer := shape.Rect(v2.XY(oSize.X, oSize.Y), round).Extrude(oSize.Z)
	inner := shape.Rect(v2.XY(iSize.X, iSize.Y), round).Extrude(iSize.Z)
	box := outer.Cut(inner)

	// add some internal walls
	yOfs := oSize.Y * 0.2
	wall := solid.Box(v3.XYZ(oSize.X, wallThickness, oSize.Z), 0)
	wall0 := wall.Translate(v3.Y(yOfs))
	wall1 := wall.Translate(v3.Y(-yOfs))
	box = box.Union(wall0, wall1)

	lidHeight := 0.5*oSize.Z - wallThickness

	if upper {
		box = box.CutPlane(v3.Z(lidHeight), v3.Z(1))
	} else {
		box = box.CutPlane(v3.Z(lidHeight), v3.Z(-1))
	}

	// angled tabs
	tabSize := v3.XYZ(2.5*wallThickness, wallThickness, wallThickness)
	tab := obj.NewAngleTab(tabSize, clearance)
	xOfs := oSize.X * 0.25
	mSet := []solid.M44{
		solid.Translate3d(v3.XYZ(xOfs, yOfs, lidHeight)),
		solid.Translate3d(v3.XYZ(xOfs, -yOfs, lidHeight)),
		solid.Translate3d(v3.XYZ(-xOfs, yOfs, lidHeight)),
		solid.Translate3d(v3.XYZ(-xOfs, -yOfs, lidHeight)),
	}
	withTabs := obj.AddTabs(box, tab, upper, mSet)

	// screw tabs
	l := oSize.Z * 0.35
	k := obj.ScrewTab{
		Length:     l,                   // length of pillar
		Radius:     0.8 * wallThickness, // radius of pillar
		Round:      true,                // round the bottom of the pillar
		HoleUpper:  wallThickness,       // length of upper hole
		HoleLower:  0.8 * l,             // length of lower hole
		HoleRadius: 1,                   // radius of hole
	}
	tab2 := obj.NewScrewTab(k)
	xOfs = 0.5*oSize.X - wallThickness
	yOfs = 0.5*oSize.Y - wallThickness
	mSet = []solid.M44{
		solid.Translate3d(v3.XYZ(xOfs, yOfs, lidHeight)),
		solid.Translate3d(v3.XYZ(-xOfs, yOfs, lidHeight)),
		solid.Translate3d(v3.XYZ(xOfs, -yOfs, lidHeight)),
		solid.Translate3d(v3.XYZ(-xOfs, -yOfs, lidHeight)),
	}
	return obj.AddTabs(withTabs, tab2, upper, mSet)
}

func box0(upper bool) *solid.Solid {
	oSize := v3.XYZ(40, 40, 20)
	iSize := oSize.SubScalar(2.0 * wallThickness)

	outer := solid.Box(oSize, round)
	inner := solid.Box(iSize, round)

	box := outer.Cut(inner)
	lidHeight := oSize.Z * 0.25

	if upper {
		box = box.CutPlane(v3.Z(lidHeight), v3.Z(1))
	} else {
		box = box.CutPlane(v3.Z(lidHeight), v3.Z(-1))
	}

	tabSize := v3.XYZ(3.0*wallThickness, 0.5*wallThickness, wallThickness)
	tab := obj.NewStraightTab(tabSize, clearance)

	xOfs := 0.5 * (iSize.X + wallThickness)
	yOfs := 0.5 * (iSize.Y + wallThickness)

	mSet := []solid.M44{
		solid.Translate3d(v3.XZ(xOfs, lidHeight)).Mul(solid.RotateZMatrix(90)),
		solid.Translate3d(v3.XZ(-xOfs, lidHeight)).Mul(solid.RotateZMatrix(90)),
		solid.Translate3d(v3.YZ(yOfs, lidHeight)),
		solid.Translate3d(v3.YZ(-yOfs, lidHeight)),
	}

	return obj.AddTabs(box, tab, upper, mSet)
}

func main() {
	box0(true).ScaleUniform(shrink).STL("box0_upper.stl", 3.0)
	box0(false).ScaleUniform(shrink).STL("box0_lower.stl", 3.0)
	box1(true).ScaleUniform(shrink).STL("box1_upper.stl", 3.0)
	box1(false).ScaleUniform(shrink).STL("box1_lower.stl", 3.0)
}
