package obj_test

import (
	"testing"

	"github.com/snowbldr/fluent-sdfx/obj"
	"github.com/snowbldr/fluent-sdfx/shape"
	"github.com/snowbldr/fluent-sdfx/solid"
	v2 "github.com/snowbldr/fluent-sdfx/vec/v2"
	v2i "github.com/snowbldr/fluent-sdfx/vec/v2i"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	v3i "github.com/snowbldr/fluent-sdfx/vec/v3i"
)

// nonEmptySolid asserts that the solid is non-nil and has a positive-volume bbox.
func nonEmptySolid(t *testing.T, s *solid.Solid) {
	t.Helper()
	if s == nil {
		t.Fatal("nil solid")
	}
	sz := s.Bounds().Size()
	if sz.X <= 0 || sz.Y <= 0 || sz.Z <= 0 {
		t.Fatalf("solid has non-positive bbox size: %+v", sz)
	}
}

// nonEmptyShape asserts that the shape is non-nil and has a positive-area bbox.
func nonEmptyShape(t *testing.T, s *shape.Shape) {
	t.Helper()
	if s == nil {
		t.Fatal("nil shape")
	}
	sz := s.Bounds().Size()
	if sz.X <= 0 || sz.Y <= 0 {
		t.Fatalf("shape has non-positive bbox size: %+v", sz)
	}
}

// expectPanic runs fn and fails the test if fn does not panic.
func expectPanic(t *testing.T, name string, fn func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("%s: expected panic", name)
		}
	}()
	fn()
}

// ---------- bolt / nut / threaded cylinder ----------

func TestBolt(t *testing.T) {
	t.Run("hex", func(t *testing.T) {
		nonEmptySolid(t, obj.Bolt(obj.BoltParms{
			Thread: "M8x1.25", Style: "hex", TotalLength: 30, ShankLength: 5,
		}))
	})
	t.Run("knurl", func(t *testing.T) {
		nonEmptySolid(t, obj.Bolt(obj.BoltParms{
			Thread: "M5x0.8", Style: "knurl", TotalLength: 20, ShankLength: 4,
		}))
	})
	t.Run("unc_no_thread", func(t *testing.T) {
		// shank length == total length means zero-length thread branch
		nonEmptySolid(t, obj.Bolt(obj.BoltParms{
			Thread: "unc_1/4", Style: "hex", TotalLength: 10, ShankLength: 10,
		}))
	})
	expectPanic(t, "bad style", func() {
		obj.Bolt(obj.BoltParms{Thread: "M5x0.8", Style: "bogus", TotalLength: 10})
	})
	expectPanic(t, "bad thread", func() {
		obj.Bolt(obj.BoltParms{Thread: "ZZZ", Style: "hex", TotalLength: 10})
	})
}

func TestNut(t *testing.T) {
	nonEmptySolid(t, obj.Nut(obj.NutParms{Thread: "M5x0.8", Style: "hex"}))
	nonEmptySolid(t, obj.Nut(obj.NutParms{Thread: "M4x0.7", Style: "knurl"}))
	expectPanic(t, "nut bad style", func() {
		obj.Nut(obj.NutParms{Thread: "M4x0.7", Style: "no"})
	})
}

func TestThreadedCylinder(t *testing.T) {
	nonEmptySolid(t, obj.ThreadedCylinder(obj.ThreadedCylinderParms{
		Height: 12, Diameter: 10, Thread: "M5x0.8",
	}))
	expectPanic(t, "bad thread", func() {
		obj.ThreadedCylinder(obj.ThreadedCylinderParms{Height: 1, Diameter: 1, Thread: "x"})
	})
}

// ---------- washer ----------

func TestWasher(t *testing.T) {
	p := obj.WasherParms{Thickness: 2, InnerRadius: 4, OuterRadius: 8}
	nonEmptyShape(t, obj.Washer2D(p))
	nonEmptySolid(t, obj.Washer3D(p))
	expectPanic(t, "washer bad", func() {
		obj.Washer2D(obj.WasherParms{InnerRadius: 8, OuterRadius: 4})
	})
}

// ---------- standoff ----------

