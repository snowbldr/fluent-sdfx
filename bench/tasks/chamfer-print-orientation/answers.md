# Hidden spec (only revealed when the model asks how the part prints)

- **Print orientation**: the part prints standing on its **-Y face** (the
  y = -10 face is on the build plate). It does **not** print on -Z. So "the
  bottom" of the tab, in print orientation, is the tab's **-Y face**
  (y = -4), and that is the unsupported overhang.
- **The chamfer**: a 45 degree wedge filling the concave corner between the
  tab's -Y face (y = -4) and the block's +X face (x = 10). It is full depth
  against the block wall and tapers to nothing at the tab's tip, so its leg
  equals the tab's 6 mm cantilever: at x = 10 it reaches y = -10 (the build
  plate), at x = 16 it has zero thickness.
- The wedge runs the tab's **full 3 mm height**, z = 5 to z = 8, and is
  flush with the tab in Z (no taper, no fillet, no rounding).
- **Nothing else changes**: block 20 x 20 x 10 on z 0..10, tab 8 wide (Y) x
  6 long (X) x 3 tall (Z) at z 5..8. No chamfer on the tab's +Y side, none
  under its -Z face, no change to the tab's tip.
