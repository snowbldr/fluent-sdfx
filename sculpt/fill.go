package sculpt

import (
	"fmt"
	"math"
	"sort"

	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
)

// flood walks the lattice from the corner at start over corners that share
// its side (material or air), six-connected. It returns the visited mask,
// how many corners it reached, and whether the walk touched the lattice
// boundary, which for an air walk means the region is open to the outside.
func (g *Grid) flood(start int, inside bool) (visited []bool, count int, hitBoundary bool) {
	visited = make([]bool, len(g.vals))
	nx1, ny1 := g.nx+1, g.ny+1
	stack := []int{start}
	visited[start] = true
	for len(stack) > 0 {
		id := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		count++
		i := id % nx1
		j := (id / nx1) % ny1
		k := id / (nx1 * ny1)
		if i == 0 || j == 0 || k == 0 || i == g.nx || j == g.ny || k == g.nz {
			hitBoundary = true
		}
		try := func(nb int) {
			if !visited[nb] && (g.vals[nb] < 0) == inside {
				visited[nb] = true
				stack = append(stack, nb)
			}
		}
		if i > 0 {
			try(id - 1)
		}
		if i < g.nx {
			try(id + 1)
		}
		if j > 0 {
			try(id - nx1)
		}
		if j < g.ny {
			try(id + nx1)
		}
		if k > 0 {
			try(id - nx1*ny1)
		}
		if k < g.nz {
			try(id + nx1*ny1)
		}
	}
	return visited, count, hitBoundary
}

// seedCorner finds the corner at seed and checks it is on the expected side
// of the surface.
func (g *Grid) seedCorner(verb string, seed v3.Vec, wantInside bool) int {
	id, ok := g.corner(seed)
	if !ok {
		panic(fmt.Sprintf("sculpt.%s: seed (%.3g, %.3g, %.3g) is outside the sampled region", verb, seed.X, seed.Y, seed.Z))
	}
	if (g.vals[id] < 0) != wantInside {
		side := "in air"
		if !wantInside {
			side = "inside material"
		}
		panic(fmt.Sprintf("sculpt.%s: seed (%.3g, %.3g, %.3g) is %s (field %.3g); it must be on the other side of the surface", verb, seed.X, seed.Y, seed.Z, side, g.vals[id]))
	}
	return id
}

// Enclosed reports whether the air at seed is sealed off from the outside:
// whether a bucket fill started there would stop at walls on every side.
// Panics if seed is inside material.
func (g *Grid) Enclosed(seed v3.Vec) bool {
	_, _, open := g.flood(g.seedCorner("Enclosed", seed, false), false)
	return !open
}

// Fill is the bucket fill: the air connected to seed becomes material, out
// to the walls that bound it in every direction. The use is a void that
// has to be solid, such as the inside of a shelled form that will be cast
// or moulded, or a pocket a blend left behind a surface.
//
// The region must be enclosed. An opening to the outside, however small,
// would let the fill consume all the air in the grid, so Fill panics
// instead; Enclosed asks the same question without panicking. To fill a
// pocket that is open on one side, cap the opening first (union a plate
// over it), fill, then cut the cap away. Panics if seed is inside material.
func (g *Grid) Fill(seed v3.Vec) *Grid {
	start := g.seedCorner("Fill", seed, false)
	visited, _, open := g.flood(start, false)
	if open {
		panic(fmt.Sprintf("sculpt.Fill: the air at (%.3g, %.3g, %.3g) is open to the outside, so the fill would consume everything; cap the opening first, or check with Enclosed", seed.X, seed.Y, seed.Z))
	}
	out := g.clone()
	for i, v := range visited {
		if v {
			out.vals[i] = float32(-g.h)
		}
	}
	return out.Redistance()
}

// Component is one connected body of material in a Grid.
type Component struct {
	Cells  int     // interior corners
	Volume float64 // mm³, Cells times the cell volume
	Bounds v3.Box  // of the interior corners
	Deep   v3.Vec  // the corner farthest inside: a seed that is surely in this body
}

// Components lists the connected bodies of material, largest first. A body
// is corners of material joined face to face; a single corner of material
// is a body. Two bodies that touch only at an edge or a vertex count as
// separate, which is also how they would print.
func (g *Grid) Components() []Component {
	labelled := make([]bool, len(g.vals))
	var out []Component
	for id, v := range g.vals {
		if v >= 0 || labelled[id] {
			continue
		}
		visited, count, _ := g.flood(id, true)
		c := Component{Cells: count, Volume: float64(count) * g.h * g.h * g.h}
		deep := float32(0)
		first := true
		for i, in := range visited {
			if !in {
				continue
			}
			labelled[i] = true
			p := g.posOf(i)
			if first {
				c.Bounds = v3.Box{Min: p, Max: p}
				first = false
			} else {
				c.Bounds = c.Bounds.Include(p)
			}
			if g.vals[i] < deep {
				deep = g.vals[i]
				c.Deep = p
			}
		}
		out = append(out, c)
	}
	sort.Slice(out, func(a, b int) bool { return out[a].Cells > out[b].Cells })
	return out
}

// posOf is the world position of the corner with flat index id.
func (g *Grid) posOf(id int) v3.Vec {
	nx1, ny1 := g.nx+1, g.ny+1
	return g.pos(id%nx1, (id/nx1)%ny1, id/(nx1*ny1))
}

// Keep returns only the body of material that contains seed; every other
// body becomes air. The use is dropping what is not attached: a bead a
// blend did not reach, a sliver a cut left floating, the islands a
// sculpted form sheds. Panics if seed is in air.
func (g *Grid) Keep(seed v3.Vec) *Grid {
	start := g.seedCorner("Keep", seed, true)
	visited, _, _ := g.flood(start, true)
	out := g.clone()
	changed := false
	for i, v := range g.vals {
		if v < 0 && !visited[i] {
			out.vals[i] = float32(g.h)
			changed = true
		}
	}
	if !changed {
		return out
	}
	return out.Redistance()
}

// Smooth rounds every edge and fills every crease of the form at radius r:
// the result is what a ball of radius r can roll over on the outside and
// on the inside. It is two morphological passes with a redistance between
// each step. The first shrinks by r and grows back, which rounds convex
// edges and removes anything thinner than 2r. The second grows by r and
// shrinks back, which fills concave creases and closes any gap narrower
// than 2r. Flat faces and surfaces already gentler than r are unchanged.
//
// That loss of thin features is the meaning of the operation, not a side
// effect: a burr vanishes, and so does a stem thinner than 2r. Choose r
// below half the thinnest feature that must survive, and above the cell
// size, or there is nothing to round.
func (g *Grid) Smooth(r float64) *Grid {
	if r <= 0 {
		panic(fmt.Sprintf("sculpt.Smooth: radius must be positive, got %g", r))
	}
	if r < g.h {
		panic(fmt.Sprintf("sculpt.Smooth: radius %g is below the cell size %g; the grid cannot represent the rounding", r, g.h))
	}
	w := g.pad(int(math.Ceil(r/g.h)) + 1).Redistance()
	w = w.shift(r).Redistance()
	w = w.shift(-r).Redistance()
	w = w.shift(-r).Redistance()
	w = w.shift(r).Redistance()
	return w
}
