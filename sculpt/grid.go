package sculpt

import (
	"fmt"
	"math"
	"runtime"
	"sync"

	"github.com/snowbldr/fluent-sdfx/solid"
	v3 "github.com/snowbldr/fluent-sdfx/vec/v3"
	"github.com/snowbldr/sdfx/sdf"
	v3sdf "github.com/snowbldr/sdfx/vec/v3"
)

// Grid is a solid's signed distance sampled at the corners of a cubic
// lattice. It is the working material for the grid verbs: Redistance,
// Fill, Keep, Components and Smooth all read and write the corner samples,
// and Solid returns the lattice as a solid again by trilinear
// interpolation.
//
// Every verb returns a new Grid and leaves the receiver unchanged.
type Grid struct {
	vals       []float32 // (nx+1)(ny+1)(nz+1) corner samples, x fastest
	nx, ny, nz int       // cells per axis
	min        v3.Vec    // position of corner (0,0,0)
	h          float64   // cell edge, mm
	src        v3.Box    // bounds reported by Solid()
}

// samplePad is how many cells of air surround the solid's bounds when it
// is sampled, so the surface never touches the lattice boundary and the
// flood walks have an outside to reach.
const samplePad = 3

// Sample evaluates s at the corners of a cubic lattice with `cells` cells
// along the longest side of its bounds, padded by three cells of air on
// every side. The fill runs on all CPUs.
//
// The cell size is what every later verb is accurate to: choose cells for
// the smallest feature that must survive, not for the size of the part.
// A 128-cell grid on a 40 mm part is a 0.3 mm cell and about two million
// samples.
func Sample(s *solid.Solid, cells int) *Grid {
	if cells < 2 {
		panic(fmt.Sprintf("sculpt.Sample: cells must be at least 2, got %d", cells))
	}
	bb := s.Bounds().Box
	h := bb.Size().MaxComponent() / float64(cells)
	if h <= 0 {
		panic("sculpt.Sample: the solid's bounds have no extent")
	}
	g := newGrid(bb, h, samplePad)
	g.src = bb
	g.fill(func(p v3.Vec) float64 { return s.Evaluate(p.Raw()) })
	return g
}

// newGrid lays out an empty lattice of cell size h covering bb with pad
// cells beyond it on every side. The lattice is centred on bb and rounded
// out to whole cells.
func newGrid(bb v3.Box, h float64, pad int) *Grid {
	size := bb.Size()
	nx := int(math.Ceil(size.X/h-1e-9)) + 2*pad
	ny := int(math.Ceil(size.Y/h-1e-9)) + 2*pad
	nz := int(math.Ceil(size.Z/h-1e-9)) + 2*pad
	full := v3.XYZ(float64(nx)*h, float64(ny)*h, float64(nz)*h)
	return &Grid{
		vals: make([]float32, (nx+1)*(ny+1)*(nz+1)),
		nx:   nx, ny: ny, nz: nz,
		min: bb.Center().Sub(full.MulScalar(0.5)),
		h:   h,
		src: bb,
	}
}

// CellSize is the lattice spacing in mm: the resolution every grid verb
// works to.
func (g *Grid) CellSize() float64 { return g.h }

// Cells is the number of cells along each axis.
func (g *Grid) Cells() [3]int { return [3]int{g.nx, g.ny, g.nz} }

func (g *Grid) idx(i, j, k int) int { return (k*(g.ny+1)+j)*(g.nx+1) + i }

// pos is the world position of corner (i,j,k).
func (g *Grid) pos(i, j, k int) v3.Vec {
	return g.min.Add(v3.XYZ(float64(i)*g.h, float64(j)*g.h, float64(k)*g.h))
}

// corner returns the lattice corner nearest p and whether p lies inside the
// lattice at all.
func (g *Grid) corner(p v3.Vec) (int, bool) {
	f := p.Sub(g.min).DivScalar(g.h)
	i, j, k := int(math.Round(f.X)), int(math.Round(f.Y)), int(math.Round(f.Z))
	if i < 0 || j < 0 || k < 0 || i > g.nx || j > g.ny || k > g.nz {
		return 0, false
	}
	return g.idx(i, j, k), true
}

func (g *Grid) clone() *Grid {
	c := *g
	c.vals = make([]float32, len(g.vals))
	copy(c.vals, g.vals)
	return &c
}

