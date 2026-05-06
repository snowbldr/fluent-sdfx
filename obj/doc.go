// Package obj provides parametric helper constructors for parts you'd
// otherwise build by hand. Each helper takes a Parms struct and returns a
// ready-to-use *solid.Solid (or *shape.Shape for 2D).
//
// Index by category:
//
// Fasteners
//   - Bolt, Nut, Washer3D — hex/socket/cap-head bolts, matching nuts, washers
//   - HexHead3D, KnurledHead3D, SocketHead3D — bare bolt heads
//   - CounterBoredHole3D, CounterSunkHole3D — drilled holes with relief
//
// Threads
//   - Thread profiles: ANSIThread, AcmeThread, BSPThread, BuckleThread,
//     ISOThread, KnuckleThread, NPTThread, PlasticThread, etc.
//
// Enclosures and panels
//   - Panel2D, Panel3D — flat panels with corner radii and hole patterns
//   - PanelBox — box assembled from 6 panels with finger joints
//   - Standoff3D — PCB-mounting pillars with hole + chamfer
//   - EuroRackPanel3D — Eurorack synth module panels (HP × U)
//   - PanelHole3D — generic panel cutout with optional indent
//
// Mechanical
//   - Angle3D — angle bracket with chamfered root
//   - InvoluteGear, BevelGear — gear profiles (2D, then Extrude)
//   - Pipe3D, ChamferedPipe3D — through-bore tubing
//   - PCBHole3D — PCB through-holes with annular ring
//   - DrainCover — drain covers with optional cross-bar
//
// Storage
//   - GridfinityBase, GridfinityBody — Gridfinity-compatible bases and bodies
//
// Miscellaneous
//   - Arrow3D, BoltCircle3D — utility shapes
//   - Servo (HitecHs225, etc.) — servo body models
//   - Text3D — extruded text from a font
//
// See the API reference (https://snowbldr.github.io/fluent-sdfx/api-reference/)
// for the full list and per-helper Parms struct documentation.
package obj
