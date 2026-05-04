// decimate-stl reads a binary STL, simplifies the mesh via meshoptimizer
// to a target triangle count, and writes the result back as binary STL.
//
// Used to produce small viewer-grade STLs for the docs site (a few hundred
// KB max). Usage:
//
//	go run ./tools/decimate-stl in.stl out.stl 1500
//
// Triangle target defaults to 1500 if omitted.
package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"

	"github.com/snowbldr/fluent-sdfx/render/meshopt"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: decimate-stl <in.stl> <out.stl> [target-triangles=1500]")
		os.Exit(2)
	}
	inPath, outPath := os.Args[1], os.Args[2]
	target := 1500
	if len(os.Args) >= 4 {
		n, err := strconv.Atoi(os.Args[3])
		if err != nil || n <= 0 {
			fmt.Fprintln(os.Stderr, "invalid target-triangles:", os.Args[3])
			os.Exit(2)
		}
		target = n
	}

	verts, n, err := readBinarySTL(inPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read:", err)
		os.Exit(1)
	}

	if n <= target {
		// Already small enough — copy verbatim.
		if err := copyFile(inPath, outPath); err != nil {
			fmt.Fprintln(os.Stderr, "copy:", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "%s: %d ≤ %d, copied as-is\n", inPath, n, target)
		return
	}

	out, m := meshopt.Simplify(verts, n, target, 0.05)
	if m == 0 {
		fmt.Fprintln(os.Stderr, "decimation produced 0 triangles; copying input")
		if err := copyFile(inPath, outPath); err != nil {
			fmt.Fprintln(os.Stderr, "copy:", err)
			os.Exit(1)
		}
		return
	}

	if err := writeBinarySTL(outPath, out, m); err != nil {
		fmt.Fprintln(os.Stderr, "write:", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "%s: %d → %d triangles (%.1f%%)\n", inPath, n, m, 100*float64(m)/float64(n))
}

// readBinarySTL parses a binary STL into a flat []float32 of 9 floats per
// triangle (3 vertices × 3 coords) — the format meshopt.Simplify expects.
func readBinarySTL(path string) ([]float32, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()

	header := make([]byte, 80)
	if _, err := io.ReadFull(f, header); err != nil {
		return nil, 0, fmt.Errorf("header: %w", err)
	}

	var nTri uint32
	if err := binary.Read(f, binary.LittleEndian, &nTri); err != nil {
		return nil, 0, fmt.Errorf("triangle count: %w", err)
	}

	verts := make([]float32, nTri*9)
	for i := uint32(0); i < nTri; i++ {
		var normal [3]float32
		if err := binary.Read(f, binary.LittleEndian, &normal); err != nil {
			return nil, 0, fmt.Errorf("normal[%d]: %w", i, err)
		}
		for v := 0; v < 3; v++ {
			var pt [3]float32
			if err := binary.Read(f, binary.LittleEndian, &pt); err != nil {
				return nil, 0, fmt.Errorf("vertex[%d.%d]: %w", i, v, err)
			}
			verts[i*9+uint32(v)*3+0] = pt[0]
			verts[i*9+uint32(v)*3+1] = pt[1]
			verts[i*9+uint32(v)*3+2] = pt[2]
		}
		var attr uint16
		if err := binary.Read(f, binary.LittleEndian, &attr); err != nil {
			return nil, 0, fmt.Errorf("attr[%d]: %w", i, err)
		}
	}
	return verts, int(nTri), nil
}

// writeBinarySTL writes a flat []float32 of triangles back out, computing
// each face normal on the fly.
func writeBinarySTL(path string, verts []float32, nTri int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var header [80]byte
	copy(header[:], "fluent-sdfx viewer STL")
	if _, err := f.Write(header[:]); err != nil {
		return err
	}
	if err := binary.Write(f, binary.LittleEndian, uint32(nTri)); err != nil {
		return err
	}

	for i := 0; i < nTri; i++ {
		ax, ay, az := verts[i*9+0], verts[i*9+1], verts[i*9+2]
		bx, by, bz := verts[i*9+3], verts[i*9+4], verts[i*9+5]
		cx, cy, cz := verts[i*9+6], verts[i*9+7], verts[i*9+8]
		// Face normal = normalize((b-a) × (c-a))
		ux, uy, uz := bx-ax, by-ay, bz-az
		vx, vy, vz := cx-ax, cy-ay, cz-az
		nx := uy*vz - uz*vy
		ny := uz*vx - ux*vz
		nz := ux*vy - uy*vx
		l := float32(math.Sqrt(float64(nx*nx + ny*ny + nz*nz)))
		if l > 0 {
			nx, ny, nz = nx/l, ny/l, nz/l
		}
		if err := binary.Write(f, binary.LittleEndian, [3]float32{nx, ny, nz}); err != nil {
			return err
		}
		if err := binary.Write(f, binary.LittleEndian, [3]float32{ax, ay, az}); err != nil {
			return err
		}
		if err := binary.Write(f, binary.LittleEndian, [3]float32{bx, by, bz}); err != nil {
			return err
		}
		if err := binary.Write(f, binary.LittleEndian, [3]float32{cx, cy, cz}); err != nil {
			return err
		}
		if err := binary.Write(f, binary.LittleEndian, uint16(0)); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