// fill evaluates f at every corner, one z-slab per goroutine.
func (g *Grid) fill(f func(v3.Vec) float64) {
	slabs := make(chan int, g.nz+1)
	for k := 0; k <= g.nz; k++ {
		slabs <- k
	}
	close(slabs)
	var wg sync.WaitGroup
	for w := 0; w < runtime.NumCPU(); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := range slabs {
				for j := 0; j <= g.ny; j++ {
					for i := 0; i <= g.nx; i++ {
						g.vals[g.idx(i, j, k)] = float32(f(g.pos(i, j, k)))
					}
				}
			}
		}()
	}
	wg.Wait()
}

// shift adds delta to every sample: a positive delta shrinks the solid by
// delta, a negative one grows it. The result is a distance bound, not a
// distance, until it is redistanced.
func (g *Grid) shift(delta float64) *Grid {
	out := g.clone()
	d := float32(delta)
	for i := range out.vals {
		out.vals[i] += d
	}
	return out
}

// pad returns the grid with n more cells of air on every side. The new
// corners hold a lower bound on their distance to the surface, so the
// result is safe to sample from, and Redistance makes them exact.
func (g *Grid) pad(n int) *Grid {
	if n <= 0 {
		return g.clone()
	}
	out := &Grid{
		vals: make([]float32, (g.nx+2*n+1)*(g.ny+2*n+1)*(g.nz+2*n+1)),
		nx:   g.nx + 2*n, ny: g.ny + 2*n, nz: g.nz + 2*n,
		min: g.min.Sub(v3.XYZ(float64(n)*g.h, float64(n)*g.h, float64(n)*g.h)),
		h:   g.h,
		src: g.src,
	}
	edge := g.boundaryMin()
	lo := g.min
	hi := g.pos(g.nx, g.ny, g.nz)
	for k := 0; k <= out.nz; k++ {
		for j := 0; j <= out.ny; j++ {
			for i := 0; i <= out.nx; i++ {
				si, sj, sk := i-n, j-n, k-n
				if si >= 0 && sj >= 0 && sk >= 0 && si <= g.nx && sj <= g.ny && sk <= g.nz {
					out.vals[out.idx(i, j, k)] = g.vals[g.idx(si, sj, sk)]
					continue
				}
				p := out.pos(i, j, k)
				q := p.Clamp(lo, hi)
				out.vals[out.idx(i, j, k)] = float32(p.Sub(q).Length() + edge)
			}
		}
	}
	return out
}

// boundaryMin is the smallest sample on the six faces of the lattice,
// floored at zero: the least distance any boundary point has to the
// surface.
func (g *Grid) boundaryMin() float64 {
	m := math.Inf(1)
	at := func(i, j, k int) {
		if v := float64(g.vals[g.idx(i, j, k)]); v < m {
			m = v
		}
	}
	for k := 0; k <= g.nz; k++ {
		for j := 0; j <= g.ny; j++ {
			at(0, j, k)
			at(g.nx, j, k)
		}
		for i := 0; i <= g.nx; i++ {
			at(i, 0, k)
			at(i, g.ny, k)
		}
	}
	for j := 0; j <= g.ny; j++ {
		for i := 0; i <= g.nx; i++ {
			at(i, j, 0)
			at(i, j, g.nz)
		}
	}
	return math.Max(m, 0)
}

// Volume is the material volume in mm³. Each corner contributes the
// fraction of a cell its sample says is material: all of it half a cell or
// more inside, none of it half a cell or more outside, and a linear share
// in between, so a corner exactly on the surface counts half. That makes
// the figure accurate to a fraction of a cell of surface, enough to tell a
// filled void from an empty one and a dropped bead from a kept one.
func (g *Grid) Volume() float64 {
	sum := 0.0
	for _, v := range g.vals {
		frac := 0.5 - float64(v)/g.h
		if frac > 1 {
			frac = 1
		} else if frac < 0 {
			frac = 0
		}
		sum += frac
	}
	return sum * g.h * g.h * g.h
}

// Solid returns the grid as a solid: trilinear interpolation between the
// corner samples inside the lattice, and a lower bound on the distance
// outside it. The reported bounds are those of the solid the grid was
// sampled from.
func (g *Grid) Solid() *solid.Solid {
	return solid.Wrap(&gridSDF{g: g, edge: g.boundaryMin(), lo: g.min, hi: g.pos(g.nx, g.ny, g.nz)})
}

type gridSDF struct {
	g      *Grid
	edge   float64
	lo, hi v3.Vec
}

func (s *gridSDF) BoundingBox() sdf.Box3 { return s.g.src.SDF() }