func TestStandoff(t *testing.T) {
	nonEmptySolid(t, obj.Standoff3D(obj.StandoffParms{
		PillarHeight:   10,
		PillarDiameter: 6,
		HoleDepth:      8,
		HoleDiameter:   3,
		NumberWebs:     4,
		WebHeight:      4,
		WebDiameter:    8,
		WebWidth:       1,
	}))
	// negative HoleDepth -> support stub branch
	nonEmptySolid(t, obj.Standoff3D(obj.StandoffParms{
		PillarHeight:   10,
		PillarDiameter: 6,
		HoleDepth:      -2,
		HoleDiameter:   3,
		NumberWebs:     0,
		WebHeight:      4,
		WebDiameter:    8,
		WebWidth:       1,
	}))
}

// ---------- angle ----------

func TestAngle(t *testing.T) {
	p := obj.AngleParams{
		X:          obj.AngleLeg{Length: 20, Thickness: 2},
		Y:          obj.AngleLeg{Length: 20, Thickness: 2},
		RootRadius: 1,
		Length:     10,
	}
	nonEmptyShape(t, obj.Angle2D(p))
	nonEmptySolid(t, obj.Angle3D(p))
	expectPanic(t, "angle bad", func() {
		obj.Angle2D(obj.AngleParams{X: obj.AngleLeg{Length: -1, Thickness: 1}})
	})
}

// ---------- chamfer ----------

func TestChamferedCylinder(t *testing.T) {
	cyl := solid.Cylinder(10, 4, 0)
	nonEmptySolid(t, obj.ChamferedCylinder(cyl, 0.5, 0.5))
}

// ---------- display ----------

func TestDisplay(t *testing.T) {
	p := obj.DisplayParms{
		Window:          v2.XY(40, 30),
		Rounding:        2,
		Supports:        v2.XY(50, 40),
		SupportHeight:   4,
		SupportDiameter: 5,
		HoleDiameter:    2,
		Offset:          v2.XY(0, 0),
		Thickness:       2,
	}
	nonEmptySolid(t, obj.Display(p, false))
	nonEmptySolid(t, obj.Display(p, true))
	pCS := p
	pCS.Countersunk = true
	nonEmptySolid(t, obj.Display(pCS, false))
}

// ---------- drain cover ----------

func TestDrainCover(t *testing.T) {
	p := obj.DrainCoverParms{
		WallDiameter:   60,
		WallHeight:     10,
		WallThickness:  2,
		WallDraft:      0.05,
		OuterWidth:     5,
		InnerWidth:     2,
		CoverThickness: 3,
		GrateNumber:    6,
		GrateWidth:     2,
		GrateDraft:     0.05,
		CrossBarWidth:  1,
		CrossBarWeb:    true,
	}
	nonEmptySolid(t, obj.DrainCover(p))
	// No crossbar web variant
	p.CrossBarWeb = false
	nonEmptySolid(t, obj.DrainCover(p))
}

// ---------- drone ----------

func TestDroneMotorArm(t *testing.T) {
	p := obj.DroneArmParms{
		MotorSize:     v2.XY(28, 18),
		MotorMount:    v3.XYZ(20, 20, 3),
		RotorCavity:   v2.XY(24, 2),
		WallThickness: 2,
		SideClearance: 1,
		MountHeight:   1.4,
		ArmHeight:     0.9,
		ArmLength:     30,
	}
	nonEmptySolid(t, obj.DroneMotorArm(p))

	socket := obj.DroneArmSocketParms{
		Arm:       &p,
		Size:      v3.XYZ(40, 30, 30),
		Clearance: 0.3,
		Stop:      4,
	}
	nonEmptySolid(t, obj.DroneMotorArmSocket(socket))
}

// ---------- finger button ----------

func TestFingerButton2D(t *testing.T) {
	nonEmptyShape(t, obj.FingerButton2D(obj.FingerButtonParms{
		Width: 4, Gap: 0.5, Length: 10,
	}))
}

// ---------- gear (involute) ----------

