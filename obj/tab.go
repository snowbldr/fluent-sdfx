package obj

import (
	"github.com/deadsy/sdfx/obj"
	"github.com/deadsy/sdfx/sdf"
	v3sdf "github.com/deadsy/sdfx/vec/v3"
	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// Tab is an interface for tab geometry used when splitting a solid across a plane.
type Tab = obj.Tab

// ScrewTab configures a screw-type tab.
type ScrewTab = obj.ScrewTab

// NewStraightTab creates a straight (rectangular) tab with the given size and clearance.
func NewStraightTab(size v3.Vec, clearance float64) Tab {
	t, err := obj.NewStraightTab(v3sdf.Vec(size), clearance)
	if err != nil {
		panic(err)
	}
	return t
}

// NewAngleTab creates an angled (triangular) tab with the given size and clearance.
func NewAngleTab(size v3.Vec, clearance float64) Tab {
	t, err := obj.NewAngleTab(v3sdf.Vec(size), clearance)
	if err != nil {
		panic(err)
	}
	return t
}

// NewScrewTab creates a screw-style tab.
func NewScrewTab(p ScrewTab) Tab {
	t, err := obj.NewScrewTab(&p)
	if err != nil {
		panic(err)
	}
	return t
}

// AddTabs adds tab geometry to s at the given transforms. If upper is true
// the tab is added to the upper half; otherwise the lower half.
func AddTabs(s *solid.Solid, tab Tab, upper bool, mset []solid.M44) *solid.Solid {
	sdfM := make([]sdf.M44, len(mset))
	for i, m := range mset {
		sdfM[i] = sdf.M44(m)
	}
	return solid.Wrap(obj.AddTabs(s.SDF3, tab, upper, sdfM))
}