func (s *gridSDF) Evaluate(p v3sdf.Vec) float64 {
	pp := v3.Vec(p)
	q := pp.Clamp(s.lo, s.hi)
	if q != pp {
		// Any path to the surface crosses the lattice boundary, so the
		// distance is at least the distance to the boundary plus the least
		// boundary value; and by the triangle inequality at least the value
		// at the clamped point less the distance to it. Take the larger.
		out := pp.Sub(q).Length()
		return math.Max(out+s.edge, s.g.interp(q)-out)
	}
	return s.g.interp(pp)
}

// interp is the trilinear reconstruction at a point inside the lattice.
func (g *Grid) interp(p v3.Vec) float64 {
	f := p.Sub(g.min).DivScalar(g.h)
	i, j, k := int(f.X), int(f.Y), int(f.Z)
	if i >= g.nx {
		i = g.nx - 1
	}
	if j >= g.ny {
		j = g.ny - 1
	}
	if k >= g.nz {
		k = g.nz - 1
	}
	dx, dy, dz := f.X-float64(i), f.Y-float64(j), f.Z-float64(k)
	at := func(i, j, k int) float64 { return float64(g.vals[g.idx(i, j, k)]) }
	c00 := at(i, j, k)*(1-dx) + at(i+1, j, k)*dx
	c10 := at(i, j+1, k)*(1-dx) + at(i+1, j+1, k)*dx
	c01 := at(i, j, k+1)*(1-dx) + at(i+1, j, k+1)*dx
	c11 := at(i, j+1, k+1)*(1-dx) + at(i+1, j+1, k+1)*dx
	c0 := c00*(1-dy) + c10*dy
	c1 := c01*(1-dy) + c11*dy
	return c0*(1-dz) + c1*dz
}

// touchesOtherSide reports whether any of the 26 corners around (i,j,k)
// has the opposite sign.
func (g *Grid) touchesOtherSide(i, j, k int) bool {
	in := g.vals[g.idx(i, j, k)] < 0
	for dk := -1; dk <= 1; dk++ {
		kk := k + dk
		if kk < 0 || kk > g.nz {
			continue
		}
		for dj := -1; dj <= 1; dj++ {
			jj := j + dj
			if jj < 0 || jj > g.ny {
				continue
			}
			for di := -1; di <= 1; di++ {
				ii := i + di
				if ii < 0 || ii > g.nx {
					continue
				}
				if (g.vals[g.idx(ii, jj, kk)] < 0) != in {
					return true
				}
			}
		}
	}
	return false
}

// gradient is the field's slope at corner (i,j,k) by central differences,
// one-sided at the lattice boundary, in field units per mm.
func (g *Grid) gradient(i, j, k int) v3.Vec {
	at := func(i, j, k int) float64 { return float64(g.vals[g.idx(i, j, k)]) }
	diff := func(lo, hi int, n int, f func(int) float64) float64 {
		a, b := lo, hi
		if a < 0 {
			a = 0
		}
		if b > n {
			b = n
		}
		return (f(b) - f(a)) / (float64(b-a) * g.h)
	}
	return v3.XYZ(
		diff(i-1, i+1, g.nx, func(x int) float64 { return at(x, j, k) }),
		diff(j-1, j+1, g.ny, func(y int) float64 { return at(i, y, k) }),
		diff(k-1, k+1, g.nz, func(z int) float64 { return at(i, j, z) }),
	)
}

// Redistance recomputes every sample as the true signed distance to the
// surface the samples currently describe. The surface is located where
// neighbouring samples change sign, by linear interpolation along the
// lattice edge between them; those crossing points are the seeds. Every
// corner then finds its nearest seed with a separable Euclidean distance
// transform and takes its distance to that seed's crossing point. The sign
// of each sample is kept.
//
// This is what makes the rest of the package honest. A field that has
// been through a Pull, a twist, a smooth blend or a difference is only a
// bound on the distance; a shell or an offset built on it comes out the
// wrong thickness. After Redistance the magnitudes are right to within a
// fraction of a cell everywhere, deep interior included.
//
// A grid with no surface in it (all material or all air) is returned
// unchanged.
func (g *Grid) Redistance() *Grid {
	f, cpOf, cps, band := g.seeds()
	if !band {
		return g.clone()
	}
	src := make([]int32, len(f))
	for i := range src {
		src[i] = int32(i)
	}
	distanceTransform(f, src, g.nx+1, g.ny+1, g.nz+1)
	best := g.refine(src, cpOf, cps)
	out := g.clone()
	for i, v := range g.vals {
		d := math.Sqrt(best[i])
		if v < 0 {
			d = -d
		}
		out.vals[i] = float32(d)
	}
	return out
}