func TestInvoluteGear(t *testing.T) {
	nonEmptyShape(t, obj.InvoluteGear(obj.InvoluteGearParms{
		NumberTeeth:      20,
		Module:           1,
		PressureAngleDeg: 20,
		Backlash:         0.05,
		Clearance:        0.1,
		RingWidth:        1,
		Facets:           6,
	}))
}

// ---------- geneva ----------

func TestGeneva2D(t *testing.T) {
	driver, driven := obj.Geneva2D(obj.GenevaParms{
		NumSectors:     6,
		CenterDistance: 30,
		DriverRadius:   12,
		DrivenRadius:   20,
		PinRadius:      2,
		Clearance:      0.1,
	})
	nonEmptyShape(t, driver)
	nonEmptyShape(t, driven)
	expectPanic(t, "geneva bad", func() {
		obj.Geneva2D(obj.GenevaParms{NumSectors: 1})
	})
}

// ---------- gridfinity ----------

func TestGridfinityBaseExtra(t *testing.T) {
	// Smoke variants distinct from the example: no magnet/hole.
	nonEmptySolid(t, obj.GridfinityBase(obj.GridfinityBaseParms{
		Size: v2i.XY(2, 2), Magnet: false, Hole: false,
	}))
}

func TestGridfinityBody(t *testing.T) {
	nonEmptySolid(t, obj.GridfinityBody(obj.GridfinityBodyParms{
		Size: v3i.XYZ(2, 2, 3), Empty: false, Hole: false,
	}))
	nonEmptySolid(t, obj.GridfinityBody(obj.GridfinityBodyParms{
		Size: v3i.XYZ(1, 1, 2), Empty: true, Hole: true,
	}))
}

// ---------- hex ----------

func TestHex(t *testing.T) {
	nonEmptyShape(t, obj.Hex2D(5, 0.5))
	nonEmptySolid(t, obj.Hex3D(5, 4, 0.3))
	t.Run("HexHead_tb", func(t *testing.T) {
		nonEmptySolid(t, obj.HexHead3D(5, 4, "tb"))
	})
	t.Run("HexHead_t", func(t *testing.T) {
		nonEmptySolid(t, obj.HexHead3D(5, 4, "t"))
	})
	t.Run("HexHead_b", func(t *testing.T) {
		nonEmptySolid(t, obj.HexHead3D(5, 4, "b"))
	})
	t.Run("HexHead_none", func(t *testing.T) {
		nonEmptySolid(t, obj.HexHead3D(5, 4, ""))
	})
}

// ---------- hole / circle grille / counter holes / bolt circles / keyed ----------

func TestCircleGrille(t *testing.T) {
	p := obj.CircleGrilleParms{
		HoleDiameter: 2, GrilleDiameter: 30, RadialSpacing: 0.4, TangentialSpacing: 0.4, Thickness: 2,
	}
	nonEmptyShape(t, obj.CircleGrille2D(p))
	nonEmptySolid(t, obj.CircleGrille3D(p))
}

func TestCountersAndChamfers(t *testing.T) {
	nonEmptySolid(t, obj.CounterBoredHole3D(20, 2, 4, 5))
	nonEmptySolid(t, obj.ChamferedHole3D(20, 2, 1))
	nonEmptySolid(t, obj.CounterSunkHole3D(20, 2))
}

func TestBoltCircle(t *testing.T) {
	nonEmptyShape(t, obj.BoltCircle2D(2, 10, 6))
	nonEmptySolid(t, obj.BoltCircle3D(5, 2, 10, 6))
}

func TestKeyedHole(t *testing.T) {
	t.Run("1key", func(t *testing.T) {
		p := obj.KeyedHoleParms{Diameter: 10, KeySize: 0.7, NumKeys: 1, Thickness: 3}
		nonEmptyShape(t, obj.KeyedHole2D(p))
		nonEmptySolid(t, obj.KeyedHole3D(p))
	})
	t.Run("2keys", func(t *testing.T) {
		p := obj.KeyedHoleParms{Diameter: 10, KeySize: 0.4, NumKeys: 2, Thickness: 3}
		nonEmptyShape(t, obj.KeyedHole2D(p))
		nonEmptySolid(t, obj.KeyedHole3D(p))
	})
	expectPanic(t, "keyed hole bad", func() {
		obj.KeyedHole2D(obj.KeyedHoleParms{Diameter: 10, NumKeys: 5})
	})
}

