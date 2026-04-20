package main

import (
	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/solid"
	"github.com/snowbldr/fluent-sdfx/units"
)

func drainCover(k obj.DrainCoverParms) *solid.Solid {
	return obj.DrainCover(k)
}

func vent2() *solid.Solid {
	return drainCover(obj.DrainCoverParms{
		WallDiameter:   1.9 * units.MillimetresPerInch,
		WallHeight:     0.5 * units.MillimetresPerInch,
		WallThickness:  0.125 * units.MillimetresPerInch,
		WallDraft:      0,
		OuterWidth:     0.2 * units.MillimetresPerInch,
		InnerWidth:     0.18 * units.MillimetresPerInch,
		CoverThickness: 0.125 * units.MillimetresPerInch,
		GrateNumber:    8,
		GrateWidth:     1.1,
		GrateDraft:     0,
		CrossBarWidth:  0,
		CrossBarWeb:    false,
	})
}

func drain4() *solid.Solid {
	return drainCover(obj.DrainCoverParms{
		WallDiameter:   3.9 * units.MillimetresPerInch,
		WallHeight:     0.8 * units.MillimetresPerInch,
		WallThickness:  0.2 * units.MillimetresPerInch,
		WallDraft:      units.DtoR(2.0),
		OuterWidth:     0.4 * units.MillimetresPerInch,
		InnerWidth:     0.3 * units.MillimetresPerInch,
		CoverThickness: 0.2 * units.MillimetresPerInch,
		GrateNumber:    8,
		GrateWidth:     1.1,
		GrateDraft:     units.DtoR(8.0),
		CrossBarWidth:  0.8,
		CrossBarWeb:    false,
	})
}

func drain6() *solid.Solid {
	return drainCover(obj.DrainCoverParms{
		WallDiameter:   5.8 * units.MillimetresPerInch,
		WallHeight:     0.8 * units.MillimetresPerInch,
		WallThickness:  0.2 * units.MillimetresPerInch,
		WallDraft:      units.DtoR(2.0),
		OuterWidth:     0.4 * units.MillimetresPerInch,
		InnerWidth:     0.3 * units.MillimetresPerInch,
		CoverThickness: 0.3 * units.MillimetresPerInch,
		GrateNumber:    9,
		GrateWidth:     1.0,
		GrateDraft:     units.DtoR(8.0),
		CrossBarWidth:  1.8,
		CrossBarWeb:    true,
	})
}

func drain12() *solid.Solid {
	return drainCover(obj.DrainCoverParms{
		WallDiameter:   11.8 * units.MillimetresPerInch,
		WallHeight:     1.0 * units.MillimetresPerInch,
		WallThickness:  0.3 * units.MillimetresPerInch,
		WallDraft:      units.DtoR(2.0),
		OuterWidth:     0.8 * units.MillimetresPerInch,
		InnerWidth:     0.5 * units.MillimetresPerInch,
		CoverThickness: 0.3 * units.MillimetresPerInch,
		GrateNumber:    10,
		GrateWidth:     1.0,
		GrateDraft:     units.DtoR(8.0),
		CrossBarWidth:  1.5,
		CrossBarWeb:    true,
	})
}

func main() {
	vent2().STL("vent2.stl", 3.0)
	drain4().STL("drain4.stl", 3.0)
	drain6().STL("drain6.stl", 3.0)
	drain12().STL("drain12.stl", 4.0)
}