// refine corrects the transform's choice of seed. The transform ranks
// seeds by |p−q|² + f(q), the distance to the seed corner lifted by its own
// distance to the surface, which is only an estimate of the distance to the
// seed's surface point and can pick a neighbour's seed over the right one
// by a fraction of a cell. Two sweeps over the lattice, forward and back,
// let every corner compare its seed with those of the 26 corners around it
// and keep the nearest surface point. On return src holds the corrected
// seeds and the result is each corner's squared distance to its surface
// point, in mm².
func (g *Grid) refine(src []int32, cpOf []int32, cps []v3.Vec) []float64 {
	n := len(src)
	best := make([]float64, n)
	for i := range src {
		best[i] = g.posOf(i).Sub(cps[cpOf[src[i]]]).Length2()
	}
	nx1, ny1 := g.nx+1, g.ny+1
	// The 13 neighbours that precede a corner in sweep order; negated,
	// the 13 that follow it.
	type off struct{ di, dj, dk int }
	var before []off
	for dk := -1; dk <= 1; dk++ {
		for dj := -1; dj <= 1; dj++ {
			for di := -1; di <= 1; di++ {
				if dk < 0 || (dk == 0 && dj < 0) || (dk == 0 && dj == 0 && di < 0) {
					before = append(before, off{di, dj, dk})
				}
			}
		}
	}
	visit := func(i, j, k int, offs []off) {
		id := g.idx(i, j, k)
		p := g.pos(i, j, k)
		for _, o := range offs {
			ii, jj, kk := i+o.di, j+o.dj, k+o.dk
			if ii < 0 || jj < 0 || kk < 0 || ii > g.nx || jj > g.ny || kk > g.nz {
				continue
			}
			cand := src[(kk*ny1+jj)*nx1+ii]
			if cand == src[id] {
				continue
			}
			if d := p.Sub(cps[cpOf[cand]]).Length2(); d < best[id] {
				best[id] = d
				src[id] = cand
			}
		}
	}
	for k := 0; k <= g.nz; k++ {
		for j := 0; j <= g.ny; j++ {
			for i := 0; i <= g.nx; i++ {
				visit(i, j, k, before)
			}
		}
	}
	after := make([]off, len(before))
	for i, o := range before {
		after[i] = off{-o.di, -o.dj, -o.dk}
	}
	for k := g.nz; k >= 0; k-- {
		for j := g.ny; j >= 0; j-- {
			for i := g.nx; i >= 0; i-- {
				visit(i, j, k, after)
			}
		}
	}
	return best
}

// seeds locates the surface. For every corner within reach of a sign
// change it returns, in f, the squared distance (in cells) to the nearest
// point of the surface it can see and, through cpOf into cps, that point;
// f is +Inf and cpOf −1 elsewhere. band is false when there is no surface.
func (g *Grid) seeds() (f []float64, cpOf []int32, cps []v3.Vec, band bool) {
	n := len(g.vals)
	inf := math.Inf(1)
	f = make([]float64, n)
	cpOf = make([]int32, n)
	for i := range f {
		f[i] = inf
		cpOf[i] = -1
	}
	nx1, ny1 := g.nx+1, g.ny+1
	// seed records the crossing between corners a and b (b is a's neighbour
	// along axis) for whichever of the two it is nearer than its current
	// crossing.
	seed := func(a, b int, axis v3.Vec) {
		va, vb := float64(g.vals[a]), float64(g.vals[b])
		if (va < 0) == (vb < 0) {
			return
		}
		band = true
		t := math.Abs(va) / (math.Abs(va) + math.Abs(vb)) // cells from a to the crossing
		cp := g.posOf(a).Add(axis.MulScalar(t * g.h))
		if t*t < f[a] {
			f[a] = t * t
			cpOf[a] = int32(len(cps))
			cps = append(cps, cp)
		}
		if (1-t)*(1-t) < f[b] {
			f[b] = (1 - t) * (1 - t)
			cpOf[b] = int32(len(cps))
			cps = append(cps, cp)
		}
	}
	ex, ey, ez := v3.X(1), v3.Y(1), v3.Z(1)
	for k := 0; k <= g.nz; k++ {
		for j := 0; j <= g.ny; j++ {
			for i := 0; i <= g.nx; i++ {
				id := g.idx(i, j, k)
				if i < g.nx {
					seed(id, id+1, ex)
				}
				if j < g.ny {
					seed(id, id+nx1, ey)
				}
				if k < g.nz {
					seed(id, id+nx1*ny1, ez)
				}
			}
		}
	}
	if !band {
		return
	}
	// An axis crossing is a point on the surface, so its distance is an
	// upper bound, and a loose one where the surface runs obliquely to the
	// lattice. Where the field's slope at a corner is well defined, project
	// along the gradient instead and keep whichever is nearer. This also
	// seeds corners that only touch the other side diagonally: on a sharp
	// convex edge no lattice edge crosses the surface, and without these
	// seeds everything in the fan outside the edge would measure to the
	// nearest face instead.
	limit := 3 * g.h * g.h // (√3 h)²: a diagonal neighbour is this far at most
	for k := 0; k <= g.nz; k++ {
		for j := 0; j <= g.ny; j++ {
			for i := 0; i <= g.nx; i++ {
				id := g.idx(i, j, k)
				if cpOf[id] < 0 && !g.touchesOtherSide(i, j, k) {
					continue
				}
				grad := g.gradient(i, j, k)
				l2 := grad.Length2()
				if l2 == 0 {
					continue
				}
				va := float64(g.vals[id])
				d2 := va * va / l2
				if d2 > limit || d2/(g.h*g.h) >= f[id] {
					continue
				}
				f[id] = d2 / (g.h * g.h)
				cpOf[id] = int32(len(cps))
				cps = append(cps, g.pos(i, j, k).Sub(grad.MulScalar(va/l2)))
			}
		}
	}
	return
}