// ---------- keyway ----------

func TestKeyway(t *testing.T) {
	t.Run("internal_key", func(t *testing.T) {
		p := obj.KeywayParameters{ShaftRadius: 5, KeyRadius: 4, KeyWidth: 1.5, ShaftLength: 10}
		nonEmptyShape(t, obj.Keyway2D(p))
		nonEmptySolid(t, obj.Keyway3D(p))
	})
	t.Run("external_key", func(t *testing.T) {
		p := obj.KeywayParameters{ShaftRadius: 5, KeyRadius: 6, KeyWidth: 1.5, ShaftLength: 10}
		nonEmptyShape(t, obj.Keyway2D(p))
	})
	expectPanic(t, "keyway bad", func() {
		obj.Keyway2D(obj.KeywayParameters{ShaftRadius: 0})
	})
}

// ---------- knurl ----------

func TestKnurl(t *testing.T) {
	nonEmptySolid(t, obj.Knurl3D(obj.KnurlParms{
		Length: 10, Radius: 6, Pitch: 1.5, Height: 0.3, ThetaDeg: 45,
	}))
	nonEmptySolid(t, obj.KnurledHead3D(6, 4, 1.2))
}

// ---------- panel ----------

func TestPanel(t *testing.T) {
	p := obj.PanelParms{
		Size:         v2.XY(70, 90),
		CornerRadius: 5,
		HoleDiameter: 3.5,
		HoleMargin:   [4]float64{5, 5, 5, 5},
		HolePattern:  [4]string{"x", "x", "x", "x"},
		Thickness:    3,
		Ridge:        v2.XY(2, 2),
	}
	nonEmptyShape(t, obj.Panel2D(p))
	nonEmptySolid(t, obj.Panel3D(p))
}

func TestEuroRackPanel(t *testing.T) {
	p := obj.EuroRackParms{
		U: 3, HP: 10, CornerRadius: 1.5, HoleDiameter: 0, Thickness: 2, Ridge: true,
	}
	nonEmptyShape(t, obj.EuroRackPanel2D(p))
	nonEmptySolid(t, obj.EuroRackPanel3D(p))
	// small HP path
	p2 := obj.EuroRackParms{U: 3, HP: 5, CornerRadius: 1.5, Thickness: 2}
	nonEmptyShape(t, obj.EuroRackPanel2D(p2))
}

func TestPanelHole3D(t *testing.T) {
	nonEmptySolid(t, obj.PanelHole3D(obj.PanelHoleParms{
		Diameter:    8,
		Thickness:   3,
		Indent:      v3.XYZ(2, 1, 5),
		Offset:      6,
		Orientation: 0,
	}))
}

// ---------- panel box ----------

func TestPanelBox3D(t *testing.T) {
	parts := obj.PanelBox3D(obj.PanelBoxParms{
		Size:       v3.XYZ(60, 40, 100),
		Wall:       2,
		Panel:      2,
		Rounding:   1.5,
		FrontInset: 1,
		BackInset:  1,
		Clearance:  0.05,
		Hole:       0,
		SideTabs:   "tb",
	})
	if len(parts) == 0 {
		t.Fatal("no parts")
	}
	for i, s := range parts {
		if s == nil {
			t.Fatalf("nil part %d", i)
		}
	}
}

// ---------- pipe ----------

func TestPipe(t *testing.T) {
	p := obj.PipeLookup("sch40:1", "inch")
	if p == nil || p.Outer <= p.Inner {
		t.Fatalf("PipeLookup unexpected: %+v", p)
	}
	nonEmptySolid(t, obj.Pipe3D(p.Outer, p.Inner, 30))
	nonEmptySolid(t, obj.StdPipe3D("sch40:1", "inch", 30))
	cfg := [6]bool{true, true, false, false, true, false}
	nonEmptySolid(t, obj.PipeConnector3D(obj.PipeConnectorParms{
		Length:        20,
		OuterRadius:   10,
		InnerRadius:   8,
		RecessDepth:   3,
		RecessWidth:   1,
		Configuration: cfg,
	}))
	nonEmptySolid(t, obj.StdPipeConnector3D("sch40:1", "inch", 50, cfg))
	expectPanic(t, "pipe bad lookup", func() {
		obj.PipeLookup("nosuch", "inch")
	})
}

