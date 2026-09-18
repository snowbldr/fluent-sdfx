# Hidden spec (only revealed when the model asks)

"Triangular" means a **solid corner gusset** (a corbel), not a triangular
spoke section.

- The gusset is a wedge that fills the corner UNDER each spoke, between the
  spoke's underside (z = 10) and the hub wall (r = 6).
- It runs the spoke's **full length**, right out to the tip, so every layer
  of the spoke has something under it.
- Its sloped face is at **45 degrees**: from the spoke's outer tip
  (x = 16, z = 10) down and in to the hub wall (x = 6, z = 0). The gusset's
  root therefore reaches the bottom of the hub.
- The part prints with the **hub axis vertical**, hub bottom on the plate.
  45 degrees is the self-support limit, hence the angle.
- The spokes' own **rectangular cross-section does not change**: still 3 wide
  (tangential), 3 tall, 10 long past the hub surface, underside at z = 10,
  top at z = 13, square corners.
- The hub does not change: r 6, z 0..20.
- The gusset is the same width as the spoke (3) and is merged into the hub,
  not left floating next to it.
- Anything else: your call, use a sensible default.