// distanceTransform finds, for every corner of a lattice, the seed corner
// minimising |p−q|² + f(q), where f is the seed's own squared cost in cell
// units (+Inf for a non-seed). On return f holds that minimum and src the
// index of the seed that achieved it. It is the separable lower-envelope
// algorithm of Felzenszwalb and Huttenlocher, one pass per axis with the
// argmin carried through, rows spread across all CPUs.
func distanceTransform(f []float64, src []int32, nx, ny, nz int) {
	pass := func(n, stride, rows int, rowStart func(r int) int) {
		work := make(chan int, rows)
		for r := 0; r < rows; r++ {
			work <- r
		}
		close(work)
		var wg sync.WaitGroup
		for w := 0; w < runtime.NumCPU(); w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				in := make([]float64, n)
				inSrc := make([]int32, n)
				out := make([]float64, n)
				arg := make([]int, n)
				v := make([]int, n)
				z := make([]float64, n+1)
				for r := range work {
					base := rowStart(r)
					for q := 0; q < n; q++ {
						in[q] = f[base+q*stride]
						inSrc[q] = src[base+q*stride]
					}
					dt1d(in, out, arg, v, z)
					for q := 0; q < n; q++ {
						f[base+q*stride] = out[q]
						if !math.IsInf(out[q], 1) {
							src[base+q*stride] = inSrc[arg[q]]
						}
					}
				}
			}()
		}
		wg.Wait()
	}
	// Along x: one row per (j,k).
	pass(nx, 1, ny*nz, func(r int) int { return r * nx })
	// Along y: one row per (i,k).
	pass(ny, nx, nx*nz, func(r int) int { i, k := r%nx, r/nx; return k*nx*ny + i })
	// Along z: one row per (i,j).
	pass(nz, nx*ny, nx*ny, func(r int) int { return r })
}

// dt1d is the one-dimensional squared distance transform with arbitrary
// seed costs: d[q] = min over p of (q−p)² + f[p], with arg[q] the p that
// achieved it. Infinite f are skipped. v and z are scratch of length n
// and n+1.
func dt1d(f, d []float64, arg, v []int, z []float64) {
	n := len(f)
	start := -1
	for q := 0; q < n; q++ {
		if !math.IsInf(f[q], 1) {
			start = q
			break
		}
	}
	if start < 0 {
		for q := range d {
			d[q] = math.Inf(1)
			arg[q] = q
		}
		return
	}
	k := 0
	v[0] = start
	z[0] = math.Inf(-1)
	z[1] = math.Inf(1)
	for q := start + 1; q < n; q++ {
		if math.IsInf(f[q], 1) {
			continue
		}
		fq := f[q] + float64(q)*float64(q)
		for {
			vk := v[k]
			s := (fq - (f[vk] + float64(vk)*float64(vk))) / float64(2*(q-vk))
			if s <= z[k] {
				k--
				continue
			}
			k++
			v[k] = q
			z[k] = s
			z[k+1] = math.Inf(1)
			break
		}
	}
	k = 0
	for q := 0; q < n; q++ {
		for z[k+1] < float64(q) {
			k++
		}
		dq := float64(q - v[k])
		d[q] = dq*dq + f[v[k]]
		arg[q] = v[k]
	}
}