// ---------- servo ----------

func TestServo(t *testing.T) {
	p := obj.ServoLookup("standard")
	if p == nil {
		t.Fatal("nil servo lookup")
	}
	nonEmptySolid(t, obj.Servo3D(*p))
	nonEmptyShape(t, obj.Servo2D(*p, -1))
	nonEmptyShape(t, obj.Servo2D(*p, 1.5))
	nonEmptyShape(t, obj.ServoHorn(obj.ServoHornParms{
		CenterRadius: 1.5, NumHoles: 4, CircleRadius: 8, HoleRadius: 0.8,
	}))
	expectPanic(t, "servo missing", func() { obj.ServoLookup("nope") })
	expectPanic(t, "servo horn bad", func() {
		obj.ServoHorn(obj.ServoHornParms{CenterRadius: -1})
	})
}

// ---------- shape (trapezoid / triangle) ----------

func TestIsoceles(t *testing.T) {
	nonEmptyShape(t, obj.IsocelesTrapezoid2D(20, 10, 8))
	nonEmptyShape(t, obj.IsocelesTriangle2D(10, 8))
}

// ---------- spring ----------

func TestSpring(t *testing.T) {
	p := obj.SpringParms{
		Width: 4, Height: 3, WallThickness: 0.6, Diameter: 4, NumSections: 3,
		Boss: [2]float64{2, 2},
	}
	if obj.SpringLength(p) <= 0 {
		t.Fatal("non-positive spring length")
	}
	nonEmptyShape(t, obj.Spring2D(p))
	nonEmptySolid(t, obj.Spring3D(p))
	expectPanic(t, "spring bad", func() {
		obj.Spring2D(obj.SpringParms{NumSections: 0})
	})
}

// ---------- truncated rect pyramid ----------

func TestTruncRectPyramid3D(t *testing.T) {
	nonEmptySolid(t, obj.TruncRectPyramid3D(obj.TruncRectPyramidParms{
		Size:         v3.XYZ(20, 10, 8),
		BaseAngleDeg: 70,
		BaseRadius:   1,
		RoundRadius:  0.5,
	}))
}

// ---------- arrow / axes ----------

func TestArrow(t *testing.T) {
	p := obj.ArrowParms{
		Axis:  [2]float64{20, 1},
		Head:  [2]float64{4, 2},
		Tail:  [2]float64{4, 2},
		Style: "cb",
	}
	nonEmptySolid(t, obj.Arrow3D(p))
	nonEmptySolid(t, obj.DirectedArrow3D(p, v3.XYZ(0, 0, 0), v3.XYZ(10, 10, 10)))
	nonEmptySolid(t, obj.Axes3D(v3.XYZ(0, 0, 0), v3.XYZ(10, 10, 10)))
	expectPanic(t, "arrow bad", func() {
		obj.Arrow3D(obj.ArrowParms{Style: "abcd"})
	})
}

// ---------- tabs ----------

func TestTabs(t *testing.T) {
	straight := obj.NewStraightTab(v3.XYZ(4, 4, 2), 0.1)
	angle := obj.NewAngleTab(v3.XYZ(4, 4, 2), 0.1)
	screw := obj.NewScrewTab(obj.ScrewTab{
		Length: 6, Radius: 2, Round: true, HoleUpper: 4, HoleLower: 2, HoleRadius: 0.8,
	})
	if straight == nil || angle == nil || screw == nil {
		t.Fatal("nil tab")
	}
	base := solid.Box(v3.XYZ(20, 20, 10), 0)
	mset := []solid.M44{solid.Identity3d()}
	nonEmptySolid(t, obj.AddTabs(base, straight, true, mset))
	nonEmptySolid(t, obj.AddTabs(base, angle, false, mset))
	nonEmptySolid(t, obj.AddTabs(base, screw, false, mset))
}
